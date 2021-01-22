package mocks

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// Channels todo doc
type Channels struct {
	memory.ChannelMemory
}

// GetChannel todo doc
func (chs *Channels) GetChannel(query string) (*meta.Channel, error) {
	return &meta.Channel{}, nil
}

// CreateChannel todo doc
func (chs *Channels) CreateChannel(ch *meta.Channel, context string) error {
	return nil
}

// DeleteChannel todo doc
func (chs *Channels) DeleteChannel(query string) error {
	return nil
}

// UpdateChannel todo doc
func (chs *Channels) UpdateChannel(ch *meta.Channel, query string) error {
	return nil
}
