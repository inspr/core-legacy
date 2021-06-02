package nodes

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/inspr/inspr/pkg/auth"
	"github.com/inspr/inspr/pkg/environment"
	"github.com/inspr/inspr/pkg/meta"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
	"github.com/inspr/inspr/pkg/operator/k8s"
	"github.com/inspr/inspr/pkg/sidecars/models"
	"github.com/inspr/inspr/pkg/utils"
	"go.uber.org/zap"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (no *NodeOperator) toSecret(app *meta.App) *kubeSecret {
	logger.Info("creating secret")
	scope, err := metautils.JoinScopes(app.Meta.Parent, app.Meta.Name)
	if err != nil {
		logger.Error("invalid scope", zap.Any("error", err))
		return nil
	}

	payload := auth.Payload{
		UID: app.Meta.UUID,
		Permissions: map[string][]string{
			app.Spec.Auth.Scope: app.Spec.Auth.Permissions,
		},
		Refresh:    []byte(scope),
		RefreshURL: fmt.Sprintf("%v/refreshController", os.Getenv("INSPR_INSPRD_ADDRESS")),
	}

	token, err := no.auth.Tokenize(payload)
	if err != nil {
		logger.Error("unable to tokenize", zap.Any("error", err))
		return nil
	}

	return &kubeSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name: toDeploymentName(app),
		},
		Data: map[string][]byte{
			"INSPR_CONTROLLER_TOKEN": token,
			"INSPR_CONTROLLER_SCOPE": []byte(app.Spec.Auth.Scope),
		},
	}
}

func withSecretDefinition(app *meta.App) k8s.ContainerOption {
	env := corev1.EnvFromSource{
		SecretRef: &corev1.SecretEnvSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: toDeploymentName(app),
			},
		},
	}
	return k8s.ContainerWithEnvFrom(env)
}

// withBoundary adds the boundary configuration to the kubernetes' deployment environment variables
func (no *NodeOperator) withBoundary(app *meta.App) k8s.ContainerOption {
	scope, _ := metautils.JoinScopes(app.Meta.Parent, app.Meta.Name)
	if _, err := no.memory.Apps().Get(scope); err != nil {
		return nil
	}
	return func(c *corev1.Container) {

		input := app.Spec.Boundary.Input
		output := app.Spec.Boundary.Output
		channels := input.Union(output)

		// label name to be used in the service

		resolves, err := no.memory.Apps().ResolveBoundary(app)
		if err != nil {
			logger.Error("unable to resolve Node boundaries",
				zap.Any("boundaries", app.Spec.Boundary))
			panic(err)
		}

		inputEnv := input.Map(func(boundary string) string {
			return no.returnChannelBroker(resolves[boundary])
		})

		outputEnv := output.Map(func(boundary string) string {
			return no.returnChannelBroker(resolves[boundary])
		})

		env := utils.EnvironmentMap{
			"INSPR_INPUT_CHANNELS":  inputEnv.Join(";"),
			"INSPR_OUTPUT_CHANNELS": outputEnv.Join(";"),
		}

		logger.Debug("resolving Node Boundary in the cluster")
		channels.Map(func(boundary string) string {
			resolved := resolves[boundary]
			parent, chName, _ := metautils.RemoveLastPartInScope(resolved)
			ch, _ := no.memory.Channels().Get(parent, chName)
			ct, _ := no.memory.Types().Get(parent, ch.Spec.Type)
			resolved = "INSPR_" + ch.Meta.UUID
			env[resolved+"_SCHEMA"] = ct.Schema
			env[boundary+"_RESOLVED"] = resolved
			return boundary
		})

		c.Env = append(c.Env, env.ParseToK8sArrEnv()...)
	}
}

func withNodeID(app *meta.App) k8s.ContainerOption {
	return k8s.ContainerWithEnv(corev1.EnvVar{
		Name:  "INSPR_APP_ID",
		Value: toAppID(app),
	})
}

// withLBSidecarPorts adds the load balancer sidecar ports if they are defined in the dApp definitions.
// On kubernetes, this overrides the defined configuration on the configmap
func withLBSidecarPorts(app *meta.App) k8s.ContainerOption {
	return func(c *corev1.Container) {
		lbWritePort := app.Spec.Node.Spec.SidecarPort.LBWrite
		lbReadPort := app.Spec.Node.Spec.SidecarPort.LBRead

		if lbWritePort > 0 {
			c.Env = append(c.Env, corev1.EnvVar{
				Name:  "INSPR_LBSIDECAR_WRITE_PORT",
				Value: strconv.Itoa(lbWritePort),
			})
		}
		if lbReadPort > 0 {
			c.Env = append(c.Env, corev1.EnvVar{
				Name:  "INSPR_LBSIDECAR_READ_PORT",
				Value: strconv.Itoa(lbReadPort),
			})
		}
		// The Sidecar Client read port must be added here
	}
}

