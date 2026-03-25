package ironic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

func TestNetworkingDeploymentName(t *testing.T) {
	ironic := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-ironic",
		},
	}

	name := NetworkingDeploymentName(ironic)
	assert.Equal(t, "test-ironic-networking", name)
}

func TestNetworkingServiceName(t *testing.T) {
	ironic := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-ironic",
		},
	}

	name := NetworkingServiceName(ironic)
	assert.Equal(t, "test-ironic-networking-service", name)
}

func TestNetworkingServiceEndpoint(t *testing.T) {
	ironic := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ironic",
			Namespace: "test-namespace",
		},
	}

	endpoint := NetworkingServiceEndpoint(ironic)
	assert.Equal(t, "test-ironic-networking-service.test-namespace.svc.cluster.local", endpoint)
}

func TestBuildNetworkingDeployment(t *testing.T) {
	testCases := []struct {
		Scenario                       string
		Ironic                         *metal3api.Ironic
		APISecret                      *corev1.Secret
		TLSSecret                      *corev1.Secret
		SwitchConfigSecret             *corev1.Secret
		SwitchCredentialsSecret        *corev1.Secret
		TrustedCAConfigMap             *corev1.ConfigMap
		CustomIronicImage              string
		ExpectTLSVolume                bool
		ExpectTrustedCAVolume          bool
		ExpectSwitchCredentialsVolume  bool
		ExpectedSwitchConfigSecretName string
	}{
		{
			Scenario: "basic networking deployment",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-namespace",
				},
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled:       true,
						RPCPort:       6190,
						SwitchDrivers: []string{"generic-switch"},
					},
				},
			},
			SwitchConfigSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "test-ironic-switch-config",
					ResourceVersion: "123",
				},
			},
			APISecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-api-cert",
				},
			},
			ExpectTLSVolume:       false,
			ExpectTrustedCAVolume: false,
		},
		{
			Scenario: "networking deployment with TLS",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-namespace",
				},
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled:       true,
						RPCPort:       6190,
						SwitchDrivers: []string{"generic-switch"},
					},
					TLS: metal3api.TLS{
						CertificateName: "test-tls-cert",
					},
				},
			},
			APISecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-api-cert",
				},
			},
			TLSSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-tls-cert",
				},
			},
			SwitchConfigSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "test-ironic-switch-config",
					ResourceVersion: "123",
				},
			},
			ExpectTLSVolume:       true,
			ExpectTrustedCAVolume: false,
		},
		{
			Scenario: "networking deployment with TLS and trusted CA",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-namespace",
				},
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled:       true,
						RPCPort:       6190,
						SwitchDrivers: []string{"generic-switch"},
					},
					TLS: metal3api.TLS{
						CertificateName: "test-tls-cert",
						TrustedCAName:   "test-trusted-ca",
					},
				},
			},
			APISecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-api-cert",
				},
			},
			TLSSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-tls-cert",
				},
			},
			TrustedCAConfigMap: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-trusted-ca",
				},
			},
			SwitchConfigSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "test-ironic-switch-config",
					ResourceVersion: "123",
				},
			},
			ExpectTLSVolume:       true,
			ExpectTrustedCAVolume: true,
		},
		{
			Scenario: "networking deployment with switch credentials secret",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-namespace",
				},
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled:                     true,
						RPCPort:                     6190,
						SwitchDrivers:               []string{"generic-switch"},
						SwitchCredentialsSecretName: "my-switch-creds",
					},
				},
			},
			SwitchConfigSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "test-ironic-switch-config",
					ResourceVersion: "123",
				},
			},
			SwitchCredentialsSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "my-switch-creds",
					ResourceVersion: "456",
				},
			},
			APISecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-api-cert",
				},
			},
			ExpectSwitchCredentialsVolume: true,
		},
		{
			Scenario: "networking deployment with custom ironic image",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-namespace",
				},
				Spec: metal3api.IronicSpec{
					Images: metal3api.Images{
						Ironic: "myorg/myironic:custom",
					},
					NetworkingService: &metal3api.NetworkingService{
						Enabled:       true,
						RPCPort:       6190,
						SwitchDrivers: []string{"generic-switch"},
					},
				},
			},
			SwitchConfigSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "test-ironic-switch-config",
					ResourceVersion: "123",
				},
			},
			APISecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-api-cert",
				},
			},
			CustomIronicImage: "myorg/myironic:custom",
		},
		{
			Scenario: "networking deployment with custom switchConfigSecretName",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-namespace",
				},
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled:                true,
						RPCPort:                6190,
						SwitchDrivers:          []string{"generic-switch"},
						SwitchConfigSecretName: "custom-switch-config",
					},
				},
			},
			SwitchConfigSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "custom-switch-config",
					ResourceVersion: "789",
				},
			},
			APISecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-api-cert",
				},
			},
			ExpectedSwitchConfigSecretName: "custom-switch-config",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			ironicImage := "quay.io/metal3-io/ironic:latest"
			if tc.CustomIronicImage != "" {
				ironicImage = tc.CustomIronicImage
			}
			cctx := ControllerContext{
				VersionInfo: VersionInfo{
					IronicImage: ironicImage,
				},
			}

			resources := Resources{
				Ironic:                  tc.Ironic,
				APISecret:               tc.APISecret,
				TLSSecret:               tc.TLSSecret,
				SwitchConfigSecret:      tc.SwitchConfigSecret,
				SwitchCredentialsSecret: tc.SwitchCredentialsSecret,
				TrustedCAConfigMap:      tc.TrustedCAConfigMap,
			}

			deployment := BuildNetworkingDeployment(cctx, resources)

			// Verify deployment metadata
			assert.Equal(t, "test-ironic-networking", deployment.Name)
			assert.Equal(t, "test-namespace", deployment.Namespace)
			assert.Equal(t, "ironic-networking", deployment.Labels["app"])
			assert.Equal(t, "test-ironic", deployment.Labels["ironic.metal3.io/instance"])

			// Verify replica count
			require.NotNil(t, deployment.Spec.Replicas)
			assert.Equal(t, int32(1), *deployment.Spec.Replicas)

			// Verify pod template labels
			assert.Equal(t, deployment.Labels, deployment.Spec.Template.Labels)

			// Verify pod annotations include switch config secret version
			assert.Contains(t, deployment.Spec.Template.Annotations, "ironic.metal3.io/switch-config-version")

			// Verify TLS secret version annotation when TLS is enabled
			if tc.ExpectTLSVolume {
				assert.Contains(t, deployment.Spec.Template.Annotations, "ironic.metal3.io/tls-secret-version")
			} else {
				assert.NotContains(t, deployment.Spec.Template.Annotations, "ironic.metal3.io/tls-secret-version")
			}

			// Verify selector matches labels
			assert.Equal(t, deployment.Labels, deployment.Spec.Selector.MatchLabels)

			// Verify container
			require.Len(t, deployment.Spec.Template.Spec.Containers, 1)
			container := deployment.Spec.Template.Spec.Containers[0]
			assert.Equal(t, "ironic-networking", container.Name)
			if tc.CustomIronicImage != "" {
				assert.Equal(t, tc.CustomIronicImage, container.Image)
			} else {
				assert.Equal(t, "quay.io/metal3-io/ironic:latest", container.Image)
			}
			assert.Equal(t, []string{"/bin/runironic-networking"}, container.Command)

			// Verify container ports
			require.Len(t, container.Ports, 1)
			assert.Equal(t, "networking-rpc", container.Ports[0].Name)
			assert.Equal(t, int32(6190), container.Ports[0].ContainerPort)

			// Verify container environment variables
			envMap := make(map[string]string)
			for _, env := range container.Env {
				envMap[env.Name] = env.Value
			}
			assert.Equal(t, "true", envMap["IRONIC_NETWORKING_ENABLED"])
			assert.Equal(t, "0.0.0.0", envMap["IRONIC_NETWORKING_JSON_RPC_HOST"])
			assert.Equal(t, "6190", envMap["IRONIC_NETWORKING_JSON_RPC_PORT"])
			assert.Equal(t, "/etc/ironic/networking/configs/switch-configs.conf", envMap["IRONIC_NETWORKING_SWITCH_CONFIGS"])
			assert.Equal(t, "generic-switch", envMap["IRONIC_NETWORKING_ENABLED_SWITCH_DRIVERS"])

			// Verify security context
			require.NotNil(t, container.SecurityContext)
			assert.Equal(t, int64(997), *container.SecurityContext.RunAsUser)
			assert.Equal(t, int64(994), *container.SecurityContext.RunAsGroup)
			require.NotNil(t, container.SecurityContext.Capabilities)
			assert.Contains(t, container.SecurityContext.Capabilities.Drop, corev1.Capability("ALL"))

			// Verify hostNetwork is false (no privileged networking)
			assert.False(t, deployment.Spec.Template.Spec.HostNetwork)

			// Verify volumes
			volumeMap := make(map[string]corev1.Volume)
			for _, vol := range deployment.Spec.Template.Spec.Volumes {
				volumeMap[vol.Name] = vol
			}

			// Networking pod should NOT have ironic-shared or database volumes
			assert.NotContains(t, volumeMap, "ironic-shared")

			// ironic-auth is required for JSON-RPC authentication
			assert.Contains(t, volumeMap, "ironic-auth")

			// Always present: switch-config
			assert.Contains(t, volumeMap, "switch-config")
			if tc.ExpectedSwitchConfigSecretName != "" {
				assert.Equal(t, tc.ExpectedSwitchConfigSecretName, volumeMap["switch-config"].Secret.SecretName)
			} else {
				assert.Equal(t, "test-ironic-switch-config", volumeMap["switch-config"].Secret.SecretName)
			}

			// Switch credentials volume
			if tc.ExpectSwitchCredentialsVolume {
				assert.Contains(t, volumeMap, "switch-credentials")
				assert.Equal(t, tc.Ironic.Spec.NetworkingService.SwitchCredentialsSecretName,
					volumeMap["switch-credentials"].Secret.SecretName)
				assert.Contains(t, deployment.Spec.Template.Annotations, "ironic.metal3.io/switch-credentials-version")
			} else {
				assert.NotContains(t, volumeMap, "switch-credentials")
			}

			if tc.ExpectTLSVolume {
				assert.Contains(t, volumeMap, "cert-ironic")
				assert.Equal(t, "test-tls-cert", volumeMap["cert-ironic"].Secret.SecretName)
			}

			if tc.ExpectTrustedCAVolume {
				assert.Contains(t, volumeMap, "trusted-ca")
				assert.Equal(t, "test-trusted-ca", volumeMap["trusted-ca"].ConfigMap.Name)
			} else {
				assert.NotContains(t, volumeMap, "trusted-ca")
			}

			// Verify volume mounts
			findMount := func(mounts []corev1.VolumeMount, name, path string) *corev1.VolumeMount {
				for _, mount := range mounts {
					if mount.Name == name && mount.MountPath == path {
						return &mount
					}
				}
				return nil
			}

			// Networking pod should NOT have ironic-shared or vmedia mounts
			assert.Nil(t, findMount(container.VolumeMounts, "ironic-shared", "/shared"),
				"ironic-shared should not be mounted in networking pod")
			assert.Nil(t, findMount(container.VolumeMounts, "cert-ironic", "/certs/vmedia"),
				"vmedia cert should not be mounted in networking pod")

			// Always present: switch-config
			switchConfigMount := findMount(container.VolumeMounts, "switch-config", "/etc/ironic/networking/configs")
			require.NotNil(t, switchConfigMount, "Expected switch-config volume to be mounted at /etc/ironic/networking")
			assert.True(t, switchConfigMount.ReadOnly)

			if tc.ExpectSwitchCredentialsVolume {
				credsMount := findMount(container.VolumeMounts, "switch-credentials", "/etc/ironic/networking/credentials")
				require.NotNil(t, credsMount, "Expected switch-credentials volume to be mounted at /etc/ironic/networking/credentials")
				assert.True(t, credsMount.ReadOnly)
			}

			if tc.ExpectTLSVolume {
				// Only /certs/ironic, NOT /certs/vmedia
				ironicCertMount := findMount(container.VolumeMounts, "cert-ironic", "/certs/ironic")
				require.NotNil(t, ironicCertMount, "Expected cert-ironic volume to be mounted at /certs/ironic")
				assert.True(t, ironicCertMount.ReadOnly)
			}

			if tc.ExpectTrustedCAVolume {
				trustedCAMount := findMount(container.VolumeMounts, "trusted-ca", "/certs/ca/trusted")
				require.NotNil(t, trustedCAMount, "Expected trusted-ca volume to be mounted at /certs/ca/trusted")
				assert.True(t, trustedCAMount.ReadOnly)
			}
		})
	}
}

