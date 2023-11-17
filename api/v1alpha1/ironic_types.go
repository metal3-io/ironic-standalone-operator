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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Inspection defines inspection settings
type Inspection struct {
	// Collectors is a list of inspection collectors to enable.
	Collectors []string `json:"collectors,omitempty"`

	// List of interfaces to inspect for VLANs.
	VLANInterfaces []string `json:"vlanInterfaces,omitempty"`
}

type DHCP struct {
	// DNSAddress is the IP address of the DNS server to pass to hosts via DHCP.
	// Must not be set together with ServeDNS.
	// +optional
	DNSAddress string `json:"dnsAddress,omitempty"`

	// GatewayAddress is the IP address of the gateway to pass to hosts via DHCP.
	// +optional
	GatewayAddress string `json:"gatewayAddress,omitempty"`

	// Hosts is a set of DHCP host records to pass to dnsmasq.
	// Check the dnsmasq documentation on dhcp-host for an explanation of the format.
	// There is no API-side validation. Most users will leave this unset.
	// +optional
	Hosts []string `json:"hosts,omitempty"`

	// Ignore is set of dnsmasq tags to ignore and not provide any DHCP.
	// Check the dnsmasq documentation on dhcp-ignore for an explanation of the format.
	// There is no API-side validation. Most users will leave this unset.
	// +optional
	Ignore []string `json:"ignore,omitempty"`

	// NetworkCIDR is a CIRD of the provisioning network. Required.
	NetworkCIDR string `json:"networkCIDR,omitempty"`

	// RangeBegin is the first IP that can be given to hosts. Must be inside NetworkCIDR.
	// If not set, the 10th IP from NetworkCIDR is used (e.g. .10 for /24).
	// +optional
	RangeBegin string `json:"rangeBegin,omitempty"`

	// RangeEnd is the last IP that can be given to hosts. Must be inside NetworkCIDR.
	// If not set, the 2nd IP from the end of NetworkCIDR is used (e.g. .253 for /24).
	// +optional
	RangeEnd string `json:"rangeEnd,omitempty"`

	// ServeDNS is set to true to pass the provisioning host as the DNS server on the provisioning network.
	// Must not be set together with DNSAddress.
	// +optional
	ServeDNS bool `json:"serveDNS,omitempty"`
}

// Networking defines networking settings for Ironic
type Networking struct {
	// APIPort is the public port used for Ironic.
	// +kubebuilder:default=6385
	// +kubebuilder:validation:Minimum=1
	// +optional
	APIPort int32 `json:"apiPort,omitempty"`

	// BindInterface makes Ironic API bound to only one interface.
	// +optional
	BindInterface bool `json:"bindInterface,omitempty"`

	// DHCP is a configuration of DHCP for the network boot service (dnsmasq).
	// The service is only deployed when this is set.
	DHCP *DHCP `json:"dhcp,omitempty"`

	// ExternalIP is used for accessing API and the image server from remote hosts.
	// This settings only applies to virtual media deployments.
	// +optional
	ExternalIP string `json:"externalIP,omitempty"`

	// ImageServerPort is the public port used for serving images.
	// +kubebuilder:default=6180
	// +kubebuilder:validation:Minimum=1
	// +optional
	ImageServerPort int32 `json:"imageServerPort,omitempty"`

	// ImageServerTLSPort is the public port used for serving virtual media images over TLS.
	// +kubebuilder:default=6183
	// +kubebuilder:validation:Minimum=1
	// +optional
	ImageServerTLSPort int32 `json:"imageServerTLSPort,omitempty"`

	// Interface is a Linux network device to listen on.
	// Detected from IPAddress if missing.
	// +optional
	Interface string `json:"interface,omitempty"`

	// IPAddress is the main IP address to listen on and use for communication.
	// Detected from Interface if missing. Cannot be provided for a distributed architecture.
	// +optional
	IPAddress string `json:"ipAddress,omitempty"`

	// MACAddresses can be provided to make the start script pick the interface matching any of these addresses.
	// Only set if no other options can be used.
	// +optional
	MACAddresses []string `json:"macAddresses,omitempty"`
}

type Images struct {
	// Ironic is the Ironic image (including httpd).
	// +kubebuilder:default=quay.io/metal3-io/ironic
	// +kubebuilder:validation:MinLength=1
	// +optional
	Ironic string `json:"ironic,omitempty"`

	// RamdiskDownloader is the image to be used at pod initialization to download the IPA ramdisk.
	// +kubebuilder:default=quay.io/metal3-io/ironic-ipa-downloader
	// +optional
	RamdiskDownloader string `json:"ramdiskDownloader,omitempty"`
}

// IronicSpec defines the desired state of Ironic
type IronicSpec struct {
	// CredentialsRef is a reference to the secret with Ironic API credentials.
	// A new secret will be created if this field is empty.
	// +optional
	CredentialsRef corev1.LocalObjectReference `json:"credentialsRef,omitempty"`

	// DatabaseRef defines database settings for Ironic.
	// If missing, a local SQLite database will be used. Must be provided for a distributed architecture.
	// +optional
	DatabaseRef corev1.LocalObjectReference `json:"databaseRef,omitempty"`

	// DisableVirtualMediaTLS turns off TLS on the virtual media server,
	// which may be required for hardware that cannot accept HTTPS links.
	// +optional
	DisableVirtualMediaTLS bool `json:"disableVirtualMediaTLS,omitempty"`

	// Distributed causes Ironic to be deployed as a DaemonSet on control plane nodes instead of a deployment with 1 replica.
	// Requires database to be installed and linked to DatabaseRef.
	// EXPERIMENTAL: do not use (validation will fail)!
	// +optional
	Distributed bool `json:"distributed,omitempty"`

	// Images is a collection of container images to deploy from.
	Images Images `json:"images,omitempty"`

	// Inspection defines inspection settings
	Inspection Inspection `json:"inspection,omitempty"`

	// Networking defines networking settings for Ironic.
	// +optional
	Networking Networking `json:"networking,omitempty"`

	// TLSSecretName is a reference to the secret with the database TLS certificate.
	// +optional
	TLSRef corev1.LocalObjectReference `json:"tlsRef,omitempty"`

	// RamdiskSSHKey is the contents of the public key to inject into the ramdisk for debugging purposes.
	// +optional
	RamdiskSSHKey string `json:"ramdiskSSHKey,omitempty"`
}

// IronicStatus defines the observed state of Ironic
type IronicStatus struct {
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
