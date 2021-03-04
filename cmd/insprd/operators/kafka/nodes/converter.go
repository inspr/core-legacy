package nodes

import (
	"fmt"
	"strings"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/tree"
	kafkasc "gitlab.inspr.dev/inspr/core/cmd/sidecars/kafka/client"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/utils"

	kubeApp "k8s.io/api/apps/v1"
	kubeCore "k8s.io/api/core/v1"
	kubeMeta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func baseEnvironment(app *meta.App) utils.EnvironmentMap {
	input := app.Spec.Boundary.Input
	output := app.Spec.Boundary.Output
	channels := input.Union(output)

	// label name to be used in the service
	appDeployName := toDeploymentName(environment.GetInsprEnvironment(), app)

	inputEnv := input.Join(";")
	outputEnv := output.Join(";")
	env := utils.EnvironmentMap{
		"INSPR_INPUT_CHANNELS":    inputEnv,
		"INSPR_OUTPUT_CHANNELS":   outputEnv,
		"INSPR_SIDECAR_IMAGE":     environment.GetSidecarImage(),
		"INSPR_APP_ID":            appDeployName,
		"INSPR_APP_CTX":           app.Meta.Parent,
		"INSPR_ENV":               environment.GetInsprEnvironment(),
		"KAFKA_BOOTSTRAP_SERVERS": kafkasc.GetEnvironment().KafkaBootstrapServers,
		"KAFKA_AUTO_OFFSET_RESET": kafkasc.GetEnvironment().KafkaAutoOffsetReset,
	}
	channels.Map(func(s string) string {
		ch, _ := tree.GetTreeMemory().Channels().Get(app.Meta.Parent, s)
		ct, _ := tree.GetTreeMemory().ChannelTypes().Get(app.Meta.Parent, ch.Spec.Type)
		env[s+"_SCHEMA"] = ct.Schema
		return s
	})
	return env
}

// dAppToDeployment translates the DApp
func dAppToDeployment(app *meta.App) *kubeApp.Deployment {

	sidecarEnvironment := baseEnvironment(app)

	nodeKubeEnv := append(app.Spec.Node.Spec.Environment.ParseToK8sArrEnv(), kubeCore.EnvVar{
		Name: "INSPR_UNIX_SOCKET",
		ValueFrom: &kubeCore.EnvVarSource{
			FieldRef: &kubeCore.ObjectFieldSelector{
				FieldPath: "metadata.name",
			},
		},
	})

	appDeployName := toDeploymentName(environment.GetInsprEnvironment(), app)
	nodeContainer := kubeCore.Container{
		Name:  appDeployName,
		Image: app.Spec.Node.Spec.Image,
		// parse from master env var to kube env vars
		VolumeMounts: []kubeCore.VolumeMount{
			{
				Name:      appDeployName + "-volume",
				MountPath: "/inspr",
			},
		},
		Env: nodeKubeEnv,
	}

	sidecarKubeEnv := append(sidecarEnvironment.ParseToK8sArrEnv(), kubeCore.EnvVar{
		Name: "INSPR_UNIX_SOCKET",
		ValueFrom: &kubeCore.EnvVarSource{
			FieldRef: &kubeCore.ObjectFieldSelector{
				FieldPath: "metadata.name",
			},
		},
	})

	sidecarContainer := kubeCore.Container{
		Name:  appDeployName + "-sidecar",
		Image: environment.GetSidecarImage(),
		VolumeMounts: []kubeCore.VolumeMount{
			{
				Name:      appDeployName + "-volume",
				MountPath: "/inspr",
			},
		},
		Env: sidecarKubeEnv,
	}

	volume := kubeCore.Volume{
		Name: appDeployName + "-volume",
		VolumeSource: kubeCore.VolumeSource{
			EmptyDir: &kubeCore.EmptyDirVolumeSource{
				Medium: kubeCore.StorageMediumMemory,
			},
		},
	}

	appLabels := map[string]string{"app": appDeployName}

	return &kubeApp.Deployment{
		ObjectMeta: kubeMeta.ObjectMeta{
			Name:   appDeployName,
			Labels: appLabels,
		},
		Spec: kubeApp.DeploymentSpec{
			Selector: &kubeMeta.LabelSelector{
				MatchLabels: appLabels,
			},
			Template: kubeCore.PodTemplateSpec{
				ObjectMeta: kubeMeta.ObjectMeta{
					Labels: appLabels,
				},
				Spec: kubeCore.PodSpec{
					Volumes: []kubeCore.Volume{volume},
					Containers: []kubeCore.Container{
						sidecarContainer,
						nodeContainer,
					},
				},
			},
		},
	}
}

// toDeployment - receives the context of an app and it's context
// creates a unique deployment name to be used in the k8s deploy
func toDeploymentName(envPath string, app *meta.App) string {
	var arr utils.StringArray
	if envPath != "" {

		arr = utils.StringArray{
			"inspr",
			envPath,
			app.Meta.Parent,
			app.Meta.Name,
		}
	} else {
		arr = utils.StringArray{
			"inspr",
			app.Meta.Parent,
			app.Meta.Name,
		}
	}
	return arr.Join("-")
}

// intToint32 - converts an integer to a *int32
func intToint32(v int) *int32 {
	t := int32(v)
	return &t
}

func toNode(kdep *kubeApp.Deployment) (meta.Node, error) {
	var err error
	node := meta.Node{}
	fmt.Println(kdep.ObjectMeta.Name)
	node.Meta.Name, err = toNodeName(kdep.ObjectMeta.Name)
	if err != nil {
		return meta.Node{}, err
	}
	node.Spec.Image = kdep.Spec.Template.Spec.Containers[0].Image
	node.Spec.Environment = utils.ParseFromK8sEnviroment(kdep.Spec.Template.Spec.Containers[0].Env)
	return node, nil
}

func toNodeName(deployName string) (string, error) {
	strs := strings.Split(deployName, "-")
	if len(strs) < 3 {
		return "", ierrors.NewError().Message("invalid deployment name").Build()
	}
	return strs[2], nil
}
