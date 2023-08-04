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
	"errors"
	"fmt"
	"net"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
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

func setDefaults(ironic *IronicSpec) {}

//+kubebuilder:webhook:path=/validate-metal3-io-v1alpha1-ironic,mutating=false,failurePolicy=fail,sideEffects=None,groups=metal3.io,resources=ironics,verbs=create;update,versions=v1alpha1,name=vironic.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Ironic{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Ironic) ValidateCreate() (warnings admission.Warnings, err error) {
	ironiclog.Info("validate create", "name", r.Name)
	return nil, validateIronic(&r.Spec, nil)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Ironic) ValidateUpdate(old runtime.Object) (warnings admission.Warnings, err error) {
	ironiclog.Info("validate update", "name", r.Name)
	return nil, validateIronic(&r.Spec, &old.(*Ironic).Spec)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Ironic) ValidateDelete() (warnings admission.Warnings, err error) {
	return nil, nil
}

func validateIronic(ironic *IronicSpec, old *IronicSpec) error {
	if ironic.APISecretName == "" {
		return errors.New("apiSecretName is required")
	}

	if ironic.Distributed && ironic.DatabaseName == "" {
		return errors.New("database is required for distributed architecture")
	}

	if old != nil && old.DatabaseName != "" && old.DatabaseName != ironic.DatabaseName {
		return errors.New("cannot change to a new database or remove it")
	}

	if ironic.Networking.IPAddress != "" && net.ParseIP(ironic.Networking.IPAddress) == nil {
		return fmt.Errorf("%s is not a valid IP address", ironic.Networking.IPAddress)
	}

	// TODO(dtantsur): implement and remove (comment out for local testing)
	if ironic.Distributed {
		return errors.New("distributed architecture is experimental, please do not use")
	}

	return nil
}
