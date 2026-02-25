package ironic

import (
	"testing"

	"github.com/stretchr/testify/assert"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

func TestValidateIronic(t *testing.T) {
	testCases := []struct {
		Scenario string

		Ironic        metal3api.IronicSpec
		OldIronic     *metal3api.IronicSpec
		ExpectedError string
	}{
		{
			Scenario: "empty",
		},
		{
			Scenario: "with database",
			Ironic: metal3api.IronicSpec{
				Database: &metal3api.Database{
					CredentialsName: "test",
					Host:            "example.com",
					Name:            "ironic",
				},
			},
		},
		{
			Scenario: "adding database",
			Ironic: metal3api.IronicSpec{
				Database: &metal3api.Database{
					CredentialsName: "test",
					Host:            "example.com",
					Name:            "ironic",
				},
			},
			OldIronic: &metal3api.IronicSpec{},
		},
		{
			Scenario: "removing database",
			Ironic:   metal3api.IronicSpec{},
			OldIronic: &metal3api.IronicSpec{
				Database: &metal3api.Database{
					CredentialsName: "test",
					Host:            "example.com",
					Name:            "ironic",
				},
			},
		},
		{
			Scenario: "with ipAddress",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					IPAddress: "192.168.0.2",
				},
			},
		},
		{
			Scenario: "with ipAddress-v6",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					IPAddress: "2001:db8::2",
				},
			},
		},
		{
			Scenario: "bad ipAddress",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					IPAddress: "banana",
				},
			},
			ExpectedError: "banana is not a valid IP address",
		},
		{
			Scenario: "bad externalIP",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					ExternalIP: "banana",
				},
			},
			ExpectedError: "banana is not a valid IP address",
		},
		{
			Scenario: "HA needs database",
			Ironic: metal3api.IronicSpec{
				HighAvailability: true,
			},
			ExpectedError: "database is required",
		},
		{
			Scenario: "no ipAddress with HA",
			Ironic: metal3api.IronicSpec{
				Database: &metal3api.Database{
					CredentialsName: "test",
					Host:            "example.com",
					Name:            "ironic",
				},
				Networking: metal3api.Networking{
					IPAddress: "192.168.0.1",
				},
				HighAvailability: true,
			},
			ExpectedError: "ipAddress makes no sense",
		},
		{
			Scenario: "HA disabled",
			Ironic: metal3api.IronicSpec{
				Database: &metal3api.Database{
					CredentialsName: "test",
					Host:            "example.com",
					Name:            "ironic",
				},
				HighAvailability: true,
			},
			ExpectedError: "highly available architecture is disabled",
		},
		{
			Scenario: "no DHCP with HA",
			Ironic: metal3api.IronicSpec{
				Database: &metal3api.Database{
					CredentialsName: "test",
					Host:            "example.com",
					Name:            "ironic",
				},
				HighAvailability: true,
				Networking: metal3api.Networking{
					DHCP: &metal3api.DHCP{},
				},
			},
			ExpectedError: "DHCP support is not implemented",
		},
		{
			Scenario: "With Keepalived, no DHCP",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					Interface:        "eth0",
					IPAddress:        "192.0.2.2",
					IPAddressManager: metal3api.IPAddressManagerKeepalived,
				},
			},
		},
		{
			Scenario: "With Keepalived and DHCP",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					DHCP: &metal3api.DHCP{
						NetworkCIDR: "192.0.2.1/24",
						RangeBegin:  "192.0.2.10",
						RangeEnd:    "192.0.2.200",
					},
					Interface:        "eth0",
					IPAddress:        "192.0.2.2",
					IPAddressManager: metal3api.IPAddressManagerKeepalived,
				},
			},
		},
		{
			Scenario: "Keepalived requires Interface",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					IPAddress:        "192.0.2.2",
					IPAddressManager: metal3api.IPAddressManagerKeepalived,
				},
			},
			ExpectedError: "keepalived requires specifying both ipAddress and interface",
		},
		{
			Scenario: "Keepalived requires IPAddress",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					Interface:        "eth0",
					IPAddressManager: metal3api.IPAddressManagerKeepalived,
				},
			},
			ExpectedError: "keepalived requires specifying both ipAddress and interface",
		},
		{
			Scenario: "Keepalived exclusive with HA",
			Ironic: metal3api.IronicSpec{
				Database: &metal3api.Database{
					CredentialsName: "test",
					Host:            "example.com",
					Name:            "ironic",
				},
				HighAvailability: true,
				Networking: metal3api.Networking{
					Interface:        "eth0",
					IPAddress:        "192.0.2.2",
					IPAddressManager: metal3api.IPAddressManagerKeepalived,
				},
			},
			// NOTE(dtantsur): the expected error here is shadowed by the prior validation.
			// I'm keeping this test in place to ensure that *some* validation failure happens.
			ExpectedError: "ipAddress makes no sense with highly available architecture",
		},
		{
			Scenario: "with version",
			Ironic: metal3api.IronicSpec{
				Version: "32.0",
			},
		},
		{
			Scenario: "with invalid version",
			Ironic: metal3api.IronicSpec{
				Version: "banana",
			},
			ExpectedError: "invalid version banana, expected MAJOR.MINOR",
		},
		{
			Scenario: "with unsupported version",
			Ironic: metal3api.IronicSpec{
				Version: "42.42",
			},
			ExpectedError: "version 42.42 is not supported, supported versions are 32.0, 33.0, 34.0, latest",
		},
		{
			Scenario: "change existing database config",
			Ironic: metal3api.IronicSpec{
				Database: &metal3api.Database{
					CredentialsName: "newtest",
					Host:            "newexample.com",
					Name:            "newironic",
				},
			},
			OldIronic: &metal3api.IronicSpec{
				Database: &metal3api.Database{
					CredentialsName: "oldtest",
					Host:            "oldexample.com",
					Name:            "oldironic",
				},
			},
			ExpectedError: "cannot change to a new database",
		},
		{
			Scenario: "incomplete database config",
			Ironic: metal3api.IronicSpec{
				Database: &metal3api.Database{
					CredentialsName: "test",
					Host:            "example.com",
				},
			},
			ExpectedError: "credentialsName, host and name are required on database",
		},
		{
			Scenario: "configure RPC when HA is disabled",
			Ironic: metal3api.IronicSpec{
				TLS: metal3api.TLS{
					InsecureRPC: func() *bool { b := true; return &b }(),
				},
				HighAvailability: false,
			},
			ExpectedError: "insecureRPC makes no sense without highAvailability",
		},
		{
			Scenario: "DHCP without networking identity",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					DHCP: &metal3api.DHCP{
						NetworkCIDR: "192.168.1.0/24",
						RangeBegin:  "192.168.1.10",
						RangeEnd:    "192.168.1.100",
					},
				},
			},
			ExpectedError: "networking: at least one of ipAddress, interface or macAddresses is required when DHCP is used",
		},
		{
			Scenario: "serveDNS and dnsAddress configured simultaneously",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					Interface: "eth0",
					DHCP: &metal3api.DHCP{
						NetworkCIDR: "192.168.1.0/24",
						RangeBegin:  "192.168.1.10",
						RangeEnd:    "192.168.1.100",
						ServeDNS:    true,
						DNSAddress:  "8.8.8.8",
					},
				},
			},
			ExpectedError: "networking.dhcp.dnsAddress cannot set together with serveDNS",
		},
		{
			Scenario: "DHCP rangeBegin outside CDIR",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					Interface: "eth0",
					DHCP: &metal3api.DHCP{
						NetworkCIDR: "192.168.1.0/24",
						RangeBegin:  "10.0.0.10",
						RangeEnd:    "192.168.1.100",
					},
				},
			},
			ExpectedError: "10.0.0.10 is not in networking.dhcp.networkCIDR",
		},
		{
			Scenario: "Provisioning IP address not in the CIDR",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					IPAddress: "10.0.0.10",
					DHCP: &metal3api.DHCP{
						NetworkCIDR: "192.168.1.0/24",
						RangeBegin:  "192.168.1.10",
						RangeEnd:    "192.168.1.100",
					},
				},
			},
			ExpectedError: "networking.dhcp.networkCIDR must contain networking.ipAddress",
		},
		{
			Scenario: "invalid IP provided for dnsAddress",
			Ironic: metal3api.IronicSpec{
				Networking: metal3api.Networking{
					Interface: "eth0",
					DHCP: &metal3api.DHCP{
						NetworkCIDR: "192.168.1.0/24",
						RangeBegin:  "192.168.1.10",
						RangeEnd:    "192.168.1.100",
						DNSAddress:  "not-an-ip",
					},
				},
			},
			ExpectedError: "not-an-ip is not a valid IP address",
		},
		{
			Scenario: "HA incompatible with ServiceMonitor",
			Ironic: metal3api.IronicSpec{
				Database: &metal3api.Database{
					CredentialsName: "test",
					Host:            "example.com",
					Name:            "ironic",
				},
				HighAvailability: true,
				PrometheusExporter: &metal3api.PrometheusExporter{
					Enabled: true,
				},
			},
			ExpectedError: "ServiceMonitor support is currently incompatible with the highly available architecture",
		},
		{
			Scenario: "valid agent images single architecture x86_64",
			Ironic: metal3api.IronicSpec{
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Architecture: metal3api.ArchX86_64,
							Kernel:       "file:///shared/html/images/ipa.x86_64.kernel",
							Initramfs:    "file:///shared/html/images/ipa.x86_64.initramfs",
						},
					},
				},
			},
		},
		{
			Scenario: "valid agent images single architecture aarch64",
			Ironic: metal3api.IronicSpec{
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Architecture: metal3api.ArchAarch64,
							Kernel:       "file:///shared/html/images/ipa.aarch64.kernel",
							Initramfs:    "file:///shared/html/images/ipa.aarch64.initramfs",
						},
					},
				},
			},
		},
		{
			Scenario: "valid agent images multiple architectures",
			Ironic: metal3api.IronicSpec{
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Architecture: metal3api.ArchX86_64,
							Kernel:       "file:///shared/html/images/ipa.x86_64.kernel",
							Initramfs:    "file:///shared/html/images/ipa.x86_64.initramfs",
						},
						{
							Architecture: metal3api.ArchAarch64,
							Kernel:       "file:///shared/html/images/ipa.aarch64.kernel",
							Initramfs:    "file:///shared/html/images/ipa.aarch64.initramfs",
						},
					},
				},
			},
		},
		{
			Scenario: "agent images empty architecture",
			Ironic: metal3api.IronicSpec{
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Kernel:    "file:///shared/html/images/ipa.kernel",
							Initramfs: "file:///shared/html/images/ipa.initramfs",
						},
					},
				},
			},
			ExpectedError: "overrides.agentImages[0]: architecture is required",
		},
		{
			Scenario: "agent images empty kernel",
			Ironic: metal3api.IronicSpec{
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Architecture: metal3api.ArchX86_64,
							Initramfs:    "file:///shared/html/images/ipa.initramfs",
						},
					},
				},
			},
			ExpectedError: "overrides.agentImages[0]: kernel is required",
		},
		{
			Scenario: "agent images empty initramfs",
			Ironic: metal3api.IronicSpec{
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Architecture: metal3api.ArchX86_64,
							Kernel:       "file:///shared/html/images/ipa.kernel",
						},
					},
				},
			},
			ExpectedError: "overrides.agentImages[0]: initramfs is required",
		},
		{
			Scenario: "agent images duplicate architecture",
			Ironic: metal3api.IronicSpec{
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Architecture: metal3api.ArchX86_64,
							Kernel:       "file:///shared/html/images/ipa.x86_64.kernel",
							Initramfs:    "file:///shared/html/images/ipa.x86_64.initramfs",
						},
						{
							Architecture: metal3api.ArchX86_64,
							Kernel:       "file:///shared/html/images/ipa.x86_64.v2.kernel",
							Initramfs:    "file:///shared/html/images/ipa.x86_64.v2.initramfs",
						},
					},
				},
			},
			ExpectedError: "overrides.agentImages: duplicate architecture \"x86_64\"",
		},
		{
			Scenario: "agent images with http URL",
			Ironic: metal3api.IronicSpec{
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Architecture: metal3api.ArchX86_64,
							Kernel:       "http://example.com/ipa.kernel",
							Initramfs:    "http://example.com/ipa.initramfs",
						},
					},
				},
			},
		},
		{
			Scenario: "agent images with https URL",
			Ironic: metal3api.IronicSpec{
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Architecture: metal3api.ArchX86_64,
							Kernel:       "https://example.com/ipa.kernel",
							Initramfs:    "https://example.com/ipa.initramfs",
						},
					},
				},
			},
		},
		{
			Scenario: "agent images with invalid kernel URL",
			Ironic: metal3api.IronicSpec{
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Architecture: metal3api.ArchX86_64,
							Kernel:       "://invalid-url",
							Initramfs:    "file:///shared/html/images/ipa.initramfs",
						},
					},
				},
			},
			ExpectedError: "overrides.agentImages[0].kernel: invalid URL format",
		},
		{
			Scenario: "agent images with invalid initramfs URL",
			Ironic: metal3api.IronicSpec{
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Architecture: metal3api.ArchX86_64,
							Kernel:       "file:///shared/html/images/ipa.kernel",
							Initramfs:    "not a url",
						},
					},
				},
			},
			ExpectedError: "overrides.agentImages[0].initramfs: unsupported protocol",
		},
		{
			Scenario: "agent images with unsupported kernel protocol",
			Ironic: metal3api.IronicSpec{
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Architecture: metal3api.ArchX86_64,
							Kernel:       "ftp://example.com/ipa.kernel",
							Initramfs:    "file:///shared/html/images/ipa.initramfs",
						},
					},
				},
			},
			ExpectedError: "overrides.agentImages[0].kernel: unsupported protocol \"ftp\"",
		},
		{
			Scenario: "agent images with unsupported initramfs protocol",
			Ironic: metal3api.IronicSpec{
				Overrides: &metal3api.Overrides{
					AgentImages: []metal3api.AgentImages{
						{
							Architecture: metal3api.ArchX86_64,
							Kernel:       "file:///shared/html/images/ipa.kernel",
							Initramfs:    "ssh://example.com/ipa.initramfs",
						},
					},
				},
			},
			ExpectedError: "overrides.agentImages[0].initramfs: unsupported protocol \"ssh\"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			if tc.OldIronic == nil {
				tc.OldIronic = &tc.Ironic
			}

			err := ValidateIronic(&tc.Ironic, tc.OldIronic)
			if tc.ExpectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.ExpectedError)
			}
		})
	}
}
