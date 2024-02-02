package controllers

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
	"github.com/metal3-io/ironic-standalone-operator/pkg/ironic"
)

func ensureFinalizer(cctx ironic.ControllerContext, obj client.Object) (bool, error) {
	changed := controllerutil.AddFinalizer(obj, IronicFinalizer)
	if changed {
		cctx.Logger.Info(fmt.Sprintf("adding finalizer %s", IronicFinalizer))
		err := cctx.Client.Update(cctx.Context, obj)
		if err != nil {
			return false, fmt.Errorf("failed to add finalizer: %w", err)
		}
		return true, nil
	}

	return false, nil
}

func removeFinalizer(cctx ironic.ControllerContext, obj client.Object) (bool, error) {
	changed := controllerutil.RemoveFinalizer(obj, IronicFinalizer)
	if changed {
		cctx.Logger.Info(fmt.Sprintf("removing finalizer %s", IronicFinalizer))
		err := cctx.Client.Update(cctx.Context, obj)
		if err != nil {
			return false, fmt.Errorf("failed to remove finalizer: %w", err)
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
		return nil, fmt.Errorf("could not load ironic configuration %s: %w", name, err)
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
		return nil, fmt.Errorf("could not load database configuration %s: %w", name, err)
	}

	return db, nil
}

func generateSecret(cctx ironic.ControllerContext, owner metav1.Object, meta *metav1.ObjectMeta, name string, extraFields bool) (secret *corev1.Secret, err error) {
	secret, err = ironic.GenerateSecret(meta, name, extraFields)
	if err != nil {
		return
	}

	err = controllerutil.SetControllerReference(owner, secret, cctx.Scheme)
	if err != nil {
		return
	}

	err = cctx.Client.Create(cctx.Context, secret)
	if err != nil {
		err = fmt.Errorf("cannot create an API credentials secret: %w", err)
		return
	}

	return
}
