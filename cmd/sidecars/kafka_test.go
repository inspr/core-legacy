package sidecars

import (
	"reflect"
	"testing"

	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/operator/k8s"
	"inspr.dev/inspr/pkg/sidecars/models"
	corev1 "k8s.io/api/core/v1"
)

// constants used for the tests
const (
	testBootstrap      = "bootstrap"
	testAutoOff        = "autooff"
	testSidecarImage   = "image"
	testKafkaInsprAddr = "insprdPort"
)

var testPorts = models.SidecarConnections{
	InPort:  00,
	OutPort: 01,
}

func extractContainerOpts(
	opts k8s.ContainerOption,
	envvars []corev1.EnvVar,
) k8s.ContainerOption {
	return opts
}

// This test covers all kafka.go methods
func TestKafkaToDeployment(t *testing.T) {
	deploymentKafkaConfig := KafkaConfig{
		BootstrapServers: testBootstrap,
		AutoOffsetReset:  testAutoOff,
		SidecarImage:     testSidecarImage,
		KafkaInsprAddr:   testKafkaInsprAddr,
	}
	deploymentDApp := meta.App{
		Meta: meta.Metadata{
			Name:   "dapp",
			Parent: "dapp1.dapp2",
			UUID:   "dappUUID",
		},
	}

	type args struct {
		config KafkaConfig
		dapp   meta.App
	}
	tests := []struct {
		name string
		args args
		want k8s.DeploymentOption
	}{
		{
			name: "kafkaToDeployment_base_test",
			args: args{
				config: KafkaConfig{
					BootstrapServers: testBootstrap,
					AutoOffsetReset:  testAutoOff,
					SidecarImage:     testSidecarImage,
					KafkaInsprAddr:   testKafkaInsprAddr,
				},
				dapp: meta.App{},
			},
			want: k8s.WithContainer(
				k8s.NewContainer(
					"sidecar-kafka-dappUUID",
					testSidecarImage,
					InsprAppIDConfig(&deploymentDApp),
					KafkaEnvConfig(deploymentKafkaConfig),
					extractContainerOpts(
						KafkaSidecarConfig(deploymentKafkaConfig, &testPorts),
					),
					k8s.ContainerWithPullPolicy(corev1.PullAlways),
				),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KafkaToDeployment(tt.args.config)

			gotDepOption, _ := got(&deploymentDApp, &testPorts, nil)

			gotDeploy := k8s.NewDeployment("", k8s.WithContainer(gotDepOption))
			wantDeploy := k8s.NewDeployment("", tt.want)

			if !reflect.DeepEqual(gotDeploy, wantDeploy) {
				t.Errorf("KafkaToDeployment() = %v, want %v",
					gotDeploy, wantDeploy)
			}
		})
	}
}
