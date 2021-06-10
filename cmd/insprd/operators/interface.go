package operators

import (
	"context"

	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/brokers"
)

// NodeOperatorInterface is the interface that allows to obtain or change
// node information inside a deployment
type NodeOperatorInterface interface {
	CreateNode(ctx context.Context, app *meta.App) (*meta.Node, error)
	UpdateNode(ctx context.Context, app *meta.App) (*meta.Node, error)
	DeleteNode(ctx context.Context, scope string, name string) error
}

// ChannelOperatorInterface is responsible for handling the following methods
//
// 	- `Get`: returns a channel from the DApp of the given context
//	- `GetAll`: return all channels from the DApp of the given context
// 	- `Create`: creates a channel in the DApp of the given context
// 	- `Update`: updates a channel in the DApp of the given context
// 	- `Delete`: deletes a channel of the specified name in DApp of the given context
type ChannelOperatorInterface interface {
	Get(ctx context.Context, scope string, name string) (*meta.Channel, error)
	Create(ctx context.Context, scope string, channel *meta.Channel) error
	Update(ctx context.Context, scope string, channel *meta.Channel) error
	Delete(ctx context.Context, scope string, name string) error
}

// OperatorInterface is an interface for inspr runtime operators
//
// To implement the interface you need to create two implementations,
// a node implementation, that creates nodes from inspr in the given runtime
// and a channel implementation, that creates channels from inspr in the given
// runtime.
type OperatorInterface interface {
	Nodes() NodeOperatorInterface
	Channels() ChannelOperatorInterface
	SetBrokerOperator(config brokers.BrokerConfiguration) error
}
