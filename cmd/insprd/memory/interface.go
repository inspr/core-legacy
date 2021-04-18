// Package memory provides simple interfaces for the
// in memory management of the cluster.
package memory

import (
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils/diff"
)

// ChannelMemory is the interface that allows to obtain
// or change information related to the stored state of
// the Channels in the cluster
type ChannelMemory interface {
	TransactionInterface
	ChannelGetInterface
	Create(context string, ch *meta.Channel) error
	Delete(context string, chName string) error
	Update(context string, ch *meta.Channel) error
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
	Create(context string, app *meta.App) error
	Delete(query string) error
	Update(query string, app *meta.App) error
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
	Create(context string, ct *meta.ChannelType) error
	Delete(context string, ctName string) error
	Update(context string, ct *meta.ChannelType) error
}

// ChannelTypeGetInterface is an interface to get channel types from memory
type ChannelTypeGetInterface interface {
	Get(context string, ctName string) (*meta.ChannelType, error)
}

// AliasMemory is an interface to get alias types from memory
type AliasMemory interface {
	TransactionInterface
	AliasGetInterface
	Create(query string, targetBoundary string, alias *meta.Alias) error
	Update(context string, aliasKey string, alias *meta.Alias) error
	Delete(context string, aliasKey string) error
}

// AliasGetInterface is an interface to get alias types from memory
type AliasGetInterface interface {
	Get(context string, aliasKey string) (*meta.Alias, error)
}

// Manager is the interface that allows the management
// of the current state of the cluster. Permiting the
// modification of Channels, DApps and ChannelTypes
type Manager interface {
	TransactionInterface
	Apps() AppMemory
	Channels() ChannelMemory
	ChannelTypes() ChannelTypeMemory
	Alias() AliasMemory
	Root() GetInterface
}

// GetInterface is an interface to get components from memory
type GetInterface interface {
	Apps() AppGetInterface
	Channels() ChannelGetInterface
	ChannelTypes() ChannelTypeGetInterface
	Alias() AliasGetInterface
}

// TransactionInterface makes transactions on a Memory manager
type TransactionInterface interface {
	Commit()
	GetTransactionChanges() (diff.Changelog, error)
	InitTransaction()
	Cancel()
}
