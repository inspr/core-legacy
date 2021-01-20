package models

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// MemManager is the api struct with the necessary implementations
// to satisfy the interface used in the routes established
type MemManager struct {
	ChannelType meta.ChannelType
	Channel     meta.Channel
	DApp        meta.AppSpec
}

// Apps returns manager's DApp
func (mm *MemManager) Apps() memory.AppMemory {
	return &mm.DApp
}

// Channels returns manager's DApp
func (mm *MemManager) Channels() memory.ChannelMemory {
	return &mm.Channel
}

// ChannelTypes returns manager's DApp
func (mm *MemManager) ChannelTypes() memory.ChannelTypeMemory {
	return &mm.ChannelType
}
