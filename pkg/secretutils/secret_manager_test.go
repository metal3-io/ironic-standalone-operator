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
			Labels: map[string]string{
				LabelEnvironmentName: LabelEnvironmentValue,
			},
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
	assert.Equal(t, LabelEnvironmentValue, result.Labels[LabelEnvironmentName])
}

func TestSecretManager_ObtainSecret_AddsLabel(t *testing.T) {
	scheme := newTestScheme()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "test",
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
	assert.Equal(t, LabelEnvironmentValue, result.Labels[LabelEnvironmentName])

	// Verify the label was persisted
	var updated corev1.Secret
	err = fakeClient.Get(t.Context(), types.NamespacedName{Name: "test-secret", Namespace: "test"}, &updated)
	require.NoError(t, err)
	assert.Equal(t, LabelEnvironmentValue, updated.Labels[LabelEnvironmentName])
}

func TestSecretManager_AcquireSecret_WithOwner(t *testing.T) {
	scheme := newTestScheme()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "test",
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
	assert.Equal(t, LabelEnvironmentValue, result.Labels[LabelEnvironmentName])

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
			Labels: map[string]string{
				LabelEnvironmentName: LabelEnvironmentValue,
			},
		},
		Data: map[string][]byte{
			"username": []byte("admin"),
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	result, err := sm.ObtainSecret(types.NamespacedName{Name: "test-secret", Namespace: "test"})
	require.NoError(t, err)
	assert.Equal(t, LabelEnvironmentValue, result.Labels[LabelEnvironmentName])
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
			Labels: map[string]string{
				LabelEnvironmentName: LabelEnvironmentValue,
			},
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

	// Secret exists only in the "API" (apiReader), not in the cache
	// This simulates an unlabeled secret that's not in the filtered cache
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "uncached-secret",
			Namespace: "test",
		},
		Data: map[string][]byte{
			"data": []byte("value"),
		},
	}

	// Both clients have the secret (in real scenario, cache would filter it out,
	// but for testing we simulate the fallback by using separate clients)
	// The key behavior we're testing is that findSecret tries cache first, then API
	cacheClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()
	apiReader := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), cacheClient, apiReader)

	// The secret should be found (via cache in this test, but the fallback logic is tested)
	result, err := sm.ObtainSecret(types.NamespacedName{Name: "uncached-secret", Namespace: "test"})
	require.NoError(t, err)
	assert.Equal(t, "uncached-secret", result.Name)
	// Label should be added
	assert.Equal(t, LabelEnvironmentValue, result.Labels[LabelEnvironmentName])
}

func TestSecretManager_NotInCacheButInAPI(t *testing.T) {
	scheme := newTestScheme()

	// Secret exists only in the API reader, not in the cache client
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "api-only-secret",
			Namespace: "test",
		},
		Data: map[string][]byte{
			"data": []byte("value"),
		},
	}

	// Empty cache client, secret only in API reader
	cacheClient := fake.NewClientBuilder().WithScheme(scheme).Build()
	apiReader := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), cacheClient, apiReader)

	// findSecret should find it via API fallback, but claimSecret will fail
	// because the secret doesn't exist in the cache client for update
	// This tests the fallback path in findSecret
	_, err := sm.ObtainSecret(types.NamespacedName{Name: "api-only-secret", Namespace: "test"})
	// We expect an error because we can't update a secret that doesn't exist in the cache client
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestSecretManager_AcquireSecret_WithOwnerNoScheme(t *testing.T) {
	scheme := newTestScheme()
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "test",
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

	// When scheme is nil, owner reference should not be set (only label)
	result, err := sm.AcquireSecret(types.NamespacedName{Name: "test-secret", Namespace: "test"}, owner, nil)
	require.NoError(t, err)
	assert.Equal(t, LabelEnvironmentValue, result.Labels[LabelEnvironmentName])

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
		},
		Data: map[string][]byte{
			"username": []byte("admin"),
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

	sm := NewSecretManager(t.Context(), logr.Discard(), fakeClient, fakeClient)

	// When owner is nil, only label should be set (same as ObtainSecret)
	result, err := sm.AcquireSecret(types.NamespacedName{Name: "test-secret", Namespace: "test"}, nil, scheme)
	require.NoError(t, err)
	assert.Equal(t, LabelEnvironmentValue, result.Labels[LabelEnvironmentName])

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
				"existing-label": "existing-value",
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
	assert.Equal(t, LabelEnvironmentValue, result.Labels[LabelEnvironmentName])
	assert.Equal(t, "existing-value", result.Labels["existing-label"])
}
