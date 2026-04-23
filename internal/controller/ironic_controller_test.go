/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/events"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
	"github.com/metal3-io/ironic-standalone-operator/pkg/ironic"
)

func newTestScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	_ = metal3api.AddToScheme(scheme)
	return scheme
}

func newTestReconciler(scheme *runtime.Scheme, fakeClient *fake.ClientBuilder, recorder *events.FakeRecorder) *IronicReconciler {
	cl := fakeClient.Build()
	return &IronicReconciler{
		Client:        cl,
		APIReader:     cl,
		Scheme:        scheme,
		Log:           logr.Discard(),
		EventRecorder: recorder,
	}
}

func newTestIronic() *metal3api.Ironic {
	return &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ironic",
			Namespace: "test-ns",
			UID:       "test-uid-12345",
		},
	}
}

func newTestControllerContext(t *testing.T, scheme *runtime.Scheme, cl client.Client) ironic.ControllerContext {
	t.Helper()

	return ironic.ControllerContext{
		Context: t.Context(),
		Client:  cl,
		Scheme:  scheme,
		Logger:  logr.Discard(),
	}
}

func environmentLabels() map[string]string {
	return map[string]string{
		metal3api.LabelEnvironmentName: metal3api.LabelEnvironmentValue,
	}
}

func drainEvents(recorder *events.FakeRecorder) []string {
	var result []string
	for {
		select {
		case evt := <-recorder.Events:
			result = append(result, evt)
		default:
			return result
		}
	}
}

func TestGetAndUpdateSecret_NotFound_EmitsInvalidLinkedResourceEvent(t *testing.T) {
	scheme := newTestScheme()
	recorder := events.NewFakeRecorder(10)
	ironicObj := newTestIronic()

	r := newTestReconciler(scheme, fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(ironicObj).WithObjects(ironicObj), recorder)
	cctx := newTestControllerContext(t, scheme, r.Client)

	_, requeue, err := r.getAndUpdateSecret(cctx, ironicObj, "missing-secret")

	require.Error(t, err)
	assert.True(t, requeue)
	assert.Contains(t, err.Error(), "cannot load secret")

	evts := drainEvents(recorder)
	require.Len(t, evts, 1)
	assert.Contains(t, evts[0], "InvalidLinkedResource")
	assert.Contains(t, evts[0], "secret test-ns/missing-secret not found")
}

func TestGetAndUpdateSecret_Found_NoEvent(t *testing.T) {
	scheme := newTestScheme()
	recorder := events.NewFakeRecorder(10)
	ironicObj := newTestIronic()

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "existing-secret",
			Namespace: "test-ns",
			Labels:    environmentLabels(),
		},
		Data: map[string][]byte{
			"username": []byte("admin"),
			"password": []byte("pass"),
		},
	}

	r := newTestReconciler(scheme, fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret, ironicObj), recorder)
	cctx := newTestControllerContext(t, scheme, r.Client)

	result, requeue, err := r.getAndUpdateSecret(cctx, ironicObj, "existing-secret")

	require.NoError(t, err)
	assert.False(t, requeue)
	assert.Equal(t, "existing-secret", result.Name)

	evts := drainEvents(recorder)
	assert.Empty(t, evts, "no events should be emitted for a found secret")
}

func TestGetAndUpdateSecret_MissingLabel_EmitsInvalidLinkedResourceEvent(t *testing.T) {
	scheme := newTestScheme()
	recorder := events.NewFakeRecorder(10)
	ironicObj := newTestIronic()

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "unlabeled-secret",
			Namespace: "test-ns",
		},
		Data: map[string][]byte{"data": []byte("value")},
	}

	r := newTestReconciler(scheme, fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(ironicObj).WithObjects(secret, ironicObj), recorder)
	cctx := newTestControllerContext(t, scheme, r.Client)

	_, requeue, err := r.getAndUpdateSecret(cctx, ironicObj, "unlabeled-secret")

	require.Error(t, err)
	assert.True(t, requeue)

	evts := drainEvents(recorder)
	assert.NotEmpty(t, evts, "InvalidLinkedResource should be emitted for label errors")
	foundInvalid := false
	for _, evt := range evts {
		assert.Contains(t, evt, "InvalidLinkedResource", "invalid linked resource reason must be emitted")
		foundInvalid = true
		break
	}
	assert.True(t, foundInvalid, "InvalidLinkedResource event should be present")
}

