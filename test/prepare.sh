#!/bin/bash

set -eux -o pipefail

IMG="${IMG:-controller:test}"

for image in ironic mariadb ironic-ipa-downloader; do
    docker pull "quay.io/metal3-io/${image}"
    kind load docker-image "quay.io/metal3-io/${image}"
done

make docker-build IMG="${IMG}"
kind load docker-image "${IMG}"
make install deploy IMG="${IMG}"

kubectl wait --for=condition=Available --timeout=60s \
    -n ironic-standalone-operator-system deployment/ironic-standalone-operator-controller-manager
