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
	"slices"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
	"github.com/metal3-io/ironic-standalone-operator/pkg/ironic"
)

const (
	IronicFinalizer string = "ironic.metal3.io"
)

// IronicDatabaseReconciler reconciles a IronicDatabase object
type IronicDatabaseReconciler struct {
	client.Client
	KubeClient kubernetes.Interface
	Scheme     *runtime.Scheme
	Log        logr.Logger
}

//+kubebuilder:rbac:groups=ironic.metal3.io,resources=ironicdatabases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ironic.metal3.io,resources=ironicdatabases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ironic.metal3.io,resources=ironicdatabases/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;update;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;update;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *IronicDatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log.WithValues("IronicDatabase", req.NamespacedName)
	logger.Info("starting reconcile")

	cctx := ironic.ControllerContext{
		Context:    ctx,
		Client:     r.Client,
		KubeClient: r.KubeClient,
		Scheme:     r.Scheme,
		Logger:     logger,
	}

	db, err := getDatabase(cctx, req.NamespacedName)
	if db == nil || err != nil {
		return ctrl.Result{}, err
	}

	changed, err := r.handleDatabase(cctx, db)
	if err != nil {
		return ctrl.Result{}, err
	}
	if changed {
		return ctrl.Result{Requeue: true}, nil
	}

	logger.Info("object has been fully reconciled")
	return ctrl.Result{}, nil
}

func (r *IronicDatabaseReconciler) handleDatabase(cctx ironic.ControllerContext, db *metal3api.IronicDatabase) (requeue bool, err error) {
	if db.DeletionTimestamp.IsZero() {
		requeue, err = ensureFinalizer(cctx, db)
		if requeue || err != nil {
			return
		}
	} else {
		return r.cleanUp(cctx, db)
	}

	if db.Spec.CredentialsName == "" {
		apiSecret, err := generateSecret(cctx, db, &db.ObjectMeta, "database", false)
		if err != nil {
			return true, err
		}

		cctx.Logger.Info("updating database to use the newly generated secret", "Secret", apiSecret.Name)
		db.Spec.CredentialsName = apiSecret.Name
		err = cctx.Client.Update(cctx.Context, db)
		if err != nil {
			return true, fmt.Errorf("cannot update the new API credentials secret: %w", err)
		}

		return true, nil
	}

	status, err := ironic.EnsureDatabase(cctx, db)
	newStatus := db.Status.DeepCopy()
	if err != nil {
		cctx.Logger.Error(err, "potentially transient error, will retry")
		return
	}

	requeue = setConditionsFromStatus(cctx, status, &newStatus.Conditions, db.Generation, "database")

	if !apiequality.Semantic.DeepEqual(newStatus, &db.Status) {
		cctx.Logger.Info("updating status", "Status", newStatus)
		requeue = true
		db.Status = *newStatus
		err = cctx.Client.Status().Update(cctx.Context, db)
	}
	return
}

func (r *IronicDatabaseReconciler) cleanUp(cctx ironic.ControllerContext, db *metal3api.IronicDatabase) (bool, error) {
	if !slices.Contains(db.Finalizers, IronicFinalizer) {
		return false, nil
	}

	// Only remove the database after Ironic is removed!
	err := ironic.RemoveDatabase(cctx, db)
	if err != nil {
		return false, err
	}

	// This must be the last action.
	return removeFinalizer(cctx, db)
}

// SetupWithManager sets up the controller with the Manager.
func (r *IronicDatabaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metal3api.IronicDatabase{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
