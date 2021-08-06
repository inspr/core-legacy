package sidecars

import (
	"os"
	"strconv"

	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/brokers"
	"inspr.dev/inspr/pkg/operator/k8s"
	"inspr.dev/inspr/pkg/sidecars/models"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
)

// KafkaConfig configurations used to create the KafkaSidecar
type KafkaConfig struct {
	BootstrapServers string `yaml:"bootstrapServers"`
	AutoOffsetReset  string `yaml:"autoOffsetReset"`
	SidecarImage     string `yaml:"sidecarImage"`
	// KafkaInsprAddr is the port used in the insprd service of your cluster
	KafkaInsprAddr string `yaml:"sidecarAddr"`
}

//Broker is a BrokerConfiguration interface method, it returns the broker name for this config type
func (kc KafkaConfig) Broker() string {
	return brokers.Kafka
}

// KafkaToDeployment receives a the KafkaConfig variable as a parameter and returns a
// SidecarFactory function that is used to subscribe to the sidecarFactory
func KafkaToDeployment(config KafkaConfig) models.SidecarFactory {
	// Handles defaults values in case any of the kafkaConfig variables are empty
	if config.KafkaInsprAddr == "" {
		config.KafkaInsprAddr = "http://localhost"
	}

	os.Setenv("INSPR_SIDECAR_KAFKA_BOOTSTRAP_SERVERS", config.BootstrapServers)
	return func(app *meta.App, conn *models.SidecarConnections, opts ...k8s.ContainerOption) (corev1.Container, []corev1.EnvVar) {
		envVars, kafkAddr := KafkaSidecarConfig(config, conn)

		return k8s.NewContainer(
			"sidecar-kafka-"+app.Meta.UUID, // deployment name
			config.SidecarImage,            // image url
			getKafkaContainerOptions(app, config, envVars, opts)...,
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

func getKafkaContainerOptions(app *meta.App, config KafkaConfig, envVars k8s.ContainerOption,
	opts []k8s.ContainerOption) []k8s.ContainerOption {

	stdOptions := []k8s.ContainerOption{
		InsprAppIDConfig(app),
		KafkaEnvConfig(config),
		envVars,
		k8s.ContainerWithPullPolicy(corev1.PullAlways),
		k8s.ContainerWithPorts(v1.ContainerPort{
			Name:          "tcp-kfk-metrics",
			ContainerPort: 16001,
		}),
	}

	return append(stdOptions, opts...)
}
