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
          BindInterface makes Ironic API bound to only one interface.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#ironicspecnetworkingdhcp">dhcp</a></b></td>
        <td>object</td>
        <td>
          DHCP is a configuration of DHCP for the network boot service (dnsmasq).
The service is only deployed when this is set.
This setting is currently incompatible with the highly available architecture.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>externalIP</b></td>
        <td>string</td>
        <td>
          ExternalIP is used for accessing API and the image server from remote hosts.
This settings only applies to virtual media deployments. The IP will not be accessed from the cluster itself.<br/>
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
        <td><b>interface</b></td>
        <td>string</td>
        <td>
          Interface is a Linux network device to listen on.
Detected from IPAddress if missing.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ipAddress</b></td>
        <td>string</td>
        <td>
          IPAddress is the main IP address to listen on and use for communication.
Detected from Interface if missing. Cannot be provided for a highly available architecture.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ipAddressManager</b></td>
        <td>enum</td>
        <td>
          Configures the way the provided IP address will be managed on the provided interface.
By default, the IP address is expected to be already present.
Use "keepalived" to start a Keepalived container managing the IP address.
Warning: keepalived is not compatible with the highly available architecture.<br/>
          <br/>
            <i>Enum</i>: , keepalived<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>macAddresses</b></td>
        <td>[]string</td>
        <td>
          MACAddresses can be provided to make the start script pick the interface matching any of these addresses.
Only set if no other options can be used.<br/>
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
          Resources resize policy for the container.<br/>
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
This requires the ProcMountType feature flag to be enabled.
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
          Resources resize policy for the container.<br/>
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
This requires the ProcMountType feature flag to be enabled.
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
the first key will be used and a warning will be logged for other keys.<br/>
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