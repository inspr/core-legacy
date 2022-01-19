package mocks

import (
	"context"

	"inspr.dev/inspr/pkg/controller"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils/diff"
)

// AliasMock mock structure for the operations of the controller.Aliass()
type AliasMock struct {
	err error
}

// NewAliasMock exports a mock of the Alias.interface
func NewAliasMock(err error) controller.AliasInterface {
	return &AliasMock{
		err: err,
	}
}

// Get is the AliasMock Get
func (am *AliasMock) Get(ctx context.Context, scope, name string) (*meta.Alias, error) {
	if am.err != nil {
		return &meta.Alias{}, am.err
	}
	return &meta.Alias{}, nil
}

// Create is the AliasMock Create
func (am *AliasMock) Create(ctx context.Context, scope string, alias *meta.Alias, dryRun bool) (diff.Changelog, error) {
	if am.err != nil {
		return diff.Changelog{}, am.err
	}
	return diff.Changelog{}, nil
}

// Delete is the AliasMock Delete
func (am *AliasMock) Delete(ctx context.Context, scope, name string, dryRun bool) (diff.Changelog, error) {
	if am.err != nil {
		return diff.Changelog{}, am.err
	}
	return diff.Changelog{}, nil
}

// Update is the AliasMock Update
func (am *AliasMock) Update(ctx context.Context, scope string, alias *meta.Alias, dryRun bool) (diff.Changelog, error) {
	if am.err != nil {
		return diff.Changelog{}, am.err
	}
	return diff.Changelog{}, nil
}
