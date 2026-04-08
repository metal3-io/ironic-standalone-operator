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

Additionally, CI runs `make modules` to ensure `go.mod`/`go.sum` are tidy
across all three Go modules. Run it locally if you changed dependencies:

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

Available labels: `no-params`, `tls`, `upgrade`, `ha`, `database`,
`keepalived-dnsmasq`.

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
