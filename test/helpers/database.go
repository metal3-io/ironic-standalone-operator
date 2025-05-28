package helpers

import (
	"context"
	"fmt"
	"os"

	. "github.com/onsi/gomega"

	mariadbapi "github.com/mariadb-operator/mariadb-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

var (
	mariadbName      string
	mariadbNamespace string
)

func LoadDatabaseParams() {
	mariadbName = os.Getenv("MARIADB_NAME")
	mariadbNamespace = os.Getenv("MARIADB_NAMESPACE")
}

func CreateDatabase(ctx context.Context, k8sClient client.Client, name types.NamespacedName) *metal3api.Database {
	mariadbRef := mariadbapi.MariaDBRef{
		ObjectReference: mariadbapi.ObjectReference{
			Name:      mariadbName,
			Namespace: mariadbNamespace,
		},
		WaitForIt: true,
	}
	// In these tests namespace is more unique than name
	dbUserName := name.Namespace

	database := &mariadbapi.Database{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dbUserName,
			Namespace: name.Namespace,
		},
		Spec: mariadbapi.DatabaseSpec{
			MariaDBRef:   mariadbRef,
			CharacterSet: "utf8",
			Collate:      "utf8_general_ci",
			SQLTemplate: mariadbapi.SQLTemplate{
				CleanupPolicy: ptr.To(mariadbapi.CleanupPolicyDelete),
			},
		},
	}
	err := k8sClient.Create(ctx, database)
	Expect(err).NotTo(HaveOccurred())

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dbUserName,
			Namespace: name.Namespace,
		},
		Data: map[string][]byte{
			corev1.BasicAuthUsernameKey: []byte(dbUserName),
			corev1.BasicAuthPasswordKey: []byte("test-password"),
		},
		Type: corev1.SecretTypeBasicAuth,
	}
	err = k8sClient.Create(ctx, secret)
	Expect(err).NotTo(HaveOccurred())

	user := &mariadbapi.User{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secret.Name,
			Namespace: name.Namespace,
		},
		Spec: mariadbapi.UserSpec{
			MariaDBRef: mariadbRef,
			PasswordSecretKeyRef: &mariadbapi.SecretKeySelector{
				LocalObjectReference: mariadbapi.LocalObjectReference{
					Name: secret.Name,
				},
				Key: corev1.BasicAuthPasswordKey,
			},
			Host: "%",
			SQLTemplate: mariadbapi.SQLTemplate{
				CleanupPolicy: ptr.To(mariadbapi.CleanupPolicyDelete),
			},
			MaxUserConnections: 100,
		},
	}
	err = k8sClient.Create(ctx, user)
	Expect(err).NotTo(HaveOccurred())

	grant := &mariadbapi.Grant{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name.Name,
			Namespace: name.Namespace,
		},
		Spec: mariadbapi.GrantSpec{
			MariaDBRef: mariadbRef,
			Privileges: []string{"ALL PRIVILEGES"},
			Database:   database.Name,
			Table:      "*",
			Username:   dbUserName,
			Host:       ptr.To("%"),
			SQLTemplate: mariadbapi.SQLTemplate{
				CleanupPolicy: ptr.To(mariadbapi.CleanupPolicyDelete),
			},
		},
	}
	err = k8sClient.Create(ctx, grant)
	Expect(err).NotTo(HaveOccurred())

	return &metal3api.Database{
		CredentialsName: secret.Name,
		Host:            fmt.Sprintf("%s.%s.svc", mariadbName, mariadbNamespace),
		Name:            database.Name,
	}
}
