package sidecars

import (
	"reflect"
	"testing"

	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/operator/k8s"
	"github.com/inspr/inspr/pkg/sidecar/models"
	corev1 "k8s.io/api/core/v1"
)

// constants used for the tests
const (
	testBootstrap       = "bootstrap"
	testAutoOff         = "autooff"
	testSidecarImage    = "image"
	testKafkaInsprdName = "insprdName"
	testKafkaNamespace  = "insprdNamespace"
	testKafkaInsprPort  = "insprdPort"
)

var testPorts = models.SidecarConnections{
	InPort:  00,
	OutPort: 01,
}

func TestKafkaToDeployment(t *testing.T) {
	deploymentKafkaConfig := KafkaConfig{
		bootstrapServers: testBootstrap,
		autoOffsetReset:  testAutoOff,
		sidecarImage:     testSidecarImage,
		kafkaInsprdName:  testKafkaInsprdName,
		kafkaInsprPort:   testKafkaInsprPort,
		kafkaNamespace:   testKafkaNamespace,
		ports:            testPorts,
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
					bootstrapServers: testBootstrap,
					autoOffsetReset:  testAutoOff,
					sidecarImage:     testSidecarImage,
					kafkaInsprdName:  testKafkaInsprdName,
					kafkaInsprPort:   testKafkaInsprPort,
					kafkaNamespace:   testKafkaNamespace,
					ports:            testPorts,
				},
				dapp: meta.App{},
			},
			want: k8s.WithContainer(
				k8s.NewContainer(
					"sidecar-kafka-dappUUID",
					testSidecarImage,
					insprAppIDConfig(&deploymentDApp),
					kafkaConfig(deploymentKafkaConfig),
					sidecarConfig(deploymentKafkaConfig),
					k8s.ContainerWithPullPolicy(corev1.PullAlways),
				),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KafkaToDeployment(tt.args.config)

			gotDepOption := got(&deploymentDApp, testPorts)

			gotDeploy := k8s.NewDeployment("", gotDepOption)
			wantDeploy := k8s.NewDeployment("", tt.want)

			if !reflect.DeepEqual(gotDeploy, wantDeploy) {
				t.Errorf("KafkaToDeployment() = %v, want %v",
					gotDeploy, wantDeploy)
			}
		})
	}
}

func Test_kafkaConfig(t *testing.T) {
	type args struct {
		config KafkaConfig
	}
	tests := []struct {
		name string
		args args
		want k8s.ContainerOption
	}{
		{
			name: "kafkaConfig_base_testing",
			args: args{config: KafkaConfig{
				bootstrapServers: testBootstrap,
				autoOffsetReset:  testAutoOff,
				sidecarImage:     testSidecarImage,
			}},
			want: k8s.ContainerWithEnv(
				corev1.EnvVar{
					Name:  "KAFKA_BOOTSTRAP_SERVERS",
					Value: testBootstrap,
				},
				corev1.EnvVar{
					Name:  "KAFKA_AUTO_OFFSET_RESET",
					Value: testAutoOff,
				},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := kafkaConfig(tt.args.config)

			gotContainer := k8s.NewContainer("", "", got)
			wantContainer := k8s.NewContainer("", "", tt.want)

			if !reflect.DeepEqual(gotContainer, wantContainer) {
				t.Errorf("kafkaConfig() = %v, want %v",
					gotContainer,
					wantContainer)
			}
		})
	}
}

func Test_sidecarConfig(t *testing.T) {
	type args struct {
		config KafkaConfig
	}
	tests := []struct {
		name string
		args args
		want k8s.ContainerOption
	}{
		{
			name: "kafkaConfig_base_testing",
			args: args{config: KafkaConfig{
				bootstrapServers: testBootstrap,
				autoOffsetReset:  testAutoOff,
				sidecarImage:     testSidecarImage,
				kafkaInsprdName:  testKafkaInsprdName,
				kafkaInsprPort:   testKafkaInsprPort,
				kafkaNamespace:   testKafkaNamespace,
				ports:            testPorts,
			}},
			want: k8s.ContainerWithEnv(
				corev1.EnvVar{
					Name:  "INSPR_SIDECAR_READ_PORT",
					Value: string(testPorts.OutPort),
				},
				corev1.EnvVar{
					Name:  "INSPR_SIDECAR_WRITE_PORT",
					Value: string(testPorts.InPort),
				},
				corev1.EnvVar{
					Name:  "INSPR_SIDECAR_PORT",
					Value: testKafkaInsprPort,
				},
				corev1.EnvVar{
					Name:  "INSPR_INSPRD_ADDRESS",
					Value: "http://" + testKafkaInsprdName + "." + testKafkaNamespace + "." + "svc:" + testKafkaInsprPort,
				},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sidecarConfig(tt.args.config)

			gotContainer := k8s.NewContainer("", "", got)
			wantContainer := k8s.NewContainer("", "", tt.want)

			if !reflect.DeepEqual(gotContainer, wantContainer) {
				t.Errorf("sidecarConfig() = %v, want %v",
					gotContainer, wantContainer)
			}
		})
	}
}
