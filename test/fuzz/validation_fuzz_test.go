package fuzz

import (
	"strings"
	"testing"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
	"github.com/metal3-io/ironic-standalone-operator/pkg/ironic"
)

// FuzzValidateIronic detects panics in ValidateIronic for arbitrary IP
// addresses, CIDR ranges, DHCP ranges, and image URLs. Validation errors
// are expected; the only invariant is that no panic occurs.
func FuzzValidateIronic(f *testing.F) {
	type seed struct {
		ipAddress    string
		externalIP   string
		networkCIDR  string
		rangeBegin   string
		rangeEnd     string
		dnsAddress   string
		kernelURL    string
		initramfsURL string
	}

	seeds := []seed{
		// Valid seeds
		{"192.168.1.1", "", "192.168.1.0/24", "192.168.1.10", "192.168.1.100", "", "http://example.com/ipa.kernel", "http://example.com/ipa.initramfs"},
		{"fd00::1", "", "fd00::/64", "fd00::10", "fd00::100", "", "http://[::1]/ipa.kernel", "http://[::1]/ipa.initramfs"},
		{"10.0.0.1", "10.0.0.1", "", "", "", "", "https://example.com/ipa.kernel", "https://example.com/ipa.initramfs"},
		{"172.16.0.1", "", "172.16.0.0/24", "172.16.0.10", "172.16.0.200", "172.16.0.1", "file:///shared/ipa.kernel", "file:///shared/ipa.initramfs"},
		{"192.168.0.1", "", "192.168.0.0/16", "192.168.0.10", "192.168.0.254", "", "oci://registry.example.com/ipa:latest", "oci://registry.example.com/ipa-initramfs:latest"},
		{"", "", "", "", "", "", "", ""},
		{"127.0.0.1", "", "127.0.0.0/8", "127.0.0.10", "127.0.0.254", "", "http://localhost/k", "http://localhost/i"},
		{"10.10.10.1", "", "10.10.10.0/30", "10.10.10.1", "10.10.10.2", "", "https://mirror.example.com/ipa.kernel", "https://mirror.example.com/ipa.initramfs"},
		{"192.168.100.1", "", "192.168.100.0/24", "192.168.100.50", "192.168.100.200", "8.8.8.8", "http://10.0.0.1/ipa.kernel", "http://10.0.0.1/ipa.initramfs"},
		{"10.20.30.1", "", "10.20.30.0/24", "10.20.30.10", "10.20.30.100", "", "https://registry.example.com:8443/ipa.kernel", "https://registry.example.com:8443/ipa.initramfs"},
		{"172.20.0.1", "", "172.20.0.0/24", "172.20.0.10", "172.20.0.200", "", "", ""},
		{"192.168.1.5", "", "192.168.1.5/24", "192.168.1.10", "192.168.1.100", "", "", ""},
		{"fd00::1", "", "fd00::/120", "fd00::a", "fd00::64", "", "", ""},

		// Invalid seeds
		{"not-an-ip", "also-not-ip", "not-a-cidr", "bad", "bad", "bad", "://broken", "ftp://unsupported"},
		{"192.168.5.1", "", "192.168.5.0/24", "10.0.0.1", "10.0.0.100", "", "", ""},
		{"10.1.1.1", "", "", "", "", "", "file://relative/path/ipa.kernel", "file://relative/path/ipa.initramfs"},
		{"10.2.2.1", "", "", "", "", "", "http:///ipa.kernel", "http:///ipa.initramfs"},
		{"10.3.3.1", "", "", "", "", "", "http://example.com/\x00ipa.kernel", "http://example.com/\x00ipa.initramfs"},
		{"10.4.4.1", "", "10.4.4.0/24", "10.4.4.10", "10.4.4.100", "", "http://" + strings.Repeat("a", 4096) + ".example.com/ipa.kernel", "http://10.4.4.1/ipa.initramfs"},
	}

	for _, s := range seeds {
		f.Add(s.ipAddress, s.externalIP, s.networkCIDR, s.rangeBegin, s.rangeEnd, s.dnsAddress, s.kernelURL, s.initramfsURL)
	}

	f.Fuzz(func(_ *testing.T,
		ipAddress string,
		externalIP string,
		networkCIDR string,
		rangeBegin string,
		rangeEnd string,
		dnsAddress string,
		kernelURL string,
		initramfsURL string,
	) {
		spec := &metal3api.IronicSpec{
			Networking: metal3api.Networking{
				IPAddress:  ipAddress,
				ExternalIP: externalIP,
			},
		}

		if networkCIDR != "" {
			spec.Networking.DHCP = &metal3api.DHCP{
				NetworkCIDR: networkCIDR,
				RangeBegin:  rangeBegin,
				RangeEnd:    rangeEnd,
				DNSAddress:  dnsAddress,
			}
		}

		if kernelURL != "" || initramfsURL != "" {
			spec.Overrides = &metal3api.Overrides{
				AgentImages: []metal3api.AgentImages{
					{
						Kernel:    kernelURL,
						Initramfs: initramfsURL,
					},
				},
			}
		}

		_ = ironic.ValidateIronic(spec, nil)
	})
}
