package controllers

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	metal3api "github.com/metal3-io/ironic-operator/api/v1alpha1"
	"github.com/metal3-io/ironic-operator/pkg/ironic"
)

func ensureFinalizer(cctx ironic.ControllerContext, obj client.Object) (bool, error) {
	changed := controllerutil.AddFinalizer(obj, metal3api.IronicFinalizer)
	if changed {
		err := cctx.Client.Update(cctx.Context, obj)
		if err != nil {
			return false, errors.Wrap(err, "failed to add finalizer")
		}
		return true, nil
	}

	return false, nil
}

func removeFinalizer(cctx ironic.ControllerContext, obj client.Object) (bool, error) {
	changed := controllerutil.RemoveFinalizer(obj, metal3api.IronicFinalizer)
	if changed {
		err := cctx.Client.Update(cctx.Context, obj)
		if err != nil {
			return false, errors.Wrap(err, "failed to remove finalizer")
		}
		return true, nil
	}

	return false, nil
}

func setCondition(cctx ironic.ControllerContext, conditions *[]metav1.Condition, generation int64, status metal3api.IronicStatusConditionType, value bool, reason, message string) {
	condStatus := metav1.ConditionFalse
	if value {
		condStatus = metav1.ConditionTrue
	}
	cond := metav1.Condition{
		Type:               string(status),
		Status:             condStatus,
		ObservedGeneration: generation,
		Reason:             reason,
		Message:            message,
	}
	cctx.Logger.Info("recording condition change", "Condition", cond)
	meta.SetStatusCondition(conditions, cond)
}
