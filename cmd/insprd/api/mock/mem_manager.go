package repos

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
)

// MemManager is the api struct with the necessary implementations
// to satisfy the interface used in the routes established
type MemManager struct {
	channelType ChannelTypes
	channel     Channels
	app         Apps
	memory.Manager
}

// Init configurates the Memory Manager
func (mm *MemManager) Init() {
	mm.channelType = ChannelTypes{}
	mm.channel = Channels{}
	mm.app = Apps{}
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
