package mocks

import (
	"context"

	"github.com/inspr/inspr/pkg/controller"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/utils/diff"
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
func (cm *TypeMock) Get(ctx context.Context, context string, ctName string) (*meta.Type, error) {
	if cm.err != nil {
		return &meta.Type{}, cm.err
	}
	return &meta.Type{}, nil
}

// Create is the TypeMock Create
func (cm *TypeMock) Create(ctx context.Context, context string, ct *meta.Type, dryRun bool) (diff.Changelog, error) {
	if cm.err != nil {
		return diff.Changelog{}, cm.err
	}
	return diff.Changelog{}, nil
}

// Delete is the TypeMock Delete
func (cm *TypeMock) Delete(ctx context.Context, context string, ctName string, dryRun bool) (diff.Changelog, error) {
	if cm.err != nil {
		return diff.Changelog{}, cm.err
	}
	return diff.Changelog{}, nil
}

// Update is the TypeMock Update
func (cm *TypeMock) Update(ctx context.Context, context string, ct *meta.Type, dryRun bool) (diff.Changelog, error) {
	if cm.err != nil {
		return diff.Changelog{}, cm.err
	}
	return diff.Changelog{}, nil
}
