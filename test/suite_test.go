/*

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/apiversions"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/httpbasic"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/noauth"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/conductors"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/nodes"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	yaml "sigs.k8s.io/yaml"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

// NOTE(dtantsur): these two constants refer to the Ironic API version (which
// is different from the version of Ironic itself). Versions are incremented
// every time the API is changed. The listing of all versions is here:
// https://docs.openstack.org/ironic/latest/contributor/webapi-version-history.html
const (
	// NOTE(dtantsur): latest is now at least 1.95, so we can rely on this
	// value to check that specifying Version: 27.0 actually installs 27.0
	apiVersionIn270 = "1.94"
	apiVersionIn280 = "1.95"
	// Update this periodically to make sure we're installing the latest version by default
	knownAPIMinorVersion = 95
)

var ctx context.Context
var k8sClient client.Client
var clientset *kubernetes.Clientset

var ironicCertPEM []byte
var ironicKeyPEM []byte

var customImage string
var customImageVersion string
var customDatabaseImage string

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Functional tests")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	ctx = context.TODO()

	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
	}
	Expect(kubeconfigPath).To(BeAnExistingFile(), "Failed to get the kubeconfig file for the cluster")

	var err error

	err = metal3api.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	Expect(err).NotTo(HaveOccurred())

	clientset, err = kubernetes.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	ironicCertPEM, err = os.ReadFile(os.Getenv("IRONIC_CERT_FILE"))
	Expect(err).NotTo(HaveOccurred())
	ironicKeyPEM, err = os.ReadFile(os.Getenv("IRONIC_KEY_FILE"))
	Expect(err).NotTo(HaveOccurred())

	customImage = os.Getenv("IRONIC_CUSTOM_IMAGE")
	customImageVersion = os.Getenv("IRONIC_CUSTOM_VERSION")
	customDatabaseImage = os.Getenv("MARIADB_CUSTOM_IMAGE")
})

func addHTTPTransport(serviceClient *gophercloud.ServiceClient) *gophercloud.ServiceClient {
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(ironicCertPEM)

	tlsConfig := &tls.Config{
		RootCAs:    certPool,
		MinVersion: tls.VersionTLS13,
	}
	serviceClient.HTTPClient.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	return serviceClient
}

func NewNoAuthClient(endpoint string) (*gophercloud.ServiceClient, error) {
	serviceClient, err := noauth.NewBareMetalNoAuth(noauth.EndpointOpts{
		IronicEndpoint: endpoint,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create an Ironic client: %w", err)
	}

	return addHTTPTransport(serviceClient), nil
}

func NewHTTPBasicClient(endpoint string, secret *corev1.Secret) (*gophercloud.ServiceClient, error) {
	serviceClient, err := httpbasic.NewBareMetalHTTPBasic(httpbasic.EndpointOpts{
		IronicEndpoint:     endpoint,
		IronicUser:         string(secret.Data[corev1.BasicAuthUsernameKey]),
		IronicUserPassword: string(secret.Data[corev1.BasicAuthPasswordKey]),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create an Ironic client: %w", err)
	}

	return addHTTPTransport(serviceClient), nil
}

func logResources(ironic *metal3api.Ironic, suffix string) {
	deployName := fmt.Sprintf("%s-service", ironic.Name)
	if ironic.Spec.HighAvailability {
		deploy, err := clientset.AppsV1().DaemonSets(ironic.Namespace).Get(ctx, deployName, metav1.GetOptions{})
		if err == nil {
			GinkgoWriter.Printf(".. status of daemon set: %+v\n", deploy.Status)
			writeYAML(deploy, deploy.Namespace, deploy.Name, "daemonset"+suffix)
		} else if !k8serrors.IsNotFound(err) {
			Expect(err).NotTo(HaveOccurred())
		}
	} else {
		deploy, err := clientset.AppsV1().Deployments(ironic.Namespace).Get(ctx, deployName, metav1.GetOptions{})
		if err == nil {
			GinkgoWriter.Printf(".. status of deployment: %+v\n", deploy.Status)
			writeYAML(deploy, deploy.Namespace, deploy.Name, "deployment"+suffix)
		} else if !k8serrors.IsNotFound(err) {
			Expect(err).NotTo(HaveOccurred())
		}
	}

	pods, err := clientset.CoreV1().Pods(ironic.Namespace).List(ctx, metav1.ListOptions{})
	Expect(err).NotTo(HaveOccurred())

	for _, pod := range pods.Items {
		GinkgoWriter.Printf("... status of pod %s: %+v\n", pod.Name, pod.Status)
	}
}

func WaitForIronic(name types.NamespacedName) *metal3api.Ironic {
	ironic := &metal3api.Ironic{}

	By("waiting for Ironic deployment")

	Eventually(func() bool {
		err := k8sClient.Get(ctx, name, ironic)
		Expect(err).NotTo(HaveOccurred())

		writeYAML(ironic, ironic.Namespace, ironic.Name, "ironic")
		GinkgoWriter.Printf("Current status of Ironic: %+v\n", ironic.Status)

		cond := meta.FindStatusCondition(ironic.Status.Conditions, string(metal3api.IronicStatusReady))
		if cond != nil {
			expectedVersion := ironic.Spec.Version
			if expectedVersion != "" {
				Expect(ironic.Status.RequestedVersion).To(Equal(expectedVersion))
			} else {
				Expect(ironic.Status.RequestedVersion).ToNot(BeEmpty())
			}

			if cond.Status == metav1.ConditionTrue {
				if expectedVersion != "" {
					Expect(ironic.Status.InstalledVersion).To(Equal(expectedVersion))
				} else {
					Expect(ironic.Status.InstalledVersion).ToNot(BeEmpty())
				}
				return true
			}
		}

		logResources(ironic, "")
		return false
	}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(BeTrue())

	return ironic
}

func WaitForUpgrade(name types.NamespacedName, fromVersion, toVersion string) *metal3api.Ironic {
	ironic := &metal3api.Ironic{}
	suffix := "-upgraded-" + toVersion

	By("waiting for Ironic deployment")

	Eventually(func() bool {
		err := k8sClient.Get(ctx, name, ironic)
		Expect(err).NotTo(HaveOccurred())

		writeYAML(ironic, ironic.Namespace, ironic.Name, "ironic")
		GinkgoWriter.Printf("Current status of Ironic: %+v\n", ironic.Status)

		cond := meta.FindStatusCondition(ironic.Status.Conditions, string(metal3api.IronicStatusReady))
		Expect(cond).ToNot(BeNil(), "on upgrade, the Ready condition must always be present")

		Expect(ironic.Spec.Version).To(Equal(toVersion), "unexpected Version in the Spec")
		upgradeAcknowledged := ironic.Status.RequestedVersion == toVersion
		upgraded := ironic.Status.InstalledVersion == toVersion
		if upgraded {
			Expect(upgradeAcknowledged).To(BeTrue(), "InstalledVersion set before RequestedVersion")
		}

		if upgradeAcknowledged && cond.Status == metav1.ConditionTrue {
			Expect(upgraded).To(BeTrue(), "the Ready condition set before InstalledVersion")
			logResources(ironic, suffix)
			return true
		} else {
			Expect(upgraded).To(BeFalse(), "InstalledVersion set before the Ready condition")
		}

		logResources(ironic, suffix)
		return false
	}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(BeTrue())

	return ironic
}

func WaitForIronicFailure(name types.NamespacedName, message string) *metal3api.Ironic {
	ironic := &metal3api.Ironic{}

	By("waiting for Ironic deployment to fail")

	Eventually(func() bool {
		err := k8sClient.Get(ctx, name, ironic)
		Expect(err).NotTo(HaveOccurred())

		writeYAML(ironic, ironic.Namespace, ironic.Name, "ironic")
		GinkgoWriter.Printf("Current status of Ironic: %+v\n", ironic.Status)

		cond := meta.FindStatusCondition(ironic.Status.Conditions, string(metal3api.IronicStatusReady))
		if cond == nil {
			GinkgoWriter.Printf("No Ready condition yet\n")
			return false
		}
		Expect(cond.Status).To(Equal(metav1.ConditionFalse), "Unexpected Ready status")
		if cond.Reason != metal3api.IronicReasonFailed {
			return false
		}

		if !strings.Contains(cond.Message, message) {
			GinkgoWriter.Printf("Different error for now: %s\n", cond.Message)
			return false
		}

		return true
	}).WithTimeout(15 * time.Minute).WithPolling(10 * time.Second).Should(BeTrue())

	return ironic
}

func writeYAML(obj interface{}, namespace, name, typ string) {
	fileDir := fmt.Sprintf("%s/%s", os.Getenv("LOGDIR"), namespace)
	err := os.MkdirAll(fileDir, 0755)
	Expect(err).NotTo(HaveOccurred())

	fileName := fmt.Sprintf("%s/%s_%s.yaml", fileDir, typ, name)
	yamlData, err := yaml.Marshal(obj)
	Expect(err).NotTo(HaveOccurred())
	err = os.WriteFile(fileName, yamlData, 0600)
	Expect(err).NotTo(HaveOccurred())
}

type TestAssumptions struct {
	// Assume that TLS is used
	withTLS bool

	// Assume that this secret is used for accessing the API
	apiSecret *corev1.Secret

	// Verify that an active conductor with this name exists
	activeConductor string

	// Verify that the maximum available version equals this one
	maxAPIVersion string
}

func verifyAPIVersion(ctx context.Context, cli *gophercloud.ServiceClient, assumptions TestAssumptions) {
	By("checking Ironic API versions")

	version, err := apiversions.Get(ctx, cli, "v1").Extract()
	Expect(err).NotTo(HaveOccurred())
	if assumptions.maxAPIVersion != "" {
		Expect(version.Version).To(Equal(assumptions.maxAPIVersion))
	} else if customImage == "" && customImageVersion == "" {
		// NOTE(dtantsur): we cannot make any assumptions about the provided image and version,
		// so this check only runs when they are not set.
		var minorVersion int
		_, err = fmt.Sscanf(version.Version, "1.%d", &minorVersion)
		Expect(err).NotTo(HaveOccurred())
		Expect(minorVersion).To(BeNumerically(">=", knownAPIMinorVersion))
	}
}

func VerifyIronic(ironic *metal3api.Ironic, assumptions TestAssumptions) {
	writeYAML(ironic, ironic.Namespace, ironic.Name, "ironic")

	proto := "http"
	if assumptions.withTLS {
		proto = "https"
	}
	ironicURL := fmt.Sprintf("%s://%s:6385", proto, os.Getenv("IRONIC_IP"))

	By("checking the service")

	svc, err := clientset.CoreV1().Services(ironic.Namespace).Get(ctx, ironic.Name, metav1.GetOptions{})
	GinkgoWriter.Printf("Ironic service: %+v\n", svc)
	Expect(err).NotTo(HaveOccurred())

	writeYAML(svc, svc.Namespace, svc.Name, "service")

	Expect(svc.Spec.ClusterIPs).ToNot(BeEmpty())

	By("fetching the authentication secret")

	secret := assumptions.apiSecret
	if secret == nil {
		secret, err = clientset.CoreV1().Secrets(ironic.Namespace).Get(ctx, ironic.Spec.APICredentialsName, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())
	} else {
		Expect(ironic.Spec.APICredentialsName).To(Equal(secret.Name))
	}

	By("checking Ironic authentication")

	// Do not let the test get stuck here in case of connection issues
	withTimeout, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	cli, err := NewNoAuthClient(ironicURL)
	Expect(err).NotTo(HaveOccurred())
	_, err = nodes.List(cli, nodes.ListOpts{}).AllPages(withTimeout)
	Expect(err).To(HaveOccurred())

	cli, err = NewHTTPBasicClient(ironicURL, secret)
	Expect(err).NotTo(HaveOccurred())
	cli.Microversion = "1.81" // minimum version supported by BMO

	verifyAPIVersion(withTimeout, cli, assumptions)

	By("checking Ironic conductor list")

	conductorPager, err := conductors.List(cli, conductors.ListOpts{}).AllPages(withTimeout)
	Expect(err).NotTo(HaveOccurred())
	conductors, err := conductors.ExtractConductors(conductorPager)
	Expect(err).NotTo(HaveOccurred())

	var aliveConductors []string
	for _, cond := range conductors {
		if cond.Alive {
			aliveConductors = append(aliveConductors, cond.Hostname)
		}
	}
	if assumptions.activeConductor != "" {
		Expect(aliveConductors).To(ContainElement(assumptions.activeConductor))
	} else {
		Expect(aliveConductors).NotTo(BeEmpty())
	}

	By("creating and deleting a lot of Nodes")

	withTimeout, cancel = context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	drivers := []string{"ipmi", "redfish"}
	var nodeUUIDs []string
	for idx := 0; idx < 100; idx++ {
		node, err := nodes.Create(withTimeout, cli, nodes.CreateOpts{
			Driver: drivers[rand.Intn(len(drivers))], //nolint:gosec // weak crypto is ok in tests
			Name:   fmt.Sprintf("node-%d", idx),
		}).Extract()
		Expect(err).NotTo(HaveOccurred())
		nodeUUIDs = append(nodeUUIDs, node.UUID)
	}

	for _, nodeUUID := range nodeUUIDs {
		err = nodes.Delete(withTimeout, cli, nodeUUID).ExtractErr()
		Expect(err).NotTo(HaveOccurred())
	}
}

func writeContainerLogs(pod *corev1.Pod, containerName, logDir string) {
	podLogOpts := corev1.PodLogOptions{Container: containerName}
	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(ctx)
	Expect(err).NotTo(HaveOccurred())
	defer podLogs.Close()

	targetFileName := fmt.Sprintf("%s/%s.log", logDir, containerName)
	logFile, err := os.Create(targetFileName)
	Expect(err).NotTo(HaveOccurred())
	defer logFile.Close()

	_, err = io.Copy(logFile, podLogs)
	Expect(err).NotTo(HaveOccurred())
}

func CollectLogs(namespace string) {
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	Expect(err).NotTo(HaveOccurred())

	for _, pod := range pods.Items {
		logDir := fmt.Sprintf("%s/%s/pod_%s", os.Getenv("LOGDIR"), namespace, pod.Name)
		err = os.MkdirAll(logDir, 0755)
		Expect(err).NotTo(HaveOccurred())

		writeYAML(&pod, namespace, pod.Name, "pod")

		for _, cont := range pod.Spec.Containers {
			writeContainerLogs(&pod, cont.Name, logDir)
		}
	}
}

func DeleteAndWait(ironic *metal3api.Ironic) {
	By("deleting Ironic")

	name := types.NamespacedName{
		Name:      ironic.Name,
		Namespace: ironic.Namespace,
	}
	err := k8sClient.Delete(ctx, ironic)
	Expect(err).NotTo(HaveOccurred())

	Eventually(func() bool {
		err := k8sClient.Get(ctx, name, ironic)
		if err != nil && k8serrors.IsNotFound(err) {
			return true
		}
		GinkgoWriter.Println("Ironic", name, "not deleted yet, error is", err)
		Expect(err).NotTo(HaveOccurred())
		return false
	}).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(BeTrue())
}

func buildIronic(name types.NamespacedName, spec metal3api.IronicSpec) *metal3api.Ironic {
	result := &metal3api.Ironic{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name.Name,
			Namespace: name.Namespace,
		},
		Spec: spec,
	}

	if customImage != "" {
		result.Spec.Images.Ironic = customImage
	}
	if customImageVersion != "" {
		result.Spec.Version = customImageVersion
	}

	return result
}

func buildDatabase(name types.NamespacedName) *metal3api.IronicDatabase {
	result := &metal3api.IronicDatabase{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-db", name.Name),
			Namespace: name.Namespace,
		},
	}

	if customDatabaseImage != "" {
		result.Spec.Image = customDatabaseImage
	}

	return result
}

var _ = Describe("Ironic object tests", func() {
	var namespace string

	BeforeEach(func() {
		namespace = fmt.Sprintf("test-%s", CurrentSpecReport().LeafNodeLabels[0])
		nsSpec := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}

		_, err := clientset.CoreV1().Namespaces().Create(ctx, nsSpec, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			_ = clientset.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
			Eventually(func() bool {
				_, err := clientset.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
				if err != nil && k8serrors.IsNotFound(err) {
					return true
				}
				GinkgoWriter.Println("Namespace", namespace, "not deleted yet, error is", err)
				Expect(err).NotTo(HaveOccurred())
				return false
			}).WithTimeout(2 * time.Minute).WithPolling(5 * time.Second).Should(BeTrue())
		})
	})

	It("creates Ironic without any parameters", Label("no-params"), func() {
		name := types.NamespacedName{
			Name:      "test-ironic",
			Namespace: namespace,
		}

		ironic := buildIronic(name, metal3api.IronicSpec{})
		err := k8sClient.Create(ctx, ironic)
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			CollectLogs(namespace)
			DeleteAndWait(ironic)
		})

		ironic = WaitForIronic(name)
		VerifyIronic(ironic, TestAssumptions{})
	})

	It("creates Ironic with provided credentials", Label("api-secret"), func() {
		By("creating a failing Ironic with non-existent credentials")

		name := types.NamespacedName{
			Name:      "test-ironic",
			Namespace: namespace,
		}

		ironic := buildIronic(name, metal3api.IronicSpec{
			APICredentialsName: "banana",
		})
		err := k8sClient.Create(ctx, ironic)
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			CollectLogs(namespace)
			DeleteAndWait(ironic)
		})

		_ = WaitForIronicFailure(name, fmt.Sprintf("API credentials secret %s/banana not found", namespace))

		By("creating the secret and recovering the Ironic")

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-api", name.Name),
				Namespace: namespace,
			},
			Data: map[string][]byte{
				corev1.BasicAuthUsernameKey: []byte("admin"),
				corev1.BasicAuthPasswordKey: []byte("test-password"),
			},
			Type: corev1.SecretTypeBasicAuth,
		}
		err = k8sClient.Create(ctx, secret)
		Expect(err).NotTo(HaveOccurred())

		patch := client.MergeFrom(ironic.DeepCopy())
		ironic.Spec.APICredentialsName = secret.Name
		err = k8sClient.Patch(ctx, ironic, patch)
		Expect(err).NotTo(HaveOccurred())

		ironic = WaitForIronic(name)
		VerifyIronic(ironic, TestAssumptions{apiSecret: secret})
	})

	It("creates Ironic with TLS", Label("tls"), func() {
		name := types.NamespacedName{
			Name:      "test-ironic",
			Namespace: namespace,
		}

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-api", name.Name),
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

		ironic := buildIronic(name, metal3api.IronicSpec{
			TLS: metal3api.TLS{
				CertificateName: secret.Name,
			},
		})
		err = k8sClient.Create(ctx, ironic)
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			CollectLogs(namespace)
			DeleteAndWait(ironic)
		})

		ironic = WaitForIronic(name)
		VerifyIronic(ironic, TestAssumptions{withTLS: true})
	})

	It("creates Ironic 27.0 and upgrades to 28.0", Label("v270-to-280"), func() {
		if customImage != "" || customImageVersion != "" {
			Skip("skipping because a custom image is provided")
		}

		name := types.NamespacedName{
			Name:      "test-ironic",
			Namespace: namespace,
		}

		ironic := buildIronic(name, metal3api.IronicSpec{
			Version: "27.0",
		})
		err := k8sClient.Create(ctx, ironic)
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			CollectLogs(namespace)
			DeleteAndWait(ironic)
		})

		ironic = WaitForIronic(name)
		VerifyIronic(ironic, TestAssumptions{maxAPIVersion: apiVersionIn270})

		By("upgrading to Ironic 28.0")

		patch := client.MergeFrom(ironic.DeepCopy())
		ironic.Spec.Version = "28.0"
		err = k8sClient.Patch(ctx, ironic, patch)
		Expect(err).NotTo(HaveOccurred())

		ironic = WaitForUpgrade(name, "27.0", "28.0")
		VerifyIronic(ironic, TestAssumptions{maxAPIVersion: apiVersionIn280})
	})

	It("creates Ironic 28.0 and upgrades to latest", Label("v280-to-latest"), func() {
		if customImage != "" || customImageVersion != "" {
			Skip("skipping because a custom image is provided")
		}

		name := types.NamespacedName{
			Name:      "test-ironic",
			Namespace: namespace,
		}

		ironic := buildIronic(name, metal3api.IronicSpec{
			Version: "28.0",
		})
		err := k8sClient.Create(ctx, ironic)
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			CollectLogs(namespace)
			DeleteAndWait(ironic)
		})

		ironic = WaitForIronic(name)
		VerifyIronic(ironic, TestAssumptions{maxAPIVersion: apiVersionIn280})

		By("upgrading to Ironic latest")

		patch := client.MergeFrom(ironic.DeepCopy())
		ironic.Spec.Version = "latest"
		err = k8sClient.Patch(ctx, ironic, patch)
		Expect(err).NotTo(HaveOccurred())

		ironic = WaitForUpgrade(name, "28.0", "latest")
		VerifyIronic(ironic, TestAssumptions{})
	})

	It("creates Ironic with keepalived and DHCP", Label("keepalived-dnsmasq"), func() {
		name := types.NamespacedName{
			Name:      "test-ironic",
			Namespace: namespace,
		}

		ironic := buildIronic(name, metal3api.IronicSpec{
			Networking: metal3api.Networking{
				DHCP: &metal3api.DHCP{
					NetworkCIDR: os.Getenv("PROVISIONING_CIDR"),
				},
				Interface:        os.Getenv("PROVISIONING_INTERFACE"),
				IPAddress:        os.Getenv("PROVISIONING_IP"),
				IPAddressManager: metal3api.IPAddressManagerKeepalived,
			},
		})
		err := k8sClient.Create(ctx, ironic)
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			CollectLogs(namespace)
			DeleteAndWait(ironic)
		})

		ironic = WaitForIronic(name)
		VerifyIronic(ironic, TestAssumptions{})
	})

	It("creates Ironic with non-existent database", Label("non-existent-database"), func() {
		name := types.NamespacedName{
			Name:      "test-ironic",
			Namespace: namespace,
		}

		ironic := buildIronic(name, metal3api.IronicSpec{
			DatabaseName: "banana",
		})
		err := k8sClient.Create(ctx, ironic)
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			CollectLogs(namespace)
			DeleteAndWait(ironic)
		})

		_ = WaitForIronicFailure(name, fmt.Sprintf("database %s/banana not found", namespace))
	})

	It("creates highly available Ironic", Label("high-availability-no-provnet"), func() {
		name := types.NamespacedName{
			Name:      "test-ironic",
			Namespace: namespace,
		}

		ironicDB := buildDatabase(name)

		err := k8sClient.Create(ctx, ironicDB)
		Expect(err).NotTo(HaveOccurred())

		ironic := buildIronic(name, metal3api.IronicSpec{
			DatabaseName:     ironicDB.Name,
			HighAvailability: true,
		})
		err = k8sClient.Create(ctx, ironic)
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			CollectLogs(namespace)
			DeleteAndWait(ironic)
		})

		ironic = WaitForIronic(name)
		VerifyIronic(ironic, TestAssumptions{})
	})

	It("creates Ironic with extraConfig", Label("extra-config"), func() {
		name := types.NamespacedName{
			Name:      "test-ironic",
			Namespace: namespace,
		}

		conductorName := "test-conductor"

		ironic := buildIronic(name, metal3api.IronicSpec{
			ExtraConfig: []metal3api.ExtraConfig{
				{
					Name:  "host",
					Value: conductorName,
				},
			},
		})
		err := k8sClient.Create(ctx, ironic)
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			CollectLogs(namespace)
			DeleteAndWait(ironic)
		})

		ironic = WaitForIronic(name)
		VerifyIronic(ironic, TestAssumptions{activeConductor: conductorName})
	})
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
})
