package controller

import (
	"context"

	"inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/auth"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils/diff"
)

// ChannelInterface is the interface that allows to obtain
// or change information related to the stored state of
// the Channels in the cluster
type ChannelInterface interface {
	Get(ctx context.Context, scope, name string) (*meta.Channel, error)
	Create(ctx context.Context, scope string, ch *meta.Channel, dryRun bool) (diff.Changelog, error)
	Delete(ctx context.Context, scope, name string, dryRun bool) (diff.Changelog, error)
	Update(ctx context.Context, scope string, ch *meta.Channel, dryRun bool) (diff.Changelog, error)
}

// AppInterface is the interface that allows to obtain or
// change information related to the current state of
// the DApps in the cluster
type AppInterface interface {
	Get(ctx context.Context, scope string) (*meta.App, error)
	Create(ctx context.Context, scope string, app *meta.App, dryRun bool) (diff.Changelog, error)
	Delete(ctx context.Context, scope string, dryRun bool) (diff.Changelog, error)
	Update(ctx context.Context, scope string, app *meta.App, dryRun bool) (diff.Changelog, error)
}

// TypeInterface is the interface that allows to
// obtain or change information related to the current
// state of the Types in the cluster
type TypeInterface interface {
	Get(ctx context.Context, scope, name string) (*meta.Type, error)
	Create(ctx context.Context, scope string, t *meta.Type, dryRun bool) (diff.Changelog, error)
	Delete(ctx context.Context, scope, name string, dryRun bool) (diff.Changelog, error)
	Update(ctx context.Context, scope string, t *meta.Type, dryRun bool) (diff.Changelog, error)
}

// AuthorizationInterface is the interface that allows to
// obtain information related to the authorization necessary
// to make changes in structures inside of the cluster
type AuthorizationInterface interface {
	GenerateToken(ctx context.Context, payload auth.Payload) (string, error)
	Init(ctx context.Context, key string) (string, error)
}

// AliasInterface is the interface that allows to
// obtain or change information related to the current
// state of the Alias in the cluster
type AliasInterface interface {
	Get(ctx context.Context, scope, key string) (*meta.Alias, error)
	Create(ctx context.Context, scope string, alias *meta.Alias, dryRun bool) (diff.Changelog, error)
	Delete(ctx context.Context, scope, name string, dryRun bool) (diff.Changelog, error)
	Update(ctx context.Context, scope string, alias *meta.Alias, dryRun bool) (diff.Changelog, error)
}

// BrokersInterface is the interface that allows to
// obtain or change information related to the current
// cluster's message brokers.
type BrokersInterface interface {
	Get(ctx context.Context) (*models.BrokersDI, error)
	Create(ctx context.Context, name string, config []byte) error
}

// Interface is the interface that allows the management
// of the current state of the cluster. Permiting the
// modification of Channels, DApps and Types
type Interface interface {
	Channels() ChannelInterface
	Apps() AppInterface
	Types() TypeInterface
	Authorization() AuthorizationInterface
	Alias() AliasInterface
	Brokers() BrokersInterface
}
