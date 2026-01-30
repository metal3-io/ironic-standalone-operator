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

	"github.com/go-logr/logr"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
	"github.com/metal3-io/ironic-standalone-operator/pkg/ironic"
	"github.com/metal3-io/ironic-standalone-operator/pkg/secretutils"
)

// IronicReconciler reconciles a Ironic object.
type IronicReconciler struct {
	client.Client
	KubeClient  kubernetes.Interface
	APIReader   client.Reader
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
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;update
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups=monitoring.coreos.com,resources=servicemonitors,verbs=get;list;watch;create;update;patch;delete

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

func (r *IronicReconciler) setNotReady(cctx ironic.ControllerContext, ironicConf *metal3api.Ironic, reason, message string) error {
	setCondition(&ironicConf.Status.Conditions, ironicConf.Generation,
		false, reason, message)

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
			return requeue, err
		}
	} else {
		return r.cleanUp(cctx, ironicConf)
	}

	versionInfo, err := cctx.VersionInfo.WithIronicOverrides(ironicConf)
	if err != nil {
		// This condition requires a user's intervention
		_ = r.setNotReady(cctx, ironicConf, metal3api.IronicReasonFailed, err.Error())
		return true, err
	}
	cctx.VersionInfo = versionInfo

	actuallyRequestedVersion := cctx.VersionInfo.InstalledVersion.String()
	if actuallyRequestedVersion != ironicConf.Status.InstalledVersion && actuallyRequestedVersion != ironicConf.Status.RequestedVersion {
		// Ironic does not support downgrades when a real external database is used.
		if ironicConf.Status.InstalledVersion != "" && ironicConf.Spec.Database != nil {
			var parsedVersion metal3api.Version
			parsedVersion, err = metal3api.ParseVersion(ironicConf.Status.InstalledVersion)
			if err != nil {
				return false, err
			}

			// NOTE(dtantsur): allow upgrades from latest to a numeric version.
			// This requires a user to ensure that the new version is actually newer, but it's a valid scenario.
			if !parsedVersion.IsLatest() && cctx.VersionInfo.InstalledVersion.Compare(parsedVersion) < 0 {
				cctx.Logger.Info("refusing to downgrade Ironic", "InstalledVersion", ironicConf.Status.InstalledVersion, "RequestedVersion", actuallyRequestedVersion)
				_ = r.setNotReady(cctx, ironicConf, metal3api.IronicReasonFailed, "Ironic does not support downgrades with an external database")
				return false, nil
			}
		}
		cctx.Logger.Info("new version requested", "InstalledVersion", ironicConf.Status.InstalledVersion, "RequestedVersion", actuallyRequestedVersion)
		ironicConf.Status.RequestedVersion = actuallyRequestedVersion
		err = r.setNotReady(cctx, ironicConf, metal3api.IronicReasonInProgress, "new version requested")
		if err != nil {
			return requeue, err
		}
	}

	apiSecret, requeue, err := r.ensureAPISecret(cctx, ironicConf)
	if requeue || err != nil {
		return requeue, err
	}

	var tlsSecret *corev1.Secret
	if tlsSecretName := ironicConf.Spec.TLS.CertificateName; tlsSecretName != "" {
		tlsSecret, requeue, err = r.getAndUpdateSecret(cctx, ironicConf, tlsSecretName)
		if requeue || err != nil {
			return requeue, err
		}
	}

	var bmcCASecret *corev1.Secret
	var bmcCAConfigMap *corev1.ConfigMap
	if bmcCARef := ironic.GetBMCCA(&ironicConf.Spec.TLS); bmcCARef != nil {
		switch bmcCARef.Kind {
		case metal3api.ResourceKindSecret:
			bmcCASecret, requeue, err = r.getAndUpdateSecret(cctx, ironicConf, bmcCARef.Name)
		case metal3api.ResourceKindConfigMap:
			bmcCAConfigMap, requeue, err = r.getConfigMap(cctx, ironicConf, bmcCARef.Name)
		}
		if requeue || err != nil {
			return requeue, err
		}
	}

	var trustedCASecret *corev1.Secret
	var trustedCAConfigMap *corev1.ConfigMap
	if trustedCARef := ironic.GetTrustedCA(&ironicConf.Spec.TLS); trustedCARef != nil {
		switch trustedCARef.Kind {
		case metal3api.ResourceKindSecret:
			trustedCASecret, requeue, err = r.getAndUpdateSecret(cctx, ironicConf, trustedCARef.Name)
		case metal3api.ResourceKindConfigMap:
			trustedCAConfigMap, requeue, err = r.getConfigMap(cctx, ironicConf, trustedCARef.Name)
		}
		if requeue || err != nil {
			return requeue, err
		}
	}

	resources := ironic.Resources{
		Ironic:             ironicConf,
		APISecret:          apiSecret,
		TLSSecret:          tlsSecret,
		BMCCASecret:        bmcCASecret,
		BMCCAConfigMap:     bmcCAConfigMap,
		TrustedCASecret:    trustedCASecret,
		TrustedCAConfigMap: trustedCAConfigMap,
	}

	status, err := ironic.EnsureIronic(cctx, resources)
	if err != nil {
		cctx.Logger.Error(err, "potentially transient error, will retry")
		return requeue, err
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
	return requeue, err
}

