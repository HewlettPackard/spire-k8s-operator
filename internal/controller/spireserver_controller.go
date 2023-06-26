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

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	spirev1 "github.com/glcp/spire-k8s-operator/api/v1"
)

// SpireServerReconciler reconciles a SpireServer object
type SpireServerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=spire.hpe.com,resources=spireservers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=spire.hpe.com,resources=spireservers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=spire.hpe.com,resources=spireservers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SpireServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *SpireServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	logger := log.Log.WithValues("spireServer", req.NamespacedName)

	// TODO(user): your logic here
	spireServer := &spirev1.SpireServer{}
	clusterRoles := r.spireRoleDeployment(spireServer, req.NamespacedName.String())

	err := r.Create(ctx, clusterRoles)
	if err != nil {
		logger.Error(err, "Failed to create", "Namespace", clusterRoles.Namespace, "Name", clusterRoles.Name)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *SpireServerReconciler) spireRoleDeployment(m *spirev1.SpireServer, namespace string) *rbacv1.Role {
	rules := rbacv1.PolicyRule{
		Verbs:     []string{"patch", "get", "list"},
		Resources: []string{"configmap"},
		APIGroups: []string{""},
	}

	clusterRole := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spire-server-configmap-role",
			Namespace: namespace,
		},
		Rules: []rbacv1.PolicyRule{
			rules,
		},
	}
	return clusterRole
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpireServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&spirev1.SpireServer{}).
		Complete(r)
}
