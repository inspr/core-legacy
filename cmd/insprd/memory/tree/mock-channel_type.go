package tree

import (
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// ChannelTypeMockManager mocks a channelType interface for testing
type ChannelTypeMockManager struct {
	*MockManager
}

// CreateChannelType mocks a channelType method for testing
func (ctm *ChannelTypeMockManager) CreateChannelType(ct *meta.ChannelType, context string) error {
	return nil
}

// Get mocks a channelType method for testing
func (ctm *ChannelTypeMockManager) Get(context string, ctName string) (*meta.ChannelType, error) {
	return nil, nil
}

// DeleteChannelType mocks a channelType method for testing
func (ctm *ChannelTypeMockManager) DeleteChannelType(context string, ctName string) error {
	return nil
}

// UpdateChannelType mocks a channelType method for testing
func (ctm *ChannelTypeMockManager) UpdateChannelType(ct *meta.ChannelType, query string) error {
	return nil
}
