package ironic

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

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
		// This is currently set unconditionally by IrSO itself and will eventually be replaced by a proper ironic-image variable.
		"OS_JSON_RPC__PORT": "6189",
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
				if podTemplate.Spec.Containers[i].Name == "ironic" {
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

func TestPrometheusExporterEnvVars(t *testing.T) {
	testCases := []struct {
		name                   string
		prometheusExporter     *metal3api.PrometheusExporter
		expectedSendSensorData string
		expectedSensorInterval string
	}{
		{
			name: "PrometheusExporter enabled with default interval",
			prometheusExporter: &metal3api.PrometheusExporter{
				Enabled:                  true,
				SensorCollectionInterval: 0, // Should default to 60
			},
			expectedSendSensorData: "true",
			expectedSensorInterval: "60",
		},
		{
			name: "PrometheusExporter enabled with custom interval",
			prometheusExporter: &metal3api.PrometheusExporter{
				Enabled:                  true,
				SensorCollectionInterval: 120,
			},
			expectedSendSensorData: "true",
			expectedSensorInterval: "120",
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
				if cont.Name == "ironic" {
					ironicContainer = &contCopy
				}
				if cont.Name == "ironic-prometheus-exporter" {
					exporterContainer = &contCopy
				}
			}
			require.NotNil(t, ironicContainer, "ironic container not found")
			if expectExporter {
				require.NotNil(t, exporterContainer, "ironic-prometheus-exporter container not found")
				assert.Len(t, exporterContainer.Env, 1)
				assert.Equal(t, "FLASK_RUN_PORT", exporterContainer.Env[0].Name)
				assert.Equal(t, "9608", exporterContainer.Env[0].Value)
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
