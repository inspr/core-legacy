package mocks

import (
	"context"

	"inspr.dev/inspr/pkg/controller"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils/diff"
)

// ChannelMock mock structure for the operations of the controller.Channels()
type ChannelMock struct {
	err error
}

// NewChannelMock exports a mock of the channel.interface
func NewChannelMock(err error) controller.ChannelInterface {
	return &ChannelMock{err: err}
}

// Get is the channelmock Get
func (cm *ChannelMock) Get(
	ctx context.Context,
	scope, name string,
) (*meta.Channel, error) {
	if cm.err != nil {
		return &meta.Channel{}, cm.err
	}
	return &meta.Channel{}, nil
}

// Create is the channelmock Create
func (cm *ChannelMock) Create(
	ctx context.Context,
	scope string,
	ch *meta.Channel,
	dryRun bool,
) (diff.Changelog, error) {
	if cm.err != nil {
		return diff.Changelog{}, cm.err
	}
	return diff.Changelog{}, nil
}

// Delete is the channelmock Delete
func (cm *ChannelMock) Delete(
	ctx context.Context,
	scope, name string,
	dryRun bool,
) (diff.Changelog, error) {
	if cm.err != nil {
		return diff.Changelog{}, cm.err
	}
	return diff.Changelog{}, nil
}

// Update is the channelmock Update
func (cm *ChannelMock) Update(
	ctx context.Context,
	scope string,
	ch *meta.Channel,
	dryRun bool,
) (diff.Changelog, error) {
	if cm.err != nil {
		return diff.Changelog{}, cm.err
	}
	return diff.Changelog{}, nil
}
