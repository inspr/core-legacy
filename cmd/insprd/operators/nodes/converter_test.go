package nodes

import (
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"inspr.dev/inspr/cmd/insprd/memory/fake"
	"inspr.dev/inspr/cmd/insprd/memory/tree"
	apimodels "inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/auth"
	authmock "inspr.dev/inspr/pkg/auth/mocks"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/utils"
	kubeCore "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kfake "k8s.io/client-go/kubernetes/fake"
)

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

func TestNodeOperator_withBoundary(t *testing.T) { //TODO needs more testing, longer creates, updates, and deletes
	mem := fake.MockTreeMemory(nil)
	mem.InitTransaction()
	mem.Channels().Create("", &meta.Channel{
		Meta: meta.Metadata{
			Name: "channel1",
			UUID: "channel1_UUID",
		},
		Spec: meta.ChannelSpec{
			Type:           "channel1type",
			SelectedBroker: "someBroker",
		},
	}, nil)
	mem.Channels().Create("", &meta.Channel{
		Meta: meta.Metadata{
			Name: "channel2",
			UUID: "channel2_UUID",
		},
		Spec: meta.ChannelSpec{
			Type:           "channel2type",
			SelectedBroker: "someBroker",
		},
	}, nil)

	mem.Types().Create("", &meta.Type{
		Meta: meta.Metadata{
			Name: "channel1type",
		},
		Schema: "channel1type",
	})
	mem.Types().Create("", &meta.Type{
		Meta: meta.Metadata{
			Name: "channel2type",
		},
		Schema: "channel2type",
	})

	type fields struct {
		clientSet kubernetes.Interface
		memory    tree.Manager
		auth      auth.Auth
	}
	type args struct {
		app *meta.App
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *kubeCore.Container
	}{
		{
			name: "no input channels",
			fields: fields{
				clientSet: kfake.NewSimpleClientset(),
				memory:    mem,
			},
			args: args{
				app: &meta.App{
					Meta: meta.Metadata{
						Name: "app1",
					},
					Spec: meta.AppSpec{
						Boundary: meta.AppBoundary{
							Channels: meta.Boundary{
								Input: []string{
									"channel1",
									"channel2",
								},
							},
						},
						Channels: map[string]*meta.Channel{
							"channel1": {
								Meta: meta.Metadata{Name: "channel1"},
								Spec: meta.ChannelSpec{SelectedBroker: "*"},
							},
							"channel2": {
								Meta: meta.Metadata{Name: "channel2"},
								Spec: meta.ChannelSpec{SelectedBroker: "someBroker"},
							},
						},
					},
				},
			},
			want: &kubeCore.Container{
				Env: []kubeCore.EnvVar{
					{
						Name:  "INSPR_INPUT_CHANNELS",
						Value: "channel1@someBroker;channel2@someBroker",
					},
					{
						Name:  "INSPR_OUTPUT_CHANNELS",
						Value: "",
					},
					{
						Name:  "INSPR_channel1_UUID_SCHEMA",
						Value: "channel1type",
					},
					{
						Name:  "channel1_RESOLVED",
						Value: "INSPR_channel1_UUID",
					},
					{
						Name:  "INSPR_channel2_UUID_SCHEMA",
						Value: "channel2type",
					},
					{
						Name:  "channel2_RESOLVED",
						Value: "INSPR_channel2_UUID",
					},
				},
			},
		},
		{
			name: "no output channels",
			fields: fields{
				clientSet: kfake.NewSimpleClientset(),
				memory:    mem,
			},
			args: args{
				app: &meta.App{
					Meta: meta.Metadata{
						Name: "app2",
					},
					Spec: meta.AppSpec{
						Boundary: meta.AppBoundary{
							Channels: meta.Boundary{
								Output: []string{
									"channel1",
									"channel2",
								},
							},
						},
						Channels: map[string]*meta.Channel{
							"channel1": {
								Meta: meta.Metadata{Name: "channel1"},
								Spec: meta.ChannelSpec{SelectedBroker: "someBroker"},
							},
							"channel2": {
								Meta: meta.Metadata{Name: "channel2"},
								Spec: meta.ChannelSpec{SelectedBroker: "someBroker"},
							},
						},
					},
				},
			},
			want: &kubeCore.Container{
				Env: []kubeCore.EnvVar{
					{
						Name:  "INSPR_OUTPUT_CHANNELS",
						Value: "channel1@someBroker;channel2@someBroker",
					},
					{
						Name:  "INSPR_INPUT_CHANNELS",
						Value: "",
					},
					{
						Name:  "INSPR_channel1_UUID_SCHEMA",
						Value: "channel1type",
					},
					{
						Name:  "channel1_RESOLVED",
						Value: "INSPR_channel1_UUID",
					},
					{
						Name:  "INSPR_channel2_UUID_SCHEMA",
						Value: "channel2type",
					},
					{
						Name:  "channel2_RESOLVED",
						Value: "INSPR_channel2_UUID",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem.Apps().Create("", tt.args.app, &apimodels.BrokersDI{
				Available: []string{"some_broker"},
				Default:   "some_broker",
			})
			no := &NodeOperator{
				clientSet: tt.fields.clientSet,
				memory:    tt.fields.memory,
				auth:      tt.fields.auth,
			}
			got := &kubeCore.Container{}
			option := no.withBoundary(tt.args.app, false)
			option(got)
			if !cmp.Equal(got, tt.want, cmp.Comparer(func(a1, a2 []kubeCore.EnvVar) bool {
				a1cmp, a2cmp := envVarArr(a1), envVarArr(a2)
				sort.Sort(a1cmp)
				sort.Sort(a2cmp)

				return cmp.Equal(a1cmp, a2cmp)
			})) {
				t.Errorf("TestNodeOperator_withBoundary got = \n%v, \nwant \n%v", got, tt.want)
			}

		})
	}
}

type envVarArr []kubeCore.EnvVar

func (a envVarArr) Len() int {
	return len(a)
}
func (a envVarArr) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a envVarArr) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

func Test_withNodeID(t *testing.T) {
	type args struct {
		app *meta.App
	}
	tests := []struct {
		name string
		args args
		want *kubeCore.Container
	}{
		{
			name: "correct injection",
			args: args{
				app: &meta.App{
					Meta: meta.Metadata{
						Name: "app1",
						UUID: "UUID",
					},
				},
			},
			want: &kubeCore.Container{
				Env: []kubeCore.EnvVar{
					{
						Name:  "INSPR_APP_ID",
						Value: "app1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := withNodeID(tt.args.app)
			got := &kubeCore.Container{}
			option(got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("withNodeID() got = %v, want = %v", got, tt.want)
			}

		})
	}
}

func Test_withLBSidecarPorts(t *testing.T) {
	type args struct {
		app *meta.App
	}
	tests := []struct {
		name string
		args args
		want *kubeCore.Container
	}{
		{
			name: "only write port",
			args: args{
				app: &meta.App{
					Spec: meta.AppSpec{
						Node: meta.Node{
							Spec: meta.NodeSpec{
								SidecarPort: meta.SidecarPort{
									LBWrite: 1234,
								},
							},
						},
					},
				},
			},
			want: &kubeCore.Container{
				Env: []kubeCore.EnvVar{
					{
						Name:  "INSPR_LBSIDECAR_WRITE_PORT",
						Value: "1234",
					},
				},
			},
		},
		{
			name: "only read port",
			args: args{
				app: &meta.App{
					Spec: meta.AppSpec{
						Node: meta.Node{
							Spec: meta.NodeSpec{
								SidecarPort: meta.SidecarPort{
									LBRead: 1234,
								},
							},
						},
					},
				},
			},
			want: &kubeCore.Container{
				Env: []kubeCore.EnvVar{
					{
						Name:  "INSPR_LBSIDECAR_READ_PORT",
						Value: "1234",
					},
				},
			},
		},
		{
			name: "both ports",
			args: args{
				app: &meta.App{
					Spec: meta.AppSpec{
						Node: meta.Node{
							Spec: meta.NodeSpec{
								SidecarPort: meta.SidecarPort{
									LBRead:  1234,
									LBWrite: 1234,
								},
							},
						},
					},
				},
			},
			want: &kubeCore.Container{
				Env: []kubeCore.EnvVar{
					{
						Name:  "INSPR_LBSIDECAR_WRITE_PORT",
						Value: "1234",
					},
					{
						Name:  "INSPR_LBSIDECAR_READ_PORT",
						Value: "1234",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			option := overwritePortEnvs(tt.args.app)
			got := &kubeCore.Container{}
			option(got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("withLBSidecarPorts() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestNodeOperator_withLBSidecarImage(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name string
		env  string
		args args
		want *kubeCore.Container
	}{
		{
			name: "correct sidecar image",
			want: &kubeCore.Container{
				Image: "inspr-sidecar-image",
			},
			env: "inspr-sidecar-image",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("INSPR_LBSIDECAR_IMAGE", tt.env)
			defer os.Unsetenv("INSPR_LBSIDECAR_IMAGE")
			no := &NodeOperator{}
			option := no.withLBSidecarImage()
			got := &kubeCore.Container{}
			option(got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("withLBSidecarPorts() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func Test_withLBSidecarConfiguration(t *testing.T) {
	tests := []struct {
		name   string
		want   *kubeCore.Container
		before func()
		after  func()
	}{
		{
			before: func() {
				os.Setenv("INSPR_LBSIDECAR_CONFIGMAP", "inspr-lbsidecar-configuration")
			},
			name: "correct configmap configuration",
			want: &kubeCore.Container{
				EnvFrom: []kubeCore.EnvFromSource{
					{
						ConfigMapRef: &kubeCore.ConfigMapEnvSource{
							LocalObjectReference: kubeCore.LocalObjectReference{
								Name: "inspr-lbsidecar-configuration",
							},
						},
					},
				},
			},
			after: func() {
				os.Unsetenv("INSPR_LBSIDECAR_CONFIGMAP_NAME")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()
			defer tt.after()
			option := withLBSidecarConfiguration()
			got := &kubeCore.Container{}
			option(got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("withLBSidecarConfiguration() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestNodeOperator_toSecret(t *testing.T) {

	type fields struct {
		clientSet kubernetes.Interface
		memory    tree.Manager
		auth      auth.Auth
	}
	type args struct {
		app *meta.App
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *kubeSecret
	}{
		{
			name: "correct secret definition",
			fields: fields{
				clientSet: kfake.NewSimpleClientset(),
				auth:      authmock.NewMockAuth(nil),
			},
			args: args{
				app: &meta.App{
					Meta: meta.Metadata{
						Name: "app1",
						UUID: "app1_UUID",
					},
					Spec: meta.AppSpec{

						Auth: meta.AppAuth{
							Scope: "scope1",
							Permissions: utils.StringArray{
								"create:dapp",
								"delete:dapp",
							},
						},
					},
				},
			},
			want: &kubeSecret{
				ObjectMeta: v1.ObjectMeta{
					Name: "node-app1_UUID",
				},
				Data: map[string][]byte{
					"INSPR_CONTROLLER_TOKEN": []byte("mock"),
					"INSPR_CONTROLLER_SCOPE": []byte("scope1"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			no := &NodeOperator{
				clientSet: tt.fields.clientSet,
				memory:    tt.fields.memory,
				auth:      tt.fields.auth,
			}
			if got := no.toSecret(tt.args.app); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeOperator.toSecret() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_WithNodePort(t *testing.T) {
	tests := []struct {
		name string
		want *kubeCore.Container
	}{
		{
			name: "injection",
			want: &kubeCore.Container{
				Ports: []kubeCore.ContainerPort{
					{
						Name:          "tcp-nd-metrics",
						ContainerPort: 16002,
						Protocol:      kubeCore.ProtocolTCP,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := withNodePort()
			got := &kubeCore.Container{}
			options(got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("withNodePort() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func Test_WithLBPort(t *testing.T) {
	tests := []struct {
		name string
		want *kubeCore.Container
	}{
		{
			name: "injection",
			want: &kubeCore.Container{
				Ports: []kubeCore.ContainerPort{
					{
						Name:          "tcp-lbs-metrics",
						ContainerPort: 16000,
						Protocol:      kubeCore.ProtocolTCP,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := withLBPort()
			got := &kubeCore.Container{}
			options(got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("withLBPort() got = %v, want = %v", got, tt.want)
			}
		})
	}
}
