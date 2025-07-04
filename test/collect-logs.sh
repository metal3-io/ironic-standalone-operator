#!/usr/bin/env bash

# NOTE(dtantsur): do not use -e, commands can fail if the test breaks early
set -ux

REPO_ROOT=$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")
cd "${REPO_ROOT}"

LOGDIR="${LOGDIR:-/tmp/logs}"
CLUSTER_TYPE="${CLUSTER_TYPE:-kind}"

. test/testing.env

if [[ "${CLUSTER_TYPE}" == "minikube" ]]; then
    minikube logs --file "${LOGDIR}/minikube.log"
fi

kubectl get node -o wide > "${LOGDIR}/nodes.txt"
kubectl get deploy -o wide -A > "${LOGDIR}/deployments.txt"
kubectl get daemonset -o wide -A > "${LOGDIR}/daemonsets.txt"
kubectl get service -o wide -A > "${LOGDIR}/services.txt"

collect_from_ns() {
    local ns="$1"
    local dest="${LOGDIR}/${2:-"${ns}"}"
    mkdir -p "${dest}/pod"

    kubectl get all -n "${ns}" > "${dest}/all.txt"
    kubectl get pod -o yaml -n "${ns}" > "${dest}/pods.yaml"
    kubectl get deployment -o yaml -n "${ns}" > "${dest}/deployments.yaml"
    kubectl get daemonset -o yaml -n "${ns}" > "${dest}/daemonsets.yaml"

    for pod in $(kubectl get pod -o name -n "${ns}"); do
        kubectl describe "${pod}" -n "${ns}" >"${dest}/${pod}.txt"
        kubectl logs "${pod}" -n "${ns}" >"${dest}/${pod}.log" 2>&1
    done
}

collect_from_ns ironic-standalone-operator-system controller
collect_from_ns "${MARIADB_NAMESPACE}" ironic-database
collect_from_ns mariadb-operator
collect_from_ns kube-system
