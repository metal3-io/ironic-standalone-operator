#!/usr/bin/env bash

set -eux -o pipefail

REPO_ROOT=$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")
cd "${REPO_ROOT}/test"

LOGDIR="${LOGDIR:-/tmp/logs}"
TEST_TIMEOUT="${TEST_TIMEOUT:-60m}"

. testing.env

mkdir -p "${LOGDIR}"

exec go test -timeout "${TEST_TIMEOUT}"
