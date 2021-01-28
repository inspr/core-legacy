package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// TreeMockManager mocks a tree structure for testing
type TreeMockManager struct {
	root   *meta.App
	appErr error
	mockC  bool
	mockCT bool
	mockA  bool
}

// Channels mocks a channel interface for testing
func (tmm *TreeMockManager) Channels() memory.ChannelMemory {
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
func (tmm *TreeMockManager) ChannelTypes() memory.ChannelTypeMemory {
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
func (tmm *TreeMockManager) Apps() memory.AppMemory {
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
