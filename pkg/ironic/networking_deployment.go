package ironic

import (
	"fmt"
	"maps"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

const (
	switchConfigKey = "switch-configs.conf"

	// managedSecretLabel is set on secrets created by the operator to distinguish
	// them from externally-provided secrets. EnsureSwitchConfigSecretDeleted only
	// deletes secrets that carry this label.
	managedSecretLabel = "ironic.metal3.io/managed"
)

// NetworkingDeploymentName returns the name of the networking service deployment.
func NetworkingDeploymentName(ironic *metal3api.Ironic) string {
	return ironic.Name + "-networking"
}

// NetworkingServiceName returns the name of the networking service.
func NetworkingServiceName(ironic *metal3api.Ironic) string {
	return ironic.Name + "-networking-service"
}

// networkingLabels returns common labels for networking service resources.
func networkingLabels(ironic *metal3api.Ironic) map[string]string {
	return map[string]string{
		"app":                       "ironic-networking",
		"ironic.metal3.io/instance": ironic.Name,
	}
}

// buildNetworkingContainerEnv builds the environment variables for the networking container.
func buildNetworkingContainerEnv(resources Resources) []corev1.EnvVar {
	envVars := []corev1.EnvVar{
		{
			Name:  "IRONIC_NETWORKING_ENABLED",
			Value: "true",
		},
		{
			Name:  "IRONIC_NETWORKING_JSON_RPC_HOST",
			Value: "0.0.0.0",
		},
		{
			Name:  "IRONIC_NETWORKING_JSON_RPC_PORT",
			Value: strconv.Itoa(int(resources.Ironic.Spec.NetworkingService.RPCPort)),
		},
		{
			Name:  "IRONIC_NETWORKING_SWITCH_CONFIGS",
			Value: "/etc/ironic/networking/configs/switch-configs.conf",
		},
	}

	envVars = appendListOfStringsEnv(envVars, "IRONIC_NETWORKING_ENABLED_SWITCH_DRIVERS",
		resources.Ironic.Spec.NetworkingService.SwitchDrivers, ",")

	// Set the password for RPC
	if resources.APISecret != nil {
		envVars = append(envVars,
			corev1.EnvVar{
				Name: "IRONIC_HTPASSWD",
				ValueFrom: &corev1.EnvVarSource{
					SecretKeyRef: &corev1.SecretKeySelector{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: resources.APISecret.Name,
						},
						Key: htpasswdKey,
					},
				},
			},
		)
	}

	// Get Ironic IP for service catalog
	if resources.Ironic.Spec.Networking.IPAddress != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name:  "IRONIC_IP",
			Value: resources.Ironic.Spec.Networking.IPAddress,
		})
	}

	// Add TLS settings if enabled
	if resources.TLSSecret != nil {
		if insecureRPC := resources.Ironic.Spec.TLS.InsecureRPC; insecureRPC != nil && *insecureRPC {
			envVars = append(envVars, corev1.EnvVar{
				Name:  "IRONIC_INSECURE",
				Value: "true",
			})
		}
	}

	return envVars
}

// buildNetworkingContainer builds the networking service container.
func buildNetworkingContainer(cctx ControllerContext, resources Resources, mounts []corev1.VolumeMount) corev1.Container {
	image := resources.Ironic.Spec.Images.Ironic
	if image == "" {
		image = cctx.VersionInfo.IronicImage
	}

	return corev1.Container{
		Name:            "ironic-networking",
		Image:           image,
		ImagePullPolicy: corev1.PullAlways,
		Command:         []string{"/bin/runironic-networking"},
		Env:             buildNetworkingContainerEnv(resources),
		Ports: []corev1.ContainerPort{
			{
				Name:          ironicNetworkingPortName,
				ContainerPort: resources.Ironic.Spec.NetworkingService.RPCPort,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		VolumeMounts: mounts,
		LivenessProbe: newProbe(corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt32(resources.Ironic.Spec.NetworkingService.RPCPort),
			},
		}),
		ReadinessProbe: newProbe(corev1.ProbeHandler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt32(resources.Ironic.Spec.NetworkingService.RPCPort),
			},
		}),
		SecurityContext: &corev1.SecurityContext{
			RunAsUser:  ptr.To(ironicUser),
			RunAsGroup: ptr.To(ironicGroup),
			Capabilities: &corev1.Capabilities{
				Drop: []corev1.Capability{"ALL"},
			},
		},
	}
}

