# Ironic Standalone Operator (IrSO) - AI Coding Agent Instructions

This file provides comprehensive instructions for AI coding agents working on
the Ironic Standalone Operator project. It covers architecture, conventions,
tooling, CI/CD, and behavioral guidelines.

## Table of Contents

- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Development Workflows](#development-workflows)
- [Makefile Reference](#makefile-reference)
- [Hack Scripts Reference](#hack-scripts-reference)
- [CI/CD and GitHub Actions](#cicd-and-github-actions)
- [Code Patterns and Conventions](#code-patterns-and-conventions)
- [Testing Guidelines](#testing-guidelines)
- [Ironic Integration](#ironic-integration)
- [Integration Points](#integration-points)
- [Common Pitfalls](#common-pitfalls)
- [AI Agent Behavioral Guidelines](#ai-agent-behavioral-guidelines)

---

## Project Overview

The Ironic Standalone Operator (IrSO) is a Kubernetes operator for deploying
and managing [OpenStack Ironic](https://ironicbaremetal.org/) services for
Metal3. It provides declarative management of Ironic deployments via custom
resources, handling configuration, scaling, lifecycle management, version
upgrades, and high availability of Ironic components.

**Key URLs:**

- Repository: <https://github.com/metal3-io/ironic-standalone-operator>
- Container Image: `quay.io/metal3-io/ironic-standalone-operator`
- Documentation: <https://book.metal3.io/irso/introduction>
- API Reference: [docs/api.md](docs/api.md)

### Project Goals

IrSO aims to provide a **Kubernetes-native deployment mechanism for Ironic**:

1. **Declarative Configuration** - Define complete Ironic deployments via CRDs
   instead of manual manifests, Helm charts, or scripts.

2. **Version Management** - Supports multiple Ironic versions with automated
   upgrade handling including database migrations.

3. **High Availability** - Optional HA mode deploying Ironic as a DaemonSet on
   control plane nodes with shared MariaDB database and JSON RPC.

4. **Deep Ironic Integration** - Manages all Ironic components: API server,
   conductor, httpd (for images and virtual media), dnsmasq (DHCP/TFTP),
   keepalived (VIP management), and IPA ramdisk downloader.

5. **TLS Support** - Full TLS support for Ironic API, virtual media server,
   JSON RPC, and database connections.

6. **Production Ready** - Comprehensive functional test suite running on real
   Kubernetes clusters (Minikube, Kind) testing deployments, upgrades, HA, and
   various configurations.

---

## Architecture

### Core Components

1. **Custom Resources (CRDs)** - Located in `api/v1alpha1/`:
   - `Ironic` - Main resource defining an Ironic deployment
     - Configures Ironic API, Conductor, database connection
     - Manages dnsmasq, httpd, keepalived containers
     - Handles TLS certificates, API credentials, and networking
     - Supports version selection and upgrade orchestration

2. **Controllers** - Located in `internal/controller/`:
   - `IronicReconciler` - Manages complete Ironic deployment lifecycle
     - Creates/updates Deployments (single replica) or DaemonSets (HA)
     - Handles Services, Secrets, and Jobs (database migrations)
     - Manages version upgrades with pre/post-upgrade jobs
     - Orchestrates multi-container pods (ironic, httpd, dnsmasq, etc.)

3. **Ironic Package** - Located in `pkg/ironic/`:
   - `ironic.go` - Main deployment logic (Deployment/DaemonSet/Service)
   - `containers.go` - Container definitions and environment variables
   - `upgrades.go` - Database migration job handling
   - `validation.go` - Spec validation logic
   - `version.go` - Version management and image resolution
   - `secrets.go` - API credentials and htpasswd generation
   - `local.go` - Local (non-Kubernetes) Ironic runner support

4. **Webhooks** - Located in `internal/webhook/v1alpha1/`:
   - Validation webhooks for Ironic spec
   - Immutability checks for certain fields

### Resource Relationships

```text
Ironic (CRD)
   ↓
   ├── Secret (API credentials, htpasswd)
   ├── Secret (TLS certificate) [optional]
   ├── Secret (BMC CA certificate) [optional]
   ├── Database reference → MariaDB [optional, required for HA]
   │
   ├── Deployment (single replica mode)
   │   └── Pod
   │       ├── init: ramdisk-downloader [optional]
   │       ├── ironic (API + Conductor)
   │       ├── httpd (image server, TLS proxy)
   │       ├── ramdisk-logs
   │       ├── dnsmasq [optional, DHCP/TFTP]
   │       └── keepalived [optional, VIP]
   │
   └── DaemonSet (HA mode on control plane)
       └── Pods (one per control plane node)
           └── [same containers as above]
```

### Directory Structure

```text
ironic-standalone-operator/
├── api/v1alpha1/           # CRD type definitions (separate Go module)
│   ├── ironic_types.go     # Main Ironic spec/status definitions
│   ├── common.go           # Shared types (Labels, conditions)
│   ├── features.go         # Feature gates
│   └── version.go          # Version parsing and comparison
├── cmd/
│   ├── main.go             # Operator entrypoint
│   └── run-local-ironic/   # Local Ironic runner binary
├── config/                 # Kustomize manifests
│   ├── crd/bases/          # Generated CRD YAMLs
│   ├── default/            # Default deployment (with webhooks)
│   ├── manager/            # Controller manager deployment
│   ├── rbac/               # RBAC manifests
│   ├── samples/            # Example Ironic resources
│   └── webhook/            # Webhook configuration
├── docs/                   # Documentation
│   ├── api.md              # Generated API reference
│   └── releasing.md        # Release process documentation
├── hack/                   # Build and CI scripts
│   ├── boilerplate.go.txt  # License header template
│   └── tools/              # Tool dependencies
├── internal/
│   ├── controller/         # Controller implementations
│   └── webhook/            # Webhook implementations
├── pkg/ironic/             # Core Ironic deployment logic
├── releasenotes/           # Release notes (triggers releases)
└── test/                   # Functional tests (separate Go module)
    ├── helpers/            # Test helper functions
    ├── local-ironic/       # Local Ironic test scenarios
    └── suite_test.go       # Main functional test suite
```

### Multi-Module Structure

IrSO uses multiple Go modules:

- **Root module** (`go.mod`) - Main operator code
- **API module** (`api/go.mod`) - CRD types, can be imported independently
- **Test module** (`test/go.mod`) - Functional tests with additional dependencies

---

## Development Workflows

### Quick Start Commands

```bash
# Full verification (generate + fmt + vet + test)
make test

# Generate code and manifests after API changes
make generate manifests

# Run unit tests only
make test

# Run linters
make lint

# Build all binaries (manager + run-local-ironic)
make build

# Build manager binary only
make build-manager

# Verify go modules are tidy
make modules
```

### Local Development

```bash
# Run controller locally (requires KUBECONFIG)
make run

# Build and run local Ironic (without Kubernetes)
make build-run-local-ironic
./bin/run-local-ironic --help
```

### Docker Build

```bash
# Build container image
make docker-build IMG=quay.io/metal3-io/ironic-standalone-operator:dev

# Build with debug symbols
make docker-build-debug
```

### Deployment

```bash
# Install CRDs into cluster
make install

# Deploy operator to cluster
make deploy

# Undeploy operator
make undeploy

# Uninstall CRDs
make uninstall
```

---

## Makefile Reference

### Testing Targets

| Target | Description |
|--------|-------------|
| `make test` | Run generate + fmt + vet + unit tests |

### Build Targets

| Target | Description |
|--------|-------------|
| `make build` | Build all binaries (manager + run-local-ironic) |
| `make build-manager` | Build manager binary to `bin/manager` |
| `make build-run-local-ironic` | Build local Ironic runner |

### Code Generation Targets

| Target | Description |
|--------|-------------|
| `make generate` | Generate DeepCopy methods |
| `make manifests` | Generate CRDs, RBAC, webhooks, and API docs |
| `make fmt` | Run go fmt on all modules |
| `make vet` | Run go vet on all modules |

### Linting Targets

| Target | Description |
|--------|-------------|
| `make lint` | Run golangci-lint on all modules |

### Module Management

| Target | Description |
|--------|-------------|
| `make modules` | Run go mod tidy and verify on all modules |

### Deployment Targets

| Target | Description |
|--------|-------------|
| `make install` | Install CRDs into cluster |
| `make uninstall` | Uninstall CRDs from cluster |
| `make deploy` | Deploy operator to cluster |
| `make undeploy` | Undeploy operator from cluster |

### Docker Targets

| Target | Description |
|--------|-------------|
| `make docker-build` | Build docker image |
| `make docker-build-debug` | Build docker image with debug symbols |
| `make docker-buildx` | Build multi-platform image |

### Release Targets

| Target | Description |
|--------|-------------|
| `make release` | Full release process (requires RELEASE_TAG) |
| `make release-manifests` | Build release manifests to `out/` |
| `make release-notes` | Generate release notes (requires RELEASE_TAG) |

### Utility Targets

| Target | Description |
|--------|-------------|
| `make go-version` | Print Go version used for builds |
| `make help` | Display all available targets |

---

## Hack Scripts Reference

Scripts in `hack/` support both local execution and containerized CI runs.

### Available Scripts

| Script | Purpose |
|--------|---------|
| `gen-api-doc.sh` | Generate API documentation (docs/api.md) |
| `gomod.sh` | Verify go.mod files are tidy |
| `markdownlint.sh` | Lint markdown files (containerized) |
| `manifestlint.sh` | Validate Kubernetes manifests with kubeconform |
| `verify-release.sh` | Comprehensive release verification |
| `ensure-go.sh` | Ensure correct Go version is installed |

### Container Execution Pattern

Scripts support containerized execution via environment variables:

```bash
IS_CONTAINER="${IS_CONTAINER:-false}"
CONTAINER_RUNTIME="${CONTAINER_RUNTIME:-podman}"

if [ "${IS_CONTAINER}" != "false" ]; then
    # Run the actual logic
else
    "${CONTAINER_RUNTIME}" run --rm \
        --env IS_CONTAINER=TRUE \
        --volume "${PWD}:/workdir:ro,z" \
        <image> /workdir/hack/<script>.sh "$@"
fi
```

---

## CI/CD and GitHub Actions

IrSO uses GitHub Actions for CI/CD (not Prow like other Metal3 projects).
Workflows are defined in `.github/workflows/`.

### Pull Request Workflows

| Workflow | Trigger | What it does |
|----------|---------|--------------|
| `build.yml` | PR | Unit tests, container build, manifest verification |
| `golangci-lint.yml` | PR | Go linting for root, api, and test modules |
| `functional.yml` | PR | Functional tests on Minikube (regular, upgrade, HA) |
| `local-ironic.yml` | PR | Local Ironic tests (empty, provnet scenarios) |
| `pr-verifier.yaml` | PR | PR content verification (uses project-infra) |
| `pr-link-check.yml` | PR | Check for broken links |

### Functional Test Matrix

The `functional.yml` workflow runs three test configurations:

| Name | Label Filter | Minikube Args | Tests |
|------|--------------|---------------|-------|
| `regular` | `!upgrade && !ha` | (default) | Basic, TLS, credentials |
| `upgrade` | `upgrade && !ha` | (default) | Version upgrades |
| `ha` | `ha` | `--ha` | HA with MariaDB |

### Release Workflow

The `release.yaml` workflow is triggered when release notes are merged to main:

1. **push_release_tags** - Creates git tags (`v0.x.y`, `api/v0.x.y`, `test/v0.x.y`)
2. **release** - Creates draft GitHub release with artifacts
3. **build_irso** - Builds and pushes container image to Quay

### Other Workflows

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `scheduled-link-check.yml` | Schedule | Weekly broken link check |
| `scheduled-osv-scan.yml` | Schedule | Vulnerability scanning |
| `build-images-action.yml` | Manual/Release | Container image builds |
| `dependabot.yml` | Schedule | Dependency updates |

---

## Code Patterns and Conventions

### Go Code Style

Go code follows golangci-lint rules in `.golangci.yaml`:

- Import aliasing is enforced (see `importas` section in `.golangci.yaml`)
- License headers are required (see `hack/boilerplate.go.txt`)
- Go version is specified in `Makefile` (`GO_VERSION`)

### YAML Conventions

Kubernetes manifests use 2-space indentation (Kubernetes standard).
Validated by kubeconform via `hack/manifestlint.sh`.

### Shell Script Conventions

**Required settings at script start:**

```bash
#!/bin/sh
set -eux
```

### Markdown Conventions

Configuration in `.markdownlint-cli2.yaml`. Auto-generated files
(`docs/api.md`) are excluded from linting.

### API Modifications Workflow

1. Edit types in `api/v1alpha1/*_types.go`
2. Run `make generate manifests` to regenerate:
   - DeepCopy functions
   - CRDs in `config/crd/bases/`
   - RBAC in `config/rbac/`
   - API documentation in `docs/api.md`
3. Update webhooks if validation changes
4. Run `make test` to verify

### Controller Patterns

**Reconciliation flow in `IronicReconciler`:**

```go
func (r *IronicReconciler) Reconcile(ctx, req) (ctrl.Result, error) {
    // 1. Fetch Ironic resource
    // 2. Handle deletion (finalizers)
    // 3. Validate and resolve version
    // 4. Ensure API credentials secret
    // 5. Load TLS and BMC CA secrets
    // 6. Call pkg/ironic.EnsureIronic() for deployment
    // 7. Update status conditions
}
```

**EnsureIronic flow in `pkg/ironic`:**

```go
func EnsureIronic(cctx, resources) (Status, error) {
    // 1. Validate spec
    // 2. Run pre-upgrade job (if database configured)
    // 3. Create/update Deployment or DaemonSet
    // 4. Create/update Service
    // 5. Run post-upgrade job (if database configured)
}
```

### Status Conditions

IrSO uses standard Kubernetes conditions:

- `Ready` - Ironic is fully deployed and operational
- Reasons: `InProgress`, `Ready`, `Failed`

### Version Management

Supported Ironic versions are defined in
[`api/v1alpha1/ironic_types.go`](api/v1alpha1/ironic_types.go) in the
`SupportedVersions` map. This mapping must be updated when adding support for
new Ironic versions.

---

## Testing Guidelines

### Unit Tests

Unit tests use standard Go testing in each package. Run with:

```bash
make test
```

### Functional Tests

Functional tests in `test/` use Ginkgo/Gomega and test against real clusters:

```go
var _ = Describe("Ironic object tests", func() {
    It("creates Ironic without any parameters", Label("no-params"), func() {
        // Test implementation
    })
})
```

### Running Functional Tests

```bash
# Prepare test environment (on Kind/Minikube)
./test/prepare.sh

# Run all tests
./test/run.sh

# Run specific tests by label
LABEL_FILTER="!upgrade && !ha" ./test/run.sh

# Collect logs
./test/collect-logs.sh
```

### Test Labels

| Label | Description |
|-------|-------------|
| `no-params` | Basic Ironic without configuration |
| `api-secret` | Custom API credentials |
| `tls` | TLS certificate configuration |
| `upgrade` | Version upgrade tests |
| `ha` | High availability tests |
| `database` | MariaDB integration |
| `keepalived-dnsmasq` | Network boot services |
| `extra-config` | Custom Ironic configuration |
| `disabled-downloader` | Without IPA downloader |

### Local Ironic Tests

Tests in `test/local-ironic/` verify the `run-local-ironic` binary:

```bash
./test/local-ironic/run.sh empty
./test/local-ironic/run.sh provnet
./test/local-ironic/tear-down.sh <scenario>
```

### Test Environment Variables

| Variable | Description |
|----------|-------------|
| `IRONIC_CERT_FILE` | TLS certificate file (required) |
| `IRONIC_KEY_FILE` | TLS key file (required) |
| `LOGDIR` | Directory for test logs |
| `IRONIC_CUSTOM_IMAGE` | Custom Ironic image |
| `IRONIC_CUSTOM_VERSION` | Custom Ironic version |
| `CLUSTER_TYPE` | `kind` or `minikube` |

---

## Ironic Integration

IrSO is tightly integrated with [OpenStack Ironic](https://ironicbaremetal.org/)
and the [ironic-image](https://github.com/metal3-io/ironic-image) container.

### Container Images Used

| Image | Purpose |
|-------|---------|
| `quay.io/metal3-io/ironic` | Main Ironic (API, Conductor, httpd) |
| `quay.io/metal3-io/ironic-ipa-downloader` | IPA ramdisk downloader |
| `quay.io/metal3-io/keepalived` | VIP management (optional) |

### Ironic Components Managed

- **Ironic API** - REST API for bare metal management
- **Ironic Conductor** - Orchestrates provisioning workflows
- **httpd** - Serves IPA images and virtual media
- **dnsmasq** - DHCP and TFTP for PXE boot (optional)
- **Keepalived** - Virtual IP management (optional)

### Database Support

- **SQLite** (default) - Ephemeral, state lost on pod restart
- **MariaDB** (recommended for HA) - Persistent, required for multi-replica

### Ironic Version Branches

IrSO supports specific Ironic versions mapped to
[ironic-image](https://github.com/metal3-io/ironic-image) branches. The
version-to-branch mapping is maintained in
[`api/v1alpha1/ironic_types.go`](api/v1alpha1/ironic_types.go) (`SupportedVersions`).

For Ironic API version history, see the
[Ironic API versions listing](https://docs.openstack.org/ironic/latest/contributor/webapi-version-history.html).

### Upgrade Handling

When a database is configured, IrSO runs database migration jobs:

1. **Pre-upgrade job** - `ironic-dbsync --config-file /etc/ironic/ironic.conf upgrade`
2. **Deployment/DaemonSet update** - Rolling update of Ironic pods
3. **Post-upgrade job** - Online data migrations

---

## Integration Points

### With Baremetal Operator (BMO)

- IrSO is an alternative to BMO's built-in Ironic deployment
- BMO can discover and use Ironic deployed by IrSO
- More declarative and Kubernetes-native than manual manifests

### With CAPM3

- CAPM3 (Cluster API Provider Metal3) can use Ironic deployed by IrSO
- IrSO handles Ironic lifecycle independently of cluster provisioning

### Standalone Use

- Can deploy Ironic without BMO or CAPM3
- Useful for non-Cluster API bare metal management
- Simplified Ironic deployment for development/testing

### E2E Testing

**Important:** IrSO has its own comprehensive functional test suite running via
GitHub Actions. Unlike IPAM (which relies on CAPM3 for e2e), IrSO tests are
self-contained and run on Minikube/Kind clusters.

The tests verify:

- Basic deployment scenarios
- TLS and authentication
- Version upgrades (with and without database)
- High availability with MariaDB
- DHCP/Keepalived configurations
- Custom Ironic configuration options

---

## Common Pitfalls

1. **Secret Dependencies** - Ironic needs API credentials and TLS secrets before
   starting. The controller waits for secrets and reports clear error messages.

2. **Network Configuration** - Interface names vary by environment. Use
   `networking.macAddresses` if interface names are unreliable.

3. **Database Initialization** - MariaDB must be ready before Ironic starts.
   Use MariaDB Operator for reliable database management.

4. **Port Conflicts** - Default ports (6385, 6180, 6183, 6189) may conflict.
   All ports are configurable via the Ironic spec.

5. **HA Requirements** - High availability requires:
   - Database (MariaDB) configuration
   - Multiple control plane nodes
   - Feature gate `HighAvailability` enabled

6. **Downgrade Prevention** - Downgrades with an external database are blocked
   as Ironic doesn't support database schema rollback.

7. **Forgetting `make generate manifests`** - After API changes, always
   regenerate code and manifests.

8. **Multi-Module Structure** - Remember to run `make modules` to tidy all
   three Go modules (root, api, test).

---

## AI Agent Behavioral Guidelines

### Critical Rules

1. **Use single commands** - Do not concatenate multiple commands with `&&` or
   `;` to avoid interactive permission prompts from the user. Run one command
   at a time.

2. **Be strategic with output filtering** - Use `head`, `tail`, or `grep` when
   output is clearly excessive (e.g., large logs), but prefer full output for
   smaller results to avoid missing context.

3. **Challenge assumptions** - Do not take user statements as granted. If you
   have evidence against an assumption, present it respectfully with
   references.

4. **Search for latest versions** - When suggesting dependencies, libraries,
   or tools, always verify and use the latest stable versions.

5. **Security first** - Take security very seriously. Review code for:
   - Hardcoded credentials
   - Insecure defaults
   - Missing input validation
   - Privilege escalation risks

6. **Pin dependencies by SHA** - All external dependencies must be SHA pinned
   when possible (container images, GitHub Actions, downloaded binaries).
   This prevents supply chain attacks and ensures reproducible builds.

7. **Provide references** - Back up suggestions with links to documentation,
   issues, PRs, or code examples from the repository.

8. **Follow existing conventions** - Match the style of existing code:
   - Shell scripts: use patterns from `hack/` scripts
   - Go code: follow golangci-lint rules in `.golangci.yaml`
   - Markdown: follow `.markdownlint-cli2.yaml` rules
   - License headers: use templates from `hack/boilerplate.go.txt`

### Before Making Changes

1. Run `make lint` to understand current linting rules
2. Run `make test` to verify baseline test status
3. Check existing patterns in similar files
4. Verify Go version matches `Makefile` (`GO_VERSION` variable)

### When Modifying Code

1. Make minimal, surgical changes
2. Run `make generate manifests` after API changes
3. Run `make test` before submitting
4. Update documentation if behavior changes
5. Add tests for new functionality

### When Debugging CI Failures

1. Check GitHub Actions workflow definitions in `.github/workflows/`
2. Run the same commands locally (e.g., `make lint`, `make test`)
3. For functional tests, check test labels and filters
4. Use the exact container images from CI when possible

### Commit Guidelines

- Sign commits with `-s` flag (DCO required)
- Use conventional commit prefixes:
  - ✨ `:sparkles:` - New feature
  - 🐛 `:bug:` - Bug fix
  - 📖 `:book:` - Documentation
  - 🌱 `:seedling:` - Other changes
  - ⚠️ `:warning:` - Breaking changes
  - 🚀 `:rocket:` - Release

---

## Git and Release Information

- **Branches**: `main` (development), `release-X.Y` (stable releases)
- **DCO Required**: All commits must be signed off (`git commit -s`)
- **PR Labels**: ⚠️ breaking, ✨ feature, 🐛 bug, 📖 docs, 🌱 other
- **Release Process**: See [docs/releasing.md](./docs/releasing.md)
- **Tags**: Primary (`v0.x.y`), API module (`api/v0.x.y`), Test module (`test/v0.x.y`)

---

## Additional Resources

### Related Projects

- [Ironic](https://ironicbaremetal.org/) - OpenStack bare metal provisioning
- [ironic-image](https://github.com/metal3-io/ironic-image) - Container images
- [BMO](https://github.com/metal3-io/baremetal-operator) - Baremetal Operator
- [CAPM3](https://github.com/metal3-io/cluster-api-provider-metal3) - Cluster
  API Provider Metal3
- [Metal3 Docs](https://book.metal3.io) - Project documentation

### Issue Tracking

- Issues: <https://github.com/metal3-io/ironic-standalone-operator/issues>
- Good first issues: [good first issue label](https://github.com/metal3-io/ironic-standalone-operator/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22)

---

## Quick Reference Card

```bash
# Most common commands
make test              # Full verification (generate + fmt + vet + test)
make lint              # Linting only
make generate manifests # Regenerate code and manifests
make modules           # Tidy all go modules
make build             # Build all binaries
make docker-build      # Build container

# Development
make run               # Run controller locally
make install           # Deploy CRDs to cluster
make deploy            # Deploy operator to cluster

# Functional tests (requires cluster)
./test/prepare.sh      # Prepare test environment
./test/run.sh          # Run functional tests
LABEL_FILTER="tls" ./test/run.sh  # Run specific tests
./test/collect-logs.sh # Collect test logs

# Local Ironic (no Kubernetes)
make build-run-local-ironic
./bin/run-local-ironic --help

# Hack scripts (containerized)
./hack/markdownlint.sh # Markdown linting
./hack/manifestlint.sh # K8s manifest validation
./hack/gomod.sh        # Verify go modules

# Release verification
RELEASE_TAG=v0.x.y ./hack/verify-release.sh
```
