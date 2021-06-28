package tree

import (
	apimodels "inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/meta"
)

// ChannelMockManager mocks a channel interface for testing
type ChannelMockManager struct {
	*MockManager
}

// Get mocks a channel method for testing
func (cmm *ChannelMockManager) Get(scope, name string) (*meta.Channel, error) {
	return nil, nil
}

// Create mocks a channel method for testing
func (cmm *ChannelMockManager) Create(scope string, ch *meta.Channel, brokers *apimodels.BrokersDI) error {
	return nil
}

// Delete mocks a channel method for testing
func (cmm *ChannelMockManager) Delete(scope, name string) error {
	return nil
}

// Update mocks a channel method for testing
func (cmm *ChannelMockManager) Update(query string, ch *meta.Channel) error {
	return nil
}
