package nodes

import (
	"context"
	"errors"
	"os"
	"testing"

	"inspr.dev/inspr/cmd/insprd/memory/tree"
	authmock "inspr.dev/inspr/pkg/auth/mocks"
	"inspr.dev/inspr/pkg/environment"
	"inspr.dev/inspr/pkg/meta"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	kubeApp "k8s.io/api/apps/v1"
	kubeCore "k8s.io/api/core/v1"
	kubeMeta "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8stesting "k8s.io/client-go/testing"
)

// ignores unused code for this file in the staticcheck
//lint:file-ignore U1000 Ignore all unused code

type k8sKind string

const (
	k8sService    k8sKind = "services"
	k8sSecret     k8sKind = "secrets"
	k8sDeployment k8sKind = "deployments"
)

func mockK8sClientSet(verbsToDepAndErr map[k8sKind]map[string]struct {
	runtime.Object
	error
}) kubernetes.Interface {
	environment.SetMockEnv()
	client := &fake.Clientset{}

	os.Setenv("NODES_APPS_NAMESPACE", "default.node.opr")
	os.Setenv("INSPR_SIDECAR_KAFKA_BOOTSTRAP_SERVERS", "kafka.default.svc:9092")
	os.Setenv("INSPR_SIDECAR_KAFKA_AUTO_OFFSET_RESET", "earliest")
	for verb, res := range verbsToDepAndErr {
		for resource, act := range res {
			client.Fake.AddReactor(resource, string(verb), func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
				mysar := act.Object
				return true, mysar, act.error
			})
		}
	}
	return client
}

