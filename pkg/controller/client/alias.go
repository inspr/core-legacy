package client

import (
	"context"
	"net/http"

	"inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils/diff"
	"inspr.dev/inspr/pkg/rest"
	"inspr.dev/inspr/pkg/rest/request"
)

// AliasClient is a client for manipulating Alias structures in Insprd
type AliasClient struct {
	reqClient *request.Client
}

// Get gets a alias from the Insprd
// The scope refers to the app of the given alias, represented with a dot separated query
// such as app1.app2. The name is the name of the alias
func (ac *AliasClient) Get(ctx context.Context, scope, name string) (*meta.Alias, error) {
	aliasQuery := models.AliasQueryDI{
		Name: name,
	}

	var resp meta.Alias

	err := ac.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/alias", http.MethodGet, aliasQuery, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Create creates an alias inside the Insprd
// The scope refers to the app of the given alias, represented with a dot separated query
// such as app1.app2.
func (ac *AliasClient) Create(ctx context.Context, scope string, alias *meta.Alias, dryRun bool) (diff.Changelog, error) {
	aliasQuery := models.AliasDI{
		Alias:  *alias,
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := ac.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/alias", http.MethodPost, aliasQuery, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a alias inside the Insprd
// The scope refers to the app of the given alias, represented with a dot separated query
// such as app1.app2. The name is the name of the alias to be deleted.
func (ac *AliasClient) Delete(ctx context.Context, scope, name string, dryRun bool) (diff.Changelog, error) {
	aliasQuery := models.AliasQueryDI{
		Name:   name,
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := ac.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/alias", http.MethodDelete, aliasQuery, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update updates a alias inside the Insprd
// The scope refers to the app of the given alias, represented with a dot separated query
// such as app1.app2. Works similarly to the Create method.
func (ac *AliasClient) Update(ctx context.Context, scope string, alias *meta.Alias, dryRun bool) (diff.Changelog, error) {
	aliasQuery := models.AliasDI{
		Alias:  *alias,
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := ac.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/alias", http.MethodPut, aliasQuery, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
