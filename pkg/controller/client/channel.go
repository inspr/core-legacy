package client

import (
	"context"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
	"gitlab.inspr.dev/inspr/core/pkg/utils/diff"
)

// ChannelClient interacts with channels on the Insprd
type ChannelClient struct {
	c *request.Client
}

// Get gets a channel from the Insprd
//
// The context refers to the parent app of the given channel, represented with a dot separated query
// such as app1.app2
//
// The name is the name of the channel. So to search for a channel inside app1 with the name channel1 you
// would call cc.Get(context.Background(), "app1", "channel1")
func (cc *ChannelClient) Get(ctx context.Context, context string, name string) (*meta.Channel, error) {
	cdi := models.ChannelQueryDI{
		Ctx:    context,
		ChName: name,
		Valid:  true,
	}

	var resp meta.Channel

	err := cc.c.Send(ctx, "/channels", "GET", cdi, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Create creates a channel inside the Insprd
//
// The context refers to the parent app of the given channel, represented with a dot separated query
// such as **app1.app2**
//
// The channel information such as name and etc will be inferred from the given channel's metadata.
//
// So to create a channel inside app1 with the name channel1 you
// would call cc.Create(context.Background(), "app1", &meta.Channel{...})
func (cc *ChannelClient) Create(ctx context.Context, context string, ch *meta.Channel) (diff.Changelog, error) {
	cdi := models.ChannelDI{
		Ctx:     context,
		Channel: *ch,
		Valid:   true,
	}

	var resp diff.Changelog
	err := cc.c.Send(ctx, "/channels", "POST", cdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a channel inside the Insprd
//
// The context refers to the parent app of the given channel, represented with a dot separated query
// such as **app1.app2**
//
// The name is the name of the channel to be deleted.
//
// So to delete a channel inside app1 with the name channel1 you
// would call cc.Delete(context.Background(), "app1", "channel1")
func (cc *ChannelClient) Delete(ctx context.Context, context string, name string) (diff.Changelog, error) {
	cdi := models.ChannelQueryDI{
		Ctx:    context,
		ChName: name,
		Valid:  true,
	}

	var resp diff.Changelog
	err := cc.c.Send(ctx, "/channels", "DELETE", cdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update creates a channel inside the Insprd
//
// The context refers to the parent app of the given channel, represented with a dot separated query
// such as **app1.app2**
//
// The channel information such as name and etc will be inferred from the given channel's metadata.
//
// So to update a channel inside app1 with the name channel1 you
// would call cc.Update(context.Background(), "app1", &meta.Channel{...})
func (cc *ChannelClient) Update(ctx context.Context, context string, ch *meta.Channel) (diff.Changelog, error) {
	cdi := models.ChannelDI{
		Ctx:     context,
		Channel: *ch,
		Valid:   true,
	}

	var resp diff.Changelog
	err := cc.c.Send(ctx, "/channels", "PUT", cdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
