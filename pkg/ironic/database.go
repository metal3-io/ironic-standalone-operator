package ironic

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	metal3api "github.com/metal3-io/ironic-operator/api/v1alpha1"
)

const (
	databasePort = 3306
	databaseUser = 27
)

func databaseDeploymentName(db *metal3api.IronicDatabase) string {
	return fmt.Sprintf("%s-database", db.Name)
}

func DatabaseDNSName(db *metal3api.IronicDatabase, domain string) string {
	if domain != "" && domain[0] != '.' {
		domain = fmt.Sprintf(".%s", domain)
	}
	return fmt.Sprintf("%s.%s.%s%s:%d", databaseDeploymentName(db), db.Namespace, serviceDNSSuffix, domain, databasePort)
}

func commonDatabaseVars(db *metal3api.IronicDatabase) []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name:  "MARIADB_DATABASE",
			Value: db.Name,
		},
		{
			Name: "MARIADB_USER",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: db.Spec.CredentialsRef,
					Key:                  "username",
				},
			},
		},
		{
			Name: "MARIADB_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: db.Spec.CredentialsRef,
					Key:                  "password",
				},
			},
		},
	}

}

func newDatabasePodTemplate(db *metal3api.IronicDatabase) corev1.PodTemplateSpec {
	volumes := []corev1.Volume{}
	mounts := []corev1.VolumeMount{}

	if db.Spec.TLSRef.Name != "" {
		volumes = append(volumes, corev1.Volume{
			Name: "cert-mariadb",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: db.Spec.TLSRef.Name,
				},
			},
		})
		mounts = append(mounts, corev1.VolumeMount{
			Name:      "cert-mariadb",
			MountPath: "/certs/mariadb",
		})
	}

	envVars := commonDatabaseVars(db)
	envVars = append(envVars, []corev1.EnvVar{
		{
			Name:  "MARIADB_HOST",
			Value: "%",
		},
		{
			Name:  "RESTART_CONTAINER_CERTIFICATE_UPDATED",
			Value: "true",
		},
	}...)

	probe := newProbe(corev1.ProbeHandler{
		Exec: &corev1.ExecAction{
			Command: []string{"sh", "-c", "mysqladmin status -u$(printenv MARIADB_USER) -p$(printenv MARIADB_PASSWORD)"},
		},
	})

	containers := []corev1.Container{
		{
			Name:            "mariadb",
			Image:           db.Spec.Image,
			ImagePullPolicy: corev1.PullAlways,
			Env:             envVars,
			VolumeMounts:    mounts,
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:  pointer.Int64(databaseUser),
				RunAsGroup: pointer.Int64(databaseUser),
			},
			LivenessProbe:  probe,
			ReadinessProbe: probe,
		},
	}

	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{metal3api.IronicOperatorLabel: databaseDeploymentName(db)},
		},
		Spec: corev1.PodSpec{
			Containers: containers,
			Volumes:    volumes,
		},
	}
}

func ensureDatabaseDeployment(cctx ControllerContext, db *metal3api.IronicDatabase) (metal3api.IronicStatusConditionType, error) {
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: databaseDeploymentName(db), Namespace: db.Namespace},
	}
	_, err := controllerutil.CreateOrUpdate(cctx.Context, cctx.Client, deploy, func() error {
		if deploy.ObjectMeta.CreationTimestamp.IsZero() {
			cctx.Logger.Info("creating a new deployment")
			matchLabels := map[string]string{metal3api.IronicOperatorLabel: databaseDeploymentName(db)}
			deploy.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: matchLabels,
			}
			deploy.Spec.Replicas = pointer.Int32(1)
		}
		deploy.Spec.Template = newDatabasePodTemplate(db)
		return controllerutil.SetControllerReference(db, deploy, cctx.Scheme)
	})
	if err != nil {
		return metal3api.IronicStatusProgressing, err
	}
	return getDeploymentStatus(deploy)
}

func ensureDatabaseService(cctx ControllerContext, db *metal3api.IronicDatabase) (metal3api.IronicStatusConditionType, error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: databaseDeploymentName(db), Namespace: db.Namespace},
	}
	_, err := controllerutil.CreateOrUpdate(cctx.Context, cctx.Client, service, func() error {
		if service.ObjectMeta.Labels == nil {
			cctx.Logger.Info("creating a new service")
			service.ObjectMeta.Labels = make(map[string]string)
		}
		service.ObjectMeta.Labels[metal3api.IronicOperatorLabel] = databaseDeploymentName(db)

		service.Spec.Selector = map[string]string{metal3api.IronicOperatorLabel: databaseDeploymentName(db)}
		service.Spec.Ports = []corev1.ServicePort{
			{
				Protocol: corev1.ProtocolTCP,
				Port:     databasePort,
			},
		}
		service.Spec.Type = corev1.ServiceTypeClusterIP

		return controllerutil.SetControllerReference(db, service, cctx.Scheme)
	})
	if err != nil || len(service.Spec.ClusterIPs) == 0 {
		return metal3api.IronicStatusProgressing, err
	}
	return metal3api.IronicStatusAvailable, nil
}

// EnsureDatabase ensures MariaDB is running with the current configuration.
func EnsureDatabase(cctx ControllerContext, db *metal3api.IronicDatabase) (status metal3api.IronicStatusConditionType, err error) {
	status, err = ensureDatabaseDeployment(cctx, db)
	if err != nil || status != metal3api.IronicStatusAvailable {
		return
	}

	return ensureDatabaseService(cctx, db)
}

// RemoveDatabase removes the MariaDB database.
func RemoveDatabase(cctx ControllerContext, db *metal3api.IronicDatabase) error {
	return nil // rely on ownership-based clean up
}
