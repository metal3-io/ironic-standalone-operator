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

// Database defines database settings
type Database struct {
	// Image is the MariaDB image.
	// +kubebuilder:default=quay.io/metal3-io/mariadb
	// +optional
	Image string `json:"image,omitempty"`

	// ExternalIP can be set to use an existing MariaDB installation instead of a managed one.
	ExternalIP string `json:"externalIP,omitempty"`

	// CredentialsSecretName is the name of the secret with database credentials.
	CredentialsSecretName string `json:"credentialsSecretName,omitempty"`

	// TLSSecretName is the name of the secret with the database TLS certificate.
	// +optional
	TLSSecretName string `json:"tlsSecretName,omitempty"`
}

// Inspection defines inspection settings
type Inspection struct {
	// Collectors is a list of inspection collectors to enable.
	Collectors []string `json:"collectors,omitempty"`

	// List of interfaces to inspect for VLANs.
	VLANInterfaces []string `json:"vlanInterfaces,omitempty"`
}

// Networking defines networking settings for Ironic
type Networking struct {
	// Interface is a Linux network device to listen on.
	// +optional
	Interface string `json:"interface,omitempty"`

	// MACAddresses can be provided to make the start script pick the interface matching any of these addresses.
	// +optional
	MACAddresses []string `json:"macAddresses,omitempty"`

	// BindInterface makes Ironic API bound to only one interface.
	BindInterface bool `json:"bindInterface,omitempty"`

	// ExternalIP is used for accessing API and the image server from remote hosts.
	ExternalIP string `json:"externalIP,omitempty"`
}

// IronicSpec defines the desired state of Ironic
type IronicSpec struct {
	// Image is the Ironic image (including httpd).
	// +kubebuilder:default=quay.io/metal3-io/ironic
	// +optional
	Image string `json:"image,omitempty"`

	// Size is the desired count of Ironic instances.
	// +kubebuilder:validation:Minimum=1
	// +optional
	Size int32 `json:"size,omitempty"`

	// Database defines database settings for Ironic. Mandatory when size is more than 1.
	Database *Database `json:"database,omitempty"`

	// Inspection defines inspection settings
	Inspection Inspection `json:"inspection,omitempty"`

	// Networking defines networking settings for Ironic.
	// At least one of the subfield should be provided.
	// +optional
	Networking Networking `json:"networking,omitempty"`

	// APIPort is the public port used for Ironic.
	// +kubebuilder:default=6385
	// +kubebuilder:validation:Minimum=1
	// +optional
	APIPort int32 `json:"apiPort,omitempty"`

	// ImageServerPort is the public port used for serving images.
	// +kubebuilder:default=8088
	// +kubebuilder:validation:Minimum=1
	// +optional
	ImageServerPort int32 `json:"imageServerPort,omitempty"`

	// ImageServerTLSPort is the public port used for serving virtual media images over TLS.
	// Setting it to 0 disables TLS for virtual media.
	// +kubebuilder:default=8089
	// +kubebuilder:validation:Minimum=0
	// +optional
	ImageServerTLSPort int32 `json:"imageServerTLSPort,omitempty"`

	// APISecretName is the name of the secret with Ironic API credentials.
	APISecretName string `json:"apiSecretName"`

	// TLSSecretName is the name of the secret with the API TLS certificate.
	// +optional
	TLSSecretName string `json:"tlsSecretName,omitempty"`

	// RamdiskSSHKey is the contents of the public key to inject into the ramdisk for debugging purposes.
	RamdiskSSHKey string `json:"ramdiskSSHKey,omitempty"`
}

type IronicStatusConditionType string

const (
	// Available indicates that Ironic is fully available
	IronicStatusAvailable IronicStatusConditionType = "Available"

	// Progressing indicates that Ironic deployment is in progress
	IronicStatusProgressing IronicStatusConditionType = "Progressing"

	// Degraded indicates that Ironic deployment cannot progress
	IronicStatusDegraded IronicStatusConditionType = "Degraded"
)

// IronicStatus defines the observed state of Ironic
type IronicStatus struct {
	// IronicEndpoints is the available API endpoints of Ironic.
	// +optional
	IronicEndpoints []string `json:"ironicEndpoint,omitempty"`

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

// Ironic is the Schema for the ironics API
type Ironic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IronicSpec   `json:"spec,omitempty"`
	Status IronicStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IronicList contains a list of Ironic
type IronicList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Ironic `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Ironic{}, &IronicList{})
}
