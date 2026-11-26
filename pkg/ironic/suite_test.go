package ironic

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

var (
	testEnv     *envtest.Environment
	testEnvOnce sync.Once
	errTestEnv  error
	testClient  client.Client
	testScheme  = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(testScheme))
	utilruntime.Must(metal3api.AddToScheme(testScheme))
}

func TestMain(m *testing.M) {
	code := m.Run()
	if testEnv != nil {
		if err := testEnv.Stop(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to stop envtest: %v\n", err)
			code = 1
		}
	}
	os.Exit(code)
}

// newEnvtestClient returns a client for a real API server (and etcd) shared by
// all tests in this package. The server is started on first use. The test is
// skipped if KUBEBUILDER_ASSETS is not set, `make unit` takes care of that.
func newEnvtestClient(t *testing.T) client.Client {
	t.Helper()

	if os.Getenv("KUBEBUILDER_ASSETS") == "" {
		t.Skip("KUBEBUILDER_ASSETS is not set, run the tests via `make unit`")
	}

	testEnvOnce.Do(func() {
		testEnv = &envtest.Environment{Scheme: testScheme}
		cfg, err := testEnv.Start()
		if err != nil {
			errTestEnv = err
			return
		}
		testClient, errTestEnv = client.New(cfg, client.Options{Scheme: testScheme})
	})
	require.NoError(t, errTestEnv)

	return testClient
}
