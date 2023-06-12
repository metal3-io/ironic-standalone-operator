/*
Copyright 2023.

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

package v1alpha1

import (
	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var ironiclog = logf.Log.WithName("ironic-resource")

func (r *Ironic) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-metal3-io-v1alpha1-ironic,mutating=true,failurePolicy=fail,sideEffects=None,groups=metal3.io,resources=ironics,verbs=create;update,versions=v1alpha1,name=mironic.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Ironic{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Ironic) Default() {
	ironiclog.Info("default", "name", r.Name)
	setDefaults(&r.Spec)
}

func setDefaults(ironic *IronicSpec) {
	if ironic.APIPort == 0 {
		ironic.APIPort = DefaultAPIPort
	}

	if ironic.ImageServerPort == 0 {
		ironic.ImageServerPort = DefaultImageServerPort
	}

	if ironic.ImageServerTLSPort == 0 && !ironic.DisableVirtualMediaTLS {
		ironic.ImageServerTLSPort = DefaultImageServerTLSPort
	}

	if ironic.Image == "" {
		ironic.Image = DefaultIronicImage
	}

	if ironic.Size > 1 && ironic.Database == nil {
		ironic.Database = new(Database)
	}

	if ironic.Database != nil && ironic.Database.Image == "" {
		ironic.Database.Image = DefaultDatabaseImage
	}
}

//+kubebuilder:webhook:path=/validate-metal3-io-v1alpha1-ironic,mutating=false,failurePolicy=fail,sideEffects=None,groups=metal3.io,resources=ironics,verbs=create;update,versions=v1alpha1,name=vironic.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Ironic{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Ironic) ValidateCreate() error {
	ironiclog.Info("validate create", "name", r.Name)
	return validate(&r.Spec)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Ironic) ValidateUpdate(old runtime.Object) error {
	ironiclog.Info("validate update", "name", r.Name)
	return validate(&r.Spec)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Ironic) ValidateDelete() error {
	return nil
}

func validate(ironic *IronicSpec) error {
	if ironic.Database != nil && ironic.Database.ExternalIP != "" && ironic.Database.CredentialsSecretName == "" {
		return errors.New("external database requires credentials")
	}

	return nil
}
