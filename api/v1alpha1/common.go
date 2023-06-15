package v1alpha1

type IronicStatusConditionType string

const (
	// Available indicates that Ironic is fully available
	IronicStatusAvailable IronicStatusConditionType = "Available"

	// Progressing indicates that Ironic deployment is in progress
	IronicStatusProgressing IronicStatusConditionType = "Progressing"

	IronicOperatorLabel = "metal3.io/ironic-operator"
)
