package ironic

import (
	"errors"
	"fmt"
	"net/netip"
	"reflect"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

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

func ValidateDHCP(ironic *metal3api.IronicSpec) error {
	dhcp := ironic.Networking.DHCP
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
	if dhcp.RangeBegin == "" || dhcp.RangeEnd == "" {
		return errors.New("networking.dhcp: rangeBegin and rangeEnd are required")
	}

	provCIDR, err := netip.ParsePrefix(dhcp.NetworkCIDR)
	if err != nil {
		return fmt.Errorf("networking.dhcp.networkCIDR is invalid: %w", err)
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

	// Check that the provisioning IP is in the CIDR
	if ironic.Networking.IPAddress != "" {
		provIP, _ := netip.ParseAddr(ironic.Networking.IPAddress)
		if !provCIDR.Contains(provIP) {
			return errors.New("networking.dhcp.networkCIDR must contain networking.ipAddress")
		}
	}

	return nil
}

func ValidateIronic(ironic *metal3api.IronicSpec, old *metal3api.IronicSpec) error {
	if ironic.HighAvailability && ironic.Database == nil && ironic.DatabaseName == "" {
		return errors.New("database is required for highly available architecture")
	}

	if old != nil && old.Database != nil && ironic.Database != nil && !reflect.DeepEqual(old.Database, ironic.Database) {
		return errors.New("cannot change to a new database")
	}

	if old != nil && old.DatabaseName != "" && ironic.DatabaseName != "" && old.DatabaseName != ironic.DatabaseName {
		return errors.New("cannot change to a new database")
	}

	if ironic.Database != nil {
		if ironic.DatabaseName != "" {
			return errors.New("databaseName and database cannot be used together")
		}
		if ironic.Database.CredentialsName == "" || ironic.Database.Host == "" || ironic.Database.Name == "" {
			return errors.New("credentialsName, host and name are required on database")
		}
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

	if ironic.Networking.DHCP != nil {
		if ironic.HighAvailability {
			return errors.New("DHCP support is not implemented in the highly available architecture")
		}

		if err := ValidateDHCP(ironic); err != nil {
			return err
		}
	}

	if ironic.Networking.IPAddressManager == metal3api.IPAddressManagerKeepalived {
		if ironic.HighAvailability {
			return errors.New("networking: keepalived is not compatible with the highly available architecture")
		}
		if ironic.Networking.IPAddress == "" || ironic.Networking.Interface == "" {
			return errors.New("networking: keepalived requires specifying both ipAddress and interface")
		}
	}

	if ironic.HighAvailability && !metal3api.CurrentFeatureGate.Enabled(metal3api.FeatureHighAvailability) {
		return errors.New("highly available architecture is disabled via feature gate")
	}

	if ironic.Version != "" {
		if err := metal3api.ValidateVersion(ironic.Version); err != nil {
			return err
		}
	}

	return nil
}