func TestGetConfigMap_NotFound_EmitsInvalidLinkedResourceEvent(t *testing.T) {
	scheme := newTestScheme()
	recorder := events.NewFakeRecorder(10)
	ironicObj := newTestIronic()

	r := newTestReconciler(scheme, fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(ironicObj).WithObjects(ironicObj), recorder)
	cctx := newTestControllerContext(t, scheme, r.Client)

	_, requeue, err := r.getConfigMap(cctx, ironicObj, "missing-configmap")

	require.Error(t, err)
	assert.True(t, requeue)
	assert.Contains(t, err.Error(), "cannot load configmap")

	evts := drainEvents(recorder)
	require.Len(t, evts, 1)
	assert.Contains(t, evts[0], "InvalidLinkedResource")
	assert.Contains(t, evts[0], "configmap test-ns/missing-configmap not found")
}

func TestGetConfigMap_Found_NoEvent(t *testing.T) {
	scheme := newTestScheme()
	recorder := events.NewFakeRecorder(10)
	ironicObj := newTestIronic()

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "existing-cm",
			Namespace: "test-ns",
			Labels:    environmentLabels(),
		},
		Data: map[string]string{"ca.crt": "cert-data"},
	}

	r := newTestReconciler(scheme, fake.NewClientBuilder().WithScheme(scheme).WithObjects(configMap, ironicObj), recorder)
	cctx := newTestControllerContext(t, scheme, r.Client)

	result, requeue, err := r.getConfigMap(cctx, ironicObj, "existing-cm")

	require.NoError(t, err)
	assert.False(t, requeue)
	assert.Equal(t, "existing-cm", result.Name)

	evts := drainEvents(recorder)
	assert.Empty(t, evts, "no events should be emitted for a found configmap")
}

func TestGetConfigMap_MissingLabel_EmitsInvalidLinkedResourceEvent(t *testing.T) {
	scheme := newTestScheme()
	recorder := events.NewFakeRecorder(10)
	ironicObj := newTestIronic()

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "unlabeled-cm",
			Namespace: "test-ns",
		},
		Data: map[string]string{"data": "value"},
	}

	r := newTestReconciler(scheme, fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(ironicObj).WithObjects(configMap, ironicObj), recorder)
	cctx := newTestControllerContext(t, scheme, r.Client)

	_, requeue, err := r.getConfigMap(cctx, ironicObj, "unlabeled-cm")

	require.Error(t, err)
	assert.True(t, requeue)

	evts := drainEvents(recorder)
	assert.NotEmpty(t, evts, "InvalidLinkedResource should be emitted for label errors")
	foundInvalid := false
	for _, evt := range evts {
		assert.Contains(t, evt, "InvalidLinkedResource", "invalid linked resource reason must be emitted")
		foundInvalid = true
		break
	}
	assert.True(t, foundInvalid, "InvalidLinkedResource event should be present")
}

func TestEnsureAPISecret_GeneratesSecret_EmitsAPISecretCreatedEvent(t *testing.T) {
	scheme := newTestScheme()
	recorder := events.NewFakeRecorder(10)
	ironicObj := newTestIronic()

	r := newTestReconciler(scheme, fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(ironicObj).WithObjects(ironicObj), recorder)
	cctx := newTestControllerContext(t, scheme, r.Client)

	apiSecret, requeue, err := r.ensureAPISecret(cctx, ironicObj)

	require.NoError(t, err)
	assert.True(t, requeue)
	assert.NotNil(t, apiSecret)
	assert.NotEmpty(t, ironicObj.Spec.APICredentialsName)

	evts := drainEvents(recorder)
	require.Len(t, evts, 1)
	assert.Contains(t, evts[0], "APISecretCreated")
	assert.Contains(t, evts[0], "Created new API credentials secret")
}

func TestUpdateIronicStatus_TransitionToReady_EmitsIronicReadyEvent(t *testing.T) {
	scheme := newTestScheme()
	recorder := events.NewFakeRecorder(10)
	ironicObj := newTestIronic()

	r := newTestReconciler(scheme, fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(ironicObj).WithObjects(ironicObj), recorder)
	cctx := newTestControllerContext(t, scheme, r.Client)

	readyStatus := ironic.Status{Ready: true}
	requeue, err := r.updateIronicStatus(cctx, ironicObj, readyStatus, "latest")

	require.NoError(t, err)
	assert.False(t, requeue)

	evts := drainEvents(recorder)
	require.Len(t, evts, 1)
	assert.Contains(t, evts[0], "IronicReady")
	assert.Contains(t, evts[0], "Ironic deployment is now ready")
}

