package ironic

import (
	"errors"
	"fmt"
	"net/netip"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

const (
	ironicPortName    = "ironic-api"
	imagesPortName    = "image-svc"
	imagesTLSPortName = "image-svc-tls"

	ironicUser  = 997
	ironicGroup = 994

	authDir   = "/auth"
	certsDir  = "/certs"
	sharedDir = "/shared"
)

func buildCommonEnvVars(ironic *metal3api.Ironic) []corev1.EnvVar {
	result := []corev1.EnvVar{
		{
			Name:  "RESTART_CONTAINER_CERTIFICATE_UPDATED",
			Value: "true",
		},
		{
			Name:  "IRONIC_LISTEN_PORT",
			Value: strconv.Itoa(int(ironic.Spec.Networking.APIPort)),
		},
		{
			Name:  "HTTP_PORT",
			Value: strconv.Itoa(int(ironic.Spec.Networking.ImageServerPort)),
		},
		{
			Name:  "LISTEN_ALL_INTERFACES",
			Value: strconv.FormatBool(!ironic.Spec.Networking.BindInterface),
		},
		{
			Name:  "USE_IRONIC_INSPECTOR",
			Value: "false",
		},
	}

	networkingProvided := false
	if ironic.Spec.Networking.IPAddress != "" {
		result = append(result,
			corev1.EnvVar{
				Name:  "PROVISIONING_IP",
				Value: ironic.Spec.Networking.IPAddress,
			},
		)
		networkingProvided = true
	}
	if ironic.Spec.Networking.Interface != "" {
		result = append(result,
			corev1.EnvVar{
				Name:  "PROVISIONING_INTERFACE",
				Value: ironic.Spec.Networking.Interface,
			},
		)
		networkingProvided = true
	}
	if len(ironic.Spec.Networking.MACAddresses) > 0 {
		result = append(result,
			corev1.EnvVar{
				Name:  "PROVISIONING_MACS",
				Value: strings.Join(ironic.Spec.Networking.MACAddresses, ","),
			},
		)
		networkingProvided = true
	}
	if !networkingProvided {
		result = append(result,
			corev1.EnvVar{
				Name: "PROVISIONING_IP",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "status.hostIP",
					},
				},
			},
		)
	}

	if ironic.Spec.TLSRef.Name != "" {
		// Ironic will listen on a Unix socket, httpd will be responsible for serving HTTPS.
		result = append(result, []corev1.EnvVar{
			{
				Name:  "IRONIC_PRIVATE_PORT",
				Value: "unix",
			},
			{
				Name:  "IRONIC_REVERSE_PROXY_SETUP",
				Value: "true",
			},
		}...)

		if !ironic.Spec.DisableVirtualMediaTLS {
			result = append(result,
				corev1.EnvVar{
					Name:  "VMEDIA_TLS_PORT",
					Value: strconv.Itoa(int(ironic.Spec.Networking.ImageServerTLSPort)),
				},
			)
		}
	}

	result = appendStringEnv(result,
		"IRONIC_KERNEL_PARAMS", strings.Trim(ironic.Spec.RamdiskExtraKernelParams, " \t\n\r"))

	result = appendStringEnv(result,
		"IRONIC_RAMDISK_SSH_KEY", strings.Trim(ironic.Spec.RamdiskSSHKey, " \t\n\r"))

	result = appendListOfStringsEnv(result,
		"IRONIC_IPA_COLLECTORS", ironic.Spec.Inspection.Collectors, ",")

	result = appendListOfStringsEnv(result,
		"IRONIC_INSPECTOR_VLAN_INTERFACES", ironic.Spec.Inspection.VLANInterfaces, ",")

	return result
}

