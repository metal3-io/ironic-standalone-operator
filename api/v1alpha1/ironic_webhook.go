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
	"net/netip"

	"go4.org/netipx"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var ironiclog = logf.Log.WithName("webhooks").WithName("Ironic")

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

func SetDHCPDefaults(dhcp *DHCP) {
	provCIDR, err := netip.ParsePrefix(dhcp.NetworkCIDR)
	if err != nil {
		// Let the validation hook do the actual validation
		return
	}

	provIP := provCIDR.Addr()
	if dhcp.RangeBegin == "" {
		firstIP := provIP
		for i := 0; i < 10; i++ {
			firstIP = firstIP.Next()
		}
		if firstIP.IsValid() && provCIDR.Contains(firstIP) {
			dhcp.RangeBegin = firstIP.String()
		}
	}
	if dhcp.RangeEnd == "" {
		lastIP := netipx.PrefixLastIP(provCIDR).Prev().Prev()
		if lastIP.IsValid() && provCIDR.Contains(lastIP) {
			dhcp.RangeEnd = lastIP.String()
		}
	}
}

func setDefaults(ironic *IronicSpec) {
	if dhcp := ironic.Networking.DHCP; dhcp != nil {
		SetDHCPDefaults(dhcp)
	}
}

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

func validateIP(ip string) error {
	if ip == "" {
		return nil
	}

	if _, err := netip.ParseAddr(ip); err != nil {
		return fmt.Errorf("%s is not a valid IP address: %w", ip, err)
	}

	return nil
}

func validateIPinPrefix(ip string, prefix netip.Prefix) error {
	if ip == "" {
		return nil
	}

	parsed, err := netip.ParseAddr(ip)
	if err != nil {
		return fmt.Errorf("%s is not a valid IP address: %w", ip, err)
	}

	if !prefix.Contains(parsed) {
		return fmt.Errorf("%s is not in networking.dhcp.networkCIDR", ip)
	}

	return nil
}

func ValidateDHCP(ironic *IronicSpec, dhcp *DHCP) error {
	hasNetworking := ironic.Networking.IPAddress != "" || ironic.Networking.Interface != "" || len(ironic.Networking.MACAddresses) > 0
	if !hasNetworking {
		return errors.New("networking: at least one of ipAddress, interface or macAddresses is required when DHCP is used")
	}
	if dhcp.NetworkCIDR == "" {
		return errors.New("networking.dhcp.networkCIRD is required when DHCP is used")
	}
	if dhcp.ServeDNS && dhcp.DNSAddress != "" {
		return errors.New("networking.dhcp.dnsAddress cannot set together with serveDNS")
	}

	provIP, _ := netip.ParseAddr(ironic.Networking.IPAddress)
	provCIDR, err := netip.ParsePrefix(dhcp.NetworkCIDR)
	if err != nil {
		return fmt.Errorf("networking.dhcp.networkCIDR is invalid: %w", err)
	}

	if !provCIDR.Contains(provIP) {
		return errors.New("networking.dhcp.networkCIDR must contain networking.ipAddress")
	}

	if err := validateIPinPrefix(dhcp.RangeBegin, provCIDR); err != nil {
		return err
	}

	if err := validateIPinPrefix(dhcp.RangeEnd, provCIDR); err != nil {
		return err
	}

	if err := validateIP(dhcp.DNSAddress); err != nil {
		return err
	}

	if err := validateIP(dhcp.GatewayAddress); err != nil {
		return err
	}

	// These are supposed to be populated by the webhook
	if dhcp.RangeBegin == "" || dhcp.RangeEnd == "" {
		return errors.New("firstIP and lastIP are not set and could not be automatically populated")
	}

	return nil
}

func validateIronic(ironic *IronicSpec, old *IronicSpec) error {
	if ironic.HighAvailability && ironic.DatabaseRef.Name == "" {
		return errors.New("database is required for highly available architecture")
	}

	if old != nil && old.DatabaseRef.Name != "" && old.DatabaseRef.Name != ironic.DatabaseRef.Name {
		return errors.New("cannot change to a new database or remove it")
	}

	if err := validateIP(ironic.Networking.IPAddress); err != nil {
		return err
	}

	if err := validateIP(ironic.Networking.ExternalIP); err != nil {
		return err
	}

	if ironic.HighAvailability && ironic.Networking.IPAddress != "" {
		return errors.New("networking.ipAddress makes no sense with highly available architecture")
	}

	if dhcp := ironic.Networking.DHCP; dhcp != nil {
		if err := ValidateDHCP(ironic, dhcp); err != nil {
			return err
		}
	}

	// TODO(dtantsur): implement and remove (comment out for local testing)
	if ironic.HighAvailability {
		return errors.New("highly available architecture is experimental, please do not use")
	}

	return nil
}
