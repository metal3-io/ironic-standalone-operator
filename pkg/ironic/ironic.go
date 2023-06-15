package ironic

import (
	metal3api "github.com/metal3-io/ironic-operator/api/v1alpha1"
)

// EnsureIronic removes all bits of the Ironic deployment.
func EnsureIronic(cctx ControllerContext, ironic *metal3api.Ironic, db *metal3api.IronicDatabase) (status metal3api.IronicStatusConditionType, endpoints []string, err error) {
	return metal3api.IronicStatusProgressing, nil, nil
}

// RemoveIronic removes all bits of the Ironic deployment.
func RemoveIronic(cctx ControllerContext, ironic *metal3api.Ironic) error {
	return nil
}
