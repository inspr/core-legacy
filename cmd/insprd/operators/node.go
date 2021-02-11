package operators

import (
	"context"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// NodeOperatorInterface is the interface that allows to obtain or change
// node information inside a deployment
type NodeOperatorInterface interface {
	CreateNode(ctx context.Context, context string, node *meta.Node) error
	GetNode(ctx context.Context) (*meta.Node, error)
	UpdateNode(ctx context.Context, context string, node *meta.Node) error
	DeleteNode(ctx context.Context, context, nodeName string) error
	GetAllNodes() []*meta.Node
}
