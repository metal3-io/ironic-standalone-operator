# API Reference

Packages:

- [ironic.metal3.io/v1alpha1](#ironicmetal3iov1alpha1)

# ironic.metal3.io/v1alpha1

Resource Types:

- [Ironic](#ironic)




## Ironic
<sup><sup>[↩ Parent](#ironicmetal3iov1alpha1 )</sup></sup>






Ironic is the Schema for the ironics API.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>ironic.metal3.io/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>Ironic</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#ironicspec">spec</a></b></td>
        <td>object</td>
        <td>
          IronicSpec defines the desired state of Ironic.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicstatus">status</a></b></td>
        <td>object</td>
        <td>
          IronicStatus defines the observed state of Ironic.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec
<sup><sup>[↩ Parent](#ironic)</sup></sup>



IronicSpec defines the desired state of Ironic.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>apiCredentialsName</b></td>
        <td>string</td>
        <td>
          APICredentialsName is a reference to the secret with Ironic API credentials.
A new secret will be created if this field is empty.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecdatabase">database</a></b></td>
        <td>object</td>
        <td>
          Database is a reference to a MariaDB database to use for persisting Ironic data.
Must be provided for a highly available architecture, optional otherwise.
If missing, a local SQLite database will be used, and the Ironic state will be reset on each pod restart.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecdeployramdisk">deployRamdisk</a></b></td>
        <td>object</td>
        <td>
          DeployRamdisk defines settings for the provisioning/inspection ramdisk based on Ironic Python Agent.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecextraconfigindex">extraConfig</a></b></td>
        <td>[]object</td>
        <td>
          ExtraConfig allows overriding any Ironic configuration options.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>highAvailability</b></td>
        <td>boolean</td>
        <td>
          HighAvailability causes Ironic to be deployed as a DaemonSet on control plane nodes instead of a deployment with 1 replica.
Requires database to be installed and linked in the Database field.
DHCP support is not yet implemented in the highly available architecture.
Requires the HighAvailability feature gate to be set.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecimages">images</a></b></td>
        <td>object</td>
        <td>
          Images is a collection of container images to deploy from.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecinspection">inspection</a></b></td>
        <td>object</td>
        <td>
          Inspection defines inspection settings.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecnetworking">networking</a></b></td>
        <td>object</td>
        <td>
          Networking defines networking settings for Ironic.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecnetworkingservice">networkingService</a></b></td>
        <td>object</td>
        <td>
          NetworkingService provides configuration for the Ironic Networking Service<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>nodeSelector</b></td>
        <td>map[string]string</td>
        <td>
          NodeSelector is a selector which must be true for the Ironic pod to fit on a node.
Selector which must match a node's labels for the vmi to be scheduled on that node.
More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverrides">overrides</a></b></td>
        <td>object</td>
        <td>
          Overrides for the generated Deployment or Daemon Set.
EXPERIMENTAL: requires feature gate Overrides.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecprometheusexporter">prometheusExporter</a></b></td>
        <td>object</td>
        <td>
          PrometheusExporter configures sensor data collection and Prometheus metrics export.
When enabled, this configures Ironic to collect sensor data and deploys the
ironic-prometheus-exporter container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspectls">tls</a></b></td>
        <td>object</td>
        <td>
          TLS defines TLS-related settings for various network interactions.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>version</b></td>
        <td>string</td>
        <td>
          Version is the version of Ironic to be installed.
Must be either "latest" or a MAJOR.MINOR pair, e.g. "27.0".
The default version depends on the operator branch.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.database
<sup><sup>[↩ Parent](#ironicspec)</sup></sup>



Database is a reference to a MariaDB database to use for persisting Ironic data.
Must be provided for a highly available architecture, optional otherwise.
If missing, a local SQLite database will be used, and the Ironic state will be reset on each pod restart.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>credentialsName</b></td>
        <td>string</td>
        <td>
          Name of a secret with database credentials.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          IP address or host name of the database instance.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Database name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>tlsCertificateName</b></td>
        <td>string</td>
        <td>
          Name of a secret with the a TLS certificate or a CA for verifying the database host.
If unset, Ironic will request an unencrypted connections, which is insecure,
and the server configuration may forbid it.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.deployRamdisk
<sup><sup>[↩ Parent](#ironicspec)</sup></sup>



DeployRamdisk defines settings for the provisioning/inspection ramdisk based on Ironic Python Agent.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>disableDownloader</b></td>
        <td>boolean</td>
        <td>
          DisableDownloader tells the operator not to start the IPA downloader as the init container.
The user will be responsible for providing the right image to BareMetal Operator.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>extraKernelParams</b></td>
        <td>string</td>
        <td>
          ExtraKernelParams is a string with kernel parameters to pass to the provisioning/inspection ramdisk.
Will not take effect if the host uses a pre-built ISO (either through its PreprovisioningImage or via the DEPLOY_ISO_URL baremetal-operator parameter).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sshKey</b></td>
        <td>string</td>
        <td>
          SSHKey is the contents of the public key to inject into the ramdisk for debugging purposes.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.extraConfig[index]
<sup><sup>[↩ Parent](#ironicspec)</sup></sup>



ExtraConfig allows overriding any Ironic configuration options.
See the entire listing of available options in the Ironic documentation:
https://docs.openstack.org/ironic/latest/configuration/config.html
(note that some options may not be available in earlier releases).

Warning: modifying arbitrary options may cause your Ironic installation to
fail or misbehave. Do not modify anything you don't understand well.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>group</b></td>
        <td>string</td>
        <td>
          The group that config belongs to.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The name of the config.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The value of the config.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.images
<sup><sup>[↩ Parent](#ironicspec)</sup></sup>



Images is a collection of container images to deploy from.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>deployRamdiskBranch</b></td>
        <td>string</td>
        <td>
          DeployRamdiskBranch is the branch of IPA to download. The main branch is used by default.
Not used if deployRamdisk.disableDownloader is true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>deployRamdiskDownloader</b></td>
        <td>string</td>
        <td>
          DeployRamdiskDownloader is the image to be used at pod initialization to download the IPA ramdisk.
Not used if deployRamdisk.disableDownloader is true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ironic</b></td>
        <td>string</td>
        <td>
          Ironic is the Ironic image (including httpd).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keepalived</b></td>
        <td>string</td>
        <td>
          Keepalived is the Keepalived image.
Not used if networking.ipAddressManager is not set to keepalived.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.inspection
<sup><sup>[↩ Parent](#ironicspec)</sup></sup>



Inspection defines inspection settings.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>collectors</b></td>
        <td>[]string</td>
        <td>
          Collectors is a list of inspection collectors to enable.
See https://docs.openstack.org/ironic-python-agent/latest/admin/how_it_works.html#inspection-data for details.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>vlanInterfaces</b></td>
        <td>[]string</td>
        <td>
          List of interfaces to inspect for VLANs.
This can be interface names (to collect all VLANs using LLDP) or pairs <interface>.<vlan ID>.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.networking
<sup><sup>[↩ Parent](#ironicspec)</sup></sup>



Networking defines networking settings for Ironic.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>apiPort</b></td>
        <td>integer</td>
        <td>
          APIPort is the public port used for Ironic.<br/>
          <br/>
            <i>Format</i>: int32<br/>
            <i>Default</i>: 6385<br/>
            <i>Minimum</i>: 1<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>bindInterface</b></td>
        <td>boolean</td>
        <td>
          BindInterface makes Ironic API bound to only one interface.
Requires DisableHostNetwork to be false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecnetworkingdhcp">dhcp</a></b></td>
        <td>object</td>
        <td>
          DHCP is a configuration of DHCP for the network boot service (dnsmasq).
The service is only deployed when this is set.
This setting is currently incompatible with the highly available architecture.
Requires DisableHostNetwork to be false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>disableHostNetwork</b></td>
        <td>boolean</td>
        <td>
          DisableHostNetwork disables the use of host networking for Ironic pods.
Disabling host networking makes network boot impossible.
This should only be used with virtual media deployments.<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>externalCallbackURL</b></td>
        <td>string</td>
        <td>
          ExternalCallbackURL for Ironic API server.
Set this option when your Ironic API server is not directly accessible.
Setting this option, will override URL set by networking.ingress.host.
Must be set together with networking.imageServerExternalURL or networking.ingress
This should only be used with virtual media deployments.
Cannot be set at the same time with networking.externalIP.<br/>
          <br/>
            <i>Format</i>: uri<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>externalIP</b></td>
        <td>string</td>
        <td>
          ExternalIP is used for accessing API and the image server from remote hosts.
This setting only applies to virtual media deployments. The IP will not be accessed from the cluster itself.
Cannot be set at the same time with networking.ingress, networking.externalCallbackURL, or networking.imageServerExternalURL.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>imageServerExternalURL</b></td>
        <td>string</td>
        <td>
          ImageServerExternalURL is to set external HTTP URL for Image server.
Set this option when your image server is not directly accessible.
Setting this option, will override URL set by networking.ingress.host.
Must be set together with networking.externalCallbackURL or networking.ingress
Cannot be set at the same time with networking.externalIP.
This should only be used with virtual media deployments.<br/>
          <br/>
            <i>Format</i>: uri<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>imageServerPort</b></td>
        <td>integer</td>
        <td>
          ImageServerPort is the public port used for serving images.<br/>
          <br/>
            <i>Format</i>: int32<br/>
            <i>Default</i>: 6180<br/>
            <i>Minimum</i>: 1<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>imageServerTLSPort</b></td>
        <td>integer</td>
        <td>
          ImageServerTLSPort is the public port used for serving virtual media images over TLS.<br/>
          <br/>
            <i>Format</i>: int32<br/>
            <i>Default</i>: 6183<br/>
            <i>Minimum</i>: 1<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecnetworkingingress">ingress</a></b></td>
        <td>object</td>
        <td>
          Configure Ingress resource for Ironic services.
Set this option when you are planning to deploy Ironic in a public cluster and willing to use Ingress instead of IP address on the Host Network.
The API and the image server will be accessible via the hostname specified in the ingress configuration.
This should only be used with virtual media deployments.
Cannot be set at the same time with networking.externalIP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>interface</b></td>
        <td>string</td>
        <td>
          Interface is a Linux network device to listen on.
Detected from IPAddress if missing.
Requires DisableHostNetwork to be false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ipAddress</b></td>
        <td>string</td>
        <td>
          IPAddress is the main IP address to listen on and use for communication.
Detected from Interface if missing. Cannot be provided for a highly available architecture.
Requires DisableHostNetwork to be false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ipAddressManager</b></td>
        <td>enum</td>
        <td>
          Configures the way the provided IP address will be managed on the provided interface.
By default, the IP address is expected to be already present.
Use "keepalived" to start a Keepalived container managing the IP address.
Warning: keepalived is not compatible with the highly available architecture.

Deprecated: Use the keepalived field instead.<br/>
          <br/>
            <i>Enum</i>: , keepalived<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecnetworkingkeepalived">keepalived</a></b></td>
        <td>object</td>
        <td>
          Keepalived configures Keepalived to manage virtual IPs on the specified interfaces.
When enabled, a Keepalived container will be started to manage the main ipAddress
on the main interface, plus any additional VIPs listed in additionalVIPs.
Cannot be used together with ipAddressManager.
Requires DisableHostNetwork to be false.
Warning: keepalived is not compatible with the highly available architecture.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>macAddresses</b></td>
        <td>[]string</td>
        <td>
          MACAddresses can be provided to make the start script pick the interface matching any of these addresses.
Only set if no other options can be used.
Requires DisableHostNetwork to be false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>prometheusExporterPort</b></td>
        <td>integer</td>
        <td>
          PrometheusExporterPort is the port used for the Ironic Prometheus Exporter metrics endpoint.
Only used when spec.prometheusExporter.enabled is true.<br/>
          <br/>
            <i>Format</i>: int32<br/>
            <i>Default</i>: 9608<br/>
            <i>Minimum</i>: 1<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>rpcPort</b></td>
        <td>integer</td>
        <td>
          RPCPort is the internal RPC port used for Ironic.
Only change this if the default value causes a conflict on your deployment.<br/>
          <br/>
            <i>Format</i>: int32<br/>
            <i>Default</i>: 6189<br/>
            <i>Minimum</i>: 1<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.networking.dhcp
<sup><sup>[↩ Parent](#ironicspecnetworking)</sup></sup>



DHCP is a configuration of DHCP for the network boot service (dnsmasq).
The service is only deployed when this is set.
This setting is currently incompatible with the highly available architecture.
Requires DisableHostNetwork to be false.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>dnsAddress</b></td>
        <td>string</td>
        <td>
          DNSAddress is the IP address of the DNS server to pass to hosts via DHCP.
Must not be set together with ServeDNS.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>gatewayAddress</b></td>
        <td>string</td>
        <td>
          GatewayAddress is the IP address of the gateway to pass to hosts via DHCP.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hosts</b></td>
        <td>[]string</td>
        <td>
          Hosts is a set of DHCP host records to pass to dnsmasq.
Check the dnsmasq documentation on dhcp-host for an explanation of the format.
There is no API-side validation. Most users will leave this unset.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ignore</b></td>
        <td>[]string</td>
        <td>
          Ignore is set of dnsmasq tags to ignore and not provide any DHCP.
Check the dnsmasq documentation on dhcp-ignore for an explanation of the format.
There is no API-side validation. Most users will leave this unset.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>networkCIDR</b></td>
        <td>string</td>
        <td>
          NetworkCIDR is a CIDR of the provisioning network. Required.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>rangeBegin</b></td>
        <td>string</td>
        <td>
          RangeBegin is the first IP that can be given to hosts. Must be inside NetworkCIDR.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>rangeEnd</b></td>
        <td>string</td>
        <td>
          RangeEnd is the last IP that can be given to hosts. Must be inside NetworkCIDR.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>serveDNS</b></td>
        <td>boolean</td>
        <td>
          ServeDNS is set to true to pass the provisioning host as the DNS server on the provisioning network.
Must not be set together with DNSAddress.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.networking.ingress
<sup><sup>[↩ Parent](#ironicspecnetworking)</sup></sup>



Configure Ingress resource for Ironic services.
Set this option when you are planning to deploy Ironic in a public cluster and willing to use Ingress instead of IP address on the Host Network.
The API and the image server will be accessible via the hostname specified in the ingress configuration.
This should only be used with virtual media deployments.
Cannot be set at the same time with networking.externalIP.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>annotations</b></td>
        <td>map[string]string</td>
        <td>
          Annotations to be added to Ingress resource<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host is the fully qualified domain name of a network host.
This defines the hostname that the Ingress resource will route traffic for.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ingressClassName</b></td>
        <td>string</td>
        <td>
          IngressClass of Ingress resource<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.networking.keepalived
<sup><sup>[↩ Parent](#ironicspecnetworking)</sup></sup>



Keepalived configures Keepalived to manage virtual IPs on the specified interfaces.
When enabled, a Keepalived container will be started to manage the main ipAddress
on the main interface, plus any additional VIPs listed in additionalVIPs.
Cannot be used together with ipAddressManager.
Requires DisableHostNetwork to be false.
Warning: keepalived is not compatible with the highly available architecture.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Enabled indicates whether Keepalived should be started to manage the virtual IP.
When enabled, the main ipAddress and interface from the networking configuration
are always included automatically.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#ironicspecnetworkingkeepalivedadditionalvipsindex">additionalVIPs</a></b></td>
        <td>[]object</td>
        <td>
          AdditionalVIPs is a list of additional virtual IPs to be managed by Keepalived,
beyond the main ipAddress/interface from the networking configuration.
Use this when you need Keepalived to manage IPs on additional network interfaces.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.networking.keepalived.additionalVIPs[index]
<sup><sup>[↩ Parent](#ironicspecnetworkingkeepalived)</sup></sup>



KeepalivedIP defines a virtual IP address to be managed by Keepalived.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>interface</b></td>
        <td>string</td>
        <td>
          Interface is the Linux network interface on which to manage the IP.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>ipAddress</b></td>
        <td>string</td>
        <td>
          IPAddress is the virtual IP address to manage.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>prefix</b></td>
        <td>integer</td>
        <td>
          Prefix is the prefix length of the IP address (e.g. 24 for a /24 subnet).
When not set, keepalived uses /32 for IPv4 and /128 for IPv6.<br/>
          <br/>
            <i>Format</i>: int32<br/>
            <i>Minimum</i>: 0<br/>
            <i>Maximum</i>: 128<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.networkingService
<sup><sup>[↩ Parent](#ironicspec)</sup></sup>



NetworkingService provides configuration for the Ironic Networking Service

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Enabled enables the Ironic Networking Service integration<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecnetworkingserviceprovidernetworksindex">providerNetworks</a></b></td>
        <td>[]object</td>
        <td>
          ProviderNetworks defines the provider network configurations for Ironic<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>switchConfigSecretName</b></td>
        <td>string</td>
        <td>
          SwitchConfigSecretName optionally specifies the name of the secret containing
switch configuration. If not specified, defaults to "<ironic-name>-switch-config".
The secret must have the environment label to be used by the operator.
The operator creates an empty secret if it does not exist, and the Baremetal
Operator populates it from BareMetalSwitch CRDs.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>switchCredentialsSecretName</b></td>
        <td>string</td>
        <td>
          SwitchCredentialsSecretName optionally specifies the name of the secret containing
additional switch credentials (e.g., SSH private keys for switches using publickey
authentication). If not specified, defaults to "<ironic-name>-switch-credentials".
The secret must have the environment label to be used by the operator.
The operator creates an empty secret if it does not exist, and the Baremetal
Operator populates it from BareMetalSwitch CRDs.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.networkingService.providerNetworks[index]
<sup><sup>[↩ Parent](#ironicspecnetworkingservice)</sup></sup>



ProviderNetworkConfig defines the network configuration for Ironic service operations.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>mode</b></td>
        <td>enum</td>
        <td>
          Mode specifies the switch port mode for service operations<br/>
          <br/>
            <i>Enum</i>: access, trunk, hybrid<br/>
            <i>Default</i>: access<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>nativeVLAN</b></td>
        <td>integer</td>
        <td>
          NativeVLAN specifies the native VLAN ID for service operations<br/>
          <br/>
            <i>Format</i>: int32<br/>
            <i>Minimum</i>: 1<br/>
            <i>Maximum</i>: 4094<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>enum</td>
        <td>
          Type specifies which provider network this configuration applies to.<br/>
          <br/>
            <i>Enum</i>: idle, inspection, cleaning, rescuing, servicing, provisioning<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>allowedVLANs</b></td>
        <td>[]string</td>
        <td>
          AllowedVLANs specifies the list of allowed VLANs for trunk/hybrid modes.
Each entry can be a single VLAN ID (e.g., "100") or a range (e.g., "100-200").<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mtu</b></td>
        <td>integer</td>
        <td>
          MTU defines the MTU to be applied to the switch port configuration.<br/>
          <br/>
            <i>Format</i>: int32<br/>
            <i>Minimum</i>: 1280<br/>
            <i>Maximum</i>: 9216<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides
<sup><sup>[↩ Parent](#ironicspec)</sup></sup>



Overrides for the generated Deployment or Daemon Set.
EXPERIMENTAL: requires feature gate Overrides.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesagentimagesindex">agentImages</a></b></td>
        <td>[]object</td>
        <td>
          AgentImages overrides the default IPA images with custom per-architecture images.
Consider setting deployRamdisk.disableDownloader=true when using custom images.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>annotations</b></td>
        <td>map[string]string</td>
        <td>
          Extra annotations to add to each pod (including upgrade jobs).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindex">containers</a></b></td>
        <td>[]object</td>
        <td>
          Containers to append to the main Ironic pod.
If a container name matches an existing container, the existing container is replaced.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverrideshttpdlivenessprobe">httpdLivenessProbe</a></b></td>
        <td>object</td>
        <td>
          HttpdLivenessProbe overrides the httpd container liveness probe.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverrideshttpdreadinessprobe">httpdReadinessProbe</a></b></td>
        <td>object</td>
        <td>
          HttpdReadinessProbe overrides the httpd container readiness probe.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindex">initContainers</a></b></td>
        <td>[]object</td>
        <td>
          InitContainers to append to the main Ironic pod.
If a container name matches an existing init container, the existing init container is replaced.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>labels</b></td>
        <td>map[string]string</td>
        <td>
          Extra labels to add to each pod (including upgrade jobs).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindex">volumes</a></b></td>
        <td>[]object</td>
        <td>
          Volumes to add to the main Ironic pod (and upgrade jobs).
Use this together with volumeMounts on overridden containers to mount
additional ConfigMaps, Secrets or PersistentVolumeClaims.
If a volume name matches an existing volume, the existing volume is replaced.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.agentImages[index]
<sup><sup>[↩ Parent](#ironicspecoverrides)</sup></sup>



AgentImages defines a single IPA (Ironic Python Agent) image configuration.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>initramfs</b></td>
        <td>string</td>
        <td>
          Initramfs is the URL of the IPA initramfs/ramdisk image.
Supported schemes: file://, http://, https://, oci://.
file:// URLs must use absolute paths (e.g. "file:///shared/html/images/ironic-python-agent.initramfs").<br/>
          <br/>
            <i>Format</i>: uri<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>kernel</b></td>
        <td>string</td>
        <td>
          Kernel is the URL of the IPA kernel image.
Supported schemes: file://, http://, https://, oci://.
file:// URLs must use absolute paths (e.g. "file:///shared/html/images/ironic-python-agent.kernel").<br/>
          <br/>
            <i>Format</i>: uri<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>architecture</b></td>
        <td>enum</td>
        <td>
          Architecture is the target CPU architecture.
When empty, sets the default DEPLOY_KERNEL_URL/DEPLOY_RAMDISK_URL.<br/>
          <br/>
            <i>Enum</i>: , x86_64, aarch64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index]
<sup><sup>[↩ Parent](#ironicspecoverrides)</sup></sup>



A single application container that you want to run within a pod.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the container specified as a DNS_LABEL.
Each container in a pod must have a unique name (DNS_LABEL).
Cannot be updated.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>args</b></td>
        <td>[]string</td>
        <td>
          Arguments to the entrypoint.
The container image's CMD is used if this is not provided.
Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced
to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will
produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless
of whether the variable exists or not. Cannot be updated.
More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Entrypoint array. Not executed within a shell.
The container image's ENTRYPOINT is used if this is not provided.
Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced
to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will
produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless
of whether the variable exists or not. Cannot be updated.
More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexenvindex">env</a></b></td>
        <td>[]object</td>
        <td>
          List of environment variables to set in the container.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexenvfromindex">envFrom</a></b></td>
        <td>[]object</td>
        <td>
          List of sources to populate environment variables in the container.
The keys defined within a source may consist of any printable ASCII characters except '='.
When a key exists in multiple
sources, the value associated with the last source will take precedence.
Values defined by an Env with a duplicate key will take precedence.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>image</b></td>
        <td>string</td>
        <td>
          Container image name.
More info: https://kubernetes.io/docs/concepts/containers/images
This field is optional to allow higher level config management to default or override
container images in workload controllers like Deployments and StatefulSets.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>imagePullPolicy</b></td>
        <td>string</td>
        <td>
          Image pull policy.
One of Always, Never, IfNotPresent.
Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
Cannot be updated.
More info: https://kubernetes.io/docs/concepts/containers/images#updating-images<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlifecycle">lifecycle</a></b></td>
        <td>object</td>
        <td>
          Actions that the management system should take in response to container lifecycle events.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlivenessprobe">livenessProbe</a></b></td>
        <td>object</td>
        <td>
          Periodic probe of container liveness.
Container will be restarted if the probe fails.
Cannot be updated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexportsindex">ports</a></b></td>
        <td>[]object</td>
        <td>
          List of ports to expose from the container. Not specifying a port here
DOES NOT prevent that port from being exposed. Any port which is
listening on the default "0.0.0.0" address inside a container will be
accessible from the network.
Modifying this array with strategic merge patch may corrupt the data.
For more information See https://github.com/kubernetes/kubernetes/issues/108255.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexreadinessprobe">readinessProbe</a></b></td>
        <td>object</td>
        <td>
          Periodic probe of container service readiness.
Container will be removed from service endpoints if the probe fails.
Cannot be updated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexresizepolicyindex">resizePolicy</a></b></td>
        <td>[]object</td>
        <td>
          Resources resize policy for the container.
This field cannot be set on ephemeral containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexresources">resources</a></b></td>
        <td>object</td>
        <td>
          Compute Resources required by this container.
Cannot be updated.
More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>restartPolicy</b></td>
        <td>string</td>
        <td>
          RestartPolicy defines the restart behavior of individual containers in a pod.
This overrides the pod-level restart policy. When this field is not specified,
the restart behavior is defined by the Pod's restart policy and the container type.
Additionally, setting the RestartPolicy as "Always" for the init container will
have the following effect:
this init container will be continually restarted on
exit until all regular containers have terminated. Once all regular
containers have completed, all init containers with restartPolicy "Always"
will be shut down. This lifecycle differs from normal init containers and
is often referred to as a "sidecar" container. Although this init
container still starts in the init container sequence, it does not wait
for the container to complete before proceeding to the next init
container. Instead, the next init container starts immediately after this
init container is started, or after any startupProbe has successfully
completed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexrestartpolicyrulesindex">restartPolicyRules</a></b></td>
        <td>[]object</td>
        <td>
          Represents a list of rules to be checked to determine if the
container should be restarted on exit. The rules are evaluated in
order. Once a rule matches a container exit condition, the remaining
rules are ignored. If no rule matches the container exit condition,
the Container-level restart policy determines the whether the container
is restarted or not. Constraints on the rules:
- At most 20 rules are allowed.
- Rules can have the same action.
- Identical rules are not forbidden in validations.
When rules are specified, container MUST set RestartPolicy explicitly
even it if matches the Pod's RestartPolicy.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexsecuritycontext">securityContext</a></b></td>
        <td>object</td>
        <td>
          SecurityContext defines the security options the container should be run with.
If set, the fields of SecurityContext override the equivalent fields of PodSecurityContext.
More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexstartupprobe">startupProbe</a></b></td>
        <td>object</td>
        <td>
          StartupProbe indicates that the Pod has successfully initialized.
If specified, no other probes are executed until this completes successfully.
If this probe fails, the Pod will be restarted, just as if the livenessProbe failed.
This can be used to provide different probe parameters at the beginning of a Pod's lifecycle,
when it might take a long time to load data or warm a cache, than during steady-state operation.
This cannot be updated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>stdin</b></td>
        <td>boolean</td>
        <td>
          Whether this container should allocate a buffer for stdin in the container runtime. If this
is not set, reads from stdin in the container will always result in EOF.
Default is false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>stdinOnce</b></td>
        <td>boolean</td>
        <td>
          Whether the container runtime should close the stdin channel after it has been opened by
a single attach. When stdin is true the stdin stream will remain open across multiple attach
sessions. If stdinOnce is set to true, stdin is opened on container start, is empty until the
first client attaches to stdin, and then remains open and accepts data until the client disconnects,
at which time stdin is closed and remains closed until the container is restarted. If this
flag is false, a container processes that reads from stdin will never receive an EOF.
Default is false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationMessagePath</b></td>
        <td>string</td>
        <td>
          Optional: Path at which the file to which the container's termination message
will be written is mounted into the container's filesystem.
Message written is intended to be brief final status, such as an assertion failure message.
Will be truncated by the node if greater than 4096 bytes. The total message length across
all containers will be limited to 12kb.
Defaults to /dev/termination-log.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationMessagePolicy</b></td>
        <td>string</td>
        <td>
          Indicate how the termination message should be populated. File will use the contents of
terminationMessagePath to populate the container status message on both success and failure.
FallbackToLogsOnError will use the last chunk of container log output if the termination
message file is empty and the container exited with an error.
The log output is limited to 2048 bytes or 80 lines, whichever is smaller.
Defaults to File.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tty</b></td>
        <td>boolean</td>
        <td>
          Whether this container should allocate a TTY for itself, also requires 'stdin' to be true.
Default is false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexvolumedevicesindex">volumeDevices</a></b></td>
        <td>[]object</td>
        <td>
          volumeDevices is the list of block devices to be used by the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexvolumemountsindex">volumeMounts</a></b></td>
        <td>[]object</td>
        <td>
          Pod volumes to mount into the container's filesystem.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>workingDir</b></td>
        <td>string</td>
        <td>
          Container's working directory.
If not specified, the container runtime's default will be used, which
might be configured in the container image.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].env[index]
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindex)</sup></sup>



EnvVar represents an environment variable present in a Container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the environment variable.
May consist of any printable ASCII characters except '='.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Variable references $(VAR_NAME) are expanded
using the previously defined environment variables in the container and
any service environment variables. If a variable cannot be resolved,
the reference in the input string will be unchanged. Double $$ are reduced
to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e.
"$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)".
Escaped references will never be expanded, regardless of whether the variable
exists or not.
Defaults to "".<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexenvindexvaluefrom">valueFrom</a></b></td>
        <td>object</td>
        <td>
          Source for the environment variable's value. Cannot be used if value is not empty.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].env[index].valueFrom
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexenvindex)</sup></sup>



Source for the environment variable's value. Cannot be used if value is not empty.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridescontainersindexenvindexvaluefromconfigmapkeyref">configMapKeyRef</a></b></td>
        <td>object</td>
        <td>
          Selects a key of a ConfigMap.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexenvindexvaluefromfieldref">fieldRef</a></b></td>
        <td>object</td>
        <td>
          Selects a field of the pod: supports metadata.name, metadata.namespace, `metadata.labels['<KEY>']`, `metadata.annotations['<KEY>']`,
spec.nodeName, spec.serviceAccountName, status.hostIP, status.podIP, status.podIPs.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexenvindexvaluefromfilekeyref">fileKeyRef</a></b></td>
        <td>object</td>
        <td>
          FileKeyRef selects a key of the env file.
Requires the EnvFiles feature gate to be enabled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexenvindexvaluefromresourcefieldref">resourceFieldRef</a></b></td>
        <td>object</td>
        <td>
          Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.memory, limits.ephemeral-storage, requests.cpu, requests.memory and requests.ephemeral-storage) are currently supported.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexenvindexvaluefromsecretkeyref">secretKeyRef</a></b></td>
        <td>object</td>
        <td>
          Selects a key of a secret in the pod's namespace<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].env[index].valueFrom.configMapKeyRef
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexenvindexvaluefrom)</sup></sup>



Selects a key of a ConfigMap.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The key to select.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the ConfigMap or its key must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].env[index].valueFrom.fieldRef
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexenvindexvaluefrom)</sup></sup>



Selects a field of the pod: supports metadata.name, metadata.namespace, `metadata.labels['<KEY>']`, `metadata.annotations['<KEY>']`,
spec.nodeName, spec.serviceAccountName, status.hostIP, status.podIP, status.podIPs.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>fieldPath</b></td>
        <td>string</td>
        <td>
          Path of the field to select in the specified API version.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>apiVersion</b></td>
        <td>string</td>
        <td>
          Version of the schema the FieldPath is written in terms of, defaults to "v1".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].env[index].valueFrom.fileKeyRef
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexenvindexvaluefrom)</sup></sup>



FileKeyRef selects a key of the env file.
Requires the EnvFiles feature gate to be enabled.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The key within the env file. An invalid key will prevent the pod from starting.
The keys defined within a source may consist of any printable ASCII characters except '='.
During Alpha stage of the EnvFiles feature gate, the key size is limited to 128 characters.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          The path within the volume from which to select the file.
Must be relative and may not contain the '..' path or start with '..'.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>volumeName</b></td>
        <td>string</td>
        <td>
          The name of the volume mount containing the env file.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the file or its key must be defined. If the file or key
does not exist, then the env var is not published.
If optional is set to true and the specified key does not exist,
the environment variable will not be set in the Pod's containers.

If optional is set to false and the specified key does not exist,
an error will be returned during Pod creation.<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].env[index].valueFrom.resourceFieldRef
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexenvindexvaluefrom)</sup></sup>



Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.memory, limits.ephemeral-storage, requests.cpu, requests.memory and requests.ephemeral-storage) are currently supported.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resource</b></td>
        <td>string</td>
        <td>
          Required: resource to select<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>containerName</b></td>
        <td>string</td>
        <td>
          Container name: required for volumes, optional for env vars<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>divisor</b></td>
        <td>int or string</td>
        <td>
          Specifies the output format of the exposed resources, defaults to "1"<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].env[index].valueFrom.secretKeyRef
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexenvindexvaluefrom)</sup></sup>



Selects a key of a secret in the pod's namespace

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The key of the secret to select from.  Must be a valid secret key.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the Secret or its key must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].envFrom[index]
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindex)</sup></sup>



EnvFromSource represents the source of a set of ConfigMaps or Secrets

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridescontainersindexenvfromindexconfigmapref">configMapRef</a></b></td>
        <td>object</td>
        <td>
          The ConfigMap to select from<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>prefix</b></td>
        <td>string</td>
        <td>
          Optional text to prepend to the name of each environment variable.
May consist of any printable ASCII characters except '='.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexenvfromindexsecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          The Secret to select from<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].envFrom[index].configMapRef
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexenvfromindex)</sup></sup>



