package client

import (
	"context"

	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/utils/diff"
	"github.com/inspr/inspr/pkg/rest"
	"github.com/inspr/inspr/pkg/rest/request"
)

// TypeClient interacts with types on the Insprd
type TypeClient struct {
	reqClient *request.Client
}

// Get gets a Type from the Insprd
//
// The scope refers to the parent app of the given channel type, represented with a dot separated query
// such as app1.app2
//
// The name is the name of the channel type. So to search for a channel type inside app1 with the name type1 you
// would call tc.Get(context.Background(), "app1", "type1")
func (tc *TypeClient) Get(ctx context.Context, scope string, name string) (*meta.Type, error) {
	tdi := models.TypeQueryDI{
		TypeName: name,
	}
	var resp meta.Type

	err := tc.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/types", "GET", tdi, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Create creates a Type inside the Insprd
//
// The scope refers to the parent app of the given channel type, represented with a dot separated query
// such as **app1.app2**
//
// The Type information such as name and etc will be inferred from the given Type's metadata.
//
// So to create a channel type inside app1 with the name type1 you
// would call tc.Create(context.Background(), "app1", &meta.Type{...})
func (tc *TypeClient) Create(ctx context.Context, scope string, ch *meta.Type, dryRun bool) (diff.Changelog, error) {
	tdi := models.TypeDI{
		Type:   *ch,
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := tc.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/types", "POST", tdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a Type inside the Insprd
//
// The scope refers to the parent app of the given channel type, represented with a dot separated query
// such as **app1.app2**
//
// The name is the name of the Type to be deleted.
//
// So to delete a channel type inside app1 with the name type1 you
// would call tc.Delete(context.Background(), "app1", "type1")
func (tc *TypeClient) Delete(ctx context.Context, scope string, name string, dryRun bool) (diff.Changelog, error) {
	tdi := models.TypeQueryDI{
		TypeName: name,
		DryRun:   dryRun,
	}
	var resp diff.Changelog

	err := tc.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/types", "DELETE", tdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update updates a Type inside the Insprd
//
// The scope refers to the parent app of the given channel type, represented with a dot separated query
// such as **app1.app2**
//
// The Type information such as name and etc will be inferred from the given Type's metadata.
//
// So to update a channel type inside app1 with the name type1 you
// would call tc.Create(context.Background(), "app1", &meta.Type{...})
func (tc *TypeClient) Update(ctx context.Context, scope string, ch *meta.Type, dryRun bool) (diff.Changelog, error) {
	tdi := models.TypeDI{
		Type:   *ch,
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := tc.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/types", "PUT", tdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
