package client

import (
	"context"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
	"gitlab.inspr.dev/inspr/core/pkg/rest/request"
)

// AliasClient interacts with Aliases on the Insprd
type AliasClient struct {
	c *request.Client
}

// Get gets a alias from the Insprd
//
// The context refers to the parent app of the given alias, represented with a dot separated query
// such as app1.app2
//
// The key is the key of the alias. So to search for a alias inside app1 with the key myKey you
// would call ac.Get(context.Background(), "app1", "myKey")
func (ac *AliasClient) Get(ctx context.Context, context, key string) (*meta.Alias, error) {
	aliasQuery := models.AliasQueryDI{
		Ctx: context,
		Key: key,
	}

	var resp meta.Alias

	err := ac.c.Send(ctx, "/alias", "GET", aliasQuery, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

///////////////// UPDATE THIS DOC //////////////////////////////
// Create creates a alias inside the Insprd
//
// The context refers to the parent app of the given alias, represented with a dot separated query
// such as **app1.app2**
//
// The alias information such as name and etc will be inferred from the given alias metadata.
//
// So to create a channel type inside app1 with the name channeltype1 you
// would call ctc.Create(context.Background(), "app1", &meta.ChannelType{...})
func (ac *AliasClient) Create(ctx context.Context, context string, target string, alias *meta.Alias) (diff.Changelog, error) {
	aliasQuery := models.AliasDI{
		Ctx:    context,
		Target: target,
		Alias:  *alias,
	}

	var resp diff.Changelog
	err := ac.c.Send(ctx, "/alias", "POST", aliasQuery, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a alias inside the Insprd
//
// The context refers to the parent app of the given alias, represented with a dot separated query
// such as **app1.app2**
//
// The key is the key of the alias to be deleted.
//
// So to delete a alias inside app1 with the key alias1 you
// would call ac.Delete(context.Background(), "app1", "alias1")
func (ac *AliasClient) Delete(ctx context.Context, context, key string) (diff.Changelog, error) {
	aliasQuery := models.AliasQueryDI{
		Ctx: context,
		Key: key,
	}

	var resp diff.Changelog
	err := ac.c.Send(ctx, "/alias", "DELETE", aliasQuery, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update updates a alias inside the Insprd
//
// The context refers to the parent app of the given alias, represented with a dot separated query
// such as **app1.app2**
//
// The alias information will be inferred from the given alias metadata.
//
// So to update a alias inside app1 with the key myalias you
// would call ac.Create(context.Background(), "app1", &meta.Alias{...})
func (ac *AliasClient) Update(ctx context.Context, context string, target string, alias *meta.Alias) (diff.Changelog, error) {
	aliasQuery := models.AliasDI{
		Ctx:    context,
		Target: target,
		Alias:  *alias,
	}

	var resp diff.Changelog
	err := ac.c.Send(ctx, "/alias", "PUT", aliasQuery, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
