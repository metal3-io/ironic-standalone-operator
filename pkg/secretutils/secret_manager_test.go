package secretutils

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

func newTestScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = metal3api.AddToScheme(scheme)
	return scheme
}

func environmentLabels() map[string]string {
	return map[string]string{
		metal3api.LabelEnvironmentName: metal3api.LabelEnvironmentValue,
	}
}

func TestSecretManager_ObtainSecret_NotFound(t *testing.T) {
	scheme := newTestScheme()
	fakeClient := fake.NewClientBuilder().WithScheme(scheme).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	_, err := sm.ObtainSecret(types.NamespacedName{Name: "nonexistent", Namespace: "test"})
	require.Error(t, err)
}

func TestSecretManager_ObtainSecret_FoundInCache(t *testing.T) {
	scheme := newTestScheme()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "test",
			Labels:    environmentLabels(),
		},
		Data: map[string][]byte{
			"username": []byte("admin"),
			"password": []byte("secret"),
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	result, err := sm.ObtainSecret(types.NamespacedName{Name: "test-secret", Namespace: "test"})
	require.NoError(t, err)
	assert.Equal(t, "test-secret", result.Name)
	assert.Equal(t, metal3api.LabelEnvironmentValue, result.Labels[metal3api.LabelEnvironmentName])
}

func TestSecretManager_ObtainSecret_LabelAlreadySet(t *testing.T) {
	scheme := newTestScheme()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "test",
			Labels:    environmentLabels(),
		},
		Data: map[string][]byte{
			"username": []byte("admin"),
			"password": []byte("secret"),
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	result, err := sm.ObtainSecret(types.NamespacedName{Name: "test-secret", Namespace: "test"})
	require.NoError(t, err)
	assert.Equal(t, metal3api.LabelEnvironmentValue, result.Labels[metal3api.LabelEnvironmentName])

	// Verify the secret was not modified (no extra labels, no ownerRefs)
	var updated corev1.Secret
	err = fakeClient.Get(t.Context(), types.NamespacedName{Name: "test-secret", Namespace: "test"}, &updated)
	require.NoError(t, err)
	assert.Equal(t, metal3api.LabelEnvironmentValue, updated.Labels[metal3api.LabelEnvironmentName])
	assert.Empty(t, updated.OwnerReferences)
}

func TestSecretManager_AcquireSecret_WithOwner(t *testing.T) {
	scheme := newTestScheme()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "test",
			Labels:    environmentLabels(),
		},
		Data: map[string][]byte{
			"username": []byte("admin"),
			"password": []byte("secret"),
		},
	}

	owner := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ironic",
			Namespace: "test",
			UID:       "test-uid-12345",
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret, owner).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	result, err := sm.AcquireSecret(types.NamespacedName{Name: "test-secret", Namespace: "test"}, owner, scheme)
	require.NoError(t, err)
	assert.Equal(t, metal3api.LabelEnvironmentValue, result.Labels[metal3api.LabelEnvironmentName])

	// Verify owner reference was added
	var updated corev1.Secret
	err = fakeClient.Get(t.Context(), types.NamespacedName{Name: "test-secret", Namespace: "test"}, &updated)
	require.NoError(t, err)
	assert.Len(t, updated.OwnerReferences, 1)
	assert.Equal(t, owner.UID, updated.OwnerReferences[0].UID)
}

func TestSecretManager_AcquireSecret_AlreadyLabeled(t *testing.T) {
	scheme := newTestScheme()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "test",
			Labels:    environmentLabels(),
		},
		Data: map[string][]byte{
			"username": []byte("admin"),
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	result, err := sm.ObtainSecret(types.NamespacedName{Name: "test-secret", Namespace: "test"})
	require.NoError(t, err)
	assert.Equal(t, metal3api.LabelEnvironmentValue, result.Labels[metal3api.LabelEnvironmentName])
}

func TestSecretManager_AcquireSecret_AlreadyOwned(t *testing.T) {
	scheme := newTestScheme()

	owner := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ironic",
			Namespace: "test",
			UID:       "test-uid-12345",
		},
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "test",
			Labels:    environmentLabels(),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "ironic.metal3.io/v1alpha1",
					Kind:       "Ironic",
					Name:       owner.Name,
					UID:        owner.UID,
				},
			},
		},
		Data: map[string][]byte{
			"username": []byte("admin"),
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret, owner).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	result, err := sm.AcquireSecret(types.NamespacedName{Name: "test-secret", Namespace: "test"}, owner, scheme)
	require.NoError(t, err)
	assert.Equal(t, "test-secret", result.Name)

	// Should still have exactly one owner reference
	var updated corev1.Secret
	err = fakeClient.Get(t.Context(), types.NamespacedName{Name: "test-secret", Namespace: "test"}, &updated)
	require.NoError(t, err)
	assert.Len(t, updated.OwnerReferences, 1)
}

