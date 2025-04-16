#!/usr/bin/env bash

# NOTE(dtantsur): do not use -e, commands can fail if the test breaks early
set -ux

REPO_ROOT=$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")
cd "${REPO_ROOT}"

LOGDIR="${LOGDIR:-/tmp/logs}"
CLUSTER_TYPE="${CLUSTER_TYPE:-kind}"

mkdir -p "${LOGDIR}/controller/" "${LOGDIR}/mariadb/pod"

. test/testing.env

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

kubectl get all -n "${MARIADB_NAMESPACE}" > "${LOGDIR}/mariadb/all.txt"
kubectl get pod -o yaml -n "${MARIADB_NAMESPACE}" > "${LOGDIR}/mariadb/pods.yaml"
for pod in $(kubectl get pod -o name -n "${MARIADB_NAMESPACE}"); do
    kubectl logs "${pod}" -n "${MARIADB_NAMESPACE}" >"${LOGDIR}/mariadb/${pod}.log" 2>&1
done
