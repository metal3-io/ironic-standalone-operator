package ironic

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metal3api "github.com/metal3-io/ironic-operator/api/v1alpha1"
)

const (
	databaseDeploymentName string = "metal3-database"
)

// RemoveDatabase removes the MariaDB database.
func RemoveDatabase(ctx context.Context, kubeClient kubernetes.Interface, ironic *metal3api.Ironic) error {
	return client.IgnoreNotFound(kubeClient.AppsV1().Deployments(ironic.Namespace).Delete(context.Background(), databaseDeploymentName, metav1.DeleteOptions{}))
}
