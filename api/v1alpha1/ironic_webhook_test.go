package v1alpha1

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
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
			dhcp := &DHCP{}
			cidr := netip.MustParsePrefix(tc.Prefix)

			setDHCPDefaults(dhcp, cidr)
			assert.Equal(t, tc.ExpectedFirst, dhcp.FirstIP)
			assert.Equal(t, tc.ExpectedLast, dhcp.LastIP)
		})
	}
}
