package helpers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

func VerifyNetworkingDeploymentExists(ctx context.Context, clientset *kubernetes.Clientset, namespace, ironicName string) {
	By("verifying networking deployment exists")

	deployName := ironicName + "-networking"
	deploy, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deployName, metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())

	// Check labels
	Expect(deploy.Labels["app"]).To(Equal("ironic-networking"))
	Expect(deploy.Labels["ironic.metal3.io/instance"]).To(Equal(ironicName))

	// Check replicas
	Expect(deploy.Spec.Replicas).NotTo(BeNil())
	Expect(*deploy.Spec.Replicas).To(Equal(int32(1)))

	// Check container
	Expect(deploy.Spec.Template.Spec.Containers).To(HaveLen(1))
	container := deploy.Spec.Template.Spec.Containers[0]
	Expect(container.Name).To(Equal("ironic-networking"))
	Expect(container.Command).To(Equal([]string{"/bin/runironic-networking"}))
	Expect(container.Ports).To(HaveLen(1))
	Expect(container.Ports[0].Name).To(Equal("networking-rpc"))

	// hostNetwork should be false
	Expect(deploy.Spec.Template.Spec.HostNetwork).To(BeFalse())
}

func VerifyNetworkingServiceExists(ctx context.Context, clientset *kubernetes.Clientset, namespace, ironicName string) {
	By("verifying networking service exists")

	svcName := ironicName + "-networking-service"
	svc, err := clientset.CoreV1().Services(namespace).Get(ctx, svcName, metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())

	Expect(svc.Spec.Type).To(Equal(corev1.ServiceTypeClusterIP))
	Expect(svc.Spec.Ports).To(HaveLen(1))
	Expect(svc.Spec.Ports[0].Port).To(Equal(int32(6190)))
}

func VerifySwitchConfigSecretExists(ctx context.Context, clientset *kubernetes.Clientset, namespace, ironicName string) {
	By("verifying switch config secret exists")

	secretName := ironicName + "-switch-config"
	secret, err := clientset.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())

	Expect(secret.Labels["ironic.metal3.io/managed"]).To(Equal("true"))
	Expect(secret.Labels[metal3api.LabelEnvironmentName]).To(Equal(metal3api.LabelEnvironmentValue))
	Expect(secret.Data).To(HaveKey("switch-configs.conf"))
}

func VerifySwitchCredentialsSecretExists(ctx context.Context, clientset *kubernetes.Clientset, namespace, ironicName string) {
	By("verifying switch credentials secret exists")

	secretName := ironicName + "-switch-credentials"
	secret, err := clientset.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())

	Expect(secret.Labels["ironic.metal3.io/managed"]).To(Equal("true"))
	Expect(secret.Labels[metal3api.LabelEnvironmentName]).To(Equal(metal3api.LabelEnvironmentValue))
}

func VerifySwitchSecretsGone(ctx context.Context, clientset *kubernetes.Clientset, namespace, ironicName string) {
	By("verifying switch secrets are cleaned up")

	configName := ironicName + "-switch-config"
	credsName := ironicName + "-switch-credentials"

	Eventually(func() bool {
		_, configErr := clientset.CoreV1().Secrets(namespace).Get(ctx, configName, metav1.GetOptions{})
		_, credsErr := clientset.CoreV1().Secrets(namespace).Get(ctx, credsName, metav1.GetOptions{})
		return (configErr != nil && k8serrors.IsNotFound(configErr)) &&
			(credsErr != nil && k8serrors.IsNotFound(credsErr))
	}).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(BeTrue())
}

func VerifyNetworkingEnvVarsOnIronic(ctx context.Context, clientset *kubernetes.Clientset, namespace, ironicName string) {
	By("verifying networking env vars on ironic container")

	deployName := ironicName + "-service"
	deploy, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deployName, metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())

	var ironicContainer *corev1.Container
	for i := range deploy.Spec.Template.Spec.Containers {
		if deploy.Spec.Template.Spec.Containers[i].Name == "ironic" {
			ironicContainer = &deploy.Spec.Template.Spec.Containers[i]
			break
		}
	}
	Expect(ironicContainer).NotTo(BeNil(), "ironic container not found")

	envMap := make(map[string]string)
	for _, env := range ironicContainer.Env {
		envMap[env.Name] = env.Value
	}

	Expect(envMap).To(HaveKeyWithValue("IRONIC_NETWORKING_ENABLED", "true"))
	Expect(envMap).To(HaveKey("IRONIC_NETWORKING_JSON_RPC_HOST"))
	Expect(envMap).To(HaveKeyWithValue("IRONIC_DEFAULT_NETWORK_INTERFACE", "ironic-networking"))
}

func VerifyNetworkingResourcesGone(ctx context.Context, clientset *kubernetes.Clientset, namespace, ironicName string) {
	By("verifying networking resources are cleaned up")

	deployName := ironicName + "-networking"
	svcName := ironicName + "-networking-service"

	Eventually(func() bool {
		_, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deployName, metav1.GetOptions{})
		if err != nil && k8serrors.IsNotFound(err) {
			_, svcErr := clientset.CoreV1().Services(namespace).Get(ctx, svcName, metav1.GetOptions{})
			return svcErr != nil && k8serrors.IsNotFound(svcErr)
		}
		return false
	}).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(BeTrue())
}

func GetNetworkingPodAnnotation(ctx context.Context, clientset *kubernetes.Clientset, namespace, ironicName, key string) string {
	deployName := ironicName + "-networking"
	deploy, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deployName, metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())
	return deploy.Spec.Template.Annotations[key]
}

func VerifyNetworkingDeploymentNotExists(ctx context.Context, clientset *kubernetes.Clientset, namespace, ironicName string) {
	deployName := ironicName + "-networking"
	_, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deployName, metav1.GetOptions{})
	Expect(k8serrors.IsNotFound(err)).To(BeTrue(), "expected networking deployment to not exist")
}

func VerifyNetworkingServiceNotExists(ctx context.Context, clientset *kubernetes.Clientset, namespace, ironicName string) {
	svcName := ironicName + "-networking-service"
	_, err := clientset.CoreV1().Services(namespace).Get(ctx, svcName, metav1.GetOptions{})
	Expect(k8serrors.IsNotFound(err)).To(BeTrue(), "expected networking service to not exist")
}
