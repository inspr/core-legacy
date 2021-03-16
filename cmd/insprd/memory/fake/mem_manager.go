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
	alias       Alias
}

// LookupMemManager mocks getter for roots
type LookupMemManager MemManager

// Apps mocks an app getter
func (l LookupMemManager) Apps() memory.AppGetInterface {
	return &l.app
}

// Channels mocks a channel getter
func (l LookupMemManager) Channels() memory.ChannelGetInterface {
	return &l.channel
}

// ChannelTypes mocks a channel type getter
func (l LookupMemManager) ChannelTypes() memory.ChannelTypeGetInterface {
	return &l.channelType
}

// Alias mocks a alias getter
func (l LookupMemManager) Alias() memory.AliasGetInterface {
	return &l.alias
}

// MockMemoryManager mock exported with propagated error through the functions
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
		alias: Alias{
			fail:  failErr,
			alias: make(map[string]*meta.Alias),
		},
	}
}

// Root mocks a root getter interface
func (mm *MemManager) Root() memory.GetInterface {
	return (*LookupMemManager)(mm)
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

// Alias returns manager's Alias
func (mm *MemManager) Alias() memory.AliasMemory {
	return &mm.alias
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
