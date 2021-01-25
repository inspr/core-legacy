package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// TreeMockManager mock tree manager
type TreeMockManager struct {
	root   *meta.App
	appErr error
	mockC  bool
	mockCT bool
	mockA  bool
}

// Channels MockChannel
func (tmm *TreeMockManager) Channels() memory.ChannelMemory {
	if tmm.mockC {
		return nil
	}
	return &ChannelMemoryManager{
		root: tmm.root,
	}
}

// ChannelTypes Mock channel types
func (tmm *TreeMockManager) ChannelTypes() memory.ChannelTypeMemory {
	if tmm.mockCT {
		return nil
	}
	return nil
}

// Apps Mock Apps
func (tmm *TreeMockManager) Apps() memory.AppMemory {
	if tmm.mockA {
		return &MockAppManager{
			root: tmm.root,
			err:  tmm.appErr,
		}
	}
	return nil //AppMemoryManager
}
