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
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;update
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

	var bmcSecret *corev1.Secret
	if bmcSecretName := ironicConf.Spec.TLS.BMCCAName; bmcSecretName != "" {
		bmcSecret, requeue, err = r.getAndUpdateSecret(cctx, ironicConf, bmcSecretName)
		if requeue || err != nil {
			return requeue, err
		}
	}

	// Ensure switch config secret if networking service is enabled
	var switchConfigSecret *corev1.Secret
	if ironicConf.IsNetworkingServiceEnabled() {
		err = ironic.EnsureSwitchConfigSecret(cctx, ironicConf)
		if err != nil {
			cctx.Logger.Error(err, "failed to ensure switch config secret")
			return requeue, err
		}
		// Fetch the switch config secret so we can include its version in pod annotations
		switchConfigSecretName := ironic.SwitchConfigSecretName(ironicConf)
		switchConfigSecret, requeue, err = r.getAndUpdateSecret(cctx, ironicConf, switchConfigSecretName)
		if requeue || err != nil {
			return requeue, err
		}
	} else {
		err = ironic.EnsureSwitchConfigSecretDeleted(cctx, ironicConf)
		if err != nil {
			cctx.Logger.Error(err, "failed to ensure switch config secret deleted")
			return requeue, err
		}
	}

	// Fetch switch credentials secret if specified
	var switchCredentialsSecret *corev1.Secret
	if ironicConf.IsNetworkingServiceEnabled() &&
		ironicConf.Spec.NetworkingService.SwitchCredentialsSecretName != "" {
		switchCredentialsSecretName := ironicConf.Spec.NetworkingService.SwitchCredentialsSecretName
		switchCredentialsSecret, requeue, err = r.getAndUpdateSecret(cctx, ironicConf, switchCredentialsSecretName)
		if requeue || err != nil {
			return requeue, err
		}
	}

	var trustedCAConfigMap *corev1.ConfigMap
	if trustedCAConfigMapName := ironicConf.Spec.TLS.TrustedCAName; trustedCAConfigMapName != "" {
		trustedCAConfigMap, requeue, err = r.getConfigMap(cctx, ironicConf, trustedCAConfigMapName)
		if requeue || err != nil {
			return requeue, err
		}
	}

	resources := ironic.Resources{
		Ironic:                    ironicConf,
		APISecret:                 apiSecret,
		TLSSecret:                 tlsSecret,
		BMCCASecret:               bmcSecret,
		SwitchConfigSecret:        switchConfigSecret,
		SwitchCredentialsSecret:   switchCredentialsSecret,
		TrustedCAConfigMap:        trustedCAConfigMap,
		NetworkingServiceEndpoint: ironic.GetNetworkingServiceEndpoint(ironicConf),
	}

	// Manage networking service deployment
	if ironicConf.IsNetworkingServiceEnabled() {
		if ironicConf.Spec.NetworkingService.Endpoint == "" {
			// Operator-managed networking service deployment
			err = r.ensureNetworkingDeployment(cctx, resources)
			if err != nil {
				cctx.Logger.Error(err, "failed to ensure networking deployment")
				return requeue, err
			}
			err = r.ensureNetworkingService(cctx, ironicConf)
			if err != nil {
				cctx.Logger.Error(err, "failed to ensure networking service")
				return requeue, err
			}

			// Wait for the networking deployment to be ready before proceeding
			// so that Ironic doesn't start before its networking service is available.
			networkingDeploy := &appsv1.Deployment{}
			err = cctx.Client.Get(cctx.Context, client.ObjectKey{
				Name:      ironic.NetworkingDeploymentName(ironicConf),
				Namespace: ironicConf.Namespace,
			}, networkingDeploy)
			if err != nil {
				return requeue, err
			}
			if networkingDeploy.Status.ReadyReplicas < 1 {
				cctx.Logger.Info("waiting for networking service deployment to become ready")
				return true, nil
			}
		} else {
			// External networking service - use the provided endpoint
			// Clean up networking deployment/service if using an external service
			err = r.deleteNetworkingResources(cctx, ironicConf)
			if err != nil {
				cctx.Logger.Error(err, "failed to delete networking resources")
				return requeue, err
			}
		}
	} else {
		// Clean up networking deployment/service if disabled
		err = r.deleteNetworkingResources(cctx, ironicConf)
		if err != nil {
			cctx.Logger.Error(err, "failed to delete networking resources")
			return requeue, err
		}
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

	if err := r.deleteNetworkingResources(cctx, ironicConf); err != nil {
		return false, err
	}

	if err := ironic.EnsureSwitchConfigSecretDeleted(cctx, ironicConf); err != nil {
		return false, err
	}

	// This must be the last action.
	return removeFinalizer(cctx, ironicConf)
}

