// Package memory provides simple interfaces for the
// in memory management of the cluster.
package memory

import "gitlab.inspr.dev/inspr/core/pkg/meta"

// ChannelMemory is the interface that allows to obtain
// or change information related to the stored state of
// the Channels in the cluster
type ChannelMemory interface {
	GetChannel(ref string) (*meta.Channel, error)
	CreateChannel(ch *meta.Channel) error
	DeleteChannel(ref string) error
	UpdateChannel(ch *meta.Channel, ref string) error
}

// AppMemory is the interface that allows to obtain or
// change information related to the current state of
// the DApps in the cluster
type AppMemory interface {
	GetApp(ref string) (*meta.AppSpec, error)
	CreateApp(app *meta.AppSpec) error
	DeleteApp(ref string) error
	UpdateApp(app *meta.AppSpec, ref string) error
}

// ChannelTypeMemory is the interface that allows to
// obtain or change information related to the current
// state of the ChannelTypes in the cluster
type ChannelTypeMemory interface {
	GetChannelType(ref string) (*meta.ChannelType, error)
	CreateChannelType(ct *meta.ChannelType) error
	DeleteChannelType(ref string) error
	UpdateChannelType(ct *meta.ChannelType, ref string) error
}

// Manager is the interface that allows the management
// of the current state of the cluster. Permiting the
// modification of Channels, DApps and ChannelTypes
type Manager interface {
	Channels() ChannelMemory
	Apps() AppMemory
	ChannelTypes() ChannelTypeMemory
}
