package ironic

import (
	"fmt"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

const (
	defaultRegistry = "quay.io/metal3-io"
)

var (
	// NOTE(dtantsur): defaultVersion must be updated after branching.
	defaultVersion                = metal3api.VersionLatest
	defaultMariaDBImage           = defaultRegistry + "/mariadb:latest"
	defaultRamdiskDownloaderImage = defaultRegistry + "/ironic-ipa-downloader:latest"
	defaultKeepalivedImage        = defaultRegistry + "/keepalived:latest"

	versionUpgradeScripts      = metal3api.Version290
	versionMountDatabaseSecret = metal3api.Version290
	versionDataMounts          = metal3api.Version290
)

type VersionInfo struct {
	InstalledVersion       metal3api.Version
	IronicImage            string
	MariaDBImage           string
	RamdiskDownloaderImage string
	AgentBranch            string
	KeepalivedImage        string
}

// Creates a version info from images and version.
func NewVersionInfo(ironicImages metal3api.Images, ironicVersion string, databaseImage string) (result VersionInfo, err error) {
	if ironicVersion != "" {
		parsedVersion, err := metal3api.ParseVersion(ironicVersion)
		if err != nil {
			return VersionInfo{}, err
		}
		result.InstalledVersion = parsedVersion
	} else {
		result.InstalledVersion = defaultVersion
	}
	tag := metal3api.SupportedVersions[result.InstalledVersion]

	if ironicImages.Ironic != "" {
		result.IronicImage = ironicImages.Ironic
	} else {
		result.IronicImage = fmt.Sprintf("%s/ironic:%s", defaultRegistry, tag)
	}

	if ironicImages.DeployRamdiskDownloader != "" {
		result.RamdiskDownloaderImage = ironicImages.DeployRamdiskDownloader
	} else {
		result.RamdiskDownloaderImage = defaultRamdiskDownloaderImage
	}

	if ironicImages.Keepalived != "" {
		result.KeepalivedImage = ironicImages.Keepalived
	} else {
		result.KeepalivedImage = defaultKeepalivedImage
	}

	if ironicImages.DeployRamdiskBranch != "" {
		result.AgentBranch = ironicImages.DeployRamdiskBranch
	}

	if databaseImage != "" {
		result.MariaDBImage = databaseImage
	} else {
		result.MariaDBImage = defaultMariaDBImage
	}

	return
}

// Takes VersionInfo with defaults from the configuration and applies any overrides from the Ironic object.
// Explicit images from the Images object take priority. Otherwise, the defaults are taken from the hardcoded defaults for the given version.
func (versionInfo VersionInfo) WithIronicOverrides(ironic *metal3api.Ironic) (VersionInfo, error) {
	images := &ironic.Spec.Images

	if ironic.Spec.Version != "" {
		parsedVersion, err := metal3api.ParseVersion(ironic.Spec.Version)
		if err != nil {
			return VersionInfo{}, err
		}
		versionInfo.InstalledVersion = parsedVersion

		// NOTE(dtantsur): a non-default version requires a different default image
		if images.Ironic == "" {
			tag := metal3api.SupportedVersions[parsedVersion]
			versionInfo.IronicImage = fmt.Sprintf("%s/ironic:%s", defaultRegistry, tag)
		}
	}

	if images.DeployRamdiskBranch != "" {
		versionInfo.AgentBranch = images.DeployRamdiskBranch
	}

	if images.DeployRamdiskDownloader != "" {
		versionInfo.RamdiskDownloaderImage = images.DeployRamdiskDownloader
	}

	if images.Ironic != "" {
		versionInfo.IronicImage = images.Ironic
	}

	if images.Keepalived != "" {
		versionInfo.KeepalivedImage = images.Keepalived
	}

	return versionInfo, nil
}

// Takes VersionInfo with defaults from the configuration and applies any overrides from the IronicDatabase object.
func (versionInfo VersionInfo) WithIronicDatabaseOverrides(db *metal3api.IronicDatabase) VersionInfo {
	if db.Spec.Image != "" {
		versionInfo.MariaDBImage = db.Spec.Image
	}

	return versionInfo
}
