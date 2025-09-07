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
	assert.NotEqual(t, "", podTemplate.Annotations["ironic.metal3.io/api-secret-version"])
	assert.Equal(t, "my-label", podTemplate.Labels["label.example.com"])
	assert.Equal(t, "test", podTemplate.Labels[metal3api.IronicServiceLabel])
}
