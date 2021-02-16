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
	TransactionInterface
	GetChannel(context string, chName string) (*meta.Channel, error)
	CreateChannel(context string, ch *meta.Channel) error
	DeleteChannel(context string, chName string) error
	UpdateChannel(context string, ch *meta.Channel) error
}

// AppMemory is the interface that allows to obtain or
// change information related to the current state of
// the DApps in the cluster
type AppMemory interface {
	TransactionInterface
	GetApp(query string) (*meta.App, error)
	CreateApp(context string, app *meta.App) error
	DeleteApp(query string) error
	UpdateApp(query string, app *meta.App) error
}

// ChannelTypeMemory is the interface that allows to
// obtain or change information related to the current
// state of the ChannelTypes in the cluster
type ChannelTypeMemory interface {
	TransactionInterface
	GetChannelType(context string, ctName string) (*meta.ChannelType, error)
	CreateChannelType(context string, ct *meta.ChannelType) error
	DeleteChannelType(context string, ctName string) error
	UpdateChannelType(context string, ct *meta.ChannelType) error
}

// Manager is the interface that allows the management
// of the current state of the cluster. Permiting the
// modification of Channels, DApps and ChannelTypes
type Manager interface {
	TransactionInterface
	Apps() AppMemory
	Channels() ChannelMemory
	ChannelTypes() ChannelTypeMemory
}

// TransactionInterface makes transactions on a Memory manager
type TransactionInterface interface {
	Commit()
	GetTransactionChanges() (diff.Changelog, error)
	InitTransaction()
	Cancel()
}
