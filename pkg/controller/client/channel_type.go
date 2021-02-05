package client

import (
	"context"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

type ChannelTypeClient struct {
	c *rest.Client
}

func (ctc *ChannelTypeClient) GetChannelType(ctx context.Context, context string, ctName string) (*meta.ChannelType, error) {
	panic("not implemented") // TODO: Implement
}

func (ctc *ChannelTypeClient) CreateChannelType(ctx context.Context, ct *meta.ChannelType, context string) error {
	panic("not implemented") // TODO: Implement
}

func (ctc *ChannelTypeClient) DeleteChannelType(ctx context.Context, context string, ctName string) error {
	panic("not implemented") // TODO: Implement
}

func (ctc *ChannelTypeClient) UpdateChannelType(ctx context.Context, ct *meta.ChannelType, context string) error {
	panic("not implemented") // TODO: Implement
}
