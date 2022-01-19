package tree

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	apimodels "inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	metautils "inspr.dev/inspr/pkg/meta/utils"
	"inspr.dev/inspr/pkg/utils"
)

// AppMemoryManager implements the App interface
// and provides methos for operating on dApps
type AppMemoryManager struct {
	*treeMemoryManager
	logger *zap.Logger
}

// Apps is a MemoryManager method that provides an access point for Apps
func (tmm *treeMemoryManager) Apps() AppMemory {
	logger.Debug("recovering dapp memory manager")
	return &AppMemoryManager{
		treeMemoryManager: tmm,
		logger:            logger.With(zap.String("subSection", "dapp")),
	}
}

// Get receives a query string (format = 'x.y.z') and iterates through the
// memory tree until it finds the dApp which name is equal to the last query element.
// The root dApp is returned if the query string is an empty string.
// If the specified dApp is found, it is returned. Otherwise, returns an error.
func (amm *AppMemoryManager) Get(query string) (*meta.App, error) {
	l := amm.logger.With(zap.String("operation", "get"), zap.String("query", query))
	l.Debug("received dapp get request")

	if amm.root == nil {
		l.Error("root of the tree is nil, no transaction started")
		return nil, ierrors.New(
			"root of the tree is nil, no transaction started",
		).InternalServer()
	}
	if query == "" {
		return amm.root, nil
	}

	err := ierrors.New("dApp not found for given query: %s", query).NotFound()

	reference := strings.Split(query, ".")
	nextApp := amm.root
	for _, element := range reference {
		if nextApp.Spec.Apps == nil {
			l.Debug("dApp for given query is null", zap.String("app", element))
			return nil, err
		}
		nextApp = nextApp.Spec.Apps[element]
		if nextApp == nil {
			l.Debug("unable to find dApp for given query", zap.String("app", element))
			return nil, err
		}
	}
	l.Debug("dapp found")
	return nextApp, nil
}

// Create instantiates a new dApp in the given scope.
// If the dApp's information is invalid, returns an error. The same goes for an invalid scope.
// In case of scope being an empty string, the dApp is created inside the root dApp.
func (amm *AppMemoryManager) Create(scope string, app *meta.App, brokers *apimodels.BrokersDI) error {
	l := amm.logger.With(
		zap.String("operation", "create"),
		zap.String("dApp", app.Meta.Name),
		zap.String("scope", scope),
	)
	l.Debug("received dapp creation request")

	parentApp, err := amm.Get(scope)
	if err != nil {
		return err
	}

	if _, ok := parentApp.Spec.Apps[app.Meta.Name]; ok {
		l.Debug("dapp already exists - refusing request")
		return ierrors.New(
			"this app already exists in parentApp",
		).InvalidApp()
	}

	l.Debug("checking dApp structure")
	appErr := amm.checkApp(app, parentApp, brokers)
	if appErr != nil {
		l.Debug("dapp invalid - refusing request")
		return appErr
	}

	l.Debug("adding dApp to the memory tree")
	amm.addAppInTree(app, parentApp)

	l.Debug("trying to resolve dApp boundaries")
	appErr = amm.recursiveBoundaryValidation(app)
	if appErr != nil {
		l.Debug("unable to resolve dApps boundaries - refusing request")
		return appErr
	}

	l.Debug("updating connected dApps and Aliases to resolved Channels")
	amm.connectAppsBoundaries(app)
	return nil
}

// Delete receives a query and searches for the specified dApp through the tree.
// If the dApp is found and it doesn't have any dApps insite of it, it's deleted.
// If it has other dApps inside of itself, those dApps are deleted recursively.
// Channels and Types inside the dApps to be deleted are also deleted
// dApp's reference inside of it's parent is also deleted.
// In case of dApp not found an error is returned.
func (amm *AppMemoryManager) Delete(query string) error {
	l := amm.logger.With(
		zap.String("operation", "delete"),
		zap.String("query", query),
	)
	l.Debug("trying to delete a dApp", zap.String("dApp query", query))

	if query == "" {
		l.Debug("unable to delete root dApp - refusing request")
		return ierrors.New("can't delete root dApp").BadRequest()
	}

	l.Debug("getting dApp to be deleted")
	app, err := amm.Get(query)
	if err != nil {
		return err
	}
	parent, errParent := getParentApp(query, amm.treeMemoryManager)
	if errParent != nil {
		return errParent
	}

	l.Debug("updating Channels to which the dApp was connected")

	amm.removeFromParentBoundary(app, parent)

	l.Debug("removing dApp from its parent")

	delete(parent.Spec.Apps, app.Meta.Name)

	return nil
}

