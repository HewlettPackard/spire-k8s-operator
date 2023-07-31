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
	// The trust domain associated with the SPIRE server
	TrustDomain string `json:"trustDomain"`

	// The port on which the SPIRE server listens to agents
	Port int `json:"port"`

	// The node attestor plugins the SPIRE server uses
	NodeAttestors []string `json:"nodeAttestors"`

	// Indicates whether the generated keys are stored on disk or in memory
	KeyStorage string `json:"keyStorage"`

	// Number of replicas for SPIRE server
	Replicas int `json:"replicas"`

	// Indicates how server data should be stored (sqlite3, mysql, or postgres)
	DataStore string `json:"dataStore"`

	// Connection string for the datastore
	ConnectionString string `json:"connectionString"`
}

// SpireServerStatus defines the observed state of SpireServer
type SpireServerStatus struct {
	// Indicates whether the SPIRE server is in an error state (ERROR), initializing (INIT), live (LIVE), or ready (READY)
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
