package controller

import (
	"context"
	"testing"

	// "sigs.k8s.io/controller-runtime/pkg/client/fake"

	spirev1 "github.com/glcp/spire-k8s-operator/api/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var agentReconciler = &SpireAgentReconciler{
	Client: &MockClient{
		CreateFn: func(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
			// Handle the create logic in the mock client
			return nil
		},
	},
	Scheme: scheme.Scheme,
}

func TestSpireAgentController(t *testing.T) {
	// Create the objects needed for the test

	spireagent := &spirev1.SpireAgent{}
	spireServiceNamespace := "test-namespace"
	agentServiceAccount := agentReconciler.agentServiceAccountDeployment(spireServiceNamespace)
	agentClusterRoles := agentReconciler.agentClusterRoleDeployment()
	agentClusterRoleBinding := agentReconciler.agentClusterRoleBindingDeployment(spireServiceNamespace)
	agentConfigMap := agentReconciler.agentConfigMapDeployment(spireagent, spireServiceNamespace)
	agentDaemonSet := agentReconciler.agentDaemonSetDeployment(spireagent, spireServiceNamespace)

	// Call the method you want to test
	// Assert the expected behavior

	if agentServiceAccount.Namespace != spireServiceNamespace {
		t.Errorf("Expected namespace %s, got %s", spireServiceNamespace, agentServiceAccount.Namespace)
	}
	if agentClusterRoles.Namespace != "" {
		t.Errorf("Expected namespace %s, got %s", spireServiceNamespace, agentClusterRoles.Namespace)
	}
	if agentClusterRoleBinding.Namespace != "" {
		t.Errorf("Expected namespace \"\", got %s", agentClusterRoleBinding.Namespace)
	}
	if agentConfigMap.Namespace != spireServiceNamespace {
		t.Errorf("Expected namespace %s, got %s", spireServiceNamespace, agentConfigMap.Namespace)
	}
	if agentDaemonSet.Namespace != spireServiceNamespace {
		t.Errorf("Expected namespace %s, got %s", spireServiceNamespace, agentDaemonSet.Namespace)
	}
}