func mockK8sClientset(verb string, dep kubeApp.Deployment, svc kubeCore.Service, erro error, serviceErr error) kubernetes.Interface {
	environment.SetMockEnv()
	client := &fake.Clientset{}
	client.Fake.AddReactor(verb, "deployments", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		mysar := &dep
		return true, mysar, erro
	})
	client.Fake.AddReactor(verb, "services", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		mysar := &svc
		return true, mysar, serviceErr
	})
	os.Setenv("NODES_APPS_NAMESPACE", "default.node.opr")
	os.Setenv("INSPR_SIDECAR_KAFKA_BOOTSTRAP_SERVERS", "kafka.default.svc:9092")
	os.Setenv("INSPR_SIDECAR_KAFKA_AUTO_OFFSET_RESET", "earliest")
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
	os.Setenv("INSPR_SIDECAR_KAFKA_BOOTSTRAP_SERVERS", "kafka.default.svc:9092")
	os.Setenv("INSPR_SIDECAR_KAFKA_AUTO_OFFSET_RESET", "earliest")
	return client
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
				clientSet: mockK8sClientSet(
					map[k8sKind]map[string]struct {
						runtime.Object
						error
					}{
						k8sService: {
							"create": {
								error:  nil,
								Object: mockService(),
							},
						},
						k8sDeployment: {
							"create": {
								error:  nil,
								Object: mockDeployment(),
							},
						},
					},
				),
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
			name: "K8s failure in k8s create",
			fields: fields{
				clientSet: mockK8sClientSet(
					map[k8sKind]map[string]struct {
						runtime.Object
						error
					}{
						k8sDeployment: {
							"create": {
								error:  errors.New("error in createion"),
								Object: nil,
							},
						},
						k8sService: {
							"create": {
								error:  nil,
								Object: mockService(),
							},
						},
					},
				),
			},
			args: args{
				ctx: context.Background(),
				app: &meta.App{},
			},
			wantErr: true,
		},

		{
			name: "K8s failure in k8s create service",
			fields: fields{
				clientSet: mockK8sClientSet(
					map[k8sKind]map[string]struct {
						runtime.Object
						error
					}{
						k8sService: {
							"create": {
								error:  errors.New("this is a bad creation error"),
								Object: mockService(),
							},
						},
						k8sDeployment: {
							"create": {
								error:  nil,
								Object: mockDeployment(),
							},
						},
					},
				),
			},
			args: args{
				ctx: context.Background(),
				app: &meta.App{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nop := &NodeOperator{
				clientSet: tt.fields.clientSet,
				auth:      authmock.NewMockAuth(nil),
				memory:    tree.GetTreeMemory(),
			}
			_, err := nop.CreateNode(tt.args.ctx, tt.args.app)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeOperator.CreateNode() error = %v, wantErr %v", err, tt.wantErr)
				return
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
				clientSet: mockK8sClientSet(
					map[k8sKind]map[string]struct {
						runtime.Object
						error
					}{
						k8sService: {
							"update": {
								error:  nil,
								Object: mockService(),
							},
						},
						k8sDeployment: {
							"update": {
								error:  nil,
								Object: mockDeployment(),
							},
						},
					},
				),
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
			name: "K8s failure in k8s update",
			fields: fields{
				clientSet: mockK8sClientSet(
					map[k8sKind]map[string]struct {
						runtime.Object
						error
					}{
						k8sDeployment: {
							"update": {
								error:  errors.New("error in createion"),
								Object: nil,
							},
						},
						k8sService: {
							"update": {
								error:  nil,
								Object: mockService(),
							},
						},
					},
				),
			},
			args: args{
				ctx: context.Background(),
				app: &meta.App{},
			},
			wantErr: true,
		},

		{
			name: "K8s failure in k8s update service",
			fields: fields{
				clientSet: mockK8sClientSet(
					map[k8sKind]map[string]struct {
						runtime.Object
						error
					}{
						k8sService: {
							"update": {
								error:  errors.New("this is a bad creation error"),
								Object: mockService(),
							},
						},
						k8sDeployment: {
							"update": {
								error:  nil,
								Object: mockDeployment(),
							},
						},
					},
				),
			},
			args: args{
				ctx: context.Background(),
				app: &meta.App{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nop := &NodeOperator{
				clientSet: tt.fields.clientSet,
				auth:      authmock.NewMockAuth(nil),
				memory:    tree.GetTreeMemory(),
			}
			_, err := nop.UpdateNode(tt.args.ctx, tt.args.app)
			if (err != nil) != tt.wantErr {
				t.Errorf("NodeOperator.UpdateNode() error = %v, wantErr %v", err, tt.wantErr)
				return
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
				clientSet: mockK8sClientSet(
					map[k8sKind]map[string]struct {
						runtime.Object
						error
					}{
						k8sService: {
							"delete": {
								error:  nil,
								Object: mockService(),
							},
						},
						k8sDeployment: {
							"delete": {
								error:  nil,
								Object: mockDeployment(),
							},
						},
					},
				),
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "K8s invalid delete",
			fields: fields{
				clientSet: mockK8sClientSet(
					map[k8sKind]map[string]struct {
						runtime.Object
						error
					}{
						k8sService: {
							"delete": {
								error:  nil,
								Object: mockService(),
							},
						},
						k8sDeployment: {
							"delete": {
								error:  errors.New("error in deleting deployment"),
								Object: mockDeployment(),
							},
						},
					},
				)},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "K8s invalid delete service",
			fields: fields{
				clientSet: mockK8sClientSet(
					map[k8sKind]map[string]struct {
						runtime.Object
						error
					}{
						k8sService: {
							"delete": {
								error:  errors.New("error in deleting service"),
								Object: mockService(),
							},
						},
						k8sDeployment: {
							"delete": {
								error:  nil,
								Object: mockDeployment(),
							},
						},
					},
				)},
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
				auth:      authmock.NewMockAuth(nil),
			}
			if err := nop.DeleteNode(tt.args.ctx, tt.args.nodeContext, tt.args.nodeName); (err != nil) != tt.wantErr {
				t.Errorf("NodeOperator.DeleteNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	mem.Cancel()
}

func mockService() *kubeCore.Service {
	return &kubeCore.Service{
		ObjectMeta: kubeMeta.ObjectMeta{
			Name: "name-name-name",
		},
		Spec: kubeCore.ServiceSpec{
			Selector: map[string]string{
				"app": "name",
			},
		},
	}
}

func mockDeployment() *kubeApp.Deployment {
	return &kubeApp.Deployment{
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
