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
	IronicAppLabel     = IronicLabelPrefix + "/app"
	IronicServiceLabel = IronicLabelPrefix + "/ironic"
	IronicVersionLabel = IronicLabelPrefix + "/version"
)

// ResourceReference references a ConfigMap or Secret resource.
type ResourceReference struct {
	// Name of the resource.
	Name string `json:"name"`

	// Kind of the resource (ConfigMap or Secret).
	// +kubebuilder:validation:Enum=ConfigMap;Secret
	Kind string `json:"kind"`
}

// ResourceReferenceWithKey references a ConfigMap or Secret resource and
// targets a specific key from it.
type ResourceReferenceWithKey struct {
	ResourceReference `json:",inline"`

	// Key within the resource to use. If not specified and the resource contains multiple keys,
	// the first key will be used and a warning will be logged for other keys.
	// +optional
	Key string `json:"key,omitempty"`
}
