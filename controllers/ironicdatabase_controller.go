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

package controllers

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	metal3api "github.com/metal3-io/ironic-operator/api/v1alpha1"
	"github.com/metal3-io/ironic-operator/pkg/ironic"
)

// IronicDatabaseReconciler reconciles a IronicDatabase object
type IronicDatabaseReconciler struct {
	client.Client
	KubeClient kubernetes.Interface
	Scheme     *runtime.Scheme
}

//+kubebuilder:rbac:groups=metal3.io,resources=ironicdatabases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=metal3.io,resources=ironicdatabases/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=metal3.io,resources=ironicdatabases/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;update;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;update;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *IronicDatabaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("ironicDatabase", req.NamespacedName)
	logger.Info("starting reconcile")

	db, err := r.getDatabase(ctx, req)
	if db == nil || err != nil {
		return ctrl.Result{}, err
	}

	var changed bool
	if db.DeletionTimestamp.IsZero() {
		changed, err = ensureFinalizer(ctx, r.Client, db)
	} else {
		changed, err = r.cleanUp(ctx, r.KubeClient, db)
	}

	if err != nil {
		return ctrl.Result{}, err
	}
	if changed {
		return ctrl.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

func (r *IronicDatabaseReconciler) getDatabase(ctx context.Context, req ctrl.Request) (*metal3api.IronicDatabase, error) {
	db := &metal3api.IronicDatabase{}
	err := r.Get(ctx, req.NamespacedName, db)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "could not load ironic configuration %s", req.NamespacedName)
	}

	return db, nil
}

func (r *IronicDatabaseReconciler) cleanUp(ctx context.Context, kubeClient kubernetes.Interface, db *metal3api.IronicDatabase) (bool, error) {
	if !slices.Contains(db.Finalizers, metal3api.IronicFinalizer) {
		return false, nil
	}

	// Only remove the database after Ironic is removed!
	err := ironic.RemoveDatabase(ctx, kubeClient, db)
	if err != nil {
		return false, err
	}

	// This must be the last action.
	return removeFinalizer(ctx, r.Client, db)
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
