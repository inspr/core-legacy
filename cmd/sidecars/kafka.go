package sidecars

import (
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/operator/k8s"
	"github.com/inspr/inspr/pkg/sidecar/models"
	corev1 "k8s.io/api/core/v1"
)

// KafkaConfig configurations used to create the KafkaSidecar
type KafkaConfig struct {
	bootstrapServers string
	autoOffsetReset  string
	sidecarImage     string

	// insprdName is the name used in the insprd service of your cluster
	kafkaInsprdName string
	// insprdPort is the port used in the insprd service of your cluster
	kafkaInsprPort string
	// namespace is the release namespace in which the insprd is located
	kafkaNamespace string
	// ports of the
	ports models.SidecarConnections
}

// KafkaToDeployment receives a the KafkaConfig variable as a parameter and returns a
// SidecarFactory function that is used to subscribe to the sidecarFactory
func KafkaToDeployment(config KafkaConfig) models.SidecarFactory {
	return func(app *meta.App, conn *models.SidecarConnections) k8s.DeploymentOption {
		return k8s.WithContainer(
			k8s.NewContainer(
				"sidecar-kafka-"+app.Meta.UUID, // deployment name
				config.sidecarImage,            // image url

				// label to the dApp associated with it
				insprAppIDConfig(app),
				kafkaConfig(config),
				sidecarConfig(config),
				k8s.ContainerWithPullPolicy(corev1.PullAlways),
			),
		)
	}
}

// kafkaConfig adds teh necessary env variables to configure kafka
func kafkaConfig(config KafkaConfig) k8s.ContainerOption {
	return k8s.ContainerWithEnv(
		corev1.EnvVar{
			Name:  "KAFKA_BOOTSTRAP_SERVERS",
			Value: config.bootstrapServers,
		},
		corev1.EnvVar{
			Name:  "KAFKA_AUTO_OFFSET_RESET",
			Value: config.autoOffsetReset,
		},
	)
}

// sidecarConfig adds the necessary env variables to configure the sidecar in the cluster
func sidecarConfig(config KafkaConfig) k8s.ContainerOption {
	return k8s.ContainerWithEnv(
		corev1.EnvVar{
			Name:  "INSPR_SIDECAR_READ_PORT",
			Value: string(config.ports.OutPort),
		},
		corev1.EnvVar{
			Name:  "INSPR_SIDECAR_WRITE_PORT",
			Value: string(config.ports.InPort),
		},
		corev1.EnvVar{
			Name:  "INSPR_SIDECAR_PORT",
			Value: config.kafkaInsprPort,
		},
		corev1.EnvVar{
			Name:  "INSPR_INSPRD_ADDRESS",
			Value: "http://" + config.kafkaInsprdName + "." + config.kafkaNamespace + "." + "svc:" + config.kafkaInsprPort,
		},
	)
}
