package ironic

import (
	"context"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

func serviceMonitorName(ironic *metal3api.Ironic) string {
	return ironic.Name + "-metrics"
}

// ensureServiceMonitor ensures the ServiceMonitor exists when PrometheusExporter is enabled.
func ensureServiceMonitor(cctx ControllerContext, ironic *metal3api.Ironic) (Status, error) {
	// If PrometheusExporter is disabled or ServiceMonitor is explicitly disabled, ensure ServiceMonitor is deleted
	if ironic.Spec.PrometheusExporter == nil || !ironic.Spec.PrometheusExporter.Enabled || ironic.Spec.PrometheusExporter.DisableServiceMonitor {
		return removeServiceMonitor(cctx, ironic)
	}

	sm := &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceMonitorName(ironic),
			Namespace: ironic.Namespace,
		},
	}

	result, err := controllerutil.CreateOrUpdate(cctx.Context, cctx.Client, sm, func() error {
		if sm.Labels == nil {
			sm.Labels = make(map[string]string, 1)
		}
		sm.Labels[metal3api.IronicServiceLabel] = ironic.Name

		sm.Spec.Endpoints = []monitoringv1.Endpoint{
			{
				Port: metricsPortName,
				// TODO(dtantsur): TLS support?
				// NOTE: The ironic-prometheus-exporter binds to localhost (127.0.0.1) by default
				// to prevent unauthenticated access from the host network. This means that
				// ServiceMonitor-based scraping from a remote Prometheus instance will NOT work
				// unless the exporter is explicitly configured to listen on all interfaces.
				Scheme: ptr.To(monitoringv1.SchemeHTTP),
				Path:   "/metrics",
			},
		}
		sm.Spec.Selector = metav1.LabelSelector{
			MatchLabels: map[string]string{
				metal3api.IronicServiceLabel: ironic.Name,
			},
		}

		return controllerutil.SetControllerReference(ironic, sm, cctx.Scheme)
	})

	if result != controllerutil.OperationResultNone {
		cctx.Logger.Info("ServiceMonitor", "ServiceMonitor", sm.Name, "Status", result)
		return updated()
	}
	if err != nil {
		return transientError(err)
	}

	return ready()
}

// removeServiceMonitor deletes the ServiceMonitor.
func removeServiceMonitor(cctx ControllerContext, ironic *metal3api.Ironic) (Status, error) {
	sm := &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceMonitorName(ironic),
			Namespace: ironic.Namespace,
		},
	}

	err := cctx.Client.Delete(context.Background(), sm)
	// Ignore NotFound errors and NoMatchError (API not available)
	if err == nil || k8serrors.IsNotFound(err) || meta.IsNoMatchError(err) {
		return ready()
	}
	return transientError(err)
}
