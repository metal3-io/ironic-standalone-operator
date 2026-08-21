#!/usr/bin/env bash

# NOTE(dtantsur): do not use -e, commands can fail if the test breaks early
set -ux

LOGDIR="${LOGDIR:-/tmp/logs}"

# NOTE(dtantsur): podman tends to hang occasionally, work around
run_podman() {
    local attempts=3
    while true; do
        sudo timeout 15s podman "$@"
        local result="$?"
        # exit code 124 is timeout
        if [[ "${result}" == "124" ]] && (( attempts > 0 )); then
            attempts=$((attempts-1))
        else
            exit "${result}"
        fi
    done
}

run_podman ps --all | tee "${LOGDIR}/containers.txt"
for cid in $(run_podman ps --quiet); do
    run_podman inspect "${cid}" | tee "${LOGDIR}/${cid}.txt" > /dev/null
    run_podman logs "${cid}" 2>&1 | tee "${LOGDIR}/${cid}.log" > /dev/null
done
sudo chown -R "${USER}" "${LOGDIR}"
