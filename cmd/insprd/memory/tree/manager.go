package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

type TreeMemoryManager struct {
	root *meta.App
}

var dappTree *TreeMemoryManager

func GetTreeMemory() memory.Manager {
	if dappTree == nil {
		dappTree = newTreeMemory()
	}
	return dappTree
}

func newTreeMemory() *TreeMemoryManager {
	return &TreeMemoryManager{
		root: &meta.App{},
	}
}

func (tmm *TreeMemoryManager) Channels() memory.ChannelMemory {
	return nil
}

func (tmm *TreeMemoryManager) Apps() memory.AppMemory {
	return nil
}
