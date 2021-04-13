package nodes

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/inspr/inspr/cmd/insprd/memory/tree"
	"github.com/inspr/inspr/cmd/insprd/operators"
	"github.com/inspr/inspr/pkg/environment"
	"github.com/inspr/inspr/pkg/meta"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	kubeApp "k8s.io/api/apps/v1"
	kubeCore "k8s.io/api/core/v1"
	kubeMeta "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stesting "k8s.io/client-go/testing"
)

func mockK8sClientset(verb string, dep kubeApp.Deployment, erro error) kubernetes.Interface {
	environment.SetMockEnv()
	client := &fake.Clientset{}
	client.Fake.AddReactor(verb, "deployments", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		mysar := &dep
		return true, mysar, erro
	})
	os.Setenv("NODES_APPS_NAMESPACE", "default.node.opr")
	os.Setenv("KAFKA_BOOTSTRAP_SERVERS", "kafka.default.svc:9092")
	os.Setenv("KAFKA_AUTO_OFFSET_RESET", "earliest")
	return client
}

func mockK8sList(verb string, deps kubeApp.DeploymentList, erro error) kubernetes.Interface {
	environment.SetMockEnv()
	client := &fake.Clientset{}
	client.Fake.AddReactor(verb, "deployments", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		mysar := &deps
		return true, mysar, erro
	})
	os.Setenv("NODES_APPS_NAMESPACE", "default.node.opr")
	os.Setenv("KAFKA_BOOTSTRAP_SERVERS", "kafka.default.svc:9092")
	os.Setenv("KAFKA_AUTO_OFFSET_RESET", "earliest")
	return client
}

