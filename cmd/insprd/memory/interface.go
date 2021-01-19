package memory

import "gitlab.inspr.dev/inspr/core/pkg/meta"

// ChannelMemory In memory channel methods
type ChannelMemory interface {
	GetChannel(ref string) (*meta.Channel, error)
	CreateChannel(ch *meta.Channel) error
	DeleteChannel(ref string) error
	UpdateChannel(ch *meta.Channel, ref string) error
}

// AppMemory In memory app methods
type AppMemory interface {
	GetApp(ref string) (*meta.DApp, error)
	CreateApp(app *meta.DApp) error
	DeleteApp(ref string) error
	UpdateApp(app *meta.DApp, ref string) error
}

// ChannelTypeMemory In memory channelType methods
type ChannelTypeMemory interface {
	GetChannelType(ref string) (*meta.ChannelType, error)
	CreateChannelType(ct *meta.ChannelType) error
	DeleteChannelType(ref string) error
	UpdateChannelType(ct *meta.ChannelType, ref string) error
}

// Interface In memory general methods
type Interface interface {
	Channel() ChannelMemory
	App() AppMemory
	ChannelType() ChannelTypeMemory
}
