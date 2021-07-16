package tree

import (
	"inspr.dev/inspr/pkg/meta/utils/diff"
)

// MockManager mocks a tree structure for testing
type MockManager struct {
	*treeMemoryManager
	appErr error
	mockC  bool
	mockCT bool
	mockA  bool
}

// Channels mocks a channel interface for testing
func (tmm *MockManager) Channels() ChannelMemory {
	if tmm.mockC {
		return &ChannelMockManager{
			MockManager: tmm,
		}
	}
	return &ChannelMemoryManager{
		treeMemoryManager: tmm.treeMemoryManager,
		logger:            logger,
	}
}

// Types mocks a Type interface for testing
func (tmm *MockManager) Types() TypeMemory {
	if tmm.mockCT {
		return &TypeMockManager{
			MockManager: tmm,
		}
	}
	return &TypeMemoryManager{
		treeMemoryManager: tmm.treeMemoryManager,
		logger:            logger,
	}
}

// Apps mocks an app interface for testing
func (tmm *MockManager) Apps() AppMemory {
	if tmm.mockA {
		return &MockAppManager{
			MockManager: tmm,
			err:         tmm.appErr,
		}
	}
	return &AppMemoryManager{
		treeMemoryManager: tmm.treeMemoryManager,
		logger:            logger,
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

// Tree mock interface structure
func (tmm *MockManager) Tree() GetInterface {
	return &PermTreeGetter{
		tmm.root,
		logger,
	}
}
