package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// TreeMemoryManager defines a memory manager interface
type TreeMemoryManager struct {
	root *meta.App
}

var tree memory.Manager

// GetTreeMemory returns a memory manager interface
func GetTreeMemory() memory.Manager {
	if tree == nil {
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
