package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

type TreeMemoryManager struct {
	root *meta.App
}

func NewTreeMemory() memory.Manager {
	return &TreeMemoryManager{
		root: &meta.App{},
	}
}

func (tmm *TreeMemoryManager) Channels() memory.ChannelMemory {
	return nil
}

func (tmm *TreeMemoryManager) ChannelTypes() memory.ChannelTypeMemory {
	return nil
}
