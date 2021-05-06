package tree

import (
	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/pkg/meta/utils/diff"
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
			MockManager: tmm,
		}
	}
	return &ChannelMemoryManager{
		MemoryManager: tmm.MemoryManager,
	}
}

// Types mocks a Type interface for testing
func (tmm *MockManager) Types() memory.TypeMemory {
	if tmm.mockCT {
		return &TypeMockManager{
			MockManager: tmm,
		}
	}
	return &TypeMemoryManager{
		MemoryManager: tmm.MemoryManager,
	}
}

// Apps mocks an app interface for testing
func (tmm *MockManager) Apps() memory.AppMemory {
	if tmm.mockA {
		return &MockAppManager{
			MockManager: tmm,
			err:         tmm.appErr,
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

// Root mock interface structure
func (tmm *MockManager) Root() memory.GetInterface {
	return &RootGetter{
		tmm.root,
	}
}
