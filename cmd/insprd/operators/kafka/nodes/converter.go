package nodes

import (
	"strings"

	kafkasc "github.com/inspr/inspr/cmd/sidecars/kafka/client"
	"github.com/inspr/inspr/pkg/environment"
	"github.com/inspr/inspr/pkg/meta"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
	"github.com/inspr/inspr/pkg/utils"
	"go.uber.org/zap"

	kubeApp "k8s.io/api/apps/v1"
	kubeCore "k8s.io/api/core/v1"
	kubeMeta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (no *NodeOperator) baseEnvironment(app *meta.App) utils.EnvironmentMap {
	logger.Debug("getting necessary environment variables for Node structure deployment")
	input := app.Spec.Boundary.Input
	output := app.Spec.Boundary.Output
	channels := input.Union(output)

	// label name to be used in the service
	appID := toAppID(app)
	inputEnv := input.Join(";")
	outputEnv := output.Join(";")
	env := utils.EnvironmentMap{
		"INSPR_INPUT_CHANNELS":    inputEnv,
		"INSPR_OUTPUT_CHANNELS":   outputEnv,
		"INSPR_SIDECAR_IMAGE":     environment.GetSidecarImage(),
		"INSPR_APP_ID":            appID,
		"INSPR_APP_CTX":           app.Meta.Parent,
		"INSPR_ENV":               environment.GetInsprEnvironment(),
		"KAFKA_BOOTSTRAP_SERVERS": kafkasc.GetEnvironment().KafkaBootstrapServers,
		"KAFKA_AUTO_OFFSET_RESET": kafkasc.GetEnvironment().KafkaAutoOffsetReset,
	}

	resolves, err := no.memory.Apps().ResolveBoundary(app)
	if err != nil {
		logger.Error("unable to resolve Node boundaries",
			zap.Any("boundaries", app.Spec.Boundary))
		panic(err)
	}

	logger.Debug("resolving Node Boundary in the cluster")
	channels.Map(func(boundary string) string {
		resolved := resolves[boundary]
		parent, chName, _ := metautils.RemoveLastPartInScope(resolved)
		ch, _ := no.memory.Channels().Get(parent, chName)
		ct, _ := no.memory.ChannelTypes().Get(parent, ch.Spec.Type)
		resolved = "INSPR_" + ch.Meta.UUID
		env[resolved+"_SCHEMA"] = ct.Schema
		env[boundary+"_RESOLVED"] = resolved
		return boundary
	})
	return env
}

// dAppToDeployment translates the DApp
func (no *NodeOperator) dAppToDeployment(app *meta.App) *kubeApp.Deployment {
	logger.Debug("converting a dApp structure to a k8s deployment")

	sidecarEnvironment := no.baseEnvironment(app)

	logger.Debug("defining Node's env vars")
	nodeKubeEnv := append(app.Spec.Node.Spec.Environment.ParseToK8sArrEnv(), kubeCore.EnvVar{
		Name: "INSPR_UNIX_SOCKET",
		ValueFrom: &kubeCore.EnvVarSource{
			FieldRef: &kubeCore.ObjectFieldSelector{
				FieldPath: "metadata.name",
			},
		},
	})

	logger.Debug("defining Node's k8s container data")
	appDeployName := toDeploymentName(app)
	appID := toAppID(app)
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

	logger.Debug("defining Sidecars's k8s env vars")
	sidecarKubeEnv := append(sidecarEnvironment.ParseToK8sArrEnv(), kubeCore.EnvVar{
		Name: "INSPR_UNIX_SOCKET",
		ValueFrom: &kubeCore.EnvVarSource{
			FieldRef: &kubeCore.ObjectFieldSelector{
				FieldPath: "metadata.name",
			},
		},
	})

	logger.Debug("defining Sidecar's k8s container data")
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

	appLabels := map[string]string{"app": appID}
	replicas := new(int32)

	if app.Spec.Node.Spec.Replicas == 0 {
		app.Spec.Node.Spec.Replicas = 1
	}

	*replicas = int32(app.Spec.Node.Spec.Replicas)

	logger.Debug("building and returning the k8s complete deployment")
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

func dappToService(app *meta.App) *kubeCore.Service {
	appID := toAppID(app)
	appDeployName := toDeploymentName(app)
	appLabels := map[string]string{"app": appID}

	svc := &kubeCore.Service{
		ObjectMeta: kubeMeta.ObjectMeta{
			Name: appDeployName,
		},
		Spec: kubeCore.ServiceSpec{
			Ports: func() (ports []kubeCore.ServicePort) {
				for _, port := range app.Spec.Node.Spec.Ports {
					ports = append(ports, kubeCore.ServicePort{
						Port:       int32(port.Port),
						TargetPort: intstr.FromInt(port.TargetPort),
					})
				}
				return
			}(),
			Selector: appLabels,
		},
	}
	return svc
}

// toDeployment - creates the kubernetes deployment name from the app
func toDeploymentName(app *meta.App) string {

	return "node-" + app.Meta.UUID
}

// toAppID - creates the kubernetes deployment name from the app
func toAppID(app *meta.App) string {
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
