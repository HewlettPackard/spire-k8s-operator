package controller

import (
	"context"
	"strconv"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

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
		It("Should create SPIRE server Service Account", func() {
			By("By creating SPIRE server Service Account with static config")
			ctx := context.Background()
			testServiceAccount := &corev1.ServiceAccount{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ServiceAccount",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "spire-server",
					Namespace: "default",
				},
			}
			Expect(k8sClient.Create(ctx, testServiceAccount)).Should(Succeed())

			serviceAccountLookupKey := types.NamespacedName{Name: "spire-server", Namespace: "default"}
			createdServiceAccount := &corev1.ServiceAccount{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, serviceAccountLookupKey, createdServiceAccount)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Expect(createdServiceAccount.ObjectMeta.Name).Should(Equal("spire-server"))
			Expect(createdServiceAccount.ObjectMeta.Namespace).Should(Equal("default"))
		})
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

		It("Should create SPIRE server StatefulSet", func() {
			By("By creating SPIRE server StatefulSet with static config")
			ctx := context.Background()

			var numReplicas int32 = int32(2)
			labelSelector := metav1.LabelSelector{MatchLabels: map[string]string{"app": "spire-server"}}
			volMount1 := corev1.VolumeMount{
				Name:      "spire-config",
				MountPath: "/run/spire/config",
				ReadOnly:  true,
			}
			volMount2 := corev1.VolumeMount{
				Name:      "spire-data",
				MountPath: "/run/spire/data",
				ReadOnly:  false,
			}
			livenessProbe := corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{
					Path: "/live", Port: intstr.IntOrString{IntVal: 8080}, Scheme: "HTTP"}},
				FailureThreshold:    2,
				SuccessThreshold:    1,
				InitialDelaySeconds: 15,
				PeriodSeconds:       60,
				TimeoutSeconds:      3,
			}
			readinessProbe := corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{
					Path: "/ready", Port: intstr.IntOrString{IntVal: 8080}, Scheme: "HTTP"}},
				InitialDelaySeconds: 5,
				TimeoutSeconds:      1,
				PeriodSeconds:       5,
				SuccessThreshold:    1,
				FailureThreshold:    3,
			}
			podVolume := corev1.Volume{
				Name: "spire-config",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{Name: "spire-config-map"},
					},
				},
			}
			containerSpec := corev1.Container{
				Name:           "spire-server",
				Image:          "ghcr.io/spiffe/spire-server:1.5.1",
				Args:           []string{"-config", "/run/spire/config/server.conf"},
				Ports:          []corev1.ContainerPort{{ContainerPort: 8081}},
				VolumeMounts:   []corev1.VolumeMount{volMount1, volMount2},
				LivenessProbe:  &livenessProbe,
				ReadinessProbe: &readinessProbe,
			}
			podSpec := corev1.PodSpec{
				ServiceAccountName: "spire-server",
				Containers:         []corev1.Container{containerSpec},
				Volumes:            []corev1.Volume{podVolume},
			}

			volClaimTemplate := corev1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "spire-data",
					Namespace: "default",
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					Resources: corev1.ResourceRequirements{
						Requests: map[corev1.ResourceName]resource.Quantity{
							corev1.ResourceStorage: resource.MustParse("1Gi"),
						},
					},
				},
			}
			statefulSetSpec := appsv1.StatefulSetSpec{
				Replicas: &numReplicas,
				Selector: &labelSelector,
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: "default",
						Labels:    map[string]string{"app": "spire-server"},
					},
					Spec: podSpec,
				},
				VolumeClaimTemplates: []corev1.PersistentVolumeClaim{volClaimTemplate},
			}
			spireStatefulSet := &appsv1.StatefulSet{
				TypeMeta: metav1.TypeMeta{
					Kind:       "StatefulSet",
					APIVersion: "apps/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "spire-server",
					Namespace: "default",
					Labels:    map[string]string{"app": "spire-server"},
				},
				Spec: statefulSetSpec,
			}
			Expect(k8sClient.Create(ctx, spireStatefulSet)).Should(Succeed())

			/*
				After creating this StatefulSet, let's check that the Spire StatefulSet's Spec fields match what we passed in or not.
				Note that, because the k8s apiserver may not have finished creating a StatefulSet after our `Create()` call from earlier,
				we will use Gomega’s Eventually() testing function instead of Expect() to give the apiserver an opportunity to finish
				creating our Spire StatefulSet.
				`Eventually()` will repeatedly run the function provided as an argument every interval seconds until
				(a) the function’s output matches what’s expected in the subsequent `Should()` call, or
				(b) the number of attempts * interval period exceed the provided timeout value.
				In the examples below, timeout and interval are Go Duration values of our choosing.
			*/
			statefulSetLookupKey := types.NamespacedName{Name: "spire-server", Namespace: "default"}
			createdStatefulSet := &appsv1.StatefulSet{}

			// We'll need to retry getting this newly created Service, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, statefulSetLookupKey, createdStatefulSet)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			Expect(createdStatefulSet.ObjectMeta.Name).Should(Equal("spire-server"))
			Expect(createdStatefulSet.ObjectMeta.Namespace).Should(Equal("default"))
			// Now let us see if the expectation matches or not
			Expect(*createdStatefulSet.Spec.Replicas).Should(Equal(int32(2)))
			//check for storage volume creation
			Expect(createdStatefulSet.Spec.VolumeClaimTemplates[0].Spec.Resources.Requests[corev1.ResourceStorage]).Should(Equal(resource.MustParse("1Gi")))
			Expect(*createdStatefulSet.Spec.Template.Spec.Containers[0].LivenessProbe).Should(Equal(livenessProbe))
			Expect(*createdStatefulSet.Spec.Template.Spec.Containers[0].ReadinessProbe).Should(Equal(readinessProbe))
			Expect(createdStatefulSet.Spec.Template.Spec.Containers[0].VolumeMounts[0]).Should(Equal(volMount1))
			Expect(createdStatefulSet.Spec.Template.Spec.Containers[0].VolumeMounts[1]).Should(Equal(volMount2))
			Expect(createdStatefulSet.Spec.Template.Spec.Volumes[0].VolumeSource.ConfigMap.LocalObjectReference.Name).Should(Equal("spire-config-map"))
			Expect(*createdStatefulSet.Spec.Selector).Should(Equal(labelSelector))
		})
	})

	serverTypeMeta := metav1.TypeMeta{
		APIVersion: "spire.hpe.com/v1",
		Kind:       "SpireServer",
	}
	serverObjectMeta := metav1.ObjectMeta{
		Name:      "invalid-spire-server",
		Namespace: "default",
	}
	Context("When creating SPIRE server with invalid/empty trust domain", func() {
		var ctx = context.Background()
		var spireServer *spirev1.SpireServer
		BeforeEach(func() {
			spireServer = &spirev1.SpireServer{
				TypeMeta:   serverTypeMeta,
				ObjectMeta: serverObjectMeta,
				Spec: spirev1.SpireServerSpec{
					TrustDomain:   "",
					Port:          8081,
					NodeAttestors: []spirev1.NodeAttestor{{Name: "k8s_sat"}},
					KeyStorage:    "disk",
					Replicas:      1,
				},
			}
			Expect(k8sClient.Create(ctx, spireServer)).ShouldNot(Succeed())
		})
		It("should not create the CRD instance", func() {
			serverLookupKey := types.NamespacedName{Name: spireServer.Name, Namespace: spireServer.Namespace}
			createdSpireServer := &spirev1.SpireServer{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, serverLookupKey, createdSpireServer)
				return err != nil
			}, timeout, interval).Should(BeTrue())
		})
	})
	Context("When creating SPIRE Server with unsupported node attestors", func() {
		var ctx = context.Background()
		var spireServer *spirev1.SpireServer
		BeforeEach(func() {
			spireServer = &spirev1.SpireServer{
				TypeMeta:   serverTypeMeta,
				ObjectMeta: serverObjectMeta,
				Spec: spirev1.SpireServerSpec{
					TrustDomain:   "example.org",
					Port:          8081,
					NodeAttestors: []spirev1.NodeAttestor{{Name: "k8s_sat"}, {Name: "aws_iid"}},
					KeyStorage:    "disk",
					Replicas:      1,
				},
			}
			Expect(k8sClient.Create(ctx, spireServer)).ShouldNot(Succeed())
		})
		It("should not create the CRD instance", func() {
			serverLookupKey := types.NamespacedName{Name: spireServer.Name, Namespace: spireServer.Namespace}
			createdSpireServer := &spirev1.SpireServer{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, serverLookupKey, createdSpireServer)
				return err != nil
			}, timeout, interval).Should(BeTrue())
		})
	})
	Context("When creating SPIRE Server with incorrect key storage", func() {
		var ctx = context.Background()
		var spireServer *spirev1.SpireServer
		BeforeEach(func() {
			spireServer = &spirev1.SpireServer{
				TypeMeta:   serverTypeMeta,
				ObjectMeta: serverObjectMeta,
				Spec: spirev1.SpireServerSpec{
					TrustDomain:   "example.org",
					Port:          8081,
					NodeAttestors: []spirev1.NodeAttestor{{Name: "k8s_sat"}},
					KeyStorage:    "drive",
					Replicas:      1,
				},
			}
			Expect(k8sClient.Create(ctx, spireServer)).ShouldNot(Succeed())
		})
		It("should not create the CRD instance", func() {
			serverLookupKey := types.NamespacedName{Name: spireServer.Name, Namespace: spireServer.Namespace}
			createdSpireServer := &spirev1.SpireServer{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, serverLookupKey, createdSpireServer)
				return err != nil
			}, timeout, interval).Should(BeTrue())
		})
	})
	Context("When creating SPIRE Server with invalid port number", func() {
		var ctx = context.Background()
		var spireServer *spirev1.SpireServer
		BeforeEach(func() {
			spireServer = &spirev1.SpireServer{
				TypeMeta:   serverTypeMeta,
				ObjectMeta: serverObjectMeta,
				Spec: spirev1.SpireServerSpec{
					TrustDomain:   "example.org",
					Port:          -1,
					NodeAttestors: []spirev1.NodeAttestor{{Name: "k8s_sat"}},
					KeyStorage:    "disk",
					Replicas:      1,
				},
			}
			Expect(k8sClient.Create(ctx, spireServer)).ShouldNot(Succeed())
		})
		It("should not create the CRD instance", func() {
			serverLookupKey := types.NamespacedName{Name: spireServer.Name, Namespace: spireServer.Namespace}
			createdSpireServer := &spirev1.SpireServer{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, serverLookupKey, createdSpireServer)
				return err != nil
			}, timeout, interval).Should(BeTrue())
		})
	})
	Context("When creating SPIRE Server with invalid replicas", func() {
		var ctx = context.Background()
		var spireServer *spirev1.SpireServer
		BeforeEach(func() {
			spireServer = &spirev1.SpireServer{
				TypeMeta:   serverTypeMeta,
				ObjectMeta: serverObjectMeta,
				Spec: spirev1.SpireServerSpec{
					TrustDomain:   "example.org",
					Port:          8081,
					NodeAttestors: []spirev1.NodeAttestor{{Name: "k8s_sat"}},
					KeyStorage:    "disk",
					Replicas:      -1,
				},
			}
			Expect(k8sClient.Create(ctx, spireServer)).ShouldNot(Succeed())
		})
		It("should not create the CRD instance", func() {
			serverLookupKey := types.NamespacedName{Name: spireServer.Name, Namespace: spireServer.Namespace}
			createdSpireServer := &spirev1.SpireServer{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, serverLookupKey, createdSpireServer)
				return err != nil
			}, timeout, interval).Should(BeTrue())
		})
	})
	Context("When creating SPIRE Server with invalid datastore", func() {
		var ctx = context.Background()
		var spireServer *spirev1.SpireServer
		BeforeEach(func() {
			spireServer = &spirev1.SpireServer{
				TypeMeta:   serverTypeMeta,
				ObjectMeta: serverObjectMeta,
				Spec: spirev1.SpireServerSpec{
					TrustDomain:   "example.org",
					Port:          8081,
					NodeAttestors: []spirev1.NodeAttestor{{Name: "k8s_sat"}},
					KeyStorage:    "disk",
					Replicas:      1,
					DataStore:     "cloud",
				},
			}
			Expect(k8sClient.Create(ctx, spireServer)).ShouldNot(Succeed())
		})
		It("should not create the CRD instance", func() {
			serverLookupKey := types.NamespacedName{Name: spireServer.Name, Namespace: spireServer.Namespace}
			createdSpireServer := &spirev1.SpireServer{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, serverLookupKey, createdSpireServer)
				return err != nil
			}, timeout, interval).Should(BeTrue())
		})
	})
	Context("When creating SPIRE Server with empty connection string", func() {
		var ctx = context.Background()
		var spireServer *spirev1.SpireServer
		BeforeEach(func() {
			spireServer = &spirev1.SpireServer{
				TypeMeta:   serverTypeMeta,
				ObjectMeta: serverObjectMeta,
				Spec: spirev1.SpireServerSpec{
					TrustDomain:      "example.org",
					Port:             8081,
					NodeAttestors:    []spirev1.NodeAttestor{{Name: "k8s_sat"}},
					KeyStorage:       "disk",
					Replicas:         1,
					DataStore:        "sqlite3",
					ConnectionString: "",
				},
			}
			Expect(k8sClient.Create(ctx, spireServer)).ShouldNot(Succeed())
		})
		It("should not create the CRD instance", func() {
			serverLookupKey := types.NamespacedName{Name: spireServer.Name, Namespace: spireServer.Namespace}
			createdSpireServer := &spirev1.SpireServer{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, serverLookupKey, createdSpireServer)
				return err != nil
			}, timeout, interval).Should(BeTrue())
		})
	})
	Context("When creating SPIRE Server with sqlite3 datastore and > 1 replicas", func() {
		var ctx = context.Background()
		var spireServer *spirev1.SpireServer
		BeforeEach(func() {
			spireServer = &spirev1.SpireServer{
				TypeMeta:   serverTypeMeta,
				ObjectMeta: serverObjectMeta,
				Spec: spirev1.SpireServerSpec{
					TrustDomain:      "example.org",
					Port:             8081,
					NodeAttestors:    []spirev1.NodeAttestor{{Name: "k8s_sat"}},
					KeyStorage:       "disk",
					Replicas:         3,
					DataStore:        "sqlite3",
					ConnectionString: "/run/spire/data/datastore.sqlite3",
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
	Context("When creating SPIRE Server with all valid fields", func() {
		var ctx = context.Background()
		var spireServer *spirev1.SpireServer
		BeforeEach(func() {
			spireServer = &spirev1.SpireServer{
				TypeMeta:   serverTypeMeta,
				ObjectMeta: serverObjectMeta,
				Spec: spirev1.SpireServerSpec{
					TrustDomain:      "example.org",
					Port:             8081,
					NodeAttestors:    []spirev1.NodeAttestor{{Name: "k8s_sat"}},
					KeyStorage:       "disk",
					Replicas:         1,
					DataStore:        "sqlite3",
					ConnectionString: "/run/spire/data/datastore.sqlite3",
				},
			}
			Expect(k8sClient.Create(ctx, spireServer)).Should(Succeed())
		})
		It("should create and NOT delete the CRD instance", func() {
			serverLookupKey := types.NamespacedName{Name: spireServer.Name, Namespace: spireServer.Namespace}
			createdSpireServer := &spirev1.SpireServer{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, serverLookupKey, createdSpireServer)
				return err == nil
			}, timeout, interval).Should(BeTrue())
		})
	})
})
