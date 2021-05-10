package fake

import (
	"github.com/inspr/inspr/pkg/meta/utils/diff"

	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/pkg/meta"
)

// MemManager is the api struct with the necessary implementations
// to satisfy the interface used in the routes established
type MemManager struct {
	itype   Types // inspr type
	channel Channels
	app     Apps
	alias   Alias
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

// Types mocks a Type getter
func (l LookupMemManager) Types() memory.TypeGetInterface {
	return &l.itype
}

// Alias mocks a alias getter
func (l LookupMemManager) Alias() memory.AliasGetInterface {
	return &l.alias
}

// MockMemoryManager mock exported with propagated error through the functions
func MockMemoryManager(failErr error) memory.Manager {
	return &MemManager{
		itype: Types{
			fail:  failErr,
			Types: make(map[string]*meta.Type),
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

// Types returns manager's DApp
func (mm *MemManager) Types() memory.TypeMemory {
	return &mm.itype
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
