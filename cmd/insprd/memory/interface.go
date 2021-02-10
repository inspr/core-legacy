// Package memory provides simple interfaces for the
// in memory management of the cluster.
package memory

import (
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/utils/diff"
)

// ChannelMemory is the interface that allows to obtain
// or change information related to the stored state of
// the Channels in the cluster
type ChannelMemory interface {
	GetChannel(context string, chName string) (*meta.Channel, error)
	CreateChannel(context string, ch *meta.Channel) error
	DeleteChannel(context string, chName string) error
	UpdateChannel(context string, ch *meta.Channel) error
}

// AppMemory is the interface that allows to obtain or
// change information related to the current state of
// the DApps in the cluster
type AppMemory interface {
	GetApp(query string) (*meta.App, error)
	CreateApp(app *meta.App, context string) error
	DeleteApp(query string) error
	UpdateApp(app *meta.App, query string) error
}

// ChannelTypeMemory is the interface that allows to
// obtain or change information related to the current
// state of the ChannelTypes in the cluster
type ChannelTypeMemory interface {
	GetChannelType(context string, ctName string) (*meta.ChannelType, error)
	CreateChannelType(ct *meta.ChannelType, context string) error
	DeleteChannelType(context string, ctName string) error
	UpdateChannelType(ct *meta.ChannelType, context string) error
}

// Manager is the interface that allows the management
// of the current state of the cluster. Permiting the
// modification of Channels, DApps and ChannelTypes
type Manager interface {
	Apps() AppMemory
	Cancel()
	Channels() ChannelMemory
	ChannelTypes() ChannelTypeMemory
	Commit()
	GetTransactionChanges() (diff.Changelog, error)
	InitTransaction()
}
