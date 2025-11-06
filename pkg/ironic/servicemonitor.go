package ironic

import (
	"context"
	"reflect"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

func serviceMonitorName(ironic *metal3api.Ironic) string {
	return ironic.Name + "-metrics"
}

// newServiceMonitor creates a ServiceMonitor for Ironic Prometheus Exporter metrics.
func newServiceMonitor(ironic *metal3api.Ironic) *monitoringv1.ServiceMonitor {
	labels := map[string]string{
		metal3api.IronicServiceLabel: ironic.Name,
	}

	return &monitoringv1.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceMonitorName(ironic),
			Namespace: ironic.Namespace,
			Labels:    labels,
		},
		Spec: monitoringv1.ServiceMonitorSpec{
			Endpoints: []monitoringv1.Endpoint{
				{
					Port: metricsPortName,
				},
			},
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					metal3api.IronicServiceLabel: ironic.Name,
				},
			},
		},
	}
}

// ensureServiceMonitor ensures the ServiceMonitor exists when PrometheusExporter is enabled.
func ensureServiceMonitor(cctx ControllerContext, ironic *metal3api.Ironic) (Status, error) {
	// If PrometheusExporter is disabled or ServiceMonitor is explicitly disabled, ensure ServiceMonitor is deleted
	if ironic.Spec.PrometheusExporter == nil || !ironic.Spec.PrometheusExporter.Enabled || ironic.Spec.PrometheusExporter.DisableServiceMonitor {
		return removeServiceMonitor(cctx, ironic)
	}

	sm := newServiceMonitor(ironic)
	if err := controllerutil.SetControllerReference(ironic, sm, cctx.Scheme); err != nil {
		return transientError(err)
	}

	existing := &monitoringv1.ServiceMonitor{}
	err := cctx.Client.Get(cctx.Context, client.ObjectKey{
		Namespace: ironic.Namespace,
		Name:      serviceMonitorName(ironic),
	}, existing)

	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			// ServiceMonitor doesn't exist, create it
			cctx.Logger.Info("creating ServiceMonitor", "ServiceMonitor", sm.Name)
			err = cctx.Client.Create(cctx.Context, sm)
			if err != nil {
				return transientError(err)
			}
			return updated()
		}
		return transientError(err)
	}

	// ServiceMonitor exists, update if needed
	if !reflect.DeepEqual(&sm.Spec, &existing.Spec) {
		existing.Spec = sm.Spec
		err = cctx.Client.Update(cctx.Context, existing)
		if err != nil {
			return transientError(err)
		}
		return updated()
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
