package controller

import (
	"context"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/utils/diff"
)

// ChannelInterface is the interface that allows to obtain
// or change information related to the stored state of
// the Channels in the cluster
type ChannelInterface interface {
	Get(ctx context.Context, context string, chName string) (*meta.Channel, error)
	Create(ctx context.Context, context string, ch *meta.Channel) (diff.Changelog, error)
	Delete(ctx context.Context, context string, chName string) (diff.Changelog, error)
	Update(ctx context.Context, context string, ch *meta.Channel) (diff.Changelog, error)
}

// AppInterface is the interface that allows to obtain or
// change information related to the current state of
// the DApps in the cluster
type AppInterface interface {
	Get(ctx context.Context, query string) (*meta.App, error)
	Create(ctx context.Context, context string, app *meta.App) (diff.Changelog, error)
	Delete(ctx context.Context, query string) (diff.Changelog, error)
	Update(ctx context.Context, query string, app *meta.App) (diff.Changelog, error)
}

// ChannelTypeInterface is the interface that allows to
// obtain or change information related to the current
// state of the ChannelTypes in the cluster
type ChannelTypeInterface interface {
	Get(ctx context.Context, context string, ctName string) (*meta.ChannelType, error)
	Create(ctx context.Context, context string, ct *meta.ChannelType) (diff.Changelog, error)
	Delete(ctx context.Context, context string, ctName string) (diff.Changelog, error)
	Update(ctx context.Context, context string, ct *meta.ChannelType) (diff.Changelog, error)
}

// Interface is the interface that allows the management
// of the current state of the cluster. Permiting the
// modification of Channels, DApps and ChannelTypes
type Interface interface {
	Channels() ChannelInterface
	Apps() AppInterface
	ChannelTypes() ChannelTypeInterface
}
