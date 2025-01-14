package v1alpha1

type IronicStatusConditionType string

const (
	// Ready indicates that Ironic is fully available.
	IronicStatusReady IronicStatusConditionType = "Ready"

	IronicReasonFailed     = "DeploymentFailed"
	IronicReasonInProgress = "DeploymentInProgress"
	IronicReasonAvailable  = "DeploymentAvailable"

	IronicLabelPrefix = "ironic.metal3.io"
)

var (
	IronicAppLabel      = IronicLabelPrefix + "/app"
	IronicDatabaseLabel = IronicLabelPrefix + "/db"
	IronicVersionLabel  = IronicLabelPrefix + "/version"
)
