package controllers

import (
	"context"

	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	metal3api "github.com/metal3-io/ironic-operator/api/v1alpha1"
)

func ensureFinalizer(ctx context.Context, cli client.Client, obj client.Object) (bool, error) {
	changed := controllerutil.AddFinalizer(obj, metal3api.IronicFinalizer)
	if changed {
		err := cli.Update(ctx, obj)
		if err != nil {
			return false, errors.Wrap(err, "failed to add finalizer")
		}
		return true, nil
	}

	return false, nil
}

func removeFinalizer(ctx context.Context, cli client.Client, obj client.Object) (bool, error) {
	changed := controllerutil.RemoveFinalizer(obj, metal3api.IronicFinalizer)
	if changed {
		err := cli.Update(ctx, obj)
		if err != nil {
			return false, errors.Wrap(err, "failed to remove finalizer")
		}
		return true, nil
	}

	return false, nil
}
