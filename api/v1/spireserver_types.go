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
	// +kubebuilder:validation:Required

	// +kubebuilder:validation:Pattern="[a-z0-9._-]{1,255}"
	TrustDomain string `json:"trustDomain"`

	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=65535
	Port int `json:"port"`

	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:UniqueItems=true
	NodeAttestors []string `json:"nodeAttestors"`

	// +kubebuilder:validation:Enum=disk;memory
	KeyStorage string `json:"keyStorage"`

	// +kubebuilder:validation:Minimum=1
	Replicas int `json:"replicas"`

	// +kubebuilder:validation:Enum=sqlite3;postgres;mysql
	DataStore string `json:"dataStore"`

	// +kubebuilder:validation:MinLength=1
	ConnectionString string `json:"connectionString"`
}

// SpireServerStatus defines the observed state of SpireServer
type SpireServerStatus struct {
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
