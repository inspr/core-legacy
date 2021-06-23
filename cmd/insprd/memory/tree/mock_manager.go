package tree

import (
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils/diff"
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
func (tmm *MockManager) Channels() ChannelMemory {
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
func (tmm *MockManager) Types() TypeMemory {
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
func (tmm *MockManager) Apps() AppMemory {
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
func (tmm *MockManager) Root() GetInterface {
	return &RootGetter{
		tmm.root,
	}
}

// SetMockedTree receives a mock manager that has the configs of the
// tree structure to be mocked and used in tests where tree access is needed
func SetMockedTree(root *meta.App, appErr error, mockC, mockA, mockT bool) {
	setTree(&MockManager{
		MemoryManager: &MemoryManager{
			root: root,
			tree: root,
		},
		appErr: appErr,
		mockC:  mockC,
		mockA:  mockA,
		mockCT: mockT,
	})
}
