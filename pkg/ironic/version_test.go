package ironic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			},
		},
		{
			Scenario: "latest version",

			Ironic: metal3api.Ironic{
				Spec: metal3api.IronicSpec{
					Version: "38.0",
				},
			},

			Expected: VersionInfo{
				InstalledVersion:       metal3api.Version380,
				IronicImage:            "quay.io/metal3-io/ironic:release-38.0",
				KeepalivedImage:        "quay.io/metal3-io/keepalived:latest",
				RamdiskDownloaderImage: "quay.io/metal3-io/ironic-ipa-downloader:latest",
			},
		},
		{
			Scenario: "older version",

			Ironic: metal3api.Ironic{
				Spec: metal3api.IronicSpec{
					Version: "37.0",
				},
			},

			Expected: VersionInfo{
				InstalledVersion:       metal3api.Version370,
				IronicImage:            "quay.io/metal3-io/ironic:release-37.0",
				KeepalivedImage:        "quay.io/metal3-io/keepalived:latest",
				RamdiskDownloaderImage: "quay.io/metal3-io/ironic-ipa-downloader:latest",
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
			defaults, err := NewVersionInfo(tc.DefaultIronicImages, tc.DefaultIronicVersion)
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
			name:          "PrometheusExporter with version 35.0",
			version:       metal3api.Version350,
			enabled:       true,
			expectedError: "",
		},
		{
			name:          "PrometheusExporter with version 37.0",
			version:       metal3api.Version370,
			enabled:       true,
			expectedError: "",
		},
		{
			name:          "PrometheusExporter with version 38.0",
			version:       metal3api.Version380,
			enabled:       true,
			expectedError: "",
		},
		{
			name:          "PrometheusExporter with latest version",
			version:       metal3api.VersionLatest,
			enabled:       true,
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

			err := CheckVersion(resources, tc.version)
			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMultiRangeDHCPVersionCheck(t *testing.T) {
	testCases := []struct {
		name          string
		version       metal3api.Version
		dhcp          *metal3api.DHCP
		expectedError string
	}{
		{
			name:    "no DHCP at all",
			version: metal3api.Version350,
		},
		{
			name:    "DHCP without ExtraRanges on older version (allowed)",
			version: metal3api.Version350,
			dhcp: &metal3api.DHCP{
				NetworkCIDR: "192.0.2.0/24",
				RangeBegin:  "192.0.2.10",
				RangeEnd:    "192.0.2.100",
			},
		},
		{
			name:    "ExtraRanges on latest version",
			version: metal3api.VersionLatest,
			dhcp: &metal3api.DHCP{
				NetworkCIDR: "192.0.2.0/24",
				RangeBegin:  "192.0.2.10",
				RangeEnd:    "192.0.2.100",
				ExtraRanges: []metal3api.DHCPRange{
					{NetworkCIDR: "198.51.100.0/24", RangeBegin: "198.51.100.10", RangeEnd: "198.51.100.100"},
				},
			},
		},
		{
			name:    "ExtraRanges on 37.0 is allowed",
			version: metal3api.Version370,
			dhcp: &metal3api.DHCP{
				NetworkCIDR: "192.0.2.0/24",
				RangeBegin:  "192.0.2.10",
				RangeEnd:    "192.0.2.100",
				ExtraRanges: []metal3api.DHCPRange{
					{NetworkCIDR: "198.51.100.0/24", RangeBegin: "198.51.100.10", RangeEnd: "198.51.100.100"},
				},
			},
		},
		{
			name:    "ExtraRanges on 35.0 is rejected",
			version: metal3api.Version350,
			dhcp: &metal3api.DHCP{
				NetworkCIDR: "192.0.2.0/24",
				RangeBegin:  "192.0.2.10",
				RangeEnd:    "192.0.2.100",
				ExtraRanges: []metal3api.DHCPRange{
					{NetworkCIDR: "198.51.100.0/24", RangeBegin: "198.51.100.10", RangeEnd: "198.51.100.100"},
				},
			},
			expectedError: "networking.dhcp.extraRanges requires Ironic 37.0 or newer",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ironic := &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test"},
				Spec: metal3api.IronicSpec{
					Networking: metal3api.Networking{DHCP: tc.dhcp},
				},
			}
			err := CheckVersion(Resources{Ironic: ironic}, tc.version)
			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