func buildNetworkingVolumesAndMounts(resources Resources) (volumes []corev1.Volume, mounts []corev1.VolumeMount) {
	// RPC auth credentials (the ironic-image startup script reads credentials
	// from /auth/ironic-rpc for JSON-RPC authentication)
	if resources.APISecret != nil {
		volumes = append(volumes, corev1.Volume{
			Name: "ironic-auth",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  resources.APISecret.Name,
					DefaultMode: ptr.To(corev1.SecretVolumeSourceDefaultMode),
				},
			},
		})
		mounts = append(mounts, corev1.VolumeMount{
			Name:      "ironic-auth",
			MountPath: authDir + "/ironic-rpc",
		})
	}

	// Switch config (always required for networking service)
	volumes = append(volumes, corev1.Volume{
		Name: "switch-config",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  SwitchConfigSecretName(resources.Ironic),
				DefaultMode: ptr.To(corev1.SecretVolumeSourceDefaultMode),
			},
		},
	})
	mounts = append(mounts, corev1.VolumeMount{
		Name:      "switch-config",
		MountPath: "/etc/ironic/networking/configs",
		ReadOnly:  true,
	})

	// Switch credentials (optional)
	if resources.Ironic.Spec.NetworkingService != nil &&
		resources.Ironic.Spec.NetworkingService.SwitchCredentialsSecretName != "" {
		volumes = append(volumes, corev1.Volume{
			Name: "switch-credentials",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  resources.Ironic.Spec.NetworkingService.SwitchCredentialsSecretName,
					DefaultMode: ptr.To(corev1.SecretVolumeSourceDefaultMode),
				},
			},
		})
		mounts = append(mounts, corev1.VolumeMount{
			Name:      "switch-credentials",
			MountPath: "/etc/ironic/networking/credentials",
			ReadOnly:  true,
		})
	}

	// TLS certificate (only /certs/ironic, not vmedia — networking service
	// doesn't serve virtual media)
	if resources.TLSSecret != nil {
		volumes = append(volumes, corev1.Volume{
			Name: "cert-ironic",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  resources.TLSSecret.Name,
					DefaultMode: ptr.To(corev1.SecretVolumeSourceDefaultMode),
				},
			},
		})
		mounts = append(mounts, corev1.VolumeMount{
			Name:      "cert-ironic",
			MountPath: certsDir + "/ironic",
			ReadOnly:  true,
		})
	}

	// Trusted CA bundle
	if resources.TrustedCAConfigMap != nil {
		volumes = append(volumes, corev1.Volume{
			Name: "trusted-ca",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: resources.TrustedCAConfigMap.Name,
					},
					DefaultMode: ptr.To(corev1.ConfigMapVolumeSourceDefaultMode),
				},
			},
		})
		mounts = append(mounts, corev1.VolumeMount{
			Name:      "trusted-ca",
			MountPath: certsDir + "/ca/trusted",
			ReadOnly:  true,
		})
	}

	return volumes, mounts
}

// BuildNetworkingDeployment builds the Deployment for the networking service.
func BuildNetworkingDeployment(cctx ControllerContext, resources Resources) *appsv1.Deployment {
	labels := networkingLabels(resources.Ironic)

	volumes, mounts := buildNetworkingVolumesAndMounts(resources)

	// Pod annotations include secret versions for auto-restart
	annotations := make(map[string]string)
	if resources.SwitchConfigSecret != nil {
		maps.Copy(annotations, secretVersionAnnotations("switch-config", resources.SwitchConfigSecret))
	}
	if resources.SwitchCredentialsSecret != nil {
		maps.Copy(annotations, secretVersionAnnotations("switch-credentials", resources.SwitchCredentialsSecret))
	}
	// Track TLS secret so networking pod restarts when TLS cert is updated
	if resources.TLSSecret != nil {
		maps.Copy(annotations, secretVersionAnnotations("tls-secret", resources.TLSSecret))
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NetworkingDeploymentName(resources.Ironic),
			Namespace: resources.Ironic.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptr.To[int32](1),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						buildNetworkingContainer(cctx, resources, mounts),
					},
					Volumes:      volumes,
					NodeSelector: resources.Ironic.Spec.NodeSelector,
				},
			},
		},
	}

	return deployment
}

