package nodes

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	mfake "gitlab.inspr.dev/inspr/core/cmd/insprd/memory/fake"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/tree"
	kafkasc "gitlab.inspr.dev/inspr/core/cmd/sidecars/kafka/client"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	metautils "gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
	kubeApp "k8s.io/api/apps/v1"
	kubeCore "k8s.io/api/core/v1"
	kubeMeta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestInsprDAppToK8sDeployment(t *testing.T) {
	environment.SetMockEnv()

	os.Setenv("KAFKA_BOOTSTRAP_SERVERS", "kafka.default.svc:9092")
	os.Setenv("KAFKA_AUTO_OFFSET_RESET", "earliest")
	testApp := meta.App{
		Meta: meta.Metadata{
			Name:      "mock_app",
			Reference: "ref",
			Parent:    "parent",
			UUID:      "parent-UUID",
		},
		Spec: meta.AppSpec{
			Node: meta.Node{
				Meta: meta.Metadata{
					Name:      "mock_node",
					Reference: "ref",
					Parent:    "mock_app",
					UUID:      "parent-UUID",
				},
				Spec: meta.NodeSpec{
					Image:    "nodeImage",
					Replicas: 3,
					Environment: map[string]string{
						"key_1": "value_1",
						"key_2": "value_2",
						"key_3": "value_3",
					},
				},
			},
		},
	}

	outputChannels := strings.Join(testApp.Spec.Boundary.Output, ";")
	inputChannels := strings.Join(testApp.Spec.Boundary.Input, ";")
	appID, _ := metautils.JoinScopes(environment.GetInsprAppContext(), testApp.Meta.Name)
	testEnv := map[string]string{
		"INSPR_INPUT_CHANNELS":  inputChannels,
		"INSPR_CHANNEL_SIDECAR": environment.GetSidecarImage(),
		"INSPR_APPS_TLS":        "true",
		"INSPR_OUTPUT_CHANNELS": outputChannels,
		"INSPR_APP_ID":          appID,
	}

	appDeployName := toDeploymentName(&testApp)
	appID = toAppID(&testApp)
	type args struct {
		app *meta.App
	}

	replicasHelper := int32(testApp.Spec.Node.Spec.Replicas)

	tests := []struct {
		name string
		args args
		want *kubeApp.Deployment
	}{
		{
			name: "successful",
			args: args{&testApp},
			want: &kubeApp.Deployment{
				ObjectMeta: kubeMeta.ObjectMeta{
					Name:   appDeployName,
					Labels: map[string]string{"app": appID},
				},
				Spec: kubeApp.DeploymentSpec{
					Selector: &kubeMeta.LabelSelector{
						MatchLabels: map[string]string{
							"app": appID,
						},
					},
					Replicas: &replicasHelper,
					Template: kubeCore.PodTemplateSpec{

						ObjectMeta: kubeMeta.ObjectMeta{
							Labels: map[string]string{
								"app": appID,
							},
						},
						Spec: kubeCore.PodSpec{
							Volumes: []kubeCore.Volume{
								{
									Name: appDeployName + "-volume",
									VolumeSource: kubeCore.VolumeSource{
										EmptyDir: &kubeCore.EmptyDirVolumeSource{
											Medium: kubeCore.StorageMediumMemory,
										},
									},
								},
							},
							Containers: []kubeCore.Container{
								{
									Name:  appDeployName,
									Image: testApp.Spec.Node.Spec.Image,
									// parse from master env var to kube env vars
									VolumeMounts: []kubeCore.VolumeMount{
										{
											Name:      appDeployName + "-volume",
											MountPath: "/inspr",
										},
									},
									Env: append(utils.EnvironmentMap(testApp.Spec.Node.Spec.Environment).ParseToK8sArrEnv(),
										kubeCore.EnvVar{
											Name: "INSPR_UNIX_SOCKET",
											ValueFrom: &kubeCore.EnvVarSource{
												FieldRef: &kubeCore.ObjectFieldSelector{
													FieldPath: "metadata.name",
												},
											},
										}),
								},
								{
									Name:  appDeployName + "-sidecar",
									Image: environment.GetSidecarImage(),
									VolumeMounts: []kubeCore.VolumeMount{
										{
											Name:      appDeployName + "-volume",
											MountPath: "/inspr",
										},
									},
									Env: append(utils.EnvironmentMap(testEnv).ParseToK8sArrEnv(), kubeCore.EnvVar{
										Name: "INSPR_UNIX_SOCKET",
										ValueFrom: &kubeCore.EnvVarSource{
											FieldRef: &kubeCore.ObjectFieldSelector{
												FieldPath: "metadata.name",
											},
										},
									}),
								},
							},
						},
					},
				},
			},
		},
	}
	op := &NodeOperator{
		memory:    mfake.MockMemoryManager(nil),
		clientSet: fake.NewSimpleClientset(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := op.dAppToDeployment(tt.args.app); !cmp.Equal(got, tt.want, utils.GetMapCompareOptions()) {
				t.Errorf("InsprDAppToK8sDeployment() = \n%v, \nwant \n%v", got, tt.want)
			}
		})
	}
	environment.UnsetMockEnv()
}

