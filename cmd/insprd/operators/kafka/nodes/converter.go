package nodes

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/inspr/inspr/pkg/auth"
	"github.com/inspr/inspr/pkg/environment"
	"github.com/inspr/inspr/pkg/meta"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
	"github.com/inspr/inspr/pkg/operator/k8s"
	"github.com/inspr/inspr/pkg/utils"
	"go.uber.org/zap"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (no *NodeOperator) toSecret(app *meta.App) *kubeSecret {
	log.Println("creating secret")
	scope, err := metautils.JoinScopes(app.Meta.Parent, app.Meta.Name)
	if err != nil {
		log.Printf("err = %+v\n", err)
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
		log.Printf("err = %+v\n", err)
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
		inputEnv := input.Join(";")
		outputEnv := output.Join(";")
		env := utils.EnvironmentMap{
			"INSPR_INPUT_CHANNELS":  inputEnv,
			"INSPR_OUTPUT_CHANNELS": outputEnv,
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

// withSidecarPorts adds the sidecar ports if they are defined in the dApp definitions.
// On kubernetes, this onverrides the defined configuration on the configmap
func withSidecarPorts(app *meta.App) k8s.ContainerOption {
	return func(c *corev1.Container) {
		writePort := app.Spec.Node.Spec.SidecarPort.Write
		readPort := app.Spec.Node.Spec.SidecarPort.Read

		if writePort > 0 {
			c.Env = append(c.Env, corev1.EnvVar{
				Name:  "INSPR_SIDECAR_WRITE_PORT",
				Value: strconv.Itoa(writePort),
			})
		}
		if readPort > 0 {
			c.Env = append(c.Env, corev1.EnvVar{
				Name:  "INSPR_SIDECAR_READ_PORT",
				Value: strconv.Itoa(readPort),
			})
		}
	}
}

// withSidecarImage adds the sidecar image to the dApp
func (no *NodeOperator) withSidecarImage(app *meta.App) k8s.ContainerOption {
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
	log.Println("constructing deployment")

	return (*kubeDeploy)(
		k8s.NewDeployment(
			appDeployName,
			k8s.WithLabels(appLabels),
			k8s.WithContainer(
				k8s.NewContainer(
					appDeployName,
					app.Spec.Node.Spec.Image,
					withSidecarPorts(app),
					withSecretDefinition(app),
					withSidecarConfiguration(),
				),
				k8s.NewContainer(
					appDeployName+"-sidecar",
					"",
					no.withSidecarImage(app),
					no.withBoundary(app),
					withSidecarPorts(app),
					withKafkaConfiguration(),
					withSidecarConfiguration(),
					withNodeID(app),
					k8s.ContainerWithPullPolicy(corev1.PullAlways),
				),
			),
		))
}

var sidecarPort int32

func withKafkaConfiguration() k8s.ContainerOption {
	return k8s.ContainerWithEnvFrom(
		corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "inspr-kafka-configuration",
				},
			},
		},
	)
}

func withSidecarConfiguration() k8s.ContainerOption {
	return k8s.ContainerWithEnvFrom(
		corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "inspr-sidecar-configuration",
				},
			},
		},
	)
}

func dappToService(app *meta.App) *kubeService {
	temp, _ := strconv.Atoi(os.Getenv("INSPR_SIDECAR_PORT"))
	sidecarPort = int32(temp)
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
					Name:       "sidecar-port",
					Port:       sidecarPort,
					TargetPort: intstr.FromInt(int(sidecarPort)),
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