func buildIronicEnvVars(ironic *metal3api.Ironic, db *metal3api.IronicDatabase, htpasswd string, domain string) []corev1.EnvVar {
	result := buildCommonEnvVars(ironic)
	result = append(result, []corev1.EnvVar{
		{
			Name:  "IRONIC_USE_MARIADB",
			Value: strconv.FormatBool(db != nil),
		},
		{
			Name:  "IRONIC_EXPOSE_JSON_RPC",
			Value: strconv.FormatBool(ironic.Spec.Distributed),
		},
	}...)

	if db != nil {
		result = append(result, commonDatabaseVars(db)...)
		result = append(result,
			corev1.EnvVar{
				Name:  "MARIADB_HOST",
				Value: DatabaseDNSName(db, domain),
			},
		)
	}

	if ironic.Spec.Distributed {
		result = append(result, []corev1.EnvVar{
			// NOTE(dtantsur): this is not strictly correct but is required for JSON RPC authentication
			{
				Name:  "IRONIC_DEPLOYMENT",
				Value: "Conductor",
			},
			{
				Name:  "IRONIC_INSECURE",
				Value: strconv.FormatBool(ironic.Spec.DisableRPCHostValidation),
			},
		}...)
	}

	// When TLS is used, httpd is responsible for authentication.
	// When JSON RPC is enabled, the password is required for it as well.
	if htpasswd != "" && (ironic.Spec.TLSRef.Name == "" || ironic.Spec.Distributed) {
		result = append(result,
			corev1.EnvVar{
				Name: "IRONIC_HTPASSWD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: htpasswd,
						},
						Key: htpasswdKey,
					},
				},
			},
		)
	}

	result = appendStringEnv(result, "IRONIC_EXTERNAL_IP", ironic.Spec.Networking.ExternalIP)

	return result
}

func buildHttpdEnvVars(ironic *metal3api.Ironic, htpasswd string) []corev1.EnvVar {
	result := buildCommonEnvVars(ironic)

	// When TLS is used, httpd is responsible for authentication
	if htpasswd != "" && ironic.Spec.TLSRef.Name != "" {
		result = append(result,
			corev1.EnvVar{
				Name: "IRONIC_HTPASSWD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: htpasswd,
						},
						Key: htpasswdKey,
					},
				},
			},
		)
	}

	return result
}

func buildIronicVolumesAndMounts(ironic *metal3api.Ironic, db *metal3api.IronicDatabase) (volumes []corev1.Volume, mounts []corev1.VolumeMount) {
	volumes = []corev1.Volume{
		{
			Name: "ironic-shared",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
		{
			Name: "ironic-auth",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: ironic.Spec.CredentialsRef.Name,
				},
			},
		},
	}
	mounts = []corev1.VolumeMount{
		{
			Name:      "ironic-shared",
			MountPath: sharedDir,
		},
	}
	if ironic.Spec.Distributed {
		mounts = append(mounts, corev1.VolumeMount{
			Name:      "ironic-auth",
			MountPath: authDir + "/ironic-rpc",
		})
	}

	if ironic.Spec.TLSRef.Name != "" {
		volumes = append(volumes,
			corev1.Volume{
				Name: "cert-ironic",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: ironic.Spec.TLSRef.Name,
					},
				},
			},
		)
		mounts = append(mounts,
			corev1.VolumeMount{
				Name:      "cert-ironic",
				MountPath: certsDir + "/ironic",
				ReadOnly:  true,
			},
		)
		if !ironic.Spec.DisableVirtualMediaTLS {
			mounts = append(mounts,
				corev1.VolumeMount{
					Name:      "cert-ironic",
					MountPath: certsDir + "/vmedia",
					ReadOnly:  true,
				},
			)
		}
	}

	if db != nil && db.Spec.TLSRef.Name != "" {
		volumes = append(volumes, corev1.Volume{
			Name: "cert-mariadb",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: db.Spec.TLSRef.Name,
				},
			},
		})
		mounts = append(mounts, corev1.VolumeMount{
			Name:      "cert-mariadb",
			MountPath: certsDir + "/ca/mariadb",
		})
	}

	return
}

