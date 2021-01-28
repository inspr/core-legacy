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
	return nil
}

// ChannelTypes mocks a channelType interface for testing
func (tmm *TreeMockManager) ChannelTypes() memory.ChannelTypeMemory {
	if tmm.mockCT {
		return &ChannelTypeMockManager{
			root: tmm.root,
		}
	}
	return nil
}

// Apps mocks an app interface for testing
func (tmm *TreeMockManager) Apps() memory.AppMemory {
	if tmm.mockA {
		return nil
	}
	return &AppMemoryManager{
		root: tmm.root,
	}
}
