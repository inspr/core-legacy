// Package memory provides simple interfaces for the
// in memory management of the cluster.
package memory

import (
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
)

// ChannelMemory is the interface that allows to obtain
// or change information related to the stored state of
// the Channels in the cluster
type ChannelMemory interface {
	TransactionInterface
	ChannelGetInterface
	CreateChannel(context string, ch *meta.Channel) error
	DeleteChannel(context string, chName string) error
	UpdateChannel(context string, ch *meta.Channel) error
}

// ChannelGetInterface is an interface to get channels from memory
type ChannelGetInterface interface {
	Get(context string, ctName string) (*meta.Channel, error)
}

// AppMemory is the interface that allows to obtain or
// change information related to the current state of
// the DApps in the cluster
type AppMemory interface {
	TransactionInterface
	AppGetInterface
	CreateApp(context string, app *meta.App) error
	DeleteApp(query string) error
	UpdateApp(query string, app *meta.App) error
	ResolveBoundary(app *meta.App) (map[string]string, error)
}

// AppGetInterface is an interface to get apps from memory
type AppGetInterface interface {
	Get(query string) (*meta.App, error)
}

// ChannelTypeMemory is the interface that allows to
// obtain or change information related to the current
// state of the ChannelTypes in the cluster
type ChannelTypeMemory interface {
	TransactionInterface
	ChannelTypeGetInterface
	CreateChannelType(context string, ct *meta.ChannelType) error
	DeleteChannelType(context string, ctName string) error
	UpdateChannelType(context string, ct *meta.ChannelType) error
}

// ChannelTypeGetInterface is an interface to get channel types from memory
type ChannelTypeGetInterface interface {
	Get(context string, ctName string) (*meta.ChannelType, error)
}

type AliasMemory interface {
	TransactionInterface
	AliasGetInterface
	CreateAlias(context string, targetBoundary string, targetChannel string) error
}

// Manager is the interface that allows the management
// of the current state of the cluster. Permiting the
// modification of Channels, DApps and ChannelTypes
type Manager interface {
	TransactionInterface
	Apps() AppMemory
	Channels() ChannelMemory
	ChannelTypes() ChannelTypeMemory
	Root() GetInterface
}

// GetInterface is an interface to get components from memory
type GetInterface interface {
	Apps() AppGetInterface
	Channels() ChannelGetInterface
	ChannelTypes() ChannelTypeGetInterface
}

// TransactionInterface makes transactions on a Memory manager
type TransactionInterface interface {
	Commit()
	GetTransactionChanges() (diff.Changelog, error)
	InitTransaction()
	Cancel()
}
