package ironic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
			Scenario: "older version",

			Ironic: metal3api.Ironic{
				Spec: metal3api.IronicSpec{
					Version: "30.0",
				},
			},

			Expected: VersionInfo{
				InstalledVersion:       metal3api.Version300,
				IronicImage:            "quay.io/metal3-io/ironic:release-30.0",
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
