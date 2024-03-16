#!/bin/bash

set -x

IMG="${IMG:-controller:test}"

clean_up() {
    cd "$(dirname $0)/.."
    if [[ -z "${NOCLEAN:-}" ]]; then
        make undeploy uninstall
    fi
}

on_exit() {
    EXIT=$?
    set +e
    clean_up
    exit $EXIT
}

clean_up
trap on_exit EXIT

set -eu -o pipefail

make docker-build IMG="${IMG}"
kind load docker-image "${IMG}"
make install deploy IMG="${IMG}"

kubectl wait --for=condition=Available --timeout=60s \
    -n ironic-standalone-operator-system deployment/ironic-standalone-operator-controller-manager
cd test && go test
