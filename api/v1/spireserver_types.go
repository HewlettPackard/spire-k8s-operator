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

// SpireServerSpec defines the desired state of SpireServer
type SpireServerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	TrustDomain string `json:"trustDomain"`

	Port int `json:"port"`

	NodeAttestors []string `json:"nodeAttestors"`

	KeyStorage string `json:"keyStorage"`

	Replicas int `json:"replicas"`

	DataStore string `json:"dataStore"`

	ConnectionString string `json:"connectionString"`

	// The path to the trusted CA bundle on disk for the x509pop node attestor
	// +kubebuilder:validation:Optional
	CABundlePath string `json:"caBundlePath"`

	// A list of trusted CAs in ssh authorized_keys format for the sshpop node attestor
	// +kubebuilder:validation:Optional
	CertAuthorities []string `json:"certAuthorities"`

	// A file that contains a list of trusted CAs in ssh authorized_keys format for the sshpop node attestor
	// +kubebuilder:validation:Optional
	CertAuthoritiesPath string `json:"certAuthoritiesPath"`
}

// SpireServerStatus defines the observed state of SpireServer
type SpireServerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Health string `json:"health"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Health",type=string,JSONPath=`.status.health`

// SpireServer is the Schema for the spireservers API
type SpireServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpireServerSpec   `json:"spec,omitempty"`
	Status SpireServerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SpireServerList contains a list of SpireServer
type SpireServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SpireServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SpireServer{}, &SpireServerList{})
}
