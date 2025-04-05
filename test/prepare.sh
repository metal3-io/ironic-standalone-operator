#!/usr/bin/env bash

set -eux -o pipefail

REPO_ROOT=$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")
cd "${REPO_ROOT}"

IMG="${IMG:-localhost/controller:test}"
LOGDIR="${LOGDIR:-/tmp/logs}"
CONTAINER_RUNTIME="${CONTAINER_RUNTIME:-}"
CERT_MANAGER_VERSION="${CERT_MANAGER_VERSION:-1.16.1}"
CLUSTER_TYPE="${CLUSTER_TYPE:-kind}"

. test/testing.env

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

image_load() {
    local image="$1"
    local archive="$(mktemp --suffix=.tar)"
    "${CONTAINER_RUNTIME}" save "${image}" > "${archive}"
    if [[ "${CLUSTER_TYPE}" == "kind" ]]; then
        kind load image-archive -v 2 "${archive}"
    else
        minikube image load --logtostderr "${archive}"
    fi
    rm -f "${archive}"
}

for image in ironic mariadb ironic-ipa-downloader; do
    if [[ "${CLUSTER_TYPE}" == "kind" ]]; then
        "${CONTAINER_RUNTIME}" pull "quay.io/metal3-io/${image}"
        image_load "quay.io/metal3-io/${image}"
    else
        minikube image pull --logtostderr "quay.io/metal3-io/${image}"
    fi
done

# Building and installing the operator

"${CONTAINER_RUNTIME}" build -t "${IMG}" . 2>&1 | tee "${LOGDIR}/docker-build.log"
image_load "${IMG}" 2>&1

make install deploy IMG="${IMG}" DEPLOY_TARGET=testing

kubectl wait --for=condition=Available --timeout=60s \
    -n ironic-standalone-operator-system deployment/ironic-standalone-operator-controller-manager

# Preparing the TLS certificate

openssl req -x509 -new -subj "/CN=ironic" \
    -addext "subjectAltName = IP:${IRONIC_IP},IP:${PROVISIONING_IP}" \
    -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -nodes \
    -keyout "${IRONIC_KEY_FILE}" -out "${IRONIC_CERT_FILE}"
