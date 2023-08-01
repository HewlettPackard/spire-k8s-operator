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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SpireAgentSpec defines the desired state of SpireAgent
type SpireAgentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	TrustDomain string `json:"trustDomain"`

	NodeAttestor string `json:"nodeAttestor"`

	WorkloadAttestors []string `json:"workloadAttestors"`

	KeyStorage string `json:"keyStorage"`

	ServerPort int `json:"serverPort"`

	// The path to the private key on disk for the x509pop node attestor
	// +kubebuilder:validation:Optional
	PrivateKeyPath string `json:"privateKeyPath"`

	// The path to the certificate bundle on disk for the x509pop node attestor
	// +kubebuilder:validation:Optional
	CertificatePath string `json:"certificatePath"`

	// The path to the private key on disk in openssh format for the sshpop node attestor
	// +kubebuilder:validation:Optional
	HostKeyPath string `json:"hostKeyPath"`

	// The path to the certificate on disk in openssh format for the sshpop node attestor
	// +kubebuilder:validation:Optional
	HostCertPath string `json:"hostCertPath"`
}

// SpireAgentStatus defines the observed state of SpireAgent
type SpireAgentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SpireAgent is the Schema for the spireagents API
type SpireAgent struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpireAgentSpec   `json:"spec,omitempty"`
	Status SpireAgentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SpireAgentList contains a list of SpireAgent
type SpireAgentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SpireAgent `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SpireAgent{}, &SpireAgentList{})
}