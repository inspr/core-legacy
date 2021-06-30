package fake

import (
	"inspr.dev/inspr/cmd/insprd/memory/tree"
	"inspr.dev/inspr/pkg/meta/utils/diff"

	"inspr.dev/inspr/pkg/meta"
)

// TreeMemoryMock is the api struct with the necessary implementations
// to satisfy the interface used in the routes established
type TreeMemoryMock struct {
	insprType Types // inspr type
	channel   Channels
	app       Apps
	alias     Alias
}

// LookupMemManager mocks getter for roots
type LookupMemManager TreeMemoryMock

// Apps mocks an app getter
func (l LookupMemManager) Apps() tree.AppGetInterface {
	return &l.app
}

// Channels mocks a channel getter
func (l LookupMemManager) Channels() tree.ChannelGetInterface {
	return &l.channel
}

// Types mocks a Type getter
func (l LookupMemManager) Types() tree.TypeGetInterface {
	return &l.insprType
}

// Alias mocks a alias getter
func (l LookupMemManager) Alias() tree.AliasGetInterface {
	return &l.alias
}

// MockTreeMemory mock exported with propagated error through the functions
func MockTreeMemory(failErr error) tree.Manager {
	return &TreeMemoryMock{
		insprType: Types{
			fail:       failErr,
			insprTypes: make(map[string]*meta.Type),
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

// Perm mocks a root getter interface
func (mm *TreeMemoryMock) Perm() tree.GetInterface {
	return (*LookupMemManager)(mm)
}

// Apps returns manager of DApps
func (mm *TreeMemoryMock) Apps() tree.AppMemory {
	return &mm.app
}

// Channels returns manager's DApp
func (mm *TreeMemoryMock) Channels() tree.ChannelMemory {
	return &mm.channel
}

// Types returns manager's DApp
func (mm *TreeMemoryMock) Types() tree.TypeMemory {
	return &mm.insprType
}

// Alias returns manager's Alias
func (mm *TreeMemoryMock) Alias() tree.AliasMemory {
	return &mm.alias
}

//InitTransaction mock interface structure
func (mm *TreeMemoryMock) InitTransaction() {}

//Commit mock interface structure
func (mm *TreeMemoryMock) Commit() {}

//Cancel mock interface structure
func (mm *TreeMemoryMock) Cancel() {}

//GetTransactionChanges mock interface structure
func (mm *TreeMemoryMock) GetTransactionChanges() (diff.Changelog, error) {
	return diff.Changelog{}, nil
}
