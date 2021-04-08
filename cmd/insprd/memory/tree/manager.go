package tree

import (
	"sync"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
	"go.uber.org/zap"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = zap.NewDevelopment(zap.Fields(zap.String("section", "memory-tree")))
}

// MemoryManager defines a memory manager interface
type MemoryManager struct {
	root *meta.App
	tree *meta.App
	sync.Mutex
}

var dapptree memory.Manager

// GetTreeMemory returns a memory manager interface
func GetTreeMemory() memory.Manager {
	if dapptree == nil {
		setTree(newTreeMemory())
	}
	return dapptree
}

func newTreeMemory() *MemoryManager {
	return &MemoryManager{
		tree: &meta.App{
			Meta: meta.Metadata{
				Annotations: map[string]string{},
			},
			Spec: meta.AppSpec{
				Apps:         map[string]*meta.App{},
				Channels:     map[string]*meta.Channel{},
				ChannelTypes: map[string]*meta.ChannelType{},
				Aliases:      map[string]*meta.Alias{},
			},
		},
	}
}

func setTree(tmm memory.Manager) {
	dapptree = tmm
}

//InitTransaction copies and reserves the current tree structure so that changes can be reversed
func (mm *MemoryManager) InitTransaction() {
	mm.Lock()
	utils.DeepCopy(mm.tree, &mm.root)
}

//Commit applies changes from a transaction in to the tree structure
func (mm *MemoryManager) Commit() {
	defer mm.Unlock()
	mm.tree = mm.root
	mm.root = nil
}

//Cancel discarts changes made in the last transaction
func (mm *MemoryManager) Cancel() {
	defer mm.Unlock()
	mm.root = nil
}

//GetTransactionChanges returns the changelog resulting from the current transaction.
func (mm *MemoryManager) GetTransactionChanges() (diff.Changelog, error) {
	cl, err := diff.Diff(mm.tree, mm.root)
	return cl, err
}

// RootGetter is a structure that gets components from the root, without the current changes.
type RootGetter struct {
	tree *meta.App
}

// Apps returns a getter for apps on the root.
func (t *RootGetter) Apps() memory.AppGetInterface {
	return &AppRootGetter{
		tree: t.tree,
	}
}

// Channels returns a getter for channels on the root.
func (t *RootGetter) Channels() memory.ChannelGetInterface {
	return &ChannelRootGetter{}
}

// ChannelTypes returns a getter for channel types on the root
func (t *RootGetter) ChannelTypes() memory.ChannelTypeGetInterface {
	return &ChannelTypeRootGetter{}
}

// Alias returns a getter for alias on the root
func (t *RootGetter) Alias() memory.AliasGetInterface {
	return &AliasRootGetter{}
}

// Root returns a getter for objects on the root of the tree, without the current changes.
func (mm *MemoryManager) Root() memory.GetInterface {
	return &RootGetter{
		tree: mm.tree,
	}
}
