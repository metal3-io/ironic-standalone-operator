# Testing

This document describes how to test IrSO locally before submitting a pull
request. All checks listed under [mandatory](#mandatory-checks) are enforced
by CI and must pass before a PR can be merged.

## Container runtime

Several commands require a container runtime. Set `CONTAINER_RUNTIME` to
`podman` or `docker` depending on what is available in your environment:

```bash
export CONTAINER_RUNTIME=podman
```

## Mandatory checks

Run these before every pull request:

```bash
make test          # generate + fmt + vet + unit tests
make lint          # golangci-lint across all modules
```

`make test` regenerates auto-generated code and manifests, checks formatting,
runs `go vet`, and executes the unit test suite.

If you have changed anything under `api/`, you can separately update the
manifests with:

```bash
make generate manifests
```

and commit the resulting changes. CI will verify that generated files are
up-to-date.

Additionally, the build process requires `go.mod`/`go.sum` to be tidy
across all three Go modules. Run this locally if you changed dependencies:

```bash
make modules
```

## Optional linters

These linters run in CI but are only relevant when you modify the
corresponding file types:

| Command | When to run |
|---------|-------------|
| `./hack/markdownlint.sh` | After editing Markdown files |
| `./hack/manifestlint.sh` | After editing Kubernetes manifests under `config/` |
| `./hack/shellcheck.sh` | After editing shell scripts |

The hack scripts run inside containers automatically, so they match the CI
environment exactly.

CI also runs **yamllint** and **link checking** on pull requests; these use
shared workflows and do not have local equivalents.

## Functional tests

Functional tests deploy the operator on a local Kubernetes cluster (Minikube
or Kind) and exercise it end-to-end. They are **not required** for every PR
but are useful when you change reconciliation logic or deployment manifests.
CI runs them automatically.

### Prerequisites

- Minikube or Kind
- A container runtime (podman or docker)
- Go (matching the version in the Makefile)
- Helm (downloaded automatically by `prepare.sh`)

### Running

```bash
# Set up dependencies (cert-manager, MariaDB operator, build and deploy IrSO)
./test/prepare.sh

# Run all functional tests
./test/run.sh
```

### Running a subset of tests

Use the `LABEL_FILTER` environment variable to select tests by label:

```bash
LABEL_FILTER="tls" ./test/run.sh
```

To see all available labels, search for `Label(` in `test/suite_test.go`.

### Configuration

| Variable | Default | Purpose |
|----------|---------|---------|
| `LABEL_FILTER` | *(empty -- run all)* | Ginkgo label filter expression |
| `LOGDIR` | `/tmp/logs` | Where logs and JUnit reports are written |
| `TEST_TIMEOUT` | `90m` | Go test timeout |
| `IMG` | `localhost/controller:test` | Operator image to build and deploy |

After a test run, collect diagnostic logs with:

```bash
./test/collect-logs.sh
```

### Extending the test suite

Functional tests live in `test/suite_test.go` and use
[Ginkgo v2](https://onsi.github.io/ginkgo/) with
[Gomega](https://onsi.github.io/gomega/). Helper utilities are in
`test/helpers/`.

#### Test structure

All tests are defined inside a single `Describe("Ironic object tests")` block.
Each `It` block must have a unique **label** as its first `Label(...)` argument.
The `BeforeEach` creates a Kubernetes namespace named `test-<first-label>`,
so every test runs in isolation.

A typical test follows this pattern:

```go
It("creates Ironic with my feature", Label("my-feature"), func() {
    name := types.NamespacedName{
        Name:      "test-ironic",
        Namespace: namespace,
    }

    ironic := helpers.NewIronic(ctx, k8sClient, name, metal3api.IronicSpec{
        // Set the fields relevant to your feature
    })
    DeferCleanup(func() {
        CollectLogs(namespace)
        DeleteAndWait(ironic)
    })

    ironic = WaitForIronic(name)
    VerifyIronic(ironic, TestAssumptions{
        // Enable the checks relevant to your feature
    })
})
```

#### Key helpers

| Function / Type | Location | Purpose |
|-----------------|----------|---------|
| `helpers.NewIronic` | `test/helpers/resources.go` | Create an Ironic CR with the given spec |
| `helpers.NewAuthSecret` | `test/helpers/resources.go` | Create a BasicAuth secret |
| `helpers.NewTLSSecret` | `test/helpers/resources.go` | Create a TLS secret |
| `helpers.CreateDatabase` | `test/helpers/database.go` | Set up MariaDB database, user, and grant |
| `WaitForIronic` | `test/suite_test.go` | Poll until the Ironic CR reaches Ready |
| `WaitForUpgrade` | `test/suite_test.go` | Poll until an upgrade completes |
| `WaitForIronicFailure` | `test/suite_test.go` | Poll until the Ready condition reports a specific error |
| `VerifyIronic` | `test/suite_test.go` | Run verification checks against a running Ironic |
| `DeleteAndWait` | `test/suite_test.go` | Delete an Ironic CR and wait for it to be gone |
| `CollectLogs` | `test/suite_test.go` | Save pod logs and resource YAML to `$LOGDIR` |
| `TestAssumptions` | `test/suite_test.go` | Struct controlling which checks `VerifyIronic` runs |

#### Adding a new test

1. Add a new `It` block inside the existing `Describe` in
   `test/suite_test.go`. Give it a descriptive name and a unique label.
1. If your test needs new verification logic, add fields to the
   `TestAssumptions` struct and handle them in `VerifyIronic` (or in
   dedicated `verify*` functions).
1. If you need reusable resource creation (new Secret types, ConfigMaps,
   etc.), add helpers to `test/helpers/`.
1. Run your test in isolation first:

   ```bash
   LABEL_FILTER="my-feature" ./test/run.sh
   ```
