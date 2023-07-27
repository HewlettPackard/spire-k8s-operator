/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"errors"
	"strings"
	"golang.org/x/exp/slices"
	"strconv"
  
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	spirev1 "github.com/glcp/spire-k8s-operator/api/v1"
)

// SpireAgentReconciler reconciles a SpireAgent object
type SpireAgentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	supportedWorkloadAttestors = []string{"k8s", "unix", "docker", "systemd", "windows"}
)

//+kubebuilder:rbac:groups=spire.hpe.com,resources=spireagents,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=spire.hpe.com,resources=spireagents/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=spire.hpe.com,resources=spireagents/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SpireAgent object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *SpireAgentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	logger := log.Log.WithValues("SpireAgent", req.NamespacedName)

	agent := &spirev1.SpireAgent{}

	// fetching SPIRE Agent instance
	err := r.Get(ctx, req.NamespacedName, agent)
	if err != nil {
		if apiErrors.IsNotFound(err) {
			logger.Error(err, "SPIRE Agent not found.")
			return ctrl.Result{}, nil
		}

		logger.Error(err, "Failed to get SPIRE Agent instance.")
		return ctrl.Result{}, nil
	}

	err = validateAgentYaml(agent, r, ctx)
	if err != nil {
		logger.Error(err, "Failed to validate YAML file so cannot deploy SPIRE agent. Deleting old instance of CRD.")
		err = r.Delete(ctx, agent)
		return ctrl.Result{}, err
	}
  
	clusterRole := r.agentClusterRoleDeployment()
	clusterRoleBinding := r.agentClusterRoleBindingDeployment(req.Namespace)
	serviceAccount := r.agentServiceAccountDeployment(req.Namespace)
	agentConfigMap := r.agentConfigMapDeployment(agent, req.Namespace)
	agentDaemonSet := r.agentDaemonSetDeployment(agent, req.Namespace)

	components := map[string]interface{}{
		"serviceAccount":     serviceAccount,
		"clusterRole":        clusterRole,
		"clusterRoleBinding": clusterRoleBinding,
		"agentConfigMap":     agentConfigMap,
		"agentDaemonSet":     agentDaemonSet,
	}

	for key, value := range components {
		err := r.Create(ctx, value.(client.Object))
		result, createError := checkIfFailToCreate(err, key, logger)
		if createError != nil {
			err = createError
			return result, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *SpireAgentReconciler) agentClusterRoleDeployment() *rbacv1.ClusterRole {
	rules := rbacv1.PolicyRule{
		Verbs:     []string{"get"},
		Resources: []string{"pods", "nodes", "nodes/proxy"},
		APIGroups: []string{""},
	}

	clusterRole := &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "spire-agent-cluster-role",
		},
		Rules: []rbacv1.PolicyRule{
			rules,
		},
	}

	return clusterRole
}

func (r *SpireAgentReconciler) agentClusterRoleBindingDeployment(namespace string) *rbacv1.ClusterRoleBinding {
	subject := rbacv1.Subject{
		Kind:      "ServiceAccount",
		Name:      "spire-agent",
		Namespace: namespace,
	}

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "spire-agent-cluster-role-binding",
		},
		Subjects: []rbacv1.Subject{
			subject,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "spire-agent-cluster-role",
		},
	}

	return clusterRoleBinding
}

func validateAgentYaml(a *spirev1.SpireAgent, r *SpireAgentReconciler, ctx context.Context) error {
	invalidTrustDomain := false
	checkTrustDomain(a.Spec.TrustDomain, &invalidTrustDomain)

	if invalidTrustDomain {
		return errors.New("trust domain is invalid")
	}

	if a.Spec.ServerPort != serverPort {
		return errors.New("the inputted port does not correspond to a SPIRE server")
	}

	if !(slices.Contains(serverNodeAttestors, a.Spec.NodeAttestor)) {
		return errors.New("the inputted node attestor is not supported by the server")
	}

	for _, currWLAttestor := range a.Spec.WorkloadAttestors {
		if !(slices.Contains(supportedWorkloadAttestors, currWLAttestor)) {
			return errors.New("incorrect workload attestors list inputted: at least one of the specified workload attestors is not supported")
		}
	}

	if !((strings.Compare("disk", strings.ToLower(a.Spec.KeyStorage)) == 0) || (strings.Compare("memory", strings.ToLower(a.Spec.KeyStorage)) == 0)) {
		return errors.New("generated key storage is only supported on disk or in memory")
	}

	return nil
}

func (r *SpireAgentReconciler) agentServiceAccountDeployment(namespace string) *corev1.ServiceAccount {
	serviceAccount := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spire-agent",
			Namespace: namespace,
		},
	}
	return serviceAccount
}

