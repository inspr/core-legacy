package k8soperator

import (
	"context"
	"errors"

	pb "gitlab.inspr.dev/inspr/core/pkg/meta"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"gitlab.inspr.dev/inspr/core/cmd/operator/k8s/builder"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ServerStruct is used to implement operator.OperatorServer.
type ServerStruct struct {
	pb.UnimplementedNodeOperatorServer
	kubeClient  kubernetes.Interface
	namespace   string
	replicas    int32
	deployments map[string]*appsv1.Deployment
	nodes       map[string]*pb.Node
}

// ServerObject is the object to store server data
var ServerObject ServerStruct

// Server is an interface defined to expose server's functions
type Server interface {
	CreateNode(ctx context.Context, in *pb.Node) (*pb.NodeReply, error)
	DeleteNode(ctx context.Context, in *pb.NodeDescription) (*pb.NodeReply, error)
	UpdateNode(ctx context.Context, in *pb.NodeWithDescription) (*pb.NodeReply, error)
	GetNode(ctx context.Context, in *pb.NodeDescription) (*pb.Node, error)
	ListNodes(ctx context.Context, in *pb.Stub) (*pb.NodeArray, error)
	UpdateNodeStatus(ctx context.Context, in *pb.Stub) (*pb.NodeArray, error)
}

func newNodeReply(err error, value string) (*pb.NodeReply, error) {
	return &pb.NodeReply{
		Error: err.Error(),
		Value: value,
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

	ServerObject = ServerStruct{
		kubeClient: clientSet,
		replicas:   1,
		namespace:  "inspr-apps",
	}

	return &ServerObject, nil
}

func (s *ServerStruct) createKubeDeployment(node *pb.Node) *appsv1.Deployment {
	nodeName := node.Metadata.Name
	labels := map[string]string{"app": nodeName}
	return builder.NewPod().
		WithObjectMetadata(nodeName, s.namespace, labels).
		WithPodSelectorMatchLabels(labels).
		WithPodTemplateObjectMetadata(nodeName, s.namespace, labels).
		WithPodTemplateSpec(nodeName, node.Spec.Image).
		GetDeployment()
}

// CreateNode implements server.CreateNode
func (s *ServerStruct) CreateNode(ctx context.Context, in *pb.Node) (*pb.NodeReply, error) {
	namespace := s.namespace
	kubeDeployment := s.createKubeDeployment(in)
	deployment, err := s.kubeClient.AppsV1().Deployments(namespace).Create(kubeDeployment)
	if err != nil {
		return newNodeReply(err, "")
	}

	nodeName := in.Metadata.Name
	s.deployments[nodeName] = deployment
	s.nodes[nodeName] = in

	return newNodeReply(nil, nodeName)
}

// DeleteNode implements server.CreateNode
func (s *ServerStruct) DeleteNode(ctx context.Context, in *pb.NodeDescription) (*pb.NodeReply, error) {
	nodeName := in.NodeDescription
	err := s.kubeClient.AppsV1().Deployments(s.namespace).Delete(nodeName, nil)
	if err != nil {
		return newNodeReply(err, "")
	}
	return newNodeReply(nil, nodeName)
}

// UpdateNode deletes and create a new pod, due to the specs used
func (s *ServerStruct) UpdateNode(ctx context.Context, in *pb.NodeWithDescription) (*pb.NodeReply, error) {
	nodeDescription := pb.NodeDescription{
		NodeDescription: in.GetNodeDescription(),
	}

	reply, err := s.DeleteNode(ctx, &nodeDescription)
	if err != nil {
		return reply, err
	}

	reply, err = s.CreateNode(ctx, in.TargetNode)
	return reply, err
}

// GetNode recover the node used to create a pod
func (s *ServerStruct) GetNode(ctx context.Context, in *pb.NodeDescription) (*pb.Node, error) {
	_, ok := s.nodes[in.GetNodeDescription()]
	if !ok {
		return s.nodes[in.GetNodeDescription()], nil
	}
	return nil, errors.New("Node Not Found")
}

// ListNodes list all nodes running in a cluster
func (s *ServerStruct) ListNodes(ctx context.Context, in *pb.Stub) (*pb.NodeArray, error) {
	var nodeArray *pb.NodeArray
	for _, node := range s.nodes {
		nodeArray.Node = append(nodeArray.Node, node)
	}
	return nodeArray, nil
}

// UpdateNodeStatus returns the up to date pod status for each node
func (s *ServerStruct) UpdateNodeStatus(ctx context.Context, in *pb.Stub) (*pb.NodeArray, error) {

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

	var nodeArray *pb.NodeArray
	for _, pod := range pods.Items {
		if value, ok := s.nodes[pod.Name]; ok {
			nodeArray.Node = append(nodeArray.Node, value)
		} else {
			nodeArray.Node = append(nodeArray.Node, nil)
		}

		nodeArray.Status = append(nodeArray.Status, pod.Status.Message)
	}
	return nodeArray, nil
}
