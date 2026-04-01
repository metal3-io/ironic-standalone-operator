# Run controller locally on kind

Use this flow when developing the controller while targeting a local Kubernetes
cluster created with [kind](https://kind.sigs.k8s.io/).

## Prerequisites

- `go`
- `kubectl`
- `kind`
- `make`

## Steps

1. Create a local cluster:

   ```bash
   kind create cluster --name irso-dev
   ```

1. Point kubectl to the kind context and confirm it:

   ```bash
   kubectl config use-context kind-irso-dev
   kubectl cluster-info
   ```

1. Install CRDs:

   ```bash
   make install
   ```

1. Run the controller locally (outside the cluster):

   ```bash
   make run RUNARGS="--webhook-port=0 --leader-elect=false"
   ```

1. (Optional) Apply the sample Ironic resource:

   ```bash
   kubectl apply -f config/samples/v1alpha1_ironic.yaml
   ```

1. Verify resources:

   ```bash
   kubectl get crd | grep ironics.ironic.metal3.io
   kubectl get ironic -A
   ```

   To stop the controller, press `Ctrl+C` in the terminal where `make run` is
running. 

1. Delete the kind cluster:

   ```bash
   kind delete cluster --name irso-dev
   ```
