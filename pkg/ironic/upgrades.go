package ironic

import (
	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func upgradeJobRequired(cctx ControllerContext, ironic *metal3api.Ironic, db *metal3api.Database) bool {
	return db != nil && cctx.VersionInfo.InstalledVersion.String() != ironic.Status.InstalledVersion
}

func newMigrationTemplate(cctx ControllerContext, ironic *metal3api.Ironic, database *metal3api.Database, script string) corev1.PodTemplateSpec {
	volumes, mounts := databaseClientMounts(database)

	envVars := databaseClientEnvVars(database)
	// NOTE(dtantsur): we don't really care about PROVISIONING_IP when running upgrade scripts,
	// but the Ironic configuration script requires it
	envVars = append(envVars,
		corev1.EnvVar{
			Name: "PROVISIONING_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "status.podIP",
				},
			},
		},
	)

	containers := []corev1.Container{
		{
			Name:         script,
			Image:        cctx.VersionInfo.IronicImage,
			Command:      []string{"/bin/run" + script},
			Env:          envVars,
			VolumeMounts: mounts,
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:  ptr.To(ironicUser),
				RunAsGroup: ptr.To(ironicGroup),
				Capabilities: &corev1.Capabilities{
					Drop: []corev1.Capability{"ALL"},
				},
			},
		},
	}
	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				metal3api.IronicAppLabel:     ironicDeploymentName(ironic),
				metal3api.IronicVersionLabel: cctx.VersionInfo.InstalledVersion.String(),
			},
		},
		Spec: corev1.PodSpec{
			Containers: containers,
			Volumes:    volumes,
			// https://kubernetes.io/docs/concepts/workloads/controllers/job/#pod-backoff-failure-policy
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
}

func newPreUpgradePodTemplate(cctx ControllerContext, ironic *metal3api.Ironic, database *metal3api.Database) (corev1.PodTemplateSpec, error) {
	return newMigrationTemplate(cctx, ironic, database, "database-upgrade"), nil
}

func newPostUpgradePodTemplate(cctx ControllerContext, ironic *metal3api.Ironic, database *metal3api.Database) (corev1.PodTemplateSpec, error) {
	return newMigrationTemplate(cctx, ironic, database, "online-data-migrations"), nil
}