func TestBuildNetworkingService(t *testing.T) {
	testCases := []struct {
		Scenario        string
		Ironic          *metal3api.Ironic
		ExpectedPort    int32
		ExpectedService string
	}{
		{
			Scenario: "default port",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-namespace",
				},
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled: true,
						RPCPort: 6190,
					},
				},
			},
			ExpectedPort:    6190,
			ExpectedService: "test-ironic-networking-service",
		},
		{
			Scenario: "custom port",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-ironic",
					Namespace: "my-namespace",
				},
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled: true,
						RPCPort: 7190,
					},
				},
			},
			ExpectedPort:    7190,
			ExpectedService: "my-ironic-networking-service",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			service := BuildNetworkingService(tc.Ironic)

			// Verify service metadata
			assert.Equal(t, tc.ExpectedService, service.Name)
			assert.Equal(t, tc.Ironic.Namespace, service.Namespace)
			assert.Equal(t, "ironic-networking", service.Labels["app"])
			assert.Equal(t, tc.Ironic.Name, service.Labels["ironic.metal3.io/instance"])

			// Verify service type
			assert.Equal(t, corev1.ServiceTypeClusterIP, service.Spec.Type)

			// Verify selector matches deployment labels
			expectedLabels := map[string]string{
				"app":                       "ironic-networking",
				"ironic.metal3.io/instance": tc.Ironic.Name,
			}
			assert.Equal(t, expectedLabels, service.Spec.Selector)

			// Verify ports
			require.Len(t, service.Spec.Ports, 1)
			assert.Equal(t, "networking-rpc", service.Spec.Ports[0].Name)
			assert.Equal(t, corev1.ProtocolTCP, service.Spec.Ports[0].Protocol)
			assert.Equal(t, tc.ExpectedPort, service.Spec.Ports[0].Port)
		})
	}
}

