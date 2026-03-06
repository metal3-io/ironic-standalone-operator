#!/usr/bin/env bash

set -eux -o pipefail

REPO_ROOT=$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")
cd "${REPO_ROOT}/test"

LOGDIR="${LOGDIR:-/tmp/logs}"
JUNIT_OUTPUT="${JUNIT_OUTPUT:-${LOGDIR}/report.xml}"
TEST_TIMEOUT="${TEST_TIMEOUT:-90m}"
IPA_SERVER_PORT=8089
IPA_DIR="/tmp/ipa"

# shellcheck disable=SC1091
. testing.env

mkdir -p "${LOGDIR}"

# Start local IPA server (downloaded by prepare.sh)
python3 -m http.server "${IPA_SERVER_PORT}" --directory "${IPA_DIR}" &>/dev/null &
IPA_SERVER_PID=$!
trap 'kill ${IPA_SERVER_PID} 2>/dev/null || true' EXIT

declare -a EXTRA_ARGS
if [[ -n "${LABEL_FILTER:-}" ]]; then
    EXTRA_ARGS=(--ginkgo.label-filter "${LABEL_FILTER}")
fi

go test --ginkgo.vv --ginkgo.junit-report "${JUNIT_OUTPUT}" -timeout "${TEST_TIMEOUT}" \
    --ginkgo.fail-on-empty "${EXTRA_ARGS[@]}"
