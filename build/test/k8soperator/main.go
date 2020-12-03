package main

import (
	"context"
	"log"
	"time"

	pb "gitlab.inspr.dev/inspr/core/pkg/meta"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50000"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNodeOperatorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// Create a node into the cluster
	r, err := c.CreateNode(ctx, &pb.Node{
		Metadata: &pb.Metadata{
			Name:      "app01",
			Reference: "inspr",
		},
		Spec: &pb.NodeSpec{
			Image: "docker.io/hello-world:latest",
		},
	})
	if err != nil {
		log.Fatalf("could not create app01: %v", err)
	} else {
		log.Println(r.GetError(), r.GetValue())
	}

	// List the nodes at the cluster
	nodeArray, err := c.ListNodes(ctx, &pb.Stub{})

	if err != nil {
		log.Fatalf("could not create app01: %v", err)
	} else {
		log.Println(nodeArray.Node[0])
	}

	// Update the current node into the cluster
	r, err = c.UpdateNode(ctx, &pb.NodeWithDescription{
		TargetNode: &pb.Node{
			Metadata: &pb.Metadata{
				Name:      "app01",
				Reference: "inspr",
				Parent:    "inspr",
			},
			Spec: &pb.NodeSpec{
				Image: "docker.io/hello-world:latest",
			},
		},
		NodeDescription: "app01",
	})
	if err != nil {
		log.Fatalf("could not create app01: %v", err)
	} else {
		log.Println(r.GetError(), r.GetValue())
	}

	// Get the node at the cluster
	node, err := c.GetNode(ctx, &pb.NodeDescription{NodeDescription: "app01"})

	if err != nil {
		log.Fatalf("could not create app01: %v", err)
	} else {
		log.Println(node.Metadata.GetParent())
	}

	// Get the current state of a pod in the cluster
	nodesStatus, err := c.UpdateNodeStatus(ctx, &pb.Stub{})

	if err != nil {
		log.Fatalf("could not create app01: %v", err)
	} else {
		if len(nodesStatus.Node) > 0 {
			log.Println(nodesStatus.Node[0])
			log.Println(nodesStatus.Status[0])
		}
	}

	// Delete the node of the cluster
	r, err = c.DeleteNode(ctx, &pb.NodeDescription{NodeDescription: "app01"})
	if err != nil {
		log.Fatalf("could not create app01: %v", err)
	} else {
		log.Println(r.GetError(), r.GetValue())
	}
	// List the nodes at the cluster
	nodeArray, err = c.ListNodes(ctx, &pb.Stub{})

	if err != nil {
		log.Fatalf("could not create app01: %v", err)
	} else {
		log.Println("Lenght is: ", len(nodeArray.Node))
	}
}
