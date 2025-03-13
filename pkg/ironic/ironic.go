package ironic

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

func ironicDeploymentName(ironic *metal3api.Ironic) string {
	return fmt.Sprintf("%s-service", ironic.Name)
}

func ensureIronicDaemonSet(cctx ControllerContext, ironic *metal3api.Ironic, db *metal3api.IronicDatabase, apiSecret *corev1.Secret) (Status, error) {
	template, err := newIronicPodTemplate(cctx, ironic, db, apiSecret, cctx.Domain)
	if err != nil {
		return transientError(err)
	}

	deploy := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ironicDeploymentName(ironic),
			Namespace: ironic.Namespace,
		},
	}
	result, err := controllerutil.CreateOrUpdate(cctx.Context, cctx.Client, deploy, func() error {
		if deploy.ObjectMeta.CreationTimestamp.IsZero() {
			cctx.Logger.Info("creating a new ironic daemon set")
		}
		if deploy.Labels == nil {
			deploy.Labels = make(map[string]string, 2)
		}
		deploy.Labels[metal3api.IronicAppLabel] = ironicDeploymentName(ironic)
		deploy.Labels[metal3api.IronicVersionLabel] = cctx.VersionInfo.InstalledVersion.String()

		matchLabels := map[string]string{metal3api.IronicAppLabel: ironicDeploymentName(ironic)}
		deploy.Spec.Selector = &metav1.LabelSelector{MatchLabels: matchLabels}
		mergePodTemplates(&deploy.Spec.Template, template)
		return controllerutil.SetControllerReference(ironic, deploy, cctx.Scheme)
	})
	if err != nil {
		return transientError(err)
	}
	if result != controllerutil.OperationResultNone {
		cctx.Logger.Info("ironic daemon set", "DaemonSet", deploy.Name, "Status", result)
		return updated()
	}

	return getDaemonSetStatus(cctx, deploy)
}

func ensureIronicDeployment(cctx ControllerContext, ironic *metal3api.Ironic, db *metal3api.IronicDatabase, apiSecret *corev1.Secret) (Status, error) {
	template, err := newIronicPodTemplate(cctx, ironic, db, apiSecret, cctx.Domain)
	if err != nil {
		return transientError(err)
	}

	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ironicDeploymentName(ironic),
			Namespace: ironic.Namespace,
		},
	}
	result, err := controllerutil.CreateOrUpdate(cctx.Context, cctx.Client, deploy, func() error {
		if deploy.ObjectMeta.CreationTimestamp.IsZero() {
			cctx.Logger.Info("creating a new ironic deployment")
		}
		if deploy.Labels == nil {
			deploy.Labels = make(map[string]string, 2)
		}
		deploy.Labels[metal3api.IronicAppLabel] = ironicDeploymentName(ironic)
		deploy.Labels[metal3api.IronicVersionLabel] = cctx.VersionInfo.InstalledVersion.String()

		matchLabels := map[string]string{metal3api.IronicAppLabel: ironicDeploymentName(ironic)}
		deploy.Spec.Selector = &metav1.LabelSelector{MatchLabels: matchLabels}
		deploy.Spec.Replicas = ptr.To(int32(1))
		mergePodTemplates(&deploy.Spec.Template, template)
		// We cannot run two copies of Ironic in parallel
		deploy.Spec.Strategy = appsv1.DeploymentStrategy{
			Type: appsv1.RecreateDeploymentStrategyType,
		}
		return controllerutil.SetControllerReference(ironic, deploy, cctx.Scheme)
	})
	if err != nil {
		return transientError(err)
	}
	if result != controllerutil.OperationResultNone {
		cctx.Logger.Info("ironic deployment", "Deployment", deploy.Name, "Status", result)
		return updated()
	}

	return getDeploymentStatus(cctx, deploy)
}

func ensureIronicService(cctx ControllerContext, ironic *metal3api.Ironic) (Status, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: ironic.Name, Namespace: ironic.Namespace},
	}
	exposedPort := 80
	if ironic.Spec.TLS.CertificateName != "" {
		exposedPort = 443
	}
	result, err := controllerutil.CreateOrUpdate(cctx.Context, cctx.Client, service, func() error {
		if service.ObjectMeta.Labels == nil {
			cctx.Logger.Info("creating a new ironic service")
			service.ObjectMeta.Labels = make(map[string]string)
		}
		service.ObjectMeta.Labels[metal3api.IronicAppLabel] = ironicDeploymentName(ironic)

		service.Spec.Selector = map[string]string{metal3api.IronicAppLabel: ironicDeploymentName(ironic)}
		service.Spec.Ports = []corev1.ServicePort{
			{
				Protocol:   corev1.ProtocolTCP,
				Port:       int32(exposedPort),
				TargetPort: intstr.FromString(ironicPortName),
			},
		}
		service.Spec.Type = corev1.ServiceTypeClusterIP

		return controllerutil.SetControllerReference(ironic, service, cctx.Scheme)
	})
	if result != controllerutil.OperationResultNone {
		cctx.Logger.Info("ironic service", "Service", service.Name, "Status", result)
		return updated()
	}
	if err != nil {
		return transientError(err)
	}

	return getServiceStatus(service)
}

