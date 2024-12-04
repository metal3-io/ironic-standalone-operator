package ironic

import metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"

const (
	installedVersion string = "latest"
	defaultRegistry  string = "quay.io/metal3-io/"
)

var (
	defaultIronicImage            = defaultRegistry + "ironic:" + installedVersion
	defaultMariaDBImage           = defaultRegistry + "mariadb:latest"
	defaultRamdiskDownloaderImage = defaultRegistry + "ironic-ipa-downloader:latest"
	defaultKeepalivedImage        = defaultRegistry + "keepalived:latest"
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

func (versionInfo VersionInfo) withIronicOverrides(images *metal3api.Images) VersionInfo {
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
	return versionInfo
}

func (versionInfo VersionInfo) WithDefaults() VersionInfo {
	if versionInfo.InstalledVersion == "" {
		versionInfo.InstalledVersion = installedVersion
	}
	if versionInfo.IronicImage == "" {
		versionInfo.IronicImage = defaultIronicImage
	}
	if versionInfo.MariaDBImage == "" {
		versionInfo.MariaDBImage = defaultMariaDBImage
	}
	if versionInfo.RamdiskDownloaderImage == "" {
		versionInfo.RamdiskDownloaderImage = defaultRamdiskDownloaderImage
	}
	if versionInfo.KeepalivedImage == "" {
		versionInfo.KeepalivedImage = defaultKeepalivedImage
	}
	return versionInfo
}
