package mocks

import (
	"context"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
)

type channelMock struct {
	err error
}

// NewChannelMock exports a mock of the channel.interface
func NewChannelMock(err error) *channelMock {
	return &channelMock{err: err}
}

func (cm *channelMock) Get(ctx context.Context, context string, chName string) (*meta.Channel, error) {
	if cm.err != nil {
		return &meta.Channel{}, cm.err
	}
	return &meta.Channel{}, nil
}
func (cm *channelMock) Create(ctx context.Context, context string, ch *meta.Channel, dryRun bool) (diff.Changelog, error) {
	if cm.err != nil {
		return diff.Changelog{}, cm.err
	}
	return diff.Changelog{}, nil
}
func (cm *channelMock) Delete(ctx context.Context, context string, chName string, dryRun bool) (diff.Changelog, error) {
	if cm.err != nil {
		return diff.Changelog{}, cm.err
	}
	return diff.Changelog{}, nil
}
func (cm *channelMock) Update(ctx context.Context, context string, ch *meta.Channel, dryRun bool) (diff.Changelog, error) {
	if cm.err != nil {
		return diff.Changelog{}, cm.err
	}
	return diff.Changelog{}, nil
}
