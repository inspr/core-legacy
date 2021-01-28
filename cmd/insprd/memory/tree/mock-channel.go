package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// ChannelMockManager mocks a channel interface for testing
type ChannelMockManager struct {
	root *meta.App
}

// Channels returns a mocked channel interface for testing
func (tmm *TreeMemoryManager) Channels() memory.ChannelMemory {
	return &ChannelMockManager{
		root: tmm.root,
	}
}

// GetChannel mocks a channel method for testing
func (cmm *ChannelMockManager) GetChannel(context string, chName string) (*meta.Channel, error) {
	return nil, nil
}

// CreateChannel mocks a channel method for testing
func (cmm *ChannelMockManager) CreateChannel(ch *meta.Channel, context string) error {
	return nil
}

// DeleteChannel mocks a channel method for testing
func (cmm *ChannelMockManager) DeleteChannel(context string, chName string) error {
	return nil
}

// UpdateChannel mocks a channel method for testing
func (cmm *ChannelMockManager) UpdateChannel(ch *meta.Channel, query string) error {
	return nil
}
