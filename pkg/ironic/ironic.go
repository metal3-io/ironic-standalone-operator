package ironic

import (
	metal3api "github.com/metal3-io/ironic-operator/api/v1alpha1"
)

// RemoveIronic removes all bits of the Ironic deployment.
func RemoveIronic(cctx ControllerContext, ironic *metal3api.Ironic) error {
	return nil
}
