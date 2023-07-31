package controller

import (
	"context"
	"testing"

	spirev1 "github.com/glcp/spire-k8s-operator/api/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	serverTypeMeta = metav1.TypeMeta{
		APIVersion: "spire.hpe.com/v1",
		Kind:       "SpireServer",
	}
	serverObjectMeta = metav1.ObjectMeta{
		Name:      "valid-spire-server",
		Namespace: "default",
	}

	s = &spirev1.SpireServer{
		TypeMeta:   serverTypeMeta,
		ObjectMeta: serverObjectMeta,
		Spec: spirev1.SpireServerSpec{
			TrustDomain:   "example.org",
			Port:          8081,
			NodeAttestors: []string{"k8s_sat"},
			KeyStorage:    "disk",
			Replicas:      1,
		},
	}
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

func createReconciler() *SpireServerReconciler {
	reconciler := &SpireServerReconciler{
		Client: &MockClient{
			CreateFn: func(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
				// Handle the create logic in the mock client
				return nil
			},
		},
		Scheme: scheme.Scheme,
	}
	return reconciler
}

func TestSpireserverController(t *testing.T) {
	reconciler := createReconciler()

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

func TestValidNameSpaceConfigMap(t *testing.T) {
	reconcilerForConfigMap := createReconciler()
	configMap := reconcilerForConfigMap.spireConfigMapDeployment(s, "default")
	assert.Equal(t, configMap.Namespace, "default", "Namespaces should be the same.")
}

func TestInvalidNameSpaceServiceAccount(t *testing.T) {
	reconcilerForConfigMap := createReconciler()
	configMap := reconcilerForConfigMap.spireConfigMapDeployment(s, "namespace1")
	assert.NotEqual(t, configMap.Namespace, "namespace2", "Namespaces should not be the same.")
}

func TestEmptyNameSpaceServiceAccount(t *testing.T) {
	reconcilerForConfigMap := createReconciler()
	configMap := reconcilerForConfigMap.spireConfigMapDeployment(s, "")
	assert.Equal(t, configMap.Namespace, "", "Namespace should be empty.")
}

func TestValidConfigMapSingleAttestor(t *testing.T) {
	reconcilerForConfigMap := createReconciler()
	configMap := reconcilerForConfigMap.spireConfigMapDeployment(s, "default")

	assert.Contains(t, configMap.Data["server.conf"], "NodeAttestor \"k8s_sat\"")

	assert.Contains(t, configMap.Data["server.conf"], "trust_domain = \"example.org\"")
	assert.Contains(t, configMap.Data["server.conf"], "bind_port = \"8081\"")
	assert.Contains(t, configMap.Data["server.conf"], "KeyManager \"disk\"")
}

func TestValidConfigMapMultipleAttestors(t *testing.T) {
	s2 := &spirev1.SpireServer{
		TypeMeta:   serverTypeMeta,
		ObjectMeta: serverObjectMeta,
		Spec: spirev1.SpireServerSpec{
			TrustDomain:   "example.org",
			Port:          8081,
			NodeAttestors: []string{"k8s_sat", "join_token", "k8s_psat"},
			KeyStorage:    "disk",
			Replicas:      1,
		}}

	reconcilerForConfigMap := createReconciler()
	configMap := reconcilerForConfigMap.spireConfigMapDeployment(s2, "default")

	assert.Contains(t, configMap.Data["server.conf"], "NodeAttestor \"k8s_sat\"")
	assert.Contains(t, configMap.Data["server.conf"], "NodeAttestor \"join_token\"")
	assert.Contains(t, configMap.Data["server.conf"], "NodeAttestor \"k8s_psat\"")

	assert.Contains(t, configMap.Data["server.conf"], "trust_domain = \"example.org\"")
	assert.Contains(t, configMap.Data["server.conf"], "bind_port = \"8081\"")
	assert.Contains(t, configMap.Data["server.conf"], "KeyManager \"disk\"")
}
