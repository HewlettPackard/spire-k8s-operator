package controller

import (
	"context"
	"testing"

	// "sigs.k8s.io/controller-runtime/pkg/client/fake"

	spirev1 "github.com/glcp/spire-k8s-operator/api/v1"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var reconciler = &SpireServerReconciler{
	Client: &MockClient{
		CreateFn: func(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
			// Handle the create logic in the mock client
			return nil
		},
	},
	Scheme: scheme.Scheme,
}

var (
	serverTypeMeta = metav1.TypeMeta{
		APIVersion: "spire.hpe.com/v1",
		Kind:       "SpireServer",
	}

	serverObjectMeta = metav1.ObjectMeta{
		Name:      "valid-spire-server",
		Namespace: "default",
	}

	mockSpireServer = createSpireServer("example.org", 8081, []spirev1.NodeAttestor{{Name: "k8s_sat"}}, "disk", 1)
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

func createSpireServer(trustDomain string, port int, nodeAttestors []spirev1.NodeAttestor, keyStorage string, replicas int) *spirev1.SpireServer {
	return &spirev1.SpireServer{
		TypeMeta:   serverTypeMeta,
		ObjectMeta: serverObjectMeta,
		Spec: spirev1.SpireServerSpec{
			TrustDomain:   trustDomain,
			Port:          port,
			NodeAttestors: nodeAttestors,
			KeyStorage:    keyStorage,
			Replicas:      replicas,
		},
	}
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

func TestValidDNSStringTrustDomain(t *testing.T) {
	_, dnsValue := dns.IsDomainName("prod.acme.com")
	validTrustDomain := checkTrustDomain("prod.acme.com")
	if validTrustDomain == nil && dnsValue {
		trustDomainValid := true
		assert.Equal(t, trustDomainValid, dnsValue, "There should be no error with \"prod.acme.com\"")
	}
}

func TestInvalidDNSStringTrustDomain(t *testing.T) {
	_, dnsValue := dns.IsDomainName("prod@acme.com")
	validTrustDomain := checkTrustDomain("prod@acme.com")
	assert.Equal(t, dnsValue, true, "Even though it is an invalid DNS, function only checks for 8 bit protocol and character limit.")
	assert.NotEqual(t, validTrustDomain, nil, "There should be an error with \"prod@acme.com\"")
}

func TestValidDNSNumberTrustDomain(t *testing.T) {
	_, dnsValue := dns.IsDomainName("8-8-8-8")
	validTrustDomain := checkTrustDomain("8-8-8-8")
	if validTrustDomain == nil && dnsValue {
		trustDomainValid := true
		assert.Equal(t, trustDomainValid, dnsValue, "There should be no error with \"8-8-8-8\"")
	}
}

func TestInvalidDNSNumberTrustDomain(t *testing.T) {
	_, dnsValue := dns.IsDomainName("8*8*8*8")
	validTrustDomain := checkTrustDomain("8*8*8*8")
	assert.Equal(t, dnsValue, true, "Even though it is an invalid DNS, function only checks for 8 bit protocol and character limit.")
	assert.NotEqual(t, validTrustDomain, nil, "There should be an error with \"8*8*8*8\"")
}

func TestValidStringTrustDomain(t *testing.T) {
	_, dnsValue := dns.IsDomainName("thisisatrustdomain")
	validTrustDomain := checkTrustDomain("thisisatrustdomain")
	if validTrustDomain == nil && dnsValue {
		trustDomainValid := true
		assert.Equal(t, trustDomainValid, dnsValue, "There should be no error with \"thisisatrustdomain\"")
	}
}

func TestInvalidStringTrustDomain(t *testing.T) {
	_, dnsValue := dns.IsDomainName("this is an invalid trust domain")
	invalidTrustDomain := checkTrustDomain("this is an invalid trust domain")
	assert.Equal(t, dnsValue, true, "Even though it is an invalid DNS, function only checks for 8 bit protocol and character limit.")
	assert.NotEqual(t, invalidTrustDomain, nil, "There should be an error with \"this is an invalid trust domain\"")
}
func TestValidNumberTrustDomain(t *testing.T) {
	_, dnsValue := dns.IsDomainName("393939")
	validTrustDomain := checkTrustDomain("393939")
	if validTrustDomain == nil && dnsValue {
		trustDomainValid := true
		assert.Equal(t, trustDomainValid, dnsValue, "There should be no error with \"393939\"")
	}
}

func TestInvalidNumberTrustDomain(t *testing.T) {
	_, dnsValue := dns.IsDomainName("*10001")
	invalidTrustDomain := checkTrustDomain("*10001")
	assert.Equal(t, dnsValue, true, "Even though it is an invalid DNS, function only checks for 8 bit protocol and character limit.")
	assert.NotEqual(t, invalidTrustDomain, nil, "There should be an error with \"*10001\"")
}

func TestValidNameSpaceConfigMap(t *testing.T) {
	configMap := reconciler.spireConfigMapDeployment(mockSpireServer, "default")
	assert.Equal(t, configMap.Namespace, "default", "Namespaces should be the same.")
}

func TestInvalidNameSpaceConfigMap(t *testing.T) {
	configMap := reconciler.spireConfigMapDeployment(mockSpireServer, "namespace1")
	assert.NotEqual(t, configMap.Namespace, "namespace2", "Namespaces should not be the same.")
}

func TestEmptyNameSpaceConfigMap(t *testing.T) {
	configMap := reconciler.spireConfigMapDeployment(mockSpireServer, "")
	assert.Equal(t, configMap.Namespace, "", "Namespace should be empty.")
}

func TestValidConfigMapSingleAttestor(t *testing.T) {
	configMap := reconciler.spireConfigMapDeployment(mockSpireServer, "default")

	assert.Contains(t, configMap.Data["server.conf"], "NodeAttestor \"k8s_sat\"")

	assert.Contains(t, configMap.Data["server.conf"], "trust_domain = \"example.org\"")
	assert.Contains(t, configMap.Data["server.conf"], "bind_port = \"8081\"")
	assert.Contains(t, configMap.Data["server.conf"], "KeyManager \"disk\"")
}

func TestValidConfigMapMultipleAttestors(t *testing.T) {
	mockSpireServer2 := createSpireServer("example.org", 8081, []spirev1.NodeAttestor{{Name: "k8s_sat"}, {Name: "join_token"}, {Name: "k8s_psat"}}, "disk", 1)

	configMap := reconciler.spireConfigMapDeployment(mockSpireServer2, "default")

	assert.Contains(t, configMap.Data["server.conf"], "NodeAttestor \"k8s_sat\"")
	assert.Contains(t, configMap.Data["server.conf"], "NodeAttestor \"join_token\"")
	assert.Contains(t, configMap.Data["server.conf"], "NodeAttestor \"k8s_psat\"")

	assert.Contains(t, configMap.Data["server.conf"], "trust_domain = \"example.org\"")
	assert.Contains(t, configMap.Data["server.conf"], "bind_port = \"8081\"")
	assert.Contains(t, configMap.Data["server.conf"], "KeyManager \"disk\"")
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
	serviceAccount := reconciler.createServiceAccount("")
	assert.Equal(t, serviceAccount.Namespace, "", "Namespaces should be empty.")
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

func TestValidNameSpaceRoles(t *testing.T) {
	roles := reconciler.spireRoleDeployment("default")
	assert.Equal(t, roles.Namespace, "default")
	assert.Equal(t, roles.Kind, "Role")
	assert.Equal(t, roles.Name, "spire-server-configmap-role")
	assert.Equal(t, roles.APIVersion, "rbac.authorization.k8s.io/v1")
	assert.Equal(t, roles.Rules[0].Verbs, []string{"patch", "get", "list"})
	assert.Equal(t, roles.Rules[0].Resources, []string{"configmaps"})
	assert.Equal(t, roles.Rules[0].APIGroups, []string{""})
}

func TestInvalidNameSpaceRoles(t *testing.T) {
	roles := reconciler.spireRoleDeployment("default1")
	assert.NotEqual(t, roles.Namespace, "default2")
}

func TestEmptyNameSpaceRoles(t *testing.T) {
	roles := reconciler.spireRoleDeployment("")
	assert.Equal(t, roles.Namespace, "")
}

func TestValidNameSpaceRoleBinding(t *testing.T) {
	roleBinding := reconciler.spireRoleBindingDeployment("default")
	assert.Equal(t, roleBinding.Namespace, "default")
	assert.Equal(t, roleBinding.Kind, "RoleBinding")
	assert.Equal(t, roleBinding.APIVersion, "rbac.authorization.k8s.io/v1")
	assert.Equal(t, roleBinding.Name, "spire-server-configmap-role-binding")
	assert.Equal(t, roleBinding.RoleRef.Kind, "Role")
	assert.Equal(t, roleBinding.RoleRef.Name, "spire-server-configmap-role")
	assert.Equal(t, roleBinding.RoleRef.APIGroup, "rbac.authorization.k8s.io")
	assert.Equal(t, roleBinding.Subjects[0].Kind, "ServiceAccount")
	assert.Equal(t, roleBinding.Subjects[0].Name, "spire-server")
	assert.Equal(t, roleBinding.Subjects[0].Namespace, "default")
}

func TestInvalidNameSpaceRoleBinding(t *testing.T) {
	roleBinding := reconciler.spireRoleBindingDeployment("default1")
	assert.NotEqual(t, roleBinding.Namespace, "default2")
}

func TestEmptyNameSpaceRoleBinding(t *testing.T) {
	roleBinding := reconciler.spireRoleBindingDeployment("")
	assert.Equal(t, roleBinding.Namespace, "")
}

func TestValidNameSpaceClusterRoles(t *testing.T) {
	clusterRoles := reconciler.spireClusterRoleDeployment("default")
	assert.Equal(t, clusterRoles.Namespace, "")
	assert.Equal(t, clusterRoles.Kind, "ClusterRole")
	assert.Equal(t, clusterRoles.Name, "spire-server-trust-role")
	assert.Equal(t, clusterRoles.APIVersion, "rbac.authorization.k8s.io/v1")
	assert.Equal(t, clusterRoles.Rules[0].Verbs, []string{"create"})
	assert.Equal(t, clusterRoles.Rules[0].Resources, []string{"tokenreviews"})
	assert.Equal(t, clusterRoles.Rules[0].APIGroups, []string{"authentication.k8s.io"})
}

func TestInvalidNameSpaceClusterRoles(t *testing.T) {
	clusterRoles := reconciler.spireClusterRoleDeployment("default1")
	assert.Equal(t, clusterRoles.Namespace, "")
}

func TestEmptyNameSpaceClusterRoles(t *testing.T) {
	clusterRoles := reconciler.spireClusterRoleDeployment("")
	assert.Equal(t, clusterRoles.Namespace, "")
}

func TestValidNameSpaceClusterRoleBinding(t *testing.T) {
	clusterRoleBinding := reconciler.spireClusterRoleBindingDeployment("default")
	assert.Equal(t, clusterRoleBinding.Kind, "ClusterRoleBinding")
	assert.Equal(t, clusterRoleBinding.APIVersion, "rbac.authorization.k8s.io/v1")
	assert.Equal(t, clusterRoleBinding.Name, "spire-server-trust-role-binding")
	assert.Equal(t, clusterRoleBinding.RoleRef.Kind, "ClusterRole")
	assert.Equal(t, clusterRoleBinding.RoleRef.Name, "spire-server-trust-role")
	assert.Equal(t, clusterRoleBinding.RoleRef.APIGroup, "rbac.authorization.k8s.io")
	assert.Equal(t, clusterRoleBinding.Subjects[0].Kind, "ServiceAccount")
	assert.Equal(t, clusterRoleBinding.Subjects[0].Name, "spire-server")
	assert.Equal(t, clusterRoleBinding.Subjects[0].Namespace, "default")
}

func TestInvalidNameSpaceClusterRoleBinding(t *testing.T) {
	clusterRoleBinding := reconciler.spireClusterRoleBindingDeployment("default1")
	assert.Equal(t, clusterRoleBinding.Namespace, "")
}

func TestEmptyNameSpaceClusterRoleBinding(t *testing.T) {
	clusterRoleBinding := reconciler.spireClusterRoleBindingDeployment("")
	assert.Equal(t, clusterRoleBinding.Namespace, "")
}