The ConfigMap to select from

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the ConfigMap must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].envFrom[index].secretRef
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexenvfromindex)</sup></sup>



The Secret to select from

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the Secret must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].lifecycle
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindex)</sup></sup>



Actions that the management system should take in response to container lifecycle events.
Cannot be updated.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlifecyclepoststart">postStart</a></b></td>
        <td>object</td>
        <td>
          PostStart is called immediately after a container is created. If the handler fails,
the container is terminated and restarted according to its restart policy.
Other management of the container blocks until the hook completes.
More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlifecycleprestop">preStop</a></b></td>
        <td>object</td>
        <td>
          PreStop is called immediately before a container is terminated due to an
API request or management event such as liveness/startup probe failure,
preemption, resource contention, etc. The handler is not called if the
container crashes or exits. The Pod's termination grace period countdown begins before the
PreStop hook is executed. Regardless of the outcome of the handler, the
container will eventually terminate within the Pod's termination grace
period (unless delayed by finalizers). Other management of the container blocks until the hook completes
or until the termination grace period is reached.
More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>stopSignal</b></td>
        <td>string</td>
        <td>
          StopSignal defines which signal will be sent to a container when it is being stopped.
If not specified, the default is defined by the container runtime in use.
StopSignal can only be set for Pods with a non-empty .spec.os.name<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].lifecycle.postStart
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlifecycle)</sup></sup>



PostStart is called immediately after a container is created. If the handler fails,
the container is terminated and restarted according to its restart policy.
Other management of the container blocks until the hook completes.
More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlifecyclepoststartexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlifecyclepoststarthttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlifecyclepoststartsleep">sleep</a></b></td>
        <td>object</td>
        <td>
          Sleep represents a duration that the container should sleep.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlifecyclepoststarttcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility. There is no validation of this field and
lifecycle hooks will fail at runtime when it is specified.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].lifecycle.postStart.exec
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlifecyclepoststart)</sup></sup>



Exec specifies a command to execute in the container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command  is root ('/') in the container's filesystem. The command is simply exec'd, it is
not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use
a shell, you need to explicitly call out to that shell.
Exit status of 0 is treated as live/healthy and non-zero is unhealthy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].lifecycle.postStart.httpGet
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlifecyclepoststart)</sup></sup>



HTTPGet specifies an HTTP GET request to perform.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP. You probably want to set
"Host" in httpHeaders instead.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlifecyclepoststarthttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].lifecycle.postStart.httpGet.httpHeaders[index]
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlifecyclepoststarthttpget)</sup></sup>



HTTPHeader describes a custom header to be used in HTTP probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The header field name.
This will be canonicalized upon output, so case-variant names will be understood as the same header.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].lifecycle.postStart.sleep
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlifecyclepoststart)</sup></sup>



Sleep represents a duration that the container should sleep.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>seconds</b></td>
        <td>integer</td>
        <td>
          Seconds is the number of seconds to sleep.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].lifecycle.postStart.tcpSocket
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlifecyclepoststart)</sup></sup>



Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility. There is no validation of this field and
lifecycle hooks will fail at runtime when it is specified.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].lifecycle.preStop
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlifecycle)</sup></sup>



PreStop is called immediately before a container is terminated due to an
API request or management event such as liveness/startup probe failure,
preemption, resource contention, etc. The handler is not called if the
container crashes or exits. The Pod's termination grace period countdown begins before the
PreStop hook is executed. Regardless of the outcome of the handler, the
container will eventually terminate within the Pod's termination grace
period (unless delayed by finalizers). Other management of the container blocks until the hook completes
or until the termination grace period is reached.
More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlifecycleprestopexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlifecycleprestophttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlifecycleprestopsleep">sleep</a></b></td>
        <td>object</td>
        <td>
          Sleep represents a duration that the container should sleep.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlifecycleprestoptcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility. There is no validation of this field and
lifecycle hooks will fail at runtime when it is specified.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].lifecycle.preStop.exec
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlifecycleprestop)</sup></sup>



Exec specifies a command to execute in the container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command  is root ('/') in the container's filesystem. The command is simply exec'd, it is
not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use
a shell, you need to explicitly call out to that shell.
Exit status of 0 is treated as live/healthy and non-zero is unhealthy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].lifecycle.preStop.httpGet
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlifecycleprestop)</sup></sup>



HTTPGet specifies an HTTP GET request to perform.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP. You probably want to set
"Host" in httpHeaders instead.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlifecycleprestophttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].lifecycle.preStop.httpGet.httpHeaders[index]
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlifecycleprestophttpget)</sup></sup>



HTTPHeader describes a custom header to be used in HTTP probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The header field name.
This will be canonicalized upon output, so case-variant names will be understood as the same header.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].lifecycle.preStop.sleep
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlifecycleprestop)</sup></sup>



Sleep represents a duration that the container should sleep.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>seconds</b></td>
        <td>integer</td>
        <td>
          Seconds is the number of seconds to sleep.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].lifecycle.preStop.tcpSocket
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlifecycleprestop)</sup></sup>



Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility. There is no validation of this field and
lifecycle hooks will fail at runtime when it is specified.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].livenessProbe
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindex)</sup></sup>



Periodic probe of container liveness.
Container will be restarted if the probe fails.
Cannot be updated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlivenessprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlivenessprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlivenessprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlivenessprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.
The grace period is the duration in seconds after the processes running in the pod are sent
a termination signal and the time when the processes are forcibly halted with a kill signal.
Set this value longer than the expected cleanup time for your process.
If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this
value overrides the value provided by the pod spec.
Value must be non-negative integer. The value zero indicates stop immediately via
the kill signal (no opportunity to shut down).
This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate.
Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].livenessProbe.exec
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlivenessprobe)</sup></sup>



Exec specifies a command to execute in the container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command  is root ('/') in the container's filesystem. The command is simply exec'd, it is
not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use
a shell, you need to explicitly call out to that shell.
Exit status of 0 is treated as live/healthy and non-zero is unhealthy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].livenessProbe.grpc
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlivenessprobe)</sup></sup>



GRPC specifies a GRPC HealthCheckRequest.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).

If this is not specified, the default behavior is defined by gRPC.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].livenessProbe.httpGet
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlivenessprobe)</sup></sup>



HTTPGet specifies an HTTP GET request to perform.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP. You probably want to set
"Host" in httpHeaders instead.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexlivenessprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].livenessProbe.httpGet.httpHeaders[index]
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlivenessprobehttpget)</sup></sup>



HTTPHeader describes a custom header to be used in HTTP probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The header field name.
This will be canonicalized upon output, so case-variant names will be understood as the same header.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].livenessProbe.tcpSocket
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexlivenessprobe)</sup></sup>



TCPSocket specifies a connection to a TCP port.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].ports[index]
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindex)</sup></sup>



ContainerPort represents a network port in a single container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>containerPort</b></td>
        <td>integer</td>
        <td>
          Number of port to expose on the pod's IP address.
This must be a valid port number, 0 < x < 65536.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>hostIP</b></td>
        <td>string</td>
        <td>
          What host IP to bind the external port to.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostPort</b></td>
        <td>integer</td>
        <td>
          Number of port to expose on the host.
If specified, this must be a valid port number, 0 < x < 65536.
If HostNetwork is specified, this must match ContainerPort.
Most containers do not need this.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          If specified, this must be an IANA_SVC_NAME and unique within the pod. Each
named port in a pod must have a unique name. Name for the port that can be
referred to by services.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>protocol</b></td>
        <td>string</td>
        <td>
          Protocol for port. Must be UDP, TCP, or SCTP.
Defaults to "TCP".<br/>
          <br/>
            <i>Default</i>: TCP<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].readinessProbe
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindex)</sup></sup>



Periodic probe of container service readiness.
Container will be removed from service endpoints if the probe fails.
Cannot be updated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridescontainersindexreadinessprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexreadinessprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexreadinessprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexreadinessprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.
The grace period is the duration in seconds after the processes running in the pod are sent
a termination signal and the time when the processes are forcibly halted with a kill signal.
Set this value longer than the expected cleanup time for your process.
If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this
value overrides the value provided by the pod spec.
Value must be non-negative integer. The value zero indicates stop immediately via
the kill signal (no opportunity to shut down).
This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate.
Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].readinessProbe.exec
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexreadinessprobe)</sup></sup>



Exec specifies a command to execute in the container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command  is root ('/') in the container's filesystem. The command is simply exec'd, it is
not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use
a shell, you need to explicitly call out to that shell.
Exit status of 0 is treated as live/healthy and non-zero is unhealthy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].readinessProbe.grpc
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexreadinessprobe)</sup></sup>



GRPC specifies a GRPC HealthCheckRequest.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).

If this is not specified, the default behavior is defined by gRPC.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].readinessProbe.httpGet
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexreadinessprobe)</sup></sup>



HTTPGet specifies an HTTP GET request to perform.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP. You probably want to set
"Host" in httpHeaders instead.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexreadinessprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].readinessProbe.httpGet.httpHeaders[index]
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexreadinessprobehttpget)</sup></sup>



HTTPHeader describes a custom header to be used in HTTP probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The header field name.
This will be canonicalized upon output, so case-variant names will be understood as the same header.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].readinessProbe.tcpSocket
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexreadinessprobe)</sup></sup>



TCPSocket specifies a connection to a TCP port.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].resizePolicy[index]
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindex)</sup></sup>



