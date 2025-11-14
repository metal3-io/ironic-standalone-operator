package ironic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

func TestWithIronicOverrides(t *testing.T) {
	testCases := []struct {
		Scenario string

		DefaultIronicImages  metal3api.Images
		DefaultIronicVersion string
		DefaultDatabaseImage string
		Ironic               metal3api.Ironic

		Expected    VersionInfo
		ExpectError string
	}{
		{
			Scenario: "only defaults",

			Expected: VersionInfo{
				// NOTE(dtantsur): this value will change on stable branches
				InstalledVersion:       metal3api.VersionLatest,
				IronicImage:            "quay.io/metal3-io/ironic:latest",
				KeepalivedImage:        "quay.io/metal3-io/keepalived:latest",
				RamdiskDownloaderImage: "quay.io/metal3-io/ironic-ipa-downloader:latest",
				MariaDBImage:           "quay.io/metal3-io/mariadb:latest",
			},
		},
		{
			Scenario: "explicit overrides",

			Ironic: metal3api.Ironic{
				Spec: metal3api.IronicSpec{
					Images: metal3api.Images{
						DeployRamdiskBranch:     "stable/x.y",
						DeployRamdiskDownloader: "myorg/ramdisk-downloader:tag",
						Ironic:                  "myorg/ironic:tag",
						Keepalived:              "myorg/keepalived:tag",
					},
				},
			},

			Expected: VersionInfo{
				AgentBranch: "stable/x.y",
				// NOTE(dtantsur): this value will change on stable branches
				InstalledVersion:       metal3api.VersionLatest,
				IronicImage:            "myorg/ironic:tag",
				KeepalivedImage:        "myorg/keepalived:tag",
				RamdiskDownloaderImage: "myorg/ramdisk-downloader:tag",
				MariaDBImage:           "quay.io/metal3-io/mariadb:latest",
			},
		},
		{
			Scenario: "latest version",

			Ironic: metal3api.Ironic{
				Spec: metal3api.IronicSpec{
					Version: "32.0",
				},
			},

			Expected: VersionInfo{
				InstalledVersion:       metal3api.Version320,
				IronicImage:            "quay.io/metal3-io/ironic:release-32.0",
				KeepalivedImage:        "quay.io/metal3-io/keepalived:latest",
				RamdiskDownloaderImage: "quay.io/metal3-io/ironic-ipa-downloader:latest",
				MariaDBImage:           "quay.io/metal3-io/mariadb:latest",
			},
		},
		{
			Scenario: "older version",

			Ironic: metal3api.Ironic{
				Spec: metal3api.IronicSpec{
					Version: "31.0",
				},
			},

			Expected: VersionInfo{
				InstalledVersion:       metal3api.Version310,
				IronicImage:            "quay.io/metal3-io/ironic:release-31.0",
				KeepalivedImage:        "quay.io/metal3-io/keepalived:latest",
				RamdiskDownloaderImage: "quay.io/metal3-io/ironic-ipa-downloader:latest",
				MariaDBImage:           "quay.io/metal3-io/mariadb:latest",
			},
		},
		{
			Scenario: "invalid version",

			Ironic: metal3api.Ironic{
				Spec: metal3api.IronicSpec{
					Version: "42",
				},
			},

			ExpectError: "invalid version 42",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			defaults, err := NewVersionInfo(tc.DefaultIronicImages, tc.DefaultIronicVersion, tc.DefaultDatabaseImage)
			require.NoError(t, err)
			result, err := defaults.WithIronicOverrides(&tc.Ironic)
			if tc.ExpectError != "" {
				assert.ErrorContains(t, err, tc.ExpectError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.Expected, result)
			}
		})
	}
}

func TestPrometheusExporterVersionCheck(t *testing.T) {
	testCases := []struct {
		name          string
		version       metal3api.Version
		enabled       bool
		expectedError string
	}{
		{
			name:          "PrometheusExporter with version 31.0",
			version:       metal3api.Version310,
			enabled:       true,
			expectedError: "",
		},
		{
			name:          "PrometheusExporter with version 32.0",
			version:       metal3api.Version320,
			enabled:       true,
			expectedError: "",
		},
		{
			name:          "PrometheusExporter with latest version",
			version:       metal3api.VersionLatest,
			enabled:       true,
			expectedError: "",
		},
		{
			name:          "PrometheusExporter with version 30.0 (too old)",
			version:       metal3api.Version300,
			enabled:       true,
			expectedError: "using prometheusExporter is only possible for Ironic 31.0 or newer",
		},
		{
			name:          "PrometheusExporter disabled with version 30.0",
			version:       metal3api.Version300,
			enabled:       false,
			expectedError: "",
		},
		{
			name:          "PrometheusExporter not configured with version 30.0",
			version:       metal3api.Version300,
			enabled:       false,
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var prometheusExporter *metal3api.PrometheusExporter
			if tc.enabled {
				prometheusExporter = &metal3api.PrometheusExporter{
					Enabled: tc.enabled,
				}
			}

			ironic := &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test",
				},
				Spec: metal3api.IronicSpec{
					PrometheusExporter: prometheusExporter,
				},
			}

			resources := Resources{
				Ironic: ironic,
			}

			err := checkVersion(resources, tc.version)
			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBMCCAVersionCheck(t *testing.T) {
	testCases := []struct {
		name          string
		version       metal3api.Version
		expectedError string
	}{
		{
			name:          "BMCCA with version 32.0",
			version:       metal3api.Version320,
			expectedError: "",
		},
		{
			name:          "BMCCA with latest version",
			version:       metal3api.VersionLatest,
			expectedError: "",
		},
		{
			name:          "BMCCA with version 31.0 (too old)",
			version:       metal3api.Version310,
			expectedError: "using tls.bmcCAName is only possible for Ironic 32.0 or newer",
		},
		{
			name:          "BMCCA with version 30.0 (too old)",
			version:       metal3api.Version300,
			expectedError: "using tls.bmcCAName is only possible for Ironic 32.0 or newer",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bmcSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bmc-ca",
					Namespace: "test",
				},
				Data: map[string][]byte{
					"ca.crt": []byte("test-ca-cert"),
				},
			}

			ironic := &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test",
				},
				Spec: metal3api.IronicSpec{
					TLS: metal3api.TLS{
						BMCCAName: "bmc-ca",
					},
				},
			}

			resources := Resources{
				Ironic:      ironic,
				BMCCASecret: bmcSecret,
			}

			err := checkVersion(resources, tc.version)

			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
