package operator

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	test "gitlab.inspr.dev/inspr/core/pkg/testutils"
	kubeApp "k8s.io/api/apps/v1"
	kubeCore "k8s.io/api/core/v1"
	kubeMeta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestInsprDAppToK8sDeployment(t *testing.T) {
	environment.SetMockEnv()
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

	inputChannels := ""
	for _, c := range testApp.Spec.Boundary.Input {
		inputChannels += c + ";"
	}
	outputChannels := ""
	for _, c := range testApp.Spec.Boundary.Output {
		outputChannels += c + ";"
	}
	testEnv := map[string]string{
		"INSPR_INPUT_CHANNELS":  inputChannels,
		"INSPR_CHANNEL_SIDECAR": environment.GetEnvironment().SidecarImage,
		"INSPR_APPS_TLS":        "true",

		"INSPR_OUTPUT_CHANNELS": outputChannels,
		"INSPR_app_ID":          environment.GetEnvironment().InsprAppContext + "." + testApp.Meta.Name,
	}

	appDeployName := toDeploymentName(environment.GetEnvironment().InsprAppContext, &testApp)

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
					Replicas: intToint32(testApp.Spec.Node.Spec.Replicas),
					Selector: &kubeMeta.LabelSelector{
						MatchLabels: map[string]string{
							"app": appDeployName,
						},
					},
					Strategy: kubeApp.DeploymentStrategy{
						Type: kubeApp.RollingUpdateDeploymentStrategyType,
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
									Name: testApp.Spec.Node.Meta.Name,
									Ports: func() []kubeCore.ContainerPort {
										return nil
									}(),

									Image: testApp.Spec.Node.Spec.Image,
									// parse from master env var to kube env vars
									ImagePullPolicy: kubeCore.PullAlways,
									VolumeMounts: []kubeCore.VolumeMount{
										{
											Name:      testApp.Spec.Node.Meta.Name + "-volume",
											MountPath: "/inspr",
										},
									},
									Env: append([]kubeCore.EnvVar{
										{
											Name: "UUID",
											ValueFrom: &kubeCore.EnvVarSource{
												FieldRef: &kubeCore.ObjectFieldSelector{
													FieldPath: "metadata.name",
												},
											},
										},
									}, parseToK8sArrEnv(testApp.Spec.Node.Spec.Environment)...), // TODO WHAT OT PUT HERE
								},
								{
									Name:            appDeployName + "-sidecar",
									Image:           environment.GetEnvironment().SidecarImage,
									ImagePullPolicy: kubeCore.PullIfNotPresent,
									VolumeMounts: []kubeCore.VolumeMount{
										{
											Name:      testApp.Spec.Node.Meta.Name + "-sidecar-volume",
											MountPath: "/inspr",
										},
									},
									Env: append(parseToK8sArrEnv(testEnv), kubeCore.EnvVar{
										Name: "UUID",
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
			if got := InsprDAppToK8sDeployment(tt.args.app); !cmp.Equal(got, tt.want, test.GetMapCompareOptions()) {
				t.Errorf("InsprDAppToK8sDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
	environment.UnsetMockEnv()
}

func Test_parseToK8sArrEnv(t *testing.T) {
	type args struct {
		arrappEnv map[string]string
	}
	tests := []struct {
		name string
		args args
		want []kubeCore.EnvVar
	}{
		{
			name: "successful_test",
			args: args{
				arrappEnv: map[string]string{
					"key_1": "value_1",
					"key_2": "value_2",
					"key_3": "value_3",
				},
			},
			want: []kubeCore.EnvVar{
				{
					Name:  "key_1",
					Value: "value_1",
				},
				{
					Name:  "key_2",
					Value: "value_2",
				},
				{
					Name:  "key_3",
					Value: "value_3",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseToK8sArrEnv(tt.args.arrappEnv)

			if !cmp.Equal(got, tt.want, test.GetMapCompareOptions()) {
				t.Errorf("parseToK8sArrEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toDeploymentName(t *testing.T) {
	testApp := meta.App{
		Meta: meta.Metadata{
			Name: "Mock-app",
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
				filePath: "root.app1.app2",
				app:      &testApp,
			},
			want: "root.app1.app2.mock-app",
		},
		{
			name: "removing_first_character_when_dot",
			args: args{
				filePath: ".root.app1.app2",
				app:      &testApp,
			},
			want: "root.app1.app2.mock-app",
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
