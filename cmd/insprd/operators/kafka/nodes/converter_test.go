package nodes

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	test "gitlab.inspr.dev/inspr/core/pkg/testutils"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
	kubeApp "k8s.io/api/apps/v1"
	kubeCore "k8s.io/api/core/v1"
	kubeMeta "k8s.io/apimachinery/pkg/apis/meta/v1"
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
			SHA256:    "sha256",
		},
		Spec: meta.AppSpec{
			Node: meta.Node{
				Meta: meta.Metadata{
					Name:      "mock_node",
					Reference: "ref",
					Parent:    "mock_app",
					SHA256:    "sha256",
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
	testEnv := map[string]string{
		"INSPR_INPUT_CHANNELS":  inputChannels,
		"INSPR_CHANNEL_SIDECAR": environment.GetEnvironment().SidecarImage,
		"INSPR_APPS_TLS":        "true",

		"INSPR_OUTPUT_CHANNELS": outputChannels,
		"INSPR_APP_ID":          environment.GetEnvironment().InsprAppContext + "." + testApp.Meta.Name,
	}

	appDeployName := toDeploymentName(environment.GetEnvironment().InsprEnvironment, &testApp)

	type args struct {
		app *meta.App
	}

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
					Labels: map[string]string{"app": appDeployName},
				},
				Spec: kubeApp.DeploymentSpec{
					Selector: &kubeMeta.LabelSelector{
						MatchLabels: map[string]string{
							"app": appDeployName,
						},
					},
					Template: kubeCore.PodTemplateSpec{

						ObjectMeta: kubeMeta.ObjectMeta{
							Labels: map[string]string{
								"app": appDeployName,
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
									Image: environment.GetEnvironment().SidecarImage,
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dAppToDeployment(tt.args.app); !cmp.Equal(got, tt.want, test.GetMapCompareOptions()) {
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
		},
	}
	type args struct {
		filePath string
		app      *meta.App
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "successful_need_replacement",
			args: args{
				filePath: "test",
				app:      &testApp,
			},
			want: "inspr-test-parent-app1",
		},
		{
			name: "removing_first_character_when_dot",
			args: args{
				filePath: "",
				app:      &testApp,
			},
			want: "inspr-parent-app1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toDeploymentName(tt.args.filePath, tt.args.app); got != tt.want {
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
