package ironic

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metal3api "github.com/metal3-io/ironic-operator/api/v1alpha1"
)

const (
	probeInitialDelay     = 1
	probeTimeout          = 5
	probeFailureThreshold = 12

	serviceDNSSuffix = "svc"
)

type ControllerContext struct {
	Context    context.Context
	Client     client.Client
	KubeClient kubernetes.Interface
	Scheme     *runtime.Scheme
	Logger     logr.Logger
}

func getDeploymentStatus(deploy *appsv1.Deployment) (metal3api.IronicStatusConditionType, error) {
	if deploy.Status.ObservedGeneration != deploy.Generation {
		return metal3api.IronicStatusProgressing, nil
	}

	var available bool
	var err error
	for _, cond := range deploy.Status.Conditions {
		if cond.Type == appsv1.DeploymentAvailable && cond.Status == corev1.ConditionTrue {
			available = true
		}
		if cond.Type == appsv1.DeploymentReplicaFailure && cond.Status == corev1.ConditionTrue {
			err = fmt.Errorf("deployment failed: %s", cond.Message)
			return metal3api.IronicStatusProgressing, err
		}
	}

	if available {
		return metal3api.IronicStatusAvailable, nil
	} else {
		return metal3api.IronicStatusProgressing, nil
	}
}

func getDaemonSetStatus(deploy *appsv1.DaemonSet) (metal3api.IronicStatusConditionType, error) {
	if deploy.Status.ObservedGeneration != deploy.Generation {
		return metal3api.IronicStatusProgressing, nil
	}

	var available bool

	// FIXME(dtantsur): the current version of appsv1 does not seem to have
	// constants for conditions types.
	// var err error
	// for _, cond := range deploy.Status.Conditions {
	// 	if cond.Type == appsv1.??? && cond.Status == corev1.ConditionTrue {
	// 		available = true
	// 	}
	// 	if cond.Type == appsv1.??? && cond.Status == corev1.ConditionTrue {
	// 		err = fmt.Errorf("deployment failed: %s", cond.Message)
	// 		return metal3api.IronicStatusProgressing, err
	// 	}
	// }
	available = deploy.Status.NumberUnavailable == 0

	if available {
		return metal3api.IronicStatusAvailable, nil
	} else {
		return metal3api.IronicStatusProgressing, nil
	}
}

func buildEndpoints(ips []string, port int, includeProto string) (endpoints []string) {
	portString := fmt.Sprint(port)
	for _, ip := range ips {
		var endpoint string
		if (includeProto == "https" && port == 443) || (includeProto == "http" && port == 80) {
			if strings.Contains(ip, ":") {
				endpoint = fmt.Sprintf("%s://[%s]", includeProto, ip) // IPv6
			} else {
				endpoint = fmt.Sprintf("%s://%s", includeProto, ip)
			}
		} else {
			endpoint = net.JoinHostPort(ip, portString)
			if includeProto != "" {
				endpoint = fmt.Sprintf("%s://%s", includeProto, endpoint)
			}
		}

		endpoints = append(endpoints, endpoint)
	}
	sort.Strings(endpoints)
	return
}

func newProbe(handler corev1.ProbeHandler) *corev1.Probe {
	return &corev1.Probe{
		ProbeHandler: handler,
		// NOTE(dtantsur): we want some delay because Ironic does not start instantly.
		// Also be conservative about failing the pod since Ironic restars are not cheap (the database is wiped).
		InitialDelaySeconds: probeInitialDelay,
		TimeoutSeconds:      probeTimeout,
		FailureThreshold:    probeFailureThreshold,
	}
}

func isReady(conditions []metav1.Condition) bool {
	return meta.IsStatusConditionTrue(conditions, string(metal3api.IronicStatusAvailable))
}
