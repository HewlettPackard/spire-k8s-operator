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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

	agentDaemonSet := r.agentDaemonSetDeployment(agent, req.Namespace)
	components := map[string]interface{}{
		"agentDaemonSet": agentDaemonSet,
	}
	serviceAccount := r.agentServiceAccountDeployment(req.Namespace)

	components := map[string]interface{}{
		"serviceAccount": serviceAccount,
	}

	for key, value := range components {
		err := r.Create(ctx, value.(client.Object))
		result, createError := checkIfFailToCreate(err, key, logger)
		if createError != nil {
			err = createError
			return result, err
		}
	}
	for key, value := range components {
		err = r.Create(ctx, value.(client.Object))
		result, createError := checkIfFailToCreate(err, key, logger)
		if createError != nil {
			err = createError
			return result, err
		}
	}
	return ctrl.Result{}, nil
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

// SetupWithManager sets up the controller with the Manager.
func (r *SpireAgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&spirev1.SpireAgent{}).
		Complete(r)
}
