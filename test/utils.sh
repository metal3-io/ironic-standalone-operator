#!/usr/bin/env bash

set -eux -o pipefail

CLUSTER_TYPE="${CLUSTER_TYPE:-kind}"

if [[ -z "${CONTAINER_RUNTIME}" ]]; then
    if command -v podman &> /dev/null;  then
        CONTAINER_RUNTIME=podman
    else
        CONTAINER_RUNTIME=docker
    fi
fi

IMAGE_NAMESPACE="${IMAGE_NAMESPACE:-quay.io/metal3-io}"

image_load() {
    local image="$1"
    local archive
    archive="$(mktemp --suffix=.tar)"
    "${CONTAINER_RUNTIME}" save "${image}" > "${archive}"
    if [[ "${CLUSTER_TYPE}" == "kind" ]]; then
        kind load image-archive -v 2 "${archive}"
    else
        minikube image load --logtostderr "${archive}"
    fi
    rm -f "${archive}"
}

image_pull() {
    local image="$1"
    if [[ "${CLUSTER_TYPE}" == "kind" ]]; then
        "${CONTAINER_RUNTIME}" pull "${IMAGE_NAMESPACE}/${image}"
        image_load "${IMAGE_NAMESPACE}/${image}"
    else
        minikube image pull --logtostderr "${IMAGE_NAMESPACE}/${image}"
    fi
}
