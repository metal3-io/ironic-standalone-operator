package ironic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	testCases := []struct {
		Scenario string

		User     string
		Password string
	}{
		{
			User:     "admin",
			Password: "pa$$w0rd",
		},
		{
			User:     "admin\n",
			Password: "pa$$w0rd\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			result, err := generateHtpasswd([]byte(tc.User), []byte(tc.Password))
			require.NoError(t, err)
			user, password, ok := strings.Cut(result, ":")
			assert.Truef(t, ok, "%s is not separated with a colon", result)
			assert.Equal(t, strings.Trim(tc.User, "\n"), user)
			assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(password), []byte(strings.Trim(tc.Password, "\n"))))
		})
	}
}

func TestSecretNeedsUpdating(t *testing.T) {
	authConfig := `
[DEFAULT]
auth_strategy = http_basic
http_basic_auth_user_file = /etc/ironic/htpasswd
[json_rpc]
auth_strategy = http_basic
auth_type = http_basic
username = admin
password = password
`
	testCases := []struct {
		Scenario string

		User            string
		Password        string
		CurrentHtpasswd string
		AuthConfig      string

		ExpectedChanged bool
	}{
		{
			Scenario: "nothing-changed",

			User:            "admin",
			Password:        "password",
			CurrentHtpasswd: "admin:$2y$05$CJozjmp4SHJjNWcJn1vVsOx4OEBQTDTVTdNFc0I.CVt5xpEZMK4pW",
			AuthConfig:      authConfig,
		},
		{
			Scenario: "newlines",

			User:            "admin\n",
			Password:        "password\n",
			CurrentHtpasswd: "admin:$2y$05$CJozjmp4SHJjNWcJn1vVsOx4OEBQTDTVTdNFc0I.CVt5xpEZMK4pW",
			AuthConfig:      authConfig,
		},
		{
			Scenario: "new-value",

			User:            "admin",
			Password:        "password",
			AuthConfig:      authConfig,
			ExpectedChanged: true,
		},
		{
			Scenario: "user-changed",

			User:            "admin2",
			Password:        "password",
			CurrentHtpasswd: "admin:$2y$05$CJozjmp4SHJjNWcJn1vVsOx4OEBQTDTVTdNFc0I.CVt5xpEZMK4pW",
			AuthConfig:      authConfig,
			ExpectedChanged: true,
		},
		{
			Scenario: "password-changed",

			User:            "admin",
			Password:        "password2",
			CurrentHtpasswd: "admin:$2y$05$CJozjmp4SHJjNWcJn1vVsOx4OEBQTDTVTdNFc0I.CVt5xpEZMK4pW",
			AuthConfig:      authConfig,
			ExpectedChanged: true,
		},
		{
			Scenario: "missing-user",

			Password:        "password",
			CurrentHtpasswd: "admin:$2y$05$CJozjmp4SHJjNWcJn1vVsOx4OEBQTDTVTdNFc0I.CVt5xpEZMK4pW",
			AuthConfig:      authConfig,
			ExpectedChanged: true,
		},
		{
			Scenario: "missing-password",

			User:            "admin",
			CurrentHtpasswd: "admin:$2y$05$CJozjmp4SHJjNWcJn1vVsOx4OEBQTDTVTdNFc0I.CVt5xpEZMK4pW",
			AuthConfig:      authConfig,
			ExpectedChanged: true,
		},
		{
			Scenario: "missing-auth-config",

			User:            "admin",
			Password:        "password",
			CurrentHtpasswd: "admin:$2y$05$CJozjmp4SHJjNWcJn1vVsOx4OEBQTDTVTdNFc0I.CVt5xpEZMK4pW",
			ExpectedChanged: true,
		},
		{
			Scenario: "outdated-auth-config",

			User:            "admin",
			Password:        "password",
			CurrentHtpasswd: "admin:$2y$05$CJozjmp4SHJjNWcJn1vVsOx4OEBQTDTVTdNFc0I.CVt5xpEZMK4pW",
			AuthConfig:      strings.ReplaceAll(authConfig, "admin", "user"),
			ExpectedChanged: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			secret := &corev1.Secret{
				Data: map[string][]byte{
					"username": []byte(tc.User),
					"password": []byte(tc.Password),
				},
			}
			if tc.CurrentHtpasswd != "" {
				secret.Data["htpasswd"] = []byte(tc.CurrentHtpasswd)
			}
			if tc.AuthConfig != "" {
				secret.Data["auth-config"] = []byte(tc.AuthConfig)
			}

			changed := secretNeedsUpdating(secret, logr.Discard())
			assert.Equal(t, tc.ExpectedChanged, changed)
		})
	}
}

func TestGenerateSecret(t *testing.T) {
	for _, tc := range []bool{true, false} {
		t.Run(fmt.Sprintf("extra-%v", tc), func(t *testing.T) {
			meta := &metav1.ObjectMeta{
				Name:      "my-ironic",
				Namespace: "test",
			}
			secret, err := GenerateSecret(meta, "foo", tc)
			require.NoError(t, err)
			require.NotNil(t, secret)
			assert.Len(t, secret.Data["password"], passwordLength)
			if tc {
				assert.NotNil(t, secret.Data["htpasswd"])
				assert.NotNil(t, secret.Data["auth-config"])
			} else {
				assert.Nil(t, secret.Data["htpasswd"])
				assert.Nil(t, secret.Data["auth-config"])
			}
			assert.Equal(t, "test", secret.Namespace)
			assert.Equal(t, "my-ironic-foo-", secret.GenerateName)
		})
	}
}
