package ironic

import (
	"strings"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

const ironicContainerName = "ironic"

func TestExpectedContainers(t *testing.T) {
	testCases := []struct {
		Scenario string

		Ironic metal3api.IronicSpec

		ExpectedContainerNames     []string
		ExpectedInitContainerNames []string
		ExpectedError              string
	}{
		{
			Scenario:                   "empty",
			ExpectedContainerNames:     []string{"httpd", "ironic", "ramdisk-logs"},
			ExpectedInitContainerNames: []string{"ramdisk-downloader"},
		},
		{
			Scenario: "Keepalived",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					Interface:        "eth0",
					IPAddress:        "192.0.2.2",
					IPAddressManager: metal3api.IPAddressManagerKeepalived,
				},
			},
			ExpectedContainerNames:     []string{"httpd", "ironic", "keepalived", "ramdisk-logs"},
			ExpectedInitContainerNames: []string{"ramdisk-downloader"},
		},
		{
			Scenario: "Keepalived and DHCP",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					DHCP: &metal3api.DHCP{
						NetworkCIDR: "192.0.2.1/24",
						RangeBegin:  "192.0.2.10",
						RangeEnd:    "192.0.2.200",
					},
					Interface:        "eth0",
					IPAddress:        "192.0.2.2",
					IPAddressManager: metal3api.IPAddressManagerKeepalived,
				},
			},
			ExpectedContainerNames:     []string{"dnsmasq", "httpd", "ironic", "keepalived", "ramdisk-logs"},
			ExpectedInitContainerNames: []string{"ramdisk-downloader"},
		},
		{
			Scenario: "Keepalived list",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					Interface: "eth0",
					IPAddress: "192.0.2.2",
					Keepalived: []metal3api.KeepalivedIP{
						{IPAddress: "192.0.2.2", Interface: "eth0"},
						{IPAddress: "192.168.1.50", Interface: "eth1"},
					},
				},
			},
			ExpectedContainerNames:     []string{"httpd", "ironic", "keepalived", "ramdisk-logs"},
			ExpectedInitContainerNames: []string{"ramdisk-downloader"},
		},
		{
			Scenario: "No ramdisk downloader",
			Ironic: metal3api.IronicSpec{
				DeployRamdisk: metal3api.DeployRamdisk{
					DisableDownloader: true,
				},
			},
			ExpectedContainerNames:     []string{"httpd", "ironic", "ramdisk-logs"},
			ExpectedInitContainerNames: []string{},
		},
		{
			Scenario: "PrometheusExporter enabled",
			Ironic: metal3api.IronicSpec{
				PrometheusExporter: &metal3api.PrometheusExporter{
					Enabled:                  true,
					SensorCollectionInterval: 60,
				},
			},
			ExpectedContainerNames:     []string{"httpd", "ironic", "ironic-prometheus-exporter", "ramdisk-logs"},
			ExpectedInitContainerNames: []string{"ramdisk-downloader"},
		},
		{
			Scenario: "PrometheusExporter disabled",
			Ironic: metal3api.IronicSpec{
				PrometheusExporter: &metal3api.PrometheusExporter{
					Enabled: false,
				},
			},
			ExpectedContainerNames:     []string{"httpd", "ironic", "ramdisk-logs"},
			ExpectedInitContainerNames: []string{"ramdisk-downloader"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			cctx := ControllerContext{}
			secret := &corev1.Secret{
				Data: map[string][]byte{
					"htpasswd": []byte("abcd"),
				},
			}
			ironic := &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test",
					Name:      "test",
				},
				Spec: tc.Ironic,
			}

			resources := Resources{Ironic: ironic, APISecret: secret}
			podTemplate, err := newIronicPodTemplate(cctx, resources)
			if tc.ExpectedError == "" {
				require.NoError(t, err)

				var containerNames []string
				for _, cont := range podTemplate.Spec.Containers {
					containerNames = append(containerNames, cont.Name)
				}

				assert.ElementsMatch(t, tc.ExpectedContainerNames, containerNames)
			} else {
				assert.ErrorContains(t, err, tc.ExpectedError)
			}
		})
	}
}

func TestImageOverrides(t *testing.T) {
	cctx := ControllerContext{}
	secret := &corev1.Secret{
		Data: map[string][]byte{"htpasswd": []byte("abcd")},
	}
	expectedImages := map[string]string{
		"httpd":              "myorg/myironic:test",
		"ironic":             "myorg/myironic:test",
		"keepalived":         "myorg/mykeepalived:test",
		"ramdisk-downloader": "myorg/mydownloader:test",
		"ramdisk-logs":       "myorg/myironic:test",
	}
	ironic := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "test",
		},
		Spec: metal3api.IronicSpec{
			Images: metal3api.Images{
				DeployRamdiskBranch:     "stable/x.y",
				DeployRamdiskDownloader: "myorg/mydownloader:test",
				Ironic:                  "myorg/myironic:test",
				Keepalived:              "myorg/mykeepalived:test",
			},
			Networking: metal3api.Networking{
				Interface:        "eth0",
				IPAddress:        "192.0.2.2",
				IPAddressManager: metal3api.IPAddressManagerKeepalived,
			},
		},
	}

	version, err := cctx.VersionInfo.WithIronicOverrides(ironic)
	require.NoError(t, err)
	cctx.VersionInfo = version

	resources := Resources{Ironic: ironic, APISecret: secret}
	podTemplate, err := newIronicPodTemplate(cctx, resources)
	require.NoError(t, err)

	images := make(map[string]string, len(expectedImages))
	var actualBranch string
	for _, cont := range podTemplate.Spec.InitContainers {
		images[cont.Name] = cont.Image
		if cont.Name == "ramdisk-downloader" {
			for _, env := range cont.Env {
				if env.Name == "IPA_BRANCH" {
					actualBranch = env.Value
				}
			}
		}
	}
	for _, cont := range podTemplate.Spec.Containers {
		images[cont.Name] = cont.Image
	}

	assert.Equal(t, expectedImages, images)
	assert.Equal(t, "stable/x.y", actualBranch)
}