func TestSecretManager_FallbackToAPIReader(t *testing.T) {
	scheme := newTestScheme()

	// Secret with the label exists only in the API reader, not in the cache.
	// This simulates a labeled secret that hasn't been synced into the cache
	// yet, and exercises the findSecret fallback path to apiReader.
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "uncached-secret",
			Namespace: "test",
			Labels:    environmentLabels(),
		},
		Data: map[string][]byte{
			"data": []byte("value"),
		},
	}

	cacheClient := fake.NewClientBuilder().WithScheme(scheme).Build()
	apiReader := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), cacheClient, apiReader)

	result, err := sm.ObtainSecret(types.NamespacedName{Name: "uncached-secret", Namespace: "test"})
	require.NoError(t, err)
	assert.Equal(t, "uncached-secret", result.Name)
	assert.Equal(t, metal3api.LabelEnvironmentValue, result.Labels[metal3api.LabelEnvironmentName])
}

func TestSecretManager_NotInCacheButInAPI(t *testing.T) {
	scheme := newTestScheme()

	// Secret exists only in the API reader, not in the cache client.
	// It does NOT have the environment label, so claimSecret must reject it.
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "api-only-secret",
			Namespace: "test",
		},
		Data: map[string][]byte{
			"data": []byte("value"),
		},
	}

	cacheClient := fake.NewClientBuilder().WithScheme(scheme).Build()
	apiReader := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), cacheClient, apiReader)

	_, err := sm.ObtainSecret(types.NamespacedName{Name: "api-only-secret", Namespace: "test"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not have the required label")
}

func TestSecretManager_AcquireSecret_WithOwnerNoScheme(t *testing.T) {
	scheme := newTestScheme()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "test",
			Labels:    environmentLabels(),
		},
		Data: map[string][]byte{
			"username": []byte("admin"),
		},
	}

	owner := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ironic",
			Namespace: "test",
			UID:       "test-uid-12345",
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret, owner).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	// When scheme is nil, owner reference should not be set (only label checked)
	result, err := sm.AcquireSecret(types.NamespacedName{Name: "test-secret", Namespace: "test"}, owner, nil)
	require.NoError(t, err)
	assert.Equal(t, metal3api.LabelEnvironmentValue, result.Labels[metal3api.LabelEnvironmentName])

	// Verify NO owner reference was added (scheme was nil)
	var updated corev1.Secret
	err = fakeClient.Get(t.Context(), types.NamespacedName{Name: "test-secret", Namespace: "test"}, &updated)
	require.NoError(t, err)
	assert.Empty(t, updated.OwnerReferences)
}

func TestSecretManager_AcquireSecret_NilOwnerWithScheme(t *testing.T) {
	scheme := newTestScheme()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "test",
			Labels:    environmentLabels(),
		},
		Data: map[string][]byte{
			"username": []byte("admin"),
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	// When owner is nil, only label should be checked (same as ObtainSecret)
	result, err := sm.AcquireSecret(types.NamespacedName{Name: "test-secret", Namespace: "test"}, nil, scheme)
	require.NoError(t, err)
	assert.Equal(t, metal3api.LabelEnvironmentValue, result.Labels[metal3api.LabelEnvironmentName])

	// Verify NO owner reference was added (owner was nil)
	var updated corev1.Secret
	err = fakeClient.Get(t.Context(), types.NamespacedName{Name: "test-secret", Namespace: "test"}, &updated)
	require.NoError(t, err)
	assert.Empty(t, updated.OwnerReferences)
}

func TestSecretManager_MultipleOwners(t *testing.T) {
	scheme := newTestScheme()

	owner1 := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ironic-1",
			Namespace: "test",
			UID:       "uid-1",
		},
	}

	owner2 := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ironic-2",
			Namespace: "test",
			UID:       "uid-2",
		},
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "shared-secret",
			Namespace: "test",
			Labels:    environmentLabels(),
		},
		Data: map[string][]byte{
			"data": []byte("value"),
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret, owner1, owner2).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	// First owner acquires the secret
	_, err := sm.AcquireSecret(types.NamespacedName{Name: "shared-secret", Namespace: "test"}, owner1, scheme)
	require.NoError(t, err)

	// Second owner also acquires the same secret
	_, err = sm.AcquireSecret(types.NamespacedName{Name: "shared-secret", Namespace: "test"}, owner2, scheme)
	require.NoError(t, err)

	// Verify both owner references exist
	var updated corev1.Secret
	err = fakeClient.Get(t.Context(), types.NamespacedName{Name: "shared-secret", Namespace: "test"}, &updated)
	require.NoError(t, err)
	assert.Len(t, updated.OwnerReferences, 2)
}

