package mocks

import (
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// Channels - mocks the implementation of the ChannelMemory interface methods
type Channels struct {
	fail error
	*MemManager
}

// GetChannel - simple mock
func (chs *Channels) GetChannel(context string, chName string) (*meta.Channel, error) {
	if chs.fail != nil {
		return &meta.Channel{}, chs.fail
	}
	return &meta.Channel{}, nil
}

// CreateChannel - simple mock
func (chs *Channels) CreateChannel(context string, ch *meta.Channel) error {
	if chs.fail != nil {
		return chs.fail
	}
	return nil
}

// DeleteChannel - simple mock
func (chs *Channels) DeleteChannel(context string, chName string) error {
	if chs.fail != nil {
		return chs.fail
	}
	return nil
}

// UpdateChannel - simple mock
func (chs *Channels) UpdateChannel(context string, ch *meta.Channel) error {
	if chs.fail != nil {
		return chs.fail
	}
	return nil
}
