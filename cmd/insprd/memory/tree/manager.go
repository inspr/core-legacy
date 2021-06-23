package tree

import (
	"sync"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils"
	"inspr.dev/inspr/pkg/meta/utils/diff"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = zap.NewProduction(zap.Fields(zap.String("section", "memory-tree")))
	// logger = zap.NewNop()
}

// treeMemoryManager defines a memory manager interface
type treeMemoryManager struct {
	root *meta.App
	tree *meta.App
	sync.Mutex
}

var dapptree Manager

// GetTreeMemory returns a memory manager interface
func GetTreeMemory() Manager {
	if dapptree == nil {
		setTree(newTreeMemory())
	}
	return dapptree
}

func newTreeMemory() *treeMemoryManager {
	return &treeMemoryManager{
		tree: &meta.App{
			Meta: meta.Metadata{
				Annotations: map[string]string{},
			},
			Spec: meta.AppSpec{
				Apps:     map[string]*meta.App{},
				Channels: map[string]*meta.Channel{},
				Types:    map[string]*meta.Type{},
				Aliases:  map[string]*meta.Alias{},
				Auth: meta.AppAuth{
					Scope:       "",
					Permissions: nil,
				},
			},
		},
	}
}

func setTree(tmm Manager) {
	dapptree = tmm
}

//InitTransaction copies and reserves the current tree structure so that changes can be reversed
func (mm *treeMemoryManager) InitTransaction() {
	mm.Lock()
	utils.DeepCopy(mm.tree, &mm.root)
}

//Commit applies changes from a transaction in to the tree structure
func (mm *treeMemoryManager) Commit() {
	defer mm.Unlock()
	mm.tree = mm.root
	mm.root = nil
}

//Cancel discarts changes made in the last transaction
func (mm *treeMemoryManager) Cancel() {
	defer mm.Unlock()
	mm.root = nil
}

//GetTransactionChanges returns the changelog resulting from the current transaction.
func (mm *treeMemoryManager) GetTransactionChanges() (diff.Changelog, error) {
	cl, err := diff.Diff(mm.tree, mm.root)
	return cl, err
}

// RootGetter is a structure that gets components from the root, without the current changes.
type RootGetter struct {
	tree *meta.App
}

// Apps returns a getter for apps on the root.
func (t *RootGetter) Apps() AppGetInterface {
	return &AppRootGetter{
		tree: t.tree,
	}
}

// Channels returns a getter for channels on the root.
func (t *RootGetter) Channels() ChannelGetInterface {
	return &ChannelRootGetter{}
}

// Types returns a getter for Types on the root
func (t *RootGetter) Types() TypeGetInterface {
	return &TypeRootGetter{}
}

// Alias returns a getter for alias on the root
func (t *RootGetter) Alias() AliasGetInterface {
	return &AliasRootGetter{}
}

// Root returns a getter for objects on the root of the tree, without the current changes.
func (mm *treeMemoryManager) Root() GetInterface {
	return &RootGetter{
		tree: mm.tree,
	}
}