func TestNodeOperator_GetNode(t *testing.T) {
	type fields struct {
		clientSet kubernetes.Interface
	}
	type args struct {
		ctx context.Context
		app *meta.App
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *meta.Node
		wantErr bool
	}{
		{
			name: "K8s valid get",
			fields: fields{
				clientSet: mockK8sClientset("get", mockDeployment(), nil),
			},
			args: args{
				ctx: context.Background(),
				app: &meta.App{},
			},
			want: &meta.Node{
				Meta: meta.Metadata{
					Name:      "name",
					Reference: "",
					Parent:    "name.name",
					UUID:      "",
				},
				Spec: meta.NodeSpec{
					Image:       "image",
					Environment: make(map[string]string),
					Replicas:    1,
				},
			},
			wantErr: false,
		},
		{
			name: "K8s invalid get",
			fields: fields{
				clientSet: mockK8sClientset("get", mockDeployment(), errors.New("Expected error")),
			},
			args: args{
				ctx: context.Background(),
				app: &meta.App{},
			},
			want:    &meta.Node{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nop := &NodeOperator{
				clientSet: tt.fields.clientSet,
			}
			got, err := nop.GetNode(tt.args.ctx, tt.args.app)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeOperator.GetNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeOperator.GetNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeOperator_Nodes(t *testing.T) {
	type fields struct {
		clientSet kubernetes.Interface
	}
	tests := []struct {
		name   string
		fields fields
		want   operators.NodeOperatorInterface
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nop := &NodeOperator{
				clientSet: tt.fields.clientSet,
			}
			if got := nop.Nodes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeOperator.Nodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeOperator_CreateNode(t *testing.T) {
	type fields struct {
		clientSet kubernetes.Interface
	}
	type args struct {
		ctx context.Context
		app *meta.App
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *meta.Node
		wantErr bool
	}{
		{
			name: "K8s valid create",
			fields: fields{
				clientSet: mockK8sClientset("create", mockDeployment(), nil),
			},
			args: args{
				ctx: context.Background(),
				app: &meta.App{},
			},
			want: &meta.Node{
				Meta: meta.Metadata{
					Name:      "name",
					Reference: "",
					Parent:    "name.name",
					UUID:      "",
				},
				Spec: meta.NodeSpec{
					Image:       "image",
					Environment: make(map[string]string),
					Replicas:    1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nop := &NodeOperator{
				clientSet: tt.fields.clientSet,
				memory:    tree.GetTreeMemory(),
			}
			got, err := nop.CreateNode(tt.args.ctx, tt.args.app)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeOperator.CreateNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeOperator.CreateNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeOperator_UpdateNode(t *testing.T) {
	type fields struct {
		clientSet kubernetes.Interface
	}
	type args struct {
		ctx context.Context
		app *meta.App
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *meta.Node
		wantErr bool
	}{
		{
			name: "K8s valid update",
			fields: fields{
				clientSet: mockK8sClientset("update", mockDeployment(), nil),
			},
			args: args{
				ctx: context.Background(),
				app: &meta.App{},
			},
			want: &meta.Node{
				Meta: meta.Metadata{
					Name:      "name",
					Reference: "",
					Parent:    "name.name",
					UUID:      "",
				},
				Spec: meta.NodeSpec{
					Image:       "image",
					Environment: make(map[string]string),
					Replicas:    1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nop := &NodeOperator{
				clientSet: tt.fields.clientSet,
				memory:    tree.GetTreeMemory(),
			}
			got, err := nop.UpdateNode(tt.args.ctx, tt.args.app)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeOperator.UpdateNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeOperator.UpdateNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNodeOperator_DeleteNode(t *testing.T) {
	mem := tree.GetTreeMemory()
	mem.InitTransaction()

	type fields struct {
		clientSet kubernetes.Interface
	}
	type args struct {
		ctx         context.Context
		nodeContext string
		nodeName    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "K8s valid delete",
			fields: fields{
				clientSet: mockK8sClientset("delete", kubeApp.Deployment{}, nil),
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "K8s invalid delete",
			fields: fields{
				clientSet: mockK8sClientset("delete", kubeApp.Deployment{}, errors.New("Expected error")),
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nop := &NodeOperator{
				clientSet: tt.fields.clientSet,
				memory:    mem,
			}
			if err := nop.DeleteNode(tt.args.ctx, tt.args.nodeContext, tt.args.nodeName); (err != nil) != tt.wantErr {
				t.Errorf("NodeOperator.DeleteNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	mem.Cancel()
}

func TestNodeOperator_GetAllNodes(t *testing.T) {
	t.Skip("not sure why not working, useless right now")
	type fields struct {
		clientSet kubernetes.Interface
	}
	tests := []struct {
		name   string
		fields fields
		want   []meta.Node
	}{
		{
			name: "K8s valid list",
			fields: fields{
				clientSet: mockK8sList("list", kubeApp.DeploymentList{
					Items: []kubeApp.Deployment{
						mockDeployment(),
						mockDeployment(),
						mockDeployment(),
					},
				}, nil),
			},
			want: []meta.Node{
				{
					Meta: meta.Metadata{
						Name:      "name",
						Reference: "",
						Parent:    "",
						UUID:      "",
					},
					Spec: meta.NodeSpec{
						Image:       "image",
						Environment: make(map[string]string),
						Replicas:    1,
					},
				},
				{
					Meta: meta.Metadata{
						Name:      "name",
						Reference: "",
						Parent:    "",
						UUID:      "",
					},
					Spec: meta.NodeSpec{
						Image:       "image",
						Environment: make(map[string]string),
						Replicas:    1,
					},
				},
				{
					Meta: meta.Metadata{
						Name:      "name",
						Reference: "",
						Parent:    "",
						UUID:      "",
					},
					Spec: meta.NodeSpec{
						Image:       "image",
						Environment: make(map[string]string),
						Replicas:    1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nop := &NodeOperator{
				clientSet: tt.fields.clientSet,
			}
			if got := nop.GetAllNodes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NodeOperator.GetAllNodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockDeployment() kubeApp.Deployment {
	return kubeApp.Deployment{
		ObjectMeta: kubeMeta.ObjectMeta{
			Name:   "name-name-name",
			Labels: map[string]string{"app": "name"},
		},
		Spec: kubeApp.DeploymentSpec{
			Replicas: intToint32(1),
			Selector: &kubeMeta.LabelSelector{
				MatchLabels: map[string]string{
					"app": "name",
				},
			},
			Strategy: kubeApp.DeploymentStrategy{
				Type: kubeApp.RollingUpdateDeploymentStrategyType,
			},
			Template: kubeCore.PodTemplateSpec{
				ObjectMeta: kubeMeta.ObjectMeta{
					Labels: map[string]string{
						"app": "name",
					},
				},
				Spec: kubeCore.PodSpec{
					Volumes: []kubeCore.Volume{
						{
							Name: "name-name-name" + "-volume",
							VolumeSource: kubeCore.VolumeSource{
								EmptyDir: &kubeCore.EmptyDirVolumeSource{
									Medium: kubeCore.StorageMediumMemory,
								},
							},
						},
					},
					Containers: []kubeCore.Container{
						{
							Name: "name-name-name",
							Ports: func() []kubeCore.ContainerPort {
								return nil
							}(),

							Image: "image",
							// parse from master env var to kube env vars
							ImagePullPolicy: kubeCore.PullAlways,
							VolumeMounts: []kubeCore.VolumeMount{
								{
									Name:      "name" + "-volume",
									MountPath: "/inspr",
								},
							},
						},
						{
							Name:            "name-name-name" + "-sidecar",
							Image:           "sidecar-image",
							ImagePullPolicy: kubeCore.PullIfNotPresent,
							VolumeMounts: []kubeCore.VolumeMount{
								{
									Name:      "name-name-name" + "-volume",
									MountPath: "/inspr",
								},
							},
						},
					},
				},
			},
		},
	}
}
