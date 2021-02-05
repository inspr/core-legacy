package client

import (
	"context"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

type ChannelClient struct {
	c *rest.Client
}

func (cc *ChannelClient) GetChannel(ctx context.Context, context string, chName string) (*meta.Channel, error) {
	cdi := models.ChannelQueryDI{
		Ctx:    context,
		ChName: chName,
		Valid:  true,
	}

	var resp meta.Channel

	err := cc.c.SendRequest(ctx, "/channel", "GET", cdi, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (cc *ChannelClient) CreateChannel(ctx context.Context, context string, ch *meta.Channel) error {
	cdi := models.ChannelDI{
		Ctx:     context,
		Channel: *ch,
		Valid:   true,
	}

	err := cc.c.SendRequest(ctx, "/channel", "POST", cdi, nil)
	if err != nil {
		return err
	}

	return nil
}

func (cc *ChannelClient) DeleteChannel(ctx context.Context, context string, chName string) error {
	cdi := models.ChannelQueryDI{
		Ctx:    context,
		ChName: chName,
		Valid:  true,
	}

	err := cc.c.SendRequest(ctx, "/channel", "DELETE", cdi, nil)
	if err != nil {
		return err
	}

	return nil
}

func (cc *ChannelClient) UpdateChannel(ctx context.Context, context string, ch *meta.Channel) error {
	cdi := models.ChannelDI{
		Ctx:     context,
		Channel: *ch,
		Valid:   true,
	}

	err := cc.c.SendRequest(ctx, "/channel", "PUT", cdi, nil)
	if err != nil {
		return err
	}

	return nil
}