func TestExpectedExtraEnvVars(t *testing.T) {
	cctx := ControllerContext{}
	secret := &corev1.Secret{
		Data: map[string][]byte{
			"htpasswd": []byte("abcd"),
		},
	}

	expectedExtraVars := map[string]string{
		"OS_PXE__BOOT_RETRY_TIMEOUT":            "1200",
		"OS_CONDUCTOR__DEPLOY_CALLBACK_TIMEOUT": "4800",
		"OS_CONDUCTOR__INSPECT_TIMEOUT":         "1800",
	}

	ironic := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "test",
		},
		Spec: metal3api.IronicSpec{
			Networking: metal3api.Networking{
				Interface:        "eth0",
				IPAddress:        "192.0.2.2",
				IPAddressManager: metal3api.IPAddressManagerKeepalived,
				RPCPort:          6189,
			},
			ExtraConfig: []metal3api.ExtraConfig{
				{
					Group: "pxe",
					Name:  "boot_retry_timeout",
					Value: "1200",
				},
				{
					Group: "conductor",
					Name:  "deploy_callback_timeout",
					Value: "4800",
				},
				{
					Group: "conductor",
					Name:  "inspect_timeout",
					Value: "1800",
				},
			},
		},
	}

	resources := Resources{Ironic: ironic, APISecret: secret}
	podTemplate, err := newIronicPodTemplate(cctx, resources)
	require.NoError(t, err)

	extraVars := make(map[string]string, len(expectedExtraVars))
	for _, env := range podTemplate.Spec.Containers[0].Env {
		if strings.Contains(env.Name, "OS_") {
			extraVars[env.Name] = env.Value
		}
	}

	assert.Equal(t, expectedExtraVars, extraVars)
}

func TestAnnotationLabelOverrides(t *testing.T) {
	cctx := ControllerContext{}
	secret := &corev1.Secret{
		Data: map[string][]byte{"htpasswd": []byte("abcd")},
	}
	ironic := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "test",
		},
		Spec: metal3api.IronicSpec{
			Overrides: &metal3api.Overrides{
				Annotations: map[string]string{
					"annotation.example.com": "my-annotation",
				},
				Labels: map[string]string{
					"label.example.com": "my-label",
				},
			},
		},
	}

	version, err := cctx.VersionInfo.WithIronicOverrides(ironic)
	require.NoError(t, err)
	cctx.VersionInfo = version

	resources := Resources{Ironic: ironic, APISecret: secret}
	podTemplate, err := newIronicPodTemplate(cctx, resources)
	require.NoError(t, err)

	assert.Equal(t, "my-annotation", podTemplate.Annotations["annotation.example.com"])
	assert.NotEmpty(t, podTemplate.Annotations["ironic.metal3.io/api-secret-version"])
	assert.Equal(t, "my-label", podTemplate.Labels["label.example.com"])
	assert.Equal(t, "test", podTemplate.Labels[metal3api.IronicServiceLabel])
}

