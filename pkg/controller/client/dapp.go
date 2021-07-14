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

// AppClient is a client for manipulating dApp structures in Insprd
type AppClient struct {
	reqClient *request.Client
}

// Get gets information of a dApp that exists in Insprd.
// The scope refers to the app itself, represented with a dot separated query
// such as app1.app2
func (ac *AppClient) Get(ctx context.Context, scope string) (*meta.App, error) {
	adi := models.AppQueryDI{}
	var resp meta.App

	err := ac.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/apps", http.MethodGet, request.DefaultHost, adi, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Create creates given dApp in Insprd
// The scope refers to the parent app where the actual app will be instantiated,
// represented with a dot separated query such as **app1.app2**.
func (ac *AppClient) Create(ctx context.Context, scope string, app *meta.App, dryRun bool) (diff.Changelog, error) {
	adi := models.AppDI{
		App:    *app,
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := ac.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/apps", http.MethodPost, request.DefaultHost, adi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a dApp that exists in Insprd.
// The scope refers to the app itself, represented with a dot separated query
// such as app1.app2
func (ac *AppClient) Delete(ctx context.Context, scope string, dryRun bool) (diff.Changelog, error) {
	adi := models.AppQueryDI{
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := ac.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/apps", http.MethodDelete, request.DefaultHost, adi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update updates a dApp in Insprd, if it exists.
// The scope refers to the parent dApp where the actual dApp is instantiated,
// represented with a dot separated query such as app1.app2
func (ac *AppClient) Update(ctx context.Context, scope string, app *meta.App, dryRun bool) (diff.Changelog, error) {
	adi := models.AppDI{
		App:    *app,
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := ac.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/apps", http.MethodPut, request.DefaultHost, adi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
