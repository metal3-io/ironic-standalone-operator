package ironic

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"maps"
	"math/big"
	"net"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-logr/logr"
	"golang.org/x/crypto/bcrypt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

func checkValidUser(user string) error {
	for i, r := range user {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '-' && r != '_' {
			return fmt.Errorf("username cannot contain symbol %v (position %d)", r, i)
		}
	}
	return nil
}

// NOTE(dtantsur): this excludes most symbols that are hard to use in shell, or
// cause errors when substituting in SQL files, or are incompatible with
// the way MariaDB password is provided in ironic.conf. Make up for it by
// generating an absurdly long password.
const (
	passwordCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"
	passwordLength  = 32
)

func generatePassword() ([]byte, error) {
	password := make([]byte, passwordLength)
	maxIdx := big.NewInt(int64(len(passwordCharset)))
	for i := range passwordLength {
		idx, err := rand.Int(rand.Reader, maxIdx)
		if err != nil {
			return nil, fmt.Errorf("cannot generate a new password: %w", err)
		}
		password[i] = passwordCharset[idx.Int64()]
	}

	return password, nil
}

func normalizeSecretValue(value []byte) string {
	return strings.Trim(string(value), " \n\r\t")
}

func generateHtpasswd(userBytes, passwordBytes []byte) (string, error) {
	// A common source of errors: an accidental line break after a password
	user := normalizeSecretValue(userBytes)
	password := normalizeSecretValue(passwordBytes)
	err := checkValidUser(user)
	if err != nil {
		return "", err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("cannot generate a hashed password: %w", err)
	}

	return fmt.Sprintf("%s:%s", user, hashed), nil
}

func htpasswdFromSecret(secret *corev1.Secret) (string, error) {
	user, ok := secret.Data["username"]
	if !ok {
		return "", fmt.Errorf("missing username in secret %s/%s", secret.Namespace, secret.Name)
	}

	password, ok := secret.Data["password"]
	if !ok {
		return "", fmt.Errorf("missing password in secret %s/%s", secret.Namespace, secret.Name)
	}

	return generateHtpasswd(user, password)
}

const (
	htpasswdKey = "htpasswd"
	// cloudsYAMLKey holds a clouds.yaml for the baremetal CLI, so the
	// secret can be mounted as /etc/openstack in a client container.
	cloudsYAMLKey = "clouds.yaml"
)

func secretNeedsUpdating(secret *corev1.Secret, logger logr.Logger) bool {
	existing := secret.Data[htpasswdKey]
	user, password, ok := strings.Cut(string(existing), ":")
	if ok && user != "" && password != "" {
		newUser, ok := secret.Data[corev1.BasicAuthUsernameKey]
		newUserString := normalizeSecretValue(newUser)
		if ok && newUserString == user {
			newPassword, ok := secret.Data[corev1.BasicAuthPasswordKey]
			newPasswordString := normalizeSecretValue(newPassword)
			if ok && bcrypt.CompareHashAndPassword([]byte(password), []byte(newPasswordString)) == nil {
				// All good, keep the secret the way it is
				return false
			}
			logger.Info("API secret needs updating: passwords don't match")
		} else {
			logger.Info("API secret needs updating: users don't match", "HtpasswdUser", user, "CurrentUser", newUserString)
		}
	} else {
		logger.Info("API secret needs updating: no or malformed htpasswd")
	}

	return true
}

func UpdateSecret(secret *corev1.Secret, logger logr.Logger) (bool, error) {
	if !secretNeedsUpdating(secret, logger) {
		return false, nil
	}

	htpasswd, err := htpasswdFromSecret(secret)
	if err != nil {
		return false, err
	}
	secret.Data[htpasswdKey] = []byte(htpasswd)
	return true, nil
}

func GenerateSecret(owner *metav1.ObjectMeta, name string, extraFields bool) (*corev1.Secret, error) {
	pwd, err := generatePassword()
	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-%s-", owner.Name, name),
			Namespace:    owner.Namespace,
			Labels: map[string]string{
				// Add the environment label so the secret is included in the filtered cache
				metal3api.LabelEnvironmentName: metal3api.LabelEnvironmentValue,
			},
		},
		Data: map[string][]byte{
			corev1.BasicAuthUsernameKey: []byte(owner.Name),
			corev1.BasicAuthPasswordKey: pwd,
		},
		Type: corev1.SecretTypeBasicAuth,
	}

	if extraFields {
		_, err = UpdateSecret(secret, logr.Discard())
		if err != nil {
			return nil, err
		}
	}

	return secret, nil
}

