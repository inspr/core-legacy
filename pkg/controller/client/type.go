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

// TypeClient is a client for manipulating Type structures in Insprd
type TypeClient struct {
	reqClient *request.Client
}

// Get gets a Type from Insprd, if it exists.
// The scope refers to the dApp in which the Type is in, represented with a dot separated query
// such as app1.app2
func (tc *TypeClient) Get(ctx context.Context, scope, name string) (*meta.Type, error) {
	tdi := models.TypeQueryDI{
		TypeName: name,
	}
	var resp meta.Type

	err := tc.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/types", http.MethodGet, request.DefaultHost, tdi, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Create creates given Type inside of Insprd
// The scope refers to the dApp in which the Type will be in, represented with a dot separated query
// such as app1.app2
func (tc *TypeClient) Create(ctx context.Context, scope string, t *meta.Type, dryRun bool) (diff.Changelog, error) {
	tdi := models.TypeDI{
		Type:   *t,
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := tc.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/types", http.MethodPost, request.DefaultHost, tdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a Type from Insprd, if exists.
// The scope refers  to the dApp in which the Type is in, represented with a dot separated query
// such as app1.app2. The name is the name of the Type to be deleted
func (tc *TypeClient) Delete(ctx context.Context, scope, name string, dryRun bool) (diff.Changelog, error) {
	tdi := models.TypeQueryDI{
		TypeName: name,
		DryRun:   dryRun,
	}
	var resp diff.Changelog

	err := tc.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/types", http.MethodDelete, request.DefaultHost, tdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update updates given Type in Insprd, if it exists.
// The scope refers to the dApp in which the Type is in, represented with a dot separated query
// such asapp1.app2
func (tc *TypeClient) Update(ctx context.Context, scope string, t *meta.Type, dryRun bool) (diff.Changelog, error) {
	tdi := models.TypeDI{
		Type:   *t,
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := tc.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/types", http.MethodPut, request.DefaultHost, tdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
