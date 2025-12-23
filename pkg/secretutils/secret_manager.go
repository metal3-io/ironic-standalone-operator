package secretutils

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// SecretManager is a type for fetching Secrets whether or not they are in the
// client cache, labelling so that they will be included in the client cache,
// and optionally setting an owner reference.
//
//nolint:containedctx // Context is intentionally stored for use throughout the manager lifecycle
type SecretManager struct {
	ctx       context.Context
	log       logr.Logger
	client    client.Client
	apiReader client.Reader
}

// NewSecretManager returns a new SecretManager.
func NewSecretManager(ctx context.Context, log logr.Logger, cacheClient client.Client, apiReader client.Reader) SecretManager {
	return SecretManager{
		ctx:       ctx,
		log:       log.WithName("secret_manager"),
		client:    cacheClient,
		apiReader: apiReader,
	}
}

func objectType(object client.Object) string {
	switch object.(type) {
	case *corev1.Secret:
		return "secret"
	case *corev1.ConfigMap:
		return "config map"
	}

	return "unknown object"
}

// findObject retrieves an object from the cache if it is available, and from the k8s API if not.
func (sm *SecretManager) findObject(key types.NamespacedName, object client.Object) error {
	// Look for object in the filtered cache
	err := sm.client.Get(sm.ctx, key, object)
	if err == nil {
		return nil
	}
	if !k8serrors.IsNotFound(err) {
		return err
	}

	// Object not in cache; check API directly for unlabelled object
	err = sm.apiReader.Get(sm.ctx, key, object)
	if err != nil {
		return err
	}

	return nil
}

// claimObject ensures that the object has a label that will ensure it is
// present in the cache (and that we can watch for changes), and optionally
// that it has a particular owner reference.
func (sm *SecretManager) claimObject(object client.Object, owner client.Object, scheme *runtime.Scheme) error {
	log := sm.log.WithValues(objectType(object), object.GetName(), "namespace", object.GetNamespace())
	needsUpdate := false

	currentLabels := object.GetLabels()
	if _, found := currentLabels[LabelEnvironmentName]; !found {
		log.Info("setting environment label")
		if currentLabels == nil {
			currentLabels = make(map[string]string, 1)
		}
		currentLabels[LabelEnvironmentName] = LabelEnvironmentValue
		object.SetLabels(currentLabels)
		needsUpdate = true
	}

	if owner != nil && scheme != nil {
		ownerLog := log.WithValues(
			"ownerKind", owner.GetObjectKind().GroupVersionKind().Kind,
			"owner", owner.GetNamespace()+"/"+owner.GetName(),
			"ownerUID", owner.GetUID())

		alreadyOwned := false
		ownerUID := owner.GetUID()
		for _, ref := range object.GetOwnerReferences() {
			if ref.UID == ownerUID {
				alreadyOwned = true
				break
			}
		}
		if !alreadyOwned {
			ownerLog.Info("setting owner reference")
			if err := controllerutil.SetOwnerReference(owner, object, scheme); err != nil {
				return fmt.Errorf("failed to set %s owner reference: %w", objectType(object), err)
			}
			needsUpdate = true
		}
	}

	if needsUpdate {
		if err := sm.client.Update(sm.ctx, object); err != nil {
			return fmt.Errorf("failed to update %s %s/%s: %w", objectType(object), object.GetNamespace(), object.GetName(), err)
		}
	}

	return nil
}

// AcquireSecret retrieves a Secret and ensures that it has a label that will
// ensure it is present in the cache (and that we can watch for changes), and
// that it has a particular owner reference.
func (sm *SecretManager) AcquireSecret(key types.NamespacedName, owner client.Object, scheme *runtime.Scheme) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := sm.findObject(key, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch secret %s in namespace %s: %w", key.Name, key.Namespace, err)
	}
	err = sm.claimObject(secret, owner, scheme)
	return secret, err
}

// ObtainSecret retrieves a Secret and ensures that it has a label that will
// ensure it is present in the cache (and that we can watch for changes).
// This version does not set owner references.
func (sm *SecretManager) ObtainSecret(key types.NamespacedName) (*corev1.Secret, error) {
	return sm.AcquireSecret(key, nil, nil)
}

// AcquireConfigMap retrieves a ConfigMap and ensures that it has a label that will
// ensure it is present in the cache (and that we can watch for changes), and
// that it has a particular owner reference.
func (sm *SecretManager) AcquireConfigMap(key types.NamespacedName, owner client.Object, scheme *runtime.Scheme) (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{}
	err := sm.findObject(key, configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch config map %s in namespace %s: %w", key.Name, key.Namespace, err)
	}
	err = sm.claimObject(configMap, owner, scheme)
	return configMap, err
}

// ObtainConfigMap retrieves a ConfigMap and ensures that it has a label that will
// ensure it is present in the cache (and that we can watch for changes).
// This version does not set owner references.
func (sm *SecretManager) ObtainConfigMap(key types.NamespacedName) (*corev1.ConfigMap, error) {
	return sm.AcquireConfigMap(key, nil, nil)
}
