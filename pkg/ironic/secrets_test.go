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
		Password string //nolint:gosec // test struct field, not a real secret
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
	testCases := []struct {
		Scenario string

		User            string
		Password        string //nolint:gosec // test struct field, not a real secret
		CurrentHtpasswd string

		ExpectedChanged bool
	}{
		{
			Scenario: "nothing-changed",

			User:            "admin",
			Password:        "password",
			CurrentHtpasswd: "admin:$2y$05$CJozjmp4SHJjNWcJn1vVsOx4OEBQTDTVTdNFc0I.CVt5xpEZMK4pW",
		},
		{
			Scenario: "newlines",

			User:            "admin\n",
			Password:        "password\n",
			CurrentHtpasswd: "admin:$2y$05$CJozjmp4SHJjNWcJn1vVsOx4OEBQTDTVTdNFc0I.CVt5xpEZMK4pW",
		},
		{
			Scenario: "new-value",

			User:            "admin",
			Password:        "password",
			ExpectedChanged: true,
		},
		{
			Scenario: "user-changed",

			User:            "admin2",
			Password:        "password",
			CurrentHtpasswd: "admin:$2y$05$CJozjmp4SHJjNWcJn1vVsOx4OEBQTDTVTdNFc0I.CVt5xpEZMK4pW",
			ExpectedChanged: true,
		},
		{
			Scenario: "password-changed",

			User:            "admin",
			Password:        "password2",
			CurrentHtpasswd: "admin:$2y$05$CJozjmp4SHJjNWcJn1vVsOx4OEBQTDTVTdNFc0I.CVt5xpEZMK4pW",
			ExpectedChanged: true,
		},
		{
			Scenario: "missing-user",

			Password:        "password",
			CurrentHtpasswd: "admin:$2y$05$CJozjmp4SHJjNWcJn1vVsOx4OEBQTDTVTdNFc0I.CVt5xpEZMK4pW",
			ExpectedChanged: true,
		},
		{
			Scenario: "missing-password",

			User:            "admin",
			CurrentHtpasswd: "admin:$2y$05$CJozjmp4SHJjNWcJn1vVsOx4OEBQTDTVTdNFc0I.CVt5xpEZMK4pW",
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
			} else {
				assert.Nil(t, secret.Data["htpasswd"])
			}
			assert.Equal(t, "test", secret.Namespace)
			assert.Equal(t, "my-ironic-foo-", secret.GenerateName)
		})
	}
}

func TestSecretVersionAnnotations(t *testing.T) {
	testCases := []struct {
		Scenario string

		SecretType string
		SecretData map[string][]byte

		ExpectedAnnotationKey string
		ExpectedHashValue     string
	}{
		{
			Scenario:   "simple-secret",
			SecretType: "api-secret",
			SecretData: map[string][]byte{
				"username": []byte("admin"),
				"password": []byte("secret123"),
			},
			ExpectedAnnotationKey: "ironic.metal3.io/api-secret-version",
			ExpectedHashValue:     "d9ab4722976729ee3be458f3bc652126e8ed4297e4fa2f05632d92317356fda8",
		},
		{
			Scenario:   "tls-secret",
			SecretType: "tls-secret",
			SecretData: map[string][]byte{
				"tls.crt": []byte("certificate-data"),
				"tls.key": []byte("key-data"),
			},
			ExpectedAnnotationKey: "ironic.metal3.io/tls-secret-version",
			ExpectedHashValue:     "41b814de6978e1e9bbaa5520ec7988c28228e77ac98422a9b7173a8ab91a835a",
		},
		{
			Scenario:   "single-field-secret",
			SecretType: "database",
			SecretData: map[string][]byte{
				"password": []byte("dbpassword"),
			},
			ExpectedAnnotationKey: "ironic.metal3.io/database-version",
			ExpectedHashValue:     "487745f63f2b3980402dd732c333d65d32f80ce7ee342d1955501b0af660bd9f",
		},
		{
			Scenario:              "empty-secret",
			SecretType:            "empty",
			SecretData:            map[string][]byte{},
			ExpectedAnnotationKey: "ironic.metal3.io/empty-version",
			ExpectedHashValue:     "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			Scenario:   "multi-field-secret",
			SecretType: "complex",
			SecretData: map[string][]byte{
				"field1": []byte("value1"),
				"field2": []byte("value2"),
				"field3": []byte("value3"),
			},
			ExpectedAnnotationKey: "ironic.metal3.io/complex-version",
			ExpectedHashValue:     "7549f265c52de5c1dea2bd192d34f5c31452bb2d45352fa8ad4469d5055b4434",
		},
		{
			Scenario:   "binary-data",
			SecretType: "binary",
			SecretData: map[string][]byte{
				"cert": {0x00, 0x01, 0x02, 0xff, 0xfe, 0xfd},
			},
			ExpectedAnnotationKey: "ironic.metal3.io/binary-version",
			ExpectedHashValue:     "ffb95ee4f4809f6893b1fe84b096325cb57a34817d2fc4bf03710fc7ec36fcdf",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Scenario, func(t *testing.T) {
			secret := &corev1.Secret{
				Data: tc.SecretData,
			}
			annotations := secretVersionAnnotations(tc.SecretType, secret)
			assert.Len(t, annotations, 1)
			assert.Equal(t, tc.ExpectedHashValue, annotations[tc.ExpectedAnnotationKey])

			// Consistency check: the same call results in the same results
			secret = &corev1.Secret{
				Data: tc.SecretData,
			}
			annotations = secretVersionAnnotations(tc.SecretType, secret)
			assert.Equal(t, tc.ExpectedHashValue, annotations[tc.ExpectedAnnotationKey])
		})
	}
}

func TestSecretVersionAnnotationsNoBinaryCollision(t *testing.T) {
	// These two secrets differ but serialize identically under a NUL-delimited
	// scheme, so they must still hash differently to trigger a pod rollout.
	hashOf := func(data map[string][]byte) string {
		secret := &corev1.Secret{Data: data}
		return secretVersionAnnotations("api", secret)["ironic.metal3.io/api-version"]
	}

	first := hashOf(map[string][]byte{"a": []byte("x"), "b\x00y": []byte("z")})
	second := hashOf(map[string][]byte{"a": []byte("x\x00b"), "y": []byte("z")})
	assert.NotEqual(t, first, second)
}

func TestSecretVersionAnnotationsNoConcatCollision(t *testing.T) {
	// Two distinct secrets whose keys and values concatenate to the same byte
	// stream ("passwordxusernameusernamey") when no field boundary is encoded.
	// They must hash differently, otherwise a changed secret would not trigger
	// a pod rollout.
	hashOf := func(data map[string][]byte) string {
		secret := &corev1.Secret{Data: data}
		return secretVersionAnnotations("api", secret)["ironic.metal3.io/api-version"]
	}

	first := hashOf(map[string][]byte{"password": []byte("xusername"), "username": []byte("y")})
	second := hashOf(map[string][]byte{"password": []byte("x"), "username": []byte("usernamey")})
	assert.NotEqual(t, first, second)
}
