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
	"regexp"
	"strings"

	apiErrors "k8s.io/apimachinery/pkg/api/errors"
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

var (
	supportedNodeAttestors = []string{"k8s_psat", "k8s_sat", "join_token"}
)

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
	logger := log.Log.WithValues("SpireServer", req.NamespacedName)

	server := &spirev1.SpireServer{}

	// fetching SPIRE Server instance
	err := r.Get(ctx, req.NamespacedName, server)
	if err != nil {
		if apiErrors.IsNotFound(err) {
			logger.Error(err, "SPIRE server not found.")
			return ctrl.Result{}, nil
		}

		logger.Error(err, "Failed to get SPIRE Server instance.")
		return ctrl.Result{}, nil
	}

	err = validateYaml(server)
	if err != nil {
		logger.Error(err, "Failed to validate YAML file so cannot deploy SPIRE server. Deleting old instance of CRD.")
		err = r.Delete(ctx, server)
		return ctrl.Result{}, err
	}

	serviceAccount := r.createServiceAccount(server, req.Namespace)
	err := r.Create(ctx, serviceAccount)
	if err != nil {
		logger.Error(err, "Failed to create", "Namespace", serviceAccount.Namespace, "Name", serviceAccount.Name)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func validateYaml(s *spirev1.SpireServer) error {
	// trust domain takes the same form as a DNS Name
	validDns, err := regexp.MatchString("^([a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9].)+[A-Za-z]{2,}$", s.Spec.TrustDomain)
	if err != nil {
		return errors.New("cannot validate DNS name for trust domain")
	} else if !validDns {
		return errors.New("trust domain is not a valid DNS name")
	}

	if !(s.Spec.Port >= 0 && s.Spec.Port <= 65535) {
		return errors.New("invalid port number") //TODO: should we restrict to other ports? This is basic for all ports.
	}

	var match bool
	for _, currAttestor := range s.Spec.NodeAttestors {
		match = false
		for _, nodeAttestor := range supportedNodeAttestors {
			if strings.Compare(currAttestor, nodeAttestor) == 0 {
				match = true
				break
			}
		}
	}

	if !match {
		return errors.New("incorrect node attestors list inputted: at least one of the specified node attestors is not supported")
	}

	if !((strings.Compare("disk", strings.ToLower(s.Spec.KeyStorage)) == 0) || (strings.Compare("memory", strings.ToLower(s.Spec.KeyStorage)) == 0)) {
		return errors.New("generated key storage is only supported on disk or in memory")
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpireServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&spirev1.SpireServer{}).
		Complete(r)
}

// CreateServiceAccount creates a service account for the SPIRE server.
func (r *SpireServerReconciler) createServiceAccount(m *spirev1.SpireServer, namespace string) *corev1.ServiceAccount {
	serviceAccount := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "spire-service-account",
			Namespace: namespace,
		},
	}
	return serviceAccount
}
