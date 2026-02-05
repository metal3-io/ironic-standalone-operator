#!/usr/bin/env bash

set -eux -o pipefail

REPO_ROOT=$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")
cd "${REPO_ROOT}"
mkdir -p bin

IMG="${IMG:-localhost/controller:test}"
LOGDIR="${LOGDIR:-/tmp/logs}"
CONTAINER_RUNTIME="${CONTAINER_RUNTIME:-}"

HELM_VERSION=3.18.6
HELM_CHECKSUM=3f43c0aa57243852dd542493a0f54f1396c0bc8ec7296bbb2c01e802010819ce
HELM_FILE="helm-v${HELM_VERSION}-linux-amd64.tar.gz"

CERT_MANAGER_VERSION="${CERT_MANAGER_VERSION:-1.18.2}"
MARIADB_OPERATOR_VERSION="${MARIADB_OPERATOR_VERSION:-25.8.3}"

# shellcheck disable=SC1091
. test/testing.env
. test/utils.sh

mkdir -p "${LOGDIR}"

pushd /tmp
curl -OL "https://get.helm.sh/${HELM_FILE}"
echo "${HELM_CHECKSUM} ${HELM_FILE}" | sha256sum -c
tar -xzf "${HELM_FILE}"
HELM="${REPO_ROOT}/bin/helm"
mv linux-amd64/helm "${HELM}"
popd

# Installing cert-manager

"${HELM}" repo add jetstack https://charts.jetstack.io --force-update
"${HELM}" install cert-manager jetstack/cert-manager \
  --namespace cert-manager --create-namespace \
  --version "v${CERT_MANAGER_VERSION}" --set crds.enabled=true

# Installing MariaDB operator

# Add MariaDB Operator repository using direct GitHub Pages URL
"${HELM}" repo add mariadb-operator https://mariadb-operator.github.io/mariadb-operator
"${HELM}" repo update mariadb-operator

# Install MariaDB Operator CRDs with retries
echo "Installing MariaDB Operator CRDs..."
for i in {1..5}; do
    if "${HELM}" install mariadb-operator-crds mariadb-operator/mariadb-operator-crds \
        --version "${MARIADB_OPERATOR_VERSION}"; then
        break
    fi
    echo "Attempt $i failed, retrying in 10s..."
    sleep 10

    if [[ "$i" -eq 5 ]]; then
        echo "ERROR: Failed to install MariaDB CRDs after 5 attempts."
        exit 1
    fi
done

# Install MariaDB Operator with retries
echo "Installing MariaDB Operator..."
for i in {1..5}; do
    if "${HELM}" install mariadb-operator mariadb-operator/mariadb-operator \
        --namespace mariadb-operator --create-namespace \
        --version "${MARIADB_OPERATOR_VERSION}" \
        --set webhook.cert.certManager.enabled=true \
        --set certController.enabled=false; then
        break
    fi
    echo "Attempt $i failed, retrying in 10s..."
    sleep 10

    if [[ "$i" -eq 5 ]]; then
        echo "ERROR: Failed to install MariaDB Operator after 5 attempts."
        exit 1
    fi
done

# Caching required images

for image in ironic ironic-ipa-downloader keepalived; do
    image_pull "${image}"
done

# Building and installing the operator

"${CONTAINER_RUNTIME}" build -t "${IMG}" . 2>&1 | tee "${LOGDIR}/docker-build.log"
image_load "${IMG}" 2>&1

make install deploy IMG="${IMG}" DEPLOY_TARGET=testing

kubectl wait --for=condition=Available --timeout=60s \
    -n ironic-standalone-operator-system deployment/ironic-standalone-operator-controller-manager

# Preparing the TLS certificate

SUBJECT_ALT_NAME="IP:${PROVISIONING_IP}"
if [[ "${CLUSTER_TYPE}" == kind ]]; then
    SUBJECT_ALT_NAME+=",IP:${IRONIC_IP}"
else
    for node_ip in $(minikube node list | awk '{ print $2; }'); do
        SUBJECT_ALT_NAME+=",IP:${node_ip}"
    done
fi
openssl req -x509 -new -subj "/CN=ironic" \
    -addext "subjectAltName = ${SUBJECT_ALT_NAME}" \
    -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 -nodes \
    -keyout "${IRONIC_KEY_FILE}" -out "${IRONIC_CERT_FILE}"

# Creating the database

kubectl create -f test/database.yaml
