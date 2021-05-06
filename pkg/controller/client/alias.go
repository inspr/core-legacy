package client

import (
	"context"

	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/utils/diff"
	"github.com/inspr/inspr/pkg/rest"
	"github.com/inspr/inspr/pkg/rest/request"
)

// AliasClient interacts with Aliases on the Insprd
type AliasClient struct {
	reqClient *request.Client
}

// Get gets a alias from the Insprd
//
// The scope refers to the parent app of the given alias, represented with a dot separated query
// such as app1.app2
//
// The key is the key of the alias. So to search for a alias inside app1 with the key aliasKey you
// would call ac.Get(context.Background(), "app1", "aliasKey")
func (ac *AliasClient) Get(ctx context.Context, scope, key string) (*meta.Alias, error) {
	aliasQuery := models.AliasQueryDI{
		Key: key,
	}

	var resp meta.Alias

	err := ac.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/alias", "GET", aliasQuery, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Create creates a alias inside the Insprd
//
// The scope refers to the parent app of the given alias, represented with a dot separated query
// such as **app1.app2**
//
// The alias information such as name and etc will be inferred from the given alias metadata.
//
// So to create a alias inside app1 with the name aliasOne you
// would call ctc.Create(context.Background(), "app1", &meta.Alias{...})
func (ac *AliasClient) Create(ctx context.Context, scope string, target string, alias *meta.Alias, dryRun bool) (diff.Changelog, error) {
	aliasQuery := models.AliasDI{
		Target: target,
		Alias:  *alias,
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := ac.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/alias", "POST", aliasQuery, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a alias inside the Insprd
//
// The scope refers to the parent app of the given alias, represented with a dot separated query
// such as **app1.app2**
//
// The key is the key of the alias to be deleted.
//
// So to delete a alias inside app1 with the key alias1 you
// would call ac.Delete(context.Background(), "app1", "alias1")
func (ac *AliasClient) Delete(ctx context.Context, scope, key string, dryRun bool) (diff.Changelog, error) {
	aliasQuery := models.AliasQueryDI{
		Key:    key,
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := ac.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/alias", "DELETE", aliasQuery, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update updates a alias inside the Insprd
//
// The scope refers to the parent app of the given alias, represented with a dot separated query
// such as **app1.app2**
//
// The alias information will be inferred from the given alias metadata.
//
// So to update a alias inside app1 with the key myalias you
// would call ac.Create(context.Background(), "app1", &meta.Alias{...})
func (ac *AliasClient) Update(ctx context.Context, scope string, target string, alias *meta.Alias, dryRun bool) (diff.Changelog, error) {
	aliasQuery := models.AliasDI{
		Target: target,
		Alias:  *alias,
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := ac.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/alias", "PUT", aliasQuery, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
