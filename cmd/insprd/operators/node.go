package operators

import (
	"context"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// NodeOperatorInterface is the interface that allows to obtain or change
// node information inside a deployment
type NodeOperatorInterface interface {
	CreateNode(ctx context.Context, app *meta.App) (*meta.Node, error)
	GetNode(ctx context.Context, app *meta.App) (*meta.Node, error)
	UpdateNode(ctx context.Context, app *meta.App) (*meta.Node, error)
	DeleteNode(ctx context.Context, nodeContext string, nodeName string) error
	GetAllNodes() []meta.Node
}