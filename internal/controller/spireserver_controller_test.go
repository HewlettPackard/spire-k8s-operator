package controller

import (
	"context"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SpireServer controller", func() {
	const (
		duration = time.Second * 10
		interval = time.Millisecond * 250
		timeout  = time.Second * 10
	)
	Context("When installing SPIRE server", func() {
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
					Path: "/live", Port: intstr.IntOrString{IntVal: 8080}}},
				FailureThreshold:    2,
				InitialDelaySeconds: 15,
				PeriodSeconds:       60,
				TimeoutSeconds:      3,
			}
			readinessProbe := corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{
					Path: "/ready", Port: intstr.IntOrString{IntVal: 8080}}},
				InitialDelaySeconds: 5,
				PeriodSeconds:       5,
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
			statefulSetLookupKey := types.NamespacedName{Name: "spire-service", Namespace: "default"}
			createdStatefulSet := &appsv1.StatefulSet{}

			// We'll need to retry getting this newly created Service, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, statefulSetLookupKey, createdStatefulSet)
				return err == nil
			}, timeout, interval).Should(BeTrue())

			// Now let us see if the expectation matches or not
			Expect(createdStatefulSet.Spec.Replicas).Should(Equal(int32(2)))
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

			// Now let us see if the expectation matches or not
			Expect(createdBundle.Name).Should(Equal("spire-bundle"))
			Expect(createdBundle.Namespace).Should(Equal("default"))
			Expect(createdBundle.Data).Should(Not(Equal(nil)))
			Expect(createdBundle.BinaryData).Should(Not(Equal(nil)))
			Expect(createdBundle.Labels).Should(Not(Equal(nil)))
			Expect(createdBundle.Annotations).Should(Not(Equal(nil)))
		})
	})
})