// Update receives a pointer to a dApp and the path to where this dApp is inside the memory tree.
// If the current dApp is found and the new structure is valid, it's updated.
// Otherwise, returns an error.
func (amm *AppMemoryManager) Update(query string, app *meta.App, brokers *apimodels.BrokersDI) error {
	l := amm.logger.With(
		zap.String("operation", "update"),
		zap.String("dapp", app.Meta.Name),
		zap.String("scope", query),
	)
	l.Debug("received request for dapp update")

	l.Debug("getting dApp to be updated")
	currentApp, err := amm.Get(query)
	if err != nil {
		return err
	}

	l.Debug("validating new dApp structure")
	if currentApp.Meta.Name != app.Meta.Name {
		l.Debug("invalid name change operation", zap.String("old-dapp", currentApp.Meta.Name))
		return ierrors.New(
			"dApp's name mustn't change when updating",
		).InvalidName()
	}
	if !nodeIsEmpty(app.Spec.Node) && !(len(app.Spec.Apps) == 0) {
		l.Debug("a Node can't contain child dApps")
		return ierrors.New(
			"dApp mustn't have a Node and other dApps at the same time",
		).InvalidApp()
	}

	parent, errParent := getParentApp(query, amm.treeMemoryManager)
	if errParent != nil {
		return errParent
	}

	appErr := amm.checkApp(app, parent, brokers)
	if appErr != nil {
		l.Debug("unable to update dApp - invalid structure")
		return appErr
	}

	l.Debug("deleting old dApp")
	amm.removeFromParentBoundary(app, parent)

	delete(parent.Spec.Apps, currentApp.Meta.Name)

	l.Debug("creating new dApp")
	amm.addAppInTree(app, parent)

	return nil
}

// AppPermTreeGetter returns a getter that gets apps from the root structure of the app, without the current changes.
// The getter does not allow changes in the structure, just visualization.
type AppPermTreeGetter struct {
	tree *meta.App
	logs *zap.Logger
}

// Get receives a query string (format = 'x.y.z') and iterates through the
// memory tree until it finds the dApp which name is equal to the last query element.
// The tree root dApp is returned if the query string is an empty string.
// If the specified dApp is found, it is returned. Otherwise, returns an error.
// This method is used to get the structure as it is in the cluster, before any modifications.
func (amm *AppPermTreeGetter) Get(query string) (*meta.App, error) {
	l := amm.logs.With(zap.String("operation", "perm-get"), zap.String("query", query))
	l.Debug("received request for dapp retrieve")

	if amm.tree == nil {
		l.Error("root of the tree is nil, no transaction started")
		return nil, ierrors.New(
			"root of the tree is nil, no transaction started",
		).InternalServer()
	}
	if query == "" {
		return amm.tree, nil
	}

	reference := strings.Split(query, ".")
	err := ierrors.New(
		"dApp not found for given query '%v'", query,
	).NotFound()

	nextApp := amm.tree
	for _, element := range reference {
		nextApp = nextApp.Spec.Apps[element]
		if nextApp == nil {
			l.Debug("unable to find dApp for given query")
			return nil, err
		}
	}

	l.Debug("dapp found")
	return nextApp, nil

}

