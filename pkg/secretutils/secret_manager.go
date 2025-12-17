package secretutils

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// findSecret retrieves a Secret from the cache if it is available, and from the
// k8s API if not.
func (sm *SecretManager) findSecret(key types.NamespacedName) (secret *corev1.Secret, err error) {
	secret = &corev1.Secret{}

	// Look for secret in the filtered cache
	err = sm.client.Get(sm.ctx, key, secret)
	if err == nil {
		return secret, nil
	}
	if !k8serrors.IsNotFound(err) {
		return nil, err
	}

	// Secret not in cache; check API directly for unlabelled Secret
	err = sm.apiReader.Get(sm.ctx, key, secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// claimSecret ensures that the Secret has a label that will ensure it is
// present in the cache (and that we can watch for changes), and optionally
// that it has a particular owner reference.
func (sm *SecretManager) claimSecret(secret *corev1.Secret, owner client.Object, scheme *runtime.Scheme) error {
	log := sm.log.WithValues("secret", secret.Name, "secretNamespace", secret.Namespace)
	needsUpdate := false

	if !metav1.HasLabel(secret.ObjectMeta, LabelEnvironmentName) {
		log.Info("setting secret environment label")
		metav1.SetMetaDataLabel(&secret.ObjectMeta, LabelEnvironmentName, LabelEnvironmentValue)
		needsUpdate = true
	}

	if owner != nil && scheme != nil {
		ownerLog := log.WithValues(
			"ownerKind", owner.GetObjectKind().GroupVersionKind().Kind,
			"owner", owner.GetNamespace()+"/"+owner.GetName(),
			"ownerUID", owner.GetUID())

		alreadyOwned := false
		ownerUID := owner.GetUID()
		for _, ref := range secret.GetOwnerReferences() {
			if ref.UID == ownerUID {
				alreadyOwned = true
				break
			}
		}
		if !alreadyOwned {
			ownerLog.Info("setting secret owner reference")
			if err := controllerutil.SetOwnerReference(owner, secret, scheme); err != nil {
				return fmt.Errorf("failed to set secret owner reference: %w", err)
			}
			needsUpdate = true
		}
	}

	if needsUpdate {
		if err := sm.client.Update(sm.ctx, secret); err != nil {
			return fmt.Errorf("failed to update secret %s in namespace %s: %w", secret.Name, secret.Namespace, err)
		}
	}

	return nil
}

// AcquireSecret retrieves a Secret and ensures that it has a label that will
// ensure it is present in the cache (and that we can watch for changes), and
// that it has a particular owner reference.
func (sm *SecretManager) AcquireSecret(key types.NamespacedName, owner client.Object, scheme *runtime.Scheme) (*corev1.Secret, error) {
	secret, err := sm.findSecret(key)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch secret %s in namespace %s: %w", key.Name, key.Namespace, err)
	}
	err = sm.claimSecret(secret, owner, scheme)
	return secret, err
}

// ObtainSecret retrieves a Secret and ensures that it has a label that will
// ensure it is present in the cache (and that we can watch for changes).
// This version does not set owner references.
func (sm *SecretManager) ObtainSecret(key types.NamespacedName) (*corev1.Secret, error) {
	return sm.AcquireSecret(key, nil, nil)
}
