/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"reflect"
	"slices"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/go-logr/logr"
	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
	"github.com/metal3-io/ironic-standalone-operator/pkg/ironic"
)

// IronicReconciler reconciles a Ironic object.
type IronicReconciler struct {
	client.Client
	KubeClient  kubernetes.Interface
	Scheme      *runtime.Scheme
	Log         logr.Logger
	Domain      string
	VersionInfo ironic.VersionInfo
}

//+kubebuilder:rbac:groups=ironic.metal3.io,resources=ironics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ironic.metal3.io,resources=ironics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ironic.metal3.io,resources=ironics/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *IronicReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("Ironic", req.NamespacedName)
	logger.Info("starting reconcile")

	cctx := ironic.ControllerContext{
		Context:     ctx,
		Client:      r.Client,
		KubeClient:  r.KubeClient,
		Scheme:      r.Scheme,
		Logger:      logger,
		Domain:      r.Domain,
		VersionInfo: r.VersionInfo,
	}

	ironicConf, err := getIronic(cctx, req.NamespacedName)
	if ironicConf == nil || err != nil {
		return ctrl.Result{}, err
	}

	changed, err := r.handleIronic(cctx, ironicConf)
	if err != nil {
		cctx.Logger.Error(err, "reconcile failed, will retry")
		return ctrl.Result{}, err
	}
	if changed {
		return ctrl.Result{Requeue: true}, nil
	}

	logger.Info("object has been fully reconciled")
	return ctrl.Result{}, nil
}

func (r *IronicReconciler) setCondition(cctx ironic.ControllerContext, ironicConf *metal3api.Ironic, value bool, reason, message string) error {
	setCondition(cctx, &ironicConf.Status.Conditions, ironicConf.Generation,
		metal3api.IronicStatusReady, value, reason, message)

	err := cctx.Client.Status().Update(cctx.Context, ironicConf)
	if err != nil {
		cctx.Logger.Error(err, "potentially transient error when updating conditions")
	}
	return err
}