func TestTrustedCAConfigMap(t *testing.T) {
	testCases := []struct {
		Scenario                string
		TrustedCAConfigMap      *corev1.ConfigMap
		ExpectVolume            bool
		ExpectVolumeMount       bool
		ExpectEnvVar            bool
		ExpectedVolumeMountPath string
		ExpectedEnvVarValue     string
	}{
		{
			Scenario: "with TrustedCAConfigMap single key",
			TrustedCAConfigMap: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "trusted-ca",
					Namespace: "test",
				},
				Data: map[string]string{
					"ca-bundle.crt": "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
				},
			},
			ExpectVolume:            true,
			ExpectVolumeMount:       true,
			ExpectEnvVar:            true,
			ExpectedVolumeMountPath: "/certs/ca/trusted",
			ExpectedEnvVarValue:     "/certs/ca/trusted/ca-bundle.crt",
		},
		{
			Scenario: "with TrustedCAConfigMap multiple keys",
			TrustedCAConfigMap: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "trusted-ca",
					Namespace: "test",
				},
				Data: map[string]string{
					"ca-bundle.crt": "-----BEGIN CERTIFICATE-----\ntest\n-----END CERTIFICATE-----",
					"extra-ca.crt":  "-----BEGIN CERTIFICATE-----\nextra\n-----END CERTIFICATE-----",
				},
			},
			ExpectVolume:            true,
			ExpectVolumeMount:       true,
			ExpectEnvVar:            true,
			ExpectedVolumeMountPath: "/certs/ca/trusted",
			// Note: The actual key used will depend on map iteration order, but we just verify it exists
			ExpectedEnvVarValue: "", // We'll check it contains /certs/ca/trusted/ prefix instead
		},
		{
			Scenario:           "without TrustedCAConfigMap",
			TrustedCAConfigMap: nil,
			ExpectVolume:       false,
			ExpectVolumeMount:  false,
			ExpectEnvVar:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			cctx := ControllerContext{}
			apiSecret := &corev1.Secret{
				Data: map[string][]byte{
					"htpasswd": []byte("abcd"),
				},
			}
			ironic := &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test",
					Name:      "test",
				},
				Spec: metal3api.IronicSpec{},
			}

			resources := Resources{
				Ironic:             ironic,
				APISecret:          apiSecret,
				TrustedCAConfigMap: tc.TrustedCAConfigMap,
			}

			podTemplate, err := newIronicPodTemplate(cctx, resources)
			require.NoError(t, err)

			// Check volume
			var foundVolume bool
			for _, vol := range podTemplate.Spec.Volumes {
				if vol.Name == "trusted-ca" {
					foundVolume = true
					if tc.ExpectVolume {
						assert.NotNil(t, vol.ConfigMap)
						assert.Equal(t, tc.TrustedCAConfigMap.Name, vol.ConfigMap.Name)
					}
					break
				}
			}
			assert.Equal(t, tc.ExpectVolume, foundVolume, "Volume existence mismatch")

			// Check volume mount on ironic container
			var ironicContainer *corev1.Container
			for i := range podTemplate.Spec.Containers {
				if podTemplate.Spec.Containers[i].Name == ironicContainerName {
					ironicContainer = &podTemplate.Spec.Containers[i]
					break
				}
			}
			require.NotNil(t, ironicContainer, "Ironic container should exist")

			var foundMount bool
			for _, mount := range ironicContainer.VolumeMounts {
				if mount.Name == "trusted-ca" {
					foundMount = true
					if tc.ExpectVolumeMount {
						assert.Equal(t, tc.ExpectedVolumeMountPath, mount.MountPath)
						assert.True(t, mount.ReadOnly)
					}
					break
				}
			}
			assert.Equal(t, tc.ExpectVolumeMount, foundMount, "Volume mount existence mismatch")

			// Check environment variable (WEBSERVER_CACERT_FILE)
			var foundWebserverCACert bool
			var webserverCACertValue string
			for _, env := range ironicContainer.Env {
				if env.Name == "WEBSERVER_CACERT_FILE" {
					foundWebserverCACert = true
					webserverCACertValue = env.Value
				}
			}
			assert.Equal(t, tc.ExpectEnvVar, foundWebserverCACert, "WEBSERVER_CACERT_FILE environment variable existence mismatch")

			if tc.ExpectEnvVar {
				if tc.ExpectedEnvVarValue != "" {
					// Exact match for single key case
					assert.Equal(t, tc.ExpectedEnvVarValue, webserverCACertValue, "WEBSERVER_CACERT_FILE value mismatch")
				} else {
					// For multiple keys case, just verify it starts with the correct prefix
					assert.Contains(t, webserverCACertValue, "/certs/ca/trusted/", "WEBSERVER_CACERT_FILE should contain /certs/ca/trusted/")
				}
			}
		})
	}
}

func TestIronicPortEnvVars(t *testing.T) {
	testCases := []struct {
		name               string
		apiPort            int32
		expectedListenPort string
		expectedAccessPort string
	}{
		{
			name:               "standard API port",
			apiPort:            6385,
			expectedListenPort: "6385",
			expectedAccessPort: "6385",
		},
		{
			name:               "custom API port",
			apiPort:            6388,
			expectedListenPort: "6388",
			expectedAccessPort: "6388",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cctx := ControllerContext{}
			secret := &corev1.Secret{
				Data: map[string][]byte{
					"htpasswd": []byte("abcd"),
				},
			}
			ironic := &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test",
					Name:      "test",
				},
				Spec: metal3api.IronicSpec{
					Networking: metal3api.Networking{
						APIPort: tc.apiPort,
					},
				},
			}

			resources := Resources{Ironic: ironic, APISecret: secret}
			podTemplate, err := newIronicPodTemplate(cctx, resources)
			require.NoError(t, err)

			// Find the ironic container
			var ironicContainer *corev1.Container
			for i := range podTemplate.Spec.Containers {
				if podTemplate.Spec.Containers[i].Name == ironicContainerName {
					ironicContainer = &podTemplate.Spec.Containers[i]
					break
				}
			}
			require.NotNil(t, ironicContainer, "ironic container should exist")

			// Check for IRONIC_LISTEN_PORT and IRONIC_ACCESS_PORT env vars
			var foundListenPort, foundAccessPort bool
			var listenPortValue, accessPortValue string
			for _, env := range ironicContainer.Env {
				if env.Name == "IRONIC_LISTEN_PORT" {
					foundListenPort = true
					listenPortValue = env.Value
				}
				if env.Name == "IRONIC_ACCESS_PORT" {
					foundAccessPort = true
					accessPortValue = env.Value
				}
			}

			assert.True(t, foundListenPort, "IRONIC_LISTEN_PORT env var should be present")
			assert.True(t, foundAccessPort, "IRONIC_ACCESS_PORT env var should be present")
			assert.Equal(t, tc.expectedListenPort, listenPortValue, "IRONIC_LISTEN_PORT value mismatch")
			assert.Equal(t, tc.expectedAccessPort, accessPortValue, "IRONIC_ACCESS_PORT value mismatch")
			// Both ports should always be the same
			assert.Equal(t, listenPortValue, accessPortValue, "IRONIC_LISTEN_PORT and IRONIC_ACCESS_PORT should have the same value")
		})
	}
}

