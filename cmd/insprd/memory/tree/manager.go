package tree

import (
	"sync"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/logs"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils"
	"inspr.dev/inspr/pkg/meta/utils/diff"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = logs.Logger(zap.Fields(zap.String("section", "memory-tree")))
}

// treeMemoryManager defines a memory manager interface
type treeMemoryManager struct {
	root *meta.App
	tree *meta.App
	sync.Mutex
}

var dapptree *treeMemoryManager

// GetTreeMemory returns a tree memory manager interface
func GetTreeMemory() Manager {
	if dapptree == nil {
		setTree(newTreeMemory())
	}
	logger.Debug("getting singleton tree memory manager")
	return dapptree
}

func newTreeMemory() *treeMemoryManager {
	logger.Info("initializing memory tree")
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

func setTree(tmm *treeMemoryManager) {
	dapptree = tmm
}

//InitTransaction copies and reserves the current tree structure so that changes can be reversed
func (tmm *treeMemoryManager) InitTransaction() {
	tmm.Lock()
	utils.DeepCopy(tmm.tree, &tmm.root)
}

//Commit applies changes from a transaction in to the tree structure
func (tmm *treeMemoryManager) Commit() {
	defer tmm.Unlock()
	tmm.tree = tmm.root
	tmm.root = nil
}

//Cancel discarts changes made in the last transaction
func (tmm *treeMemoryManager) Cancel() {
	defer tmm.Unlock()
	tmm.root = nil
}

//GetTransactionChanges returns the changelog resulting from the current transaction.
func (tmm *treeMemoryManager) GetTransactionChanges() (diff.Changelog, error) {
	cl, err := diff.Diff(tmm.tree, tmm.root)
	return cl, err
}

// PermTreeGetter is a structure that gets components from the root, without the current changes.
type PermTreeGetter struct {
	tree   *meta.App
	logger *zap.Logger
}

// Apps returns a getter for apps on the root.
func (ptg *PermTreeGetter) Apps() AppGetInterface {
	return &AppPermTreeGetter{
		tree: ptg.tree,
		logs: logger,
	}
}

// Channels returns a getter for channels on the root.
func (ptg *PermTreeGetter) Channels() ChannelGetInterface {
	return &ChannelPermTreeGetter{
		PermTreeGetter: ptg,
		logs:           logger,
	}
}

// Types returns a getter for Types on the root
func (ptg *PermTreeGetter) Types() TypeGetInterface {
	return &TypePermTreeGetter{
		PermTreeGetter: ptg,
		logs:           logger,
	}
}

// Alias returns a getter for alias on the root
func (ptg *PermTreeGetter) Alias() AliasGetInterface {
	return &AliasPermTreeGetter{
		PermTreeGetter: ptg,
		logs:           logger,
	}
}

// Perm returns a getter for objects on the tree without the current changes.
func (tmm *treeMemoryManager) Perm() GetInterface {
	return &PermTreeGetter{
		tree:   tmm.tree,
		logger: logger,
	}
}
