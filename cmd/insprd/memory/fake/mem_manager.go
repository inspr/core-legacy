package fake

import (
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// MemManager is the api struct with the necessary implementations
// to satisfy the interface used in the routes established
type MemManager struct {
	channelType ChannelTypes
	channel     Channels
	app         Apps
}

// MockMemoryManager mock exported with propagated error throught the functions
func MockMemoryManager(failErr error) memory.Manager {
	return &MemManager{
		channelType: ChannelTypes{
			fail:         failErr,
			channelTypes: make(map[string]*meta.ChannelType),
		},
		channel: Channels{
			fail:     failErr,
			channels: make(map[string]*meta.Channel),
		},
		app: Apps{
			fail: failErr,
			apps: make(map[string]*meta.App),
		},
	}
}

// Apps returns manager of DApps
func (mm *MemManager) Apps() memory.AppMemory {
	return &mm.app
}

// Channels returns manager's DApp
func (mm *MemManager) Channels() memory.ChannelMemory {
	return &mm.channel
}

// ChannelTypes returns manager's DApp
func (mm *MemManager) ChannelTypes() memory.ChannelTypeMemory {
	return &mm.channelType
}

//InitTransaction mock interface structure
func (mm *MemManager) InitTransaction() {}

//Commit mock interface structure
func (mm *MemManager) Commit() {}

//Cancel mock interface structure
func (mm *MemManager) Cancel() {}

//GetTransactionChanges mock interface structure
func (mm *MemManager) GetTransactionChanges() (diff.Changelog, error) {
	return diff.Changelog{}, nil
}