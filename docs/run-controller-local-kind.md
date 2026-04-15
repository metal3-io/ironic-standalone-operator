# Run controller locally on kind

Use this flow when you want a **fast feedback loop** while developing the
controller (e.g. debugging reconciliation logic, iterating on controller code,
watching logs/events) against a local Kubernetes cluster created with
[kind](https://kind.sigs.k8s.io/).

For full end-to-end coverage, use the functional test harness
(`./test/prepare.sh` + `./test/run.sh`). That path deploys IrSO into the cluster
and also installs extra dependencies (e.g. cert-manager, MariaDB operator),
which is great for E2E but typically much slower to iterate on.

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
