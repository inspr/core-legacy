package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

//TreeMockManager is a TreeMemoryManager mock struct
type TreeMockManager struct {
	root   *meta.App
	appErr error
	mockC  bool
	mockCT bool
	mockA  bool
}

//Channels mock
func (tmm *TreeMockManager) Channels() memory.ChannelMemory {
	if tmm.mockC {
		return nil
	}
	return nil // ChannelMemoryManager
}

//ChannelTypes mock
func (tmm *TreeMockManager) ChannelTypes() memory.ChannelTypeMemory {
	if tmm.mockCT {
		return nil
	}
	return &ChannelTypeMemoryManager{
		root: tmm.root,
	}
}

//Apps mock
func (tmm *TreeMockManager) Apps() memory.AppMemory {
	if tmm.mockA {
		return &MockAppManager{
			root: tmm.root,
			err:  tmm.appErr,
		}
	}
	return nil //AppMemoryManager
}
