package tree

import (
	apimodels "inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils/diff"
)

// ChannelMemory is the interface that allows to obtain
// or change information related to the stored state of
// the Channels in the cluster
type ChannelMemory interface {
	TransactionInterface
	ChannelGetInterface
	Create(scope string, ch *meta.Channel, brokers *apimodels.BrokersDI) error
	Delete(scope, name string) error
	Update(scope string, ch *meta.Channel) error
}

// ChannelGetInterface is an interface to get channels from memory
type ChannelGetInterface interface {
	Get(scope, name string) (*meta.Channel, error)
}

// AppMemory is the interface that allows to obtain or
// change information related to the current state of
// the DApps in the cluster
type AppMemory interface {
	TransactionInterface
	AppGetInterface
	Create(scope string, app *meta.App, brokers *apimodels.BrokersDI) error
	Delete(query string) error
	Update(query string, app *meta.App, brokers *apimodels.BrokersDI) error
	ResolveBoundary(app *meta.App, usePermTree bool) (map[string]string, error)
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
	Create(scope string, ct *meta.Type) error
	Delete(scope, name string) error
	Update(scope string, ct *meta.Type) error
}

// TypeGetInterface is an interface to get Types from memory
type TypeGetInterface interface {
	Get(scope, name string) (*meta.Type, error)
}

// AliasMemory is an interface to get alias types from memory
type AliasMemory interface {
	TransactionInterface
	AliasGetInterface
	Create(scope string, alias *meta.Alias) error
	Delete(scope, name string) error
	Update(scope string, alias *meta.Alias) error
	CheckSource(scope string, app *meta.App, alias *meta.Alias) error
	CheckDestination(app *meta.App, alias *meta.Alias) error
}

// AliasGetInterface is an interface to get alias types from memory
type AliasGetInterface interface {
	Get(scope, name string) (*meta.Alias, error)
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
	Perm() GetInterface
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
