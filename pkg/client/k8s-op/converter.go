package operator

import (
	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/meta"

	kubeApp "k8s.io/api/apps/v1"
	kubeCore "k8s.io/api/core/v1"
	kubeMeta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InsprDAppToK8sDeployment translates the DApp
func InsprDAppToK8sDeployment(app *meta.App) *kubeApp.Deployment {
	inputChannels := ""
	for _, c := range app.Spec.Boundary.Input {
		inputChannels += c + ";"
	}
	outputChannels := ""
	for _, c := range app.Spec.Boundary.Output {
		outputChannels += c + ";"
	}

	envAppLabel := toDeployment(app)

	sidecarEnvironment := map[string]string{
		"INSPR_INPUT_CHANNELS": inputChannels,
		// "INSPR_DOMAIN":          insprEnv.Domain,
		// "INSPR_SUBDOMAIN":       insprEnv.Subdomain,
		// "INSPR_APPS_SUBDOMAIN":  insprEnv.AppsSubdomain,
		// "INSPR_CHANNEL_SIDECAR": insprEnv.ChannelSidecarImage,
		// "INSPR_APPS_TLS":        "true",

		"INSPR_OUTPUT_CHANNELS": outputChannels,
		"INSPR_app_ID":          app.Meta.Name, // TODO REVIEW
		// "INSPR_REGISTRY_URL":    insprEnv.RegistryURL,
		// "INSPR_LOG_CHANNEL":     insprEnv.LogChannel,
		// "INSPR_ENVIRONMENT":     insprEnv.Environment,
	}

	return &kubeApp.Deployment{
		ObjectMeta: kubeMeta.ObjectMeta{
			Name:   envAppLabel,
			Labels: map[string]string{"app": envAppLabel},
		},
		Spec: kubeApp.DeploymentSpec{
			// Replicas: &app.Replicas, // TODO REVIEW HOW TO ADD THIS
			Selector: &kubeMeta.LabelSelector{
				MatchLabels: map[string]string{
					"app": envAppLabel,
				},
			},
			Strategy: kubeApp.DeploymentStrategy{
				Type: kubeApp.RollingUpdateDeploymentStrategyType,
			},
			Template: kubeCore.PodTemplateSpec{

				ObjectMeta: kubeMeta.ObjectMeta{
					Labels: map[string]string{
						"app": envAppLabel,
					},
				},
				Spec: kubeCore.PodSpec{
					Volumes: []kubeCore.Volume{
						{
							Name: envAppLabel + "-volume",
							VolumeSource: kubeCore.VolumeSource{
								EmptyDir: &kubeCore.EmptyDirVolumeSource{
									Medium: kubeCore.StorageMediumMemory,
								},
							},
						},
					},
					Containers: []kubeCore.Container{
						{
							Name: envAppLabel,
							Ports: func() []kubeCore.ContainerPort {
								if !app.Exposed {
									return nil
								}
								return []kubeCore.ContainerPort{{
									ContainerPort: func() int32 {
										return app.Port
									}(),
								},
								}
							}(),

							Image: app.Image,
							// parse from master env var to kube env vars
							ImagePullPolicy: kubeCore.PullAlways,
							VolumeMounts: []kubeCore.VolumeMount{
								{
									Name:      envAppLabel + "-volume",
									MountPath: "/inspr",
								},
							},
							Env: append([]kubeCore.EnvVar{
								{
									Name: "UUID",
									ValueFrom: &kubeCore.EnvVarSource{
										FieldRef: &kubeCore.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
							}, parseToK8sArrEnv(app.Environment)...),
						},
						{
							Name:            envAppLabel + "-sidecar",
							Image:           environment.GetEnvironment().KafkaSidecarImage,
							ImagePullPolicy: kubeCore.PullIfNotPresent,
							VolumeMounts: []kubeCore.VolumeMount{
								{
									Name:      envAppLabel + "-volume",
									MountPath: "/inspr",
								},
							},
							Env: append(parseToK8sArrEnv(sidecarEnvironment), kubeCore.EnvVar{
								Name: "UUID",
								ValueFrom: &kubeCore.EnvVarSource{
									FieldRef: &kubeCore.ObjectFieldSelector{
										FieldPath: "metadata.name",
									},
								},
							}),
						},
					},
				},
			},
		},
	}
}

func parseToK8sArrEnv(arrappEnv map[string]string) []kubeCore.EnvVar {
	var arrEnv []kubeCore.EnvVar
	for key, val := range arrappEnv {
		arrEnv = append(arrEnv, kubeCore.EnvVar{
			Name:  key,
			Value: val,
		})
	}
	return arrEnv
}
