package k8soperator

import (
	pb "gitlab.inspr.dev/inspr/core/pkg/operator/node"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// server is used to implement operator.OperatorServer.
type server struct {
	pb.UnimplementedOperatorServer
}

/*
// CreateNode implements operator.OperatorServer
func (s *server) CreateNode(ctx context.Context, in *pb.NodeDescription) (*pb.NodeReply, error) {
	log.Printf("Received: %v -> %v", in.GetNodeDescription(), in.GetValue())
	even = even + 1
	if even%2 == 1 {
		return &pb.OperationReply{
			Err:    "OK! OP (" + in.GetKind() + ") - " + in.GetValue() + "\nGo ahead!",
			Status: true,
		}, nil
	} else {
		return &pb.OperationReply{
			Err:    "Stop! Error to parser OP (" + in.GetKind() + ") \nPermission denied!",
			Status: false,
		}, nil
	}
}

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
