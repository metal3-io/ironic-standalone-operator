package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestDHCPDefaults(t *testing.T) {
	testCases := []struct {
		Scenario string

		Prefix        string
		ExpectedFirst string
		ExpectedLast  string
	}{
		{
			Scenario:      "v4",
			Prefix:        "10.1.42.0/24",
			ExpectedFirst: "10.1.42.10",
			ExpectedLast:  "10.1.42.253",
		},
		{
			Scenario:      "v6",
			Prefix:        "2001:db8::/112",
			ExpectedFirst: "2001:db8::a",
			ExpectedLast:  "2001:db8::fffd",
		},
		{
			Scenario: "broken",
			Prefix:   "10.1.42.0/32",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			dhcp := &DHCP{NetworkCIDR: tc.Prefix}

			SetDHCPDefaults(dhcp)
			assert.Equal(t, tc.ExpectedFirst, dhcp.RangeBegin)
			assert.Equal(t, tc.ExpectedLast, dhcp.RangeEnd)
		})
	}
}

func TestValidateIronic(t *testing.T) {
	testCases := []struct {
		Scenario string

		Ironic        IronicSpec
		OldIronic     *IronicSpec
		ExpectedError string
	}{
		{
			Scenario: "empty",
		},
		{
			Scenario: "with database",
			Ironic: IronicSpec{
				DatabaseRef: corev1.LocalObjectReference{
					Name: "db",
				},
			},
		},
		{
			Scenario: "with ipAddress",
			Ironic: IronicSpec{
				Networking: Networking{
					IPAddress: "192.168.0.2",
				},
			},
		},
		{
			Scenario: "with ipAddress-v6",
			Ironic: IronicSpec{
				Networking: Networking{
					IPAddress: "2001:db8::2",
				},
			},
		},
		{
			Scenario: "bad ipAddress",
			Ironic: IronicSpec{
				Networking: Networking{
					IPAddress: "banana",
				},
			},
			ExpectedError: "banana is not a valid IP address",
		},
		{
			Scenario: "bad externalIP",
			Ironic: IronicSpec{
				Networking: Networking{
					ExternalIP: "banana",
				},
			},
			ExpectedError: "banana is not a valid IP address",
		},
		{
			Scenario: "HA needs database",
			Ironic: IronicSpec{
				HighAvailability: true,
			},
			ExpectedError: "database is required",
		},
		{
			Scenario: "no ipAddress with HA",
			Ironic: IronicSpec{
				DatabaseRef: corev1.LocalObjectReference{
					Name: "db",
				},
				Networking: Networking{
					IPAddress: "192.168.0.1",
				},
				HighAvailability: true,
			},
			ExpectedError: "ipAddress makes no sense",
		},
		{
			Scenario: "HA disabled",
			Ironic: IronicSpec{
				DatabaseRef: corev1.LocalObjectReference{
					Name: "db",
				},
				HighAvailability: true,
			},
			ExpectedError: "highly available architecture is disabled",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			if tc.OldIronic == nil {
				tc.OldIronic = &tc.Ironic
			}

			err := validateIronic(&tc.Ironic, tc.OldIronic)
			if tc.ExpectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.ExpectedError)
			}
		})
	}
}
