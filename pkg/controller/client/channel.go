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

// ChannelClient is a client for manipulating Channel structures in Insprd
type ChannelClient struct {
	reqClient *request.Client
}

// Get gets a channel from Insprd
// The scope refers to the parent app of the given channel, represented with a dot separated query
// such as app1.app2. The name is the name of the channel.
func (cc *ChannelClient) Get(ctx context.Context, scope, name string) (*meta.Channel, error) {
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

// Create creates given channel inside of Insprd
// The scope refers to the parent app of the given channel, represented with a dot separated query
// such as app1.app2
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
// The scope refers to the parent app of the given channel, represented with a dot separated query
// such as app1.app2. The name is the name of the channel to be deleted
func (cc *ChannelClient) Delete(ctx context.Context, scope, name string, dryRun bool) (diff.Changelog, error) {
	cdi := models.ChannelQueryDI{
		ChName: name,
		DryRun: dryRun,
	}
	var resp diff.Changelog

	err := cc.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/channels", http.MethodDelete, cdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update updates given channel structure, if it exists in Insprd
// The scope refers to the parent app of the given channel, represented with a dot separated query
// such as app1.app2
func (cc *ChannelClient) Update(ctx context.Context, scope string, ch *meta.Channel, dryRun bool) (diff.Changelog, error) {
	cdi := models.ChannelDI{
		Channel: *ch,
		DryRun:  dryRun,
	}
	var resp diff.Changelog

	err := cc.reqClient.
		Header(rest.HeaderScopeKey, scope).
		Send(ctx, "/channels", http.MethodPut, cdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
