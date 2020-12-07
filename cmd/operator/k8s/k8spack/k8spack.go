package k8spack

import (
	"log"
	"net"

	k8sop "gitlab.inspr.dev/inspr/core/cmd/operator/k8s/server"
	pb "gitlab.inspr.dev/inspr/core/pkg/meta"

	"google.golang.org/grpc"
)

const (
	port = ":50000"
)

// NewK8sOperator creates a new GRPC instance of node operator
func NewK8sOperator(url string, port string) {
	address := url + ":" + port
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	_, err = k8sop.NewServer()
	if err != nil {
		panic(err)
	}

	pb.RegisterNodeOperatorServer(s, &k8sop.ServerObject)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
