#!/usr/bin/env bash

set -eux

REPO_ROOT=$(realpath "$(dirname "${BASH_SOURCE[0]}")/../..")
cd "${REPO_ROOT}"

LOGDIR="${LOGDIR:-/tmp/logs}"
SCENARIO="${1:-}"
IRONIC="test/local-ironic/scenario-${SCENARIO}.yaml"
IRONIC_IP="${IRONIC_IP:-}"

if [[ -z "${SCENARIO}" ]] || [[ ! -f "${IRONIC}" ]] ; then
    echo "FATAL: valid scenario is required"
    exit 1
fi

if [[ -z "${IRONIC_IP}" ]] && grep -qE " +ipAddress:" "${IRONIC}"; then
    IRONIC_IP=$(grep -E " +ipAddress:" "${IRONIC}" | cut -d: -f2 | tr -d ' ')
fi

if [[ -n "${IRONIC_IP}" ]]; then
    sudo ip link delete ironic type dummy || true
    sudo ip link add ironic type dummy
    trap "sudo ip link delete ironic type dummy" EXIT

    sudo ip addr add "${IRONIC_IP}/24" brd + dev ironic label ironic:0
    sudo ip link set dev ironic up
else
    # Simple test: just use localhost
    IRONIC_IP=127.0.0.1
fi

mkdir -p "${LOGDIR}"
make build-run-local-ironic
sudo ./bin/run-local-ironic --input "${IRONIC}" \
    --output "${LOGDIR}/generated-${SCENARIO}.yaml" --verbose

# podman returns before containers fully start, so give it a graceful period
ATTEMPT=0
SUCCESS=
while [[ $ATTEMPT -ne 30 ]]
do
    sleep 2
    ATTEMPT=$(($ATTEMPT + 1))

    echo "Checking Ironic: attempt $ATTEMPT / 30"
    if curl -vfL "http://${IRONIC_IP}:6385/v1"; then
        SUCCESS=true
        break
    fi
done

if [[ "${SUCCESS}" != true ]]; then
    echo "FATAL: All ATTEMPTs failed!"
    exit 2
fi
