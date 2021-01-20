package repos

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// Channels todo doc
type Channels struct {
	memory.ChannelMemory
}

// GetChannel todo doc
func (chs *Channels) GetChannel(ref string) (*meta.Channel, error) {
	return &meta.Channel{}, nil
}

// CreateChannel todo doc
func (chs *Channels) CreateChannel(ch *meta.Channel) error {
	return nil
}

// DeleteChannel todo doc
func (chs *Channels) DeleteChannel(ref string) error {
	return nil
}

// UpdateChannel todo doc
func (chs *Channels) UpdateChannel(ch *meta.Channel, ref string) error {
	return nil
}
