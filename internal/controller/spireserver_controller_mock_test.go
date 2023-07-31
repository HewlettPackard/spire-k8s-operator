package controller

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	// "sigs.k8s.io/controller-runtime/pkg/client/fake"

	spirev1 "github.com/glcp/spire-k8s-operator/api/v1"
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
			return nil
		},
	},
	Scheme: scheme.Scheme,
}

func TestSpireserverController(t *testing.T) {
	// Create the objects needed for the test

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

func TestValidTrustBundle(t *testing.T) {
	spireServiceNamespace := "sameNameSpace"
	bundle := reconciler.spireBundleDeployment(spireServiceNamespace)
	assert.Equal(t, bundle.Namespace, spireServiceNamespace, "Namespaces should be the same.")
}

func TestInvalidNameSpaceTrustBundle(t *testing.T) {
	spireServiceNamespace := "namespace1"
	bundle := reconciler.spireBundleDeployment("namespace2")
	assert.NotEqual(t, bundle.Namespace, spireServiceNamespace, "Namespaces should not be the same.")
}

func TestEmptyNameSpaceTrustBundle(t *testing.T) {
	bundle := reconciler.spireBundleDeployment("")
	assert.Equal(t, bundle.Namespace, "", "Namespaces should be empty.")
}