func (r *SpireAgentReconciler) agentDaemonSetDeployment(a *spirev1.SpireAgent, namespace string) *appsv1.DaemonSet {
	initContainer := corev1.Container{
		Name:  "init",
		Image: "cgr.dev/chainguard/wait-for-it",
		Args:  []string{"-t", "30", "spire-service:8081"},
	}

	volMount1 := corev1.VolumeMount{
		Name:      "spire-config",
		MountPath: "/run/spire/config",
		ReadOnly:  true,
	}

	volMount2 := corev1.VolumeMount{
		Name:      "spire-bundle",
		MountPath: "/run/spire/bundle",
	}

	volMount3 := corev1.VolumeMount{
		Name:      "spire-agent-socket",
		MountPath: "/run/spire/sockets",
		ReadOnly:  false,
	}

	livenessProbe := corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{
			Path: "/live", Port: intstr.IntOrString{IntVal: 8080}}},
		FailureThreshold:    2,
		InitialDelaySeconds: 15,
		PeriodSeconds:       60,
		TimeoutSeconds:      3,
	}

	readinessProbe := corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{
			Path: "/ready", Port: intstr.IntOrString{IntVal: 8080}}},
		InitialDelaySeconds: 5,
		PeriodSeconds:       5,
	}

	container := corev1.Container{
		Name:           "spire-agent",
		Image:          "ghcr.io/spiffe/spire-agent:1.5.1",
		Args:           []string{"-config", "/run/spire/config/agent.conf"},
		VolumeMounts:   []corev1.VolumeMount{volMount1, volMount2, volMount3},
		LivenessProbe:  &livenessProbe,
		ReadinessProbe: &readinessProbe,
	}

	vol1 := corev1.Volume{
		Name: "spire-config",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: "spire-agent"},
			},
		},
	}

	vol2 := corev1.Volume{
		Name: "spire-bundle",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: "spire-bundle"},
			},
		},
	}

	var hostPathType corev1.HostPathType = "DirectoryOrCreate"

	vol3 := corev1.Volume{
		Name: "spire-agent-socket",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/run/spire/sockets",
				Type: &hostPathType,
			},
		},
	}

	agentPodSpec := corev1.PodSpec{
		HostPID:            true,
		HostNetwork:        true,
		DNSPolicy:          "ClusterFirstWithHostNet",
		ServiceAccountName: "spire-agent",
		InitContainers:     []corev1.Container{initContainer},
		Containers:         []corev1.Container{container},
		Volumes:            []corev1.Volume{vol1, vol2, vol3},
	}

	daemonSetSpec := appsv1.DaemonSetSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"app": "spire-agent"},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Labels:    map[string]string{"app": "spire-agent"},
			},
			Spec: agentPodSpec,
		},
	}

	agentDaemonSet := &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "DaemonSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spire-agent",
			Namespace: namespace,
			Labels:    map[string]string{"app": "spire-agent"},
		},
		Spec: daemonSetSpec,
	}

	return agentDaemonSet
}

func (r *SpireAgentReconciler) agentConfigMapDeployment(a *spirev1.SpireAgent, namespace string) *corev1.ConfigMap {
	nodeAttestorsConfig := ""

	if strings.Compare(a.Spec.NodeAttestor, "join_token") == 0 {
		nodeAttestorsConfig += joinTokenAgentNodeAttestor()
	} else if strings.Compare(a.Spec.NodeAttestor, "k8s_sat") == 0 {
		nodeAttestorsConfig += k8sSatAgentNodeAttestor()
	} else if strings.Compare(a.Spec.NodeAttestor, "k8s_psat") == 0 {
		nodeAttestorsConfig += k8sPsatAgentNodeAttestor()
	}

	workloadAttestorsConfig := ""
	for _, wLAttestor := range a.Spec.WorkloadAttestors {
		if strings.Compare(wLAttestor, "k8s") == 0 {
			workloadAttestorsConfig += k8sWLAttestor()
		} else if strings.Compare(wLAttestor, "unix") == 0 {
			workloadAttestorsConfig += unixWLAttestor()
		} else if strings.Compare(wLAttestor, "docker") == 0 {
			workloadAttestorsConfig += dockerWLAttestor()
		} else if strings.Compare(wLAttestor, "systemd") == 0 {
			workloadAttestorsConfig += systemdWLAttestor()
		} else if strings.Compare(wLAttestor, "windows") == 0 {
			workloadAttestorsConfig += windowsWLAttestor()
		}
	}

	config := agentCreation(strconv.Itoa(a.Spec.ServerPort), a.Spec.TrustDomain) +
		pluginsAgent(nodeAttestorsConfig, a.Spec.KeyStorage, workloadAttestorsConfig) +
		healthChecks()

	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      "spire-agent",
			Namespace: namespace,
		},

		Data: map[string]string{
			"agent.conf": config,
		},
	}

	return configMap
}

func joinTokenAgentNodeAttestor() string {
	return `
	NodeAttestor "join_token" {
		plugin_data {

		}
	}`
}

func k8sSatAgentNodeAttestor() string {
	return `
	NodeAttestor "k8s_sat" {
		plugin_data {
			cluster = "demo-cluster"
		}
	}`
}

func k8sPsatAgentNodeAttestor() string {
	return `
	NodeAttestor "k8s_psat" {
		plugin_data {
			clusters = {
				"cluster" = {
		
				}
			}
		}
	}`
}

func agentCreation(port string, trustDomain string) string {
	return `
	agent {
		data_dir = "/run/spire"
		log_level = "DEBUG"
		server_address = "spire-service"
		server_port = "` + port + `"
		socket_path = "/run/spire/sockets/agent.sock"
		trust_bundle_path = "/run/spire/bundle/bundle.crt"
		trust_domain = "` + trustDomain + `"
	  }`
}

func pluginsAgent(nodeAttestorsConfig string, keyStorage string, workloadAttestorsConfig string) string {
	return `

	plugins {
		` + nodeAttestorsConfig + `
	
		KeyManager "` + keyStorage + `" {
			plugin_data {
			}
		} ` +
		workloadAttestorsConfig + `
	}`
}

func k8sWLAttestor() string {
	return `

	WorkloadAttestor "k8s" {
		plugin_data {
			skip_kubelet_verification = true
		}
	}`
}

func unixWLAttestor() string {
	return `

	WorkloadAttestor "unix" {
		plugin_data {
		}
	}`
}

func dockerWLAttestor() string {
	return `

	WorkloadAttestor "docker" {
		plugin_data {
		}
	}`
}

func systemdWLAttestor() string {
	return `

	WorkloadAttestor "systemd" {
	}`
}

func windowsWLAttestor() string {
	return `

	WorkloadAttestor "windows" {
	}`
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpireAgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&spirev1.SpireAgent{}).
		Complete(r)
}
