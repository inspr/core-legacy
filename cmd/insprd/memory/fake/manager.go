package fake

import (
	"inspr.dev/inspr/cmd/insprd/memory"
	"inspr.dev/inspr/cmd/insprd/memory/brokers"
	"inspr.dev/inspr/cmd/insprd/memory/tree"
)

// MemoryMock is the struct with the necessary implementations
// to mock the interface used to manage memory
type MemoryMock struct {
	tree    tree.Manager
	brokers brokers.Manager
}

// GetMockMemoryManager returns a mock for generic memory managing
func GetMockMemoryManager(treeErr, brokerErr error) memory.Manager {
	return &MemoryMock{
		tree:    MockTreeMemory(treeErr),
		brokers: MockBrokerMemory(brokerErr),
	}
}

// Tree returns a tree memory mocked manager
func (mm *MemoryMock) Tree() tree.Manager {
	return mm.tree
}

// Brokers returns a broker memory mocked manager
func (mm *MemoryMock) Brokers() brokers.Manager {
	return mm.brokers
}
