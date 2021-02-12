package operator

import (
	"fmt"
	"strings"

	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
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

	// pod env variables
	insprEnv := environment.GetEnvironment()
	// label name to be used in the service
	appDeployName := toDeploymentName(insprEnv.InsprEnvironment, app.Meta.Parent, app.Spec.Node.Meta.Name)

	sidecarEnvironment := map[string]string{
		"INSPR_INPUT_CHANNELS":  inputChannels,
		"INSPR_CHANNEL_SIDECAR": insprEnv.SidecarImage,
		"INSPR_APPS_TLS":        "true",

		"INSPR_OUTPUT_CHANNELS": outputChannels,
		"INSPR_APP_ID":          appDeployName,
	}

	return &kubeApp.Deployment{
		ObjectMeta: kubeMeta.ObjectMeta{
			Name:   appDeployName,
			Labels: map[string]string{"app": appDeployName},
		},
		Spec: kubeApp.DeploymentSpec{
			Replicas: intToint32(app.Spec.Node.Spec.Replicas),
			Selector: &kubeMeta.LabelSelector{
				MatchLabels: map[string]string{
					"app": appDeployName,
				},
			},
			Strategy: kubeApp.DeploymentStrategy{
				Type: kubeApp.RollingUpdateDeploymentStrategyType,
			},
			Template: kubeCore.PodTemplateSpec{
				ObjectMeta: kubeMeta.ObjectMeta{
					Labels: map[string]string{
						"app": appDeployName,
					},
				},
				Spec: kubeCore.PodSpec{
					Volumes: []kubeCore.Volume{
						{
							Name: appDeployName + "-volume",
							VolumeSource: kubeCore.VolumeSource{
								EmptyDir: &kubeCore.EmptyDirVolumeSource{
									Medium: kubeCore.StorageMediumMemory,
								},
							},
						},
					},
					Containers: []kubeCore.Container{
						{
							Name: app.Spec.Node.Meta.Name,
							Ports: func() []kubeCore.ContainerPort {
								return nil
							}(),

							Image: app.Spec.Node.Spec.Image,
							// parse from master env var to kube env vars
							ImagePullPolicy: kubeCore.PullAlways,
							VolumeMounts: []kubeCore.VolumeMount{
								{
									Name:      app.Spec.Node.Meta.Name + "-volume",
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
							}, parseToK8sArrEnv(app.Spec.Node.Spec.Environment)...),
						},
						{
							Name:            appDeployName + "-sidecar",
							Image:           insprEnv.SidecarImage,
							ImagePullPolicy: kubeCore.PullIfNotPresent,
							VolumeMounts: []kubeCore.VolumeMount{
								{
									Name:      app.Spec.Node.Meta.Name + "-sidecar-volume",
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

func parseToNodeEnviroment(envs []kubeCore.EnvVar) map[string]string {
	nodeEnv := make(map[string]string)
	for _, env := range envs {
		nodeEnv[env.Name] = env.Value
	}
	return nodeEnv
}

// toDeployment - receives the context of an app and it's context
// creates a unique deployment name to be used in the k8s deploy
func toDeploymentName(insprEnv string, context string, name string) string {
	s := fmt.Sprintf("%s.%s.%s", insprEnv, context, name)
	if s[0] == '.' {
		s = s[1:]
	}
	return strings.ToLower(s)
}

// intToint32 - converts an integer to a *int32
func intToint32(v int) *int32 {
	t := int32(v)
	return &t
}

func ToNode(kdep *kubeApp.Deployment) (*meta.Node, error) {
	var err error
	node := &meta.Node{}
	node.Meta.Name, err = toNodeName(kdep.ObjectMeta.Name)
	if err != nil {
		return &meta.Node{}, err
	}
	node.Spec.Replicas = int(*kdep.Spec.Replicas)
	node.Spec.Image = kdep.Spec.Template.Spec.Containers[0].Image
	node.Spec.Environment = parseToNodeEnviroment(kdep.Spec.Template.Spec.Containers[0].Env)
	return node, nil
}

func toNodeName(deployName string) (string, error) {
	strs := strings.Split(deployName, ".")
	if len(strs) < 3 {
		return "", ierrors.NewError().Message("invalid deployment name").Build()
	}
	return strs[2], nil
}
