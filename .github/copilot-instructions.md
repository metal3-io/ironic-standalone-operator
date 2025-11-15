# Ironic Standalone Operator - AI Coding Assistant Instructions

## Project Overview

Kubernetes operator for deploying and managing standalone Ironic
services. Provides declarative management of Ironic deployments via
custom resources, handling configuration, scaling, and lifecycle
management of Ironic components.

## Architecture

### Custom Resources (CRDs)

Located in `api/v1alpha1/`:

- `Ironic` - Main resource defining an Ironic deployment
  - Configures Ironic API, Conductor, database
  - Manages dnsmasq, httpd, and related services
  - Handles secrets, ConfigMaps, and networking

### Controllers

Located in `internal/controller/`:

- `IronicReconciler` - Manages Ironic deployment lifecycle
  - Creates/updates Deployments, Services, ConfigMaps
  - Handles Ironic configuration rendering
  - Manages database initialization
  - Orchestrates multi-container pods

## Key Features

- **Declarative Configuration** - Define Ironic via CRD instead of manual manifests
- **Configuration Management** - Auto-generates Ironic configs from spec
- **Secret Management** - Handles BMC credentials, database passwords
- **Database Support** - SQLite or MariaDB backends
- **Networking** - Configures provisioning and external networks
- **Upgrades** - Manages Ironic version upgrades

## Development Workflows

```bash
# Generate manifests after API changes
make manifests

# Run unit tests
make test

# Run linters
make lint

# Build manager binary
make build-manager

# Build container image
make docker-build IMG=quay.io/metal3-io/ironic-standalone-operator:dev

# Run locally (against cluster)
make run

# Deploy to cluster
make deploy

# Run local Ironic using operator's tooling
make build-run-local-ironic
./bin/run-local-ironic
```

**Ironic Resource Spec:**

```yaml
apiVersion: ironic.metal3.io/v1alpha1
```

## Code Patterns

**Ironic Resource Spec:**

```yaml
apiVersion: ironic.metal3.io/v1alpha1
kind: Ironic
metadata:
  name: ironic
  namespace: metal3
spec:
  databaseRef:
    name: mariadb-connection
  networking:
    interface: ens3
    dhcpRange: 172.22.0.10,172.22.0.100
  images:
    ipa:
      kernel: http://example.com/ipa.kernel
      initramfs: http://example.com/ipa.initramfs
```

**Controller Pattern:**

- Reconcile creates/updates all Ironic components
- Status reflects deployment state (Ready, Degraded, etc.)
- Handles finalizers for cleanup
- Uses owner references for garbage collection

## Integration Points

### With BMO

- Alternative to BMO's built-in Ironic deployment
- BMO can use Ironic deployed by this operator
- More declarative than manual manifests

### Standalone Use

- Can deploy Ironic without BMO
- Useful for non-CAPI use cases
- Simplified Ironic management

## Key Files

- `main.go` - Operator entrypoint
- `api/v1alpha1/ironic_types.go` - Ironic CRD definition
- `internal/controller/ironic_controller.go` - Main reconciliation logic
- `config/` - Deployment manifests and CRDs
- `tools/run_local_ironic.sh` - Standalone Ironic runner

## Common Pitfalls

1. **Resource Dependencies** - Ironic needs ConfigMaps, Secrets before
   starting
2. **Network Configuration** - Interface names vary by environment
3. **Database Initialization** - MariaDB must be ready before Ironic
   starts
4. **Image Availability** - IPA images must be accessible from Ironic
5. **Port Conflicts** - Ensure no port conflicts with existing services

## Deployment Example

```bash
# Install CRDs
make install

# Deploy operator
make deploy

# Create Ironic instance
kubectl apply -f config/samples/ironic_v1alpha1_ironic.yaml

# Check status
kubectl get ironic -n metal3
kubectl describe ironic ironic -n metal3
```

## Status

This operator provides a more Kubernetes-native way to manage Ironic
compared to manual manifests. It's in active development and may be the
preferred deployment method in future Metal3 releases.

## Version Compatibility

- Compatible with Ironic 2023.1 (Antelope) and newer
- Requires Kubernetes 1.25+
- Works with BMO v0.5.0+
