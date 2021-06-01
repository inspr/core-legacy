package sidecars

import (
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/operator/k8s"
	"github.com/inspr/inspr/pkg/sidecars/models"
	corev1 "k8s.io/api/core/v1"
)

// KafkaConfig configurations used to create the KafkaSidecar
type KafkaConfig struct {
	BootstrapServers string
	AutoOffsetReset  string
	SidecarImage     string
	// insprdPort is the port used in the insprd service of your cluster
	KafkaInsprPort string
}

// KafkaToDeployment receives a the KafkaConfig variable as a parameter and returns a
// SidecarFactory function that is used to subscribe to the sidecarFactory
func KafkaToDeployment(config KafkaConfig) models.SidecarFactory {
	return func(app *meta.App, conn *models.SidecarConnections) corev1.Container {
		return k8s.NewContainer(
			"sidecar-kafka-"+app.Meta.UUID, // deployment name
			config.SidecarImage,            // image url

			// label to the dApp associated with it
			InsprAppIDConfig(app),
			KafkaEnvConfig(config),
			KafkaSidecarConfig(config, conn),
			k8s.ContainerWithPullPolicy(corev1.PullAlways),
		)
	}
}

// KafkaEnvConfig adds the necessary env variables to configure kafka
func KafkaEnvConfig(config KafkaConfig) k8s.ContainerOption {
	return k8s.ContainerWithEnv(
		corev1.EnvVar{
			Name:  "INSPR_SIDECAR_KAFKA_BOOTSTRAP_SERVERS",
			Value: config.BootstrapServers,
		},
		corev1.EnvVar{
			Name:  "INSPR_SIDECAR_KAFKA_AUTO_OFFSET_RESET",
			Value: config.AutoOffsetReset,
		},
	)
}

// KafkaSidecarConfig adds the necessary env variables to configure the sidecar in the cluster
func KafkaSidecarConfig(config KafkaConfig, conns *models.SidecarConnections) k8s.ContainerOption {
	return k8s.ContainerWithEnv(
		corev1.EnvVar{
			Name:  "INSPR_SIDECAR_KAFKA_READ_PORT",
			Value: string(conns.OutPort),
		},
		corev1.EnvVar{
			Name:  "INSPR_SIDECAR_KAFKA_WRITE_PORT",
			Value: string(conns.InPort),
		},
		corev1.EnvVar{
			Name:  "INSPR_SIDECAR_KAFKA_PORT",
			Value: config.KafkaInsprPort,
		},
	)
}
