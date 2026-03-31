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

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

// MissingLabelError is returned when a resource does not have the required
// environment label. This error requires user intervention to add the label
// to the resource before the operator will use it.
type MissingLabelError struct {
	ObjectType string
	Namespace  string
	Name       string
}

func (e *MissingLabelError) Error() string {
	return fmt.Sprintf(
		"%s %s/%s does not have the required label %s=%s",
		e.ObjectType, e.Namespace, e.Name,
		metal3api.LabelEnvironmentName, metal3api.LabelEnvironmentValue)
}

// SecretManager is a type for fetching Secrets whether or not they are in the
// client cache, verifying that they carry the required
// environment label, and optionally setting an owner reference.
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

// claimSecret ensures that the Secret has the required environment label
// (which must be set by the user for user-provided resources, or by the
// operator for resources it creates), and optionally sets an owner reference.
func (sm *SecretManager) claimSecret(secret *corev1.Secret, owner client.Object, scheme *runtime.Scheme) error {
	log := sm.log.WithValues("secret", secret.Name, "secretNamespace", secret.Namespace)
	needsUpdate := false

	// Require the environment label to already be present. For user-provided
	// resources this acts as an opt-in: the user must label the resource before
	// the operator will use it. Operator-created resources are labelled at
	// creation time.
	currentLabels := secret.GetLabels()
	if currentLabels[metal3api.LabelEnvironmentName] != metal3api.LabelEnvironmentValue {
		return &MissingLabelError{
			ObjectType: "secret",
			Namespace:  secret.Namespace,
			Name:       secret.Name,
		}
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

// AcquireSecret retrieves a Secret, verifies it carries the required
// environment label, and ensures it has a particular owner reference.
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
