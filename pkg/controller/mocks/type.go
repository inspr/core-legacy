package mocks

import (
	"context"

	"inspr.dev/inspr/pkg/controller"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils/diff"
)

// TypeMock mock structure for the operations of the controller.Types()
type TypeMock struct {
	err error
}

// NewTypeMock exports a mock of the Type.interface
func NewTypeMock(err error) controller.TypeInterface {
	return &TypeMock{err: err}
}

// Get is the TypeMock Get
func (tm *TypeMock) Get(ctx context.Context, scope, ctName string) (*meta.Type, error) {
	if tm.err != nil {
		return &meta.Type{}, tm.err
	}
	return &meta.Type{}, nil
}

// Create is the TypeMock Create
func (tm *TypeMock) Create(ctx context.Context, scope string, ct *meta.Type, dryRun bool) (diff.Changelog, error) {
	if tm.err != nil {
		return diff.Changelog{}, tm.err
	}
	return diff.Changelog{}, nil
}

// Delete is the TypeMock Delete
func (tm *TypeMock) Delete(ctx context.Context, scope, ctName string, dryRun bool) (diff.Changelog, error) {
	if tm.err != nil {
		return diff.Changelog{}, tm.err
	}
	return diff.Changelog{}, nil
}

// Update is the TypeMock Update
func (tm *TypeMock) Update(ctx context.Context, scope string, ct *meta.Type, dryRun bool) (diff.Changelog, error) {
	if tm.err != nil {
		return diff.Changelog{}, tm.err
	}
	return diff.Changelog{}, nil
}