func TestPrometheusExporterEnvVars(t *testing.T) {
	testCases := []struct {
		name                   string
		prometheusExporter     *metal3api.PrometheusExporter
		expectedSendSensorData string
		expectedSensorInterval string
		expectedFlaskRunHost   string
	}{
		{
			name: "PrometheusExporter enabled with default interval",
			prometheusExporter: &metal3api.PrometheusExporter{
				Enabled:                  true,
				SensorCollectionInterval: 0, // Should default to 60
			},
			expectedSendSensorData: "true",
			expectedSensorInterval: "60",
			expectedFlaskRunHost:   "0.0.0.0",
		},
		{
			name: "PrometheusExporter enabled with custom interval",
			prometheusExporter: &metal3api.PrometheusExporter{
				Enabled:                  true,
				SensorCollectionInterval: 120,
			},
			expectedSendSensorData: "true",
			expectedSensorInterval: "120",
			expectedFlaskRunHost:   "0.0.0.0",
		},
		{
			name: "PrometheusExporter enabled with bindAddress 0.0.0.0",
			prometheusExporter: &metal3api.PrometheusExporter{
				Enabled:     true,
				BindAddress: "0.0.0.0",
			},
			expectedSendSensorData: "true",
			expectedSensorInterval: "60",
			expectedFlaskRunHost:   "0.0.0.0",
		},
		{
			name: "PrometheusExporter enabled with custom bindAddress",
			prometheusExporter: &metal3api.PrometheusExporter{
				Enabled:     true,
				BindAddress: "192.168.1.10",
			},
			expectedSendSensorData: "true",
			expectedSensorInterval: "60",
			expectedFlaskRunHost:   "192.168.1.10",
		},
		{
			name: "PrometheusExporter enabled with empty bindAddress defaults to wildcard",
			prometheusExporter: &metal3api.PrometheusExporter{
				Enabled: true,
			},
			expectedSendSensorData: "true",
			expectedSensorInterval: "60",
			expectedFlaskRunHost:   "0.0.0.0",
		},
		{
			name: "PrometheusExporter disabled",
			prometheusExporter: &metal3api.PrometheusExporter{
				Enabled: false,
			},
		},
		{
			name:               "PrometheusExporter nil",
			prometheusExporter: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cctx := ControllerContext{}
			secret := &corev1.Secret{
				Data: map[string][]byte{
					"htpasswd": []byte("abcd"),
				},
			}
			ironic := &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test",
					Name:      "test",
				},
				Spec: metal3api.IronicSpec{
					PrometheusExporter: tc.prometheusExporter,
				},
			}

			resources := Resources{Ironic: ironic, APISecret: secret}
			podTemplate, err := newIronicPodTemplate(cctx, resources)
			require.NoError(t, err)

			expectExporter := tc.expectedSendSensorData == "true"

			var ironicContainer *corev1.Container
			var exporterContainer *corev1.Container
			for _, cont := range podTemplate.Spec.Containers {
				contCopy := cont
				if cont.Name == ironicContainerName {
					ironicContainer = &contCopy
				}
				if cont.Name == "ironic-prometheus-exporter" {
					exporterContainer = &contCopy
				}
			}
			require.NotNil(t, ironicContainer, "ironic container not found")
			if expectExporter {
				require.NotNil(t, exporterContainer, "ironic-prometheus-exporter container not found")
				assert.Len(t, exporterContainer.Env, 2)
				assert.Equal(t, "FLASK_RUN_HOST", exporterContainer.Env[0].Name)
				assert.Equal(t, tc.expectedFlaskRunHost, exporterContainer.Env[0].Value)
				assert.Equal(t, "FLASK_RUN_PORT", exporterContainer.Env[1].Name)
				assert.Equal(t, "9608", exporterContainer.Env[1].Value)
				assert.Len(t, exporterContainer.Ports, 1)
				assert.Equal(t, int32(9608), exporterContainer.Ports[0].ContainerPort)
			}

			// Check for SEND_SENSOR_DATA env var
			var foundSendSensorData, foundSensorInterval bool
			for _, env := range ironicContainer.Env {
				if env.Name == "SEND_SENSOR_DATA" {
					foundSendSensorData = true
					assert.Equal(t, tc.expectedSendSensorData, env.Value)
				}
				if env.Name == "OS_SENSOR_DATA__INTERVAL" {
					foundSensorInterval = true
					assert.Equal(t, tc.expectedSensorInterval, env.Value)
				}
			}

			assert.Equal(t, expectExporter, foundSendSensorData,
				"SEND_SENSOR_DATA env var presence mismatch")
			assert.Equal(t, expectExporter, foundSensorInterval,
				"OS_SENSOR_DATA__INTERVAL env var presence mismatch")
		})
	}
}