func TestGetNetworkingServiceEndpoint(t *testing.T) {
	testCases := []struct {
		Scenario         string
		Ironic           *metal3api.Ironic
		ExpectedEndpoint string
	}{
		{
			Scenario: "networking service disabled",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-namespace",
				},
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled: false,
					},
				},
			},
			ExpectedEndpoint: "",
		},
		{
			Scenario: "networking service nil",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-namespace",
				},
				Spec: metal3api.IronicSpec{
					NetworkingService: nil,
				},
			},
			ExpectedEndpoint: "",
		},
		{
			Scenario: "external endpoint provided",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-namespace",
				},
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled:  true,
						Endpoint: "external-networking.example.com",
					},
				},
			},
			ExpectedEndpoint: "external-networking.example.com",
		},
		{
			Scenario: "operator-managed networking service",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-namespace",
				},
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled: true,
						RPCPort: 6190,
					},
				},
			},
			ExpectedEndpoint: "test-ironic-networking-service.test-namespace.svc.cluster.local",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			endpoint := GetNetworkingServiceEndpoint(tc.Ironic)
			assert.Equal(t, tc.ExpectedEndpoint, endpoint)
		})
	}
}

func TestBuildNetworkingContainerEnv(t *testing.T) {
	testCases := []struct {
		Scenario              string
		Ironic                *metal3api.Ironic
		APISecret             *corev1.Secret
		TLSSecret             *corev1.Secret
		ExpectedEnvs          map[string]string
		NotExpectedEnvs       []string
		ExpectHTPasswdFromRef bool
	}{
		{
			Scenario: "basic configuration",
			Ironic: &metal3api.Ironic{
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled:       true,
						RPCPort:       6190,
						SwitchDrivers: []string{"generic-switch"},
					},
				},
			},
			ExpectedEnvs: map[string]string{
				"IRONIC_NETWORKING_ENABLED":                "true",
				"IRONIC_NETWORKING_JSON_RPC_HOST":          "0.0.0.0",
				"IRONIC_NETWORKING_JSON_RPC_PORT":          "6190",
				"IRONIC_NETWORKING_SWITCH_CONFIGS":         "/etc/ironic/networking/configs/switch-configs.conf",
				"IRONIC_NETWORKING_ENABLED_SWITCH_DRIVERS": "generic-switch",
			},
			NotExpectedEnvs: []string{"IRONIC_TLS_SETUP", "IRONIC_IP"},
		},
		{
			Scenario: "with IP address",
			Ironic: &metal3api.Ironic{
				Spec: metal3api.IronicSpec{
					Networking: metal3api.Networking{
						IPAddress: "192.168.1.100",
					},
					NetworkingService: &metal3api.NetworkingService{
						Enabled: true,
						RPCPort: 6190,
					},
				},
			},
			ExpectedEnvs: map[string]string{
				"IRONIC_IP": "192.168.1.100",
			},
		},
		{
			Scenario: "with TLS",
			Ironic: &metal3api.Ironic{
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled: true,
						RPCPort: 6190,
					},
					TLS: metal3api.TLS{
						CertificateName: "test-cert",
					},
				},
			},
			TLSSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cert",
				},
			},
		},
		{
			Scenario: "with multiple switch drivers",
			Ironic: &metal3api.Ironic{
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled:       true,
						RPCPort:       6190,
						SwitchDrivers: []string{"generic-switch", "netmiko"},
					},
				},
			},
			ExpectedEnvs: map[string]string{
				"IRONIC_NETWORKING_ENABLED_SWITCH_DRIVERS": "generic-switch,netmiko",
			},
		},
		{
			Scenario: "InsecureRPC true with TLS",
			Ironic: &metal3api.Ironic{
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled: true,
						RPCPort: 6190,
					},
					TLS: metal3api.TLS{
						CertificateName: "test-cert",
						InsecureRPC:     ptr.To(true),
					},
				},
			},
			TLSSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cert",
				},
			},
			ExpectedEnvs: map[string]string{
				"IRONIC_INSECURE": "true",
			},
		},
		{
			Scenario: "InsecureRPC false with TLS",
			Ironic: &metal3api.Ironic{
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled: true,
						RPCPort: 6190,
					},
					TLS: metal3api.TLS{
						CertificateName: "test-cert",
						InsecureRPC:     ptr.To(false),
					},
				},
			},
			TLSSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-cert",
				},
			},
			NotExpectedEnvs: []string{"IRONIC_INSECURE"},
		},
		{
			Scenario: "InsecureRPC true without TLS secret",
			Ironic: &metal3api.Ironic{
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled: true,
						RPCPort: 6190,
					},
					TLS: metal3api.TLS{
						InsecureRPC: ptr.To(true),
					},
				},
			},
			NotExpectedEnvs: []string{"IRONIC_INSECURE"},
		},
		{
			Scenario: "API secret present sets IRONIC_HTPASSWD from secret ref",
			Ironic: &metal3api.Ironic{
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled: true,
						RPCPort: 6190,
					},
				},
			},
			APISecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-api-secret",
				},
			},
			ExpectHTPasswdFromRef: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			resources := Resources{
				Ironic:    tc.Ironic,
				APISecret: tc.APISecret,
				TLSSecret: tc.TLSSecret,
			}

			envVars := buildNetworkingContainerEnv(resources)

			// Build map for easier checking
			envMap := make(map[string]string)
			envVarList := make(map[string]corev1.EnvVar)
			for _, env := range envVars {
				envMap[env.Name] = env.Value
				envVarList[env.Name] = env
			}

			// Check expected env vars
			for key, expectedValue := range tc.ExpectedEnvs {
				assert.Contains(t, envMap, key, "Expected env var %s to be present", key)
				assert.Equal(t, expectedValue, envMap[key], "Expected env var %s to have value %s", key, expectedValue)
			}

			// Check that not expected env vars are not present
			for _, key := range tc.NotExpectedEnvs {
				assert.NotContains(t, envMap, key, "Did not expect env var %s to be present", key)
			}

			// Check IRONIC_HTPASSWD is set from a secret ref
			if tc.ExpectHTPasswdFromRef {
				htpasswdEnv, ok := envVarList["IRONIC_HTPASSWD"]
				require.True(t, ok, "Expected IRONIC_HTPASSWD env var to be present")
				require.NotNil(t, htpasswdEnv.ValueFrom, "Expected IRONIC_HTPASSWD to have ValueFrom")
				require.NotNil(t, htpasswdEnv.ValueFrom.SecretKeyRef, "Expected IRONIC_HTPASSWD to have SecretKeyRef")
				assert.Equal(t, tc.APISecret.Name, htpasswdEnv.ValueFrom.SecretKeyRef.Name)
				assert.Equal(t, "htpasswd", htpasswdEnv.ValueFrom.SecretKeyRef.Key)
			}
		})
	}
}

