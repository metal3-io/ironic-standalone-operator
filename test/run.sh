#!/usr/bin/env bash

set -eux -o pipefail

REPO_ROOT=$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")
cd "${REPO_ROOT}/test"

LOGDIR="${LOGDIR:-/tmp/logs}"
JUNIT_OUTPUT="${JUNIT_OUTPUT:-${LOGDIR}/report.xml}"
TEST_TIMEOUT="${TEST_TIMEOUT:-90m}"

. testing.env

mkdir -p "${LOGDIR}"

declare -a EXTRA_ARGS
if [[ -n "${LABEL_FILTER:-}" ]]; then
    EXTRA_ARGS=(--ginkgo.label-filter "${LABEL_FILTER}")
fi

exec go test --ginkgo.vv --ginkgo.junit-report "${JUNIT_OUTPUT}" -timeout "${TEST_TIMEOUT}" \
    --ginkgo.fail-on-empty "${EXTRA_ARGS[@]}"