// ResolveBoundary is the recursive method that resolves connections for dApp boundaries
// returns a map of boundary to  their respective resolved channel query
func (amm *AppMemoryManager) ResolveBoundary(app *meta.App, usePermTree bool) (map[string]string, error) {
	l := amm.logger.With(zap.String("operation", "boundary-resolution"), zap.String("dapp", app.Meta.Name))
	l.Debug("received boundary resolution request", zap.Bool("useperm", usePermTree))

	boundaries := make(map[string]string)
	unresolved := metautils.StrSet{}
	for _, bound := range app.Spec.Boundary.Input.Union(app.Spec.Boundary.Output) {
		boundaries[bound] = fmt.Sprintf("%s.%s", app.Meta.Name, bound)
		unresolved[bound] = true
	}

	var parentApp *meta.App

	if usePermTree {
		parApp, err := amm.Perm().Apps().Get(app.Meta.Parent)
		if err != nil {
			return nil, err
		}
		parentApp = parApp

	} else {
		parApp, err := amm.treeMemoryManager.Apps().Get(app.Meta.Parent)
		if err != nil {
			return nil, err
		}
		parentApp = parApp
	}

	l.Debug("recursively resolving dApp boundaries",
		zap.Any("boundaries", boundaries))

	err := amm.recursivelyResolve(parentApp, boundaries, unresolved, usePermTree)
	if err != nil {
		l.Debug("couldn't resolve boundaries for given dApp",
			zap.String("parent", parentApp.Meta.Name),
			zap.Any("boundaries", boundaries))
		return nil, err
	}
	return boundaries, nil
}

func (amm *AppMemoryManager) recursivelyResolve(app *meta.App, boundaries map[string]string, unresolved metautils.StrSet, usePermTree bool) error {
	_ = amm.logger.With(zap.String("operation", "boundary-resolution"))
	merr := ierrors.MultiError{
		Errors: []error{},
	}
	if len(unresolved) == 0 {
		return nil
	}
	for key := range unresolved {
		val := boundaries[key]
		if alias, ok := app.Spec.Aliases[val]; ok { //resolve in aliases
			val = alias.Resource //setup for alias resolve
		} else {
			_, val, _ = metautils.RemoveLastPartInScope(val) //setup for direct resolve
		}
		if ch, ok := app.Spec.Channels[val]; ok { // resolve in channels (direct or through alias)
			scope, _ := metautils.JoinScopes(app.Meta.Parent, app.Meta.Name)
			boundaries[key], _ = metautils.JoinScopes(scope, ch.Meta.Name) // if channel exists, resolve
			delete(unresolved, key)
			continue
		}
		if app.Spec.Boundary.Input.Union(app.Spec.Boundary.Output).Contains(val) { //resolve in boundaries
			boundaries[key], _ = metautils.JoinScopes(app.Meta.Name, val) // if boundary exists, setup to resolve in parernt
			continue
		}
		merr.Add(ierrors.New("invalid boundary: %s invalid", key))
		delete(unresolved, key)

	}
	if !merr.Empty() {
		// throwing erros for boundaries couldn't be resolved because of some invalid boundary
		for key := range unresolved {
			merr.Add(ierrors.New("invalid boundary: %s unresolved", key))
		}
		return &merr
	}

	var parentApp *meta.App

	if usePermTree {
		parApp, err := amm.Perm().Apps().Get(app.Meta.Parent)
		if err != nil {
			return err
		}
		parentApp = parApp

	} else {
		parApp, err := amm.treeMemoryManager.Apps().Get(app.Meta.Parent)
		if err != nil {
			return err
		}
		parentApp = parApp
	}

	return amm.recursivelyResolve(parentApp, boundaries, unresolved, usePermTree)
}

func (amm *AppMemoryManager) removeFromParentBoundary(app, parent *meta.App) {
	l := amm.logger.With(
		zap.String("operation", "boundary-removal"),
		zap.String("dApp", app.Meta.Name),
		zap.String("parent", parent.Meta.Name),
	)
	l.Debug("removing dApp from parent's Channels connected apps list")

	appBoundary := utils.StringSliceUnion(app.Spec.Boundary.Input, app.Spec.Boundary.Output)
	resolution, _ := amm.ResolveBoundary(app, false)
	for _, chName := range appBoundary {
		resolved := resolution[chName]
		_, chName, _ := metautils.RemoveLastPartInScope(resolved)
		if _, ok := parent.Spec.Channels[chName]; ok {
			parent.Spec.Channels[chName].ConnectedApps = utils.
				Remove(parent.Spec.Channels[chName].ConnectedApps, app.Meta.Name)
		}
	}
}
