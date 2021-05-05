package nodes

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	mfake "github.com/inspr/inspr/cmd/insprd/memory/fake"
	"github.com/inspr/inspr/pkg/environment"
	"github.com/inspr/inspr/pkg/meta"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
	"github.com/inspr/inspr/pkg/utils"
	kubeApp "k8s.io/api/apps/v1"
	kubeCore "k8s.io/api/core/v1"
	kubeMeta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/fake"
)

func TestInsprDAppToK8sDeployment(t *testing.T) {
	environment.SetMockEnv()
	defer environment.UnsetMockEnv()

	os.Setenv("KAFKA_BOOTSTRAP_SERVERS", "kafka.default.svc:9092")
	os.Setenv("KAFKA_AUTO_OFFSET_RESET", "earliest")
	defer os.Unsetenv("KAFKA_BOOTSTRAP_SERVERS")
	defer os.Unsetenv("KAFKA_AUTO_OFFSET_RESET")
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

func Test_dappToService(t *testing.T) {
	m := metautils.InjectUUID(meta.Metadata{
		Name:   "test-dapp",
		Parent: "",
	})
	type args struct {
		app *meta.App
	}
	tests := []struct {
		name string
		args args
		want *kubeCore.Service
	}{
		{
			name: "node with ports",
			args: args{
				app: &meta.App{
					Meta: m,
					Spec: meta.AppSpec{
						Node: meta.Node{
							Spec: meta.NodeSpec{
								Ports: []meta.NodePort{
									{
										Port:       80,
										TargetPort: 80,
									},
									{
										Port:       90,
										TargetPort: 100,
									},
								},
							},
						},
					},
				},
			},
			want: &kubeCore.Service{
				ObjectMeta: kubeMeta.ObjectMeta{
					Name: "node-" + m.UUID,
				},
				Spec: kubeCore.ServiceSpec{
					Ports: []kubeCore.ServicePort{
						{
							Port:       80,
							TargetPort: intstr.FromInt(80),
						},
						{
							Port:       90,
							TargetPort: intstr.FromInt(100),
						},
						{
							Port:       sidecarPort,
							TargetPort: intstr.FromInt(int(sidecarPort)),
						},
					},
					Selector: map[string]string{
						"app": "test-dapp",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dappToService(tt.args.app); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dappToService() = %v, want %v", got, tt.want)
			}
		})
	}
}