func TestAppendAgentImageEnvVars(t *testing.T) {
	testCases := []struct {
		name            string
		images          []metal3api.AgentImages
		expectedEnvVars map[string]string
		expectNoEnvVars bool
	}{
		{
			name:            "empty images",
			images:          nil,
			expectNoEnvVars: true,
		},
		{
			name: "single architecture x86_64",
			images: []metal3api.AgentImages{
				{
					Kernel:       "file:///shared/html/images/ipa.x86_64.kernel",
					Initramfs:    "file:///shared/html/images/ipa.x86_64.initramfs",
					Architecture: metal3api.ArchX86_64,
				},
			},
			expectedEnvVars: map[string]string{
				"DEPLOY_KERNEL_BY_ARCH":  "x86_64:file:///shared/html/images/ipa.x86_64.kernel",
				"DEPLOY_RAMDISK_BY_ARCH": "x86_64:file:///shared/html/images/ipa.x86_64.initramfs",
			},
		},
		{
			name: "single architecture aarch64",
			images: []metal3api.AgentImages{
				{
					Kernel:       "file:///shared/html/images/ipa.aarch64.kernel",
					Initramfs:    "file:///shared/html/images/ipa.aarch64.initramfs",
					Architecture: metal3api.ArchAarch64,
				},
			},
			expectedEnvVars: map[string]string{
				"DEPLOY_KERNEL_BY_ARCH":  "aarch64:file:///shared/html/images/ipa.aarch64.kernel",
				"DEPLOY_RAMDISK_BY_ARCH": "aarch64:file:///shared/html/images/ipa.aarch64.initramfs",
			},
		},
		{
			name: "multiple architectures",
			images: []metal3api.AgentImages{
				{
					Kernel:       "file:///shared/html/images/ipa.x86_64.kernel",
					Initramfs:    "file:///shared/html/images/ipa.x86_64.initramfs",
					Architecture: metal3api.ArchX86_64,
				},
				{
					Kernel:       "file:///shared/html/images/ipa.aarch64.kernel",
					Initramfs:    "file:///shared/html/images/ipa.aarch64.initramfs",
					Architecture: metal3api.ArchAarch64,
				},
			},
			expectedEnvVars: map[string]string{
				"DEPLOY_KERNEL_BY_ARCH":  "x86_64:file:///shared/html/images/ipa.x86_64.kernel,aarch64:file:///shared/html/images/ipa.aarch64.kernel",
				"DEPLOY_RAMDISK_BY_ARCH": "x86_64:file:///shared/html/images/ipa.x86_64.initramfs,aarch64:file:///shared/html/images/ipa.aarch64.initramfs",
			},
		},
		{
			name: "whitespace trimming",
			images: []metal3api.AgentImages{
				{
					Kernel:       "  file:///path/kernel  ",
					Initramfs:    "\tfile:///path/initramfs\n",
					Architecture: metal3api.ArchX86_64,
				},
			},
			expectedEnvVars: map[string]string{
				"DEPLOY_KERNEL_BY_ARCH":  "x86_64:file:///path/kernel",
				"DEPLOY_RAMDISK_BY_ARCH": "x86_64:file:///path/initramfs",
			},
		},
		{
			name: "default (empty architecture)",
			images: []metal3api.AgentImages{
				{
					Kernel:    "file:///shared/html/images/ipa.kernel",
					Initramfs: "file:///shared/html/images/ipa.initramfs",
				},
			},
			expectedEnvVars: map[string]string{
				"DEPLOY_KERNEL_URL":  "file:///shared/html/images/ipa.kernel",
				"DEPLOY_RAMDISK_URL": "file:///shared/html/images/ipa.initramfs",
			},
		},
		{
			name: "default with architecture-specific",
			images: []metal3api.AgentImages{
				{
					Kernel:    "file:///shared/html/images/ipa.kernel",
					Initramfs: "file:///shared/html/images/ipa.initramfs",
				},
				{
					Kernel:       "file:///shared/html/images/ipa.x86_64.kernel",
					Initramfs:    "file:///shared/html/images/ipa.x86_64.initramfs",
					Architecture: metal3api.ArchX86_64,
				},
			},
			expectedEnvVars: map[string]string{
				"DEPLOY_KERNEL_URL":      "file:///shared/html/images/ipa.kernel",
				"DEPLOY_RAMDISK_URL":     "file:///shared/html/images/ipa.initramfs",
				"DEPLOY_KERNEL_BY_ARCH":  "x86_64:file:///shared/html/images/ipa.x86_64.kernel",
				"DEPLOY_RAMDISK_BY_ARCH": "x86_64:file:///shared/html/images/ipa.x86_64.initramfs",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			envVars := appendAgentImageEnvVars(nil, tc.images)

			if tc.expectNoEnvVars {
				assert.Empty(t, envVars)
				return
			}

			envMap := make(map[string]string, len(envVars))
			for _, env := range envVars {
				envMap[env.Name] = env.Value
			}

			for name, expectedValue := range tc.expectedEnvVars {
				assert.Equal(t, expectedValue, envMap[name], "env var %s", name)
			}
			assert.Len(t, envVars, len(tc.expectedEnvVars))
		})
	}
}

