# Ironic Standalone Operator (IrSO) - AI Agent Instructions

Instructions for AI coding agents. For project overview, see [README.md](README.md).
For contribution guidelines, see [CONTRIBUTING.md](CONTRIBUTING.md).

## Repository Structure

| Directory | Purpose |
|-----------|---------|
| `api/v1alpha1/` | CRD types (Ironic) - separate Go module |
| `internal/controller/` | Reconciliation logic |
| `internal/webhook/` | Validation webhooks |
| `pkg/ironic/` | Core Ironic deployment logic |
| `config/` | Kustomize manifests (CRDs, RBAC, webhooks) |
| `hack/` | CI scripts (prefer Make targets locally) |
| `test/` | Functional tests - separate Go module |

## Testing Standards

CI uses GitHub Actions (`.github/workflows/`).
Run these locally before submitting PRs:

**Make targets:**

| Command | Purpose |
|---------|---------|
| `make test` | Full verification (generate + fmt + vet + unit) |
| `make generate` | Regenerate DeepCopy methods |
| `make manifests` | Regenerate CRDs, RBAC, webhooks, API docs |
| `make lint` | Go linting via golangci-lint (all modules) |
| `make modules` | Verify go.mod is tidy (all 3 modules) |

**Hack scripts** (auto-containerized, match CI exactly):

| Script | Purpose |
|--------|---------|
| `./hack/markdownlint.sh` | Markdown linting (config: `.markdownlint-cli2.yaml`) |
| `./hack/manifestlint.sh` | Kubernetes manifest validation (kubeconform) |

## Code Conventions

- **Go**: Linting rules in `.golangci.yaml`, license headers in `hack/boilerplate.go.txt`
- **Shell**: Use `set -eux`
- **Markdown**: Config in `.markdownlint-cli2.yaml`

## Key Workflows

### Modifying APIs

1. Edit `api/v1alpha1/*_types.go`
2. Run `make generate manifests`
3. Update webhooks in `internal/webhook/` if validation changes
4. Run `make test`

## Functional Testing

IrSO has self-contained functional tests on Minikube/Kind:

| Command | Purpose |
|---------|---------|
| `./test/prepare.sh` | Prepare test environment |
| `./test/run.sh` | Run all functional tests |
| `LABEL_FILTER="tls" ./test/run.sh` | Run specific tests by label |

Test labels: `no-params`, `tls`, `upgrade`, `ha`, `database`, `keepalived-dnsmasq`

## Code Review Guidelines

When reviewing pull requests:

1. **Security** - Hardcoded secrets, unpinned dependencies, missing input validation
2. **Test coverage** - New functionality should have tests
3. **Consistency** - Match existing patterns in the codebase
4. **Breaking changes** - Flag API/behavior changes affecting users

Focus on: internal/controller/, pkg/ironic/, api/, internal/webhook/.

## AI Agent Guidelines

### Before Changes

1. Run `make test` to verify baseline
2. Check patterns in similar existing files

### When Making Changes

1. Make minimal, surgical edits
2. Run `make generate manifests` after API changes
3. Run `make test` before committing
4. Add tests for new functionality

### Security Requirements

- Pin external dependencies by SHA (containers, GitHub Actions, binaries)
- No hardcoded credentials
- Validate all inputs

## Related Documentation

- [API Documentation](docs/api.md)
- [Release Process](docs/releasing.md)
- [Metal3 Book](https://book.metal3.io/irso/introduction)
