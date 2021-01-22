package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

type ChannelTypeMemoryManager struct {
	root *meta.App
}

func (tmm *TreeMemoryManager) ChannelTypes() memory.ChannelTypeMemory {

	return &ChannelTypeMemoryManager{
		root: tmm.root,
	}
}

func (ctm *ChannelTypeMemoryManager) CreateChannelType(ct *meta.ChannelType, context string) error {
	return nil
}

func (ctm *ChannelTypeMemoryManager) GetChannelType(context string, ctName string) (*meta.ChannelType, error) {
	return nil, nil
}

func (ctm *ChannelTypeMemoryManager) DeleteChannelType(context string, ctName string) error {
	return nil
}

func (ctm *ChannelTypeMemoryManager) UpdateChannelType(ct *meta.ChannelType, query string) error {
	return nil
}