func Test_toDeploymentName(t *testing.T) {
	testApp := meta.App{
		Meta: meta.Metadata{
			Name:   "app1",
			Parent: "parent",
			UUID:   "APP-UUID",
		},
	}
	type args struct {
		app *meta.App
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "successful_need_replacement",
			args: args{
				app: &testApp,
			},
			want: "node-APP-UUID",
		},
		{
			name: "complex parent",
			args: args{app: &meta.App{
				Meta: meta.Metadata{
					Name:   "app1",
					Parent: "ggp.gp.p",
					UUID:   "APP-UUID",
				},
			}},
			want: "node-APP-UUID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toDeploymentName(tt.args.app); got != tt.want {
				t.Errorf("toDeploymentName() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_intToint32(t *testing.T) {
	var x int32 = 3
	var pointer32int *int32 = &x
	type args struct {
		v int
	}
	tests := []struct {
		name string
		args args
		want *int32
	}{
		{
			name: "successful",
			args: args{v: 3},
			want: pointer32int,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := intToint32(tt.args.v)

			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("intToint32().Type() = %v, want %v", reflect.TypeOf(got), reflect.TypeOf(tt.want))
			}
			if *got != *tt.want {
				t.Errorf("intToint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseEnvironment(t *testing.T) {
	os.Setenv("KAFKA_BOOTSTRAP_SERVERS", "bootstrap")
	os.Setenv("KAFKA_AUTO_OFFSET_RESET", "earliest")
	os.Setenv("INSPR_SIDECAR_IMAGE", "sidecar")
	os.Setenv("INSPR_ENV", "environment")
	defer os.Unsetenv("KAFKA_BOOTSTRAP_SERVERS")
	defer os.Unsetenv("KAFKA_AUTO_OFFSET_RESET")
	defer os.Unsetenv("INSPR_SIDECAR_IMAGE")
	defer os.Unsetenv("INSPR_ENV")
	kafkasc.RefreshEnviromentVariables()
	mem := tree.GetTreeMemory()
	chs := map[string]*meta.Channel{
		"channel1": {
			Meta: meta.Metadata{
				Name:   "channel1",
				Parent: "parent",
				UUID:   "uuid-channel1",
			},
			Spec: meta.ChannelSpec{
				Type: "channelType1",
			},
		},
		"channel2": {
			Meta: meta.Metadata{
				Name:   "channel2",
				Parent: "parent",
				UUID:   "uuid-channel2",
			},
			Spec: meta.ChannelSpec{
				Type: "channelType2",
			},
		},
	}
	mem.InitTransaction()
	mem.Apps().Create("", &meta.App{
		Meta: meta.Metadata{
			Name: "parent",
		},
		Spec: meta.AppSpec{
			Channels: chs,
			Aliases:  map[string]*meta.Alias{},
			ChannelTypes: map[string]*meta.ChannelType{
				"channelType1": {
					Meta: meta.Metadata{
						Name:   "channelType1",
						Parent: "parent",
					},
					Schema: "schema1",
				},
				"channelType2": {
					Meta: meta.Metadata{
						Name:   "channelType2",
						Parent: "parent",
					},
					Schema: "schema2",
				},
			},
		},
	})
	mem.Commit()
	type args struct {
		app *meta.App
	}
	type fields struct {
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   utils.EnvironmentMap
	}{
		{
			name: "no boundary",
			args: args{
				app: &meta.App{
					Meta: meta.Metadata{
						Name:   "app1",
						Parent: "parent",
					},
				},
			},
			want: utils.EnvironmentMap{
				"INSPR_INPUT_CHANNELS":    "",
				"INSPR_OUTPUT_CHANNELS":   "",
				"INSPR_SIDECAR_IMAGE":     "sidecar",
				"INSPR_APP_ID":            "parent-app1",
				"INSPR_APP_CTX":           "parent",
				"INSPR_ENV":               "environment",
				"KAFKA_BOOTSTRAP_SERVERS": "bootstrap",
				"KAFKA_AUTO_OFFSET_RESET": "earliest",
			},
		},
		{
			name: "only input boundary",
			args: args{
				app: &meta.App{
					Meta: meta.Metadata{
						Name:   "app1",
						Parent: "parent",
					},
					Spec: meta.AppSpec{
						Boundary: meta.AppBoundary{
							Input: []string{
								"channel1",
								"channel2",
							},
						},
					},
				},
			},
			want: utils.EnvironmentMap{
				"INSPR_INPUT_CHANNELS":    "channel1;channel2",
				"INSPR_OUTPUT_CHANNELS":   "",
				"INSPR_SIDECAR_IMAGE":     "sidecar",
				"INSPR_APP_ID":            "parent-app1",
				"INSPR_APP_CTX":           "parent",
				"INSPR_ENV":               "environment",
				"KAFKA_BOOTSTRAP_SERVERS": "bootstrap",
				"KAFKA_AUTO_OFFSET_RESET": "earliest",
				"INSPR_" + chs["channel1"].Meta.UUID + "_SCHEMA": "schema1",
				"INSPR_" + chs["channel2"].Meta.UUID + "_SCHEMA": "schema2",
				"channel1_RESOLVED": "INSPR_" + chs["channel1"].Meta.UUID,
				"channel2_RESOLVED": "INSPR_" + chs["channel2"].Meta.UUID,
			},
		},
	}

	for _, tt := range tests {
		op := &NodeOperator{
			memory:    mem,
			clientSet: fake.NewSimpleClientset(),
		}
		t.Run(tt.name, func(t *testing.T) {
			mem.InitTransaction()
			if got := op.baseEnvironment(tt.args.app); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("baseEnvironment() = \n%v, want \n%v", got, tt.want)
			}
			mem.Cancel()
		})
	}
}