func TestHttpdProbeConfiguration(t *testing.T) {
	testCases := []struct {
		Scenario                   string
		Ironic                     metal3api.IronicSpec
		ExpectExecProbe            bool
		ExpectCurlFail             bool
		ExpectCustomProbe          bool
		ExpectCustomReadinessProbe bool
		ExpectExecReadinessProbe   bool
	}{
		{
			Scenario: "Default - no custom images",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					ImageServerPort: 8080,
				},
			},
			ExpectExecProbe: true,
			ExpectCurlFail:  true,
		},
		{
			Scenario: "Custom images - exec probe without HTTP success requirement",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					ImageServerPort: 8080,
				},
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Architecture: metal3api.ArchX86_64,
							Kernel:       "file:///custom/kernel",
							Initramfs:    "file:///custom/initramfs",
						},
					},
				},
			},
			ExpectExecProbe: true,
			ExpectCurlFail:  false,
		},
		{
			Scenario: "Custom images with downloader disabled",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					ImageServerPort: 8080,
				},
				DeployRamdisk: metal3api.DeployRamdisk{
					DisableDownloader: true,
				},
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Architecture: metal3api.ArchAarch64,
							Kernel:       "http://external.com/kernel",
							Initramfs:    "http://external.com/initramfs",
						},
					},
				},
			},
			ExpectExecProbe: true,
			ExpectCurlFail:  false,
		},
		{
			Scenario: "Custom images with explicit readiness probe override",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					ImageServerPort: 8080,
				},
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Architecture: metal3api.ArchX86_64,
							Kernel:       "file:///custom/kernel",
							Initramfs:    "file:///custom/initramfs",
						},
					},
					HttpdReadinessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							HTTPGet: &corev1.HTTPGetAction{
								Path: "/ready",
								Port: intstr.FromInt(8080),
							},
						},
					},
				},
			},
			ExpectExecProbe:            true, // liveness stays exec default
			ExpectCustomReadinessProbe: true,
		},
		{
			Scenario: "Custom images with explicit liveness probe override",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					ImageServerPort: 8080,
				},
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Architecture: metal3api.ArchX86_64,
							Kernel:       "file:///custom/kernel",
							Initramfs:    "file:///custom/initramfs",
						},
					},
					HttpdLivenessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							HTTPGet: &corev1.HTTPGetAction{
								Path: "/health",
								Port: intstr.FromInt(8080),
							},
						},
					},
				},
			},
			ExpectCustomProbe:        true,
			ExpectExecReadinessProbe: true, // readiness stays exec when only liveness is overridden
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			cctx := ControllerContext{}
			secret := &corev1.Secret{
				Data: map[string][]byte{
					"htpasswd": []byte("test"),
				},
			}
			ironic := &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test",
					Name:      "test",
				},
				Spec: tc.Ironic,
			}

			resources := Resources{Ironic: ironic, APISecret: secret}
			podTemplate, err := newIronicPodTemplate(cctx, resources)
			require.NoError(t, err)

			// Find the httpd container
			var httpdContainer *corev1.Container
			for i := range podTemplate.Spec.Containers {
				if podTemplate.Spec.Containers[i].Name == "httpd" {
					httpdContainer = &podTemplate.Spec.Containers[i]
					break
				}
			}
			require.NotNil(t, httpdContainer, "httpd container should exist")
			require.NotNil(t, httpdContainer.LivenessProbe, "liveness probe should not be nil")
			require.NotNil(t, httpdContainer.ReadinessProbe, "readiness probe should not be nil")

			if tc.ExpectCustomProbe {
				assert.NotNil(t, httpdContainer.LivenessProbe.HTTPGet)
				assert.Equal(t, "/health", httpdContainer.LivenessProbe.HTTPGet.Path)
			} else if tc.ExpectExecProbe {
				assert.NotNil(t, httpdContainer.LivenessProbe.Exec, "should have exec probe")
				assert.Nil(t, httpdContainer.LivenessProbe.HTTPGet, "should not have HTTPGet probe")
				if tc.ExpectCurlFail {
					assert.Contains(t, httpdContainer.LivenessProbe.Exec.Command, "--fail", "exec probe should include --fail for default images")
				} else {
					assert.NotContains(t, httpdContainer.LivenessProbe.Exec.Command, "--fail", "exec probe should not include --fail for custom images")
				}
			}

			if tc.ExpectCustomReadinessProbe {
				assert.NotNil(t, httpdContainer.ReadinessProbe.HTTPGet)
				assert.Equal(t, "/ready", httpdContainer.ReadinessProbe.HTTPGet.Path)
			} else if tc.ExpectExecProbe || tc.ExpectExecReadinessProbe {
				assert.NotNil(t, httpdContainer.ReadinessProbe.Exec, "should have exec readiness probe")
				assert.Nil(t, httpdContainer.ReadinessProbe.HTTPGet, "should not have HTTPGet readiness probe")
				if tc.ExpectCurlFail {
					assert.Contains(t, httpdContainer.ReadinessProbe.Exec.Command, "--fail", "exec readiness probe should include --fail for default images")
				} else {
					assert.NotContains(t, httpdContainer.ReadinessProbe.Exec.Command, "--fail", "exec readiness probe should not include --fail for custom images")
				}
			}
		})
	}
}

