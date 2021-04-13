package client

import (
	"context"

	"github.com/inspr/inspr/cmd/insprd/api/models"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/utils/diff"
	"github.com/inspr/inspr/pkg/rest/request"
)

// ChannelTypeClient interacts with channeltypes on the Insprd
type ChannelTypeClient struct {
	c *request.Client
}

// Get gets a channel type from the Insprd
//
// The context refers to the parent app of the given channel type, represented with a dot separated query
// such as app1.app2
//
// The name is the name of the channel type. So to search for a channel type inside app1 with the name channeltype1 you
// would call ctc.Get(context.Background(), "app1", "channeltype1")
func (ctc *ChannelTypeClient) Get(ctx context.Context, context string, name string) (*meta.ChannelType, error) {
	ctdi := models.ChannelTypeQueryDI{
		Ctx:    context,
		CtName: name,
	}

	var resp meta.ChannelType

	err := ctc.c.Send(ctx, "/channeltypes", "GET", ctdi, &resp)
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
// So to create a channel type inside app1 with the name channeltype1 you
// would call ctc.Create(context.Background(), "app1", &meta.ChannelType{...})
func (ctc *ChannelTypeClient) Create(ctx context.Context, context string, ch *meta.ChannelType, dryRun bool) (diff.Changelog, error) {
	ctdi := models.ChannelTypeDI{
		Ctx:         context,
		ChannelType: *ch,
		DryRun:      dryRun,
	}

	var resp diff.Changelog
	err := ctc.c.Send(ctx, "/channeltypes", "POST", ctdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a channel type inside the Insprd
//
// The context refers to the parent app of the given channel type, represented with a dot separated query
// such as **app1.app2**
//
// The name is the name of the channel type to be deleted.
//
// So to delete a channel type inside app1 with the name channeltype1 you
// would call ctc.Delete(context.Background(), "app1", "channeltype1")
func (ctc *ChannelTypeClient) Delete(ctx context.Context, context string, name string, dryRun bool) (diff.Changelog, error) {
	ctdi := models.ChannelTypeQueryDI{
		Ctx:    context,
		CtName: name,
		DryRun: dryRun,
	}

	var resp diff.Changelog
	err := ctc.c.Send(ctx, "/channeltypes", "DELETE", ctdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update creates a channel type inside the Insprd
//
// The context refers to the parent app of the given channel type, represented with a dot separated query
// such as **app1.app2**
//
// The channel type information such as name and etc will be inferred from the given channel type's metadata.
//
// So to update a channel type inside app1 with the name channeltype1 you
// would call ctc.Create(context.Background(), "app1", &meta.ChannelType{...})
func (ctc *ChannelTypeClient) Update(ctx context.Context, context string, ch *meta.ChannelType, dryRun bool) (diff.Changelog, error) {
	ctdi := models.ChannelTypeDI{
		Ctx:         context,
		ChannelType: *ch,
		DryRun:      dryRun,
	}

	var resp diff.Changelog
	err := ctc.c.Send(ctx, "/channeltypes", "PUT", ctdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
