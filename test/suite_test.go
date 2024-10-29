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
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gophercloud/gophercloud/v2"
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
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	yaml "sigs.k8s.io/yaml"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

var ctx context.Context
var cfg *rest.Config
var k8sClient client.Client
var clientset *kubernetes.Clientset

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

})

func NewNoAuthClient(endpoint string) (*gophercloud.ServiceClient, error) {
	return noauth.NewBareMetalNoAuth(noauth.EndpointOpts{
		IronicEndpoint: endpoint,
	})
}

func NewHTTPBasicClient(endpoint string, secret *corev1.Secret) (*gophercloud.ServiceClient, error) {
	return httpbasic.NewBareMetalHTTPBasic(httpbasic.EndpointOpts{
		IronicEndpoint:     endpoint,
		IronicUser:         string(secret.Data[corev1.BasicAuthUsernameKey]),
		IronicUserPassword: string(secret.Data[corev1.BasicAuthPasswordKey]),
	})
}

func WaitForIronic(name types.NamespacedName) *metal3api.Ironic {
	ironic := &metal3api.Ironic{}

	By("waiting for Ironic deployment")

	Eventually(func() bool {
		err := k8sClient.Get(ctx, name, ironic)
		Expect(err).NotTo(HaveOccurred())

		cond := meta.FindStatusCondition(ironic.Status.Conditions, string(metal3api.IronicStatusReady))
		if cond != nil && cond.Status == metav1.ConditionTrue {
			Expect(ironic.Status.InstalledVersion).ToNot(BeNil())
			return true
		}
		GinkgoWriter.Printf("Current status of Ironic: %+v\n", ironic.Status)

		deployName := fmt.Sprintf("%s-service", name.Name)
		if ironic.Spec.HighAvailability {
			deploy, err := clientset.AppsV1().DaemonSets(name.Namespace).Get(ctx, deployName, metav1.GetOptions{})
			if err == nil {
				GinkgoWriter.Printf(".. status of daemon set: %+v\n", deploy.Status)
			} else if !k8serrors.IsNotFound(err) {
				Expect(err).NotTo(HaveOccurred())
			}
		} else {
			deploy, err := clientset.AppsV1().Deployments(name.Namespace).Get(ctx, deployName, metav1.GetOptions{})
			if err == nil {
				GinkgoWriter.Printf(".. status of deployment: %+v\n", deploy.Status)
			} else if !k8serrors.IsNotFound(err) {
				Expect(err).NotTo(HaveOccurred())
			}
		}

		pods, err := clientset.CoreV1().Pods(name.Namespace).List(ctx, metav1.ListOptions{})
		Expect(err).NotTo(HaveOccurred())

		for _, pod := range pods.Items {
			GinkgoWriter.Printf("... status of pod %s: %+v\n", pod.Name, pod.Status)
		}

		return false
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
	err = os.WriteFile(fileName, yamlData, 0644)
	Expect(err).NotTo(HaveOccurred())
}

func VerifyIronic(ironic *metal3api.Ironic) {
	writeYAML(ironic, ironic.Namespace, ironic.Name, "ironic")

	By("checking the service")

	svc, err := clientset.CoreV1().Services(ironic.Namespace).Get(ctx, ironic.Name, metav1.GetOptions{})
	GinkgoWriter.Printf("Ironic service: %+v\n", svc)
	Expect(err).NotTo(HaveOccurred())

	writeYAML(svc, svc.Namespace, svc.Name, "service")

	Expect(svc.Spec.ClusterIPs).ToNot(BeEmpty())

	By("fetching the authentication secret")

	secret, err := clientset.CoreV1().Secrets(ironic.Namespace).Get(ctx, ironic.Spec.CredentialsRef.Name, metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())

	By("checking Ironic authentication")

	// Do not let the test get stuck here in case of connection issues
	withTimeout, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	cli, err := NewNoAuthClient("http://127.0.0.1:6385")
	Expect(err).NotTo(HaveOccurred())
	_, err = nodes.List(cli, nodes.ListOpts{}).AllPages(withTimeout)
	Expect(err).To(HaveOccurred())

	By("checking Ironic conductor list")

	cli, err = NewHTTPBasicClient("http://127.0.0.1:6385", secret)
	Expect(err).NotTo(HaveOccurred())
	cli.Microversion = "1.81" // minimum version supported by BMO

	conductorPager, err := conductors.List(cli, conductors.ListOpts{}).AllPages(withTimeout)
	Expect(err).NotTo(HaveOccurred())
	conductors, err := conductors.ExtractConductors(conductorPager)
	Expect(err).NotTo(HaveOccurred())
	Expect(conductors).NotTo(BeEmpty())

	By("creating and deleting a lot of Nodes")

	withTimeout, cancel = context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	drivers := []string{"ipmi", "redfish"}
	var nodeUUIDs []string
	for idx := 0; idx < 100; idx++ {
		node, err := nodes.Create(withTimeout, cli, nodes.CreateOpts{
			Driver: drivers[rand.Intn(len(drivers))],
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

func CollectLogs(namespace string) {
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	Expect(err).NotTo(HaveOccurred())

	for _, pod := range pods.Items {
		logDir := fmt.Sprintf("%s/%s/pod_%s", os.Getenv("LOGDIR"), namespace, pod.Name)
		err = os.MkdirAll(logDir, 0755)
		Expect(err).NotTo(HaveOccurred())

		writeYAML(&pod, namespace, pod.Name, "pod")

		for _, cont := range pod.Spec.Containers {
			podLogOpts := corev1.PodLogOptions{Container: cont.Name}
			req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
			podLogs, err := req.Stream(ctx)
			Expect(err).NotTo(HaveOccurred())
			defer podLogs.Close()

			logFile, err := os.Create(fmt.Sprintf("%s/%s.log", logDir, cont.Name))
			Expect(err).NotTo(HaveOccurred())
			defer logFile.Close()

			_, err = io.Copy(logFile, podLogs)
			Expect(err).NotTo(HaveOccurred())
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

		ironic := &metal3api.Ironic{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name.Name,
				Namespace: name.Namespace,
			},
		}
		err := k8sClient.Create(ctx, ironic)
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			CollectLogs(namespace)
			DeleteAndWait(ironic)
		})

		ironic = WaitForIronic(name)
		VerifyIronic(ironic)
	})

	It("creates highly available Ironic", Label("high-availability-no-provnet"), func() {
		name := types.NamespacedName{
			Name:      "test-ironic",
			Namespace: namespace,
		}

		ironicDb := &metal3api.IronicDatabase{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-db", name.Name),
				Namespace: name.Namespace,
			},
		}

		err := k8sClient.Create(ctx, ironicDb)
		Expect(err).NotTo(HaveOccurred())

		ironic := &metal3api.Ironic{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name.Name,
				Namespace: name.Namespace,
			},
			Spec: metal3api.IronicSpec{
				DatabaseRef: corev1.LocalObjectReference{
					Name: ironicDb.Name,
				},
				HighAvailability: true,
			},
		}
		err = k8sClient.Create(ctx, ironic)
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			CollectLogs(namespace)
			DeleteAndWait(ironic)
		})

		ironic = WaitForIronic(name)
		VerifyIronic(ironic)
	})
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
})
