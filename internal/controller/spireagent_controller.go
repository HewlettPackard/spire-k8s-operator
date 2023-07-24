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

	return ctrl.Result{}, nil
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
		if strings.Compare(nodeAttestor, nodeAttestor) == 0 {
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
	}

	if !match {
		return errors.New("incorrect workload attestors list inputted: at least one of the specified workload attestors is not supported")
	}

	if !((strings.Compare("disk", strings.ToLower(a.Spec.KeyStorage)) == 0) || (strings.Compare("memory", strings.ToLower(a.Spec.KeyStorage)) == 0)) {
		return errors.New("generated key storage is only supported on disk or in memory")
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpireAgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&spirev1.SpireAgent{}).
		Complete(r)
}