ContainerResizePolicy represents resource resize policy for the container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resourceName</b></td>
        <td>string</td>
        <td>
          Name of the resource to which this resource resize policy applies.
Supported values: cpu, memory.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>restartPolicy</b></td>
        <td>string</td>
        <td>
          Restart policy to apply when specified resource is resized.
If not specified, it defaults to NotRequired.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].resources
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindex)</sup></sup>



Compute Resources required by this container.
Cannot be updated.
More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridescontainersindexresourcesclaimsindex">claims</a></b></td>
        <td>[]object</td>
        <td>
          Claims lists the names of resources, defined in spec.resourceClaims,
that are used by this container.

This field depends on the
DynamicResourceAllocation feature gate.

This field is immutable. It can only be set for containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>limits</b></td>
        <td>map[string]int or string</td>
        <td>
          Limits describes the maximum amount of compute resources allowed.
More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>requests</b></td>
        <td>map[string]int or string</td>
        <td>
          Requests describes the minimum amount of compute resources required.
If Requests is omitted for a container, it defaults to Limits if that is explicitly specified,
otherwise to an implementation-defined value. Requests cannot exceed Limits.
More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].resources.claims[index]
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexresources)</sup></sup>



ResourceClaim references one entry in PodSpec.ResourceClaims.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name must match the name of one entry in pod.spec.resourceClaims of
the Pod where this field is used. It makes that resource available
inside a container.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>request</b></td>
        <td>string</td>
        <td>
          Request is the name chosen for a request in the referenced claim.
If empty, everything from the claim is made available, otherwise
only the result of this request.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].restartPolicyRules[index]
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindex)</sup></sup>



ContainerRestartRule describes how a container exit is handled.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>action</b></td>
        <td>string</td>
        <td>
          Specifies the action taken on a container exit if the requirements
are satisfied. The only possible value is "Restart" to restart the
container.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexrestartpolicyrulesindexexitcodes">exitCodes</a></b></td>
        <td>object</td>
        <td>
          Represents the exit codes to check on container exits.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].restartPolicyRules[index].exitCodes
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexrestartpolicyrulesindex)</sup></sup>



Represents the exit codes to check on container exits.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          Represents the relationship between the container exit code(s) and the
specified values. Possible values are:
- In: the requirement is satisfied if the container exit code is in the
  set of specified values.
- NotIn: the requirement is satisfied if the container exit code is
  not in the set of specified values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]integer</td>
        <td>
          Specifies the set of values to check for container exit codes.
At most 255 elements are allowed.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].securityContext
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindex)</sup></sup>



SecurityContext defines the security options the container should be run with.
If set, the fields of SecurityContext override the equivalent fields of PodSecurityContext.
More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>allowPrivilegeEscalation</b></td>
        <td>boolean</td>
        <td>
          AllowPrivilegeEscalation controls whether a process can gain more
privileges than its parent process. This bool directly controls if
the no_new_privs flag will be set on the container process.
AllowPrivilegeEscalation is true always when the container is:
1) run as Privileged
2) has CAP_SYS_ADMIN
Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexsecuritycontextapparmorprofile">appArmorProfile</a></b></td>
        <td>object</td>
        <td>
          appArmorProfile is the AppArmor options to use by this container. If set, this profile
overrides the pod's appArmorProfile.
Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexsecuritycontextcapabilities">capabilities</a></b></td>
        <td>object</td>
        <td>
          The capabilities to add/drop when running containers.
Defaults to the default set of capabilities granted by the container runtime.
Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>privileged</b></td>
        <td>boolean</td>
        <td>
          Run container in privileged mode.
Processes in privileged containers are essentially equivalent to root on the host.
Defaults to false.
Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>procMount</b></td>
        <td>string</td>
        <td>
          procMount denotes the type of proc mount to use for the containers.
The default value is Default which uses the container runtime defaults for
readonly paths and masked paths.
Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnlyRootFilesystem</b></td>
        <td>boolean</td>
        <td>
          Whether this container has a read-only root filesystem.
Default is false.
Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsGroup</b></td>
        <td>integer</td>
        <td>
          The GID to run the entrypoint of the container process.
Uses runtime default if unset.
May also be set in PodSecurityContext.  If set in both SecurityContext and
PodSecurityContext, the value specified in SecurityContext takes precedence.
Note that this field cannot be set when spec.os.name is windows.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsNonRoot</b></td>
        <td>boolean</td>
        <td>
          Indicates that the container must run as a non-root user.
If true, the Kubelet will validate the image at runtime to ensure that it
does not run as UID 0 (root) and fail to start the container if it does.
If unset or false, no such validation will be performed.
May also be set in PodSecurityContext.  If set in both SecurityContext and
PodSecurityContext, the value specified in SecurityContext takes precedence.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsUser</b></td>
        <td>integer</td>
        <td>
          The UID to run the entrypoint of the container process.
Defaults to user specified in image metadata if unspecified.
May also be set in PodSecurityContext.  If set in both SecurityContext and
PodSecurityContext, the value specified in SecurityContext takes precedence.
Note that this field cannot be set when spec.os.name is windows.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexsecuritycontextselinuxoptions">seLinuxOptions</a></b></td>
        <td>object</td>
        <td>
          The SELinux context to be applied to the container.
If unspecified, the container runtime will allocate a random SELinux context for each
container.  May also be set in PodSecurityContext.  If set in both SecurityContext and
PodSecurityContext, the value specified in SecurityContext takes precedence.
Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexsecuritycontextseccompprofile">seccompProfile</a></b></td>
        <td>object</td>
        <td>
          The seccomp options to use by this container. If seccomp options are
provided at both the pod & container level, the container options
override the pod options.
Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexsecuritycontextwindowsoptions">windowsOptions</a></b></td>
        <td>object</td>
        <td>
          The Windows specific settings applied to all containers.
If unspecified, the options from the PodSecurityContext will be used.
If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
Note that this field cannot be set when spec.os.name is linux.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].securityContext.appArmorProfile
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexsecuritycontext)</sup></sup>



appArmorProfile is the AppArmor options to use by this container. If set, this profile
overrides the pod's appArmorProfile.
Note that this field cannot be set when spec.os.name is windows.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          type indicates which kind of AppArmor profile will be applied.
Valid options are:
  Localhost - a profile pre-loaded on the node.
  RuntimeDefault - the container runtime's default profile.
  Unconfined - no AppArmor enforcement.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>localhostProfile</b></td>
        <td>string</td>
        <td>
          localhostProfile indicates a profile loaded on the node that should be used.
The profile must be preconfigured on the node to work.
Must match the loaded name of the profile.
Must be set if and only if type is "Localhost".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].securityContext.capabilities
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexsecuritycontext)</sup></sup>



The capabilities to add/drop when running containers.
Defaults to the default set of capabilities granted by the container runtime.
Note that this field cannot be set when spec.os.name is windows.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>add</b></td>
        <td>[]string</td>
        <td>
          Added capabilities<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>drop</b></td>
        <td>[]string</td>
        <td>
          Removed capabilities<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].securityContext.seLinuxOptions
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexsecuritycontext)</sup></sup>



The SELinux context to be applied to the container.
If unspecified, the container runtime will allocate a random SELinux context for each
container.  May also be set in PodSecurityContext.  If set in both SecurityContext and
PodSecurityContext, the value specified in SecurityContext takes precedence.
Note that this field cannot be set when spec.os.name is windows.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>level</b></td>
        <td>string</td>
        <td>
          Level is SELinux level label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>role</b></td>
        <td>string</td>
        <td>
          Role is a SELinux role label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type is a SELinux type label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>user</b></td>
        <td>string</td>
        <td>
          User is a SELinux user label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].securityContext.seccompProfile
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexsecuritycontext)</sup></sup>



The seccomp options to use by this container. If seccomp options are
provided at both the pod & container level, the container options
override the pod options.
Note that this field cannot be set when spec.os.name is windows.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          type indicates which kind of seccomp profile will be applied.
Valid options are:

Localhost - a profile defined in a file on the node should be used.
RuntimeDefault - the container runtime default profile should be used.
Unconfined - no profile should be applied.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>localhostProfile</b></td>
        <td>string</td>
        <td>
          localhostProfile indicates a profile defined in a file on the node should be used.
The profile must be preconfigured on the node to work.
Must be a descending path, relative to the kubelet's configured seccomp profile location.
Must be set if type is "Localhost". Must NOT be set for any other type.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].securityContext.windowsOptions
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexsecuritycontext)</sup></sup>



The Windows specific settings applied to all containers.
If unspecified, the options from the PodSecurityContext will be used.
If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
Note that this field cannot be set when spec.os.name is linux.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>gmsaCredentialSpec</b></td>
        <td>string</td>
        <td>
          GMSACredentialSpec is where the GMSA admission webhook
(https://github.com/kubernetes-sigs/windows-gmsa) inlines the contents of the
GMSA credential spec named by the GMSACredentialSpecName field.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>gmsaCredentialSpecName</b></td>
        <td>string</td>
        <td>
          GMSACredentialSpecName is the name of the GMSA credential spec to use.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostProcess</b></td>
        <td>boolean</td>
        <td>
          HostProcess determines if a container should be run as a 'Host Process' container.
All of a Pod's containers must have the same effective HostProcess value
(it is not allowed to have a mix of HostProcess containers and non-HostProcess containers).
In addition, if HostProcess is true then HostNetwork must also be set to true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsUserName</b></td>
        <td>string</td>
        <td>
          The UserName in Windows to run the entrypoint of the container process.
Defaults to the user specified in image metadata if unspecified.
May also be set in PodSecurityContext. If set in both SecurityContext and
PodSecurityContext, the value specified in SecurityContext takes precedence.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].startupProbe
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindex)</sup></sup>



StartupProbe indicates that the Pod has successfully initialized.
If specified, no other probes are executed until this completes successfully.
If this probe fails, the Pod will be restarted, just as if the livenessProbe failed.
This can be used to provide different probe parameters at the beginning of a Pod's lifecycle,
when it might take a long time to load data or warm a cache, than during steady-state operation.
This cannot be updated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridescontainersindexstartupprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexstartupprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexstartupprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexstartupprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.
The grace period is the duration in seconds after the processes running in the pod are sent
a termination signal and the time when the processes are forcibly halted with a kill signal.
Set this value longer than the expected cleanup time for your process.
If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this
value overrides the value provided by the pod spec.
Value must be non-negative integer. The value zero indicates stop immediately via
the kill signal (no opportunity to shut down).
This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate.
Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].startupProbe.exec
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexstartupprobe)</sup></sup>



Exec specifies a command to execute in the container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command  is root ('/') in the container's filesystem. The command is simply exec'd, it is
not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use
a shell, you need to explicitly call out to that shell.
Exit status of 0 is treated as live/healthy and non-zero is unhealthy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].startupProbe.grpc
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexstartupprobe)</sup></sup>



GRPC specifies a GRPC HealthCheckRequest.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).

If this is not specified, the default behavior is defined by gRPC.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].startupProbe.httpGet
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexstartupprobe)</sup></sup>



HTTPGet specifies an HTTP GET request to perform.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP. You probably want to set
"Host" in httpHeaders instead.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridescontainersindexstartupprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].startupProbe.httpGet.httpHeaders[index]
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexstartupprobehttpget)</sup></sup>



HTTPHeader describes a custom header to be used in HTTP probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The header field name.
This will be canonicalized upon output, so case-variant names will be understood as the same header.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].startupProbe.tcpSocket
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindexstartupprobe)</sup></sup>



TCPSocket specifies a connection to a TCP port.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].volumeDevices[index]
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindex)</sup></sup>



volumeDevice describes a mapping of a raw block device within a container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>devicePath</b></td>
        <td>string</td>
        <td>
          devicePath is the path inside of the container that the device will be mapped to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          name must match the name of a persistentVolumeClaim in the pod<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.containers[index].volumeMounts[index]
<sup><sup>[↩ Parent](#ironicspecoverridescontainersindex)</sup></sup>



VolumeMount describes a mounting of a Volume within a container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>mountPath</b></td>
        <td>string</td>
        <td>
          Path within the container at which the volume should be mounted.  Must
not contain ':'.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          This must match the Name of a Volume.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>mountPropagation</b></td>
        <td>string</td>
        <td>
          mountPropagation determines how mounts are propagated from the host
to container and the other way around.
When not set, MountPropagationNone is used.
This field is beta in 1.10.
When RecursiveReadOnly is set to IfPossible or to Enabled, MountPropagation must be None or unspecified
(which defaults to None).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          Mounted read-only if true, read-write otherwise (false or unspecified).
Defaults to false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>recursiveReadOnly</b></td>
        <td>string</td>
        <td>
          RecursiveReadOnly specifies whether read-only mounts should be handled
recursively.

If ReadOnly is false, this field has no meaning and must be unspecified.

If ReadOnly is true, and this field is set to Disabled, the mount is not made
recursively read-only.  If this field is set to IfPossible, the mount is made
recursively read-only, if it is supported by the container runtime.  If this
field is set to Enabled, the mount is made recursively read-only if it is
supported by the container runtime, otherwise the pod will not be started and
an error will be generated to indicate the reason.

If this field is set to IfPossible or Enabled, MountPropagation must be set to
None (or be unspecified, which defaults to None).

If this field is not specified, it is treated as an equivalent of Disabled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subPath</b></td>
        <td>string</td>
        <td>
          Path within the volume from which the container's volume should be mounted.
Defaults to "" (volume's root).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subPathExpr</b></td>
        <td>string</td>
        <td>
          Expanded path within the volume from which the container's volume should be mounted.
Behaves similarly to SubPath but environment variable references $(VAR_NAME) are expanded using the container's environment.
Defaults to "" (volume's root).
SubPathExpr and SubPath are mutually exclusive.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.httpdLivenessProbe
<sup><sup>[↩ Parent](#ironicspecoverrides)</sup></sup>



HttpdLivenessProbe overrides the httpd container liveness probe.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverrideshttpdlivenessprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverrideshttpdlivenessprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverrideshttpdlivenessprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverrideshttpdlivenessprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.
The grace period is the duration in seconds after the processes running in the pod are sent
a termination signal and the time when the processes are forcibly halted with a kill signal.
Set this value longer than the expected cleanup time for your process.
If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this
value overrides the value provided by the pod spec.
Value must be non-negative integer. The value zero indicates stop immediately via
the kill signal (no opportunity to shut down).
This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate.
Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.httpdLivenessProbe.exec
<sup><sup>[↩ Parent](#ironicspecoverrideshttpdlivenessprobe)</sup></sup>



Exec specifies a command to execute in the container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command  is root ('/') in the container's filesystem. The command is simply exec'd, it is
not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use
a shell, you need to explicitly call out to that shell.
Exit status of 0 is treated as live/healthy and non-zero is unhealthy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.httpdLivenessProbe.grpc
<sup><sup>[↩ Parent](#ironicspecoverrideshttpdlivenessprobe)</sup></sup>



GRPC specifies a GRPC HealthCheckRequest.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).

If this is not specified, the default behavior is defined by gRPC.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.httpdLivenessProbe.httpGet
<sup><sup>[↩ Parent](#ironicspecoverrideshttpdlivenessprobe)</sup></sup>



HTTPGet specifies an HTTP GET request to perform.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP. You probably want to set
"Host" in httpHeaders instead.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverrideshttpdlivenessprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.httpdLivenessProbe.httpGet.httpHeaders[index]
<sup><sup>[↩ Parent](#ironicspecoverrideshttpdlivenessprobehttpget)</sup></sup>



HTTPHeader describes a custom header to be used in HTTP probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The header field name.
This will be canonicalized upon output, so case-variant names will be understood as the same header.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.httpdLivenessProbe.tcpSocket
<sup><sup>[↩ Parent](#ironicspecoverrideshttpdlivenessprobe)</sup></sup>



TCPSocket specifies a connection to a TCP port.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.httpdReadinessProbe
<sup><sup>[↩ Parent](#ironicspecoverrides)</sup></sup>



HttpdReadinessProbe overrides the httpd container readiness probe.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverrideshttpdreadinessprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverrideshttpdreadinessprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverrideshttpdreadinessprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverrideshttpdreadinessprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.
The grace period is the duration in seconds after the processes running in the pod are sent
a termination signal and the time when the processes are forcibly halted with a kill signal.
Set this value longer than the expected cleanup time for your process.
If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this
value overrides the value provided by the pod spec.
Value must be non-negative integer. The value zero indicates stop immediately via
the kill signal (no opportunity to shut down).
This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate.
Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.httpdReadinessProbe.exec
<sup><sup>[↩ Parent](#ironicspecoverrideshttpdreadinessprobe)</sup></sup>



Exec specifies a command to execute in the container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command  is root ('/') in the container's filesystem. The command is simply exec'd, it is
not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use
a shell, you need to explicitly call out to that shell.
Exit status of 0 is treated as live/healthy and non-zero is unhealthy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.httpdReadinessProbe.grpc
<sup><sup>[↩ Parent](#ironicspecoverrideshttpdreadinessprobe)</sup></sup>



GRPC specifies a GRPC HealthCheckRequest.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).

If this is not specified, the default behavior is defined by gRPC.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.httpdReadinessProbe.httpGet
<sup><sup>[↩ Parent](#ironicspecoverrideshttpdreadinessprobe)</sup></sup>



HTTPGet specifies an HTTP GET request to perform.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP. You probably want to set
"Host" in httpHeaders instead.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverrideshttpdreadinessprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.httpdReadinessProbe.httpGet.httpHeaders[index]
<sup><sup>[↩ Parent](#ironicspecoverrideshttpdreadinessprobehttpget)</sup></sup>



HTTPHeader describes a custom header to be used in HTTP probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The header field name.
This will be canonicalized upon output, so case-variant names will be understood as the same header.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.httpdReadinessProbe.tcpSocket
<sup><sup>[↩ Parent](#ironicspecoverrideshttpdreadinessprobe)</sup></sup>



TCPSocket specifies a connection to a TCP port.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index]
<sup><sup>[↩ Parent](#ironicspecoverrides)</sup></sup>



A single application container that you want to run within a pod.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the container specified as a DNS_LABEL.
Each container in a pod must have a unique name (DNS_LABEL).
Cannot be updated.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>args</b></td>
        <td>[]string</td>
        <td>
          Arguments to the entrypoint.
The container image's CMD is used if this is not provided.
Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced
to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will
produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless
of whether the variable exists or not. Cannot be updated.
More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Entrypoint array. Not executed within a shell.
The container image's ENTRYPOINT is used if this is not provided.
Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced
to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will
produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless
of whether the variable exists or not. Cannot be updated.
More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexenvindex">env</a></b></td>
        <td>[]object</td>
        <td>
          List of environment variables to set in the container.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexenvfromindex">envFrom</a></b></td>
        <td>[]object</td>
        <td>
          List of sources to populate environment variables in the container.
The keys defined within a source may consist of any printable ASCII characters except '='.
When a key exists in multiple
sources, the value associated with the last source will take precedence.
Values defined by an Env with a duplicate key will take precedence.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>image</b></td>
        <td>string</td>
        <td>
          Container image name.
More info: https://kubernetes.io/docs/concepts/containers/images
This field is optional to allow higher level config management to default or override
container images in workload controllers like Deployments and StatefulSets.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>imagePullPolicy</b></td>
        <td>string</td>
        <td>
          Image pull policy.
One of Always, Never, IfNotPresent.
Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
Cannot be updated.
More info: https://kubernetes.io/docs/concepts/containers/images#updating-images<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlifecycle">lifecycle</a></b></td>
        <td>object</td>
        <td>
          Actions that the management system should take in response to container lifecycle events.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlivenessprobe">livenessProbe</a></b></td>
        <td>object</td>
        <td>
          Periodic probe of container liveness.
Container will be restarted if the probe fails.
Cannot be updated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexportsindex">ports</a></b></td>
        <td>[]object</td>
        <td>
          List of ports to expose from the container. Not specifying a port here
DOES NOT prevent that port from being exposed. Any port which is
listening on the default "0.0.0.0" address inside a container will be
accessible from the network.
Modifying this array with strategic merge patch may corrupt the data.
For more information See https://github.com/kubernetes/kubernetes/issues/108255.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexreadinessprobe">readinessProbe</a></b></td>
        <td>object</td>
        <td>
          Periodic probe of container service readiness.
Container will be removed from service endpoints if the probe fails.
Cannot be updated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexresizepolicyindex">resizePolicy</a></b></td>
        <td>[]object</td>
        <td>
          Resources resize policy for the container.
This field cannot be set on ephemeral containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexresources">resources</a></b></td>
        <td>object</td>
        <td>
          Compute Resources required by this container.
Cannot be updated.
More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>restartPolicy</b></td>
        <td>string</td>
        <td>
          RestartPolicy defines the restart behavior of individual containers in a pod.
This overrides the pod-level restart policy. When this field is not specified,
the restart behavior is defined by the Pod's restart policy and the container type.
Additionally, setting the RestartPolicy as "Always" for the init container will
have the following effect:
this init container will be continually restarted on
exit until all regular containers have terminated. Once all regular
containers have completed, all init containers with restartPolicy "Always"
will be shut down. This lifecycle differs from normal init containers and
is often referred to as a "sidecar" container. Although this init
container still starts in the init container sequence, it does not wait
for the container to complete before proceeding to the next init
container. Instead, the next init container starts immediately after this
init container is started, or after any startupProbe has successfully
completed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexrestartpolicyrulesindex">restartPolicyRules</a></b></td>
        <td>[]object</td>
        <td>
          Represents a list of rules to be checked to determine if the
container should be restarted on exit. The rules are evaluated in
order. Once a rule matches a container exit condition, the remaining
rules are ignored. If no rule matches the container exit condition,
the Container-level restart policy determines the whether the container
is restarted or not. Constraints on the rules:
- At most 20 rules are allowed.
- Rules can have the same action.
- Identical rules are not forbidden in validations.
When rules are specified, container MUST set RestartPolicy explicitly
even it if matches the Pod's RestartPolicy.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexsecuritycontext">securityContext</a></b></td>
        <td>object</td>
        <td>
          SecurityContext defines the security options the container should be run with.
If set, the fields of SecurityContext override the equivalent fields of PodSecurityContext.
More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexstartupprobe">startupProbe</a></b></td>
        <td>object</td>
        <td>
          StartupProbe indicates that the Pod has successfully initialized.
If specified, no other probes are executed until this completes successfully.
If this probe fails, the Pod will be restarted, just as if the livenessProbe failed.
This can be used to provide different probe parameters at the beginning of a Pod's lifecycle,
when it might take a long time to load data or warm a cache, than during steady-state operation.
This cannot be updated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>stdin</b></td>
        <td>boolean</td>
        <td>
          Whether this container should allocate a buffer for stdin in the container runtime. If this
is not set, reads from stdin in the container will always result in EOF.
Default is false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>stdinOnce</b></td>
        <td>boolean</td>
        <td>
          Whether the container runtime should close the stdin channel after it has been opened by
a single attach. When stdin is true the stdin stream will remain open across multiple attach
sessions. If stdinOnce is set to true, stdin is opened on container start, is empty until the
first client attaches to stdin, and then remains open and accepts data until the client disconnects,
at which time stdin is closed and remains closed until the container is restarted. If this
flag is false, a container processes that reads from stdin will never receive an EOF.
Default is false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationMessagePath</b></td>
        <td>string</td>
        <td>
          Optional: Path at which the file to which the container's termination message
will be written is mounted into the container's filesystem.
Message written is intended to be brief final status, such as an assertion failure message.
Will be truncated by the node if greater than 4096 bytes. The total message length across
all containers will be limited to 12kb.
Defaults to /dev/termination-log.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationMessagePolicy</b></td>
        <td>string</td>
        <td>
          Indicate how the termination message should be populated. File will use the contents of
terminationMessagePath to populate the container status message on both success and failure.
FallbackToLogsOnError will use the last chunk of container log output if the termination
message file is empty and the container exited with an error.
The log output is limited to 2048 bytes or 80 lines, whichever is smaller.
Defaults to File.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tty</b></td>
        <td>boolean</td>
        <td>
          Whether this container should allocate a TTY for itself, also requires 'stdin' to be true.
Default is false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexvolumedevicesindex">volumeDevices</a></b></td>
        <td>[]object</td>
        <td>
          volumeDevices is the list of block devices to be used by the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexvolumemountsindex">volumeMounts</a></b></td>
        <td>[]object</td>
        <td>
          Pod volumes to mount into the container's filesystem.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>workingDir</b></td>
        <td>string</td>
        <td>
          Container's working directory.
If not specified, the container runtime's default will be used, which
might be configured in the container image.
Cannot be updated.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].env[index]
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindex)</sup></sup>



EnvVar represents an environment variable present in a Container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the environment variable.
May consist of any printable ASCII characters except '='.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          Variable references $(VAR_NAME) are expanded
using the previously defined environment variables in the container and
any service environment variables. If a variable cannot be resolved,
the reference in the input string will be unchanged. Double $$ are reduced
to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e.
"$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)".
Escaped references will never be expanded, regardless of whether the variable
exists or not.
Defaults to "".<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexenvindexvaluefrom">valueFrom</a></b></td>
        <td>object</td>
        <td>
          Source for the environment variable's value. Cannot be used if value is not empty.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].env[index].valueFrom
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexenvindex)</sup></sup>



Source for the environment variable's value. Cannot be used if value is not empty.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexenvindexvaluefromconfigmapkeyref">configMapKeyRef</a></b></td>
        <td>object</td>
        <td>
          Selects a key of a ConfigMap.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexenvindexvaluefromfieldref">fieldRef</a></b></td>
        <td>object</td>
        <td>
          Selects a field of the pod: supports metadata.name, metadata.namespace, `metadata.labels['<KEY>']`, `metadata.annotations['<KEY>']`,
spec.nodeName, spec.serviceAccountName, status.hostIP, status.podIP, status.podIPs.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexenvindexvaluefromfilekeyref">fileKeyRef</a></b></td>
        <td>object</td>
        <td>
          FileKeyRef selects a key of the env file.
Requires the EnvFiles feature gate to be enabled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexenvindexvaluefromresourcefieldref">resourceFieldRef</a></b></td>
        <td>object</td>
        <td>
          Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.memory, limits.ephemeral-storage, requests.cpu, requests.memory and requests.ephemeral-storage) are currently supported.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexenvindexvaluefromsecretkeyref">secretKeyRef</a></b></td>
        <td>object</td>
        <td>
          Selects a key of a secret in the pod's namespace<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].env[index].valueFrom.configMapKeyRef
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexenvindexvaluefrom)</sup></sup>



Selects a key of a ConfigMap.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The key to select.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the ConfigMap or its key must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].env[index].valueFrom.fieldRef
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexenvindexvaluefrom)</sup></sup>



Selects a field of the pod: supports metadata.name, metadata.namespace, `metadata.labels['<KEY>']`, `metadata.annotations['<KEY>']`,
spec.nodeName, spec.serviceAccountName, status.hostIP, status.podIP, status.podIPs.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>fieldPath</b></td>
        <td>string</td>
        <td>
          Path of the field to select in the specified API version.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>apiVersion</b></td>
        <td>string</td>
        <td>
          Version of the schema the FieldPath is written in terms of, defaults to "v1".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].env[index].valueFrom.fileKeyRef
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexenvindexvaluefrom)</sup></sup>



FileKeyRef selects a key of the env file.
Requires the EnvFiles feature gate to be enabled.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The key within the env file. An invalid key will prevent the pod from starting.
The keys defined within a source may consist of any printable ASCII characters except '='.
During Alpha stage of the EnvFiles feature gate, the key size is limited to 128 characters.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          The path within the volume from which to select the file.
Must be relative and may not contain the '..' path or start with '..'.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>volumeName</b></td>
        <td>string</td>
        <td>
          The name of the volume mount containing the env file.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the file or its key must be defined. If the file or key
does not exist, then the env var is not published.
If optional is set to true and the specified key does not exist,
the environment variable will not be set in the Pod's containers.

If optional is set to false and the specified key does not exist,
an error will be returned during Pod creation.<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].env[index].valueFrom.resourceFieldRef
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexenvindexvaluefrom)</sup></sup>



Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.memory, limits.ephemeral-storage, requests.cpu, requests.memory and requests.ephemeral-storage) are currently supported.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resource</b></td>
        <td>string</td>
        <td>
          Required: resource to select<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>containerName</b></td>
        <td>string</td>
        <td>
          Container name: required for volumes, optional for env vars<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>divisor</b></td>
        <td>int or string</td>
        <td>
          Specifies the output format of the exposed resources, defaults to "1"<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].env[index].valueFrom.secretKeyRef
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexenvindexvaluefrom)</sup></sup>



