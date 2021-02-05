package client

import (
	"context"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

type AppClient struct {
	c *rest.Client
}

func (cc *AppClient) GetApp(ctx context.Context, context string, chName string) (*meta.App, error) {
	cdi := models.AppQueryDI{
		Query: context,
		Valid: true,
	}

	var resp meta.App

	err := cc.c.SendRequest(ctx, "/channel", "GET", cdi, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (cc *AppClient) CreateApp(ctx context.Context, context string, ch *meta.App) error {
	cdi := models.AppDI{
		Ctx:   context,
		App:   *ch,
		Valid: true,
	}

	err := cc.c.SendRequest(ctx, "/channel", "POST", cdi, nil)
	if err != nil {
		return err
	}

	return nil
}

func (cc *AppClient) DeleteApp(ctx context.Context, context string, chName string) error {
	cdi := models.AppQueryDI{
		Query: context,
		Valid: true,
	}

	err := cc.c.SendRequest(ctx, "/channel", "DELETE", cdi, nil)
	if err != nil {
		return err
	}

	return nil
}

func (cc *AppClient) UpdateApp(ctx context.Context, context string, ch *meta.App) error {
	cdi := models.AppDI{
		Ctx:   context,
		App:   *ch,
		Valid: true,
	}

	err := cc.c.SendRequest(ctx, "/channel", "PUT", cdi, nil)
	if err != nil {
		return err
	}

	return nil
}
