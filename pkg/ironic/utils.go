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

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

const (
	probeInitialDelay     = 1
	probeTimeout          = 5
	probeFailureThreshold = 12

	serviceDNSSuffix = "svc"
)

type ControllerContext struct {
	Context     context.Context
	Client      client.Client
	KubeClient  kubernetes.Interface
	Scheme      *runtime.Scheme
	Logger      logr.Logger
	Domain      string
	VersionInfo VersionInfo
}

func mergeContainers(target, source []corev1.Container) []corev1.Container {
	if len(source) == 0 {
		return source
	}

	if len(target) != len(source) {
		// No way to reconcile different lists, re-create
		target = make([]corev1.Container, len(source))
	}

	for idx, src := range source {
		dest := &target[idx]
		dest.Name = src.Name
		dest.Image = src.Image
		dest.Command = src.Command
		dest.Ports = src.Ports
		dest.Env = src.Env
		dest.EnvFrom = src.EnvFrom
		dest.VolumeMounts = src.VolumeMounts
		dest.SecurityContext = src.SecurityContext
		if src.LivenessProbe != nil {
			dest.LivenessProbe = updateProbe(dest.LivenessProbe, src.LivenessProbe.ProbeHandler)
		} else {
			dest.LivenessProbe = nil
		}
		if src.ReadinessProbe != nil {
			dest.ReadinessProbe = updateProbe(dest.ReadinessProbe, src.ReadinessProbe.ProbeHandler)
		} else {
			dest.ReadinessProbe = nil
		}
		dest.Lifecycle = src.Lifecycle
	}

	return target
}

// mergePodTemplates updates an existing pod template, taking care to avoid
// overriding defaulted values.
func mergePodTemplates(target *corev1.PodTemplateSpec, source corev1.PodTemplateSpec) {
	if target.ObjectMeta.Labels == nil {
		target.ObjectMeta.Labels = make(map[string]string, len(source.ObjectMeta.Labels))
	}
	for k, v := range source.ObjectMeta.Labels {
		target.ObjectMeta.Labels[k] = v
	}

	target.Spec.InitContainers = mergeContainers(target.Spec.InitContainers, source.Spec.InitContainers)
	target.Spec.Containers = mergeContainers(target.Spec.Containers, source.Spec.Containers)
	target.Spec.Volumes = source.Spec.Volumes
	target.Spec.HostNetwork = source.Spec.HostNetwork
	if source.Spec.DNSPolicy != "" {
		target.Spec.DNSPolicy = source.Spec.DNSPolicy
	}
	if source.Spec.NodeSelector != nil {
		target.Spec.NodeSelector = source.Spec.NodeSelector
	}
}

func getDeploymentStatus(cctx ControllerContext, deploy *appsv1.Deployment) (Status, error) {
	if deploy.Status.ObservedGeneration != deploy.Generation {
		cctx.Logger.Info("deployment not ready yet", "Deployment", deploy.Name,
			"Generation", deploy.Generation, "ObservedGeneration", deploy.Status.ObservedGeneration)
		return inProgress("deployment not ready yet")
	}

	var available bool
	var err error
	for _, cond := range deploy.Status.Conditions {
		if cond.Type == appsv1.DeploymentAvailable && cond.Status == corev1.ConditionTrue {
			available = true
		}
		if cond.Type == appsv1.DeploymentReplicaFailure && cond.Status == corev1.ConditionTrue {
			err = fmt.Errorf("deployment failed: %s", cond.Message)
			// TODO(dtantsur): can we determine if it's fatal or not?
			return transientError(err)
		}
	}

	if available {
		return ready()
	} else {
		cctx.Logger.Info("deployment not available yet", "Deployment", deploy.Name,
			"Conditions", deploy.Status.Conditions)
		return inProgress("deployment not available yet")
	}
}

func getDaemonSetStatus(cctx ControllerContext, deploy *appsv1.DaemonSet) (Status, error) {
	if deploy.Status.ObservedGeneration != deploy.Generation {
		cctx.Logger.Info("daemon set not ready yet", "DaemonSet", deploy.Name,
			"Generation", deploy.Generation, "ObservedGeneration", deploy.Status.ObservedGeneration)
		return inProgress("daemon set not ready yet")
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
		return ready()
	} else {
		cctx.Logger.Info("daemon set not available yet", "DaemonSet", deploy.Name,
			"NumberUnavailable", deploy.Status.NumberUnavailable)
		return inProgress(fmt.Sprintf("daemon set not available yet: %d replicas unavailable", deploy.Status.NumberUnavailable))
	}
}

func getServiceStatus(service *corev1.Service) (Status, error) {
	// TODO(dtantsur): can we check anything else?
	if len(service.Spec.ClusterIPs) == 0 {
		return inProgress("service has no cluster IPs")
	}

	return ready()
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

func updateProbe(current *corev1.Probe, handler corev1.ProbeHandler) *corev1.Probe {
	if current == nil {
		current = &corev1.Probe{}
	}
	current.ProbeHandler = handler
	// NOTE(dtantsur): we want some delay because Ironic does not start instantly.
	// Also be conservative about failing the pod since Ironic restars are not cheap (the database is wiped).
	current.InitialDelaySeconds = probeInitialDelay
	current.TimeoutSeconds = probeTimeout
	current.FailureThreshold = probeFailureThreshold
	return current
}

func newProbe(handler corev1.ProbeHandler) *corev1.Probe {
	return updateProbe(nil, handler) // TODO: remove
}

func isReady(conditions []metav1.Condition) bool {
	return meta.IsStatusConditionTrue(conditions, string(metal3api.IronicStatusReady))
}

func appendStringEnv(envVars []corev1.EnvVar, name string, value string) []corev1.EnvVar {
	if value != "" {
		return append(envVars, corev1.EnvVar{
			Name:  name,
			Value: value,
		})
	}

	return envVars
}

func appendListOfStringsEnv(envVars []corev1.EnvVar, name string, value []string, sep string) []corev1.EnvVar {
	if len(value) > 0 {
		return append(envVars, corev1.EnvVar{
			Name:  name,
			Value: strings.Join(value, sep),
		})
	}

	return envVars
}
