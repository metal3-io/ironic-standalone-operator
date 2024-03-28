#!/bin/bash

set -ux

LOGDIR="${LOGDIR:-/tmp/logs}"

mkdir -p "${LOGDIR}"

kubectl logs \
    -n ironic-standalone-operator-system deployment/ironic-standalone-operator-controller-manager \
    > "${LOGDIR}/controller.log"
