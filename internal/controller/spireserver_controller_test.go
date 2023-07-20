package controller

import (
	"context"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	spirev1 "github.com/glcp/spire-k8s-operator/api/v1"
)

var _ = Describe("SpireServer controller", func() {
	const (
		duration = time.Second * 10
		interval = time.Millisecond * 250
		timeout  = time.Second * 10
	)
	Context("When installing SPIRE server", func() {
		It("Should create SPIRE server Trust Bundle", func() {
			By("By creating SPIRE server Trust Bundle with static config")
			ctx := context.Background()
			testBundle := &corev1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "spire-bundle",
					Namespace: "default",
				},
			}
			Expect(k8sClient.Create(ctx, testBundle)).Should(Succeed())

			bundleLookupKey := types.NamespacedName{Name: "spire-bundle", Namespace: "default"}
			createdBundle := &corev1.ConfigMap{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, bundleLookupKey, createdBundle)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Expect(createdBundle.ObjectMeta.Name).Should(Equal("spire-bundle"))
			Expect(createdBundle.ObjectMeta.Namespace).Should(Equal("default"))
			Expect(createdBundle.Data).Should(Not(Equal(nil)))
			Expect(createdBundle.BinaryData).Should(Not(Equal(nil)))
			Expect(createdBundle.Labels).Should(Not(Equal(nil)))
			Expect(createdBundle.Annotations).Should(Not(Equal(nil)))
		})

		It("Should create SPIRE server service", func() {
			By("By creating SPIRE server service with static config")
			ctx := context.Background()

			serviceSpec := corev1.ServiceSpec{
				Ports: []corev1.ServicePort{{Port: int32(8081)}},
			}
			spireService := &corev1.Service{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Service",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "spire-service",
					Namespace: "default",
				},
				Spec: serviceSpec,
			}
			Expect(k8sClient.Create(ctx, spireService)).Should(Succeed())

			/*
				After creating this Service, let's check that the Spire service's Spec fields match what we passed in or not.
				Note that, because the k8s apiserver may not have finished creating a Service after our `Create()` call from earlier,
				we will use Gomega’s Eventually() testing function instead of Expect() to give the apiserver an opportunity to finish
				creating our Spire Service.
				`Eventually()` will repeatedly run the function provided as an argument every interval seconds until
				(a) the function’s output matches what’s expected in the subsequent `Should()` call, or
				(b) the number of attempts * interval period exceed the provided timeout value.
				In the examples below, timeout and interval are Go Duration values of our choosing.
			*/
			serviceLookupKey := types.NamespacedName{Name: "spire-service", Namespace: "default"}
			createdService := &corev1.Service{}

			// We'll need to retry getting this newly created Service, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, serviceLookupKey, createdService)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Now let us see if the expectation matches or not
			Expect(createdService.Spec.Ports[0].Port).Should(Equal(int32(8081)))
		})

		It("Should create SPIRE server Role, ClusterRole, and Bindings", func() {
			By("By creating SPIRE server Role, ClusterRole, and Bindings with static config")
			ctx := context.Background()

			roleRules := rbacv1.PolicyRule{
				Verbs:     []string{"patch", "get", "list"},
				Resources: []string{"configmaps"},
				APIGroups: []string{""},
			}
			clusterRoleRules := rbacv1.PolicyRule{
				Verbs:     []string{"create"},
				Resources: []string{"tokenreviews"},
				APIGroups: []string{"authentication.k8s.io"},
			}

			testRole := &rbacv1.Role{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Role",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "spire-server-configmap-role",
					Namespace: "default",
				},
				Rules: []rbacv1.PolicyRule{
					roleRules,
				},
			}

			testClusterRole := &rbacv1.ClusterRole{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterRole",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "spire-server-trust-role",
				},
				Rules: []rbacv1.PolicyRule{
					clusterRoleRules,
				},
			}

			roleBindingSubject := rbacv1.Subject{
				Kind:      "ServiceAccount",
				Name:      "spire-server",
				Namespace: "default",
			}

			testRoleBinding := &rbacv1.RoleBinding{
				TypeMeta: metav1.TypeMeta{
					Kind:       "RoleBinding",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "spire-server-configmap-role-binding",
					Namespace: "default",
				},
				Subjects: []rbacv1.Subject{
					roleBindingSubject,
				},
				RoleRef: rbacv1.RoleRef{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     "Role",
					Name:     "spire-server-configmap-role",
				},
			}

			clusterRoleBindingSubject := rbacv1.Subject{
				Kind:      "ServiceAccount",
				Name:      "spire-server",
				Namespace: "default",
			}

			testClusterRoleBinding := &rbacv1.ClusterRoleBinding{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ClusterRoleBinding",
					APIVersion: "rbac.authorization.k8s.io/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "spire-server-trust-role-binding",
				},
				Subjects: []rbacv1.Subject{
					clusterRoleBindingSubject,
				},
				RoleRef: rbacv1.RoleRef{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     "ClusterRole",
					Name:     "spire-server-trust-role",
				},
			}

			Expect(k8sClient.Create(ctx, testRole)).Should(Succeed())
			Expect(k8sClient.Create(ctx, testClusterRole)).Should(Succeed())
			Expect(k8sClient.Create(ctx, testRoleBinding)).Should(Succeed())
			Expect(k8sClient.Create(ctx, testClusterRoleBinding)).Should(Succeed())

			roleLookupKey := types.NamespacedName{Name: "spire-server-configmap-role", Namespace: "default"}
			clusterRoleLookupKey := types.NamespacedName{Name: "spire-server-trust-role", Namespace: "default"}
			roleBindingLookupKey := types.NamespacedName{Name: "spire-server-configmap-role-binding", Namespace: "default"}
			clusterRoleBindingLookupKey := types.NamespacedName{Name: "spire-server-trust-role-binding", Namespace: "default"}
			createdRole := &rbacv1.Role{}
			createdClusterRole := &rbacv1.ClusterRole{}
			createdRoleBinding := &rbacv1.RoleBinding{}
			createdClusterRoleBinding := &rbacv1.ClusterRoleBinding{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, roleLookupKey, createdRole)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, clusterRoleLookupKey, createdClusterRole)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, roleBindingLookupKey, createdRoleBinding)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, clusterRoleBindingLookupKey, createdClusterRoleBinding)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Now let us see if the expectation matches or not
			Expect(createdRole.ObjectMeta.Name).Should(Equal("spire-server-configmap-role"))
			Expect(createdRole.ObjectMeta.Namespace).Should(Equal("default"))
			Expect(createdRole.Labels).Should(Not(Equal(nil)))
			Expect(createdRole.Annotations).Should(Not(Equal(nil)))
			Expect(len(createdRole.Rules)).Should(Not(Equal(0)))
			Expect(createdRole.Rules[0].Verbs).Should(ContainElement("patch"))
			Expect(createdRole.Rules[0].Verbs).Should(ContainElement("get"))
			Expect(createdRole.Rules[0].Verbs).Should(ContainElement("list"))
			Expect(createdRole.Rules[0].Resources).Should(ContainElement("configmaps"))
			Expect(len(createdRole.Rules[0].APIGroups)).Should(Equal(1))

			Expect(createdClusterRole.ObjectMeta.Name).Should(Equal("spire-server-trust-role"))
			Expect(createdClusterRole.Labels).Should(Not(Equal(nil)))
			Expect(createdClusterRole.Annotations).Should(Not(Equal(nil)))
			Expect(len(createdClusterRole.Rules)).Should(Not(Equal(0)))
			Expect(createdClusterRole.Rules[0].Verbs).Should(ContainElement("create"))
			Expect(createdClusterRole.Rules[0].Resources).Should(ContainElement("tokenreviews"))
			Expect(createdClusterRole.Rules[0].APIGroups).Should(ContainElement("authentication.k8s.io"))

			Expect(createdRoleBinding.ObjectMeta.Name).Should(Equal("spire-server-configmap-role-binding"))
			Expect(createdRoleBinding.ObjectMeta.Namespace).Should(Equal("default"))
			Expect(createdRoleBinding.Labels).Should(Not(Equal(nil)))
			Expect(createdRoleBinding.Annotations).Should(Not(Equal(nil)))
			Expect(createdRoleBinding.RoleRef.Kind).Should(Equal("Role"))
			Expect(createdRoleBinding.RoleRef.APIGroup).Should(Equal("rbac.authorization.k8s.io"))
			Expect(createdRoleBinding.RoleRef.Name).Should(Equal("spire-server-configmap-role"))
			Expect(createdRoleBinding.Subjects).Should(Not(Equal(nil)))
			Expect(createdRoleBinding.Subjects[0].Kind).Should(Equal("ServiceAccount"))
			Expect(createdRoleBinding.Subjects[0].Name).Should(Equal("spire-server"))
			Expect(createdRoleBinding.Subjects[0].Namespace).Should(Equal("default"))

			Expect(createdClusterRoleBinding.ObjectMeta.Name).Should(Equal("spire-server-trust-role-binding"))
			Expect(createdClusterRoleBinding.Labels).Should(Not(Equal(nil)))
			Expect(createdClusterRoleBinding.Annotations).Should(Not(Equal(nil)))
			Expect(createdClusterRoleBinding.RoleRef.Kind).Should(Equal("ClusterRole"))
			Expect(createdClusterRoleBinding.RoleRef.APIGroup).Should(Equal("rbac.authorization.k8s.io"))
			Expect(createdClusterRoleBinding.RoleRef.Name).Should(Equal("spire-server-trust-role"))
			Expect(createdClusterRoleBinding.Subjects).Should(Not(Equal(nil)))
			Expect(createdClusterRoleBinding.Subjects[0].Kind).Should(Equal("ServiceAccount"))
			Expect(createdClusterRoleBinding.Subjects[0].Name).Should(Equal("spire-server"))
			Expect(createdClusterRoleBinding.Subjects[0].Namespace).Should(Equal("default"))
		})

		It("Should create SPIRE server ConfigMap", func() {
			By("By creating SPIRE server ConfigMap with static config")
			ctx := context.Background()

			var nodeAttestors []string = []string{"k8s_sat"}
			var namespace string = "default"
			var port int32 = int32(8081)
			var trustDomain string = "example.org"
			var keyStorage string = "disk"

			nodeAttestorsConfig := ""

			for _, nodeAttestor := range nodeAttestors {
				if strings.Compare(nodeAttestor, "join_token") == 0 {
					nodeAttestorsConfig += `
		
			NodeAttestor "join_token" {
				plugin_data {
		
				}
			}`
				} else if strings.Compare(nodeAttestor, "k8s_sat") == 0 {
					nodeAttestorsConfig += `
		
			NodeAttestor "k8s_sat" {
				plugin_data {
					clusters = {
						"demo-cluster" = {
							use_token_review_api_validation = true
							service_account_allow_list = ["spire:spire-agent"]
						}
					}
				}
			}`
				} else if strings.Compare(nodeAttestor, "k8s_psat") == 0 {
					nodeAttestorsConfig += `
		
			NodeAttestor "k8s_psat" {
				plugin_data {
					clusters = {
						"cluster" = {
							service_account_allow_list = ["` + namespace + `:spire-agent"]
						}
					}
				}
			}`
				}
			}

			config := `
		server {
			bind_address = "0.0.0.0"
			bind_port = "` + strconv.Itoa(int(port)) + `"
			socket_path = "/tmp/spire-server/private/api.sock"
			trust_domain = "` + trustDomain + `"
			data_dir = "/run/spire/data"
			log_level = "DEBUG"
			ca_key_type = "rsa-2048"
		
			ca_subject = {
				country = ["US"],
				organization = ["SPIFFE"],
				common_name = "",
			}
		}
		
		plugins {
			DataStore "sql" {
				plugin_data {
				  database_type = "sqlite3"
				  connection_string = "/run/spire/data/datastore.sqlite3"
				}
			}` +
				nodeAttestorsConfig + `
		
			KeyManager "` + keyStorage + `" {
				plugin_data {
					keys_path = "/run/spire/data/keys.json"
				}
			}
		
			Notifier "k8sbundle" {
				plugin_data {
					namespace = "` + namespace + `"
				}
			}
		}
		
		health_checks {
			listener_enabled = true
			bind_address = "0.0.0.0"
			bind_port = "8080"
			live_path = "/live"
			ready_path = "/ready"
		}`

			configMap := &corev1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},

				ObjectMeta: metav1.ObjectMeta{
					Name:      "spire-config-map",
					Namespace: namespace,
				},

				Data: map[string]string{
					"server.conf": config,
				},
			}

			Expect(k8sClient.Create(ctx, configMap)).Should(Succeed())

			configMapLookupKey := types.NamespacedName{Name: "spire-config-map", Namespace: "default"}
			createdConfigMap := &corev1.ConfigMap{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, configMapLookupKey, createdConfigMap)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Expect(createdConfigMap.Name).Should(Equal("spire-config-map"))
			Expect(createdConfigMap.Namespace).Should(Equal("default"))
			Expect(createdConfigMap.Data).ShouldNot(BeEmpty())
			Expect("server.conf").Should(BeKeyOf(createdConfigMap.Data))
			Expect(len(createdConfigMap.Data["server.conf"])).ShouldNot(BeZero())
		})
	})
	Context("When creating SPIRE server with invalid/empty trust domain", func() {
		var ctx = context.Background()
		var spireServer *spirev1.SpireServer

		BeforeEach(func() {
			spireServer = &spirev1.SpireServer{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "spire.hpe.com/v1",
					Kind:       "SpireServer",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "invalid-spire-server",
					Namespace: "default",
				},
				Spec: spirev1.SpireServerSpec{
					TrustDomain:   "",
					Port:          8081,
					NodeAttestors: []string{"k8s_sat"},
					KeyStorage:    "disk",
					Replicas:      1,
				},
			}

			Expect(k8sClient.Create(ctx, spireServer)).Should(Succeed())
		})

		It("should delete the CRD instance created", func() {
			serverLookupKey := types.NamespacedName{Name: spireServer.Name, Namespace: spireServer.Namespace}
			createdSpireServer := &spirev1.SpireServer{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, serverLookupKey, createdSpireServer)
				return err != nil
			}, timeout, interval).Should(BeTrue())
		})
	})
})
