package ironic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildEndpoints(t *testing.T) {
	testCases := []struct {
		Scenario string

		IPs          []string
		Port         int
		IncludeProto string

		Expected []string
	}{
		{
			Scenario: "non-standard-port-no-protocol",

			IPs:          []string{"2001:db8::42", "192.0.2.42"},
			Port:         6385,
			IncludeProto: "",

			Expected: []string{"192.0.2.42:6385", "[2001:db8::42]:6385"},
		},
		{
			Scenario: "non-standard-port-with-protocol",

			IPs:          []string{"2001:db8::42", "192.0.2.42"},
			Port:         6385,
			IncludeProto: "http",

			Expected: []string{"http://192.0.2.42:6385", "http://[2001:db8::42]:6385"},
		},
		{
			Scenario: "http-with-protocol",

			IPs:          []string{"2001:db8::42", "192.0.2.42"},
			Port:         80,
			IncludeProto: "http",

			Expected: []string{"http://192.0.2.42", "http://[2001:db8::42]"},
		},
		{
			Scenario: "https-with-protocol",

			IPs:          []string{"2001:db8::42", "192.0.2.42"},
			Port:         443,
			IncludeProto: "https",

			Expected: []string{"https://192.0.2.42", "https://[2001:db8::42]"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			result := buildEndpoints(tc.IPs, tc.Port, tc.IncludeProto)
			assert.Equal(t, tc.Expected, result)
		})
	}
}
