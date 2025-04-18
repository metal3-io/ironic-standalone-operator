#!/usr/bin/env bash

# NOTE(dtantsur): do not use -e, commands can fail if the test breaks early
set -ux

LOGDIR="${LOGDIR:-/tmp/logs}"
CLUSTER_TYPE="${CLUSTER_TYPE:-kind}"

mkdir -p "${LOGDIR}/controller/"

kubectl get -o yaml \
    -n ironic-standalone-operator-system deployment/ironic-standalone-operator-controller-manager \
    > "${LOGDIR}/controller/deployment.yaml"
kubectl get pod -o yaml \
    -n ironic-standalone-operator-system > "${LOGDIR}/controller/pods.yaml"
kubectl logs \
    -n ironic-standalone-operator-system deployment/ironic-standalone-operator-controller-manager \
    > "${LOGDIR}/controller/manager.log"
if [[ "${CLUSTER_TYPE}" == "minikube" ]]; then
    minikube logs --file "${LOGDIR}/minikube.log"
fi
