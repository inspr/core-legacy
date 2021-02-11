package tree

import (
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// ChannelMockManager mocks a channel interface for testing
type ChannelMockManager struct {
	root *meta.App
}

// GetChannel mocks a channel method for testing
func (cmm *ChannelMockManager) GetChannel(context string, chName string) (*meta.Channel, error) {
	return nil, nil
}

// CreateChannel mocks a channel method for testing
func (cmm *ChannelMockManager) CreateChannel(context string, ch *meta.Channel) error {
	return nil
}

// DeleteChannel mocks a channel method for testing
func (cmm *ChannelMockManager) DeleteChannel(context string, chName string) error {
	return nil
}

// UpdateChannel mocks a channel method for testing
func (cmm *ChannelMockManager) UpdateChannel(query string, ch *meta.Channel) error {
	return nil
}