func TestEnsureSwitchConfigSecret(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, metal3api.AddToScheme(scheme))

	testCases := []struct {
		Scenario       string
		Ironic         *metal3api.Ironic
		ExistingSecret *corev1.Secret
		ExpectCreate   bool
	}{
		{
			Scenario: "creates empty secret when none exists",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-ns",
				},
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled: true,
					},
				},
			},
			ExistingSecret: nil,
			ExpectCreate:   true,
		},
		{
			Scenario: "does not overwrite existing secret data",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-ns",
					UID:       "abc-123",
				},
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled: true,
					},
				},
			},
			ExistingSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic-switch-config",
					Namespace: "test-ns",
				},
				Data: map[string][]byte{
					"switch-configs.conf": []byte("existing config data"),
				},
			},
			ExpectCreate: false,
		},
		{
			Scenario: "creates secret with custom switchConfigSecretName",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-ns",
				},
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						Enabled:                true,
						SwitchConfigSecretName: "custom-switch-config",
					},
				},
			},
			ExistingSecret: nil,
			ExpectCreate:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			var objects []client.Object
			if tc.ExistingSecret != nil {
				objects = append(objects, tc.ExistingSecret)
			}

			c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build()
			cctx := ControllerContext{
				Context: t.Context(),
				Client:  c,
				Scheme:  scheme,
			}

			err := EnsureSwitchConfigSecret(cctx, tc.Ironic)
			require.NoError(t, err)

			// Verify the secret exists
			secretName := SwitchConfigSecretName(tc.Ironic)
			secret := &corev1.Secret{}
			err = c.Get(t.Context(), client.ObjectKey{
				Namespace: tc.Ironic.Namespace,
				Name:      secretName,
			}, secret)
			require.NoError(t, err)
			assert.Contains(t, secret.Data, "switch-configs.conf")

			if tc.ExpectCreate {
				// New secret should have empty config, owner reference, and managed label
				assert.Empty(t, string(secret.Data["switch-configs.conf"]))
				require.Len(t, secret.OwnerReferences, 1)
				assert.Equal(t, tc.Ironic.Name, secret.OwnerReferences[0].Name)
				assert.Equal(t, "true", secret.Labels["ironic.metal3.io/managed"])
			} else {
				// Existing secret should be unchanged
				assert.Equal(t, "existing config data", string(secret.Data["switch-configs.conf"]))
			}
		})
	}
}

