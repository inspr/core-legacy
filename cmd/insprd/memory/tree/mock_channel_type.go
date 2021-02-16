package tree

import (
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// ChannelTypeMockManager mocks a channelType interface for testing
type ChannelTypeMockManager struct {
	*MockManager
}

// CreateChannelType mocks a channelType method for testing
func (ctm *ChannelTypeMockManager) CreateChannelType(context string, ct *meta.ChannelType) error {
	return nil
}

// GetChannelType mocks a channelType method for testing
func (ctm *ChannelTypeMockManager) GetChannelType(context string, ctName string) (*meta.ChannelType, error) {
	return nil, nil
}

// DeleteChannelType mocks a channelType method for testing
func (ctm *ChannelTypeMockManager) DeleteChannelType(context string, ctName string) error {
	return nil
}

// UpdateChannelType mocks a channelType method for testing
func (ctm *ChannelTypeMockManager) UpdateChannelType(query string, ct *meta.ChannelType) error {
	return nil
}
