package operators

import (
	"context"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// ChannelOperatorInterface is responsible for handling the following methods
//
// 	- `Get`: returns a channel from the DApp of the given context
//	- `GetAll`: return all channels from the DApp of the given context
// 	- `Create`: creates a channel in the DApp of the given context
// 	- `Update`: updates a channel in the DApp of the given context
// 	- `Delete`: deletes a channel of the specified name in DApp of the given context
type ChannelOperatorInterface interface {
	Get(ctx context.Context, context string, name string) (*meta.Channel, error)
	GetAll(ctx context.Context, context string) ([]*meta.Channel, error)
	Create(ctx context.Context, context string, channel *meta.Channel) error
	Update(ctx context.Context, context string, channel *meta.Channel) error
	Delete(ctx context.Context, context string, name string) error
}
