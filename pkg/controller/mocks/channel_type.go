package mocks

import (
	"context"

	"gitlab.inspr.dev/inspr/core/pkg/controller"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
)

// ChannelTypeMock mock structure for the operations of the controller.ChannelTypes()
type ChannelTypeMock struct {
	err error
}

// NewChannelTypeMock exports a mock of the channelType.interface
func NewChannelTypeMock(err error) controller.ChannelTypeInterface {
	return &ChannelTypeMock{err: err}
}

// Get is the ChannelTypeMock Get
func (cm *ChannelTypeMock) Get(ctx context.Context, context string, ctName string) (*meta.ChannelType, error) {
	if cm.err != nil {
		return &meta.ChannelType{}, cm.err
	}
	return &meta.ChannelType{}, nil
}

// Create is the ChannelTypeMock Create
func (cm *ChannelTypeMock) Create(ctx context.Context, context string, ct *meta.ChannelType, dryRun bool) (diff.Changelog, error) {
	if cm.err != nil {
		return diff.Changelog{}, cm.err
	}
	return diff.Changelog{}, nil
}

// Delete is the ChannelTypeMock Delete
func (cm *ChannelTypeMock) Delete(ctx context.Context, context string, ctName string, dryRun bool) (diff.Changelog, error) {
	if cm.err != nil {
		return diff.Changelog{}, cm.err
	}
	return diff.Changelog{}, nil
}

// Update is the ChannelTypeMock Update
func (cm *ChannelTypeMock) Update(ctx context.Context, context string, ct *meta.ChannelType, dryRun bool) (diff.Changelog, error) {
	if cm.err != nil {
		return diff.Changelog{}, cm.err
	}
	return diff.Changelog{}, nil
}