func (r *IronicReconciler) handleIronic(cctx ironic.ControllerContext, ironicConf *metal3api.Ironic) (requeue bool, err error) {
	if ironicConf.DeletionTimestamp.IsZero() {
		requeue, err = ensureFinalizer(cctx, ironicConf)
		if requeue || err != nil {
			return
		}
	} else {
		return r.cleanUp(cctx, ironicConf)
	}

	versionInfo, err := cctx.VersionInfo.WithIronicOverrides(ironicConf)
	if err != nil {
		// This condition requires a user's intervention
		_ = r.setCondition(cctx, ironicConf, false, metal3api.IronicReasonFailed, err.Error())
		return true, err
	}
	cctx.VersionInfo = versionInfo

	dbConf := ironicConf.Spec.Database

	actuallyRequestedVersion := cctx.VersionInfo.InstalledVersion.String()
	if actuallyRequestedVersion != ironicConf.Status.InstalledVersion && actuallyRequestedVersion != ironicConf.Status.RequestedVersion {
		// Ironic does not support downgrades when a real external database is used.
		if ironicConf.Status.InstalledVersion != "" && dbConf != nil {
			parsedVersion, err := metal3api.ParseVersion(ironicConf.Status.InstalledVersion)
			if err != nil {
				return false, err
			}

			// NOTE(dtantsur): allow upgrades from latest to a numeric version.
			// This requires a user to ensure that the new version is actually newer, but it's a valid scenario.
			if !parsedVersion.IsLatest() && cctx.VersionInfo.InstalledVersion.Compare(parsedVersion) < 0 {
				cctx.Logger.Info("refusing to downgrade Ironic", "InstalledVersion", ironicConf.Status.InstalledVersion, "RequestedVersion", actuallyRequestedVersion)
				_ = r.setCondition(cctx, ironicConf, false, metal3api.IronicReasonFailed, "Ironic does not support downgrades with an external database")
				return false, nil
			}
		}
		cctx.Logger.Info("new version requested", "InstalledVersion", ironicConf.Status.InstalledVersion, "RequestedVersion", actuallyRequestedVersion)
		ironicConf.Status.RequestedVersion = actuallyRequestedVersion
		err = r.setCondition(cctx, ironicConf, false, metal3api.IronicReasonInProgress, "new version requested")
		if err != nil {
			return
		}
	}

	apiSecret, requeue, err := r.ensureAPISecret(cctx, ironicConf)
	if requeue || err != nil {
		return
	}

	db, requeue, err := r.ensureDatabase(cctx, ironicConf)
	if requeue || err != nil {
		return
	}

	// Compatibility logic, remove when IronicDatabase is removed
	if db != nil {
		if !meta.IsStatusConditionTrue(db.Status.Conditions, string(metal3api.IronicStatusReady)) {
			setCondition(cctx, &ironicConf.Status.Conditions, ironicConf.Generation, metal3api.IronicStatusReady,
				false, metal3api.IronicReasonInProgress, "database is not ready yet")
			err = cctx.Client.Status().Update(cctx.Context, ironicConf)
			return
		}
		dbConf = &metal3api.Database{
			CredentialsName:    db.Spec.CredentialsName,
			Host:               ironic.DatabaseDNSName(db, cctx.Domain),
			Name:               ironic.DatabaseName,
			TLSCertificateName: db.Spec.TLSCertificateName,
		}
	}

	var tlsSecret *corev1.Secret
	if tlsSecretName := ironicConf.Spec.TLS.CertificateName; tlsSecretName != "" {
		tlsSecret, requeue, err = r.getAndUpdateSecret(cctx, ironicConf, tlsSecretName)
		if requeue || err != nil {
			return
		}
	}

	resources := ironic.Resources{
		Ironic:    ironicConf,
		Database:  dbConf,
		APISecret: apiSecret,
		TLSSecret: tlsSecret,
	}

	status, err := ironic.EnsureIronic(cctx, resources)
	if err != nil {
		cctx.Logger.Error(err, "potentially transient error, will retry")
		return
	}

	newStatus := ironicConf.Status.DeepCopy()
	setConditionsFromStatus(cctx, status, &newStatus.Conditions, ironicConf.Generation, "ironic")
	requeue = status.NeedsRequeue()
	if status.IsReady() {
		newStatus.InstalledVersion = actuallyRequestedVersion
	}

	if !apiequality.Semantic.DeepEqual(newStatus, &ironicConf.Status) {
		cctx.Logger.Info("updating status", "Status", newStatus)
		ironicConf.Status = *newStatus
		err = cctx.Client.Status().Update(cctx.Context, ironicConf)
	}
	return
}

// Get a secret and update its owner references.
// Only returns a valid pointer if requeue is false and err is nil.
func (r *IronicReconciler) getAndUpdateSecret(cctx ironic.ControllerContext, ironicConf *metal3api.Ironic, secretName string) (secret *corev1.Secret, requeue bool, err error) {
	namespacedName := types.NamespacedName{
		Namespace: ironicConf.Namespace,
		Name:      secretName,
	}
	secret = &corev1.Secret{}
	err = cctx.Client.Get(cctx.Context, namespacedName, secret)
	if err != nil {
		// NotFound requires a user's intervention, so reporting it in the conditions.
		// Everything else is reported up for a retry.
		if k8serrors.IsNotFound(err) {
			message := fmt.Sprintf("secret %s/%s not found", ironicConf.Namespace, secretName)
			_ = r.setCondition(cctx, ironicConf, false, metal3api.IronicReasonFailed, message)
		}
		return nil, true, fmt.Errorf("cannot load secret %s/%s: %w", ironicConf.Namespace, secretName, err)
	}

	requeue, err = updateSecretOwners(cctx, ironicConf, secret)
	if requeue || err != nil {
		return nil, requeue, err
	}

	return secret, false, nil
}

