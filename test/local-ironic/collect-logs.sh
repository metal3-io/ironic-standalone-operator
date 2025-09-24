#!/usr/bin/env bash

# NOTE(dtantsur): do not use -e, commands can fail if the test breaks early
set -ux

LOGDIR="${LOGDIR:-/tmp/logs}"

podman ps --all > "${LOGDIR}/containers.txt"
for cid in $(podman ps --quiet); do
    podman inspect "${cid}" > "${LOGDIR}/${cid}.txt"
    podman logs "${cid}" > "${LOGDIR}/${cid}.log" 2>&1
done
