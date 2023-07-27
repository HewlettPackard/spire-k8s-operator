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

	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	rbacv1 "k8s.io/api/rbac/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

	components := map[string]interface{}{
		"serviceAccount": serviceAccount,
    "clusterRole":        clusterRole,
		"clusterRoleBinding": clusterRoleBinding,
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

	match := false
	for _, nodeAttestor := range serverNodeAttestors {
		if strings.Compare(a.Spec.NodeAttestor, nodeAttestor) == 0 {
			match = true
			break
		}
	}

	if !match {
		return errors.New("the inputted node attestor is not supported by the server")
	}

	for _, currWLAttestor := range a.Spec.WorkloadAttestors {
		match = false
		for _, wLAttestor := range supportedWorkloadAttestors {
			if strings.Compare(currWLAttestor, wLAttestor) == 0 {
				match = true
				break
			}
		}

		if !match {
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

// SetupWithManager sets up the controller with the Manager.
func (r *SpireAgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&spirev1.SpireAgent{}).
		Complete(r)
}