Selects a key of a secret in the pod's namespace

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          The key of the secret to select from.  Must be a valid secret key.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the Secret or its key must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].envFrom[index]
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindex)</sup></sup>



EnvFromSource represents the source of a set of ConfigMaps or Secrets

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexenvfromindexconfigmapref">configMapRef</a></b></td>
        <td>object</td>
        <td>
          The ConfigMap to select from<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>prefix</b></td>
        <td>string</td>
        <td>
          Optional text to prepend to the name of each environment variable.
May consist of any printable ASCII characters except '='.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexenvfromindexsecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          The Secret to select from<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].envFrom[index].configMapRef
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexenvfromindex)</sup></sup>



The ConfigMap to select from

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the ConfigMap must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].envFrom[index].secretRef
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexenvfromindex)</sup></sup>



The Secret to select from

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          Specify whether the Secret must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].lifecycle
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindex)</sup></sup>



Actions that the management system should take in response to container lifecycle events.
Cannot be updated.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlifecyclepoststart">postStart</a></b></td>
        <td>object</td>
        <td>
          PostStart is called immediately after a container is created. If the handler fails,
the container is terminated and restarted according to its restart policy.
Other management of the container blocks until the hook completes.
More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlifecycleprestop">preStop</a></b></td>
        <td>object</td>
        <td>
          PreStop is called immediately before a container is terminated due to an
API request or management event such as liveness/startup probe failure,
preemption, resource contention, etc. The handler is not called if the
container crashes or exits. The Pod's termination grace period countdown begins before the
PreStop hook is executed. Regardless of the outcome of the handler, the
container will eventually terminate within the Pod's termination grace
period (unless delayed by finalizers). Other management of the container blocks until the hook completes
or until the termination grace period is reached.
More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>stopSignal</b></td>
        <td>string</td>
        <td>
          StopSignal defines which signal will be sent to a container when it is being stopped.
If not specified, the default is defined by the container runtime in use.
StopSignal can only be set for Pods with a non-empty .spec.os.name<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].lifecycle.postStart
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlifecycle)</sup></sup>



PostStart is called immediately after a container is created. If the handler fails,
the container is terminated and restarted according to its restart policy.
Other management of the container blocks until the hook completes.
More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlifecyclepoststartexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlifecyclepoststarthttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlifecyclepoststartsleep">sleep</a></b></td>
        <td>object</td>
        <td>
          Sleep represents a duration that the container should sleep.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlifecyclepoststarttcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility. There is no validation of this field and
lifecycle hooks will fail at runtime when it is specified.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].lifecycle.postStart.exec
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlifecyclepoststart)</sup></sup>



Exec specifies a command to execute in the container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command  is root ('/') in the container's filesystem. The command is simply exec'd, it is
not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use
a shell, you need to explicitly call out to that shell.
Exit status of 0 is treated as live/healthy and non-zero is unhealthy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].lifecycle.postStart.httpGet
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlifecyclepoststart)</sup></sup>



HTTPGet specifies an HTTP GET request to perform.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP. You probably want to set
"Host" in httpHeaders instead.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlifecyclepoststarthttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].lifecycle.postStart.httpGet.httpHeaders[index]
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlifecyclepoststarthttpget)</sup></sup>



HTTPHeader describes a custom header to be used in HTTP probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The header field name.
This will be canonicalized upon output, so case-variant names will be understood as the same header.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].lifecycle.postStart.sleep
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlifecyclepoststart)</sup></sup>



Sleep represents a duration that the container should sleep.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>seconds</b></td>
        <td>integer</td>
        <td>
          Seconds is the number of seconds to sleep.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].lifecycle.postStart.tcpSocket
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlifecyclepoststart)</sup></sup>



Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility. There is no validation of this field and
lifecycle hooks will fail at runtime when it is specified.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].lifecycle.preStop
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlifecycle)</sup></sup>



PreStop is called immediately before a container is terminated due to an
API request or management event such as liveness/startup probe failure,
preemption, resource contention, etc. The handler is not called if the
container crashes or exits. The Pod's termination grace period countdown begins before the
PreStop hook is executed. Regardless of the outcome of the handler, the
container will eventually terminate within the Pod's termination grace
period (unless delayed by finalizers). Other management of the container blocks until the hook completes
or until the termination grace period is reached.
More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlifecycleprestopexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlifecycleprestophttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlifecycleprestopsleep">sleep</a></b></td>
        <td>object</td>
        <td>
          Sleep represents a duration that the container should sleep.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlifecycleprestoptcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility. There is no validation of this field and
lifecycle hooks will fail at runtime when it is specified.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].lifecycle.preStop.exec
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlifecycleprestop)</sup></sup>



Exec specifies a command to execute in the container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command  is root ('/') in the container's filesystem. The command is simply exec'd, it is
not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use
a shell, you need to explicitly call out to that shell.
Exit status of 0 is treated as live/healthy and non-zero is unhealthy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].lifecycle.preStop.httpGet
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlifecycleprestop)</sup></sup>



HTTPGet specifies an HTTP GET request to perform.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP. You probably want to set
"Host" in httpHeaders instead.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlifecycleprestophttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].lifecycle.preStop.httpGet.httpHeaders[index]
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlifecycleprestophttpget)</sup></sup>



HTTPHeader describes a custom header to be used in HTTP probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The header field name.
This will be canonicalized upon output, so case-variant names will be understood as the same header.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].lifecycle.preStop.sleep
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlifecycleprestop)</sup></sup>



Sleep represents a duration that the container should sleep.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>seconds</b></td>
        <td>integer</td>
        <td>
          Seconds is the number of seconds to sleep.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].lifecycle.preStop.tcpSocket
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlifecycleprestop)</sup></sup>



Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
for backward compatibility. There is no validation of this field and
lifecycle hooks will fail at runtime when it is specified.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].livenessProbe
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindex)</sup></sup>



Periodic probe of container liveness.
Container will be restarted if the probe fails.
Cannot be updated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlivenessprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlivenessprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlivenessprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlivenessprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.
The grace period is the duration in seconds after the processes running in the pod are sent
a termination signal and the time when the processes are forcibly halted with a kill signal.
Set this value longer than the expected cleanup time for your process.
If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this
value overrides the value provided by the pod spec.
Value must be non-negative integer. The value zero indicates stop immediately via
the kill signal (no opportunity to shut down).
This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate.
Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].livenessProbe.exec
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlivenessprobe)</sup></sup>



Exec specifies a command to execute in the container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command  is root ('/') in the container's filesystem. The command is simply exec'd, it is
not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use
a shell, you need to explicitly call out to that shell.
Exit status of 0 is treated as live/healthy and non-zero is unhealthy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].livenessProbe.grpc
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlivenessprobe)</sup></sup>



GRPC specifies a GRPC HealthCheckRequest.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).

If this is not specified, the default behavior is defined by gRPC.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].livenessProbe.httpGet
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlivenessprobe)</sup></sup>



HTTPGet specifies an HTTP GET request to perform.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP. You probably want to set
"Host" in httpHeaders instead.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexlivenessprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].livenessProbe.httpGet.httpHeaders[index]
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlivenessprobehttpget)</sup></sup>



HTTPHeader describes a custom header to be used in HTTP probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The header field name.
This will be canonicalized upon output, so case-variant names will be understood as the same header.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].livenessProbe.tcpSocket
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexlivenessprobe)</sup></sup>



TCPSocket specifies a connection to a TCP port.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].ports[index]
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindex)</sup></sup>



ContainerPort represents a network port in a single container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>containerPort</b></td>
        <td>integer</td>
        <td>
          Number of port to expose on the pod's IP address.
This must be a valid port number, 0 < x < 65536.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>hostIP</b></td>
        <td>string</td>
        <td>
          What host IP to bind the external port to.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostPort</b></td>
        <td>integer</td>
        <td>
          Number of port to expose on the host.
If specified, this must be a valid port number, 0 < x < 65536.
If HostNetwork is specified, this must match ContainerPort.
Most containers do not need this.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          If specified, this must be an IANA_SVC_NAME and unique within the pod. Each
named port in a pod must have a unique name. Name for the port that can be
referred to by services.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>protocol</b></td>
        <td>string</td>
        <td>
          Protocol for port. Must be UDP, TCP, or SCTP.
Defaults to "TCP".<br/>
          <br/>
            <i>Default</i>: TCP<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].readinessProbe
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindex)</sup></sup>



Periodic probe of container service readiness.
Container will be removed from service endpoints if the probe fails.
Cannot be updated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexreadinessprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexreadinessprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexreadinessprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexreadinessprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.
The grace period is the duration in seconds after the processes running in the pod are sent
a termination signal and the time when the processes are forcibly halted with a kill signal.
Set this value longer than the expected cleanup time for your process.
If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this
value overrides the value provided by the pod spec.
Value must be non-negative integer. The value zero indicates stop immediately via
the kill signal (no opportunity to shut down).
This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate.
Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].readinessProbe.exec
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexreadinessprobe)</sup></sup>



Exec specifies a command to execute in the container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command  is root ('/') in the container's filesystem. The command is simply exec'd, it is
not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use
a shell, you need to explicitly call out to that shell.
Exit status of 0 is treated as live/healthy and non-zero is unhealthy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].readinessProbe.grpc
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexreadinessprobe)</sup></sup>



GRPC specifies a GRPC HealthCheckRequest.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).

If this is not specified, the default behavior is defined by gRPC.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].readinessProbe.httpGet
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexreadinessprobe)</sup></sup>



HTTPGet specifies an HTTP GET request to perform.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP. You probably want to set
"Host" in httpHeaders instead.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexreadinessprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].readinessProbe.httpGet.httpHeaders[index]
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexreadinessprobehttpget)</sup></sup>



HTTPHeader describes a custom header to be used in HTTP probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The header field name.
This will be canonicalized upon output, so case-variant names will be understood as the same header.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].readinessProbe.tcpSocket
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexreadinessprobe)</sup></sup>



TCPSocket specifies a connection to a TCP port.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].resizePolicy[index]
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindex)</sup></sup>



ContainerResizePolicy represents resource resize policy for the container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resourceName</b></td>
        <td>string</td>
        <td>
          Name of the resource to which this resource resize policy applies.
Supported values: cpu, memory.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>restartPolicy</b></td>
        <td>string</td>
        <td>
          Restart policy to apply when specified resource is resized.
If not specified, it defaults to NotRequired.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].resources
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindex)</sup></sup>



Compute Resources required by this container.
Cannot be updated.
More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexresourcesclaimsindex">claims</a></b></td>
        <td>[]object</td>
        <td>
          Claims lists the names of resources, defined in spec.resourceClaims,
that are used by this container.

This field depends on the
DynamicResourceAllocation feature gate.

This field is immutable. It can only be set for containers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>limits</b></td>
        <td>map[string]int or string</td>
        <td>
          Limits describes the maximum amount of compute resources allowed.
More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>requests</b></td>
        <td>map[string]int or string</td>
        <td>
          Requests describes the minimum amount of compute resources required.
If Requests is omitted for a container, it defaults to Limits if that is explicitly specified,
otherwise to an implementation-defined value. Requests cannot exceed Limits.
More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].resources.claims[index]
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexresources)</sup></sup>



ResourceClaim references one entry in PodSpec.ResourceClaims.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name must match the name of one entry in pod.spec.resourceClaims of
the Pod where this field is used. It makes that resource available
inside a container.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>request</b></td>
        <td>string</td>
        <td>
          Request is the name chosen for a request in the referenced claim.
If empty, everything from the claim is made available, otherwise
only the result of this request.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].restartPolicyRules[index]
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindex)</sup></sup>



ContainerRestartRule describes how a container exit is handled.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>action</b></td>
        <td>string</td>
        <td>
          Specifies the action taken on a container exit if the requirements
are satisfied. The only possible value is "Restart" to restart the
container.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexrestartpolicyrulesindexexitcodes">exitCodes</a></b></td>
        <td>object</td>
        <td>
          Represents the exit codes to check on container exits.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].restartPolicyRules[index].exitCodes
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexrestartpolicyrulesindex)</sup></sup>



Represents the exit codes to check on container exits.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          Represents the relationship between the container exit code(s) and the
specified values. Possible values are:
- In: the requirement is satisfied if the container exit code is in the
  set of specified values.
- NotIn: the requirement is satisfied if the container exit code is
  not in the set of specified values.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]integer</td>
        <td>
          Specifies the set of values to check for container exit codes.
At most 255 elements are allowed.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].securityContext
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindex)</sup></sup>



SecurityContext defines the security options the container should be run with.
If set, the fields of SecurityContext override the equivalent fields of PodSecurityContext.
More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>allowPrivilegeEscalation</b></td>
        <td>boolean</td>
        <td>
          AllowPrivilegeEscalation controls whether a process can gain more
privileges than its parent process. This bool directly controls if
the no_new_privs flag will be set on the container process.
AllowPrivilegeEscalation is true always when the container is:
1) run as Privileged
2) has CAP_SYS_ADMIN
Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexsecuritycontextapparmorprofile">appArmorProfile</a></b></td>
        <td>object</td>
        <td>
          appArmorProfile is the AppArmor options to use by this container. If set, this profile
overrides the pod's appArmorProfile.
Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexsecuritycontextcapabilities">capabilities</a></b></td>
        <td>object</td>
        <td>
          The capabilities to add/drop when running containers.
Defaults to the default set of capabilities granted by the container runtime.
Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>privileged</b></td>
        <td>boolean</td>
        <td>
          Run container in privileged mode.
Processes in privileged containers are essentially equivalent to root on the host.
Defaults to false.
Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>procMount</b></td>
        <td>string</td>
        <td>
          procMount denotes the type of proc mount to use for the containers.
The default value is Default which uses the container runtime defaults for
readonly paths and masked paths.
Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnlyRootFilesystem</b></td>
        <td>boolean</td>
        <td>
          Whether this container has a read-only root filesystem.
