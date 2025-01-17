#!/bin/bash

set -eux -o pipefail

IMG="${IMG:-localhost/controller:test}"
LOGDIR="${LOGDIR:-/tmp/logs}"
CONTAINER_RUNTIME="${CONTAINER_RUNTIME:-}"
CERT_MANAGER_VERSION="${CERT_MANAGER_VERSION:-1.16.1}"

. "$(dirname "$0")/testing.env"

mkdir -p "${LOGDIR}"

if [[ -z "${CONTAINER_RUNTIME}" ]]; then
    if command -v podman &> /dev/null;  then
        CONTAINER_RUNTIME=podman
    else
        CONTAINER_RUNTIME=docker
    fi
fi

# Installing cert-manager

kubectl apply -f "https://github.com/cert-manager/cert-manager/releases/download/v${CERT_MANAGER_VERSION}/cert-manager.yaml"
kubectl wait --for=condition=Available --timeout=60s \
    -n cert-manager deployment/cert-manager
kubectl wait --for=condition=Available --timeout=60s \
    -n cert-manager deployment/cert-manager-webhook

# Caching required images

kind_load() {
    local image="$1"
    if [[ "${CONTAINER_RUNTIME}" == "podman" ]]; then
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

# Building and installing the operator

"${CONTAINER_RUNTIME}" build -t "${IMG}" . 2>&1 | tee "${LOGDIR}/docker-build.log"
kind_load "${IMG}"
make install deploy IMG="${IMG}" DEPLOY_TARGET=testing

kubectl wait --for=condition=Available --timeout=60s \
    -n ironic-standalone-operator-system deployment/ironic-standalone-operator-controller-manager

# Preparing the TLS certificate

openssl req -x509 -new -subj "/CN=ironic" \
    -addext "subjectAltName = IP:${IRONIC_IP},IP:${PROVISIONING_IP}" \
    -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -nodes \
    -keyout "${IRONIC_KEY_FILE}" -out "${IRONIC_CERT_FILE}"
