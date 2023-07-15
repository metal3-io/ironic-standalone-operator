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
	"time"

	"golang.org/x/exp/slices"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	metal3api "github.com/metal3-io/ironic-operator/api/v1alpha1"
	"github.com/metal3-io/ironic-operator/pkg/ironic"
)

// IronicReconciler reconciles a Ironic object
type IronicReconciler struct {
	client.Client
	KubeClient kubernetes.Interface
	Scheme     *runtime.Scheme
}

//+kubebuilder:rbac:groups=metal3.io,resources=ironics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=metal3.io,resources=ironics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=metal3.io,resources=ironics/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;update;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;update;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *IronicReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("Ironic", req.NamespacedName)
	logger.Info("starting reconcile")

	cctx := ironic.ControllerContext{
		Context:    ctx,
		Client:     r.Client,
		KubeClient: r.KubeClient,
		Scheme:     r.Scheme,
		Logger:     logger,
	}

	ironicConf, err := getIronic(cctx, req.NamespacedName)
	if ironicConf == nil || err != nil {
		return ctrl.Result{}, err
	}

	changed, err := r.handleIronic(cctx, ironicConf)
	if err != nil {
		return ctrl.Result{RequeueAfter: 15 * time.Second}, err
	}
	if changed {
		return ctrl.Result{Requeue: true}, nil
	}

	logger.Info("object has been fully reconciled")
	return ctrl.Result{}, nil
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

	var db *metal3api.IronicDatabase
	if ironicConf.Spec.DatabaseName != "" {
		dbName := types.NamespacedName{
			Namespace: ironicConf.Namespace,
			Name:      ironicConf.Spec.DatabaseName,
		}
		db = &metal3api.IronicDatabase{}
		err := cctx.Client.Get(cctx.Context, dbName, db)
		if err != nil {
			cctx.Logger.Error(err, "cannot load linked database")
			return true, err
		}
	}

	status, endpoints, err := ironic.EnsureIronic(cctx, ironicConf, db)
	newStatus := ironicConf.Status.DeepCopy()
	if err != nil {
		cctx.Logger.Error(err, "failed to create or update ironic")
		setCondition(cctx, &newStatus.Conditions, ironicConf.Generation, metal3api.IronicStatusAvailable, false, "DeploymentFailed", err.Error())
	} else if status != metal3api.IronicStatusAvailable {
		setCondition(cctx, &newStatus.Conditions, ironicConf.Generation, metal3api.IronicStatusAvailable, false, "DeploymentInProgress", "deployment is not ready yet")
		setCondition(cctx, &newStatus.Conditions, ironicConf.Generation, metal3api.IronicStatusProgressing, true, "DeploymentInProgress", "deployment is in progress")
	} else {
		setCondition(cctx, &newStatus.Conditions, ironicConf.Generation, metal3api.IronicStatusAvailable, true, "DeploymentAvailable", "ironic is available")
		setCondition(cctx, &newStatus.Conditions, ironicConf.Generation, metal3api.IronicStatusProgressing, false, "DeploymentAvailable", "ironic is available")
		newStatus.IronicEndpoints = endpoints
	}

	if !apiequality.Semantic.DeepEqual(newStatus, &ironicConf.Status) {
		cctx.Logger.Info("updating status", "Status", newStatus)
		ironicConf.Status = *newStatus
		err = cctx.Client.Status().Update(cctx.Context, ironicConf)
	}
	return
}

func (r *IronicReconciler) cleanUp(cctx ironic.ControllerContext, ironicConf *metal3api.Ironic) (bool, error) {
	if !slices.Contains(ironicConf.Finalizers, metal3api.IronicFinalizer) {
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
		Owns(&corev1.Secret{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.DaemonSet{}).
		Owns(&appsv1.Deployment{}).
		Owns(&metal3api.IronicDatabase{}).
		Complete(r)
}
