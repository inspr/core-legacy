package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

type ChannelMockManager struct {
	root *meta.App
}

func (tmm *TreeMemoryManager) Channels() memory.ChannelMemory {
	return &ChannelMockManager{
		root: tmm.root,
	}
}

func (cmm *ChannelMockManager) GetChannel(context string, chName string) (*meta.Channel, error) {
	return nil, nil
}
func (cmm *ChannelMockManager) CreateChannel(ch *meta.Channel, context string) error {
	return nil
}
func (cmm *ChannelMockManager) DeleteChannel(context string, chName string) error {
	return nil
}
func (cmm *ChannelMockManager) UpdateChannel(ch *meta.Channel, query string) error {
	return nil
}
