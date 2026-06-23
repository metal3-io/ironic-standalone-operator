package ironic

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

func ensureIronicIngress(cctx ControllerContext, ironic *metal3api.Ironic) (Status, error) {
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{Name: ironic.Name, Namespace: ironic.Namespace},
	}
	ingressSettings := ironic.Spec.Networking.Ingress
	result, err := controllerutil.CreateOrUpdate(cctx.Context, cctx.Client, ingress, func() error {
		if ingress.Labels == nil {
			cctx.Logger.Info("creating an ingress resource")
			ingress.Labels = make(map[string]string, 2)
		}
		ingress.Labels[metal3api.IronicServiceLabel] = ironic.Name
		ingress.Labels[metal3api.IronicVersionLabel] = cctx.VersionInfo.InstalledVersion.String()

		if ingressSettings.Annotations != nil {
			ingress.SetAnnotations(ingressSettings.Annotations)
		}
		if ingressSettings.IngressClassName != "" {
			ingress.Spec.IngressClassName = &ingressSettings.IngressClassName
		}
		ingress.Spec.TLS = []networkingv1.IngressTLS{{
			Hosts:      []string{ingressSettings.Host},
			SecretName: ironic.Name + "-ingress-tls",
		}}

		imagesPortNameIngress := imagesPortName
		if ironic.Spec.TLS.CertificateName != "" {
			imagesPortNameIngress = imagesTLSPortName
		}

		ingress.Spec.Rules = []networkingv1.IngressRule{{
			Host: ingressSettings.Host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{
						{
							Path:     "/",
							PathType: ptr.To(networkingv1.PathTypeImplementationSpecific),
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: ironic.Name,
									Port: networkingv1.ServiceBackendPort{
										Name: ironicPortName,
									},
								},
							},
						},
						{
							Path:     "/redfish",
							PathType: ptr.To(networkingv1.PathTypeImplementationSpecific),
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: ironic.Name,
									Port: networkingv1.ServiceBackendPort{
										Name: imagesPortNameIngress,
									},
								},
							},
						},
						{
							Path:     "/images",
							PathType: ptr.To(networkingv1.PathTypeImplementationSpecific),
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: ironic.Name,
									Port: networkingv1.ServiceBackendPort{
										Name: imagesPortNameIngress,
									},
								},
							},
						},
					},
				},
			},
		}}

		return controllerutil.SetControllerReference(ironic, ingress, cctx.Scheme)
	})
	if err != nil {
		return transientError(err)
	}
	if result != controllerutil.OperationResultNone {
		cctx.Logger.Info("ironic ingress", "Ingress", ingress.Name, "Status", result)
		return updated()
	}

	return ready()
}
