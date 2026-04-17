package ironic

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

const (
	protoFile = "file"
	protoOCI  = "oci"
)

const (
	minVLAN = 1
	maxVLAN = 4094
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
	urlStr = strings.TrimSpace(urlStr)
	if urlStr == "" {
		return fmt.Errorf("%s: URL must not be empty", fieldName)
	}

	parsed, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("%s: invalid URL format: %w", fieldName, err)
	}

	switch parsed.Scheme {
	case protoFile:
		if parsed.Host != "" {
			return fmt.Errorf("%s: file URL must use an absolute path (file:///...)", fieldName)
		}
		if parsed.Path == "" || parsed.Path[0] != '/' {
			return fmt.Errorf("%s: file URL must use an absolute path (file:///...)", fieldName)
		}
	case protoHTTP, protoHTTPS:
		if parsed.Host == "" {
			return fmt.Errorf("%s: %s URL must include a host", fieldName, parsed.Scheme)
		}
	case protoOCI:
		if parsed.Host == "" {
			return fmt.Errorf("%s: oci URL must include a registry host", fieldName)
		}
	default:
		return fmt.Errorf("%s: unsupported protocol %q (must be file://, http://, https://, or oci://)", fieldName, parsed.Scheme)
	}

	return nil
}

func validateAgentImages(images []metal3api.AgentImages) error {
	if len(images) == 0 {
		return nil
	}

	seenArchitectures := make(map[metal3api.CPUArchitecture]bool)

	for i, img := range images {
		if strings.TrimSpace(img.Kernel) == "" {
			return fmt.Errorf("overrides.agentImages[%d]: kernel is required", i)
		}
		if strings.TrimSpace(img.Initramfs) == "" {
			return fmt.Errorf("overrides.agentImages[%d]: initramfs is required", i)
		}

		if err := validateURL(img.Kernel, fmt.Sprintf("overrides.agentImages[%d].kernel", i)); err != nil {
			return err
		}

		if err := validateURL(img.Initramfs, fmt.Sprintf("overrides.agentImages[%d].initramfs", i)); err != nil {
			return err
		}

		if seenArchitectures[img.Architecture] {
			if img.Architecture == "" {
				return errors.New("overrides.agentImages: duplicate default (empty architecture) entry")
			}
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

	if ironic.PrometheusExporter != nil && ironic.PrometheusExporter.BindAddress != "" {
		ip := net.ParseIP(ironic.PrometheusExporter.BindAddress)
		if ip == nil {
			return fmt.Errorf("prometheusExporter: bindAddress %q is not a valid IP address", ironic.PrometheusExporter.BindAddress)
		}
	}

	if ironic.PrometheusExporter != nil && ironic.PrometheusExporter.Enabled && !ironic.PrometheusExporter.DisableServiceMonitor {
		bindAddr := ironic.PrometheusExporter.BindAddress
		if bindAddr == "" {
			bindAddr = defaultMetricsBindAddr
		}
		ip := net.ParseIP(bindAddr)
		if ip != nil && ip.IsLoopback() {
			return fmt.Errorf("ServiceMonitor is not compatible with a loopback bindAddress %q, since the metrics endpoint is not reachable from remote Prometheus instances; set a non-loopback bindAddress or set disableServiceMonitor to true", bindAddr)
		}
	}

	if ironic.HighAvailability && !metal3api.CurrentFeatureGate.Enabled(metal3api.FeatureHighAvailability) {
		return errors.New("highly available architecture is disabled via feature gate")
	}

	if !ironic.HighAvailability && ironic.TLS.InsecureRPC != nil {
		return errors.New("insecureRPC makes no sense without highAvailability")
	}

	// Validate TLS CA settings
	if err := validateCASettings(&ironic.TLS); err != nil {
		return err
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

	if ironic.NetworkingService != nil && ironic.NetworkingService.Enabled {
		if err := validateNetworkingService(ironic); err != nil {
			return err
		}
	}

	return nil
}

func validateNetworkingService(ironic *metal3api.IronicSpec) error {
	// Validate provider network configs if provided
	seenTypes := make(map[metal3api.ProviderNetworkType]bool)
	for i := range ironic.NetworkingService.ProviderNetworks {
		pn := &ironic.NetworkingService.ProviderNetworks[i]
		if seenTypes[pn.Type] {
			return fmt.Errorf("duplicate provider network type %q", pn.Type)
		}
		seenTypes[pn.Type] = true
		if err := validateProviderNetwork(pn); err != nil {
			return err
		}
	}

	return nil
}

func validateProviderNetwork(sn *metal3api.ProviderNetworkConfig) error {
	// Validate mode-specific requirements
	switch sn.Mode {
	case metal3api.SwitchportModeAccess:
		if len(sn.AllowedVLANs) > 0 {
			return errors.New("allowedVLANs cannot be set in access mode")
		}
	case metal3api.SwitchportModeTrunk, metal3api.SwitchportModeHybrid:
		if len(sn.AllowedVLANs) == 0 {
			return fmt.Errorf("allowedVLANs required for %s mode", sn.Mode)
		}
	default:
		return fmt.Errorf("invalid switchport mode: %s", sn.Mode)
	}

	for _, entry := range sn.AllowedVLANs {
		if err := validateAllowedVLANEntry(entry); err != nil {
			return err
		}
	}

	return nil
}

func validateVLANID(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%q is not a valid VLAN ID", s)
	}
	if v < minVLAN || v > maxVLAN {
		return 0, fmt.Errorf("VLAN ID %d is out of range (%d-%d)", v, minVLAN, maxVLAN)
	}
	return v, nil
}

func validateAllowedVLANEntry(entry string) error {
	parts := strings.SplitN(entry, "-", 2)
	if len(parts) == 1 {
		_, err := validateVLANID(strings.TrimSpace(parts[0]))
		return err
	}
	start, err := validateVLANID(strings.TrimSpace(parts[0]))
	if err != nil {
		return fmt.Errorf("invalid range %q: %w", entry, err)
	}
	end, err := validateVLANID(strings.TrimSpace(parts[1]))
	if err != nil {
		return fmt.Errorf("invalid range %q: %w", entry, err)
	}
	if start >= end {
		return fmt.Errorf("invalid range %q: start (%d) must be less than end (%d)", entry, start, end)
	}
	return nil
}

func validateCASettings(tls *metal3api.TLS) error {
	// Validate BMCCA
	if tls.BMCCA != nil {
		if tls.BMCCA.Name == "" {
			return errors.New("tls.bmcCA.name is required when tls.bmcCA is set")
		}
		// Both old and new fields are set - validate they're consistent
		if tls.BMCCAName != "" && (tls.BMCCA.Kind != metal3api.ResourceKindSecret || tls.BMCCA.Name != tls.BMCCAName) {
			return errors.New("tls.bmcCA and tls.bmcCAName are both set but inconsistent; use tls.bmcCA only")
		}
	}

	// Validate TrustedCA
	if tls.TrustedCA != nil {
		if tls.TrustedCA.Name == "" {
			return errors.New("tls.trustedCA.name is required when tls.trustedCA is set")
		}
		// Both old and new fields are set - validate they're consistent
		if tls.TrustedCAName != "" && (tls.TrustedCA.Kind != metal3api.ResourceKindConfigMap || tls.TrustedCA.Name != tls.TrustedCAName) {
			return errors.New("tls.trustedCA and tls.trustedCAName are both set but inconsistent; use tls.trustedCA only")
		}
	}

	return nil
}

// Validate all resources before using them. This method is a superset of
// ValidateIronic with validations that require access to linked resources.
func (resources *Resources) Validate() error {
	if err := ValidateIronic(&resources.Ironic.Spec, nil); err != nil {
		return err
	}

	if resources.Ironic.Spec.TLS.TrustedCA != nil {
		key := resources.Ironic.Spec.TLS.TrustedCA.Key
		if key != "" && !resources.hasTrustedCAKey(key) {
			return fmt.Errorf("resources referenced in tls.trustedCA does not contain the required key %s", key)
		}
	}

	return nil
}

func (resources *Resources) hasTrustedCAKey(key string) bool {
	switch {
	case resources.TrustedCASecret != nil:
		for found := range resources.TrustedCASecret.Data {
			if found == key {
				return true
			}
		}
	case resources.TrustedCAConfigMap != nil:
		for found := range resources.TrustedCAConfigMap.Data {
			if found == key {
				return true
			}
		}
	default:
		// Cannot happen in reality but just in case
		return true
	}

	return false
}