func buildIronicHttpdPorts(ironic *metal3api.Ironic) (ironicPorts []corev1.ContainerPort, httpdPorts []corev1.ContainerPort) {
	httpdPorts = []corev1.ContainerPort{
		{
			Name:          imagesPortName,
			ContainerPort: ironic.Spec.Networking.ImageServerPort,
		},
	}

	apiPort := corev1.ContainerPort{
		Name:          ironicPortName,
		ContainerPort: ironic.Spec.Networking.APIPort,
	}

	if ironic.Spec.TLSRef.Name == "" {
		ironicPorts = append(ironicPorts, apiPort)
	} else {
		httpdPorts = append(httpdPorts, apiPort)
		if !ironic.Spec.DisableVirtualMediaTLS {
			httpdPorts = append(httpdPorts, corev1.ContainerPort{
				Name:          imagesTLSPortName,
				ContainerPort: ironic.Spec.Networking.ImageServerTLSPort,
			})
		}
	}

	return
}

func buildDHCPRange(dhcp *metal3api.DHCP) string {
	prefix, err := netip.ParsePrefix(dhcp.NetworkCIDR)
	if err != nil {
		return "" // don't disable your webhooks people
	}

	return fmt.Sprintf("%s,%s,%d", dhcp.RangeBegin, dhcp.RangeEnd, prefix.Bits())
}

func buildDNSIP(dhcp *metal3api.DHCP) string {
	if dhcp.ServeDNS {
		return "provisioning" // magical value for serving DNS from the provisioning host using its settings
	}

	if dhcp.DNSAddress != "" {
		return dhcp.DNSAddress
	}

	return ""
}

func newURLProbeHandler(ironic *metal3api.Ironic, https bool, port int, path string) corev1.ProbeHandler {
	proto := "http"
	if https {
		proto = "https"
	}

	// NOTE(dtantsur): we could use HTTP GET probe but we cannot pass the certificate there.
	url := fmt.Sprintf("%s://127.0.0.1:%d%s", proto, port, path)
	return corev1.ProbeHandler{
		Exec: &corev1.ExecAction{
			Command: []string{"curl", "-sSfk", url},
		},
	}
}

func newDnsmasqContainer(ironic *metal3api.Ironic) corev1.Container {
	dhcp := ironic.Spec.Networking.DHCP

	envVars := buildCommonEnvVars(ironic)
	envVars = append(envVars, corev1.EnvVar{
		Name:  "DHCP_RANGE",
		Value: buildDHCPRange(dhcp),
	})

	envVars = appendStringEnv(envVars,
		"DNS_IP", buildDNSIP(dhcp))
	envVars = appendStringEnv(envVars,
		"GATEWAY_IP", dhcp.GatewayAddress)
	envVars = appendListOfStringsEnv(envVars,
		"DHCP_HOSTS", dhcp.Hosts, ";")
	envVars = appendListOfStringsEnv(envVars,
		"DHCP_IGNORE", dhcp.Ignore, ",")

	probe := newProbe(corev1.ProbeHandler{
		Exec: &corev1.ExecAction{
			Command: []string{"sh", "-c", "ss -lun | grep :67 && ss -lun | grep :69"},
		},
	})

	return corev1.Container{
		Name:    "dnsmasq",
		Image:   ironic.Spec.Images.Ironic,
		Command: []string{"/bin/rundnsmasq"},
		Env:     envVars,
		SecurityContext: &corev1.SecurityContext{
			RunAsUser:                pointer.Int64(ironicUser),
			RunAsGroup:               pointer.Int64(ironicGroup),
			AllowPrivilegeEscalation: pointer.Bool(true),
			Capabilities: &corev1.Capabilities{
				Drop: []corev1.Capability{"ALL"},
				Add:  []corev1.Capability{"NET_ADMIN", "NET_BIND_SERVICE", "NET_RAW"},
			},
		},
		LivenessProbe:  probe,
		ReadinessProbe: probe,
	}
}