// Get a secret and update its owner references using SecretManager.
// This ensures the secret is labeled for cache filtering and has owner references set.
// Only returns a valid pointer if requeue is false and err is nil.
func (r *IronicReconciler) getAndUpdateSecret(cctx ironic.ControllerContext, ironicConf *metal3api.Ironic, secretName string) (secret *corev1.Secret, requeue bool, err error) {
	namespacedName := types.NamespacedName{
		Namespace: ironicConf.Namespace,
		Name:      secretName,
	}

	secretManager := secretutils.NewSecretManager(cctx.Context, cctx.Logger, cctx.Client, r.APIReader)
	secret, err = secretManager.AcquireSecret(namespacedName, ironicConf, cctx.Scheme)
	if err != nil {
		// NotFound requires a user's intervention, so reporting it in the conditions.
		// Everything else is reported up for a retry.
		if k8serrors.IsNotFound(err) {
			message := fmt.Sprintf("secret %s/%s not found", ironicConf.Namespace, secretName)
			_ = r.setNotReady(cctx, ironicConf, metal3api.IronicReasonFailed, message)
		}
		return nil, true, fmt.Errorf("cannot load secret %s/%s: %w", ironicConf.Namespace, secretName, err)
	}

	return secret, false, nil
}

// Get a user-provided configmap.
// Only returns a valid pointer if requeue is false and err is nil.
func (r *IronicReconciler) getConfigMap(cctx ironic.ControllerContext, ironicConf *metal3api.Ironic, configMapName string) (configMap *corev1.ConfigMap, requeue bool, err error) {
	namespacedName := types.NamespacedName{
		Namespace: ironicConf.Namespace,
		Name:      configMapName,
	}

	secretManager := secretutils.NewSecretManager(cctx.Context, cctx.Logger, cctx.Client, r.APIReader)
	configMap, err = secretManager.ObtainConfigMap(namespacedName)
	if err != nil {
		// NotFound requires a user's intervention, so reporting it in the conditions.
		// Everything else is reported up for a retry.
		if k8serrors.IsNotFound(err) {
			message := fmt.Sprintf("configmap %s/%s not found", ironicConf.Namespace, configMapName)
			_ = r.setNotReady(cctx, ironicConf, metal3api.IronicReasonFailed, message)
		}
		return nil, true, fmt.Errorf("cannot load configmap %s/%s: %w", ironicConf.Namespace, configMapName, err)
	}

	return configMap, false, nil
}

func (r *IronicReconciler) ensureAPISecret(cctx ironic.ControllerContext, ironicConf *metal3api.Ironic) (apiSecret *corev1.Secret, requeue bool, err error) {
	if ironicConf.Spec.APICredentialsName == "" {
		apiSecret, err = generateSecret(cctx, ironicConf, &ironicConf.ObjectMeta, "service", true)
		if err != nil {
			_ = r.setNotReady(cctx, ironicConf, metal3api.IronicReasonFailed, err.Error())
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
		return apiSecret, requeue, err
	}

	apiSecret, requeue, err = r.getAndUpdateSecret(cctx, ironicConf, ironicConf.Spec.APICredentialsName)
	if requeue || err != nil {
		return nil, requeue, err
	}

	requeue, err = ironic.UpdateSecret(apiSecret, cctx.Logger)
	if err != nil {
		_ = r.setNotReady(cctx, ironicConf, metal3api.IronicReasonFailed, err.Error())
		return nil, true, err
	}
	if requeue {
		cctx.Logger.Info("updating htpasswd", "Secret", apiSecret.Name)
		err = cctx.Client.Update(cctx.Context, apiSecret)
	}

	return apiSecret, requeue, err
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
	builder := ctrl.NewControllerManagedBy(mgr).
		For(&metal3api.Ironic{}).
		Owns(&corev1.Secret{}, builder.MatchEveryOwner).
		Owns(&corev1.Service{}).
		Owns(&appsv1.DaemonSet{}).
		Owns(&appsv1.Deployment{}).
		Owns(&batchv1.Job{})

	hasServiceMonitor, err := clusterHasCRD(mgr, &monitoringv1.ServiceMonitor{})
	if err != nil {
		return err
	}

	if hasServiceMonitor {
		builder = builder.Owns(&monitoringv1.ServiceMonitor{})
	} else {
		r.Log.Info("WARNING: ServiceMonitor resources are not available and will not be reconciled")
	}

	return builder.Complete(r)
}
