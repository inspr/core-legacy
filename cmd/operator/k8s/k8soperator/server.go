package k8soperator

import (
	"context"

	pb "gitlab.inspr.dev/inspr/core/pkg/meta"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// server is used to implement operator.OperatorServer.
type server struct {
	pb.UnimplementedNodeOperatorServer
	kubeClient  kubernetes.Interface
	namespace   string
	replicas    int32
	node        *pb.Node
	deployments map[string]*appsv1.Deployment
}

// Server is an interface defined to expose server's functions
type Server interface {
	createKubeDeployment() *appsv1.Deployment
	CreateNode(ctx context.Context, in *pb.Node) (*pb.NodeReply, error)
}

// NewServer instanciate a new server
func NewServer() (Server, error) {
	config, clusterCfgErr := rest.InClusterConfig()
	if clusterCfgErr != nil {
		return nil, clusterCfgErr
	}

	clientSet, errClientSet := kubernetes.NewForConfig(config)
	if errClientSet != nil {
		return nil, errClientSet
	}
	return &server{
		kubeClient: clientSet,
		replicas:   1,
		namespace:  "inspr-apps",
	}, nil
}

func (s *server) createKubeDeployment() *appsv1.Deployment {
	nodeName := s.node.Metadata.Name + s.node.Metadata.Sha256
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nodeName,
			Namespace: s.node.Metadata.Parent,
			Labels:    map[string]string{"app": nodeName},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &s.replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": nodeName},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": nodeName},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name: nodeName,
							// parse from master env var to kube env vars
							ImagePullPolicy: apiv1.PullAlways,
							Env: []apiv1.EnvVar{
								{
									Name:  "UUID",
									Value: s.node.Metadata.Sha256,
								},
							},
						},
					},
				},
			},
		},
	}
}

// CreateNode implements server.CreateNode
func (s *server) CreateNode(ctx context.Context, in *pb.Node) (*pb.NodeReply, error) {
	s.node = in

	namespace := s.namespace
	kubeDeployment := s.createKubeDeployment()
	deployment, err := s.kubeClient.AppsV1().Deployments(namespace).Create(kubeDeployment)
	if err != nil {
		return &pb.NodeReply{
			Error: err.Error(),
			Value: "",
		}, nil
	}

	nodeName := s.node.Metadata.Name + s.node.Metadata.Sha256
	s.deployments[nodeName] = deployment

	return &pb.NodeReply{
		Error: "",
		Value: nodeName,
	}, nil
}

/*
func main() {
	//even = 0
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOperatorServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}*/

// UpdateNodeStatus returns the up to date pod status for each node
func UpdateNodeStatus() map[string]string {
	var status map[string]string
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	options := metav1.ListOptions{FieldSelector: "metadata.name=kubernetes"}
	pods, _ := clientset.CoreV1().Pods("inspr").List(options)

	for _, pod := range pods.Items {
		status[pod.Name] = pod.Status.Message
	}
	return status
}
