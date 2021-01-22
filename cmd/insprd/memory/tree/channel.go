package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

type ChannelMemoryManager struct {
	root *meta.App
}

func (tmm *TreeMemoryManager) Channels() memory.ChannelMemory {
	return &ChannelMemoryManager{
		root: tmm.root,
	}
}

func (cmm *ChannelMemoryManager) GetChannel(context string, chName string) (*meta.Channel, error) {
	return nil, nil
}
func (cmm *ChannelMemoryManager) CreateChannel(ch *meta.Channel, context string) error {
	return nil
}
func (cmm *ChannelMemoryManager) DeleteChannel(context string, chName string) error {
	return nil
}
func (cmm *ChannelMemoryManager) UpdateChannel(ch *meta.Channel, query string) error {
	return nil
}
