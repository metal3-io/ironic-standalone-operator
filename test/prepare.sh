#!/bin/bash

set -eux -o pipefail

IMG="${IMG:-localhost/controller:test}"
LOGDIR="${LOGDIR:-/tmp/logs}"
CONTAINER_RUNTIME="${CONTAINER_RUNTIME:-}"

mkdir -p "${LOGDIR}"

if [ -z "${CONTAINER_RUNTIME}" ]; then
    if podman version > /dev/null 2>&1;  then
        CONTAINER_RUNTIME=podman
    else
        CONTAINER_RUNTIME=docker
    fi
fi

kind_load() {
    local image="$1"
    if [ "${CONTAINER_RUNTIME}" = podman ]; then
        local archive="$(mktemp --suffix=.tar)"
        podman save "${image}" > "${archive}"
        kind load image-archive -v 2 "${archive}"
        rm -f "${archive}"
    else
        kind load docker-image -v 2 "${image}"
    fi
}

for image in ironic mariadb ironic-ipa-downloader; do
    "${CONTAINER_RUNTIME}" pull "quay.io/metal3-io/${image}"
    kind_load "quay.io/metal3-io/${image}"
done

"${CONTAINER_RUNTIME}" build -t "${IMG}" . 2>&1 | tee "${LOGDIR}/docker-build.log"
kind_load "${IMG}"
make install deploy IMG="${IMG}"

kubectl wait --for=condition=Available --timeout=60s \
    -n ironic-standalone-operator-system deployment/ironic-standalone-operator-controller-manager