func (r *IronicReconciler) ensureAPISecret(cctx ironic.ControllerContext, ironicConf *metal3api.Ironic) (apiSecret *corev1.Secret, requeue bool, err error) {
	if ironicConf.Spec.APICredentialsName == "" {
		apiSecret, err = generateSecret(cctx, ironicConf, &ironicConf.ObjectMeta, "service", true)
		if err != nil {
			_ = r.setCondition(cctx, ironicConf, false, metal3api.IronicReasonFailed, err.Error())
			return nil, true, err
		}

		cctx.Logger.Info("updating Ironic to use the newly generated secret", "Secret", apiSecret.Name)
		ironicConf.Spec.APICredentialsName = apiSecret.Name
		err = cctx.Client.Update(cctx.Context, ironicConf)
		if err != nil {
			// Considering this a transient error
			return nil, true, fmt.Errorf("cannot update the new API credentials secret: %w", err)
		}

		requeue = true
		return
	}

	apiSecret, requeue, err = r.getAndUpdateSecret(cctx, ironicConf, ironicConf.Spec.APICredentialsName)
	if requeue || err != nil {
		return nil, requeue, err
	}

	requeue, err = ironic.UpdateSecret(apiSecret, cctx.Logger)
	if err != nil {
		_ = r.setCondition(cctx, ironicConf, false, metal3api.IronicReasonFailed, err.Error())
		return nil, true, err
	}
	if requeue {
		cctx.Logger.Info("updating htpasswd", "Secret", apiSecret.Name)
		err = cctx.Client.Update(cctx.Context, apiSecret)
	}

	return
}

func (r *IronicReconciler) ensureDatabase(cctx ironic.ControllerContext, ironicConf *metal3api.Ironic) (db *metal3api.IronicDatabase, requeue bool, err error) {
	if ironicConf.Spec.DatabaseName == "" {
		return
	}

	dbName := types.NamespacedName{
		Namespace: ironicConf.Namespace,
		Name:      ironicConf.Spec.DatabaseName,
	}
	db = &metal3api.IronicDatabase{}
	err = cctx.Client.Get(cctx.Context, dbName, db)
	if err != nil {
		// NotFound requires a user's intervention, so reporting it in the conditions.
		// Everything else is reported up for a retry.
		if k8serrors.IsNotFound(err) {
			message := fmt.Sprintf("database %s/%s not found", ironicConf.Namespace, ironicConf.Spec.DatabaseName)
			_ = r.setCondition(cctx, ironicConf, false, metal3api.IronicReasonFailed, message)
		}
		return nil, true, fmt.Errorf("cannot load linked database %s/%s: %w", ironicConf.Namespace, ironicConf.Spec.DatabaseName, err)
	}

	oldReferences := db.GetOwnerReferences()
	err = controllerutil.SetControllerReference(ironicConf, db, cctx.Scheme)
	if err != nil {
		return nil, true, err
	}
	if !reflect.DeepEqual(oldReferences, db.GetOwnerReferences()) {
		cctx.Logger.Info("updating owner reference", "Database", db.Name)
		err = cctx.Client.Update(cctx.Context, db)
		requeue = true
	}

	return
}

func (r *IronicReconciler) cleanUp(cctx ironic.ControllerContext, ironicConf *metal3api.Ironic) (bool, error) {
	if !slices.Contains(ironicConf.Finalizers, IronicFinalizer) {
		return false, nil
	}

	err := ironic.RemoveIronic(cctx, ironicConf)
	if err != nil {
		return false, err
	}

	// This must be the last action.
	return removeFinalizer(cctx, ironicConf)
}

// SetupWithManager sets up the controller with the Manager.
func (r *IronicReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metal3api.Ironic{}).
		Owns(&corev1.Secret{}, builder.MatchEveryOwner).
		Owns(&corev1.Service{}).
		Owns(&appsv1.DaemonSet{}).
		Owns(&appsv1.Deployment{}).
		Owns(&batchv1.Job{}).
		Owns(&metal3api.IronicDatabase{}).
		Complete(r)
}
