package controllers

import (
	"github.com/pkg/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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
	meta.SetStatusCondition(conditions, cond)
}

func getIronic(cctx ironic.ControllerContext, name types.NamespacedName) (*metal3api.Ironic, error) {
	ironicConf := &metal3api.Ironic{}
	err := cctx.Client.Get(cctx.Context, name, ironicConf)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "could not load ironic configuration %s", name)
	}

	return ironicConf, nil
}

func getDatabase(cctx ironic.ControllerContext, name types.NamespacedName) (*metal3api.IronicDatabase, error) {
	db := &metal3api.IronicDatabase{}
	err := cctx.Client.Get(cctx.Context, name, db)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}
		cctx.Logger.Error(err, "unexpected error when loading the database")
		return nil, errors.Wrapf(err, "could not load ironic configuration %s", name)
	}

	return db, nil
}
