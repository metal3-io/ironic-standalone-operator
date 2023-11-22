# Deployed architecture

Ironic Operator provides two custom resources:

* [Ironic](config/crd/bases/metal3.io_ironics.yaml) describes the deployment
  of Ironic and its auxiliary [services](#services).
* [IronicDatabase](config/crd/bases/metal3.io_ironicdatabases.yaml) describes
  the deployment of MariaDB as the [storage](#storage) backend for Ironic.

## Services

An Ironic deployment always contains these three services:

* `ironic` is the main API service, as well as the conductor process that
  handles actions on bare-metal machines.
* `httpd` is the web server that serves images and configuration for iPXE and
  virtual media boot, as well as works as the HTTPS frontend for Ironic.
* `ramdisk-logs` is a script that unpacks any ramdisk logs and outputs them
  for consumption via `kubectl logs` or similar tools.

There is also a standard init container:

* `ipa-downloader` downloads images of the deployment/inspection ramdisk and
  stores them locally for easy access.

When network boot (iPXE) is enabled, another service is deployed:

* `dnsmasq` serves DHCP and functions as a PXE server for bootstrapping iPXE.

## Storage

Ironic requires a relational database to store its runtime data. By default,
a file-based SQLite database is configured. It achieves very good speed and low
footprint but is not persistent and cannot be used in an [HA setup](#ha-setup).
For these cases, it is possible to use a MariaDB instance by creating
an `IronicDatabase` object and linking it via the `databaseRef` field of an
`Ironic`. These way, three options are possible:

1. Non-HA + SQLite (the default)
2. Non-HA + MariaDB
3. HA + MariaDB

## Authentication and TLS

Ironic Operator does not allow installing services without authentication. If
the corresponding secrets are not provided, it will create ones with random
passwords and put links to them in the `credentialsRef` fields.

TLS is optional and must be enabled by providing a key+certificate pair via
a TLS secret linked in the `tlsRef` fields.

When TLS is used with a MariaDB database, ensure that Ironic can verify
the host certificate of the database. It will access a service-based URL
in the form of `<database name>-database.<namespace>.svc[.<cluster name>]`.
Here:

* `<database name>` is the name of the `IronicDatabase` object.
* `<namespace>` is the Kubernetes namespace of the both objects.
* `<cluster name>` is optional and can be passed to the operator via
  the `-cluster-name` CLI argument or the `CLUSTER_NAME` environment variable.

For example, a database called `ironic` in the namespace `test` will be
accessed as `ironic-database.test.svc` by default. If `CLUSTER_NAME` is set to
`example.com`, it will be `ironic-database.test.svc.example.com`. The TLS
certificate must be valid for this name.

## HA setup

HA setup is experimental and cannot be enabled by default (the validation
webhook will reject setting `distributed` to `true`).

The idea is to deploy ironic+httpd services in a daemon set on all control
plane nodes, while keeping dnsmasq (if enabled) on only one of them. There
are many unsolved issues with this setup - see
[ironic-operator#3](https://github.com/dtantsur/ironic-operator/issues/3).
