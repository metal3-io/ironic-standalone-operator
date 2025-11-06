package ironic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

func TestPrometheusExporterVersionCheck(t *testing.T) {
	testCases := []struct {
		name          string
		version       metal3api.Version
		enabled       bool
		expectedError string
	}{
		{
			name:          "PrometheusExporter with version 31.0",
			version:       metal3api.Version310,
			enabled:       true,
			expectedError: "",
		},
		{
			name:          "PrometheusExporter with version 32.0",
			version:       metal3api.Version320,
			enabled:       true,
			expectedError: "",
		},
		{
			name:          "PrometheusExporter with latest version",
			version:       metal3api.VersionLatest,
			enabled:       true,
			expectedError: "",
		},
		{
			name:          "PrometheusExporter with version 30.0 (too old)",
			version:       metal3api.Version300,
			enabled:       true,
			expectedError: "using prometheusExporter is only possible for Ironic 31.0 or newer",
		},
		{
			name:          "PrometheusExporter disabled with version 30.0",
			version:       metal3api.Version300,
			enabled:       false,
			expectedError: "",
		},
		{
			name:          "PrometheusExporter not configured with version 30.0",
			version:       metal3api.Version300,
			enabled:       false,
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			_ = metal3api.AddToScheme(scheme)
			_ = corev1.AddToScheme(scheme)

			var prometheusExporter *metal3api.PrometheusExporter
			if tc.name != "PrometheusExporter not configured with version 30.0" {
				prometheusExporter = &metal3api.PrometheusExporter{
					Enabled:                  tc.enabled,
					SensorCollectionInterval: 60,
				}
			}

			ironic := &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test",
				},
				Spec: metal3api.IronicSpec{
					PrometheusExporter: prometheusExporter,
				},
			}

			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "api-secret",
					Namespace: "test",
				},
				Data: map[string][]byte{
					"htpasswd": []byte("test:password"),
				},
			}

			clientBuilder := fakeclient.NewClientBuilder().WithScheme(scheme).WithObjects(ironic, secret)
			fakeClient := clientBuilder.Build()

			cctx := ControllerContext{
				Context:    t.Context(),
				Client:     fakeClient,
				KubeClient: fake.NewSimpleClientset(),
				Scheme:     scheme,
				VersionInfo: VersionInfo{
					InstalledVersion: tc.version,
					IronicImage:      "test-ironic:latest",
				},
			}

			resources := Resources{
				Ironic:    ironic,
				APISecret: secret,
			}

			status, err := EnsureIronic(cctx, resources)

			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.True(t, status.IsError())
			} else if err != nil {
				// We expect the function to proceed without version errors
				// (it may fail for other reasons like missing resources, but not version checks)
				assert.NotContains(t, err.Error(), "prometheusExporter is only possible")
			}
		})
	}
}

func TestBMCCAVersionCheck(t *testing.T) {
	testCases := []struct {
		name          string
		version       metal3api.Version
		expectedError string
	}{
		{
			name:          "BMCCA with version 32.0",
			version:       metal3api.Version320,
			expectedError: "",
		},
		{
			name:          "BMCCA with latest version",
			version:       metal3api.VersionLatest,
			expectedError: "",
		},
		{
			name:          "BMCCA with version 31.0 (too old)",
			version:       metal3api.Version310,
			expectedError: "using tls.bmcCAName is only possible for Ironic 32.0 or newer",
		},
		{
			name:          "BMCCA with version 30.0 (too old)",
			version:       metal3api.Version300,
			expectedError: "using tls.bmcCAName is only possible for Ironic 32.0 or newer",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			scheme := runtime.NewScheme()
			_ = metal3api.AddToScheme(scheme)
			_ = corev1.AddToScheme(scheme)

			apiSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "api-secret",
					Namespace: "test",
				},
				Data: map[string][]byte{
					"htpasswd": []byte("test:password"),
				},
			}

			bmcSecret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "bmc-ca",
					Namespace: "test",
				},
				Data: map[string][]byte{
					"ca.crt": []byte("test-ca-cert"),
				},
			}

			ironic := &metal3api.Ironic{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-ironic",
					Namespace: "test",
				},
				Spec: metal3api.IronicSpec{
					TLS: metal3api.TLS{
						BMCCAName: "bmc-ca",
					},
				},
			}

			clientBuilder := fakeclient.NewClientBuilder().WithScheme(scheme).WithObjects(ironic, apiSecret, bmcSecret)
			fakeClient := clientBuilder.Build()

			cctx := ControllerContext{
				Context:    t.Context(),
				Client:     fakeClient,
				KubeClient: fake.NewSimpleClientset(),
				Scheme:     scheme,
				VersionInfo: VersionInfo{
					InstalledVersion: tc.version,
					IronicImage:      "test-ironic:latest",
				},
			}

			resources := Resources{
				Ironic:      ironic,
				APISecret:   apiSecret,
				BMCCASecret: bmcSecret,
			}

			status, err := EnsureIronic(cctx, resources)

			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.True(t, status.IsError())
			} else if err != nil {
				// We expect the function to proceed without version errors
				assert.NotContains(t, err.Error(), "bmcCAName is only possible")
			}
		})
	}
}
