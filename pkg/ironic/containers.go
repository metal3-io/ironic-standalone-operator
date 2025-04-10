package ironic

import (
	"errors"
	"fmt"
	"net/netip"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

const (
	ironicPortName    = "ironic-api"
	imagesPortName    = "image-svc"
	imagesTLSPortName = "image-svc-tls"

	ironicUser      int64 = 997
	ironicGroup     int64 = 994
	keepalivedUser  int64 = 65532
	keepalivedGroup int64 = 65532

	authDir   = "/auth"
	certsDir  = "/certs"
	sharedDir = "/shared"

	knownExistingPath = "/images/ironic-python-agent.kernel"
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
						APIVersion: "v1",
						FieldPath:  "status.hostIP",
					},
				},
			},
		)
	}

	if ironic.Spec.TLS.CertificateName != "" {
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

		if !ironic.Spec.TLS.DisableVirtualMediaTLS {
			result = append(result,
				corev1.EnvVar{
					Name:  "VMEDIA_TLS_PORT",
					Value: strconv.Itoa(int(ironic.Spec.Networking.ImageServerTLSPort)),
				},
			)
		}
	}

	result = appendStringEnv(result,
		"IRONIC_KERNEL_PARAMS", strings.Trim(ironic.Spec.DeployRamdisk.ExtraKernelParams, " \t\n\r"))

	result = appendStringEnv(result,
		"IRONIC_RAMDISK_SSH_KEY", strings.Trim(ironic.Spec.DeployRamdisk.SSHKey, " \t\n\r"))

	result = appendListOfStringsEnv(result,
		"IRONIC_IPA_COLLECTORS", ironic.Spec.Inspection.Collectors, ",")

	result = appendListOfStringsEnv(result,
		"IRONIC_INSPECTOR_VLAN_INTERFACES", ironic.Spec.Inspection.VLANInterfaces, ",")

	return result
}

func buildExtraConfigVars(ironic *metal3api.Ironic) []corev1.EnvVar {
	var result []corev1.EnvVar

	for _, extraConfig := range ironic.Spec.ExtraConfig {
		// Default group value
		group := "DEFAULT"
		if extraConfig.Group != "" {
			group = extraConfig.Group
		}

		if extraConfig.Name != "" && extraConfig.Value != "" {
			name := extraConfig.Name
			value := extraConfig.Value

			result = append(result,
				corev1.EnvVar{
					Name:  fmt.Sprintf("OS_%s__%s", strings.ToUpper(group), strings.ToUpper(name)),
					Value: value,
				})
		}
	}

	return result
}

func databaseClientEnvVars(cctx ControllerContext, db *metal3api.Database) []corev1.EnvVar {
	envVars := []corev1.EnvVar{
		{
			Name:  "MARIADB_HOST",
			Value: db.Host,
		},
		{
			Name:  "MARIADB_DATABASE",
			Value: db.Name,
		},
	}

	// NOTE(dtantsur): remove when versions older than 29.0 are no longer supported
	if cctx.VersionInfo.InstalledVersion.Compare(versionMountDatabaseSecret) < 0 {
		envVars = append(envVars, []corev1.EnvVar{
			{
				Name: "MARIADB_USER",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: db.CredentialsName,
						},
						Key: "username",
					},
				},
			},
			{
				Name: "MARIADB_PASSWORD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: db.CredentialsName,
						},
						Key: "password",
					},
				},
			}}...)
	}

	return envVars
}

