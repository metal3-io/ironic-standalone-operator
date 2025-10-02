#!/usr/bin/env bash

set -eux

REPO_ROOT=$(realpath "$(dirname "${BASH_SOURCE[0]}")/../..")
cd "${REPO_ROOT}"

LOGDIR="${LOGDIR:-/tmp/logs}"
SCENARIO="${1:-}"
IRONIC="test/local-ironic/scenario-${SCENARIO}.yaml"

if [[ -z "${SCENARIO}" ]] || [[ ! -f "${IRONIC}" ]] ; then
    echo "FATAL: valid scenario is required"
    exit 1
fi

mkdir -p "${LOGDIR}"
make build-run-local-ironic
sudo ./bin/run-local-ironic --down --input "${IRONIC}" \
    --output "${LOGDIR}/generated-${SCENARIO}.yaml" --verbose
