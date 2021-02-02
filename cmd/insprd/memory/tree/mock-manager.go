package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// MockManager mocks a tree structure for testing
type MockManager struct {
	root   *meta.App
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
		root: tmm.root,
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
		root: tmm.root,
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
		root: tmm.root,
	}
}

func (tmm *MockManager) InitTransaction() error {
	return nil
}

func (tmm MockManager) Commit() {}
