// Package memory provides simple interfaces for the
// in memory management of the cluster.
package memory

import (
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/brokers"
	"github.com/inspr/inspr/pkg/meta/utils/diff"
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

// TypeMemory is the interface that allows to
// obtain or change information related to the current
// state of the Types in the cluster
type TypeMemory interface {
	TransactionInterface
	TypeGetInterface
	Create(context string, ct *meta.Type) error
	Delete(context string, ctName string) error
	Update(context string, ct *meta.Type) error
}

// TypeGetInterface is an interface to get Types from memory
type TypeGetInterface interface {
	Get(context string, ctName string) (*meta.Type, error)
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
// modification of Channels, DApps and Types
type Manager interface {
	TransactionInterface
	Apps() AppMemory
	Channels() ChannelMemory
	Types() TypeMemory
	Alias() AliasMemory
	Root() GetInterface
	Brokers() BrokerInterface
}

// GetInterface is an interface to get components from memory
type GetInterface interface {
	Apps() AppGetInterface
	Channels() ChannelGetInterface
	Types() TypeGetInterface
	Alias() AliasGetInterface
}

// TransactionInterface makes transactions on a Memory manager
type TransactionInterface interface {
	Commit()
	GetTransactionChanges() (diff.Changelog, error)
	InitTransaction()
	Cancel()
}

// BrokerInterface is the interface tht allows for interaction
// with the systems multiple brokers
type BrokerInterface interface {
	GetAll() brokers.BrokerStatusArray
	GetDefault() brokers.BrokerStatus
	Create(broker brokers.BrokerStatus, config brokers.BrokerConfiguration) error
	SetDefault(broker brokers.BrokerStatus) error
}
