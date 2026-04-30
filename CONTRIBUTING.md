# How to Contribute to Ironic Standalone Operator

> **Note**: Please read the [common Metal3 contributing guidelines](https://github.com/metal3-io/community/blob/main/CONTRIBUTING.md)
> first.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Versioning](#versioning)
   - [Codebase](#codebase)
- [Branches](#branches)
   - [Maintenance and Guarantees](#maintenance-and-guarantees)
- [Backporting Policy](#backporting-policy)
- [Release Process](#release-process)
   - [Exact Steps](#exact-steps)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Versioning

See the [common versioning and release semantics](https://github.com/metal3-io/community/blob/main/CONTRIBUTING.md#versioning-and-release-semantics)
in the Metal3 community contributing guide.

**Note**: The test module and experiments do not provide any backward
compatible guarantees.

### Codebase

Ironic Standalone Operator doesn't follow the release cadence of upstream Kubernetes.
The versioning semantics follow the common Metal3 guidelines above.

## Branches

See the [common branch structure guidelines](https://github.com/metal3-io/community/blob/main/CONTRIBUTING.md#branches)
in the Metal3 community contributing guide.

### Maintenance and Guarantees

Ironic Standalone Operator supports the most recent release for all supported
APIs and contract versions. Support in this section refers to CI support and
the ability to backport and release patch versions; the
[backport policy](#backporting-policy) is defined below.

## Backporting Policy

See the [common backporting guidelines](https://github.com/metal3-io/community/blob/main/CONTRIBUTING.md#backporting)
in the Metal3 community contributing guide.

Additionally, for Ironic Standalone Operator:

- We generally do not accept backports to Ironic Standalone Operator release
  branches that are EOL. Check the [Version support](https://github.com/metal3-io/metal3-docs/blob/main/docs/user-guide/src/version_support.md#ironic-standalone-operator)
  guide for reference.

## Release Process

See the [common release process guidelines](https://github.com/metal3-io/community/blob/main/CONTRIBUTING.md#release-process)
in the Metal3 community contributing guide.

### Exact Steps

Refer to the [releasing document](./docs/releasing.md) for the exact steps.
