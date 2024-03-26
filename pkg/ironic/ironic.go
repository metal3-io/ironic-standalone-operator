package ironic

import (
	"context"
	"errors"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
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

func ensureIronicDaemonSet(cctx ControllerContext, ironic *metal3api.Ironic, db *metal3api.IronicDatabase, apiSecret *corev1.Secret) (status metal3api.IronicStatusConditionType, err error) {
	template, err := newIronicPodTemplate(ironic, db, apiSecret, cctx.Domain)
	if err != nil {
		return metal3api.IronicStatusProgressing, err
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
			matchLabels := map[string]string{metal3api.IronicOperatorLabel: ironicDeploymentName(ironic)}
			deploy.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: matchLabels,
			}
		}
		deploy.Spec.Template = template
		return controllerutil.SetControllerReference(ironic, deploy, cctx.Scheme)
	})
	if err != nil {
		return metal3api.IronicStatusProgressing, err
	}
	if result != controllerutil.OperationResultNone {
		cctx.Logger.Info("ironic daemon set", "DaemonSet", deploy.Name, "Status", result)
	}
	return getDaemonSetStatus(cctx, deploy)
}

func ensureIronicDeployment(cctx ControllerContext, ironic *metal3api.Ironic, db *metal3api.IronicDatabase, apiSecret *corev1.Secret) (status metal3api.IronicStatusConditionType, err error) {
	template, err := newIronicPodTemplate(ironic, db, apiSecret, cctx.Domain)
	if err != nil {
		return metal3api.IronicStatusProgressing, err
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
			matchLabels := map[string]string{metal3api.IronicOperatorLabel: ironicDeploymentName(ironic)}
			deploy.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: matchLabels,
			}
			deploy.Spec.Replicas = ptr.To(int32(1))
		}
		deploy.Spec.Template = template
		// We cannot run two copies of Ironic in parallel
		deploy.Spec.Strategy = appsv1.DeploymentStrategy{
			Type: appsv1.RecreateDeploymentStrategyType,
		}
		return controllerutil.SetControllerReference(ironic, deploy, cctx.Scheme)
	})
	if err != nil {
		return metal3api.IronicStatusProgressing, err
	}
	if result != controllerutil.OperationResultNone {
		cctx.Logger.Info("ironic deployment", "Deployment", deploy.Name, "Status", result)
	}
	return getDeploymentStatus(cctx, deploy)
}

func ensureIronicService(cctx ControllerContext, ironic *metal3api.Ironic) (metal3api.IronicStatusConditionType, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: ironic.Name, Namespace: ironic.Namespace},
	}
	exposedPort := 80
	if ironic.Spec.TLSRef.Name != "" {
		exposedPort = 443
	}
	result, err := controllerutil.CreateOrUpdate(cctx.Context, cctx.Client, service, func() error {
		if service.ObjectMeta.Labels == nil {
			cctx.Logger.Info("creating a new ironic service")
			service.ObjectMeta.Labels = make(map[string]string)
		}
		service.ObjectMeta.Labels[metal3api.IronicOperatorLabel] = ironicDeploymentName(ironic)

		service.Spec.Selector = map[string]string{metal3api.IronicOperatorLabel: ironicDeploymentName(ironic)}
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
	}
	if err != nil || len(service.Spec.ClusterIPs) == 0 {
		return metal3api.IronicStatusProgressing, err
	}

	return metal3api.IronicStatusAvailable, nil
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
func EnsureIronic(cctx ControllerContext, ironic *metal3api.Ironic, db *metal3api.IronicDatabase, apiSecret *corev1.Secret) (status metal3api.IronicStatusConditionType, err error) {
	if db != nil && !isReady(db.Status.Conditions) {
		cctx.Logger.Info("database is not ready yet")
		return metal3api.IronicStatusProgressing, nil
	}

	if ironic.Spec.Distributed {
		if db == nil {
			return metal3api.IronicStatusProgressing, errors.New("database is required for a distributed deployment")
		}

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

	if err != nil {
		return
	}

	// Let the service be created while Ironic is being deployed, but do
	// not report overall success until both are done.
	serviceStatus, err := ensureIronicService(cctx, ironic)
	if err != nil || serviceStatus != metal3api.IronicStatusAvailable {
		return serviceStatus, err
	}

	return
}

// RemoveIronic removes all bits of the Ironic deployment.
func RemoveIronic(cctx ControllerContext, ironic *metal3api.Ironic) error {
	return nil // rely on ownership-based clean up
}
