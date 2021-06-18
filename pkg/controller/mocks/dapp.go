package mocks

import (
	"context"

	"inspr.dev/inspr/pkg/controller"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils/diff"
)

// AppMock mock structure for the operations of the controller.Apps()
type AppMock struct {
	err error
}

// NewAppMock exports a mock of the App.interface
func NewAppMock(err error) controller.AppInterface {
	return &AppMock{err: err}
}

// Get is the AppMock Get
func (am *AppMock) Get(ctx context.Context, scope string) (*meta.App, error) {
	if am.err != nil {
		return &meta.App{}, am.err
	}
	return &meta.App{}, nil
}

// Create is the AppMock Create
func (am *AppMock) Create(ctx context.Context, scope string, app *meta.App, dryRun bool) (diff.Changelog, error) {
	if am.err != nil {
		return diff.Changelog{}, am.err
	}
	return diff.Changelog{}, nil
}

// Delete is the AppMock Delete
func (am *AppMock) Delete(ctx context.Context, scope string, dryRun bool) (diff.Changelog, error) {
	if am.err != nil {
		return diff.Changelog{}, am.err
	}
	return diff.Changelog{}, nil
}

// Update is the AppMock Update
func (am *AppMock) Update(ctx context.Context, scope string, app *meta.App, dryRun bool) (diff.Changelog, error) {
	if am.err != nil {
		return diff.Changelog{}, am.err
	}
	return diff.Changelog{}, nil
}
