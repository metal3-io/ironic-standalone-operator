package ironic

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

// The fake client does not default anything, so mergeContainers can only be
// checked against hand-built "server" objects there. Here the Deployment goes
// through a real API server: a second pass must not report an update (which
// would requeue forever), and overrides of defaulted fields must still win.
func TestEnsureIronicDeploymentWithServerDefaults(t *testing.T) {
	c := newEnvtestClient(t)

	namespace := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{GenerateName: "irso-test-"}}
	require.NoError(t, c.Create(t.Context(), namespace))

	ironic := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: namespace.Name,
			UID:       "7b3cb6b5-7a43-4a2f-9a3c-7f5b3b1b8c10",
		},
		Spec: metal3api.IronicSpec{
			Images: metal3api.Images{
				DeployRamdiskDownloader: "quay.io/metal3-io/ironic-ipa-downloader:test",
				Ironic:                  "quay.io/metal3-io/ironic:test",
			},
		},
	}
	versionInfo, err := NewVersionInfo(ironic.Spec.Images, "")
	require.NoError(t, err)
	cctx := ControllerContext{
		Context:     t.Context(),
		Client:      c,
		Scheme:      testScheme,
		Logger:      logr.Discard(),
		VersionInfo: versionInfo,
	}
	resources := Resources{
		Ironic: ironic,
		APISecret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{Name: "test-api", Namespace: namespace.Name},
			Data:       map[string][]byte{htpasswdKey: []byte("abcd")},
		},
	}

	getContainers := func(t *testing.T) (containers map[string]corev1.Container) {
		t.Helper()

		deploy := &appsv1.Deployment{}
		require.NoError(t, c.Get(t.Context(), client.ObjectKey{Namespace: namespace.Name, Name: ironicDeploymentName(ironic)}, deploy))
		containers = make(map[string]corev1.Container, len(deploy.Spec.Template.Spec.Containers))
		for _, cont := range deploy.Spec.Template.Spec.Containers {
			containers[cont.Name] = cont
		}
		return containers
	}

	t.Run("create", func(t *testing.T) {
		status, err := ensureIronicDeployment(cctx, resources)
		require.NoError(t, err)
		assert.True(t, status.NeedsRequeue(), "the deployment must be created")

		// Make sure the API server did add the defaults that mergeContainers has to preserve.
		cont := getContainers(t)[ironicContainerName]
		assert.Equal(t, "/dev/termination-log", cont.TerminationMessagePath)
		assert.Equal(t, corev1.TerminationMessageReadFile, cont.TerminationMessagePolicy)
		assert.Equal(t, corev1.PullIfNotPresent, cont.ImagePullPolicy)
		assert.Equal(t, int32(10), cont.LivenessProbe.PeriodSeconds)
		assert.Equal(t, int32(1), cont.LivenessProbe.SuccessThreshold)
	})

	t.Run("no-changes", func(t *testing.T) {
		status, err := ensureIronicDeployment(cctx, resources)
		require.NoError(t, err)
		assert.False(t, status.NeedsRequeue(), "the deployment must not be updated again: %s", status)
	})

	t.Run("overrides-of-defaulted-fields", func(t *testing.T) {
		ironic.Spec.Overrides = &metal3api.Overrides{
			HttpdLivenessProbe: &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{Port: intstr.FromInt32(8080)},
				},
				TimeoutSeconds: 42,
			},
			Containers: []corev1.Container{{
				Name:            "ramdisk-logs",
				Image:           "quay.io/metal3-io/ironic:test",
				ImagePullPolicy: corev1.PullAlways,
			}},
		}

		status, err := ensureIronicDeployment(cctx, resources)
		require.NoError(t, err)
		assert.True(t, status.NeedsRequeue(), "the overrides must be applied")

		containers := getContainers(t)
		assert.Equal(t, corev1.PullAlways, containers["ramdisk-logs"].ImagePullPolicy)
		httpdProbe := containers[httpdContainerName].LivenessProbe
		assert.Equal(t, int32(42), httpdProbe.TimeoutSeconds)
		assert.Equal(t, corev1.URISchemeHTTP, httpdProbe.HTTPGet.Scheme)
		assert.Equal(t, "/", httpdProbe.HTTPGet.Path)
		// Defaults of the fields that were not overridden stay.
		assert.Equal(t, int32(10), httpdProbe.PeriodSeconds)
		assert.Equal(t, int32(1), httpdProbe.SuccessThreshold)
		assert.Equal(t, "/dev/termination-log", containers["ramdisk-logs"].TerminationMessagePath)

		status, err = ensureIronicDeployment(cctx, resources)
		require.NoError(t, err)
		assert.False(t, status.NeedsRequeue(), "the deployment must not be updated again: %s", status)
	})
}
