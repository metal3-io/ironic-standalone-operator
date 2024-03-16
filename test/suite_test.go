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
	"math/rand"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
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

var _ = Describe("Ironic object tests", func() {
	var namespace string

	BeforeEach(func() {
		namespace = fmt.Sprintf("test-%d", rand.Int())
		nsSpec := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}

		_, err := clientset.CoreV1().Namespaces().Create(ctx, nsSpec, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			_ = clientset.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
		})
	})

	It("creates Ironic without any parameters", func() {
		name := types.NamespacedName{
			Name:      "test-ironic",
			Namespace: namespace,
		}

		ironic := metal3api.Ironic{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name.Name,
				Namespace: name.Namespace,
			},
		}
		err := k8sClient.Create(ctx, &ironic)
		Expect(err).NotTo(HaveOccurred())

		var cond *metav1.Condition
		for attempt := 0; attempt < 10; attempt += 1 {
			err = k8sClient.Get(ctx, name, &ironic)
			Expect(err).NotTo(HaveOccurred())

			cond = meta.FindStatusCondition(ironic.Status.Conditions, string(metal3api.IronicStatusAvailable))
			if cond != nil && cond.Status == metav1.ConditionTrue {
				break
			}

			time.Sleep(5 * time.Second)
		}

		if cond == nil || cond.Status != metav1.ConditionTrue {
			Fail("Ironic deployment never succeded")
		}
	})

})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
})
