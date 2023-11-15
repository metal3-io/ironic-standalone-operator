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

const (
	DefaultAPIPort            int32  = 6385
	DefaultImageServerPort    int32  = 6180
	DefaultImageServerTLSPort int32  = 6183
	DefaultIronicImage        string = "quay.io/metal3-io/ironic"

	IronicFinalizer string = "ironic.metal3.io"
)

// Inspection defines inspection settings
type Inspection struct {
	// Collectors is a list of inspection collectors to enable.
	Collectors []string `json:"collectors,omitempty"`

	// List of interfaces to inspect for VLANs.
	VLANInterfaces []string `json:"vlanInterfaces,omitempty"`
}

type DHCP struct {
	// NetworkCIDR is a CIRD of the provisioning network.
	NetworkCIDR string `json:"networkCIDR,omitempty"`

	// FirstIP is the first IP that can be given to hosts. Must be inside NetworkCIDR.
	// If not set, the 10th IP from NetworkCIDR is used (e.g. .10 for /24).
	// +optional
	FirstIP string `json:"firstIP,omitempty"`

	// LastIP is the last IP that can be given to hosts. Must be inside NetworkCIDR.
	// If not set, the 2nd IP from the end of NetworkCIDR is used (e.g. .253 for /24).
	// +optional
	LastIP string `json:"lastIP,omitempty"`

	// ServeDNS is set to true to pass the provisioning host as the DNS server on the provisioning network.
	// Must not be set together with DNSAddress.
	// +optional
	ServeDNS bool `json:"serveDNS,omitempty"`

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
}

// Networking defines networking settings for Ironic
type Networking struct {
	// Interface is a Linux network device to listen on.
	// Detected from IPAddress if missing.
	// +optional
	Interface string `json:"interface,omitempty"`

	// IPAddress is the main IP address to listen on and use for communication.
	// Detected from Interface if missing.
	// +optional
	IPAddress string `json:"ipAddress,omitempty"`

	// MACAddresses can be provided to make the start script pick the interface matching any of these addresses.
	// Only set if no other options can be used.
	// +optional
	MACAddresses []string `json:"macAddresses,omitempty"`

	// BindInterface makes Ironic API bound to only one interface.
	// +optional
	BindInterface bool `json:"bindInterface,omitempty"`

	// ExternalIP is used for accessing API and the image server from remote hosts.
	// +optional
	ExternalIP string `json:"externalIP,omitempty"`

	// APIPort is the public port used for Ironic.
	// +kubebuilder:default=6385
	// +kubebuilder:validation:Minimum=1
	// +optional
	APIPort int32 `json:"apiPort,omitempty"`

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

	// DHCP is a configuration of DHCP for the network boot service (dnsmasq).
	// The service is only deployed when this is set.
	DHCP *DHCP `json:"dhcp,omitempty"`
}

// IronicSpec defines the desired state of Ironic
type IronicSpec struct {
	// Image is the Ironic image (including httpd).
	// +kubebuilder:default=quay.io/metal3-io/ironic
	// +kubebuilder:validation:MinLength=1
	// +optional
	Image string `json:"image,omitempty"`

	// RamdiskDownloaderImage is the image to be used at pod initialization to download the IPA ramdisk.
	// +kubebuilder:default=quay.io/metal3-io/ironic-ipa-downloader
	// +optional
	RamdiskDownloaderImage string `json:"ramdiskDownloaderImage,omitempty"`

	// Distributed causes Ironic to be deployed as a DaemonSet on control plane nodes instead of a deployment with 1 replica.
	// Requires database to be installed and linked to DatabaseName.
	// +optional
	Distributed bool `json:"distributed,omitempty"`

	// DatabaseName defines database settings for Ironic.
	DatabaseName string `json:"databaseName,omitempty"`

	// Inspection defines inspection settings
	Inspection Inspection `json:"inspection,omitempty"`

	// Networking defines networking settings for Ironic.
	// +optional
	Networking Networking `json:"networking,omitempty"`

	// DisableVirtualMediaTLS turns off TLS on the virtual media server,
	// which may be required for hardware that cannot accept HTTPS links.
	// +optional
	DisableVirtualMediaTLS bool `json:"disableVirtualMediaTLS,omitempty"`

	// APISecretName is the name of the secret with Ironic API credentials.
	APISecretName string `json:"apiSecretName"`

	// TLSSecretName is the name of the secret with the API TLS certificate.
	// +optional
	TLSSecretName string `json:"tlsSecretName,omitempty"`

	// RamdiskSSHKey is the contents of the public key to inject into the ramdisk for debugging purposes.
	// +optional
	RamdiskSSHKey string `json:"ramdiskSSHKey,omitempty"`
}

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