// withLBSidecarImage adds the sidecar image to the dApp
func (no *NodeOperator) withLBSidecarImage(app *meta.App) k8s.ContainerOption {
	return func(c *corev1.Container) {
		c.Image = environment.GetSidecarImage()
	}
}

// dAppToDeployment translates the DApp
func (no *NodeOperator) dAppToDeployment(app *meta.App) *kubeDeploy {
	appDeployName := toDeploymentName(app)
	appLabels := map[string]string{
		"inspr-app": toAppID(app),
	}
	logger.Info("constructing deployment")

	nodeContainer := createNodeContainer(app, appDeployName)
	scContainers := no.withAllSidecarsContainers(app, appDeployName)

	return (*kubeDeploy)(
		k8s.NewDeployment(
			appDeployName,
			k8s.WithLabels(appLabels),
			k8s.WithContainer(
				append(scContainers, nodeContainer)...,
			),
		))
}

func withLBSidecarConfiguration() k8s.ContainerOption {
	return k8s.ContainerWithEnvFrom(
		corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "inspr-lbsidecar-configuration",
				},
			},
		},
	)
}

var lbsidecarPort int32

func dappToService(app *meta.App) *kubeService {
	temp, _ := strconv.Atoi(os.Getenv("INSPR_LBSIDECAR_PORT"))
	lbsidecarPort = int32(temp)
	appID := toAppID(app)
	appDeployName := toDeploymentName(app)
	appLabels := map[string]string{"inspr-app": appID}

	svc := &kubeService{
		ObjectMeta: metav1.ObjectMeta{
			Name: appDeployName,
		},
		Spec: corev1.ServiceSpec{
			Ports: func() (ports []corev1.ServicePort) {
				for i, port := range app.Spec.Node.Spec.Ports {
					ports = append(ports, corev1.ServicePort{
						Name:       fmt.Sprintf("port%v", i),
						Port:       int32(port.Port),
						TargetPort: intstr.FromInt(port.TargetPort),
					})
				}
				ports = append(ports, corev1.ServicePort{
					Name:       "lbsidecar-port",
					Port:       lbsidecarPort,
					TargetPort: intstr.FromInt(int(lbsidecarPort)),
				})
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

func (no *NodeOperator) returnChannelBroker(pathToChannel string) string {
	scope, chName, err := metautils.RemoveLastPartInScope(pathToChannel)
	if err != nil {
		return ""
	}
	channel, err := no.memory.Channels().Get(scope, chName)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s_%s", chName, channel.Spec.SelectedBroker)
}

func getAvailiblePorts() *models.SidecarConnections {
	return &models.SidecarConnections{
		InPort:  42069,
		OutPort: 42070,
	}
}

func (no *NodeOperator) getAllSidecarNames(app *meta.App) utils.StringArray {
	input := app.Spec.Boundary.Input
	output := app.Spec.Boundary.Output
	channels := input.Union(output)

	resolves, err := no.memory.Apps().ResolveBoundary(app)
	if err != nil {
		logger.Error("unable to resolve Node boundaries",
			zap.Any("boundaries", app.Spec.Boundary))
		panic(err)
	}

	logger.Debug("resolving Node Boundary in the cluster")

	set, _ := metautils.MakeStrSet(channels.Map(func(boundary string) string {
		resolved := resolves[boundary]
		parent, chName, _ := metautils.RemoveLastPartInScope(resolved)
		ch, _ := no.memory.Channels().Get(parent, chName)
		return ch.Spec.SelectedBroker
	}))
	return set.ToArray()
}

func (no *NodeOperator) withAllSidecarsContainers(app *meta.App, appDeployName string) []corev1.Container {
	var containers []corev1.Container
	var sidecarAddrs []corev1.EnvVar
	for _, broker := range no.getAllSidecarNames(app) {

		factory, err := no.brokers.Factory().Get(broker)

		if err != nil {
			panic("broker not allowed")
		}

		container, addrEnvVar := factory(app,
			getAvailiblePorts(),
			no.withBoundary(app),
			withLBSidecarConfiguration())

		containers = append(containers, container)
		sidecarAddrs = append(sidecarAddrs, addrEnvVar...)
	}

	lbSidecar := k8s.NewContainer(
		appDeployName+"-lbsidecar",
		"",
		no.withLBSidecarImage(app),
		no.withBoundary(app),
		withLBSidecarPorts(app),
		withLBSidecarConfiguration(),
		k8s.ContainerWithEnv(sidecarAddrs...),
		withNodeID(app),
		k8s.ContainerWithPullPolicy(corev1.PullAlways),
	)

	containers = append(containers, lbSidecar)

	return containers
}

func createNodeContainer(app *meta.App, appDeployName string) corev1.Container {
	return k8s.NewContainer(
		appDeployName,
		app.Spec.Node.Spec.Image,
		withLBSidecarPorts(app),
		withSecretDefinition(app),
		withLBSidecarConfiguration(),
	)
}
