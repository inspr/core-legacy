package k8soperator

import (
	"context"

	pb "gitlab.inspr.dev/inspr/core/pkg/meta"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"gitlab.inspr.dev/inspr/core/cmd/operator/k8s/builder"

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
	DeleteNode(ctx context.Context, in *pb.NodeDescription) (*pb.NodeReply, error)
}

func newNodeReply(err error, value string) (*pb.NodeReply, error) {
	return &pb.NodeReply{
		Error: err.Error(),
		Value: string,
	}, err
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
	labels := map[string]string{"app": nodeName}
	return builder.NewPod().
		WithObjectMetadata(nodeName, s.namespace, labels).
		WithPodSelectorMatchLabels(labels).
		WithPodTemplateObjectMetadata(nodeName, s.namespace, labels).
		WithPodTemplateSpec(nodeName, s.node.Spec.Image).
		GetDeployment()
}

// CreateNode implements server.CreateNode
func (s *server) CreateNode(ctx context.Context, in *pb.Node) (*pb.NodeReply, error) {
	s.node = in

	namespace := s.namespace
	kubeDeployment := s.createKubeDeployment()
	deployment, err := s.kubeClient.AppsV1().Deployments(namespace).Create(kubeDeployment)
	if err != nil {
		return newNodeReply(err, "")
	}

	nodeName := s.node.Metadata.Name + s.node.Metadata.Sha256
	s.deployments[nodeName] = deployment

	return newNodeReply(nil, nodeName)
}

// DeleteNode implements server.CreateNode
func (s *server) DeleteNode(ctx context.Context, in *pb.NodeDescription) (*pb.NodeReply, error) {
	nodeName := in.NodeDescription
	err := s.kubeClient.AppsV1().Deployments(s.namespace).Delete(nodeName, nil)
	if err != nil {
		return newNodeReply(err, "")
	}
	return newNodeReply(nil, nodeName)
}

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