func newIronicPodTemplate(ironic *metal3api.Ironic, db *metal3api.IronicDatabase, apiSecret *corev1.Secret, domain string) (corev1.PodTemplateSpec, error) {
	var htpasswd string
	if apiSecret != nil {
		if len(apiSecret.Data[htpasswdKey]) == 0 {
			return corev1.PodTemplateSpec{}, errors.New("no htpasswd in the API secret")
		}

		htpasswd = apiSecret.Name
	}

	var ipaDownloaderVars []corev1.EnvVar
	ipaDownloaderVars = appendStringEnv(ipaDownloaderVars,
		"IPA_BASEURI", ironic.Spec.Images.AgentDownloadURL)
	ipaDownloaderVars = appendStringEnv(ipaDownloaderVars,
		"IPA_BRANCH", ironic.Spec.Images.AgentBranch)

	volumes, mounts := buildIronicVolumesAndMounts(ironic, db)
	sharedVolumeMount := mounts[0]
	initContainers := []corev1.Container{
		{
			Name:         "ipa-downloader",
			Image:        ironic.Spec.Images.RamdiskDownloader,
			Env:          ipaDownloaderVars,
			VolumeMounts: []corev1.VolumeMount{sharedVolumeMount},
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:  pointer.Int64(ironicUser),
				RunAsGroup: pointer.Int64(ironicGroup),
				Capabilities: &corev1.Capabilities{
					Drop: []corev1.Capability{"ALL"},
				},
			},
		},
	}

	ironicPorts, httpdPorts := buildIronicHttpdPorts(ironic)

	ironicHandler := newURLProbeHandler(ironic, ironic.Spec.TLSRef.Name != "", int(ironic.Spec.Networking.APIPort), "/v1")
	httpdHandler := newURLProbeHandler(ironic, false, int(ironic.Spec.Networking.ImageServerPort), "/images")

	containers := []corev1.Container{
		{
			Name:         "ironic",
			Image:        ironic.Spec.Images.Ironic,
			Command:      []string{"/bin/runironic"},
			Env:          buildIronicEnvVars(ironic, db, htpasswd, domain),
			VolumeMounts: mounts,
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:  pointer.Int64(ironicUser),
				RunAsGroup: pointer.Int64(ironicGroup),
				Capabilities: &corev1.Capabilities{
					Drop: []corev1.Capability{"ALL"},
				},
			},
			Ports:          ironicPorts,
			LivenessProbe:  newProbe(ironicHandler),
			ReadinessProbe: newProbe(ironicHandler),
		},
		{
			Name:         "httpd",
			Image:        ironic.Spec.Images.Ironic,
			Command:      []string{"/bin/runhttpd"},
			Env:          buildHttpdEnvVars(ironic, htpasswd),
			VolumeMounts: mounts,
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:  pointer.Int64(ironicUser),
				RunAsGroup: pointer.Int64(ironicGroup),
				Capabilities: &corev1.Capabilities{
					Drop: []corev1.Capability{"ALL"},
				},
			},
			Ports:          httpdPorts,
			LivenessProbe:  newProbe(httpdHandler),
			ReadinessProbe: newProbe(httpdHandler),
		},
		{
			Name:         "ramdisk-logs",
			Image:        ironic.Spec.Images.Ironic,
			Command:      []string{"/bin/runlogwatch.sh"},
			VolumeMounts: []corev1.VolumeMount{sharedVolumeMount},
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:  pointer.Int64(ironicUser),
				RunAsGroup: pointer.Int64(ironicGroup),
				Capabilities: &corev1.Capabilities{
					Drop: []corev1.Capability{"ALL"},
				},
			},
		},
	}
	if ironic.Spec.Networking.DHCP != nil && !ironic.Spec.Distributed {
		metal3api.SetDHCPDefaults(ironic.Spec.Networking.DHCP)
		err := metal3api.ValidateDHCP(&ironic.Spec, ironic.Spec.Networking.DHCP)
		if err != nil {
			return corev1.PodTemplateSpec{}, err
		}
		containers = append(containers, newDnsmasqContainer(ironic))
	}

	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{metal3api.IronicOperatorLabel: ironicDeploymentName(ironic)},
		},
		Spec: corev1.PodSpec{
			Containers:     containers,
			InitContainers: initContainers,
			Volumes:        volumes,
			// Ironic needs to be accessed by external machines
			HostNetwork: true,
			DNSPolicy:   corev1.DNSClusterFirstWithHostNet,
		},
	}, nil
}
