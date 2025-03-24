package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsing(t *testing.T) {
	testCases := []struct {
		Scenario string

		Value string

		Expected    Version
		ExpectError string
	}{
		{
			Scenario: "latest",
			Value:    "latest",
		},
		{
			Scenario: "valid major",
			Value:    "27.0",
			Expected: Version{Major: 27, Minor: 0},
		},
		{
			Scenario: "valid minor",
			Value:    "27.2",
			Expected: Version{Major: 27, Minor: 2},
		},
		{
			Scenario:    "invalid leading zero",
			Value:       "0.42",
			ExpectError: "invalid major version 0 in 0.42",
		},
		{
			Scenario:    "invalid major",
			Value:       "foo.42",
			ExpectError: "invalid major version foo in foo.42",
		},
		{
			Scenario:    "invalid minor",
			Value:       "42.foo",
			ExpectError: "invalid minor version foo in 42.foo",
		},
		{
			Scenario:    "invalid structure",
			Value:       "1,2",
			ExpectError: "invalid version 1,2, expected MAJOR.MINOR",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			result, err := ParseVersion(tc.Value)
			if tc.ExpectError != "" {
				assert.ErrorContains(t, err, tc.ExpectError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.Expected, result)
			}
		})
	}
}

func TestComparing(t *testing.T) {
	testCases := []struct {
		Scenario string

		First  Version
		Second Version

		Expected int
	}{
		{
			Scenario: "latest equal latest",
			First:    VersionLatest,
			Second:   VersionLatest,
			Expected: 0,
		},
		{
			Scenario: "latest > any",
			First:    VersionLatest,
			Second:   Version{Major: 99, Minor: 99},
			Expected: 1,
		},
		{
			Scenario: "equal",
			First:    Version{Major: 42, Minor: 0},
			Second:   Version{Major: 42, Minor: 0},
			Expected: 0,
		},
		{
			Scenario: "compare major",
			First:    Version{Major: 41, Minor: 99},
			Second:   Version{Major: 42, Minor: 0},
			Expected: -1,
		},
		{
			Scenario: "compare minor",
			First:    Version{Major: 42, Minor: 99},
			Second:   Version{Major: 42, Minor: 0},
			Expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			result := tc.First.Compare(tc.Second)
			assert.Equal(t, tc.Expected, result)
		})
	}
}
