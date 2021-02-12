package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/utils/diff"
)

// MockManager mocks a tree structure for testing
type MockManager struct {
	*MemoryManager
	appErr error
	mockC  bool
	mockCT bool
	mockA  bool
}

// Channels mocks a channel interface for testing
func (tmm *MockManager) Channels() memory.ChannelMemory {
	if tmm.mockC {
		return &ChannelMockManager{
			root: tmm.root,
		}
	}
	return &ChannelMemoryManager{
		MemoryManager: tmm.MemoryManager,
	}
}

// ChannelTypes mocks a channelType interface for testing
func (tmm *MockManager) ChannelTypes() memory.ChannelTypeMemory {
	if tmm.mockCT {
		return &ChannelTypeMockManager{
			root: tmm.root,
		}
	}
	return &ChannelTypeMemoryManager{
		MemoryManager: tmm.MemoryManager,
	}
}

// Apps mocks an app interface for testing
func (tmm *MockManager) Apps() memory.AppMemory {
	if tmm.mockA {
		return &MockAppManager{
			root: tmm.root,
			err:  tmm.appErr,
		}
	}
	return &AppMemoryManager{
		MemoryManager: tmm.MemoryManager,
	}
}

//InitTransaction mock interface structure
func (tmm *MockManager) InitTransaction() {}

//Commit mock interface structure
func (tmm *MockManager) Commit() {}

//Cancel mock interface structure
func (tmm *MockManager) Cancel() {}

//GetTransactionChanges mock structure
func (tmm *MockManager) GetTransactionChanges() (diff.Changelog, error) {
	return diff.Changelog{}, nil
}
