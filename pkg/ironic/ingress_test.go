package ironic

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

func TestEnsureIronicIngressKeepsForeignAnnotations(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, networkingv1.AddToScheme(scheme))
	require.NoError(t, metal3api.AddToScheme(scheme))

	ironic := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{Name: "test-ironic", Namespace: "test-ns"},
		Spec: metal3api.IronicSpec{
			Networking: metal3api.Networking{
				Ingress: &metal3api.Ingress{
					Host:        "ironic.example.com",
					Annotations: map[string]string{"cert-manager.io/cluster-issuer": "acme"},
				},
			},
		},
	}

	existing := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        ironic.Name,
			Namespace:   ironic.Namespace,
			Annotations: map[string]string{"another-controller.example.com/state": "set by someone else"},
		},
	}

	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(existing).Build()
	cctx := ControllerContext{
		Context:     context.Background(),
		Client:      cl,
		Scheme:      scheme,
		Logger:      logr.Discard(),
		VersionInfo: VersionInfo{InstalledVersion: metal3api.VersionLatest},
	}

	_, err := ensureIronicIngress(cctx, ironic)
	require.NoError(t, err)

	updated := &networkingv1.Ingress{}
	require.NoError(t, cl.Get(cctx.Context, types.NamespacedName{Name: ironic.Name, Namespace: ironic.Namespace}, updated))
	assert.Equal(t, "acme", updated.Annotations["cert-manager.io/cluster-issuer"])
	assert.Equal(t, "set by someone else", updated.Annotations["another-controller.example.com/state"])
}
