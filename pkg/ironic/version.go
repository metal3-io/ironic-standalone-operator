package ironic

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

func VersionInfoWithDefaults(versionInfo VersionInfo) VersionInfo {
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
