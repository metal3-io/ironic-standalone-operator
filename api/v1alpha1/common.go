package v1alpha1

type IronicStatusConditionType string

const (
	// Ready indicates that Ironic is fully available.
	IronicStatusReady IronicStatusConditionType = "Ready"

	IronicReasonFailed     = "DeploymentFailed"
	IronicReasonInProgress = "DeploymentInProgress"
	IronicReasonAvailable  = "DeploymentAvailable"

	IronicLabelPrefix = "ironic.metal3.io"

	// LabelEnvironmentName is the label key that must be present on user-provided
	// Secrets and ConfigMaps to grant the operator permission to use them.
	// The value must be LabelEnvironmentValue.
	LabelEnvironmentName = "environment.metal3.io/ironic-standalone-operator"

	// LabelEnvironmentValue is the required value for the LabelEnvironmentName label.
	LabelEnvironmentValue = "true"
)

var (
	IronicAppLabel     = IronicLabelPrefix + "/app"
	IronicServiceLabel = IronicLabelPrefix + "/ironic"
	IronicVersionLabel = IronicLabelPrefix + "/version"
)
