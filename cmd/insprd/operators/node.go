package operators

import (
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// NodeOperatorInterface is the interface that allows to obtain or change
// node information inside a deployment
type NodeOperatorInterface interface {
	CreateNode(context string, node *meta.Node) error
	GetNode(context string, nodeName string) (*meta.Node, error)
	UpdateNode(context string, node *meta.Node) error
	DeleteNode(context, nodeName string) error
	GetAllNodes() []*meta.Node
}

// NodeOperator defines a NodeOperatorInterface
type NodeOperator struct {
	node *meta.Node
}

// Nodes is a NodeOperatorInterface that provides methods for node manipulation
func (no *NodeOperator) Nodes() NodeOperatorInterface {
	return &NodeOperator{
		node: &meta.Node{},
	}
}

// CreateNode deploys a new node structure, if it's information is valid.
// Otherwise, returns an error
func (no *NodeOperator) CreateNode(context string, node *meta.Node) error {
	return nil
}

// GetNode returns the node with the given name, if it exists.
// Otherwise, returns an error
func (no *NodeOperator) GetNode(context string, nodeName string) (*meta.Node, error) {
	return &meta.Node{}, nil
}

// UpdateNode updates a node that already exists, if the new structure is valid.
// Otherwise, returns an error.
func (no *NodeOperator) UpdateNode(context string, node *meta.Node) error {
	return nil
}

// DeleteNode deletes node with given name, if it exists. Otherwise, returns an error
func (no *NodeOperator) DeleteNode(context, nodeName string) error {
	return nil
}

// GetAllNodes returns a list of all the active nodes in the deployment, if there are any
func (no *NodeOperator) GetAllNodes() []*meta.Node {
	return nil
}
