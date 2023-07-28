package controller

import (
	"context"
	"testing"

	// "sigs.k8s.io/controller-runtime/pkg/client/fake"

	spirev1 "github.com/glcp/spire-k8s-operator/api/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestSpireAgentController(t *testing.T) {
	// Create the objects needed for the test
	reconciler := &SpireAgentReconciler{
		Client: &MockClient{
			CreateFn: func(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
				// Handle the create logic in the mock client
				return nil
			},
		},
		Scheme: scheme.Scheme,
	}

	spireagent := &spirev1.SpireAgent{}
	spireServiceNamespace := "test-namespace"
	spireServiceAccount := reconciler.agentServiceAccountDeployment(spireServiceNamespace)
	clusterRoles := reconciler.agentClusterRoleDeployment()
	clusterRoleBinding := reconciler.agentClusterRoleBindingDeployment(spireServiceNamespace)
	agentConfigMap := reconciler.agentConfigMapDeployment(spireagent, spireServiceNamespace)
	spireService := reconciler.agentDaemonSetDeployment(spireagent, spireServiceNamespace)

	// Call the method you want to test
	// Assert the expected behavior

	if spireServiceAccount.Namespace != spireServiceNamespace {
		t.Errorf("Expected namespace %s, got %s", spireServiceNamespace, spireServiceAccount.Namespace)
	}
	if clusterRoles.Namespace != "" {
		t.Errorf("Expected namespace %s, got %s", spireServiceNamespace, clusterRoles.Namespace)
	}
	if clusterRoleBinding.Namespace != "" {
		t.Errorf("Expected namespace \"\", got %s", clusterRoleBinding.Namespace)
	}
	if agentConfigMap.Namespace != spireServiceNamespace {
		t.Errorf("Expected namespace %s, got %s", spireServiceNamespace, agentConfigMap.Namespace)
	}
	if spireService.Namespace != spireServiceNamespace {
		t.Errorf("Expected namespace %s, got %s", spireServiceNamespace, spireService.Namespace)
	}
}