Default is false.
Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsGroup</b></td>
        <td>integer</td>
        <td>
          The GID to run the entrypoint of the container process.
Uses runtime default if unset.
May also be set in PodSecurityContext.  If set in both SecurityContext and
PodSecurityContext, the value specified in SecurityContext takes precedence.
Note that this field cannot be set when spec.os.name is windows.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsNonRoot</b></td>
        <td>boolean</td>
        <td>
          Indicates that the container must run as a non-root user.
If true, the Kubelet will validate the image at runtime to ensure that it
does not run as UID 0 (root) and fail to start the container if it does.
If unset or false, no such validation will be performed.
May also be set in PodSecurityContext.  If set in both SecurityContext and
PodSecurityContext, the value specified in SecurityContext takes precedence.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsUser</b></td>
        <td>integer</td>
        <td>
          The UID to run the entrypoint of the container process.
Defaults to user specified in image metadata if unspecified.
May also be set in PodSecurityContext.  If set in both SecurityContext and
PodSecurityContext, the value specified in SecurityContext takes precedence.
Note that this field cannot be set when spec.os.name is windows.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexsecuritycontextselinuxoptions">seLinuxOptions</a></b></td>
        <td>object</td>
        <td>
          The SELinux context to be applied to the container.
If unspecified, the container runtime will allocate a random SELinux context for each
container.  May also be set in PodSecurityContext.  If set in both SecurityContext and
PodSecurityContext, the value specified in SecurityContext takes precedence.
Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexsecuritycontextseccompprofile">seccompProfile</a></b></td>
        <td>object</td>
        <td>
          The seccomp options to use by this container. If seccomp options are
provided at both the pod & container level, the container options
override the pod options.
Note that this field cannot be set when spec.os.name is windows.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexsecuritycontextwindowsoptions">windowsOptions</a></b></td>
        <td>object</td>
        <td>
          The Windows specific settings applied to all containers.
If unspecified, the options from the PodSecurityContext will be used.
If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
Note that this field cannot be set when spec.os.name is linux.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].securityContext.appArmorProfile
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexsecuritycontext)</sup></sup>



appArmorProfile is the AppArmor options to use by this container. If set, this profile
overrides the pod's appArmorProfile.
Note that this field cannot be set when spec.os.name is windows.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          type indicates which kind of AppArmor profile will be applied.
Valid options are:
  Localhost - a profile pre-loaded on the node.
  RuntimeDefault - the container runtime's default profile.
  Unconfined - no AppArmor enforcement.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>localhostProfile</b></td>
        <td>string</td>
        <td>
          localhostProfile indicates a profile loaded on the node that should be used.
The profile must be preconfigured on the node to work.
Must match the loaded name of the profile.
Must be set if and only if type is "Localhost".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].securityContext.capabilities
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexsecuritycontext)</sup></sup>



The capabilities to add/drop when running containers.
Defaults to the default set of capabilities granted by the container runtime.
Note that this field cannot be set when spec.os.name is windows.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>add</b></td>
        <td>[]string</td>
        <td>
          Added capabilities<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>drop</b></td>
        <td>[]string</td>
        <td>
          Removed capabilities<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].securityContext.seLinuxOptions
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexsecuritycontext)</sup></sup>



The SELinux context to be applied to the container.
If unspecified, the container runtime will allocate a random SELinux context for each
container.  May also be set in PodSecurityContext.  If set in both SecurityContext and
PodSecurityContext, the value specified in SecurityContext takes precedence.
Note that this field cannot be set when spec.os.name is windows.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>level</b></td>
        <td>string</td>
        <td>
          Level is SELinux level label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>role</b></td>
        <td>string</td>
        <td>
          Role is a SELinux role label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type is a SELinux type label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>user</b></td>
        <td>string</td>
        <td>
          User is a SELinux user label that applies to the container.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].securityContext.seccompProfile
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexsecuritycontext)</sup></sup>



The seccomp options to use by this container. If seccomp options are
provided at both the pod & container level, the container options
override the pod options.
Note that this field cannot be set when spec.os.name is windows.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          type indicates which kind of seccomp profile will be applied.
Valid options are:

Localhost - a profile defined in a file on the node should be used.
RuntimeDefault - the container runtime default profile should be used.
Unconfined - no profile should be applied.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>localhostProfile</b></td>
        <td>string</td>
        <td>
          localhostProfile indicates a profile defined in a file on the node should be used.
The profile must be preconfigured on the node to work.
Must be a descending path, relative to the kubelet's configured seccomp profile location.
Must be set if type is "Localhost". Must NOT be set for any other type.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].securityContext.windowsOptions
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexsecuritycontext)</sup></sup>



The Windows specific settings applied to all containers.
If unspecified, the options from the PodSecurityContext will be used.
If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
Note that this field cannot be set when spec.os.name is linux.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>gmsaCredentialSpec</b></td>
        <td>string</td>
        <td>
          GMSACredentialSpec is where the GMSA admission webhook
(https://github.com/kubernetes-sigs/windows-gmsa) inlines the contents of the
GMSA credential spec named by the GMSACredentialSpecName field.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>gmsaCredentialSpecName</b></td>
        <td>string</td>
        <td>
          GMSACredentialSpecName is the name of the GMSA credential spec to use.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>hostProcess</b></td>
        <td>boolean</td>
        <td>
          HostProcess determines if a container should be run as a 'Host Process' container.
All of a Pod's containers must have the same effective HostProcess value
(it is not allowed to have a mix of HostProcess containers and non-HostProcess containers).
In addition, if HostProcess is true then HostNetwork must also be set to true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>runAsUserName</b></td>
        <td>string</td>
        <td>
          The UserName in Windows to run the entrypoint of the container process.
Defaults to the user specified in image metadata if unspecified.
May also be set in PodSecurityContext. If set in both SecurityContext and
PodSecurityContext, the value specified in SecurityContext takes precedence.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].startupProbe
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindex)</sup></sup>



StartupProbe indicates that the Pod has successfully initialized.
If specified, no other probes are executed until this completes successfully.
If this probe fails, the Pod will be restarted, just as if the livenessProbe failed.
This can be used to provide different probe parameters at the beginning of a Pod's lifecycle,
when it might take a long time to load data or warm a cache, than during steady-state operation.
This cannot be updated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexstartupprobeexec">exec</a></b></td>
        <td>object</td>
        <td>
          Exec specifies a command to execute in the container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failureThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive failures for the probe to be considered failed after having succeeded.
Defaults to 3. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexstartupprobegrpc">grpc</a></b></td>
        <td>object</td>
        <td>
          GRPC specifies a GRPC HealthCheckRequest.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexstartupprobehttpget">httpGet</a></b></td>
        <td>object</td>
        <td>
          HTTPGet specifies an HTTP GET request to perform.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialDelaySeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after the container has started before liveness probes are initiated.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>periodSeconds</b></td>
        <td>integer</td>
        <td>
          How often (in seconds) to perform the probe.
Default to 10 seconds. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>successThreshold</b></td>
        <td>integer</td>
        <td>
          Minimum consecutive successes for the probe to be considered successful after having failed.
Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexstartupprobetcpsocket">tcpSocket</a></b></td>
        <td>object</td>
        <td>
          TCPSocket specifies a connection to a TCP port.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>terminationGracePeriodSeconds</b></td>
        <td>integer</td>
        <td>
          Optional duration in seconds the pod needs to terminate gracefully upon probe failure.
The grace period is the duration in seconds after the processes running in the pod are sent
a termination signal and the time when the processes are forcibly halted with a kill signal.
Set this value longer than the expected cleanup time for your process.
If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this
value overrides the value provided by the pod spec.
Value must be non-negative integer. The value zero indicates stop immediately via
the kill signal (no opportunity to shut down).
This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate.
Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>timeoutSeconds</b></td>
        <td>integer</td>
        <td>
          Number of seconds after which the probe times out.
Defaults to 1 second. Minimum value is 1.
More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].startupProbe.exec
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexstartupprobe)</sup></sup>



Exec specifies a command to execute in the container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>command</b></td>
        <td>[]string</td>
        <td>
          Command is the command line to execute inside the container, the working directory for the
command  is root ('/') in the container's filesystem. The command is simply exec'd, it is
not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use
a shell, you need to explicitly call out to that shell.
Exit status of 0 is treated as live/healthy and non-zero is unhealthy.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].startupProbe.grpc
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexstartupprobe)</sup></sup>



GRPC specifies a GRPC HealthCheckRequest.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>integer</td>
        <td>
          Port number of the gRPC service. Number must be in the range 1 to 65535.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>service</b></td>
        <td>string</td>
        <td>
          Service is the name of the service to place in the gRPC HealthCheckRequest
(see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).

If this is not specified, the default behavior is defined by gRPC.<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].startupProbe.httpGet
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexstartupprobe)</sup></sup>



HTTPGet specifies an HTTP GET request to perform.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Name or number of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Host name to connect to, defaults to the pod IP. You probably want to set
"Host" in httpHeaders instead.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesinitcontainersindexstartupprobehttpgethttpheadersindex">httpHeaders</a></b></td>
        <td>[]object</td>
        <td>
          Custom headers to set in the request. HTTP allows repeated headers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Path to access on the HTTP server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scheme</b></td>
        <td>string</td>
        <td>
          Scheme to use for connecting to the host.
Defaults to HTTP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].startupProbe.httpGet.httpHeaders[index]
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexstartupprobehttpget)</sup></sup>



HTTPHeader describes a custom header to be used in HTTP probes

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The header field name.
This will be canonicalized upon output, so case-variant names will be understood as the same header.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          The header field value<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].startupProbe.tcpSocket
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindexstartupprobe)</sup></sup>



TCPSocket specifies a connection to a TCP port.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>port</b></td>
        <td>int or string</td>
        <td>
          Number or name of the port to access on the container.
Number must be in the range 1 to 65535.
Name must be an IANA_SVC_NAME.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>host</b></td>
        <td>string</td>
        <td>
          Optional: Host name to connect to, defaults to the pod IP.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].volumeDevices[index]
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindex)</sup></sup>



volumeDevice describes a mapping of a raw block device within a container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>devicePath</b></td>
        <td>string</td>
        <td>
          devicePath is the path inside of the container that the device will be mapped to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          name must match the name of a persistentVolumeClaim in the pod<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.initContainers[index].volumeMounts[index]
<sup><sup>[↩ Parent](#ironicspecoverridesinitcontainersindex)</sup></sup>



VolumeMount describes a mounting of a Volume within a container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>mountPath</b></td>
        <td>string</td>
        <td>
          Path within the container at which the volume should be mounted.  Must
not contain ':'.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          This must match the Name of a Volume.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>mountPropagation</b></td>
        <td>string</td>
        <td>
          mountPropagation determines how mounts are propagated from the host
to container and the other way around.
When not set, MountPropagationNone is used.
This field is beta in 1.10.
When RecursiveReadOnly is set to IfPossible or to Enabled, MountPropagation must be None or unspecified
(which defaults to None).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          Mounted read-only if true, read-write otherwise (false or unspecified).
Defaults to false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>recursiveReadOnly</b></td>
        <td>string</td>
        <td>
          RecursiveReadOnly specifies whether read-only mounts should be handled
recursively.

If ReadOnly is false, this field has no meaning and must be unspecified.

If ReadOnly is true, and this field is set to Disabled, the mount is not made
recursively read-only.  If this field is set to IfPossible, the mount is made
recursively read-only, if it is supported by the container runtime.  If this
field is set to Enabled, the mount is made recursively read-only if it is
supported by the container runtime, otherwise the pod will not be started and
an error will be generated to indicate the reason.

If this field is set to IfPossible or Enabled, MountPropagation must be set to
None (or be unspecified, which defaults to None).

If this field is not specified, it is treated as an equivalent of Disabled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subPath</b></td>
        <td>string</td>
        <td>
          Path within the volume from which the container's volume should be mounted.
Defaults to "" (volume's root).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subPathExpr</b></td>
        <td>string</td>
        <td>
          Expanded path within the volume from which the container's volume should be mounted.
Behaves similarly to SubPath but environment variable references $(VAR_NAME) are expanded using the container's environment.
Defaults to "" (volume's root).
SubPathExpr and SubPath are mutually exclusive.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index]
<sup><sup>[↩ Parent](#ironicspecoverrides)</sup></sup>



Volume represents a named volume in a pod that may be accessed by any container in the pod.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          name of the volume.
Must be a DNS_LABEL and unique within the pod.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexawselasticblockstore">awsElasticBlockStore</a></b></td>
        <td>object</td>
        <td>
          awsElasticBlockStore represents an AWS Disk resource that is attached to a
kubelet's host machine and then exposed to the pod.
Deprecated: AWSElasticBlockStore is deprecated. All operations for the in-tree
awsElasticBlockStore type are redirected to the ebs.csi.aws.com CSI driver.
More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexazuredisk">azureDisk</a></b></td>
        <td>object</td>
        <td>
          azureDisk represents an Azure Data Disk mount on the host and bind mount to the pod.
Deprecated: AzureDisk is deprecated. All operations for the in-tree azureDisk type
are redirected to the disk.csi.azure.com CSI driver.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexazurefile">azureFile</a></b></td>
        <td>object</td>
        <td>
          azureFile represents an Azure File Service mount on the host and bind mount to the pod.
Deprecated: AzureFile is deprecated. All operations for the in-tree azureFile type
are redirected to the file.csi.azure.com CSI driver.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexcephfs">cephfs</a></b></td>
        <td>object</td>
        <td>
          cephFS represents a Ceph FS mount on the host that shares a pod's lifetime.
Deprecated: CephFS is deprecated and the in-tree cephfs type is no longer supported.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexcinder">cinder</a></b></td>
        <td>object</td>
        <td>
          cinder represents a cinder volume attached and mounted on kubelets host machine.
Deprecated: Cinder is deprecated. All operations for the in-tree cinder type
are redirected to the cinder.csi.openstack.org CSI driver.
More info: https://examples.k8s.io/mysql-cinder-pd/README.md<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexconfigmap">configMap</a></b></td>
        <td>object</td>
        <td>
          configMap represents a configMap that should populate this volume<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexcsi">csi</a></b></td>
        <td>object</td>
        <td>
          csi (Container Storage Interface) represents ephemeral storage that is handled by certain external CSI drivers.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexdownwardapi">downwardAPI</a></b></td>
        <td>object</td>
        <td>
          downwardAPI represents downward API about the pod that should populate this volume<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexemptydir">emptyDir</a></b></td>
        <td>object</td>
        <td>
          emptyDir represents a temporary directory that shares a pod's lifetime.
More info: https://kubernetes.io/docs/concepts/storage/volumes#emptydir<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexephemeral">ephemeral</a></b></td>
        <td>object</td>
        <td>
          ephemeral represents a volume that is handled by a cluster storage driver.
The volume's lifecycle is tied to the pod that defines it - it will be created before the pod starts,
and deleted when the pod is removed.

Use this if:
a) the volume is only needed while the pod runs,
b) features of normal volumes like restoring from snapshot or capacity
   tracking are needed,
c) the storage driver is specified through a storage class, and
d) the storage driver supports dynamic volume provisioning through
   a PersistentVolumeClaim (see EphemeralVolumeSource for more
   information on the connection between this volume type
   and PersistentVolumeClaim).

Use PersistentVolumeClaim or one of the vendor-specific
APIs for volumes that persist for longer than the lifecycle
of an individual pod.

Use CSI for light-weight local ephemeral volumes if the CSI driver is meant to
be used that way - see the documentation of the driver for
more information.

A pod can use both types of ephemeral volumes and
persistent volumes at the same time.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexfc">fc</a></b></td>
        <td>object</td>
        <td>
          fc represents a Fibre Channel resource that is attached to a kubelet's host machine and then exposed to the pod.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexflexvolume">flexVolume</a></b></td>
        <td>object</td>
        <td>
          flexVolume represents a generic volume resource that is
provisioned/attached using an exec based plugin.
Deprecated: FlexVolume is deprecated. Consider using a CSIDriver instead.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexflocker">flocker</a></b></td>
        <td>object</td>
        <td>
          flocker represents a Flocker volume attached to a kubelet's host machine. This depends on the Flocker control service being running.
Deprecated: Flocker is deprecated and the in-tree flocker type is no longer supported.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexgcepersistentdisk">gcePersistentDisk</a></b></td>
        <td>object</td>
        <td>
          gcePersistentDisk represents a GCE Disk resource that is attached to a
kubelet's host machine and then exposed to the pod.
Deprecated: GCEPersistentDisk is deprecated. All operations for the in-tree
gcePersistentDisk type are redirected to the pd.csi.storage.gke.io CSI driver.
More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexgitrepo">gitRepo</a></b></td>
        <td>object</td>
        <td>
          gitRepo represents a git repository at a particular revision.
Deprecated: GitRepo is deprecated. To provision a container with a git repo, mount an
EmptyDir into an InitContainer that clones the repo using git, then mount the EmptyDir
into the Pod's container.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexglusterfs">glusterfs</a></b></td>
        <td>object</td>
        <td>
          glusterfs represents a Glusterfs mount on the host that shares a pod's lifetime.
Deprecated: Glusterfs is deprecated and the in-tree glusterfs type is no longer supported.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexhostpath">hostPath</a></b></td>
        <td>object</td>
        <td>
          hostPath represents a pre-existing file or directory on the host
machine that is directly exposed to the container. This is generally
used for system agents or other privileged things that are allowed
to see the host machine. Most containers will NOT need this.
More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindeximage">image</a></b></td>
        <td>object</td>
        <td>
          image represents an OCI object (a container image or artifact) pulled and mounted on the kubelet's host machine.
The volume is resolved at pod startup depending on which PullPolicy value is provided:

- Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.
- Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.
- IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.

The volume gets re-resolved if the pod gets deleted and recreated, which means that new remote content will become available on pod recreation.
A failure to resolve or pull the image during pod startup will block containers from starting and may add significant latency. Failures will be retried using normal volume backoff and will be reported on the pod reason and message.
The types of objects that may be mounted by this volume are defined by the container runtime implementation on a host machine and at minimum must include all valid types supported by the container image field.
The OCI object gets mounted in a single directory (spec.containers[*].volumeMounts.mountPath) by merging the manifest layers in the same way as for container images.
The volume will be mounted read-only (ro) and non-executable files (noexec).
Sub path mounts for containers are not supported (spec.containers[*].volumeMounts.subpath) before 1.33.
The field spec.securityContext.fsGroupChangePolicy has no effect on this volume type.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexiscsi">iscsi</a></b></td>
        <td>object</td>
        <td>
          iscsi represents an ISCSI Disk resource that is attached to a
kubelet's host machine and then exposed to the pod.
More info: https://kubernetes.io/docs/concepts/storage/volumes/#iscsi<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexnfs">nfs</a></b></td>
        <td>object</td>
        <td>
          nfs represents an NFS mount on the host that shares a pod's lifetime
More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexpersistentvolumeclaim">persistentVolumeClaim</a></b></td>
        <td>object</td>
        <td>
          persistentVolumeClaimVolumeSource represents a reference to a
PersistentVolumeClaim in the same namespace.
More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexphotonpersistentdisk">photonPersistentDisk</a></b></td>
        <td>object</td>
        <td>
          photonPersistentDisk represents a PhotonController persistent disk attached and mounted on kubelets host machine.
Deprecated: PhotonPersistentDisk is deprecated and the in-tree photonPersistentDisk type is no longer supported.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexportworxvolume">portworxVolume</a></b></td>
        <td>object</td>
        <td>
          portworxVolume represents a portworx volume attached and mounted on kubelets host machine.
Deprecated: PortworxVolume is deprecated. All operations for the in-tree portworxVolume type
are redirected to the pxd.portworx.com CSI driver when the CSIMigrationPortworx feature-gate
is on.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexprojected">projected</a></b></td>
        <td>object</td>
        <td>
          projected items for all in one resources secrets, configmaps, and downward API<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexquobyte">quobyte</a></b></td>
        <td>object</td>
        <td>
          quobyte represents a Quobyte mount on the host that shares a pod's lifetime.
Deprecated: Quobyte is deprecated and the in-tree quobyte type is no longer supported.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexrbd">rbd</a></b></td>
        <td>object</td>
        <td>
          rbd represents a Rados Block Device mount on the host that shares a pod's lifetime.
Deprecated: RBD is deprecated and the in-tree rbd type is no longer supported.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexscaleio">scaleIO</a></b></td>
        <td>object</td>
        <td>
          scaleIO represents a ScaleIO persistent volume attached and mounted on Kubernetes nodes.
Deprecated: ScaleIO is deprecated and the in-tree scaleIO type is no longer supported.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexsecret">secret</a></b></td>
        <td>object</td>
        <td>
          secret represents a secret that should populate this volume.
More info: https://kubernetes.io/docs/concepts/storage/volumes#secret<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexstorageos">storageos</a></b></td>
        <td>object</td>
        <td>
          storageOS represents a StorageOS volume attached and mounted on Kubernetes nodes.
Deprecated: StorageOS is deprecated and the in-tree storageos type is no longer supported.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexvspherevolume">vsphereVolume</a></b></td>
        <td>object</td>
        <td>
          vsphereVolume represents a vSphere volume attached and mounted on kubelets host machine.
Deprecated: VsphereVolume is deprecated. All operations for the in-tree vsphereVolume type
are redirected to the csi.vsphere.vmware.com CSI driver.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].awsElasticBlockStore
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



