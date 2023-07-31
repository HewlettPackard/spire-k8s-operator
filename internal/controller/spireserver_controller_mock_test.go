package controller

import (
	"context"
	"testing"

	// "sigs.k8s.io/controller-runtime/pkg/client/fake"

	spirev1 "github.com/glcp/spire-k8s-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MockClient struct {
	client.Client
	CreateFn func(context.Context, client.Object, ...client.CreateOption) error
}

func (m *MockClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, obj, opts...)
	}
	return nil
}

var reconciler = &SpireServerReconciler{
	Client: &MockClient{
		CreateFn: func(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
			// Handle the create logic in the mock client
			return nil
		},
	},
	Scheme: scheme.Scheme,
}

func TestSpireserverController(t *testing.T) {

	spireserver := &spirev1.SpireServer{}
	spireServiceNamespace := "test-namespace"
	spireServiceAccount := reconciler.createServiceAccount(spireServiceNamespace)
	bundle := reconciler.spireBundleDeployment(spireServiceNamespace)
	roles := reconciler.spireRoleDeployment(spireServiceNamespace)
	roleBinding := reconciler.spireRoleBindingDeployment(spireServiceNamespace)
	clusterRoles := reconciler.spireClusterRoleDeployment(spireServiceNamespace)
	clusterRoleBinding := reconciler.spireClusterRoleBindingDeployment(spireServiceNamespace)
	serverConfigMap := reconciler.spireConfigMapDeployment(spireserver, spireServiceNamespace)
	spireStatefulSet := reconciler.spireStatefulSetDeployment(2, spireServiceNamespace)
	spireService := reconciler.spireServiceDeployment(8081, spireServiceNamespace)

	// Call the method you want to test
	// Assert the expected behavior

	if spireServiceAccount.Namespace != spireServiceNamespace {
		t.Errorf("Expected namespace %s, got %s", spireServiceNamespace, spireServiceAccount.Namespace)
	}
	if bundle.Namespace != spireServiceNamespace {
		t.Errorf("Expected namespace %s, got %s", spireServiceNamespace, bundle.Namespace)
	}
	if roles.Namespace != spireServiceNamespace {
		t.Errorf("Expected namespace %s, got %s", spireServiceNamespace, roles.Namespace)
	}
	if roleBinding.Namespace != spireServiceNamespace {
		t.Errorf("Expected namespace %s, got %s", spireServiceNamespace, roleBinding.Namespace)
	}
	if clusterRoles.Namespace != "" {
		t.Errorf("Expected namespace %s, got %s", spireServiceNamespace, clusterRoles.Namespace)
	}
	if clusterRoleBinding.Namespace != "" {
		t.Errorf("Expected namespace \"\", got %s", clusterRoleBinding.Namespace)
	}
	if serverConfigMap.Namespace != spireServiceNamespace {
		t.Errorf("Expected namespace %s, got %s", spireServiceNamespace, serverConfigMap.Namespace)
	}
	if spireStatefulSet.Namespace != spireServiceNamespace {
		t.Errorf("Expected namespace %s, got %s", spireServiceNamespace, spireStatefulSet.Namespace)
	}
	if spireService.Namespace != spireServiceNamespace {
		t.Errorf("Expected namespace %s, got %s", spireServiceNamespace, spireService.Namespace)
	}
}
func TestValidNameSpaceServiceAccount(t *testing.T) {
	spireServiceNamespace := "sameNameSpace"
	serviceAccount := reconciler.createServiceAccount(spireServiceNamespace)
	assert.Equal(t, serviceAccount.Namespace, spireServiceNamespace, "Namespaces should be the same.")
}

func TestInvalidNameSpaceServiceAccount(t *testing.T) {
	spireServiceNamespace := "namespace1"
	serviceAccount := reconciler.createServiceAccount("namespace2")
	assert.NotEqual(t, serviceAccount.Namespace, spireServiceNamespace, "Namespaces should not be the same.")
}

func TestEmptyNameSpaceServiceAccount(t *testing.T) {
	serviceAccount := reconciler.spireBundleDeployment("")
	assert.Equal(t, serviceAccount.Namespace, "", "Namespaces should be empty.")
}
