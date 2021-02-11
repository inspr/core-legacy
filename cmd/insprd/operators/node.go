package operators

import "gitlab.inspr.dev/inspr/core/pkg/meta"

// NodeOperatorInterface is the interface that allows to obtain or change
// node information inside a deployment
type NodeOperatorInterface interface {
	CreateNode(context string, node *meta.Node) error
	GetNode(context string, nodeName string) (*meta.Node, error)
	UpdateNode(context string, node *meta.Node) error
	DeleteNode(context, nodeName string) error
	GetAllNodes() []*meta.Node
}
