package mocks

import (
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// ChannelTypes - mocks the implementation of the ChannelTypeMemory interface methods
type ChannelTypes struct {
	fail error
}

// GetChannelType - simple mock
func (chType *ChannelTypes) GetChannelType(context string, ctName string) (*meta.ChannelType, error) {
	if chType.fail != nil {
		return &meta.ChannelType{}, chType.fail
	}
	return &meta.ChannelType{}, nil
}

// CreateChannelType - simple mock
func (chType *ChannelTypes) CreateChannelType(ct *meta.ChannelType, context string) error {
	if chType.fail != nil {
		return chType.fail
	}
	return nil
}

// DeleteChannelType - simple mock
func (chType *ChannelTypes) DeleteChannelType(context string, ctName string) error {
	if chType.fail != nil {
		return chType.fail
	}
	return nil
}

// UpdateChannelType - simple mock
func (chType *ChannelTypes) UpdateChannelType(ct *meta.ChannelType, context string) error {
	if chType.fail != nil {
		return chType.fail
	}
	return nil
}