// secretVersionAnnotations hashes the secret's data so that pods are rolled
// out when it changes. Keys in ignoreKeys are left out of the hash, for data
// that the pods don't consume.
func secretVersionAnnotations(secretType string, secret *corev1.Secret, ignoreKeys ...string) map[string]string {
	// Hash all data fields to detect changes
	h := sha256.New()

	// Sort keys for consistent hash ordering
	keys := slices.DeleteFunc(slices.Sorted(maps.Keys(secret.Data)), func(k string) bool {
		return slices.Contains(ignoreKeys, k)
	})

	// Hash each key and value into fixed-width digests so field boundaries are
	// unambiguous. Writing keys/values directly would let different secrets
	// concatenate to the same bytes and hash identically, skipping a rollout.
	for _, k := range keys {
		keyHash := sha256.Sum256([]byte(k))
		valueHash := sha256.Sum256(secret.Data[k])
		h.Write(keyHash[:])
		h.Write(valueHash[:])
	}

	hash := hex.EncodeToString(h.Sum(nil))

	return map[string]string{
		fmt.Sprintf("ironic.metal3.io/%s-version", secretType): hash,
	}
}

type cloudsConfig struct {
	Clouds map[string]cloudConfig `json:"clouds"`
}

type cloudConfig struct {
	AuthType                  string    `json:"auth_type"`
	Auth                      cloudAuth `json:"auth"`
	BaremetalEndpointOverride string    `json:"baremetal_endpoint_override"`
}

type cloudAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// generateCloudsYAML builds a clouds.yaml with an "internal" cloud that uses
// the Ironic service and, when the host IP is known from the spec, an
// "external" cloud that uses the host IP.
func generateCloudsYAML(ironic *metal3api.Ironic, user, password string) ([]byte, error) {
	scheme := "http"
	if ironic.Spec.TLS.CertificateName != "" {
		scheme = "https"
	}
	port := ironic.Spec.Networking.APIPort
	if port == 0 {
		port = 6385
	}
	endpoint := func(host string) string {
		return fmt.Sprintf("%s://%s", scheme, net.JoinHostPort(host, strconv.Itoa(int(port))))
	}
	cloud := func(host string) cloudConfig {
		return cloudConfig{
			AuthType:                  "http_basic",
			Auth:                      cloudAuth{Username: user, Password: password},
			BaremetalEndpointOverride: endpoint(host),
		}
	}

	config := cloudsConfig{Clouds: map[string]cloudConfig{
		"internal": cloud(fmt.Sprintf("%s.%s.svc", ironic.Name, ironic.Namespace)),
	}}
	if hostIP := ironic.Spec.Networking.ExternalIP; hostIP != "" {
		config.Clouds["external"] = cloud(hostIP)
	} else if hostIP := ironic.Spec.Networking.IPAddress; hostIP != "" {
		config.Clouds["external"] = cloud(hostIP)
	}

	return yaml.Marshal(config)
}

// UpdateCloudsYAML makes sure the API secret has an up-to-date clouds.yaml.
// It returns true if the secret was changed and needs saving.
func UpdateCloudsYAML(secret *corev1.Secret, ironic *metal3api.Ironic) (bool, error) {
	user, ok := secret.Data[corev1.BasicAuthUsernameKey]
	if !ok {
		return false, fmt.Errorf("missing username in secret %s/%s", secret.Namespace, secret.Name)
	}
	password, ok := secret.Data[corev1.BasicAuthPasswordKey]
	if !ok {
		return false, fmt.Errorf("missing password in secret %s/%s", secret.Namespace, secret.Name)
	}

	clouds, err := generateCloudsYAML(ironic, normalizeSecretValue(user), normalizeSecretValue(password))
	if err != nil {
		return false, fmt.Errorf("cannot generate clouds.yaml: %w", err)
	}
	if string(secret.Data[cloudsYAMLKey]) == string(clouds) {
		return false, nil
	}
	if secret.Data == nil {
		secret.Data = map[string][]byte{}
	}
	secret.Data[cloudsYAMLKey] = clouds
	return true, nil
}
