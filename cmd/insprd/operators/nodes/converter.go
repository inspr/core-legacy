package nodes

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/auth"
	"inspr.dev/inspr/pkg/environment"
	"inspr.dev/inspr/pkg/meta"
	metautils "inspr.dev/inspr/pkg/meta/utils"
	"inspr.dev/inspr/pkg/operator/k8s"
	"inspr.dev/inspr/pkg/sidecars/models"
	"inspr.dev/inspr/pkg/utils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var lbsidecarPort int32

func (no *NodeOperator) dappToService(app *meta.App) *kubeService {
	logger.Info("creating kubernetes service")

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

// dAppToDeployment translates the DApp to a k8s deployment
func (no *NodeOperator) dAppToDeployment(
	app *meta.App,
	usePermTree bool,
) *kubeDeployment {
	appDeployName := toDeploymentName(app)
	appLabels := map[string]string{
		"inspr-app": toAppID(app),
	}
	logger.Info("constructing deployment", zap.Bool("useperm", usePermTree))

	nodeContainer := createNodeContainer(app, appDeployName)
	scContainers := no.withAllSidecarsContainers(
		app,
		appDeployName,
		usePermTree,
	)

	return (*kubeDeployment)(
		k8s.NewDeployment(
			appDeployName,
			k8s.WithLabels(appLabels),
			k8s.WithContainer(
				append(scContainers, nodeContainer)...,
			),
			k8s.WithReplicas(app.Spec.Node.Spec.Replicas),
		))
}

func (no *NodeOperator) withAllSidecarsContainers(
	app *meta.App,
	appDeployName string,
	usePermTree bool,
) []corev1.Container {
	var containers []corev1.Container
	var sidecarAddrs []corev1.EnvVar
	for _, broker := range no.getAllSidecarBrokers(app, usePermTree) {

		factory, err := no.brokers.Factory().Get(broker)

		if err != nil {
			panic(fmt.Sprintf("broker %v not allowed: %v", broker, err))
		}

		logger.Info(
			"with all sidecars containers",
			zap.Bool("useperm", usePermTree),
		)

		container, addrEnvVar := factory(app,
			getAvailiblePorts(),
			no.withBoundary(app, usePermTree),
			k8s.ContainerWithEnv(corev1.EnvVar{
				Name:  "LOG_LEVEL",
				Value: app.Spec.LogLevel,
			}),
			withLBSidecarConfiguration())

		containers = append(containers, container)
		sidecarAddrs = append(sidecarAddrs, addrEnvVar...)
	}

	lbSidecar := k8s.NewContainer(
		appDeployName+"-lbsidecar",
		"",
		no.withLBSidecarImage(app),
		no.withBoundary(app, usePermTree),
		withLBSidecarPorts(app),
		withLBSidecarConfiguration(),
		k8s.ContainerWithEnv(sidecarAddrs...),
		withNodeID(app),
		k8s.ContainerWithEnv(corev1.EnvVar{
			Name:  "LOG_LEVEL",
			Value: app.Spec.LogLevel,
		}),
		k8s.ContainerWithPullPolicy(corev1.PullAlways),
	)

	containers = append(containers, lbSidecar)

	return containers
}

func (no *NodeOperator) getAllSidecarBrokers(
	app *meta.App,
	usePermTree bool,
) utils.StringArray {
	input := app.Spec.Boundary.Input
	output := app.Spec.Boundary.Output
	channels := input.Union(output)

	logger.Debug("resolving Node Boundary in the cluster",
		zap.String("operation", "getAllSidecarBrokers"),
		zap.Bool("useperm", usePermTree),
		zap.String("app:", app.Meta.Name),
	)

	resolves, err := no.memory.Apps().ResolveBoundary(app, usePermTree)
	if err != nil {
		logger.Error("unable to resolve Node boundaries",
			zap.Any("boundaries", app.Spec.Boundary))
		panic(err)
	}

	set, _ := metautils.MakeStrSet(channels.Map(func(boundary string) string {
		resolved := resolves[boundary]
		parent, chName, _ := metautils.RemoveLastPartInScope(resolved)
		var ch *meta.Channel
		if usePermTree {
			ch, err = no.memory.Perm().Channels().Get(parent, chName)
		} else {
			ch, err = no.memory.Channels().Get(parent, chName)
		}
		if err != nil {
			logger.Error("unable get channel for boudary resolution",
				zap.String("channel", chName))
			panic(err)
		}
		return ch.Spec.SelectedBroker
	}))
	return set.ToArray()
}

// withBoundary adds the boundary configuration to the kubernetes' deployment environment variables
func (no *NodeOperator) withBoundary(
	app *meta.App,
	usePermTree bool,
) k8s.ContainerOption {
	scope, _ := metautils.JoinScopes(app.Meta.Parent, app.Meta.Name)
	if _, err := no.memory.Apps().Get(scope); err != nil {
		return nil
	}
	return func(c *corev1.Container) {
		input := app.Spec.Boundary.Input
		output := app.Spec.Boundary.Output
		channels := input.Union(output)

		logger.Debug("with boundary", zap.Bool("useperm", usePermTree))

		resolves, err := no.memory.Apps().ResolveBoundary(app, usePermTree)
		if err != nil {
			logger.Error("unable to resolve Node boundaries",
				zap.Any("boundaries", app.Spec.Boundary))
			panic(err)
		}

		inputEnv := input.Map(func(boundary string) string {
			return no.returnChannelBroker(boundary, resolves[boundary])
		})

		outputEnv := output.Map(func(boundary string) string {
			return no.returnChannelBroker(boundary, resolves[boundary])
		})

		env := utils.EnvironmentMap{
			"INSPR_INPUT_CHANNELS":  inputEnv.Join(";"),
			"INSPR_OUTPUT_CHANNELS": outputEnv.Join(";"),
		}

		logger.Debug(
			"resolving with Node Boundary",
			zap.Bool("useperm", usePermTree),
		)
		channels.Map(func(boundary string) string {
			resolved := resolves[boundary]
			parent, chName, _ := metautils.RemoveLastPartInScope(resolved)
			var ch *meta.Channel
			var ct *meta.Type
			var cterr, cherr error
			if usePermTree {
				ch, cherr = no.memory.Perm().Channels().Get(parent, chName)
				ct, cterr = no.memory.Perm().Types().Get(parent, ch.Spec.Type)
			} else {
				ch, cherr = no.memory.Channels().Get(parent, chName)
				ct, cterr = no.memory.Types().Get(parent, ch.Spec.Type)
			}
			if cherr != nil {
				logger.Error(
					"Unable to get channel to resolve with boundary",
					zap.String("channel", chName),
				)
				panic(err)
			}
			if cterr != nil {
				logger.Error(
					"Unable to get channel type to resolve with boundary",
					zap.String("type", ch.Spec.Type),
				)
				panic(err)
			}
			resolved = "INSPR_" + ch.Meta.UUID
			env[resolved+"_SCHEMA"] = ct.Schema
			env[boundary+"_RESOLVED"] = resolved
			return boundary
		})
		logger.Debug(
			"resolved with Node Boundary",
			zap.Bool("useperm", usePermTree),
		)

		c.Env = append(c.Env, env.ParseToK8sArrEnv()...)
	}
}

// withLBSidecarImage adds the sidecar image to the dApp
func (no *NodeOperator) withLBSidecarImage(app *meta.App) k8s.ContainerOption {
	return func(c *corev1.Container) {
		c.Image = environment.GetSidecarImage()
	}
}

func (no *NodeOperator) returnChannelBroker(
	channel, pathToResolvedChannel string,
) string {
	scope, chName, err := metautils.RemoveLastPartInScope(pathToResolvedChannel)
	if err != nil {
		return ""
	}
	resolvedCh, err := no.memory.Channels().Get(scope, chName)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s@%s", channel, resolvedCh.Spec.SelectedBroker)
}

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
		Refresh: []byte(scope),
		RefreshURL: fmt.Sprintf(
			"http://%v/refreshController",
			os.Getenv("INSPR_INSPRD_ADDRESS"),
		),
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

// Auxiliar methods

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

func createNodeContainer(app *meta.App, appDeployName string) corev1.Container {
	return k8s.NewContainer(
		appDeployName,
		app.Spec.Node.Spec.Image,
		withLBSidecarPorts(app),
		withSecretDefinition(app),
		k8s.ContainerWithEnv(corev1.EnvVar{
			Name:  "LOG_LEVEL",
			Value: app.Spec.LogLevel,
		}),
		withLBSidecarConfiguration(),
	)
}

func getAvailiblePorts() *models.SidecarConnections {
	ports, err := utils.GetFreePorts(2)
	if err != nil {
		logger.Error("unable to get free ports for broker sidecar: %v",
			zap.Any("error", err))

		panic(fmt.Sprintf("error while getting free ports: %v", err))

	}
	return &models.SidecarConnections{
		InPort:  int32(ports[0]),
		OutPort: int32(ports[1]),
	}
}

func withLBSidecarConfiguration() k8s.ContainerOption {
	return k8s.ContainerWithEnvFrom(
		corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: os.Getenv("INSPR_LBSIDECAR_CONFIGMAP"),
				},
			},
		},
	)
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
	}
}

func withNodeID(app *meta.App) k8s.ContainerOption {
	return k8s.ContainerWithEnv(corev1.EnvVar{
		Name:  "INSPR_APP_ID",
		Value: toAppID(app),
	})
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