func ensureIronicUpgradeJob(cctx ControllerContext, ironic *metal3api.Ironic, phase string) (Status, error) {
	if ironic.Status.InstalledVersion == "" {
		// Nothing to upgrade: no current version
		return ready()
	}

	currentVersion, err := metal3api.ParseVersion(ironic.Status.InstalledVersion)
	if err != nil {
		return Status{Fatal: fmt.Errorf("malformed InstalledVersion in Ironic: %s", err)}, nil
	}

	if !currentVersion.IsLatest() && cctx.VersionInfo.InstalledVersion.Compare(currentVersion) <= 0 {
		// Upgrade impossible: requested version is not newer, and the current version is not "latest"
		return ready()
	}

	if !cctx.VersionInfo.InstalledVersion.HasUpgradeScripts() {
		cctx.Logger.Info("Not running upgrade scripts: the new version does not support them", "From", ironic.Status.InstalledVersion, "To", cctx.VersionInfo.InstalledVersion.String())
		return ready()
	}

	var template corev1.PodTemplateSpec
	if phase == "pre" {
		template, err = newPreUpgradePodTemplate(cctx, ironic)
	} else {
		template, err = newPostUpgradePodTemplate(cctx, ironic)
	}
	if err != nil {
		return transientError(err)
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s-%s-to-%s", ironic.Name, phase, ironic.Status.InstalledVersion, cctx.VersionInfo.InstalledVersion),
			Namespace: ironic.Namespace,
		},
	}

	result, err := controllerutil.CreateOrUpdate(cctx.Context, cctx.Client, job, func() error {
		if job.ObjectMeta.Labels == nil {
			cctx.Logger.Info(fmt.Sprintf("creating a new %s-upgrade job", phase), "From", ironic.Status.InstalledVersion, "To", cctx.VersionInfo.InstalledVersion.String())
			job.ObjectMeta.Labels = make(map[string]string)
		}
		job.ObjectMeta.Labels[metal3api.IronicAppLabel] = ironicDeploymentName(ironic)

		mergePodTemplates(&job.Spec.Template, template)

		return controllerutil.SetControllerReference(ironic, job, cctx.Scheme)
	})
	if result != controllerutil.OperationResultNone {
		cctx.Logger.Info("ironic upgrade job", "Service", job.Name, "Status", result)
		return updated()
	}
	if err != nil {
		return transientError(err)
	}

	return getJobStatus(cctx, job, fmt.Sprintf("%s-upgrade", phase))
}

func removeIronicDaemonSet(cctx ControllerContext, ironic *metal3api.Ironic) error {
	err := cctx.KubeClient.AppsV1().DaemonSets(ironic.Namespace).
		Delete(context.Background(), ironicDeploymentName(ironic), metav1.DeleteOptions{})
	return client.IgnoreNotFound(err)
}

func removeIronicDeployment(cctx ControllerContext, ironic *metal3api.Ironic) error {
	err := cctx.KubeClient.AppsV1().Deployments(ironic.Namespace).
		Delete(context.Background(), ironicDeploymentName(ironic), metav1.DeleteOptions{})
	return client.IgnoreNotFound(err)
}

// EnsureIronic deploys Ironic either as a Deployment or as a DaemonSet.
func EnsureIronic(cctx ControllerContext, ironic *metal3api.Ironic, db *metal3api.IronicDatabase, apiSecret *corev1.Secret) (status Status, err error) {
	if db != nil && !isReady(db.Status.Conditions) {
		return inProgress("database is not ready yet")
	}

	if validationErr := metal3api.ValidateIronic(&ironic.Spec, nil); validationErr != nil {
		status = Status{Fatal: validationErr}
		return
	}

	jobStatus, err := ensureIronicUpgradeJob(cctx, ironic, "pre")
	if err != nil || !jobStatus.IsReady() {
		return jobStatus, err
	}

	if ironic.Spec.HighAvailability {
		err = removeIronicDeployment(cctx, ironic)
		if err != nil {
			return
		}
		status, err = ensureIronicDaemonSet(cctx, ironic, db, apiSecret)
	} else {
		err = removeIronicDaemonSet(cctx, ironic)
		if err != nil {
			return
		}
		status, err = ensureIronicDeployment(cctx, ironic, db, apiSecret)
	}

	if err != nil || status.IsError() {
		return
	}

	// Let the service be created while Ironic is being deployed, but do
	// not report overall success until both are done.
	serviceStatus, err := ensureIronicService(cctx, ironic)
	if err != nil || !serviceStatus.IsReady() {
		return serviceStatus, err
	}

	jobStatus, err = ensureIronicUpgradeJob(cctx, ironic, "post")
	if err != nil || !jobStatus.IsReady() {
		return jobStatus, err
	}

	return
}

// RemoveIronic removes all bits of the Ironic deployment.
func RemoveIronic(cctx ControllerContext, ironic *metal3api.Ironic) error {
	return nil // rely on ownership-based clean up
}
