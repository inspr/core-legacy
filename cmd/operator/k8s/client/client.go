package k8sclient

import (
	"context"
	"log"
	"time"

	meta "gitlab.inspr.dev/inspr/core/pkg/meta"
	"google.golang.org/grpc"
)

// Metadata identify the metadata struct
type Metadata meta.Metadata

// Node identify the node struct
type Node meta.Node

// NodeSpec identify the node spec struct
type NodeSpec meta.NodeSpec

// NodeArray identify the node array struct
type NodeArray meta.NodeArray

// NodeWithDescription identify the node with its description struct
type NodeWithDescription meta.NodeWithDescription

// NodeDescription identify the node description struct
type NodeDescription meta.NodeDescription

// NodeReply identify the node reply struct
type NodeReply meta.NodeReply

type k8sClient struct {
	// The address of a cluster operator
	clusterOperatorAddrress       string
	clusterOperatorCancelFunction context.CancelFunc
	clusterOperator               meta.NodeOperatorClient
	clusterOperatorContext        context.Context
}

// K8sClient defines a k8s operator client
type K8sClient interface {
	CreateNode(node *Node) (*NodeReply, error)
	DeleteNode(nodeDescription *NodeDescription) (*NodeReply, error)
	UpdateNode(nodeWithDescription *NodeWithDescription) (*NodeReply, error)
	GetNode(nodeDescription *NodeDescription) (*Node, error)
	ListNodes() (*NodeArray, error)
	UpdateNodeStatus() (*NodeArray, error)
}

// NewK8sOpClient instanciate a new k8s operator client
func NewK8sOpClient(url string, port string) K8sClient {
	var c *k8sClient
	c.instantiateOperator(url, port)
	return c
}

func (c *k8sClient) instantiateOperator(url string, port string) {
	c.clusterOperatorAddrress = url + ":" + port
	// Set up a connection to the server.
	conn, err := grpc.Dial(c.clusterOperatorAddrress,
		grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c.clusterOperator = meta.NewNodeOperatorClient(conn)

	c.clusterOperatorContext, c.clusterOperatorCancelFunction =
		context.WithTimeout(context.Background(), time.Second)
}

func (c *k8sClient) CreateNode(node *Node) (*NodeReply, error) {
	nodeObj := meta.Node{
		Metadata: node.Metadata,
		Spec:     node.Spec,
	}
	reply, err := c.clusterOperator.CreateNode(c.clusterOperatorContext, &nodeObj)
	replyObj := NodeReply{
		Error: reply.Error,
		Value: reply.Value,
	}
	return &replyObj, err
}

func (c *k8sClient) DeleteNode(nodeDescription *NodeDescription) (*NodeReply, error) {
	nodeObj := meta.NodeDescription{
		NodeDescription: nodeDescription.NodeDescription,
	}
	reply, err := c.clusterOperator.DeleteNode(c.clusterOperatorContext, &nodeObj)
	replyObj := NodeReply{
		Error: reply.Error,
		Value: reply.Value,
	}
	return &replyObj, err
}

func (c *k8sClient) UpdateNode(nodeWithDescription *NodeWithDescription) (*NodeReply, error) {
	nodeObj := meta.NodeWithDescription{
		NodeDescription: nodeWithDescription.NodeDescription,
		TargetNode:      nodeWithDescription.TargetNode,
	}
	reply, err := c.clusterOperator.UpdateNode(c.clusterOperatorContext, &nodeObj)
	replyObj := NodeReply{
		Error: reply.Error,
		Value: reply.Value,
	}
	return &replyObj, err
}

func (c *k8sClient) GetNode(nodeDescription *NodeDescription) (*Node, error) {
	nodeDescObj := meta.NodeDescription{
		NodeDescription: nodeDescription.NodeDescription,
	}
	node, err := c.clusterOperator.GetNode(c.clusterOperatorContext, &nodeDescObj)
	nodeObj := Node{
		Metadata: node.Metadata,
		Spec:     node.Spec,
	}
	return &nodeObj, err
}

func (c *k8sClient) ListNodes() (*NodeArray, error) {
	reply, err := c.clusterOperator.ListNodes(c.clusterOperatorContext, &meta.Stub{})
	replyObj := NodeArray{
		Node:   reply.Node,
		Status: reply.Status,
	}
	return &replyObj, err
}

func (c *k8sClient) UpdateNodeStatus() (*NodeArray, error) {
	reply, err := c.clusterOperator.UpdateNodeStatus(c.clusterOperatorContext, &meta.Stub{})
	replyObj := NodeArray{
		Node:   reply.Node,
		Status: reply.Status,
	}
	return &replyObj, err
}
