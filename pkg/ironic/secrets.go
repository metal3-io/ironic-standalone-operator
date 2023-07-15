package ironic

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/go-logr/logr"
	"golang.org/x/crypto/bcrypt"
	corev1 "k8s.io/api/core/v1"
)

func checkValidUser(user string) error {
	for i, r := range user {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return fmt.Errorf("username cannot contain symbol %v (position %d)", r, i)
		}
	}
	return nil
}

func generateHtpasswd(user, password string) (string, error) {
	// A common source of errors: an accidental line break after a password
	user = strings.Trim(user, " \n\r")
	password = strings.Trim(password, " \n\r")
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

	return generateHtpasswd(string(user), string(password))
}

const htpasswdKey = "htpasswd"

func getAuthConfig(secret *corev1.Secret) string {
	user := secret.Data["username"]
	password := secret.Data["password"]
	return fmt.Sprintf(`
[DEFAULT]
auth_strategy = http_basic
http_basic_auth_user_file = /etc/ironic/htpasswd
[ironic]
auth_type = http_basic
username = %s
password = %s
[inspector]
auth_type = http_basic
username = %s
password = %s
`, user, password, user, password)
}

func secretNeedsUpdating(secret *corev1.Secret, logger logr.Logger) bool {
	existing := secret.Data[htpasswdKey]
	user, password, ok := strings.Cut(string(existing), ":")
	if ok && user != "" && password != "" {
		newUser, ok := secret.Data["username"]
		if ok && string(newUser) == user {
			newPassword, ok := secret.Data["password"]
			if ok && bcrypt.CompareHashAndPassword([]byte(password), []byte(newPassword)) == nil {
				authConfig := secret.Data["auth-config"]
				if string(authConfig) == getAuthConfig(secret) {
					// All good, keep the secret the way it is
					return false
				} else {
					logger.Info("API secret needs updating: outdated auth-config")
				}
			} else {
				logger.Info("API secret needs updating: passwords don't match")
			}
		} else {
			logger.Info("API secret needs updating: users don't match")
		}
	} else {
		logger.Info("API secret needs updating: no or malformed htpasswd")
	}

	return true
}

func UpdateSecret(cctx ControllerContext, secret *corev1.Secret) (bool, error) {
	if !secretNeedsUpdating(secret, cctx.Logger) {
		return false, nil
	}

	new, err := htpasswdFromSecret(secret)
	if err != nil {
		return false, err
	}
	secret.Data[htpasswdKey] = []byte(new)
	secret.Data["auth-config"] = []byte(getAuthConfig(secret))
	return true, nil
}
