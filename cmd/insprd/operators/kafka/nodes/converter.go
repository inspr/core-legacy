package nodes

import (
	"strings"

	kafkasc "gitlab.inspr.dev/inspr/core/cmd/sidecars/kafka/client"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	metautils "gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"gitlab.inspr.dev/inspr/core/pkg/utils"

	kubeApp "k8s.io/api/apps/v1"
	kubeCore "k8s.io/api/core/v1"
	kubeMeta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (no *NodeOperator) baseEnvironment(app *meta.App) utils.EnvironmentMap {
	input := app.Spec.Boundary.Input
	output := app.Spec.Boundary.Output
	channels := input.Union(output)

	// label name to be used in the service
	appDeployName := toDeploymentName(app)

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
	resolves, err := no.memory.Apps().ResolveBoundary(app)
	if err != nil {
		panic(err)
	}
	channels.Map(func(boundary string) string {
		resolved := resolves[boundary]
		parent, chName, _ := metautils.RemoveLastPartInScope(resolved)
		ch, _ := no.memory.Channels().Get(parent, chName)
		ct, _ := no.memory.ChannelTypes().Get(parent, ch.Spec.Type)
		env[resolved+"_SCHEMA"] = ct.Schema
		env[boundary+"_RESOLVED"] = resolved
		return boundary
	})
	return env
}

// dAppToDeployment translates the DApp
func (no *NodeOperator) dAppToDeployment(app *meta.App) *kubeApp.Deployment {

	sidecarEnvironment := no.baseEnvironment(app)

	nodeKubeEnv := append(app.Spec.Node.Spec.Environment.ParseToK8sArrEnv(), kubeCore.EnvVar{
		Name: "INSPR_UNIX_SOCKET",
		ValueFrom: &kubeCore.EnvVarSource{
			FieldRef: &kubeCore.ObjectFieldSelector{
				FieldPath: "metadata.name",
			},
		},
	})

	appDeployName := toDeploymentName(app)
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
	replicas := new(int32)

	if app.Spec.Node.Spec.Replicas == 0 {
		app.Spec.Node.Spec.Replicas = 1
	}

	*replicas = int32(app.Spec.Node.Spec.Replicas)

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
			Replicas: replicas,
		},
	}
}

// toDeployment - creates the kubernetes deployment name from the app
func toDeploymentName(app *meta.App) string {
	var depNames utils.StringArray
	depNames = strings.Split(app.Meta.Parent, ".")
	if depNames[0] == "" {
		depNames = utils.StringArray{}
	}
	depNames = append(depNames, app.Meta.Name)
	return depNames.Join("-")
}

// intToint32 - converts an integer to a *int32
func intToint32(v int) *int32 {
	t := int32(v)
	return &t
}

func toNode(kdep *kubeApp.Deployment) (meta.Node, error) {
	var err error
	node := meta.Node{}
	node.Meta.Name, err = toNodeName(kdep.ObjectMeta.Name)
	if err != nil {
		return meta.Node{}, err
	}
	node.Meta.Parent, err = toNodeParent(kdep.ObjectMeta.Name)
	if err != nil {
		return meta.Node{}, err
	}
	node.Spec.Image = kdep.Spec.Template.Spec.Containers[0].Image
	node.Spec.Environment = utils.ParseFromK8sEnvironment(kdep.Spec.Template.Spec.Containers[0].Env)
	node.Spec.Replicas = int(*kdep.Spec.Replicas)
	return node, nil
}

func toNodeName(deployName string) (string, error) {
	strs := strings.Split(deployName, "-")
	return strs[len(strs)-1], nil
}

func toNodeParent(deployName string) (string, error) {

	strs := utils.StringArray(strings.Split(deployName, "-"))
	return strs[:len(strs)-1].Join("."), nil
}
