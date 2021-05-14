package controller

import (
	"context"

	"github.com/inspr/inspr/pkg/auth"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/utils/diff"
)

// ChannelInterface is the interface that allows to obtain
// or change information related to the stored state of
// the Channels in the cluster
type ChannelInterface interface {
	Get(ctx context.Context, context string, chName string) (*meta.Channel, error)
	Create(ctx context.Context, context string, ch *meta.Channel, dryRun bool) (diff.Changelog, error)
	Delete(ctx context.Context, context string, chName string, dryRun bool) (diff.Changelog, error)
	Update(ctx context.Context, context string, ch *meta.Channel, dryRun bool) (diff.Changelog, error)
}

// AppInterface is the interface that allows to obtain or
// change information related to the current state of
// the DApps in the cluster
type AppInterface interface {
	Get(ctx context.Context, query string) (*meta.App, error)
	Create(ctx context.Context, context string, app *meta.App, dryRun bool) (diff.Changelog, error)
	Delete(ctx context.Context, query string, dryRun bool) (diff.Changelog, error)
	Update(ctx context.Context, query string, app *meta.App, dryRun bool) (diff.Changelog, error)
}

// TypeInterface is the interface that allows to
// obtain or change information related to the current
// state of the Types in the cluster
type TypeInterface interface {
	Get(ctx context.Context, context string, ctName string) (*meta.Type, error)
	Create(ctx context.Context, context string, ct *meta.Type, dryRun bool) (diff.Changelog, error)
	Delete(ctx context.Context, context string, ctName string, dryRun bool) (diff.Changelog, error)
	Update(ctx context.Context, context string, ct *meta.Type, dryRun bool) (diff.Changelog, error)
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
	Get(ctx context.Context, context, key string) (*meta.Alias, error)
	Create(ctx context.Context, context string, target string, alias *meta.Alias, dryRun bool) (diff.Changelog, error)
	Delete(ctx context.Context, context, key string, dryRun bool) (diff.Changelog, error)
	Update(ctx context.Context, context string, target string, alias *meta.Alias, dryRun bool) (diff.Changelog, error)
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
}
