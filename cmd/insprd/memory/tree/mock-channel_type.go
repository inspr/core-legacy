package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

type ChannelTypeMockManager struct {
	root *meta.App
}

func (tmm *TreeMemoryManager) ChannelTypes() memory.ChannelTypeMemory {
	return &ChannelTypeMockManager{
		root: tmm.root,
	}
}

func (ctm *ChannelTypeMockManager) CreateChannelType(ct *meta.ChannelType, context string) error {
	return nil
}

func (ctm *ChannelTypeMockManager) GetChannelType(context string, ctName string) (*meta.ChannelType, error) {
	return nil, nil
}

func (ctm *ChannelTypeMockManager) DeleteChannelType(context string, ctName string) error {
	return nil
}

func (ctm *ChannelTypeMockManager) UpdateChannelType(ct *meta.ChannelType, query string) error {
	return nil
}
