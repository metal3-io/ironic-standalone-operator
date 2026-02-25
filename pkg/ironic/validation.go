package ironic

import (
	"errors"
	"fmt"
	"net/netip"
	"net/url"
	"reflect"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

const (
	protoFile = "file"
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

func validateURL(urlStr string, fieldName string) error {
	parsed, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("%s: invalid URL format: %w", fieldName, err)
	}

	// Check for supported protocols
	switch parsed.Scheme {
	case protoFile, protoHTTP, protoHTTPS:
		return nil
	default:
		return fmt.Errorf("%s: unsupported protocol %q (must be file://, http://, or https://)", fieldName, parsed.Scheme)
	}
}

func validateAgentImages(images []metal3api.AgentImages) error {
	if len(images) == 0 {
		return nil
	}

	seenArchitectures := make(map[metal3api.CPUArchitecture]bool)

	for i, img := range images {
		if img.Architecture == "" {
			return fmt.Errorf("overrides.agentImages[%d]: architecture is required", i)
		}
		if img.Kernel == "" {
			return fmt.Errorf("overrides.agentImages[%d]: kernel is required", i)
		}
		if img.Initramfs == "" {
			return fmt.Errorf("overrides.agentImages[%d]: initramfs is required", i)
		}

		// Validate URL format and protocol for kernel
		if err := validateURL(img.Kernel, fmt.Sprintf("overrides.agentImages[%d].kernel", i)); err != nil {
			return err
		}

		// Validate URL format and protocol for initramfs
		if err := validateURL(img.Initramfs, fmt.Sprintf("overrides.agentImages[%d].initramfs", i)); err != nil {
			return err
		}

		if seenArchitectures[img.Architecture] {
			return fmt.Errorf("overrides.agentImages: duplicate architecture %q", img.Architecture)
		}
		seenArchitectures[img.Architecture] = true
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
		return errors.New("networking.dhcp.networkCIDR is required when DHCP is used")
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
	if ironic.HighAvailability && ironic.Database == nil {
		return errors.New("database is required for highly available architecture")
	}

	if old != nil && old.Database != nil && ironic.Database != nil && !reflect.DeepEqual(old.Database, ironic.Database) {
		return errors.New("cannot change to a new database")
	}

	if ironic.Database != nil && (ironic.Database.CredentialsName == "" || ironic.Database.Host == "" || ironic.Database.Name == "") {
		return errors.New("credentialsName, host and name are required on database")
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

	if ironic.HighAvailability && ironic.PrometheusExporter != nil && !ironic.PrometheusExporter.DisableServiceMonitor {
		return errors.New("ServiceMonitor support is currently incompatible with the highly available architecture")
	}

	if ironic.HighAvailability && !metal3api.CurrentFeatureGate.Enabled(metal3api.FeatureHighAvailability) {
		return errors.New("highly available architecture is disabled via feature gate")
	}

	if !ironic.HighAvailability && ironic.TLS.InsecureRPC != nil {
		return errors.New("insecureRPC makes no sense without highAvailability")
	}

	if ironic.Version != "" {
		if err := metal3api.ValidateVersion(ironic.Version); err != nil {
			return err
		}
	}

	if ironic.Overrides != nil {
		if err := validateAgentImages(ironic.Overrides.AgentImages); err != nil {
			return err
		}
	}

	return nil
}
