package nodes

import (
	"context"
	"fmt"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/tree"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"

	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type NodeOperator struct {
	clientSet kubernetes.Interface
}

func (nop *NodeOperator) retrieveKube() v1.DeploymentInterface {
	appsNamespace := GetK8SVariables().AppsNamespace
	return nop.clientSet.AppsV1().Deployments(appsNamespace)
}

// GetNodes returns the node with the given name, if it exists.
// Otherwise, returns an error
func (nop *NodeOperator) GetNodes(ctx context.Context, node *meta.node) (*meta.Node, error) {
	kube := nop.retrieveKube()
	nodeName := 
	dep, err := kube.Get(nodeName, metav1.GetOptions{})

}

// Nodes is a NodeOperatorInterface that provides methods for node manipulation
func (no *NodeOperator) Nodes() NodeOperatorInterface {
	return &NodeOperator{}
}

// CreateNode deploys a new node structure, if it's information is valid.
// Otherwise, returns an error
func (no *NodeOperator) CreateNode(ctx context.Context, context string, node *meta.Node) error {
	return nil
}

// UpdateNode updates a node that already exists, if the new structure is valid.
// Otherwise, returns an error.
func (no *NodeOperator) UpdateNode(ctx context.Context, node *meta.node) error {
	//generate dep from node
	//update deb so that it matches node
}

// DeleteNode deletes node with given name, if it exists. Otherwise, returns an error
func (no *NodeOperator) DeleteNode(ctx context.Context, context string, nodeName string) error {
	return nil
}

// GetAllNodes returns a list of all the active nodes in the deployment, if there are any
func (no *NodeOperator) GetAllNodes() []*meta.Node {
	kube := nop.retrieveKube()
	
}

func parseNodeName(node *meta.Node) string{
	return "placeholder" //TODO
}
