package mocks

import (
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// Channels - mocks the implementation of the ChannelMemory interface methods
type Channels struct {
	fail error
}

// GetChannel - simple mock
func (chs *Channels) GetChannel(query string) (*meta.Channel, error) {
	if chs.fail != nil {
		return &meta.Channel{}, chs.fail
	}
	return &meta.Channel{}, nil
}

// CreateChannel - simple mock
func (chs *Channels) CreateChannel(ch *meta.Channel, context string) error {
	if chs.fail != nil {
		return chs.fail
	}
	return nil
}

// DeleteChannel - simple mock
func (chs *Channels) DeleteChannel(query string) error {
	if chs.fail != nil {
		return chs.fail
	}
	return nil
}

// UpdateChannel - simple mock
func (chs *Channels) UpdateChannel(ch *meta.Channel, query string) error {
	if chs.fail != nil {
		return chs.fail
	}
	return nil
}
