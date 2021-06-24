package fake

import (
	"inspr.dev/inspr/cmd/insprd/memory"
	"inspr.dev/inspr/cmd/insprd/memory/brokers"
	"inspr.dev/inspr/cmd/insprd/memory/tree"
)

type MemoryMock struct {
	tree    tree.Manager
	brokers brokers.Manager
}

func GetMockMemoryManager(treeErr, brokerErr error) memory.Manager {
	return &MemoryMock{
		tree:    MockTreeMemory(treeErr),
		brokers: MockBrokerMemory(brokerErr),
	}
}

func (mm *MemoryMock) Tree() tree.Manager {
	return mm.tree
}

func (mm *MemoryMock) Brokers() brokers.Manager {
	return mm.brokers
}