func TestEnsureSwitchConfigSecretDeleted(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, metal3api.AddToScheme(scheme))

	testCases := []struct {
		Scenario       string
		Ironic         *metal3api.Ironic
		ExistingSecret *corev1.Secret
		ExpectDeleted  bool
	}{
		{
			Scenario: "delete operator-managed switch config secret",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-ns",
				},
			},
			ExistingSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic-switch-config",
					Namespace: "test-ns",
					Labels: map[string]string{
						"ironic.metal3.io/managed": "true",
					},
				},
				Data: map[string][]byte{
					"switch-configs.conf": []byte("config data"),
				},
			},
			ExpectDeleted: true,
		},
		{
			Scenario: "skip externally-managed secret",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-ns",
				},
			},
			ExistingSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic-switch-config",
					Namespace: "test-ns",
				},
				Data: map[string][]byte{
					"switch-configs.conf": []byte("external config"),
				},
			},
			ExpectDeleted: false,
		},
		{
			Scenario: "no-op when secret doesn't exist",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test-ns",
				},
			},
			ExistingSecret: nil,
			ExpectDeleted:  true, // vacuously true - nothing to check
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			var objects []client.Object
			if tc.ExistingSecret != nil {
				objects = append(objects, tc.ExistingSecret)
			}

			c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build()
			cctx := ControllerContext{
				Context: t.Context(),
				Client:  c,
				Scheme:  scheme,
			}

			err := EnsureSwitchConfigSecretDeleted(cctx, tc.Ironic)
			require.NoError(t, err)

			secretName := tc.Ironic.Name + "-switch-config"
			secret := &corev1.Secret{}
			getErr := c.Get(t.Context(), client.ObjectKey{
				Namespace: tc.Ironic.Namespace,
				Name:      secretName,
			}, secret)

			switch {
			case tc.ExistingSecret == nil:
				// Secret never existed
				assert.True(t, k8serrors.IsNotFound(getErr))
			case tc.ExpectDeleted:
				// Operator-managed secret should be gone
				assert.True(t, k8serrors.IsNotFound(getErr))
			default:
				// External secret should still exist with original data
				require.NoError(t, getErr)
				assert.Equal(t, "external config", string(secret.Data["switch-configs.conf"]))
			}
		})
	}
}

func TestSwitchConfigSecretName(t *testing.T) {
	testCases := []struct {
		Scenario     string
		Ironic       *metal3api.Ironic
		ExpectedName string
	}{
		{
			Scenario: "default name",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-ironic",
				},
			},
			ExpectedName: "my-ironic-switch-config",
		},
		{
			Scenario: "custom name via spec",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-ironic",
				},
				Spec: metal3api.IronicSpec{
					NetworkingService: &metal3api.NetworkingService{
						SwitchConfigSecretName: "custom-config",
					},
				},
			},
			ExpectedName: "custom-config",
		},
		{
			Scenario: "nil networking service",
			Ironic: &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name: "my-ironic",
				},
				Spec: metal3api.IronicSpec{
					NetworkingService: nil,
				},
			},
			ExpectedName: "my-ironic-switch-config",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			name := SwitchConfigSecretName(tc.Ironic)
			assert.Equal(t, tc.ExpectedName, name)
		})
	}
}
