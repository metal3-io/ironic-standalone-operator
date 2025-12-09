#!/usr/bin/env bash

# NOTE(dtantsur): do not use -e, commands can fail if the test breaks early
set -ux

LOGDIR="${LOGDIR:-/tmp/logs}"

sudo podman ps --all | tee "${LOGDIR}/containers.txt" > /dev/null
for cid in $(sudo podman ps --quiet); do
    sudo podman inspect "${cid}" | tee "${LOGDIR}/${cid}.txt" > /dev/null
    sudo podman logs "${cid}" 2>&1 | tee "${LOGDIR}/${cid}.log" > /dev/null
done
sudo chown -R "${USER}" "${LOGDIR}"