awsElasticBlockStore represents an AWS Disk resource that is attached to a
kubelet's host machine and then exposed to the pod.
Deprecated: AWSElasticBlockStore is deprecated. All operations for the in-tree
awsElasticBlockStore type are redirected to the ebs.csi.aws.com CSI driver.
More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>volumeID</b></td>
        <td>string</td>
        <td>
          volumeID is unique ID of the persistent disk resource in AWS (Amazon EBS volume).
More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type of the volume that you want to mount.
Tip: Ensure that the filesystem type is supported by the host operating system.
Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>partition</b></td>
        <td>integer</td>
        <td>
          partition is the partition in the volume that you want to mount.
If omitted, the default is to mount by volume name.
Examples: For volume /dev/sda1, you specify the partition as "1".
Similarly, the volume partition for /dev/sda is "0" (or you can leave the property empty).<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly value true will force the readOnly setting in VolumeMounts.
More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].azureDisk
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



azureDisk represents an Azure Data Disk mount on the host and bind mount to the pod.
Deprecated: AzureDisk is deprecated. All operations for the in-tree azureDisk type
are redirected to the disk.csi.azure.com CSI driver.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>diskName</b></td>
        <td>string</td>
        <td>
          diskName is the Name of the data disk in the blob storage<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>diskURI</b></td>
        <td>string</td>
        <td>
          diskURI is the URI of data disk in the blob storage<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>cachingMode</b></td>
        <td>string</td>
        <td>
          cachingMode is the Host Caching mode: None, Read Only, Read Write.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is Filesystem type to mount.
Must be a filesystem type supported by the host operating system.
Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.<br/>
          <br/>
            <i>Default</i>: ext4<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          kind expected values are Shared: multiple blob disks per storage account  Dedicated: single blob disk per storage account  Managed: azure managed data disk (only in managed availability set). defaults to shared<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly Defaults to false (read/write). ReadOnly here will force
the ReadOnly setting in VolumeMounts.<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].azureFile
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



azureFile represents an Azure File Service mount on the host and bind mount to the pod.
Deprecated: AzureFile is deprecated. All operations for the in-tree azureFile type
are redirected to the file.csi.azure.com CSI driver.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>secretName</b></td>
        <td>string</td>
        <td>
          secretName is the  name of secret that contains Azure Storage Account Name and Key<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>shareName</b></td>
        <td>string</td>
        <td>
          shareName is the azure share Name<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly defaults to false (read/write). ReadOnly here will force
the ReadOnly setting in VolumeMounts.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].cephfs
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



cephFS represents a Ceph FS mount on the host that shares a pod's lifetime.
Deprecated: CephFS is deprecated and the in-tree cephfs type is no longer supported.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>monitors</b></td>
        <td>[]string</td>
        <td>
          monitors is Required: Monitors is a collection of Ceph monitors
More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          path is Optional: Used as the mounted root, rather than the full Ceph tree, default is /<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly is Optional: Defaults to false (read/write). ReadOnly here will force
the ReadOnly setting in VolumeMounts.
More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>secretFile</b></td>
        <td>string</td>
        <td>
          secretFile is Optional: SecretFile is the path to key ring for User, default is /etc/ceph/user.secret
More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexcephfssecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          secretRef is Optional: SecretRef is reference to the authentication secret for User, default is empty.
More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>user</b></td>
        <td>string</td>
        <td>
          user is optional: User is the rados user name, default is admin
More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].cephfs.secretRef
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexcephfs)</sup></sup>



secretRef is Optional: SecretRef is reference to the authentication secret for User, default is empty.
More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].cinder
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



cinder represents a cinder volume attached and mounted on kubelets host machine.
Deprecated: Cinder is deprecated. All operations for the in-tree cinder type
are redirected to the cinder.csi.openstack.org CSI driver.
More info: https://examples.k8s.io/mysql-cinder-pd/README.md

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>volumeID</b></td>
        <td>string</td>
        <td>
          volumeID used to identify the volume in cinder.
More info: https://examples.k8s.io/mysql-cinder-pd/README.md<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type to mount.
Must be a filesystem type supported by the host operating system.
Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
More info: https://examples.k8s.io/mysql-cinder-pd/README.md<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly defaults to false (read/write). ReadOnly here will force
the ReadOnly setting in VolumeMounts.
More info: https://examples.k8s.io/mysql-cinder-pd/README.md<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexcindersecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          secretRef is optional: points to a secret object containing parameters used to connect
to OpenStack.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].cinder.secretRef
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexcinder)</sup></sup>



secretRef is optional: points to a secret object containing parameters used to connect
to OpenStack.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].configMap
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



configMap represents a configMap that should populate this volume

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>defaultMode</b></td>
        <td>integer</td>
        <td>
          defaultMode is optional: mode bits used to set permissions on created files by default.
Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
Defaults to 0644.
Directories within the path are not affected by this setting.
This might be in conflict with other options that affect the file
mode, like fsGroup, and the result can be other mode bits set.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexconfigmapitemsindex">items</a></b></td>
        <td>[]object</td>
        <td>
          items if unspecified, each key-value pair in the Data field of the referenced
ConfigMap will be projected into the volume as a file whose name is the
key and content is the value. If specified, the listed keys will be
projected into the specified paths, and unlisted keys will not be
present. If a key is specified which is not present in the ConfigMap,
the volume setup will error unless it is marked optional. Paths must be
relative and may not contain the '..' path or start with '..'.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          optional specify whether the ConfigMap or its keys must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].configMap.items[index]
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexconfigmap)</sup></sup>



Maps a string key to a path within a volume.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the key to project.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          path is the relative path of the file to map the key to.
May not be an absolute path.
May not contain the path element '..'.
May not start with the string '..'.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>integer</td>
        <td>
          mode is Optional: mode bits used to set permissions on this file.
Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
If not specified, the volume defaultMode will be used.
This might be in conflict with other options that affect the file
mode, like fsGroup, and the result can be other mode bits set.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].csi
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



csi (Container Storage Interface) represents ephemeral storage that is handled by certain external CSI drivers.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>driver</b></td>
        <td>string</td>
        <td>
          driver is the name of the CSI driver that handles this volume.
Consult with your admin for the correct name as registered in the cluster.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType to mount. Ex. "ext4", "xfs", "ntfs".
If not provided, the empty value is passed to the associated CSI driver
which will determine the default filesystem to apply.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexcsinodepublishsecretref">nodePublishSecretRef</a></b></td>
        <td>object</td>
        <td>
          nodePublishSecretRef is a reference to the secret object containing
sensitive information to pass to the CSI driver to complete the CSI
NodePublishVolume and NodeUnpublishVolume calls.
This field is optional, and  may be empty if no secret is required. If the
secret object contains more than one secret, all secret references are passed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly specifies a read-only configuration for the volume.
Defaults to false (read/write).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>volumeAttributes</b></td>
        <td>map[string]string</td>
        <td>
          volumeAttributes stores driver-specific properties that are passed to the CSI
driver. Consult your driver's documentation for supported values.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].csi.nodePublishSecretRef
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexcsi)</sup></sup>



nodePublishSecretRef is a reference to the secret object containing
sensitive information to pass to the CSI driver to complete the CSI
NodePublishVolume and NodeUnpublishVolume calls.
This field is optional, and  may be empty if no secret is required. If the
secret object contains more than one secret, all secret references are passed.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].downwardAPI
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



downwardAPI represents downward API about the pod that should populate this volume

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>defaultMode</b></td>
        <td>integer</td>
        <td>
          Optional: mode bits to use on created files by default. Must be a
Optional: mode bits used to set permissions on created files by default.
Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
Defaults to 0644.
Directories within the path are not affected by this setting.
This might be in conflict with other options that affect the file
mode, like fsGroup, and the result can be other mode bits set.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexdownwardapiitemsindex">items</a></b></td>
        <td>[]object</td>
        <td>
          Items is a list of downward API volume file<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].downwardAPI.items[index]
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexdownwardapi)</sup></sup>



DownwardAPIVolumeFile represents information to create the file containing the pod field

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Required: Path is  the relative path name of the file to be created. Must not be absolute or contain the '..' path. Must be utf-8 encoded. The first item of the relative path must not start with '..'<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexdownwardapiitemsindexfieldref">fieldRef</a></b></td>
        <td>object</td>
        <td>
          Required: Selects a field of the pod: only annotations, labels, name, namespace and uid are supported.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>integer</td>
        <td>
          Optional: mode bits used to set permissions on this file, must be an octal value
between 0000 and 0777 or a decimal value between 0 and 511.
YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
If not specified, the volume defaultMode will be used.
This might be in conflict with other options that affect the file
mode, like fsGroup, and the result can be other mode bits set.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexdownwardapiitemsindexresourcefieldref">resourceFieldRef</a></b></td>
        <td>object</td>
        <td>
          Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.memory, requests.cpu and requests.memory) are currently supported.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].downwardAPI.items[index].fieldRef
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexdownwardapiitemsindex)</sup></sup>



Required: Selects a field of the pod: only annotations, labels, name, namespace and uid are supported.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>fieldPath</b></td>
        <td>string</td>
        <td>
          Path of the field to select in the specified API version.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>apiVersion</b></td>
        <td>string</td>
        <td>
          Version of the schema the FieldPath is written in terms of, defaults to "v1".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].downwardAPI.items[index].resourceFieldRef
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexdownwardapiitemsindex)</sup></sup>



Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.memory, requests.cpu and requests.memory) are currently supported.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resource</b></td>
        <td>string</td>
        <td>
          Required: resource to select<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>containerName</b></td>
        <td>string</td>
        <td>
          Container name: required for volumes, optional for env vars<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>divisor</b></td>
        <td>int or string</td>
        <td>
          Specifies the output format of the exposed resources, defaults to "1"<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].emptyDir
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



emptyDir represents a temporary directory that shares a pod's lifetime.
More info: https://kubernetes.io/docs/concepts/storage/volumes#emptydir

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>medium</b></td>
        <td>string</td>
        <td>
          medium represents what type of storage medium should back this directory.
The default is "" which means to use the node's default medium.
Must be an empty string (default) or Memory.
More info: https://kubernetes.io/docs/concepts/storage/volumes#emptydir<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sizeLimit</b></td>
        <td>int or string</td>
        <td>
          sizeLimit is the total amount of local storage required for this EmptyDir volume.
The size limit is also applicable for memory medium.
The maximum usage on memory medium EmptyDir would be the minimum value between
the SizeLimit specified here and the sum of memory limits of all containers in a pod.
The default is nil which means that the limit is undefined.
More info: https://kubernetes.io/docs/concepts/storage/volumes#emptydir<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].ephemeral
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



ephemeral represents a volume that is handled by a cluster storage driver.
The volume's lifecycle is tied to the pod that defines it - it will be created before the pod starts,
and deleted when the pod is removed.

Use this if:
a) the volume is only needed while the pod runs,
b) features of normal volumes like restoring from snapshot or capacity
   tracking are needed,
c) the storage driver is specified through a storage class, and
d) the storage driver supports dynamic volume provisioning through
   a PersistentVolumeClaim (see EphemeralVolumeSource for more
   information on the connection between this volume type
   and PersistentVolumeClaim).

Use PersistentVolumeClaim or one of the vendor-specific
APIs for volumes that persist for longer than the lifecycle
of an individual pod.

Use CSI for light-weight local ephemeral volumes if the CSI driver is meant to
be used that way - see the documentation of the driver for
more information.

A pod can use both types of ephemeral volumes and
persistent volumes at the same time.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexephemeralvolumeclaimtemplate">volumeClaimTemplate</a></b></td>
        <td>object</td>
        <td>
          Will be used to create a stand-alone PVC to provision the volume.
The pod in which this EphemeralVolumeSource is embedded will be the
owner of the PVC, i.e. the PVC will be deleted together with the
pod.  The name of the PVC will be `<pod name>-<volume name>` where
`<volume name>` is the name from the `PodSpec.Volumes` array
entry. Pod validation will reject the pod if the concatenated name
is not valid for a PVC (for example, too long).

An existing PVC with that name that is not owned by the pod
will *not* be used for the pod to avoid using an unrelated
volume by mistake. Starting the pod is then blocked until
the unrelated PVC is removed. If such a pre-created PVC is
meant to be used by the pod, the PVC has to updated with an
owner reference to the pod once the pod exists. Normally
this should not be necessary, but it may be useful when
manually reconstructing a broken cluster.

This field is read-only and no changes will be made by Kubernetes
to the PVC after it has been created.

Required, must not be nil.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].ephemeral.volumeClaimTemplate
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexephemeral)</sup></sup>



Will be used to create a stand-alone PVC to provision the volume.
The pod in which this EphemeralVolumeSource is embedded will be the
owner of the PVC, i.e. the PVC will be deleted together with the
pod.  The name of the PVC will be `<pod name>-<volume name>` where
`<volume name>` is the name from the `PodSpec.Volumes` array
entry. Pod validation will reject the pod if the concatenated name
is not valid for a PVC (for example, too long).

An existing PVC with that name that is not owned by the pod
will *not* be used for the pod to avoid using an unrelated
volume by mistake. Starting the pod is then blocked until
the unrelated PVC is removed. If such a pre-created PVC is
meant to be used by the pod, the PVC has to updated with an
owner reference to the pod once the pod exists. Normally
this should not be necessary, but it may be useful when
manually reconstructing a broken cluster.

This field is read-only and no changes will be made by Kubernetes
to the PVC after it has been created.

Required, must not be nil.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexephemeralvolumeclaimtemplatespec">spec</a></b></td>
        <td>object</td>
        <td>
          The specification for the PersistentVolumeClaim. The entire content is
copied unchanged into the PVC that gets created from this
template. The same fields as in a PersistentVolumeClaim
are also valid here.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>metadata</b></td>
        <td>object</td>
        <td>
          May contain labels and annotations that will be copied into the PVC
when creating it. No other fields are allowed and will be rejected during
validation.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].ephemeral.volumeClaimTemplate.spec
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexephemeralvolumeclaimtemplate)</sup></sup>



The specification for the PersistentVolumeClaim. The entire content is
copied unchanged into the PVC that gets created from this
template. The same fields as in a PersistentVolumeClaim
are also valid here.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>accessModes</b></td>
        <td>[]string</td>
        <td>
          accessModes contains the desired access modes the volume should have.
More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexephemeralvolumeclaimtemplatespecdatasource">dataSource</a></b></td>
        <td>object</td>
        <td>
          dataSource field can be used to specify either:
* An existing VolumeSnapshot object (snapshot.storage.k8s.io/VolumeSnapshot)
* An existing PVC (PersistentVolumeClaim)
If the provisioner or an external controller can support the specified data source,
it will create a new volume based on the contents of the specified data source.
When the AnyVolumeDataSource feature gate is enabled, dataSource contents will be copied to dataSourceRef,
and dataSourceRef contents will be copied to dataSource when dataSourceRef.namespace is not specified.
If the namespace is specified, then dataSourceRef will not be copied to dataSource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexephemeralvolumeclaimtemplatespecdatasourceref">dataSourceRef</a></b></td>
        <td>object</td>
        <td>
          dataSourceRef specifies the object from which to populate the volume with data, if a non-empty
volume is desired. This may be any object from a non-empty API group (non
core object) or a PersistentVolumeClaim object.
When this field is specified, volume binding will only succeed if the type of
the specified object matches some installed volume populator or dynamic
provisioner.
This field will replace the functionality of the dataSource field and as such
if both fields are non-empty, they must have the same value. For backwards
compatibility, when namespace isn't specified in dataSourceRef,
both fields (dataSource and dataSourceRef) will be set to the same
value automatically if one of them is empty and the other is non-empty.
When namespace is specified in dataSourceRef,
dataSource isn't set to the same value and must be empty.
There are three important differences between dataSource and dataSourceRef:
* While dataSource only allows two specific types of objects, dataSourceRef
  allows any non-core object, as well as PersistentVolumeClaim objects.
* While dataSource ignores disallowed values (dropping them), dataSourceRef
  preserves all values, and generates an error if a disallowed value is
  specified.
* While dataSource only allows local objects, dataSourceRef allows objects
  in any namespaces.
(Beta) Using this field requires the AnyVolumeDataSource feature gate to be enabled.
(Alpha) Using the namespace field of dataSourceRef requires the CrossNamespaceVolumeDataSource feature gate to be enabled.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexephemeralvolumeclaimtemplatespecresources">resources</a></b></td>
        <td>object</td>
        <td>
          resources represents the minimum resources the volume should have.
Users are allowed to specify resource requirements
that are lower than previous value but must still be higher than capacity recorded in the
status field of the claim.
More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#resources<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexephemeralvolumeclaimtemplatespecselector">selector</a></b></td>
        <td>object</td>
        <td>
          selector is a label query over volumes to consider for binding.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>storageClassName</b></td>
        <td>string</td>
        <td>
          storageClassName is the name of the StorageClass required by the claim.
More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#class-1<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>volumeAttributesClassName</b></td>
        <td>string</td>
        <td>
          volumeAttributesClassName may be used to set the VolumeAttributesClass used by this claim.
If specified, the CSI driver will create or update the volume with the attributes defined
in the corresponding VolumeAttributesClass. This has a different purpose than storageClassName,
it can be changed after the claim is created. An empty string or nil value indicates that no
VolumeAttributesClass will be applied to the claim. If the claim enters an Infeasible error state,
this field can be reset to its previous value (including nil) to cancel the modification.
If the resource referred to by volumeAttributesClass does not exist, this PersistentVolumeClaim will be
set to a Pending state, as reflected by the modifyVolumeStatus field, until such as a resource
exists.
More info: https://kubernetes.io/docs/concepts/storage/volume-attributes-classes/<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>volumeMode</b></td>
        <td>string</td>
        <td>
          volumeMode defines what type of volume is required by the claim.
Value of Filesystem is implied when not included in claim spec.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>volumeName</b></td>
        <td>string</td>
        <td>
          volumeName is the binding reference to the PersistentVolume backing this claim.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].ephemeral.volumeClaimTemplate.spec.dataSource
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexephemeralvolumeclaimtemplatespec)</sup></sup>



dataSource field can be used to specify either:
* An existing VolumeSnapshot object (snapshot.storage.k8s.io/VolumeSnapshot)
* An existing PVC (PersistentVolumeClaim)
If the provisioner or an external controller can support the specified data source,
it will create a new volume based on the contents of the specified data source.
When the AnyVolumeDataSource feature gate is enabled, dataSource contents will be copied to dataSourceRef,
and dataSourceRef contents will be copied to dataSource when dataSourceRef.namespace is not specified.
If the namespace is specified, then dataSourceRef will not be copied to dataSource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind is the type of resource being referenced<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name is the name of resource being referenced<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>apiGroup</b></td>
        <td>string</td>
        <td>
          APIGroup is the group for the resource being referenced.
If APIGroup is not specified, the specified Kind must be in the core API group.
For any other third-party types, APIGroup is required.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].ephemeral.volumeClaimTemplate.spec.dataSourceRef
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexephemeralvolumeclaimtemplatespec)</sup></sup>



dataSourceRef specifies the object from which to populate the volume with data, if a non-empty
volume is desired. This may be any object from a non-empty API group (non
core object) or a PersistentVolumeClaim object.
When this field is specified, volume binding will only succeed if the type of
the specified object matches some installed volume populator or dynamic
provisioner.
This field will replace the functionality of the dataSource field and as such
if both fields are non-empty, they must have the same value. For backwards
compatibility, when namespace isn't specified in dataSourceRef,
both fields (dataSource and dataSourceRef) will be set to the same
value automatically if one of them is empty and the other is non-empty.
When namespace is specified in dataSourceRef,
dataSource isn't set to the same value and must be empty.
There are three important differences between dataSource and dataSourceRef:
* While dataSource only allows two specific types of objects, dataSourceRef
  allows any non-core object, as well as PersistentVolumeClaim objects.
* While dataSource ignores disallowed values (dropping them), dataSourceRef
  preserves all values, and generates an error if a disallowed value is
  specified.
* While dataSource only allows local objects, dataSourceRef allows objects
  in any namespaces.
(Beta) Using this field requires the AnyVolumeDataSource feature gate to be enabled.
(Alpha) Using the namespace field of dataSourceRef requires the CrossNamespaceVolumeDataSource feature gate to be enabled.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>string</td>
        <td>
          Kind is the type of resource being referenced<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name is the name of resource being referenced<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>apiGroup</b></td>
        <td>string</td>
        <td>
          APIGroup is the group for the resource being referenced.
If APIGroup is not specified, the specified Kind must be in the core API group.
For any other third-party types, APIGroup is required.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace is the namespace of resource being referenced
Note that when a namespace is specified, a gateway.networking.k8s.io/ReferenceGrant object is required in the referent namespace to allow that namespace's owner to accept the reference. See the ReferenceGrant documentation for details.
(Alpha) This field requires the CrossNamespaceVolumeDataSource feature gate to be enabled.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].ephemeral.volumeClaimTemplate.spec.resources
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexephemeralvolumeclaimtemplatespec)</sup></sup>



resources represents the minimum resources the volume should have.
Users are allowed to specify resource requirements
that are lower than previous value but must still be higher than capacity recorded in the
status field of the claim.
More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#resources

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>limits</b></td>
        <td>map[string]int or string</td>
        <td>
          Limits describes the maximum amount of compute resources allowed.
