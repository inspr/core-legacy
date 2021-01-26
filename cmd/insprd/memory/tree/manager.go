package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// TreeMemoryManager DOC TODO
type TreeMemoryManager struct {
	root *meta.App
}

var tree memory.Manager

// GetTreeMemory DOC TODO
func GetTreeMemory() memory.Manager {
	if tree == nil {
		// tree = newTreeMemory()
		setTree(newTreeMemory())
	}
	return tree
}

func newTreeMemory() *TreeMemoryManager {
	return &TreeMemoryManager{
		root: &meta.App{},
	}
}

func setTree(tmm memory.Manager) {
	tree = tmm
}

// Apps doc todo
func (tmm *TreeMemoryManager) Apps() memory.AppMemory {
	return nil
}
