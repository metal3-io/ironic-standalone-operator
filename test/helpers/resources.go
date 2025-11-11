package helpers

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

var (
	customIronicImage        string
	CustomImageVersion       string
	customIPADownloaderImage string
	customKeepalivedImage    string
)

func LoadCustomImages() {
	customIronicImage = os.Getenv("IRONIC_CUSTOM_IMAGE")
	CustomImageVersion = os.Getenv("IRONIC_CUSTOM_VERSION")
	customIPADownloaderImage = os.Getenv("IPA_DOWNLOADER_CUSTOM_IMAGE")
	customKeepalivedImage = os.Getenv("KEEPALIVED_CUSTOM_IMAGE")
}

func UsesCustomImage() bool {
	return customIronicImage != "" || CustomImageVersion != ""
}

func SkipIfCustomImage() {
	GinkgoHelper()
	if UsesCustomImage() {
		Skip("skipping because a custom image is provided")
	}
}

func NewTLSSecret(ctx context.Context, k8sClient client.Client, namespace, name string) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			corev1.TLSCertKey:       ironicCertPEM,
			corev1.TLSPrivateKeyKey: ironicKeyPEM,
		},
		Type: corev1.SecretTypeTLS,
	}
	err := k8sClient.Create(ctx, secret)
	Expect(err).NotTo(HaveOccurred())
	return secret
}

func NewAuthSecret(ctx context.Context, k8sClient client.Client, namespace, name string) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			corev1.BasicAuthUsernameKey: []byte("admin"),
			corev1.BasicAuthPasswordKey: []byte("test-password"),
		},
		Type: corev1.SecretTypeBasicAuth,
	}
	err := k8sClient.Create(ctx, secret)
	Expect(err).NotTo(HaveOccurred())
	return secret
}

func NewIronic(ctx context.Context, k8sClient client.Client, nname types.NamespacedName, spec metal3api.IronicSpec) *metal3api.Ironic {
	ironic := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nname.Name,
			Namespace: nname.Namespace,
		},
		Spec: spec,
	}

	if customIronicImage != "" {
		ironic.Spec.Images.Ironic = customIronicImage
	}
	if CustomImageVersion != "" {
		ironic.Spec.Version = CustomImageVersion
	}
	if customIPADownloaderImage != "" {
		ironic.Spec.Images.DeployRamdiskDownloader = customIPADownloaderImage
	}
	if customKeepalivedImage != "" {
		ironic.Spec.Images.Keepalived = customKeepalivedImage
	}

	err := k8sClient.Create(ctx, ironic)
	Expect(err).NotTo(HaveOccurred())
	return ironic
}