// ensureNetworkingDeployment creates or updates the networking service deployment.
func (r *IronicReconciler) ensureNetworkingDeployment(cctx ironic.ControllerContext, resources ironic.Resources) error {
	deployment := ironic.BuildNetworkingDeployment(cctx, resources)

	// Set owner reference so the deployment is cleaned up when Ironic is deleted
	if err := ctrl.SetControllerReference(resources.Ironic, deployment, r.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference on networking deployment: %w", err)
	}

	// Create or update the deployment
	existingDeployment := &appsv1.Deployment{}
	err := cctx.Client.Get(cctx.Context, client.ObjectKey{
		Name:      deployment.Name,
		Namespace: deployment.Namespace,
	}, existingDeployment)

	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Create new deployment
			cctx.Logger.Info("creating networking service deployment", "Deployment", deployment.Name)
			return cctx.Client.Create(cctx.Context, deployment)
		}
		return fmt.Errorf("failed to get networking deployment: %w", err)
	}

	// Update existing deployment
	existingDeployment.Spec = deployment.Spec
	existingDeployment.Labels = deployment.Labels
	cctx.Logger.Info("updating networking service deployment", "Deployment", deployment.Name)
	return cctx.Client.Update(cctx.Context, existingDeployment)
}

// ensureNetworkingService creates or updates the networking service.
func (r *IronicReconciler) ensureNetworkingService(cctx ironic.ControllerContext, ironicConf *metal3api.Ironic) error {
	service := ironic.BuildNetworkingService(ironicConf)

	// Set owner reference so the service is cleaned up when Ironic is deleted
	if err := ctrl.SetControllerReference(ironicConf, service, r.Scheme); err != nil {
		return fmt.Errorf("failed to set controller reference on networking service: %w", err)
	}

	// Create or update the service
	existingService := &corev1.Service{}
	err := cctx.Client.Get(cctx.Context, client.ObjectKey{
		Name:      service.Name,
		Namespace: service.Namespace,
	}, existingService)

	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Create new service
			cctx.Logger.Info("creating networking service", "Service", service.Name)
			return cctx.Client.Create(cctx.Context, service)
		}
		return fmt.Errorf("failed to get networking service: %w", err)
	}

	// Update existing service
	existingService.Spec.Ports = service.Spec.Ports
	existingService.Spec.Selector = service.Spec.Selector
	existingService.Labels = service.Labels
	cctx.Logger.Info("updating networking service", "Service", service.Name)
	return cctx.Client.Update(cctx.Context, existingService)
}

// deleteNetworkingResources deletes the networking deployment and service if they exist.
func (r *IronicReconciler) deleteNetworkingResources(cctx ironic.ControllerContext, ironicConf *metal3api.Ironic) error {
	// Delete deployment
	deployment := &appsv1.Deployment{}
	err := cctx.Client.Get(cctx.Context, client.ObjectKey{
		Name:      ironic.NetworkingDeploymentName(ironicConf),
		Namespace: ironicConf.Namespace,
	}, deployment)

	if err == nil {
		cctx.Logger.Info("deleting networking service deployment", "Deployment", deployment.Name)
		if err = cctx.Client.Delete(cctx.Context, deployment); err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete networking deployment: %w", err)
		}
	} else if !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to get networking deployment: %w", err)
	}

	// Delete service
	service := &corev1.Service{}
	err = cctx.Client.Get(cctx.Context, client.ObjectKey{
		Name:      ironic.NetworkingServiceName(ironicConf),
		Namespace: ironicConf.Namespace,
	}, service)

	if err == nil {
		cctx.Logger.Info("deleting networking service", "Service", service.Name)
		if err = cctx.Client.Delete(cctx.Context, service); err != nil && !k8serrors.IsNotFound(err) {
			return fmt.Errorf("failed to delete networking service: %w", err)
		}
	} else if !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to get networking service: %w", err)
	}

	return nil
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
