package memory

import (
	"inspr.dev/inspr/cmd/insprd/memory/brokers"
	"inspr.dev/inspr/cmd/insprd/memory/tree"
)

type memoryManager struct {
	tree    tree.Manager
	brokers brokers.Manager
}

var memManager *memoryManager

// GetMemoryManager returns a general memory manager interface
func GetMemoryManager() Manager {
	if memManager != nil {
		return memManager
	}
	return &memoryManager{
		tree:    tree.GetTreeMemory(),
		brokers: brokers.GetBrokerMemory(),
	}
}

// Tree  retrunrs the tree specific memory managing interface
func (mem *memoryManager) Tree() tree.Manager {
	return mem.tree
}

// Brokers retrunrs the broker specific memory managing interface
func (mem *memoryManager) Brokers() brokers.Manager {
	return mem.brokers
}
