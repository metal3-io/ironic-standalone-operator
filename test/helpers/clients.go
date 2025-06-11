package helpers

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	. "github.com/onsi/gomega"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/httpbasic"
	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/noauth"
	corev1 "k8s.io/api/core/v1"
)

var (
	ironicCertPEM []byte
	ironicKeyPEM  []byte
)

func LoadIronicCert() {
	var err error
	ironicCertPEM, err = os.ReadFile(os.Getenv("IRONIC_CERT_FILE"))
	Expect(err).NotTo(HaveOccurred())
	ironicKeyPEM, err = os.ReadFile(os.Getenv("IRONIC_KEY_FILE"))
	Expect(err).NotTo(HaveOccurred())
}

func addHTTPTransport(httpClient *http.Client) {
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(ironicCertPEM)

	tlsConfig := &tls.Config{
		RootCAs:    certPool,
		MinVersion: tls.VersionTLS13,
	}
	httpClient.Transport = &http.Transport{TLSClientConfig: tlsConfig}
}

func NewHTTPClient() http.Client {
	httpClient := http.Client{}
	addHTTPTransport(&httpClient)
	return httpClient
}

func GetStatusCode(ctx context.Context, httpClient *http.Client, url string) int {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, http.NoBody)
	Expect(err).NotTo(HaveOccurred())

	resp, err := httpClient.Do(req)
	Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()

	return resp.StatusCode
}

func NewNoAuthClient(endpoint string) (*gophercloud.ServiceClient, error) {
	serviceClient, err := noauth.NewBareMetalNoAuth(noauth.EndpointOpts{
		IronicEndpoint: endpoint,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot create an Ironic client: %w", err)
	}

	addHTTPTransport(&serviceClient.HTTPClient)
	return serviceClient, nil
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

	addHTTPTransport(&serviceClient.HTTPClient)
	return serviceClient, nil
}
