package ironic

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	metal3api "github.com/metal3-io/ironic-standalone-operator/api/v1alpha1"
)

// Leave enough time for potential debugging but don't hold jobs forever.
const jobTTLSeconds int32 = 24 * 3600

type upgradePhase string

const (
	// Pre-upgrade: ensure the database schema is up-to-date.
	preUpgrade upgradePhase = "pre"
	// Post-upgrade: migrate any database resources to the new version.
	// The post-upgrade job runs parallel to Ironic.
	postUpgrade upgradePhase = "post"
)

var commandPerPhase = map[upgradePhase]string{
	preUpgrade:  "database-upgrade",
	postUpgrade: "online-data-migrations",
}

func upgradeJobRequired(cctx ControllerContext, ironic *metal3api.Ironic) bool {
	return ironic.Spec.Database != nil && cctx.VersionInfo.InstalledVersion.String() != ironic.Status.InstalledVersion
}

func newMigrationTemplate(cctx ControllerContext, ironic *metal3api.Ironic, phase upgradePhase) corev1.PodTemplateSpec {
	script := commandPerPhase[phase]
	database := ironic.Spec.Database

	volumes, mounts := databaseClientMounts(database)

	// NOTE(dtantsur): these are not actually needed for upgrade scripts but are currently required by configure-ironic
	volumes = append(volumes, corev1.Volume{
		Name: "ironic-shared",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	})
	mounts = append(mounts, corev1.VolumeMount{
		Name:      "ironic-shared",
		MountPath: sharedDir,
	})

	envVars := databaseClientEnvVars(database)
	envVars = append(envVars, []corev1.EnvVar{
		{
			Name:  "IRONIC_USE_MARIADB",
			Value: "true",
		},
		// NOTE(dtantsur): we don't really care about PROVISIONING_IP when running upgrade scripts,
		// but the Ironic configuration script requires it
		{
			Name: "PROVISIONING_IP",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					APIVersion: "v1",
					FieldPath:  "status.podIP",
				},
			},
		},
	}...)

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
				ReadOnlyRootFilesystem: ptr.To(true),
			},
		},
	}
	return applyOverridesToPod(ironic.Spec.Overrides, addDataVolumes(corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				metal3api.IronicServiceLabel: ironic.Name,
				metal3api.IronicVersionLabel: cctx.VersionInfo.InstalledVersion.String(),
			},
		},
		Spec: corev1.PodSpec{
			Containers: containers,
			Volumes:    volumes,
			// https://kubernetes.io/docs/concepts/workloads/controllers/job/#pod-backoff-failure-policy
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}))
}

func ensureIronicUpgradeJob(cctx ControllerContext, resources Resources, phase upgradePhase) (Status, error) {
	if !upgradeJobRequired(cctx, resources.Ironic) {
		return ready()
	}

	fromVersion := resources.Ironic.Status.InstalledVersion
	if fromVersion == "" {
		fromVersion = "none"
	}
	toVersion := cctx.VersionInfo.InstalledVersion

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s-%s-to-%s", resources.Ironic.Name, phase, fromVersion, cctx.VersionInfo.InstalledVersion),
			Namespace: resources.Ironic.Namespace,
		},
	}

	template := newMigrationTemplate(cctx, resources.Ironic, phase)

	result, err := controllerutil.CreateOrUpdate(cctx.Context, cctx.Client, job, func() error {
		if job.ObjectMeta.Labels == nil {
			cctx.Logger.Info("creating a new upgrade job", "Phase", phase, "From", fromVersion, "To", toVersion.String())
			job.ObjectMeta.Labels = make(map[string]string, 2)
		}
		job.ObjectMeta.Labels[metal3api.IronicServiceLabel] = resources.Ironic.Name
		job.ObjectMeta.Labels[metal3api.IronicVersionLabel] = cctx.VersionInfo.InstalledVersion.String()

		job.Spec.TTLSecondsAfterFinished = ptr.To(jobTTLSeconds)
		mergePodTemplates(&job.Spec.Template, template)
		job.Spec.PodReplacementPolicy = ptr.To(batchv1.Failed)

		return controllerutil.SetControllerReference(resources.Ironic, job, cctx.Scheme)
	})
	if result != controllerutil.OperationResultNone {
		cctx.Logger.Info("ironic upgrade job", "Job", job.Name, "Status", result,
			"Phase", phase, "From", fromVersion, "To", toVersion.String())
		return updated()
	}
	if err != nil {
		return transientError(err)
	}

	status, err := getJobStatus(cctx, job, fmt.Sprintf("%s-upgrade", phase))
	if status.IsReady() && err == nil {
		cctx.Logger.Info("upgrade job succeeded", "Phase", phase, "From", fromVersion, "To", toVersion.String())
	}

	return status, err
}