func TestUpdateIronicStatus_TransitionToNotReady_EmitsIronicNotReadyEvent(t *testing.T) {
	scheme := newTestScheme()
	recorder := events.NewFakeRecorder(10)
	ironicObj := newTestIronic()
	ironicObj.Status.Conditions = []metav1.Condition{
		{
			Type:               string(metal3api.IronicStatusReady),
			Status:             metav1.ConditionTrue,
			Reason:             metal3api.IronicReasonAvailable,
			Message:            "ironic: resources are available",
			LastTransitionTime: metav1.Now(),
		},
	}

	r := newTestReconciler(scheme, fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(ironicObj).WithObjects(ironicObj), recorder)
	cctx := newTestControllerContext(t, scheme, r.Client)

	notReadyStatus := ironic.Status{Message: "deployment not available yet"}
	requeue, err := r.updateIronicStatus(cctx, ironicObj, notReadyStatus, "latest")

	require.NoError(t, err)
	assert.False(t, requeue)

	evts := drainEvents(recorder)
	require.Len(t, evts, 1)
	assert.Contains(t, evts[0], "IronicNotReady")
}

func TestUpdateIronicStatus_AlreadyReady_NoEvent(t *testing.T) {
	scheme := newTestScheme()
	recorder := events.NewFakeRecorder(10)
	ironicObj := newTestIronic()
	ironicObj.Status.Conditions = []metav1.Condition{
		{
			Type:               string(metal3api.IronicStatusReady),
			Status:             metav1.ConditionTrue,
			Reason:             metal3api.IronicReasonAvailable,
			Message:            "ironic: resources are available",
			LastTransitionTime: metav1.Now(),
		},
	}
	ironicObj.Status.InstalledVersion = "latest"

	r := newTestReconciler(scheme, fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(ironicObj).WithObjects(ironicObj), recorder)
	cctx := newTestControllerContext(t, scheme, r.Client)

	readyStatus := ironic.Status{Ready: true}
	requeue, err := r.updateIronicStatus(cctx, ironicObj, readyStatus, "latest")

	require.NoError(t, err)
	assert.False(t, requeue)

	evts := drainEvents(recorder)
	assert.Empty(t, evts, "no events should be emitted when status is already ready")
}

func TestUpdateIronicStatus_AlreadyNotReady_NoEvent(t *testing.T) {
	scheme := newTestScheme()
	recorder := events.NewFakeRecorder(10)
	ironicObj := newTestIronic()
	ironicObj.Status.Conditions = []metav1.Condition{
		{
			Type:               string(metal3api.IronicStatusReady),
			Status:             metav1.ConditionFalse,
			Reason:             metal3api.IronicReasonInProgress,
			Message:            "ironic: deployment not available yet",
			LastTransitionTime: metav1.Now(),
		},
	}

	r := newTestReconciler(scheme, fake.NewClientBuilder().WithScheme(scheme).WithStatusSubresource(ironicObj).WithObjects(ironicObj), recorder)
	cctx := newTestControllerContext(t, scheme, r.Client)

	notReadyStatus := ironic.Status{Message: "deployment not available yet"}
	requeue, err := r.updateIronicStatus(cctx, ironicObj, notReadyStatus, "latest")

	require.NoError(t, err)
	assert.False(t, requeue)

	evts := drainEvents(recorder)
	assert.Empty(t, evts, "no events should be emitted when status is already not ready")
}

func TestEnsureAPISecret_ExistingSecret_NoAPISecretCreatedEvent(t *testing.T) {
	scheme := newTestScheme()
	recorder := events.NewFakeRecorder(10)
	ironicObj := newTestIronic()
	ironicObj.Spec.APICredentialsName = "existing-api-secret"

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "existing-api-secret",
			Namespace: "test-ns",
			Labels:    environmentLabels(),
		},
		Data: map[string][]byte{
			"username": []byte("admin"),
			"password": []byte("secret"),
			"htpasswd": []byte("admin:{SHA}W6ph5Mm5Pz8GgiULbPgzG37mj9g="),
		},
	}

	r := newTestReconciler(scheme, fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret, ironicObj), recorder)
	cctx := newTestControllerContext(t, scheme, r.Client)

	result, _, err := r.ensureAPISecret(cctx, ironicObj)

	require.NoError(t, err)
	assert.NotNil(t, result)

	evts := drainEvents(recorder)
	for _, evt := range evts {
		assert.NotContains(t, evt, "APISecretCreated",
			"APISecretCreated must not be emitted when using an existing secret")
	}
}
