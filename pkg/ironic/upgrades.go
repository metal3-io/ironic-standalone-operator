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

// Leave enough time for potential debugging but don't hold jobs forever
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

func upgradeJobRequired(cctx ControllerContext, ironic *metal3api.Ironic, db *metal3api.Database) bool {
	return db != nil && cctx.VersionInfo.InstalledVersion.String() != ironic.Status.InstalledVersion
}

func newMigrationTemplate(cctx ControllerContext, ironic *metal3api.Ironic, database *metal3api.Database, phase upgradePhase) corev1.PodTemplateSpec {
	script := commandPerPhase[phase]

	volumes, mounts := databaseClientMounts(database)

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

func ensureIronicUpgradeJob(cctx ControllerContext, ironic *metal3api.Ironic, db *metal3api.Database, phase upgradePhase) (Status, error) {
	if !upgradeJobRequired(cctx, ironic, db) {
		return ready()
	}

	fromVersion := ironic.Status.InstalledVersion
	if fromVersion == "" {
		fromVersion = "none"
	}
	toVersion := cctx.VersionInfo.InstalledVersion

	// TODO(dtantsur): remove this when ironic-image < 29.0 is not supported
	if toVersion.Compare(versionUpgradeScripts) < 0 {
		cctx.Logger.Info("not running upgrade scripts: the new version does not support them", "From", fromVersion, "To", toVersion.String())
		return ready()
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s-%s-to-%s", ironic.Name, phase, fromVersion, cctx.VersionInfo.InstalledVersion),
			Namespace: ironic.Namespace,
		},
	}

	template := newMigrationTemplate(cctx, ironic, db, phase)

	result, err := controllerutil.CreateOrUpdate(cctx.Context, cctx.Client, job, func() error {
		if job.ObjectMeta.Labels == nil {
			cctx.Logger.Info("creating a new upgrade job", "Phase", phase, "From", fromVersion, "To", toVersion.String())
			job.ObjectMeta.Labels = make(map[string]string, 2)
		}
		job.ObjectMeta.Labels[metal3api.IronicAppLabel] = ironicDeploymentName(ironic)
		job.ObjectMeta.Labels[metal3api.IronicVersionLabel] = cctx.VersionInfo.InstalledVersion.String()

		job.Spec.TTLSecondsAfterFinished = ptr.To(jobTTLSeconds)
		mergePodTemplates(&job.Spec.Template, template)

		return controllerutil.SetControllerReference(ironic, job, cctx.Scheme)
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
