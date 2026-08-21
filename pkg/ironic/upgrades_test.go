package ironic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

// A Job's spec.template is immutable once created, so ensureIronicUpgradeJob
// must not attempt to change it on an already-existing job. When the desired
// template drifts (e.g. a toleration override is added) and the job has not
// started yet, it must be deleted and recreated instead of leaving the pod
// unschedulable forever.
func TestEnsureIronicUpgradeJobRecreatesOnTemplateDrift(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, batchv1.AddToScheme(scheme))
	require.NoError(t, metal3api.AddToScheme(scheme))

	c := fake.NewClientBuilder().WithScheme(scheme).Build()

	ironic := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{Name: "test-ironic", Namespace: "test-ns", UID: "test-uid"},
		Spec: metal3api.IronicSpec{
			Database: &metal3api.Database{
				Host:            "mariadb",
				Name:            "ironic",
				CredentialsName: "db-creds",
			},
		},
	}
	resources := Resources{Ironic: ironic}

	cctx := ControllerContext{
		Context: t.Context(),
		Client:  c,
		Scheme:  scheme,
		VersionInfo: VersionInfo{
			InstalledVersion: metal3api.Version380,
		},
	}

	jobKey := types.NamespacedName{Name: "test-ironic-pre-none-to-38.0", Namespace: "test-ns"}

	_, err := ensureIronicUpgradeJob(cctx, resources, preUpgrade)
	require.NoError(t, err)

	var job batchv1.Job
	require.NoError(t, c.Get(t.Context(), jobKey, &job))
	assert.Nil(t, job.Spec.Template.Spec.Tolerations)

	// A real API server always stamps CreationTimestamp on creation, which is what
	// ensureIronicUpgradeJob relies on to detect that the job already exists. The
	// fake client does not do this automatically, so it's backfilled here.
	job.CreationTimestamp = metav1.Now()
	require.NoError(t, c.Update(t.Context(), &job))

	// Simulate the user adding a toleration override while the same upgrade
	// job (same from/to version pair) still exists and has not started yet.
	ironic.Spec.Overrides = &metal3api.Overrides{
		Tolerations: []corev1.Toleration{
			{
				Key:      "node-role.kubernetes.io/control-plane",
				Operator: corev1.TolerationOpExists,
				Effect:   corev1.TaintEffectNoSchedule,
			},
		},
	}

	status, err := ensureIronicUpgradeJob(cctx, resources, preUpgrade)
	require.NoError(t, err)
	assert.True(t, status.NeedsRequeue(), "expected a requeue after deleting the stale job")

	// The stale job must be gone; ensureIronicUpgradeJob will recreate it
	// (with the new template) on the next reconcile.
	require.Error(t, c.Get(t.Context(), jobKey, &batchv1.Job{}))

	_, err = ensureIronicUpgradeJob(cctx, resources, preUpgrade)
	require.NoError(t, err)

	var recreatedJob batchv1.Job
	require.NoError(t, c.Get(t.Context(), jobKey, &recreatedJob))
	assert.Equal(t, ironic.Spec.Overrides.Tolerations, recreatedJob.Spec.Template.Spec.Tolerations)
}

// Once a job has started, its stale template must be left alone even if the
// desired template has since drifted: a running/completed job cannot be
// deleted without losing its (possibly in-progress) upgrade work.
func TestEnsureIronicUpgradeJobKeepsStartedJobOnDrift(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, batchv1.AddToScheme(scheme))
	require.NoError(t, metal3api.AddToScheme(scheme))

	c := fake.NewClientBuilder().WithScheme(scheme).Build()

	ironic := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{Name: "test-ironic", Namespace: "test-ns", UID: "test-uid"},
		Spec: metal3api.IronicSpec{
			Database: &metal3api.Database{
				Host:            "mariadb",
				Name:            "ironic",
				CredentialsName: "db-creds",
			},
		},
	}
	resources := Resources{Ironic: ironic}

	cctx := ControllerContext{
		Context: t.Context(),
		Client:  c,
		Scheme:  scheme,
		VersionInfo: VersionInfo{
			InstalledVersion: metal3api.Version380,
		},
	}

	jobKey := types.NamespacedName{Name: "test-ironic-pre-none-to-38.0", Namespace: "test-ns"}

	_, err := ensureIronicUpgradeJob(cctx, resources, preUpgrade)
	require.NoError(t, err)

	var job batchv1.Job
	require.NoError(t, c.Get(t.Context(), jobKey, &job))
	firstUID := job.UID

	// Mark the job as already started.
	job.CreationTimestamp = metav1.Now()
	require.NoError(t, c.Update(t.Context(), &job))

	startTime := metav1.Now()
	job.Status.StartTime = &startTime
	require.NoError(t, c.Status().Update(t.Context(), &job))

	var confirmJob batchv1.Job
	require.NoError(t, c.Get(t.Context(), jobKey, &confirmJob))
	require.NotNil(t, confirmJob.Status.StartTime, "test setup: job must be marked as started")

	ironic.Spec.Overrides = &metal3api.Overrides{
		Tolerations: []corev1.Toleration{
			{
				Key:      "node-role.kubernetes.io/control-plane",
				Operator: corev1.TolerationOpExists,
				Effect:   corev1.TaintEffectNoSchedule,
			},
		},
	}

	_, err = ensureIronicUpgradeJob(cctx, resources, preUpgrade)
	require.NoError(t, err)

	var unchangedJob batchv1.Job
	require.NoError(t, c.Get(t.Context(), jobKey, &unchangedJob))
	assert.Equal(t, firstUID, unchangedJob.UID, "a started job must not be deleted")
	assert.Nil(t, unchangedJob.Spec.Template.Spec.Tolerations)
}
