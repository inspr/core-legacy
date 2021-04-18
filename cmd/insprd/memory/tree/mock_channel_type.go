package tree

import (
	"inspr.dev/inspr/pkg/meta"
)

// ChannelTypeMockManager mocks a channelType interface for testing
type ChannelTypeMockManager struct {
	*MockManager
}

// Create mocks a channelType method for testing
func (ctm *ChannelTypeMockManager) Create(context string, ct *meta.ChannelType) error {
	return nil
}

// Get mocks a channelType method for testing
func (ctm *ChannelTypeMockManager) Get(context string, ctName string) (*meta.ChannelType, error) {
	return nil, nil
}

// Delete mocks a channelType method for testing
func (ctm *ChannelTypeMockManager) Delete(context string, ctName string) error {
	return nil
}

// Update mocks a channelType method for testing
func (ctm *ChannelTypeMockManager) Update(query string, ct *meta.ChannelType) error {
	return nil
}