More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>requests</b></td>
        <td>map[string]int or string</td>
        <td>
          Requests describes the minimum amount of compute resources required.
If Requests is omitted for a container, it defaults to Limits if that is explicitly specified,
otherwise to an implementation-defined value. Requests cannot exceed Limits.
More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].ephemeral.volumeClaimTemplate.spec.selector
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexephemeralvolumeclaimtemplatespec)</sup></sup>



selector is a label query over volumes to consider for binding.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexephemeralvolumeclaimtemplatespecselectormatchexpressionsindex">matchExpressions</a></b></td>
        <td>[]object</td>
        <td>
          matchExpressions is a list of label selector requirements. The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels
map is equivalent to an element of matchExpressions, whose key field is "key", the
operator is "In", and the values array contains only "value". The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].ephemeral.volumeClaimTemplate.spec.selector.matchExpressions[index]
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexephemeralvolumeclaimtemplatespecselector)</sup></sup>



A label selector requirement is a selector that contains values, a key, and an operator that
relates the key and values.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          operator represents a key's relationship to a set of values.
Valid operators are In, NotIn, Exists and DoesNotExist.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          values is an array of string values. If the operator is In or NotIn,
the values array must be non-empty. If the operator is Exists or DoesNotExist,
the values array must be empty. This array is replaced during a strategic
merge patch.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].fc
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



fc represents a Fibre Channel resource that is attached to a kubelet's host machine and then exposed to the pod.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type to mount.
Must be a filesystem type supported by the host operating system.
Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>lun</b></td>
        <td>integer</td>
        <td>
          lun is Optional: FC target lun number<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly is Optional: Defaults to false (read/write). ReadOnly here will force
the ReadOnly setting in VolumeMounts.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>targetWWNs</b></td>
        <td>[]string</td>
        <td>
          targetWWNs is Optional: FC target worldwide names (WWNs)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>wwids</b></td>
        <td>[]string</td>
        <td>
          wwids Optional: FC volume world wide identifiers (wwids)
Either wwids or combination of targetWWNs and lun must be set, but not both simultaneously.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].flexVolume
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



flexVolume represents a generic volume resource that is
provisioned/attached using an exec based plugin.
Deprecated: FlexVolume is deprecated. Consider using a CSIDriver instead.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>driver</b></td>
        <td>string</td>
        <td>
          driver is the name of the driver to use for this volume.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type to mount.
Must be a filesystem type supported by the host operating system.
Ex. "ext4", "xfs", "ntfs". The default filesystem depends on FlexVolume script.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>options</b></td>
        <td>map[string]string</td>
        <td>
          options is Optional: this field holds extra command options if any.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly is Optional: defaults to false (read/write). ReadOnly here will force
the ReadOnly setting in VolumeMounts.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexflexvolumesecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          secretRef is Optional: secretRef is reference to the secret object containing
sensitive information to pass to the plugin scripts. This may be
empty if no secret object is specified. If the secret object
contains more than one secret, all secrets are passed to the plugin
scripts.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].flexVolume.secretRef
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexflexvolume)</sup></sup>



secretRef is Optional: secretRef is reference to the secret object containing
sensitive information to pass to the plugin scripts. This may be
empty if no secret object is specified. If the secret object
contains more than one secret, all secrets are passed to the plugin
scripts.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].flocker
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



flocker represents a Flocker volume attached to a kubelet's host machine. This depends on the Flocker control service being running.
Deprecated: Flocker is deprecated and the in-tree flocker type is no longer supported.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>datasetName</b></td>
        <td>string</td>
        <td>
          datasetName is Name of the dataset stored as metadata -> name on the dataset for Flocker
should be considered as deprecated<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>datasetUUID</b></td>
        <td>string</td>
        <td>
          datasetUUID is the UUID of the dataset. This is unique identifier of a Flocker dataset<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].gcePersistentDisk
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



gcePersistentDisk represents a GCE Disk resource that is attached to a
kubelet's host machine and then exposed to the pod.
Deprecated: GCEPersistentDisk is deprecated. All operations for the in-tree
gcePersistentDisk type are redirected to the pd.csi.storage.gke.io CSI driver.
More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>pdName</b></td>
        <td>string</td>
        <td>
          pdName is unique name of the PD resource in GCE. Used to identify the disk in GCE.
More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is filesystem type of the volume that you want to mount.
Tip: Ensure that the filesystem type is supported by the host operating system.
Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>partition</b></td>
        <td>integer</td>
        <td>
          partition is the partition in the volume that you want to mount.
If omitted, the default is to mount by volume name.
Examples: For volume /dev/sda1, you specify the partition as "1".
Similarly, the volume partition for /dev/sda is "0" (or you can leave the property empty).
More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly here will force the ReadOnly setting in VolumeMounts.
Defaults to false.
More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].gitRepo
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



gitRepo represents a git repository at a particular revision.
Deprecated: GitRepo is deprecated. To provision a container with a git repo, mount an
EmptyDir into an InitContainer that clones the repo using git, then mount the EmptyDir
into the Pod's container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>repository</b></td>
        <td>string</td>
        <td>
          repository is the URL<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>directory</b></td>
        <td>string</td>
        <td>
          directory is the target directory name.
Must not contain or start with '..'.  If '.' is supplied, the volume directory will be the
git repository.  Otherwise, if specified, the volume will contain the git repository in
the subdirectory with the given name.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>revision</b></td>
        <td>string</td>
        <td>
          revision is the commit hash for the specified revision.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].glusterfs
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



glusterfs represents a Glusterfs mount on the host that shares a pod's lifetime.
Deprecated: Glusterfs is deprecated and the in-tree glusterfs type is no longer supported.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>endpoints</b></td>
        <td>string</td>
        <td>
          endpoints is the endpoint name that details Glusterfs topology.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          path is the Glusterfs volume path.
More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly here will force the Glusterfs volume to be mounted with read-only permissions.
Defaults to false.
More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].hostPath
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



hostPath represents a pre-existing file or directory on the host
machine that is directly exposed to the container. This is generally
used for system agents or other privileged things that are allowed
to see the host machine. Most containers will NOT need this.
More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          path of the directory on the host.
If the path is a symlink, it will follow the link to the real path.
More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          type for HostPath Volume
Defaults to ""
More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].image
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



image represents an OCI object (a container image or artifact) pulled and mounted on the kubelet's host machine.
The volume is resolved at pod startup depending on which PullPolicy value is provided:

- Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.
- Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.
- IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.

The volume gets re-resolved if the pod gets deleted and recreated, which means that new remote content will become available on pod recreation.
A failure to resolve or pull the image during pod startup will block containers from starting and may add significant latency. Failures will be retried using normal volume backoff and will be reported on the pod reason and message.
The types of objects that may be mounted by this volume are defined by the container runtime implementation on a host machine and at minimum must include all valid types supported by the container image field.
The OCI object gets mounted in a single directory (spec.containers[*].volumeMounts.mountPath) by merging the manifest layers in the same way as for container images.
The volume will be mounted read-only (ro) and non-executable files (noexec).
Sub path mounts for containers are not supported (spec.containers[*].volumeMounts.subpath) before 1.33.
The field spec.securityContext.fsGroupChangePolicy has no effect on this volume type.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>pullPolicy</b></td>
        <td>string</td>
        <td>
          Policy for pulling OCI objects. Possible values are:
Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.
Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.
IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.
Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reference</b></td>
        <td>string</td>
        <td>
          Required: Image or artifact reference to be used.
Behaves in the same way as pod.spec.containers[*].image.
Pull secrets will be assembled in the same way as for the container image by looking up node credentials, SA image pull secrets, and pod spec image pull secrets.
More info: https://kubernetes.io/docs/concepts/containers/images
This field is optional to allow higher level config management to default or override
container images in workload controllers like Deployments and StatefulSets.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].iscsi
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



iscsi represents an ISCSI Disk resource that is attached to a
kubelet's host machine and then exposed to the pod.
More info: https://kubernetes.io/docs/concepts/storage/volumes/#iscsi

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>iqn</b></td>
        <td>string</td>
        <td>
          iqn is the target iSCSI Qualified Name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>lun</b></td>
        <td>integer</td>
        <td>
          lun represents iSCSI Target Lun number.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>targetPortal</b></td>
        <td>string</td>
        <td>
          targetPortal is iSCSI Target Portal. The Portal is either an IP or ip_addr:port if the port
is other than default (typically TCP ports 860 and 3260).<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>chapAuthDiscovery</b></td>
        <td>boolean</td>
        <td>
          chapAuthDiscovery defines whether support iSCSI Discovery CHAP authentication<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>chapAuthSession</b></td>
        <td>boolean</td>
        <td>
          chapAuthSession defines whether support iSCSI Session CHAP authentication<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type of the volume that you want to mount.
Tip: Ensure that the filesystem type is supported by the host operating system.
Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
More info: https://kubernetes.io/docs/concepts/storage/volumes#iscsi<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initiatorName</b></td>
        <td>string</td>
        <td>
          initiatorName is the custom iSCSI Initiator Name.
If initiatorName is specified with iscsiInterface simultaneously, new iSCSI interface
<target portal>:<volume name> will be created for the connection.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>iscsiInterface</b></td>
        <td>string</td>
        <td>
          iscsiInterface is the interface Name that uses an iSCSI transport.
Defaults to 'default' (tcp).<br/>
          <br/>
            <i>Default</i>: default<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>portals</b></td>
        <td>[]string</td>
        <td>
          portals is the iSCSI Target Portal List. The portal is either an IP or ip_addr:port if the port
is other than default (typically TCP ports 860 and 3260).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly here will force the ReadOnly setting in VolumeMounts.
Defaults to false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexiscsisecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          secretRef is the CHAP Secret for iSCSI target and initiator authentication<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].iscsi.secretRef
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexiscsi)</sup></sup>



secretRef is the CHAP Secret for iSCSI target and initiator authentication

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].nfs
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



nfs represents an NFS mount on the host that shares a pod's lifetime
More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          path that is exported by the NFS server.
More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>server</b></td>
        <td>string</td>
        <td>
          server is the hostname or IP address of the NFS server.
More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly here will force the NFS export to be mounted with read-only permissions.
Defaults to false.
More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].persistentVolumeClaim
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



persistentVolumeClaimVolumeSource represents a reference to a
PersistentVolumeClaim in the same namespace.
More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>claimName</b></td>
        <td>string</td>
        <td>
          claimName is the name of a PersistentVolumeClaim in the same namespace as the pod using this volume.
More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly Will force the ReadOnly setting in VolumeMounts.
Default false.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].photonPersistentDisk
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



photonPersistentDisk represents a PhotonController persistent disk attached and mounted on kubelets host machine.
Deprecated: PhotonPersistentDisk is deprecated and the in-tree photonPersistentDisk type is no longer supported.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>pdID</b></td>
        <td>string</td>
        <td>
          pdID is the ID that identifies Photon Controller persistent disk<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type to mount.
Must be a filesystem type supported by the host operating system.
Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].portworxVolume
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



portworxVolume represents a portworx volume attached and mounted on kubelets host machine.
Deprecated: PortworxVolume is deprecated. All operations for the in-tree portworxVolume type
are redirected to the pxd.portworx.com CSI driver when the CSIMigrationPortworx feature-gate
is on.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>volumeID</b></td>
        <td>string</td>
        <td>
          volumeID uniquely identifies a Portworx volume<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fSType represents the filesystem type to mount
Must be a filesystem type supported by the host operating system.
Ex. "ext4", "xfs". Implicitly inferred to be "ext4" if unspecified.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly defaults to false (read/write). ReadOnly here will force
the ReadOnly setting in VolumeMounts.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].projected
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



projected items for all in one resources secrets, configmaps, and downward API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>defaultMode</b></td>
        <td>integer</td>
        <td>
          defaultMode are the mode bits used to set permissions on created files by default.
Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
Directories within the path are not affected by this setting.
This might be in conflict with other options that affect the file
mode, like fsGroup, and the result can be other mode bits set.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexprojectedsourcesindex">sources</a></b></td>
        <td>[]object</td>
        <td>
          sources is the list of volume projections. Each entry in this list
handles one source.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].projected.sources[index]
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexprojected)</sup></sup>



Projection that may be projected along with other supported volume types.
Exactly one of these fields must be set.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexprojectedsourcesindexclustertrustbundle">clusterTrustBundle</a></b></td>
        <td>object</td>
        <td>
          ClusterTrustBundle allows a pod to access the `.spec.trustBundle` field
of ClusterTrustBundle objects in an auto-updating file.

Alpha, gated by the ClusterTrustBundleProjection feature gate.

ClusterTrustBundle objects can either be selected by name, or by the
combination of signer name and a label selector.

Kubelet performs aggressive normalization of the PEM contents written
into the pod filesystem.  Esoteric PEM features such as inter-block
comments and block headers are stripped.  Certificates are deduplicated.
The ordering of certificates within the file is arbitrary, and Kubelet
may change the order over time.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexprojectedsourcesindexconfigmap">configMap</a></b></td>
        <td>object</td>
        <td>
          configMap information about the configMap data to project<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexprojectedsourcesindexdownwardapi">downwardAPI</a></b></td>
        <td>object</td>
        <td>
          downwardAPI information about the downwardAPI data to project<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexprojectedsourcesindexpodcertificate">podCertificate</a></b></td>
        <td>object</td>
        <td>
          Projects an auto-rotating credential bundle (private key and certificate
chain) that the pod can use either as a TLS client or server.

Kubelet generates a private key and uses it to send a
PodCertificateRequest to the named signer.  Once the signer approves the
request and issues a certificate chain, Kubelet writes the key and
certificate chain to the pod filesystem.  The pod does not start until
certificates have been issued for each podCertificate projected volume
source in its spec.

Kubelet will begin trying to rotate the certificate at the time indicated
by the signer using the PodCertificateRequest.Status.BeginRefreshAt
timestamp.

Kubelet can write a single file, indicated by the credentialBundlePath
field, or separate files, indicated by the keyPath and
certificateChainPath fields.

The credential bundle is a single file in PEM format.  The first PEM
entry is the private key (in PKCS#8 format), and the remaining PEM
entries are the certificate chain issued by the signer (typically,
signers will return their certificate chain in leaf-to-root order).

Prefer using the credential bundle format, since your application code
can read it atomically.  If you use keyPath and certificateChainPath,
your application must make two separate file reads. If these coincide
with a certificate rotation, it is possible that the private key and leaf
certificate you read may not correspond to each other.  Your application
will need to check for this condition, and re-read until they are
consistent.

The named signer controls chooses the format of the certificate it
issues; consult the signer implementation's documentation to learn how to
use the certificates it issues.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexprojectedsourcesindexsecret">secret</a></b></td>
        <td>object</td>
        <td>
          secret information about the secret data to project<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexprojectedsourcesindexserviceaccounttoken">serviceAccountToken</a></b></td>
        <td>object</td>
        <td>
          serviceAccountToken is information about the serviceAccountToken data to project<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].projected.sources[index].clusterTrustBundle
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexprojectedsourcesindex)</sup></sup>



ClusterTrustBundle allows a pod to access the `.spec.trustBundle` field
of ClusterTrustBundle objects in an auto-updating file.

Alpha, gated by the ClusterTrustBundleProjection feature gate.

ClusterTrustBundle objects can either be selected by name, or by the
combination of signer name and a label selector.

Kubelet performs aggressive normalization of the PEM contents written
into the pod filesystem.  Esoteric PEM features such as inter-block
comments and block headers are stripped.  Certificates are deduplicated.
The ordering of certificates within the file is arbitrary, and Kubelet
may change the order over time.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Relative path from the volume root to write the bundle.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexprojectedsourcesindexclustertrustbundlelabelselector">labelSelector</a></b></td>
        <td>object</td>
        <td>
          Select all ClusterTrustBundles that match this label selector.  Only has
effect if signerName is set.  Mutually-exclusive with name.  If unset,
interpreted as "match nothing".  If set but empty, interpreted as "match
everything".<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Select a single ClusterTrustBundle by object name.  Mutually-exclusive
with signerName and labelSelector.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          If true, don't block pod startup if the referenced ClusterTrustBundle(s)
aren't available.  If using name, then the named ClusterTrustBundle is
allowed not to exist.  If using signerName, then the combination of
signerName and labelSelector is allowed to match zero
ClusterTrustBundles.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>signerName</b></td>
        <td>string</td>
        <td>
          Select all ClusterTrustBundles that match this signer name.
Mutually-exclusive with name.  The contents of all selected
ClusterTrustBundles will be unified and deduplicated.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].projected.sources[index].clusterTrustBundle.labelSelector
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexprojectedsourcesindexclustertrustbundle)</sup></sup>



Select all ClusterTrustBundles that match this label selector.  Only has
effect if signerName is set.  Mutually-exclusive with name.  If unset,
interpreted as "match nothing".  If set but empty, interpreted as "match
everything".

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexprojectedsourcesindexclustertrustbundlelabelselectormatchexpressionsindex">matchExpressions</a></b></td>
        <td>[]object</td>
        <td>
          matchExpressions is a list of label selector requirements. The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>matchLabels</b></td>
        <td>map[string]string</td>
        <td>
          matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels
map is equivalent to an element of matchExpressions, whose key field is "key", the
operator is "In", and the values array contains only "value". The requirements are ANDed.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].projected.sources[index].clusterTrustBundle.labelSelector.matchExpressions[index]
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexprojectedsourcesindexclustertrustbundlelabelselector)</sup></sup>



A label selector requirement is a selector that contains values, a key, and an operator that
relates the key and values.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the label key that the selector applies to.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>operator</b></td>
        <td>string</td>
        <td>
          operator represents a key's relationship to a set of values.
Valid operators are In, NotIn, Exists and DoesNotExist.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>values</b></td>
        <td>[]string</td>
        <td>
          values is an array of string values. If the operator is In or NotIn,
the values array must be non-empty. If the operator is Exists or DoesNotExist,
the values array must be empty. This array is replaced during a strategic
merge patch.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].projected.sources[index].configMap
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexprojectedsourcesindex)</sup></sup>



configMap information about the configMap data to project

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexprojectedsourcesindexconfigmapitemsindex">items</a></b></td>
        <td>[]object</td>
        <td>
          items if unspecified, each key-value pair in the Data field of the referenced
ConfigMap will be projected into the volume as a file whose name is the
key and content is the value. If specified, the listed keys will be
projected into the specified paths, and unlisted keys will not be
present. If a key is specified which is not present in the ConfigMap,
the volume setup will error unless it is marked optional. Paths must be
relative and may not contain the '..' path or start with '..'.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          optional specify whether the ConfigMap or its keys must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].projected.sources[index].configMap.items[index]
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexprojectedsourcesindexconfigmap)</sup></sup>



Maps a string key to a path within a volume.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the key to project.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          path is the relative path of the file to map the key to.
May not be an absolute path.
May not contain the path element '..'.
May not start with the string '..'.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>integer</td>
        <td>
          mode is Optional: mode bits used to set permissions on this file.
Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
If not specified, the volume defaultMode will be used.
This might be in conflict with other options that affect the file
mode, like fsGroup, and the result can be other mode bits set.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].projected.sources[index].downwardAPI
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexprojectedsourcesindex)</sup></sup>



downwardAPI information about the downwardAPI data to project

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexprojectedsourcesindexdownwardapiitemsindex">items</a></b></td>
        <td>[]object</td>
        <td>
          Items is a list of DownwardAPIVolume file<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].projected.sources[index].downwardAPI.items[index]
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexprojectedsourcesindexdownwardapi)</sup></sup>



DownwardAPIVolumeFile represents information to create the file containing the pod field

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          Required: Path is  the relative path name of the file to be created. Must not be absolute or contain the '..' path. Must be utf-8 encoded. The first item of the relative path must not start with '..'<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexprojectedsourcesindexdownwardapiitemsindexfieldref">fieldRef</a></b></td>
        <td>object</td>
        <td>
          Required: Selects a field of the pod: only annotations, labels, name, namespace and uid are supported.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>integer</td>
        <td>
          Optional: mode bits used to set permissions on this file, must be an octal value
between 0000 and 0777 or a decimal value between 0 and 511.
YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
If not specified, the volume defaultMode will be used.
This might be in conflict with other options that affect the file
mode, like fsGroup, and the result can be other mode bits set.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexprojectedsourcesindexdownwardapiitemsindexresourcefieldref">resourceFieldRef</a></b></td>
        <td>object</td>
        <td>
          Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.memory, requests.cpu and requests.memory) are currently supported.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].projected.sources[index].downwardAPI.items[index].fieldRef
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexprojectedsourcesindexdownwardapiitemsindex)</sup></sup>



Required: Selects a field of the pod: only annotations, labels, name, namespace and uid are supported.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>fieldPath</b></td>
        <td>string</td>
        <td>
          Path of the field to select in the specified API version.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>apiVersion</b></td>
        <td>string</td>
        <td>
          Version of the schema the FieldPath is written in terms of, defaults to "v1".<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].projected.sources[index].downwardAPI.items[index].resourceFieldRef
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexprojectedsourcesindexdownwardapiitemsindex)</sup></sup>