// BuildNetworkingService builds the Service for the networking service.
func BuildNetworkingService(ironic *metal3api.Ironic) *corev1.Service {
	labels := networkingLabels(ironic)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      NetworkingServiceName(ironic),
			Namespace: ironic.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type:     corev1.ServiceTypeClusterIP,
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Name:     ironicNetworkingPortName,
					Protocol: corev1.ProtocolTCP,
					Port:     ironic.Spec.NetworkingService.RPCPort,
				},
			},
		},
	}

	return service
}

// NetworkingServiceEndpoint returns the hostname for the operator-managed networking service.
// Returns only the Service DNS name without port. The port is determined by the RPCPort
// field in the Ironic spec.
func NetworkingServiceEndpoint(ironic *metal3api.Ironic) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local", NetworkingServiceName(ironic), ironic.Namespace)
}

// GetNetworkingServiceEndpoint returns the networking service endpoint to use.
// Returns empty string if networking service is disabled.
// Returns the user-provided endpoint if specified, otherwise the operator-managed Service DNS name.
func GetNetworkingServiceEndpoint(ironic *metal3api.Ironic) string {
	if ironic.Spec.NetworkingService == nil || !ironic.Spec.NetworkingService.Enabled {
		return ""
	}

	if ironic.Spec.NetworkingService.Endpoint != "" {
		return ironic.Spec.NetworkingService.Endpoint
	}

	return NetworkingServiceEndpoint(ironic)
}

// EnsureSwitchConfigSecret ensures the switch config secret exists.
// If the secret already exists, it is left unchanged (data may be managed externally
// or by the IronicSwitch controller). If the secret does not exist, it is created
// with an empty switch-configs.conf so the networking pod can start cleanly.
// Owner references are handled by the controller's getAndUpdateSecret call that follows.
func EnsureSwitchConfigSecret(cctx ControllerContext, ironic *metal3api.Ironic) error {
	secretName := SwitchConfigSecretName(ironic)

	// Check if the secret already exists
	existing := &corev1.Secret{}
	err := cctx.Client.Get(cctx.Context, types.NamespacedName{
		Name:      secretName,
		Namespace: ironic.Namespace,
	}, existing)

	if err == nil {
		// Secret already exists, leave it alone
		return nil
	}

	if !k8serrors.IsNotFound(err) {
		return fmt.Errorf("failed to get switch config secret: %w", err)
	}

	// Create a new secret with empty config
	cctx.Logger.Info("creating switch config secret", "secret", secretName)
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: ironic.Namespace,
			Labels: map[string]string{
				managedSecretLabel: "true",
			},
		},
		Data: map[string][]byte{
			switchConfigKey: {},
		},
	}

	// Set owner reference
	if err := controllerutil.SetControllerReference(ironic, secret, cctx.Scheme); err != nil {
		return fmt.Errorf("failed to set owner reference on switch config secret: %w", err)
	}

	if err := cctx.Client.Create(cctx.Context, secret); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			// Race condition: another reconcile created it first
			return nil
		}
		return fmt.Errorf("failed to create switch config secret: %w", err)
	}

	return nil
}

// EnsureSwitchConfigSecretDeleted deletes the switch config secret only if it
// was created by the operator (identified by the managedSecretLabel). Externally-
// provided secrets are left untouched.
func EnsureSwitchConfigSecretDeleted(cctx ControllerContext, ironic *metal3api.Ironic) error {
	secretName := SwitchConfigSecretName(ironic)
	secret := &corev1.Secret{}
	if err := cctx.Client.Get(cctx.Context, types.NamespacedName{
		Name:      secretName,
		Namespace: ironic.Namespace,
	}, secret); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("failed to get switch config secret: %w", err)
	}

	// Only delete secrets that the operator created
	if secret.Labels[managedSecretLabel] != "true" {
		cctx.Logger.Info("skipping deletion of externally-managed switch config secret", "secret", secretName)
		return nil
	}

	if err := cctx.Client.Delete(cctx.Context, secret); err != nil {
		if k8serrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("failed to delete switch config secret: %w", err)
	}

	cctx.Logger.Info("deleted switch config secret", "secret", secretName)
	return nil
}
