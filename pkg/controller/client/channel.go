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

// ChannelClient interacts with channels on the Insprd
type ChannelClient struct {
	reqClient *request.Client
}

// Get gets a channel from the Insprd
//
// The scope refers to the parent app of the given channel, represented with a dot separated query
// such as app1.app2
//
// The name is the name of the channel. So to search for a channel inside app1 with the name channel1 you
// would call cc.Get(context.Background(), "app1", "channel1")
func (cc *ChannelClient) Get(ctx context.Context, scope string, name string) (*meta.Channel, error) {
	cdi := models.ChannelQueryDI{
		ChName: name,
	}
	var resp meta.Channel

	err := cc.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/channels", http.MethodGet, cdi, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Create creates a channel inside the Insprd
//
// The scope refers to the parent app of the given channel, represented with a dot separated query
// such as **app1.app2**
//
// The channel information such as name and etc will be inferred from the given channel's metadata.
//
// So to create a channel inside app1 with the name channel1 you
// would call cc.Create(context.Background(), "app1", &meta.Channel{...})
func (cc *ChannelClient) Create(ctx context.Context, scope string, ch *meta.Channel, dryRun bool) (diff.Changelog, error) {
	cdi := models.ChannelDI{
		Channel: *ch,
		DryRun:  dryRun,
	}
	var resp diff.Changelog

	err := cc.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/channels", http.MethodPost, cdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a channel inside the Insprd
//
// The scope refers to the parent app of the given channel, represented with a dot separated query
// such as **app1.app2**
//
// The name is the name of the channel to be deleted.
//
// So to delete a channel inside app1 with the name channel1 you
// would call cc.Delete(context.Background(), "app1", "channel1")
func (cc *ChannelClient) Delete(ctx context.Context, scope string, name string, dryRun bool) (diff.Changelog, error) {
	cdi := models.ChannelQueryDI{
		ChName: name,
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := cc.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/channels", "DELETE", cdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update creates a channel inside the Insprd
//
// The scope refers to the parent app of the given channel, represented with a dot separated query
// such as **app1.app2**
//
// The channel information such as name and etc will be inferred from the given channel's metadata.
//
// So to update a channel inside app1 with the name channel1 you
// would call cc.Update(context.Background(), "app1", &meta.Channel{...})
func (cc *ChannelClient) Update(ctx context.Context, scope string, ch *meta.Channel, dryRun bool) (diff.Changelog, error) {
	cdi := models.ChannelDI{
		Channel: *ch,
		DryRun:  dryRun,
	}
	var resp diff.Changelog

	err := cc.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/channels", "PUT", cdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
