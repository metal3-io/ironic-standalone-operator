# Functional tests for ironic-standalone-operator

This directory contains functional tests for IrSO that install Ironic in
different configurations and make sure it comes up correctly.

## Setup

The tests assume a Kubernetes cluster with ironic-standalone-operator is
available. A helper script `test/prepare.sh` can be used to configure the
operator and its dependencies on a Kind cluster. Then you can use `test/run.sh`
to run the tests themselves. Finally, `test/collect-logs.sh` can be used to
store the logs to `LOGDIR` (see below).

## Environment variables

Required:

- `IRONIC_CERT_FILE` - TLS certificate file to use for Ironic API
- `IRONIC_KEY_FILE` - private key file of the TLS certificate

Optional:

- `LOGDIR` - directory where logs will be placed both runing the test run and
  by `collect-logs.sh`
- `IRONIC_CUSTOM_IMAGE` - Ironic container image to use when testing
- `IRONIC_CUSTOM_VERSION` - Ironic version to use when testing
- `MARIADB_CUSTOM_IMAGE` - MariaDB container image to use with tests that
  create a database