func TestSecretManager_PreservesExistingLabels(t *testing.T) {
	scheme := newTestScheme()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "test",
			Labels: map[string]string{
				metal3api.LabelEnvironmentName: metal3api.LabelEnvironmentValue,
				"existing-label":               "existing-value",
			},
		},
		Data: map[string][]byte{
			"data": []byte("value"),
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	result, err := sm.ObtainSecret(types.NamespacedName{Name: "test-secret", Namespace: "test"})
	require.NoError(t, err)

	// Verify both labels exist
	assert.Equal(t, metal3api.LabelEnvironmentValue, result.Labels[metal3api.LabelEnvironmentName])
	assert.Equal(t, "existing-value", result.Labels["existing-label"])
}

// Rejection tests: verify that Secrets/ConfigMaps without the environment label are refused.

func TestSecretManager_ObtainSecret_NoLabel(t *testing.T) {
	scheme := newTestScheme()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "unlabeled-secret",
			Namespace: "test",
		},
		Data: map[string][]byte{
			"data": []byte("value"),
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	_, err := sm.ObtainSecret(types.NamespacedName{Name: "unlabeled-secret", Namespace: "test"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), metal3api.LabelEnvironmentName)
	assert.Contains(t, err.Error(), "does not have the required label")

	// Verify the secret was NOT modified (no label added, no ownerRef)
	var updated corev1.Secret
	err = fakeClient.Get(t.Context(), types.NamespacedName{Name: "unlabeled-secret", Namespace: "test"}, &updated)
	require.NoError(t, err)
	assert.Empty(t, updated.Labels)
	assert.Empty(t, updated.OwnerReferences)
}

func TestSecretManager_AcquireSecret_NoLabel(t *testing.T) {
	scheme := newTestScheme()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "unlabeled-secret",
			Namespace: "test",
		},
		Data: map[string][]byte{
			"data": []byte("value"),
		},
	}

	owner := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ironic",
			Namespace: "test",
			UID:       "test-uid-12345",
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret, owner).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	_, err := sm.AcquireSecret(types.NamespacedName{Name: "unlabeled-secret", Namespace: "test"}, owner, scheme)
	require.Error(t, err)
	assert.Contains(t, err.Error(), metal3api.LabelEnvironmentName)

	// Verify the secret was NOT modified (no label or ownerRef added)
	var updated corev1.Secret
	err = fakeClient.Get(t.Context(), types.NamespacedName{Name: "unlabeled-secret", Namespace: "test"}, &updated)
	require.NoError(t, err)
	assert.Empty(t, updated.Labels)
	assert.Empty(t, updated.OwnerReferences)
}

func TestSecretManager_LabelWrongValue(t *testing.T) {
	scheme := newTestScheme()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wrong-value-secret",
			Namespace: "test",
			Labels: map[string]string{
				metal3api.LabelEnvironmentName: "false",
			},
		},
		Data: map[string][]byte{
			"data": []byte("value"),
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	_, err := sm.ObtainSecret(types.NamespacedName{Name: "wrong-value-secret", Namespace: "test"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not have the required label")
}

func TestSecretManager_LabelEmptyValue(t *testing.T) {
	scheme := newTestScheme()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "empty-value-secret",
			Namespace: "test",
			Labels: map[string]string{
				metal3api.LabelEnvironmentName: "",
			},
		},
		Data: map[string][]byte{
			"data": []byte("value"),
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	_, err := sm.ObtainSecret(types.NamespacedName{Name: "empty-value-secret", Namespace: "test"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not have the required label")
}

func TestSecretManager_OtherLabelsButNoEnvironment(t *testing.T) {
	scheme := newTestScheme()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "other-labels-secret",
			Namespace: "test",
			Labels: map[string]string{
				"some-other-label": "some-value",
			},
		},
		Data: map[string][]byte{
			"data": []byte("value"),
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	_, err := sm.ObtainSecret(types.NamespacedName{Name: "other-labels-secret", Namespace: "test"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not have the required label")
}
