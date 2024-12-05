package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
)

type IronicStatusConditionType string

const (
	// Ready indicates that Ironic is fully available.
	IronicStatusReady IronicStatusConditionType = "Ready"

	IronicReasonFailed     = "DeploymentFailed"
	IronicReasonInProgress = "DeploymentInProgress"
	IronicReasonAvailable  = "DeploymentAvailable"

	IronicOperatorLabel = "metal3.io/ironic-standalone-operator"
)

type Overrides struct {
	// Extra annotations to add to each container.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Extra containers to add to the deployment or daemon set.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	Containers []corev1.Container `json:"containers,omitempty"`

	// Extra environment variables to add to each container.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	Env []corev1.EnvVar `json:"env,omitempty"`

	// Extra environment variables (with sources) to add to each container.
	// +optional
	// +listType=atomic
	EnvFrom []corev1.EnvFromSource `json:"envFrom,omitempty"`

	// Extra init containers to add to the deployment or daemon set.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	InitContainers []corev1.Container `json:"initContainers,omitempty"`

	// Extra labels to add to each container.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}
