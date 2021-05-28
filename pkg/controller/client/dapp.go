package client

import (
	"context"
	"net/http"

	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/utils/diff"
	"github.com/inspr/inspr/pkg/rest"
	"github.com/inspr/inspr/pkg/rest/request"
)

// AppClient is a client for getting and setting app information on Insprd
type AppClient struct {
	reqClient *request.Client
}

// Get gets information from an app inside the Insprd
//
// The scope refers to the app itself, represented with a dot separated query
// such as **app1.app2**.
//
// So to get an app inside app1 with the name app2 you
// would call ac.Get(context.Background(), "app1.app2")
func (ac *AppClient) Get(ctx context.Context, scope string) (*meta.App, error) {
	adi := models.AppQueryDI{}
	var resp meta.App

	err := ac.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/apps", "GET", adi, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Create creates an app inside the Insprd
//
// The scope refers to the parent app where the actual app will be instantiated,
// represented with a dot separated query such as **app1.app2**.
//
// The information of the app such as name and other metadata will be gotten from the
// definition of the app itself.
//
// So to create an app inside app1 with the name app2 you
// would call ac.Create(context.Background(), "app1", &meta.App{...})
func (ac *AppClient) Create(ctx context.Context, scope string, app *meta.App, dryRun bool) (diff.Changelog, error) {
	adi := models.AppDI{
		App:    *app,
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := ac.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/apps", http.MethodPost, adi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes an app inside the Insprd.
//
// The scope refers to the app itself, represented with a dot separated query
// such as **app1.app2**.
//
// So to delete an app inside app1 with the name app2 you
// would call ac.Delete(context.Background(), "app1.app2")
func (ac *AppClient) Delete(ctx context.Context, scope string, dryRun bool) (diff.Changelog, error) {
	adi := models.AppQueryDI{
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := ac.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/apps", "DELETE", adi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update updates an app inside the Insprd.
//
// The scope refers to the parent app where the actual app will be instantiated,
// represented with a dot separated query such as **app1.app2**.
//
// The information of the app such as name and other metadata will be gotten from the
// definition of the app itself.
//
// So to update an app inside app1 with the name app2 you
// would call ac.Update(context.Background(), "app1", &meta.App{...})
func (ac *AppClient) Update(ctx context.Context, scope string, app *meta.App, dryRun bool) (diff.Changelog, error) {
	adi := models.AppDI{
		App:    *app,
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := ac.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/apps", "PUT", adi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
