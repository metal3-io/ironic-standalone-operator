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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IronicDatabaseSpec defines the desired state of IronicDatabase.
type IronicDatabaseSpec struct {
	// CredentialsName is a reference to the secret with database credentials.
	// A new secret will be created if this field is empty.
	// +optional
	CredentialsName string `json:"credentialsName,omitempty"`

	// Image is the MariaDB image.
	// +optional
	Image string `json:"image,omitempty"`

	// TLSCertificateName is a reference to the secret with the database TLS certificate.
	// +optional
	TLSCertificateName string `json:"tlsCertificateName,omitempty"`
}

// IronicDatabaseStatus defines the observed state of IronicDatabase.
type IronicDatabaseStatus struct {
	// Conditions describe the state of the Ironic deployment.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// IronicDatabase is the Schema for the ironic database API.
//
// Deprecated: the IronicDatabase API is deprecated and will be removed soon in favour of 3rd party database operators.
type IronicDatabase struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IronicDatabaseSpec   `json:"spec,omitempty"`
	Status IronicDatabaseStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IronicDatabaseList contains a list of IronicDatabase.
type IronicDatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IronicDatabase `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &IronicDatabase{}, &IronicDatabaseList{})
}
