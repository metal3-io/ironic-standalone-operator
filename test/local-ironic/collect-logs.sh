#!/usr/bin/env bash

# NOTE(dtantsur): do not use -e, commands can fail if the test breaks early
set -ux

LOGDIR="${LOGDIR:-/tmp/logs}"

sudo podman ps --all > "${LOGDIR}/containers.txt"
for cid in $(sudo podman ps --quiet); do
    sudo podman inspect "${cid}" > "${LOGDIR}/${cid}.txt"
    sudo podman logs "${cid}" > "${LOGDIR}/${cid}.log" 2>&1
done
sudo chown -R "${USER}" "${LOGDIR}"
