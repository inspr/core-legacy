package tree

import (
	"github.com/inspr/inspr/pkg/meta"
)

// ChannelMockManager mocks a channel interface for testing
type ChannelMockManager struct {
	*MockManager
}

// Get mocks a channel method for testing
func (cmm *ChannelMockManager) Get(context string, chName string) (*meta.Channel, error) {
	return nil, nil
}

// Create mocks a channel method for testing
func (cmm *ChannelMockManager) Create(context string, ch *meta.Channel) error {
	return nil
}

// Delete mocks a channel method for testing
func (cmm *ChannelMockManager) Delete(context string, chName string) error {
	return nil
}

// Update mocks a channel method for testing
func (cmm *ChannelMockManager) Update(query string, ch *meta.Channel) error {
	return nil
}
