package sidecars

import (
	"strconv"

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
	KafkaInsprAddr string
}

// KafkaToDeployment receives a the KafkaConfig variable as a parameter and returns a
// SidecarFactory function that is used to subscribe to the sidecarFactory
func KafkaToDeployment(config KafkaConfig) models.SidecarFactory {
	return func(app *meta.App, conn *models.SidecarConnections, opts ...k8s.ContainerOption) (corev1.Container, []corev1.EnvVar) {
		envVars, kafkAddr := KafkaSidecarConfig(config, conn)

		return k8s.NewContainer(
			"sidecar-kafka-"+app.Meta.UUID, // deployment name
			config.SidecarImage,            // image url

			// label to the dApp associated with it
			InsprAppIDConfig(app),
			KafkaEnvConfig(config),
			envVars,
			k8s.ContainerWithPullPolicy(corev1.PullAlways),
			opts[0],
			opts[1],
		), kafkAddr
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
func KafkaSidecarConfig(config KafkaConfig, conns *models.SidecarConnections) (k8s.ContainerOption, []corev1.EnvVar) {
	port := corev1.EnvVar{
		Name:  "INSPR_SIDECAR_KAFKA_WRITE_PORT",
		Value: strconv.Itoa(int(conns.InPort)),
	}

	kafkaAddr := corev1.EnvVar{
		Name:  "INSPR_SIDECAR_KAFKA_ADDR",
		Value: config.KafkaInsprAddr,
	}

	return k8s.ContainerWithEnv(
			port,
		),
		[]corev1.EnvVar{port, kafkaAddr}
}
