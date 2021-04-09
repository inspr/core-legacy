package client

import (
	"context"

	"gitlab.inspr.dev/inspr/core/pkg/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

// AppClient is a client for getting and setting app information on Insprd
type AppClient struct {
	c *request.Client
}

// Get gets information from an app inside the Insprd
//
// The context refers to the app itself, represented with a dot separated query
// such as **app1.app2**.
//
// So to get an app inside app1 with the name app2 you
// would call ac.Get(context.Background(), "app1.app2")
func (ac *AppClient) Get(ctx context.Context, context string) (*meta.App, error) {
	adi := models.AppQueryDI{
		Ctx: context,
	}

	var resp meta.App

	err := ac.c.Send(ctx, "/apps", "GET", adi, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Create creates an app inside the Insprd
//
// The context refers to the parent app where the actual app will be instantiated,
// represented with a dot separated query such as **app1.app2**.
//
// The information of the app such as name and other metadata will be gotten from the
// definition of the app itself.
//
// So to create an app inside app1 with the name app2 you
// would call ac.Create(context.Background(), "app1", &meta.App{...})
func (ac *AppClient) Create(ctx context.Context, context string, app *meta.App, dryRun bool) (diff.Changelog, error) {
	adi := models.AppDI{
		Ctx:    context,
		App:    *app,
		DryRun: dryRun,
	}
	var resp diff.Changelog
	err := ac.c.Send(ctx, "/apps", "POST", adi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes an app inside the Insprd.
//
// The context refers to the app itself, represented with a dot separated query
// such as **app1.app2**.
//
// So to delete an app inside app1 with the name app2 you
// would call ac.Delete(context.Background(), "app1.app2")
func (ac *AppClient) Delete(ctx context.Context, context string, dryRun bool) (diff.Changelog, error) {
	adi := models.AppQueryDI{
		Ctx:    context,
		DryRun: dryRun,
	}
	var resp diff.Changelog
	err := ac.c.Send(ctx, "/apps", "DELETE", adi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update updates an app inside the Insprd.
//
// The context refers to the parent app where the actual app will be instantiated,
// represented with a dot separated query such as **app1.app2**.
//
// The information of the app such as name and other metadata will be gotten from the
// definition of the app itself.
//
// So to update an app inside app1 with the name app2 you
// would call ac.Update(context.Background(), "app1", &meta.App{...})
func (ac *AppClient) Update(ctx context.Context, context string, app *meta.App, dryRun bool) (diff.Changelog, error) {
	adi := models.AppDI{
		Ctx:    context,
		App:    *app,
		DryRun: dryRun,
	}

	var resp diff.Changelog
	err := ac.c.Send(ctx, "/apps", "PUT", adi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
