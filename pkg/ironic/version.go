package ironic

import (
	"fmt"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

const (
	defaultRegistry = "quay.io/metal3-io"
)

var (
	// NOTE(dtantsur): defaultVersion must be updated after branching
	defaultVersion                = metal3api.VersionLatest
	defaultMariaDBImage           = fmt.Sprintf("%s/mariadb:latest", defaultRegistry)
	defaultRamdiskDownloaderImage = fmt.Sprintf("%s/ironic-ipa-downloader:latest", defaultRegistry)
	defaultKeepalivedImage        = fmt.Sprintf("%s/keepalived:latest", defaultRegistry)
)

type VersionInfo struct {
	InstalledVersion       string
	IronicImage            string
	MariaDBImage           string
	RamdiskDownloaderImage string
	AgentBranch            string
	AgentDownloadURL       string
	KeepalivedImage        string
}

// Takes VersionInfo with defaults from the configuration and applies any overrides from the Ironic object.
// Explicit images from the Images object take priority. Otherwise, the defaults are taken from the hardcoded defaults for the given version.
func (versionInfo VersionInfo) WithIronicOverrides(ironic *metal3api.Ironic) (VersionInfo, error) {
	if ironic.Spec.Version != "" {
		versionInfo.InstalledVersion = ironic.Spec.Version
	} else if versionInfo.InstalledVersion == "" {
		versionInfo.InstalledVersion = defaultVersion.String()
	}

	parsedVersion, err := metal3api.ParseVersion(versionInfo.InstalledVersion)
	if err != nil {
		return VersionInfo{}, err
	}
	tag := metal3api.SupportedVersions[parsedVersion]

	defaults := VersionInfo{
		InstalledVersion: versionInfo.InstalledVersion,
		IronicImage:      fmt.Sprintf("%s/ironic:%s", defaultRegistry, tag),
		// MariaDBImage is not actually used here but is set for consistency.
		MariaDBImage:           defaultMariaDBImage,
		RamdiskDownloaderImage: defaultRamdiskDownloaderImage,
		KeepalivedImage:        defaultKeepalivedImage,
	}

	images := &ironic.Spec.Images

	if images.DeployRamdiskBranch != "" {
		versionInfo.AgentBranch = images.DeployRamdiskBranch
	} else if versionInfo.AgentBranch == "" {
		versionInfo.AgentBranch = defaults.AgentBranch
	}

	if images.DeployRamdiskDownloader != "" {
		versionInfo.RamdiskDownloaderImage = images.DeployRamdiskDownloader
	} else if versionInfo.RamdiskDownloaderImage == "" {
		versionInfo.RamdiskDownloaderImage = defaults.RamdiskDownloaderImage
	}

	if images.Ironic != "" {
		versionInfo.IronicImage = images.Ironic
	} else if versionInfo.IronicImage == "" {
		versionInfo.IronicImage = defaults.IronicImage
	}

	if images.Keepalived != "" {
		versionInfo.KeepalivedImage = images.Keepalived
	} else if versionInfo.KeepalivedImage == "" {
		versionInfo.KeepalivedImage = defaults.KeepalivedImage
	}

	return versionInfo, nil
}

// Takes VersionInfo with defaults from the configuration and applies any overrides from the IronicDatabase object.
func (versionInfo VersionInfo) WithIronicDatabaseOverrides(db *metal3api.IronicDatabase) VersionInfo {
	if db.Spec.Image != "" {
		versionInfo.MariaDBImage = db.Spec.Image
	} else if versionInfo.MariaDBImage == "" {
		versionInfo.MariaDBImage = defaultMariaDBImage
	}

	return versionInfo
}