func buildIronicEnvVars(cctx ControllerContext, ironic *metal3api.Ironic, db *metal3api.Database, htpasswd string) []corev1.EnvVar {
	result := buildCommonEnvVars(ironic)
	result = append(result, []corev1.EnvVar{
		{
			Name:  "IRONIC_USE_MARIADB",
			Value: strconv.FormatBool(db != nil),
		},
		{
			Name:  "IRONIC_EXPOSE_JSON_RPC",
			Value: strconv.FormatBool(ironic.Spec.HighAvailability),
		},
	}...)

	if db != nil {
		result = append(result, databaseClientEnvVars(cctx, db)...)
		// NOTE(dtantsur): upgrades are handled by a separate job
		result = append(result, corev1.EnvVar{
			Name:  "IRONIC_SKIP_DBSYNC",
			Value: "true",
		})
	}

	if ironic.Spec.HighAvailability {
		result = append(result, []corev1.EnvVar{
			// NOTE(dtantsur): this is not strictly correct but is required for JSON RPC authentication
			{
				Name:  "IRONIC_DEPLOYMENT",
				Value: "Conductor",
			},
		}...)
	}

	// When TLS is used, httpd is responsible for authentication.
	// When JSON RPC is enabled, the password is required for it as well.
	if htpasswd != "" && (ironic.Spec.TLS.CertificateName == "" || ironic.Spec.HighAvailability) {
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

	if ironic.Spec.ExtraConfig != nil {
		result = append(result, buildExtraConfigVars(ironic)...)
	}

	result = appendStringEnv(result, "IRONIC_EXTERNAL_IP", ironic.Spec.Networking.ExternalIP)

	return result
}

func buildHttpdEnvVars(ironic *metal3api.Ironic, htpasswd string) []corev1.EnvVar {
	result := buildCommonEnvVars(ironic)

	// When TLS is used, httpd is responsible for authentication
	if htpasswd != "" && ironic.Spec.TLS.CertificateName != "" {
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

func databaseClientMounts(cctx ControllerContext, db *metal3api.Database) (volumes []corev1.Volume, mounts []corev1.VolumeMount) {
	if cctx.VersionInfo.InstalledVersion.Compare(versionMountDatabaseSecret) >= 0 {
		volumes = append(volumes, corev1.Volume{
			Name: "auth-mariadb",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  db.CredentialsName,
					DefaultMode: ptr.To(corev1.SecretVolumeSourceDefaultMode),
				},
			},
		})
		mounts = append(mounts, corev1.VolumeMount{
			Name:      "auth-mariadb",
			MountPath: authDir + "/mariadb",
		})
	}

	if db.TLSCertificateName != "" {
		volumes = append(volumes, corev1.Volume{
			Name: "cert-mariadb",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  db.TLSCertificateName,
					DefaultMode: ptr.To(corev1.SecretVolumeSourceDefaultMode),
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

func buildIronicVolumesAndMounts(cctx ControllerContext, ironic *metal3api.Ironic, db *metal3api.Database) (volumes []corev1.Volume, mounts []corev1.VolumeMount) {
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
					SecretName:  ironic.Spec.APICredentialsName,
					DefaultMode: ptr.To(corev1.SecretVolumeSourceDefaultMode),
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
	if ironic.Spec.HighAvailability {
		mounts = append(mounts, corev1.VolumeMount{
			Name:      "ironic-auth",
			MountPath: authDir + "/ironic-rpc",
		})
	}

	if ironic.Spec.TLS.CertificateName != "" {
		volumes = append(volumes,
			corev1.Volume{
				Name: "cert-ironic",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName:  ironic.Spec.TLS.CertificateName,
						DefaultMode: ptr.To(corev1.SecretVolumeSourceDefaultMode),
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
		if !ironic.Spec.TLS.DisableVirtualMediaTLS {
			mounts = append(mounts,
				corev1.VolumeMount{
					Name:      "cert-ironic",
					MountPath: certsDir + "/vmedia",
					ReadOnly:  true,
				},
			)
		}
	}

	if db != nil {
		dbVolumes, dbMounts := databaseClientMounts(cctx, db)
		volumes = append(volumes, dbVolumes...)
		mounts = append(mounts, dbMounts...)
	}

	return
}

func buildIronicHttpdPorts(ironic *metal3api.Ironic) (ironicPorts []corev1.ContainerPort, httpdPorts []corev1.ContainerPort) {
	httpdPorts = []corev1.ContainerPort{
		{
			Name:          imagesPortName,
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: ironic.Spec.Networking.ImageServerPort,
			HostPort:      ironic.Spec.Networking.ImageServerPort,
		},
	}

	apiPort := corev1.ContainerPort{
		Name:          ironicPortName,
		Protocol:      corev1.ProtocolTCP,
		ContainerPort: ironic.Spec.Networking.APIPort,
		HostPort:      ironic.Spec.Networking.APIPort,
	}

	if ironic.Spec.TLS.CertificateName == "" {
		ironicPorts = append(ironicPorts, apiPort)
	} else {
		httpdPorts = append(httpdPorts, apiPort)
		if !ironic.Spec.TLS.DisableVirtualMediaTLS {
			httpdPorts = append(httpdPorts, corev1.ContainerPort{
				Name:          imagesTLSPortName,
				Protocol:      corev1.ProtocolTCP,
				ContainerPort: ironic.Spec.Networking.ImageServerTLSPort,
				HostPort:      ironic.Spec.Networking.ImageServerTLSPort,
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

func newURLProbeHandler(ironic *metal3api.Ironic, https bool, port int, path string, requiresOk bool) corev1.ProbeHandler {
	proto := "http"
	if https {
		proto = "https"
	}

	curlCmd := []string{
		"curl",
		"--head",
		"--insecure",
		// NOTE(dtantsur): prevent redirects from being treated as success
		"--location",
		"--silent",
		"--show-error",
		// NOTE(dtantsur): --head outputs headers even with --silent
		"--output", "/dev/null",
	}
	if requiresOk {
		curlCmd = append(curlCmd, "--fail")
	}

	// NOTE(dtantsur): we could use HTTP GET probe but we cannot pass the certificate there.
	url := fmt.Sprintf("%s://127.0.0.1:%d%s", proto, port, path)
	return corev1.ProbeHandler{
		Exec: &corev1.ExecAction{
			Command: append(curlCmd, url),
		},
	}
}

func newDnsmasqContainer(versionInfo VersionInfo, ironic *metal3api.Ironic) corev1.Container {
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
		Image:   versionInfo.IronicImage,
		Command: []string{"/bin/rundnsmasq"},
		Env:     envVars,
		SecurityContext: &corev1.SecurityContext{
			RunAsUser:                ptr.To(ironicUser),
			RunAsGroup:               ptr.To(ironicGroup),
			AllowPrivilegeEscalation: ptr.To(true),
			Capabilities: &corev1.Capabilities{
				Drop: []corev1.Capability{"ALL"},
				Add:  []corev1.Capability{"NET_ADMIN", "NET_BIND_SERVICE", "NET_RAW"},
			},
		},
		LivenessProbe:  probe,
		ReadinessProbe: probe,
	}
}

func newKeepalivedContainer(versionInfo VersionInfo, ironic *metal3api.Ironic) corev1.Container {
	envVars := []corev1.EnvVar{
		{
			Name:  "PROVISIONING_IP",
			Value: ironic.Spec.Networking.IPAddress,
		},
		{
			Name:  "PROVISIONING_INTERFACE",
			Value: ironic.Spec.Networking.Interface,
		},
	}

	return corev1.Container{
		Name:  "keepalived",
		Image: versionInfo.KeepalivedImage,
		Env:   envVars,
		SecurityContext: &corev1.SecurityContext{
			RunAsUser:                ptr.To(keepalivedUser),
			RunAsGroup:               ptr.To(keepalivedGroup),
			Privileged:               ptr.To(false),
			AllowPrivilegeEscalation: ptr.To(true),
			Capabilities: &corev1.Capabilities{
				Drop: []corev1.Capability{"ALL"},
				Add:  []corev1.Capability{"NET_ADMIN", "NET_BROADCAST", "NET_RAW"},
			},
		},
	}
}

func newIronicPodTemplate(cctx ControllerContext, ironic *metal3api.Ironic, db *metal3api.Database, apiSecret *corev1.Secret, domain string) (corev1.PodTemplateSpec, error) {
	var htpasswd string
	if apiSecret != nil {
		if len(apiSecret.Data[htpasswdKey]) == 0 {
			return corev1.PodTemplateSpec{}, errors.New("no htpasswd in the API secret")
		}

		htpasswd = apiSecret.Name
	}

	var ipaDownloaderVars []corev1.EnvVar
	ipaDownloaderVars = appendStringEnv(ipaDownloaderVars,
		"IPA_BRANCH", cctx.VersionInfo.AgentBranch)

	volumes, mounts := buildIronicVolumesAndMounts(cctx, ironic, db)
	sharedVolumeMount := mounts[0]

	var initContainers []corev1.Container
	if !ironic.Spec.DeployRamdisk.DisableDownloader {
		initContainers = append(initContainers, corev1.Container{
			Name:         "ramdisk-downloader",
			Image:        cctx.VersionInfo.RamdiskDownloaderImage,
			Env:          ipaDownloaderVars,
			VolumeMounts: []corev1.VolumeMount{sharedVolumeMount},
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:  ptr.To(ironicUser),
				RunAsGroup: ptr.To(ironicGroup),
				Capabilities: &corev1.Capabilities{
					Drop: []corev1.Capability{"ALL"},
				},
			},
		})
	}

	ironicPorts, httpdPorts := buildIronicHttpdPorts(ironic)

	ironicHandler := newURLProbeHandler(ironic, ironic.Spec.TLS.CertificateName != "", int(ironic.Spec.Networking.APIPort), "/v1", true)
	httpPathExpected := !ironic.Spec.DeployRamdisk.DisableDownloader
	httpdHandler := newURLProbeHandler(ironic, false, int(ironic.Spec.Networking.ImageServerPort), knownExistingPath, httpPathExpected)

	containers := []corev1.Container{
		{
			Name:         "ironic",
			Image:        cctx.VersionInfo.IronicImage,
			Command:      []string{"/bin/runironic"},
			Env:          buildIronicEnvVars(cctx, ironic, db, htpasswd),
			VolumeMounts: mounts,
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:  ptr.To(ironicUser),
				RunAsGroup: ptr.To(ironicGroup),
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
			Image:        cctx.VersionInfo.IronicImage,
			Command:      []string{"/bin/runhttpd"},
			Env:          buildHttpdEnvVars(ironic, htpasswd),
			VolumeMounts: mounts,
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:  ptr.To(ironicUser),
				RunAsGroup: ptr.To(ironicGroup),
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
			Image:        cctx.VersionInfo.IronicImage,
			Command:      []string{"/bin/runlogwatch.sh"},
			VolumeMounts: []corev1.VolumeMount{sharedVolumeMount},
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:  ptr.To(ironicUser),
				RunAsGroup: ptr.To(ironicGroup),
				Capabilities: &corev1.Capabilities{
					Drop: []corev1.Capability{"ALL"},
				},
			},
		},
	}
	if ironic.Spec.Networking.DHCP != nil && !ironic.Spec.HighAvailability {
		err := ValidateDHCP(&ironic.Spec, ironic.Spec.Networking.DHCP)
		if err != nil {
			return corev1.PodTemplateSpec{}, err
		}
		containers = append(containers, newDnsmasqContainer(cctx.VersionInfo, ironic))
	}

	if ironic.Spec.Networking.IPAddressManager == metal3api.IPAddressManagerKeepalived {
		containers = append(containers, newKeepalivedContainer(cctx.VersionInfo, ironic))
	}

	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				metal3api.IronicAppLabel:     ironicDeploymentName(ironic),
				metal3api.IronicServiceLabel: ironic.Name,
				metal3api.IronicVersionLabel: cctx.VersionInfo.InstalledVersion.String(),
			},
			Annotations: secretVersionAnnotations("api-secret", apiSecret),
		},
		Spec: corev1.PodSpec{
			Containers:     containers,
			InitContainers: initContainers,
			Volumes:        volumes,
			// Ironic needs to be accessed by external machines
			HostNetwork:  true,
			DNSPolicy:    corev1.DNSClusterFirstWithHostNet,
			NodeSelector: ironic.Spec.NodeSelector,
		},
	}, nil
}
