package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

type TreeMockManager struct {
	root   *meta.App
	appErr error
	mockC  bool
	mockCT bool
	mockA  bool
}

func (tmm *TreeMockManager) Channels() memory.ChannelMemory {
	if tmm.mockC {
		return &ChannelMockManager{
			root: tmm.root,
		}
	}
	return nil
}

func (tmm *TreeMockManager) ChannelTypes() memory.ChannelTypeMemory {
	if tmm.mockCT {
		return &ChannelTypeMockManager{
			root: tmm.root,
		}
	}
	return nil
}

func (tmm *TreeMockManager) Apps() memory.AppMemory {
	if tmm.mockA {
		return nil
	}
	return &AppMemoryManager{
		root: tmm.root,
	}
}