func TestBuildTrustedCAEnvVars(t *testing.T) {
	testCases := []struct {
		name          string
		trustedCARef  *metal3api.ResourceReferenceWithKey
		configMapData map[string]string
		secretData    map[string][]byte
		expectedKey   string
	}{
		{
			name: "ConfigMap with specific key",
			trustedCARef: &metal3api.ResourceReferenceWithKey{
				ResourceReference: metal3api.ResourceReference{
					Name: "trusted-ca",
					Kind: "ConfigMap",
				},
				Key: "custom-ca.crt",
			},
			configMapData: map[string]string{
				"custom-ca.crt": "cert1",
				"other-ca.crt":  "cert2",
			},
			expectedKey: "custom-ca.crt",
		},
		{
			name: "Secret with specific key",
			trustedCARef: &metal3api.ResourceReferenceWithKey{
				ResourceReference: metal3api.ResourceReference{
					Name: "trusted-ca-secret",
					Kind: "Secret",
				},
				Key: "tls.crt",
			},
			secretData: map[string][]byte{
				"tls.crt": []byte("cert1"),
				"ca.crt":  []byte("cert2"),
			},
			expectedKey: "tls.crt",
		},
		{
			name: "Multiple keys without Key specified - ConfigMap",
			trustedCARef: &metal3api.ResourceReferenceWithKey{
				ResourceReference: metal3api.ResourceReference{
					Name: "trusted-ca",
					Kind: "ConfigMap",
				},
			},
			configMapData: map[string]string{
				"ca1.crt": "cert1",
				"ca2.crt": "cert2",
			},
			expectedKey: "ca1.crt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := logr.New(logr.Discard().GetSink())

			cctx := ControllerContext{
				Logger: logger,
			}

			ironic := &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test",
				},
				Spec: metal3api.IronicSpec{
					TLS: metal3api.TLS{
						TrustedCA: tc.trustedCARef,
					},
				},
			}

			resources := Resources{
				Ironic: ironic,
			}

			if tc.configMapData != nil {
				resources.TrustedCAConfigMap = &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      tc.trustedCARef.Name,
						Namespace: "test",
					},
					Data: tc.configMapData,
				}
			}

			if tc.secretData != nil {
				resources.TrustedCASecret = &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      tc.trustedCARef.Name,
						Namespace: "test",
					},
					Data: tc.secretData,
				}
			}

			envVars := buildTrustedCAEnvVars(cctx, resources)

			// Should always return exactly one env var
			require.Len(t, envVars, 1, "Should return exactly one environment variable")
			assert.Equal(t, "WEBSERVER_CACERT_FILE", envVars[0].Name)

			// Verify the path contains the expected key
			expectedPath := "/certs/ca/trusted/" + tc.expectedKey
			assert.Equal(t, expectedPath, envVars[0].Value, "Environment variable value mismatch")
		})
	}
}

func TestKeepalivedEnvVars(t *testing.T) {
	testCases := []struct {
		name                     string
		ironic                   metal3api.IronicSpec
		expectedKeepalivedVIPEnv string
		expectedProvisioningIP   string
		expectedProvInterface    string
	}{
		{
			name: "legacy single-IP mode",
			ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					Interface:        "eth0",
					IPAddress:        "192.0.2.2",
					IPAddressManager: metal3api.IPAddressManagerKeepalived,
				},
			},
			expectedProvisioningIP: "192.0.2.2",
			expectedProvInterface:  "eth0",
		},
		{
			name: "multi-IP mode",
			ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					Interface: "eth0",
					IPAddress: "192.0.2.2",
					Keepalived: []metal3api.KeepalivedIP{
						{IPAddress: "192.0.2.2", Interface: "eth0"},
						{IPAddress: "192.168.1.50", Interface: "eth1"},
					},
				},
			},
			expectedKeepalivedVIPEnv: "192.0.2.2,eth0 192.168.1.50,eth1",
		},
		{
			name: "multi-IP mode single entry",
			ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					Interface: "eth0",
					IPAddress: "192.0.2.2",
					Keepalived: []metal3api.KeepalivedIP{
						{IPAddress: "192.0.2.2", Interface: "eth0"},
					},
				},
			},
			expectedKeepalivedVIPEnv: "192.0.2.2,eth0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cctx := ControllerContext{}
			secret := &corev1.Secret{
				Data: map[string][]byte{
					"htpasswd": []byte("abcd"),
				},
			}
			ironic := &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test",
					Name:      "test",
				},
				Spec: tc.ironic,
			}

			resources := Resources{Ironic: ironic, APISecret: secret}
			podTemplate, err := newIronicPodTemplate(cctx, resources)
			require.NoError(t, err)

			var keepalivedContainer *corev1.Container
			for i := range podTemplate.Spec.Containers {
				if podTemplate.Spec.Containers[i].Name == "keepalived" {
					keepalivedContainer = &podTemplate.Spec.Containers[i]
					break
				}
			}
			require.NotNil(t, keepalivedContainer, "keepalived container should exist")

			envMap := make(map[string]string)
			for _, env := range keepalivedContainer.Env {
				envMap[env.Name] = env.Value
			}

			if tc.expectedKeepalivedVIPEnv != "" {
				assert.Equal(t, tc.expectedKeepalivedVIPEnv, envMap["KEEPALIVED_VIRTUAL_IPS"])
				assert.Empty(t, envMap["PROVISIONING_IP"], "PROVISIONING_IP should not be set in multi-IP mode")
				assert.Empty(t, envMap["PROVISIONING_INTERFACE"], "PROVISIONING_INTERFACE should not be set in multi-IP mode")
			} else {
				assert.Equal(t, tc.expectedProvisioningIP, envMap["PROVISIONING_IP"])
				assert.Equal(t, tc.expectedProvInterface, envMap["PROVISIONING_INTERFACE"])
				assert.Empty(t, envMap["KEEPALIVED_VIRTUAL_IPS"], "KEEPALIVED_VIRTUAL_IPS should not be set in legacy mode")
			}
		})
	}
}

