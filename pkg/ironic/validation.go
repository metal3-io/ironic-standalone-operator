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

func validateIPinPrefix(ip string, prefix netip.Prefix, cidrField string) error {
	if ip == "" {
		return nil
	}

	parsed, err := netip.ParseAddr(ip)
	if err != nil {
		return fmt.Errorf("%s is not a valid IP address: %w", ip, err)
	}

	if !prefix.Contains(parsed) {
		return fmt.Errorf("%s is not in %s", ip, cidrField)
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

// validateRangeOrder rejects reversed DHCP ranges, which dnsmasq refuses at
// startup. Both IPs are expected to be already validated.
func validateRangeOrder(beginStr, endStr string) error {
	begin, beginErr := netip.ParseAddr(beginStr)
	end, endErr := netip.ParseAddr(endStr)
	if beginErr == nil && endErr == nil && begin.Compare(end) > 0 {
		return errors.New("rangeBegin must not be after rangeEnd")
	}
	return nil
}

func validateDHCPRange(r metal3api.DHCPRange, idx int) error {
	prefix := fmt.Sprintf("networking.dhcp.extraRanges[%d]", idx)

	if r.NetworkCIDR == "" {
		return fmt.Errorf("%s.networkCIDR is required", prefix)
	}

	cidr, err := netip.ParsePrefix(r.NetworkCIDR)
	if err != nil {
		return fmt.Errorf("%s.networkCIDR is invalid: %w", prefix, err)
	}

	if cidr.Bits() == 0 {
		return fmt.Errorf("%s.networkCIDR must have a non-zero prefix length", prefix)
	}

	if r.RangeBegin == "" || r.RangeEnd == "" {
		return fmt.Errorf("%s: rangeBegin and rangeEnd are required", prefix)
	}

	cidrField := prefix + ".networkCIDR"

	if err := validateIPinPrefix(r.RangeBegin, cidr, cidrField); err != nil {
		return fmt.Errorf("%s.rangeBegin: %w", prefix, err)
	}

	if err := validateIPinPrefix(r.RangeEnd, cidr, cidrField); err != nil {
		return fmt.Errorf("%s.rangeEnd: %w", prefix, err)
	}

	if err := validateRangeOrder(r.RangeBegin, r.RangeEnd); err != nil {
		return fmt.Errorf("%s: %w", prefix, err)
	}

	if r.GatewayAddress != "" {
		if cidr.Addr().Is6() {
			return fmt.Errorf("%s.gatewayAddress: IPv6 per-range gateway is not supported", prefix)
		}
		if err := validateIPinPrefix(r.GatewayAddress, cidr, cidrField); err != nil {
			return fmt.Errorf("%s.gatewayAddress: %w", prefix, err)
		}
	}

	return nil
}

func ValidateDHCP(ironic *metal3api.IronicSpec) error {
	dhcp := ironic.Networking.DHCP
	hasNetworking := ironic.Networking.IPAddress != "" || ironic.Networking.Interface != "" || len(ironic.Networking.MACAddresses) > 0
	if !hasNetworking {
		return errors.New("networking: at least one of ipAddress, interface or macAddresses is required when DHCP is used")
	}
	if dhcp.ServeDNS && dhcp.DNSAddress != "" {
		return errors.New("networking.dhcp.dnsAddress cannot set together with serveDNS")
	}

	if err := validateIP(dhcp.DNSAddress); err != nil {
		return err
	}

	if err := validateIP(dhcp.GatewayAddress); err != nil {
		return err
	}

	// The main range fields are all-or-nothing: leaving networkCIDR,
	// rangeBegin and rangeEnd unset disables the main range so only
	// extraRanges are served (e.g. a relay-only deployment). At least one
	// range - main or extra - must be configured.
	mainFieldsSet := 0
	for _, f := range []string{dhcp.NetworkCIDR, dhcp.RangeBegin, dhcp.RangeEnd} {
		if f != "" {
			mainFieldsSet++
		}
	}
	hasMainRange := mainFieldsSet == 3

	switch {
	case mainFieldsSet != 0 && !hasMainRange:
		return errors.New("networking.dhcp: networkCIDR, rangeBegin and rangeEnd must be set together")
	case !hasMainRange && len(dhcp.ExtraRanges) == 0:
		return errors.New("networking.dhcp: networkCIDR, rangeBegin and rangeEnd are required unless extraRanges is set")
	}

	if hasMainRange {
		provCIDR, err := netip.ParsePrefix(dhcp.NetworkCIDR)
		if err != nil {
			return fmt.Errorf("networking.dhcp.networkCIDR is invalid: %w", err)
		}

		if err := validateIPinPrefix(dhcp.RangeBegin, provCIDR, "networking.dhcp.networkCIDR"); err != nil {
			return err
		}

		if err := validateIPinPrefix(dhcp.RangeEnd, provCIDR, "networking.dhcp.networkCIDR"); err != nil {
			return err
		}

		if err := validateRangeOrder(dhcp.RangeBegin, dhcp.RangeEnd); err != nil {
			return fmt.Errorf("networking.dhcp: %w", err)
		}

		// The main range is a direct-attached subnet by definition, so the
		// provisioning IP must live in it. Subnets reached via a DHCP relay
		// belong in extraRanges.
		if ironic.Networking.IPAddress != "" {
			provIP, _ := netip.ParseAddr(ironic.Networking.IPAddress)
			if !provCIDR.Contains(provIP) {
				return errors.New("networking.dhcp.networkCIDR must contain networking.ipAddress")
			}
		}
	}

	for i, r := range dhcp.ExtraRanges {
		if err := validateDHCPRange(r, i); err != nil {
			return err
		}
	}

	return validateNoPoolOverlap(dhcp)
}

// validateNoPoolOverlap rejects DHCP address pools (main and extra ranges)
// that overlap each other: dnsmasq either refuses to start or allocates
// leases ambiguously between the overlapping pools.
func validateNoPoolOverlap(dhcp *metal3api.DHCP) error {
	type pool struct {
		begin, end netip.Addr
		field      string
	}

	pools := make([]pool, 0, len(dhcp.ExtraRanges)+1)
	addPool := func(beginStr, endStr, field string) {
		begin, err := netip.ParseAddr(beginStr)
		if err != nil {
			return
		}
		end, err := netip.ParseAddr(endStr)
		if err != nil {
			return
		}
		pools = append(pools, pool{begin: begin, end: end, field: field})
	}

	addPool(dhcp.RangeBegin, dhcp.RangeEnd, "networking.dhcp")
	for i, r := range dhcp.ExtraRanges {
		addPool(r.RangeBegin, r.RangeEnd, fmt.Sprintf("networking.dhcp.extraRanges[%d]", i))
	}

	for i := 1; i < len(pools); i++ {
		for j := range i {
			if pools[i].begin.Compare(pools[j].end) <= 0 && pools[j].begin.Compare(pools[i].end) <= 0 {
				return fmt.Errorf("DHCP pools of %s and %s overlap", pools[j].field, pools[i].field)
			}
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

	if ironic.Networking.Ingress != nil && ironic.Networking.ExternalIP != "" {
		return errors.New("networking.ingress and networking.externalIP cannot be set at the same time")
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

	if ironic.Networking.IPAddressManager == metal3api.IPAddressManagerKeepalived { //nolint:staticcheck // backward compat
		if ironic.HighAvailability {
			return errors.New("networking: keepalived is not compatible with the highly available architecture")
		}
		if ironic.Networking.IPAddress == "" || ironic.Networking.Interface == "" {
			return errors.New("networking: keepalived requires specifying both ipAddress and interface")
		}
	}

	if ironic.Networking.Keepalived != nil && ironic.Networking.Keepalived.Enabled {
		if ironic.Networking.IPAddressManager == metal3api.IPAddressManagerKeepalived { //nolint:staticcheck // backward compat
			return errors.New("networking: keepalived and ipAddressManager cannot be used together")
		}
		if ironic.HighAvailability {
			return errors.New("networking: keepalived is not compatible with the highly available architecture")
		}
		if ironic.Networking.IPAddress == "" || ironic.Networking.Interface == "" {
			return errors.New("networking: keepalived requires specifying both ipAddress and interface")
		}
		for i, entry := range ironic.Networking.Keepalived.AdditionalVIPs {
			if entry.IPAddress == "" {
				return fmt.Errorf("networking.keepalived.additionalVIPs[%d]: ipAddress is required", i)
			}
			if err := validateIP(entry.IPAddress); err != nil {
				return fmt.Errorf("networking.keepalived.additionalVIPs[%d]: %w", i, err)
			}
			if entry.Interface == "" {
				return fmt.Errorf("networking.keepalived.additionalVIPs[%d]: interface is required", i)
			}
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
			return fmt.Errorf("networkingService.providerNetworks[%d] (%s): %w", i, pn.Type, err)
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
