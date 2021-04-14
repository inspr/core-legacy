package mocks

import (
	"context"

	"gitlab.inspr.dev/inspr/core/pkg/controller"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
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
func (cm *AliasMock) Get(ctx context.Context, context, key string) (*meta.Alias, error) {
	if cm.err != nil {
		return &meta.Alias{}, cm.err
	}
	return &meta.Alias{}, nil
}

// Create is the AliasMock Create
func (cm *AliasMock) Create(ctx context.Context, context string, target string, alias *meta.Alias, dryRun bool) (diff.Changelog, error) {
	if cm.err != nil {
		return diff.Changelog{}, cm.err
	}
	return diff.Changelog{}, nil
}

// Delete is the AliasMock Delete
func (cm *AliasMock) Delete(ctx context.Context, context, key string, dryRun bool) (diff.Changelog, error) {
	if cm.err != nil {
		return diff.Changelog{}, cm.err
	}
	return diff.Changelog{}, nil
}

// Update is the AliasMock Update
func (cm *AliasMock) Update(ctx context.Context, context string, target string, alias *meta.Alias, dryRun bool) (diff.Changelog, error) {
	if cm.err != nil {
		return diff.Changelog{}, cm.err
	}
	return diff.Changelog{}, nil
}
