package client

import (
	"context"

	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/utils/diff"
	"github.com/inspr/inspr/pkg/rest/request"
)

// TypeClient interacts with Types on the Insprd
type TypeClient struct {
	c *request.Client
}

// Get gets a Type from the Insprd
//
// The context refers to the parent app of the given Type, represented with a dot separated query
// such as app1.app2
//
// The name is the name of the Type. So to search for a Type inside app1 with the name Type1 you
// would call ctc.Get(context.Background(), "app1", "Type1")
func (ctc *TypeClient) Get(ctx context.Context, context string, name string) (*meta.Type, error) {
	ctdi := models.TypeQueryDI{
		Scope:  context,
		CtName: name,
	}

	var resp meta.Type

	err := ctc.c.Send(ctx, "/Types", "GET", ctdi, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// Create creates a Type inside the Insprd
//
// The context refers to the parent app of the given Type, represented with a dot separated query
// such as **app1.app2**
//
// The Type information such as name and etc will be inferred from the given Type's metadata.
//
// So to create a Type inside app1 with the name Type1 you
// would call ctc.Create(context.Background(), "app1", &meta.Type{...})
func (ctc *TypeClient) Create(ctx context.Context, context string, ch *meta.Type, dryRun bool) (diff.Changelog, error) {
	ctdi := models.TypeDI{
		Scope:  context,
		Type:   *ch,
		DryRun: dryRun,
	}

	var resp diff.Changelog
	err := ctc.c.Send(ctx, "/Types", "POST", ctdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete deletes a Type inside the Insprd
//
// The context refers to the parent app of the given Type, represented with a dot separated query
// such as **app1.app2**
//
// The name is the name of the Type to be deleted.
//
// So to delete a Type inside app1 with the name Type1 you
// would call ctc.Delete(context.Background(), "app1", "Type1")
func (ctc *TypeClient) Delete(ctx context.Context, context string, name string, dryRun bool) (diff.Changelog, error) {
	ctdi := models.TypeQueryDI{
		Scope:  context,
		CtName: name,
		DryRun: dryRun,
	}

	var resp diff.Changelog
	err := ctc.c.Send(ctx, "/Types", "DELETE", ctdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update updates a Type inside the Insprd
//
// The context refers to the parent app of the given Type, represented with a dot separated query
// such as **app1.app2**
//
// The Type information such as name and etc will be inferred from the given Type's metadata.
//
// So to update a Type inside app1 with the name Type1 you
// would call ctc.Create(context.Background(), "app1", &meta.Type{...})
func (ctc *TypeClient) Update(ctx context.Context, context string, ch *meta.Type, dryRun bool) (diff.Changelog, error) {
	ctdi := models.TypeDI{
		Scope:  context,
		Type:   *ch,
		DryRun: dryRun,
	}

	var resp diff.Changelog
	err := ctc.c.Send(ctx, "/Types", "PUT", ctdi, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
