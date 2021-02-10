package controller

import (
	"context"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// ChannelInterface is the interface that allows to obtain
// or change information related to the stored state of
// the Channels in the cluster
type ChannelInterface interface {
	Get(ctx context.Context, context string, chName string) (*meta.Channel, error)
	Create(ctx context.Context, context string, ch *meta.Channel) error
	Delete(ctx context.Context, context string, chName string) error
	Update(ctx context.Context, context string, ch *meta.Channel) error
}

// AppInterface is the interface that allows to obtain or
// change information related to the current state of
// the DApps in the cluster
type AppInterface interface {
	Get(ctx context.Context, query string) (*meta.App, error)
	Create(ctx context.Context, context string, app *meta.App) error
	Delete(ctx context.Context, query string) error
	Update(ctx context.Context, query string, app *meta.App) error
}

// ChannelTypeInterface is the interface that allows to
// obtain or change information related to the current
// state of the ChannelTypes in the cluster
type ChannelTypeInterface interface {
	Get(ctx context.Context, context string, ctName string) (*meta.ChannelType, error)
	Create(ctx context.Context, context string, ct *meta.ChannelType) error
	Delete(ctx context.Context, context string, ctName string) error
	Update(ctx context.Context, context string, ct *meta.ChannelType) error
}

// Interface is the interface that allows the management
// of the current state of the cluster. Permiting the
// modification of Channels, DApps and ChannelTypes
type Interface interface {
	Channels() ChannelInterface
	Apps() AppInterface
	ChannelTypes() ChannelTypeInterface
}