Selects a resource of the container: only resources limits and requests
(limits.cpu, limits.memory, requests.cpu and requests.memory) are currently supported.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>resource</b></td>
        <td>string</td>
        <td>
          Required: resource to select<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>containerName</b></td>
        <td>string</td>
        <td>
          Container name: required for volumes, optional for env vars<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>divisor</b></td>
        <td>int or string</td>
        <td>
          Specifies the output format of the exposed resources, defaults to "1"<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].projected.sources[index].podCertificate
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexprojectedsourcesindex)</sup></sup>



Projects an auto-rotating credential bundle (private key and certificate
chain) that the pod can use either as a TLS client or server.

Kubelet generates a private key and uses it to send a
PodCertificateRequest to the named signer.  Once the signer approves the
request and issues a certificate chain, Kubelet writes the key and
certificate chain to the pod filesystem.  The pod does not start until
certificates have been issued for each podCertificate projected volume
source in its spec.

Kubelet will begin trying to rotate the certificate at the time indicated
by the signer using the PodCertificateRequest.Status.BeginRefreshAt
timestamp.

Kubelet can write a single file, indicated by the credentialBundlePath
field, or separate files, indicated by the keyPath and
certificateChainPath fields.

The credential bundle is a single file in PEM format.  The first PEM
entry is the private key (in PKCS#8 format), and the remaining PEM
entries are the certificate chain issued by the signer (typically,
signers will return their certificate chain in leaf-to-root order).

Prefer using the credential bundle format, since your application code
can read it atomically.  If you use keyPath and certificateChainPath,
your application must make two separate file reads. If these coincide
with a certificate rotation, it is possible that the private key and leaf
certificate you read may not correspond to each other.  Your application
will need to check for this condition, and re-read until they are
consistent.

The named signer controls chooses the format of the certificate it
issues; consult the signer implementation's documentation to learn how to
use the certificates it issues.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>keyType</b></td>
        <td>string</td>
        <td>
          The type of keypair Kubelet will generate for the pod.

Valid values are "RSA3072", "RSA4096", "ECDSAP256", "ECDSAP384",
"ECDSAP521", and "ED25519".<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>signerName</b></td>
        <td>string</td>
        <td>
          Kubelet's generated CSRs will be addressed to this signer.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>certificateChainPath</b></td>
        <td>string</td>
        <td>
          Write the certificate chain at this path in the projected volume.

Most applications should use credentialBundlePath.  When using keyPath
and certificateChainPath, your application needs to check that the key
and leaf certificate are consistent, because it is possible to read the
files mid-rotation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>credentialBundlePath</b></td>
        <td>string</td>
        <td>
          Write the credential bundle at this path in the projected volume.

The credential bundle is a single file that contains multiple PEM blocks.
The first PEM block is a PRIVATE KEY block, containing a PKCS#8 private
key.

The remaining blocks are CERTIFICATE blocks, containing the issued
certificate chain from the signer (leaf and any intermediates).

Using credentialBundlePath lets your Pod's application code make a single
atomic read that retrieves a consistent key and certificate chain.  If you
project them to separate files, your application code will need to
additionally check that the leaf certificate was issued to the key.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keyPath</b></td>
        <td>string</td>
        <td>
          Write the key at this path in the projected volume.

Most applications should use credentialBundlePath.  When using keyPath
and certificateChainPath, your application needs to check that the key
and leaf certificate are consistent, because it is possible to read the
files mid-rotation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>maxExpirationSeconds</b></td>
        <td>integer</td>
        <td>
          maxExpirationSeconds is the maximum lifetime permitted for the
certificate.

Kubelet copies this value verbatim into the PodCertificateRequests it
generates for this projection.

If omitted, kube-apiserver will set it to 86400(24 hours). kube-apiserver
will reject values shorter than 3600 (1 hour).  The maximum allowable
value is 7862400 (91 days).

The signer implementation is then free to issue a certificate with any
lifetime *shorter* than MaxExpirationSeconds, but no shorter than 3600
seconds (1 hour).  This constraint is enforced by kube-apiserver.
`kubernetes.io` signers will never issue certificates with a lifetime
longer than 24 hours.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>userAnnotations</b></td>
        <td>map[string]string</td>
        <td>
          userAnnotations allow pod authors to pass additional information to
the signer implementation.  Kubernetes does not restrict or validate this
metadata in any way.

These values are copied verbatim into the `spec.unverifiedUserAnnotations` field of
the PodCertificateRequest objects that Kubelet creates.

Entries are subject to the same validation as object metadata annotations,
with the addition that all keys must be domain-prefixed. No restrictions
are placed on values, except an overall size limitation on the entire field.

Signers should document the keys and values they support. Signers should
deny requests that contain keys they do not recognize.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].projected.sources[index].secret
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexprojectedsourcesindex)</sup></sup>



secret information about the secret data to project

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexprojectedsourcesindexsecretitemsindex">items</a></b></td>
        <td>[]object</td>
        <td>
          items if unspecified, each key-value pair in the Data field of the referenced
Secret will be projected into the volume as a file whose name is the
key and content is the value. If specified, the listed keys will be
projected into the specified paths, and unlisted keys will not be
present. If a key is specified which is not present in the Secret,
the volume setup will error unless it is marked optional. Paths must be
relative and may not contain the '..' path or start with '..'.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          optional field specify whether the Secret or its key must be defined<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].projected.sources[index].secret.items[index]
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexprojectedsourcesindexsecret)</sup></sup>



Maps a string key to a path within a volume.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the key to project.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          path is the relative path of the file to map the key to.
May not be an absolute path.
May not contain the path element '..'.
May not start with the string '..'.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>integer</td>
        <td>
          mode is Optional: mode bits used to set permissions on this file.
Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
If not specified, the volume defaultMode will be used.
This might be in conflict with other options that affect the file
mode, like fsGroup, and the result can be other mode bits set.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].projected.sources[index].serviceAccountToken
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexprojectedsourcesindex)</sup></sup>



serviceAccountToken is information about the serviceAccountToken data to project

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          path is the path relative to the mount point of the file to project the
token into.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>audience</b></td>
        <td>string</td>
        <td>
          audience is the intended audience of the token. A recipient of a token
must identify itself with an identifier specified in the audience of the
token, and otherwise should reject the token. The audience defaults to the
identifier of the apiserver.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>expirationSeconds</b></td>
        <td>integer</td>
        <td>
          expirationSeconds is the requested duration of validity of the service
account token. As the token approaches expiration, the kubelet volume
plugin will proactively rotate the service account token. The kubelet will
start trying to rotate the token if the token is older than 80 percent of
its time to live or if the token is older than 24 hours.Defaults to 1 hour
and must be at least 10 minutes.<br/>
          <br/>
            <i>Format</i>: int64<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].quobyte
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



quobyte represents a Quobyte mount on the host that shares a pod's lifetime.
Deprecated: Quobyte is deprecated and the in-tree quobyte type is no longer supported.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>registry</b></td>
        <td>string</td>
        <td>
          registry represents a single or multiple Quobyte Registry services
specified as a string as host:port pair (multiple entries are separated with commas)
which acts as the central registry for volumes<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>volume</b></td>
        <td>string</td>
        <td>
          volume is a string that references an already created Quobyte volume by name.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>group</b></td>
        <td>string</td>
        <td>
          group to map volume access to
Default is no group<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly here will force the Quobyte volume to be mounted with read-only permissions.
Defaults to false.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>tenant</b></td>
        <td>string</td>
        <td>
          tenant owning the given Quobyte volume in the Backend
Used with dynamically provisioned Quobyte volumes, value is set by the plugin<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>user</b></td>
        <td>string</td>
        <td>
          user to map volume access to
Defaults to serivceaccount user<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].rbd
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



rbd represents a Rados Block Device mount on the host that shares a pod's lifetime.
Deprecated: RBD is deprecated and the in-tree rbd type is no longer supported.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>image</b></td>
        <td>string</td>
        <td>
          image is the rados image name.
More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>monitors</b></td>
        <td>[]string</td>
        <td>
          monitors is a collection of Ceph monitors.
More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type of the volume that you want to mount.
Tip: Ensure that the filesystem type is supported by the host operating system.
Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
More info: https://kubernetes.io/docs/concepts/storage/volumes#rbd<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keyring</b></td>
        <td>string</td>
        <td>
          keyring is the path to key ring for RBDUser.
Default is /etc/ceph/keyring.
More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it<br/>
          <br/>
            <i>Default</i>: /etc/ceph/keyring<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>pool</b></td>
        <td>string</td>
        <td>
          pool is the rados pool name.
Default is rbd.
More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it<br/>
          <br/>
            <i>Default</i>: rbd<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly here will force the ReadOnly setting in VolumeMounts.
Defaults to false.
More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexrbdsecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          secretRef is name of the authentication secret for RBDUser. If provided
overrides keyring.
Default is nil.
More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>user</b></td>
        <td>string</td>
        <td>
          user is the rados user name.
Default is admin.
More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it<br/>
          <br/>
            <i>Default</i>: admin<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].rbd.secretRef
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexrbd)</sup></sup>



secretRef is name of the authentication secret for RBDUser. If provided
overrides keyring.
Default is nil.
More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].scaleIO
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



scaleIO represents a ScaleIO persistent volume attached and mounted on Kubernetes nodes.
Deprecated: ScaleIO is deprecated and the in-tree scaleIO type is no longer supported.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>gateway</b></td>
        <td>string</td>
        <td>
          gateway is the host address of the ScaleIO API Gateway.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexscaleiosecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          secretRef references to the secret for ScaleIO user and other
sensitive information. If this is not provided, Login operation will fail.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>system</b></td>
        <td>string</td>
        <td>
          system is the name of the storage system as configured in ScaleIO.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type to mount.
Must be a filesystem type supported by the host operating system.
Ex. "ext4", "xfs", "ntfs".
Default is "xfs".<br/>
          <br/>
            <i>Default</i>: xfs<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>protectionDomain</b></td>
        <td>string</td>
        <td>
          protectionDomain is the name of the ScaleIO Protection Domain for the configured storage.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly Defaults to false (read/write). ReadOnly here will force
the ReadOnly setting in VolumeMounts.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sslEnabled</b></td>
        <td>boolean</td>
        <td>
          sslEnabled Flag enable/disable SSL communication with Gateway, default false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>storageMode</b></td>
        <td>string</td>
        <td>
          storageMode indicates whether the storage for a volume should be ThickProvisioned or ThinProvisioned.
Default is ThinProvisioned.<br/>
          <br/>
            <i>Default</i>: ThinProvisioned<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>storagePool</b></td>
        <td>string</td>
        <td>
          storagePool is the ScaleIO Storage Pool associated with the protection domain.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>volumeName</b></td>
        <td>string</td>
        <td>
          volumeName is the name of a volume already created in the ScaleIO system
that is associated with this volume source.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].scaleIO.secretRef
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexscaleio)</sup></sup>



secretRef references to the secret for ScaleIO user and other
sensitive information. If this is not provided, Login operation will fail.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].secret
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



secret represents a secret that should populate this volume.
More info: https://kubernetes.io/docs/concepts/storage/volumes#secret

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>defaultMode</b></td>
        <td>integer</td>
        <td>
          defaultMode is Optional: mode bits used to set permissions on created files by default.
Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
YAML accepts both octal and decimal values, JSON requires decimal values
for mode bits. Defaults to 0644.
Directories within the path are not affected by this setting.
This might be in conflict with other options that affect the file
mode, like fsGroup, and the result can be other mode bits set.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexsecretitemsindex">items</a></b></td>
        <td>[]object</td>
        <td>
          items If unspecified, each key-value pair in the Data field of the referenced
Secret will be projected into the volume as a file whose name is the
key and content is the value. If specified, the listed keys will be
projected into the specified paths, and unlisted keys will not be
present. If a key is specified which is not present in the Secret,
the volume setup will error unless it is marked optional. Paths must be
relative and may not contain the '..' path or start with '..'.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>optional</b></td>
        <td>boolean</td>
        <td>
          optional field specify whether the Secret or its keys must be defined<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>secretName</b></td>
        <td>string</td>
        <td>
          secretName is the name of the secret in the pod's namespace to use.
More info: https://kubernetes.io/docs/concepts/storage/volumes#secret<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].secret.items[index]
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexsecret)</sup></sup>



Maps a string key to a path within a volume.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          key is the key to project.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>path</b></td>
        <td>string</td>
        <td>
          path is the relative path of the file to map the key to.
May not be an absolute path.
May not contain the path element '..'.
May not start with the string '..'.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>mode</b></td>
        <td>integer</td>
        <td>
          mode is Optional: mode bits used to set permissions on this file.
Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
If not specified, the volume defaultMode will be used.
This might be in conflict with other options that affect the file
mode, like fsGroup, and the result can be other mode bits set.<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].storageos
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



storageOS represents a StorageOS volume attached and mounted on Kubernetes nodes.
Deprecated: StorageOS is deprecated and the in-tree storageos type is no longer supported.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is the filesystem type to mount.
Must be a filesystem type supported by the host operating system.
Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>readOnly</b></td>
        <td>boolean</td>
        <td>
          readOnly defaults to false (read/write). ReadOnly here will force
the ReadOnly setting in VolumeMounts.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecoverridesvolumesindexstorageossecretref">secretRef</a></b></td>
        <td>object</td>
        <td>
          secretRef specifies the secret to use for obtaining the StorageOS API
credentials.  If not specified, default values will be attempted.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>volumeName</b></td>
        <td>string</td>
        <td>
          volumeName is the human-readable name of the StorageOS volume.  Volume
names are only unique within a namespace.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>volumeNamespace</b></td>
        <td>string</td>
        <td>
          volumeNamespace specifies the scope of the volume within StorageOS.  If no
namespace is specified then the Pod's namespace will be used.  This allows the
Kubernetes name scoping to be mirrored within StorageOS for tighter integration.
Set VolumeName to any name to override the default behaviour.
Set to "default" if you are not using namespaces within StorageOS.
Namespaces that do not pre-exist within StorageOS will be created.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].storageos.secretRef
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindexstorageos)</sup></sup>



secretRef specifies the secret to use for obtaining the StorageOS API
credentials.  If not specified, default values will be attempted.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the referent.
This field is effectively required, but due to backwards compatibility is
allowed to be empty. Instances of this type with an empty value here are
almost certainly wrong.
More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names<br/>
          <br/>
            <i>Default</i>: <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.overrides.volumes[index].vsphereVolume
<sup><sup>[↩ Parent](#ironicspecoverridesvolumesindex)</sup></sup>



vsphereVolume represents a vSphere volume attached and mounted on kubelets host machine.
Deprecated: VsphereVolume is deprecated. All operations for the in-tree vsphereVolume type
are redirected to the csi.vsphere.vmware.com CSI driver.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>volumePath</b></td>
        <td>string</td>
        <td>
          volumePath is the path that identifies vSphere volume vmdk<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>fsType</b></td>
        <td>string</td>
        <td>
          fsType is filesystem type to mount.
Must be a filesystem type supported by the host operating system.
Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>storagePolicyID</b></td>
        <td>string</td>
        <td>
          storagePolicyID is the storage Policy Based Management (SPBM) profile ID associated with the StoragePolicyName.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>storagePolicyName</b></td>
        <td>string</td>
        <td>
          storagePolicyName is the storage Policy Based Management (SPBM) profile name.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.prometheusExporter
<sup><sup>[↩ Parent](#ironicspec)</sup></sup>



PrometheusExporter configures sensor data collection and Prometheus metrics export.
When enabled, this configures Ironic to collect sensor data and deploys the
ironic-prometheus-exporter container.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>enabled</b></td>
        <td>boolean</td>
        <td>
          Enabled controls whether sensor data collection and metrics export is active.
When true, configures Ironic to collect sensor data and deploys the
ironic-prometheus-exporter container.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>bindAddress</b></td>
        <td>string</td>
        <td>
          BindAddress is the IP address the metrics endpoint listens on.
Defaults to "0.0.0.0" to listen on all interfaces.

Can be set to a specific IP address (e.g. the provisioning network IP)
to limit exposure to a particular network interface, or to "127.0.0.1"
to restrict access to the local host only (note: this makes
ServiceMonitor-based scraping impossible).<br/>
          <br/>
            <i>Default</i>: 0.0.0.0<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>disableServiceMonitor</b></td>
        <td>boolean</td>
        <td>
          DisableServiceMonitor controls whether a ServiceMonitor resource is created.
Set to true if your cluster does not have prometheus-operator installed,
or when you want to run the exporter but manage Prometheus configuration manually.

Must be set to true for a highly available deployment. In this case, every replica
provides different metrics, which must be aggregated on the consumer side.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sensorCollectionInterval</b></td>
        <td>integer</td>
        <td>
          SensorCollectionInterval defines how often (in seconds) sensor data
is collected from BMCs using Ironic. Must be at least 60 seconds.<br/>
          <br/>
            <i>Default</i>: 60<br/>
            <i>Minimum</i>: 60<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.tls
<sup><sup>[↩ Parent](#ironicspec)</sup></sup>



TLS defines TLS-related settings for various network interactions.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicspectlsbmcca">bmcCA</a></b></td>
        <td>object</td>
        <td>
          BMCCA is a reference to a ConfigMap or Secret containing the CA certificate(s)
to use when validating TLS connections to BMCs.
Supported in Ironic 32.0 or newer.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>bmcCAName</b></td>
        <td>string</td>
        <td>
          BMCCAName is a reference to the secret with the CA certificate(s)
to use when validating TLS connections to BMC's.
Supported in Ironic 32.0 or newer.

Deprecated: Use BMCCA instead. This field will be removed in a future release.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>certificateName</b></td>
        <td>string</td>
        <td>
          CertificateName is a reference to the secret with the TLS certificate.
Must contains both the certificate and the private key parts.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>disableVirtualMediaTLS</b></td>
        <td>boolean</td>
        <td>
          DisableVirtualMediaTLS turns off TLS on the virtual media server,
which may be required for hardware that cannot accept HTTPS links.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>insecureRPC</b></td>
        <td>boolean</td>
        <td>
          InsecureRPC disables TLS validation for the internal RPC.
Without it, the certificate must be valid for all IP addresses on
which Ironic replicas may end up running.
Has no effect when HighAvailability is false and requires the
HighAvailability feature gate to be set.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspectlstrustedca">trustedCA</a></b></td>
        <td>object</td>
        <td>
          TrustedCA is a reference to a ConfigMap or Secret containing the CA certificate(s)
to use when validating TLS connections to image servers and other services.
The resource should contain one or more CA certificates in PEM format.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>trustedCAName</b></td>
        <td>string</td>
        <td>
          TrustedCAName is a reference to the configmap with the CA certificate(s)
to use when validating TLS connections to image servers and other services.
The configmap should contain one or more CA certificates in PEM format.
If the configmap contains multiple keys, only the first key will be used and
a warning will be logged.

Deprecated: Use TrustedCA instead. This field will be removed in a future release.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.spec.tls.bmcCA
<sup><sup>[↩ Parent](#ironicspectls)</sup></sup>



BMCCA is a reference to a ConfigMap or Secret containing the CA certificate(s)
to use when validating TLS connections to BMCs.
Supported in Ironic 32.0 or newer.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the resource (ConfigMap or Secret).<br/>
          <br/>
            <i>Enum</i>: ConfigMap, Secret<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource.<br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### Ironic.spec.tls.trustedCA
<sup><sup>[↩ Parent](#ironicspectls)</sup></sup>



TrustedCA is a reference to a ConfigMap or Secret containing the CA certificate(s)
to use when validating TLS connections to image servers and other services.
The resource should contain one or more CA certificates in PEM format.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>kind</b></td>
        <td>enum</td>
        <td>
          Kind of the resource (ConfigMap or Secret).<br/>
          <br/>
            <i>Enum</i>: ConfigMap, Secret<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the resource.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>key</b></td>
        <td>string</td>
        <td>
          Key within the resource to use. If not specified and the resource contains multiple keys,
the first (alphabetically) key will be used and a warning will be logged for other keys.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.status
<sup><sup>[↩ Parent](#ironic)</sup></sup>



IronicStatus defines the observed state of Ironic.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#ironicstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions describe the state of the Ironic deployment.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>installedVersion</b></td>
        <td>string</td>
        <td>
          InstalledVersion identifies which version of Ironic was installed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>requestedVersion</b></td>
        <td>string</td>
        <td>
          RequestedVersion identifies which version of Ironic was last requested.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### Ironic.status.conditions[index]
<sup><sup>[↩ Parent](#ironicstatus)</sup></sup>



Condition contains details for one aspect of the current state of this API Resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          lastTransitionTime is the last time the condition transitioned from one status to another.
This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          message is a human readable message indicating details about the transition.
This may be an empty string.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          reason contains a programmatic identifier indicating the reason for the condition's last transition.
Producers of specific condition types may define expected values and meanings for this field,
and whether the values are considered a guaranteed API.
The value should be a CamelCase string.
This field may not be empty.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>enum</td>
        <td>
          status of the condition, one of True, False, Unknown.<br/>
          <br/>
            <i>Enum</i>: True, False, Unknown<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          type of condition in CamelCase or in foo.example.com/CamelCase.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>observedGeneration</b></td>
        <td>integer</td>
        <td>
          observedGeneration represents the .metadata.generation that the condition was set based upon.
For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
with respect to the current state of the instance.<br/>
          <br/>
            <i>Format</i>: int64<br/>
            <i>Minimum</i>: 0<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>