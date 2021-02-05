package controller

import (
	"context"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// ChannelInterface is the interface that allows to obtain
// or change information related to the stored state of
// the Channels in the cluster
type ChannelInterface interface {
	GetChannel(ctx context.Context, context string, chName string) (*meta.Channel, error)
	CreateChannel(ctx context.Context, context string, ch *meta.Channel) error
	DeleteChannel(ctx context.Context, context string, chName string) error
	UpdateChannel(ctx context.Context, context string, ch *meta.Channel) error
}

// AppInterface is the interface that allows to obtain or
// change information related to the current state of
// the DApps in the cluster
type AppInterface interface {
	GetApp(ctx context.Context, query string) (*meta.App, error)
	CreateApp(ctx context.Context, app *meta.App, context string) error
	DeleteApp(ctx context.Context, query string) error
	UpdateApp(ctx context.Context, app *meta.App, query string) error
}

// ChannelTypeInterface is the interface that allows to
// obtain or change information related to the current
// state of the ChannelTypes in the cluster
type ChannelTypeInterface interface {
	GetChannelType(ctx context.Context, context string, ctName string) (*meta.ChannelType, error)
	CreateChannelType(ctx context.Context, ct *meta.ChannelType, context string) error
	DeleteChannelType(ctx context.Context, context string, ctName string) error
	UpdateChannelType(ctx context.Context, ct *meta.ChannelType, context string) error
}

// Interface is the interface that allows the management
// of the current state of the cluster. Permiting the
// modification of Channels, DApps and ChannelTypes
type Interface interface {
	Channels() ChannelInterface
	Apps() AppInterface
	ChannelTypes() ChannelTypeInterface
}
