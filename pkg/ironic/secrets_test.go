package ironic

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestIsValidUser(t *testing.T) {
	testCases := []struct {
		Scenario string

		User          string
		ExpectedError bool
	}{
		{
			Scenario: "normal-user",
			User:     "admin",
		},
		{
			Scenario: "user-with-number",
			User:     "super2000",
		},
		{
			Scenario:      "has-space",
			User:          "super 2000",
			ExpectedError: true,
		},
		{
			Scenario:      "has-colon",
			User:          "super:2000",
			ExpectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			result := checkValidUser(tc.User)
			if tc.ExpectedError {
				assert.Error(t, result)
			} else {
				assert.NoError(t, result)
			}
		})
	}
}

func TestGenerateHtpasswd(t *testing.T) {
	result, err := generateHtpasswd("admin", "pa$$w0rd")
	assert.NoError(t, err)
	user, password, ok := strings.Cut(result, ":")
	assert.Truef(t, ok, "%s does not start with admin:", result)
	assert.Equal(t, "admin", user)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(password), []byte("pa$$w0rd")))
}
