package client

import (
	"context"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

// ChannelTypeClient interacts with channels on the Insprd
type ChannelTypeClient struct {
	c *request.Client
}

// Get gets a channel type from the Insprd
//
// The context refers to the parent app of the given channel type, represented with a dot separated query
// such as app1.app2
//
// The name is the name of the channel type. So to search for a channel type inside app1 with the name channel1 you
// would call cc.Get(context.Background(), "app1", "channel1")
func (cc *ChannelTypeClient) Get(ctx context.Context, context string, name string) (*meta.ChannelType, error) {
	cdi := models.ChannelTypeQueryDI{
		Ctx:    context,
		CtName: name,
		Valid:  true,
	}

	var resp meta.ChannelType

	err := cc.c.Send(ctx, "/channels", "GET", cdi, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Create creates a channel type inside the Insprd
//
// The context refers to the parent app of the given channel type, represented with a dot separated query
// such as **app1.app2**
//
// The channel type information such as name and etc will be inferred from the given channel type's metadata.
//
// So to create a channel type inside app1 with the name channel1 you
// would call cc.Create(context.Background(), "app1", &meta.{...})
func (cc *ChannelTypeClient) Create(ctx context.Context, context string, ch *meta.ChannelType) error {
	cdi := models.ChannelTypeDI{
		Ctx:         context,
		ChannelType: *ch,
		Valid:       true,
	}

	var resp interface{}
	err := cc.c.Send(ctx, "/channels", "POST", cdi, &resp)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a channel type inside the Insprd
//
// The context refers to the parent app of the given channel type, represented with a dot separated query
// such as **app1.app2**
//
// The name is the name of the channel type to be deleted.
//
// So to delete a channel type inside app1 with the name channel1 you
// would call cc.Delete(context.Background(), "app1", "channel1")
func (cc *ChannelTypeClient) Delete(ctx context.Context, context string, name string) error {
	cdi := models.ChannelTypeQueryDI{
		Ctx:    context,
		CtName: name,
		Valid:  true,
	}

	var resp interface{}
	err := cc.c.Send(ctx, "/channels", "DELETE", cdi, &resp)
	if err != nil {
		return err
	}

	return nil
}

// Update creates a channel type inside the Insprd
//
// The context refers to the parent app of the given channel type, represented with a dot separated query
// such as **app1.app2**
//
// The channel type information such as name and etc will be inferred from the given channel type's metadata.
//
// So to update a channel type inside app1 with the name channel1 you
// would call cc.Create(context.Background(), "app1", &meta.{...})
func (cc *ChannelTypeClient) Update(ctx context.Context, context string, ch *meta.ChannelType) error {
	cdi := models.ChannelTypeDI{
		Ctx:         context,
		ChannelType: *ch,
		Valid:       true,
	}

	var resp interface{}
	err := cc.c.Send(ctx, "/channels", "PUT", cdi, &resp)
	if err != nil {
		return err
	}

	return nil
}