func TestImageServerIPAddress(t *testing.T) {
	testCases := []struct {
		name            string
		ironic          metal3api.IronicSpec
		expectedHTTPURL string
		expectNoHTTPURL bool
	}{
		{
			name: "not set",
			ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					IPAddress: "192.0.2.2",
				},
			},
			expectNoHTTPURL: true,
		},
		{
			name: "set without TLS",
			ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					IPAddress:            "192.0.2.2",
					ImageServerIPAddress: "192.168.1.50",
					ImageServerPort:      6180,
				},
			},
			expectedHTTPURL: "http://192.168.1.50:6180",
		},
		{
			name: "set with TLS",
			ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					IPAddress:            "192.0.2.2",
					ImageServerIPAddress: "192.168.1.50",
					ImageServerTLSPort:   6183,
				},
				TLS: metal3api.TLS{
					CertificateName: "ironic-tls",
				},
			},
			expectedHTTPURL: "https://192.168.1.50:6183",
		},
		{
			name: "set with TLS but virtual media TLS disabled",
			ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					IPAddress:            "192.0.2.2",
					ImageServerIPAddress: "192.168.1.50",
					ImageServerPort:      6180,
				},
				TLS: metal3api.TLS{
					CertificateName:        "ironic-tls",
					DisableVirtualMediaTLS: true,
				},
			},
			expectedHTTPURL: "http://192.168.1.50:6180",
		},
		{
			name: "set with IPv6 address",
			ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					IPAddress:            "192.0.2.2",
					ImageServerIPAddress: "fd00::1",
					ImageServerPort:      6180,
				},
			},
			expectedHTTPURL: "http://[fd00::1]:6180",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cctx := ControllerContext{}
			secret := &corev1.Secret{
				Data: map[string][]byte{
					"htpasswd": []byte("abcd"),
				},
			}
			ironic := &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test",
					Name:      "test",
				},
				Spec: tc.ironic,
			}

			var tlsSecret *corev1.Secret
			if tc.ironic.TLS.CertificateName != "" {
				tlsSecret = &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name: tc.ironic.TLS.CertificateName,
					},
					Data: map[string][]byte{
						"tls.crt": []byte("cert"),
						"tls.key": []byte("key"),
					},
				}
			}

			resources := Resources{Ironic: ironic, APISecret: secret, TLSSecret: tlsSecret}
			podTemplate, err := newIronicPodTemplate(cctx, resources)
			require.NoError(t, err)

			var ironicContainer *corev1.Container
			for i := range podTemplate.Spec.Containers {
				if podTemplate.Spec.Containers[i].Name == ironicContainerName {
					ironicContainer = &podTemplate.Spec.Containers[i]
					break
				}
			}
			require.NotNil(t, ironicContainer, "ironic container should exist")

			var httpURL string
			var foundHTTPURL bool
			for _, env := range ironicContainer.Env {
				if env.Name == "IRONIC_HTTP_URL" {
					foundHTTPURL = true
					httpURL = env.Value
					break
				}
			}

			if tc.expectNoHTTPURL {
				assert.False(t, foundHTTPURL, "IRONIC_HTTP_URL should not be set")
			} else {
				assert.True(t, foundHTTPURL, "IRONIC_HTTP_URL should be set")
				assert.Equal(t, tc.expectedHTTPURL, httpURL)
			}
		})
	}
}

func TestBuildTrustedCAEnvVarsKeySelection(t *testing.T) {
	// Test that verifies the key selection logic directly
	testCases := []struct {
		name          string
		specifiedKey  string
		availableKeys []string
		expectedKey   string
	}{
		{
			name:          "Specified key exists",
			specifiedKey:  "my-ca.crt",
			availableKeys: []string{"other.crt", "my-ca.crt"},
			expectedKey:   "my-ca.crt",
		},
		{
			name:          "No key specified uses first",
			specifiedKey:  "",
			availableKeys: []string{"ignored.crt", "actual.crt"},
			expectedKey:   "actual.crt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cctx := ControllerContext{
				Logger: logr.Discard(),
			}

			trustedCARef := &metal3api.ResourceReferenceWithKey{
				ResourceReference: metal3api.ResourceReference{
					Name: "test-ca",
					Kind: "ConfigMap",
				},
				Key: tc.specifiedKey,
			}

			data := make(map[string]string)
			for _, key := range tc.availableKeys {
				data[key] = "cert-data"
			}

			ironic := &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "test",
				},
				Spec: metal3api.IronicSpec{
					TLS: metal3api.TLS{
						TrustedCA: trustedCARef,
					},
				},
			}

			resources := Resources{
				Ironic: ironic,
				TrustedCAConfigMap: &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "test-ca",
						Namespace: "test",
					},
					Data: data,
				},
			}

			envVars := buildTrustedCAEnvVars(cctx, resources)
			require.Len(t, envVars, 1)

			expectedPath := "/certs/ca/trusted/" + tc.expectedKey
			assert.Equal(t, expectedPath, envVars[0].Value)
		})
	}
}
