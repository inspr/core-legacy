package tree

import (
	"fmt"
	"strings"

	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
	"github.com/inspr/inspr/pkg/utils"
	"go.uber.org/zap"
)

// AppMemoryManager implements the App interface
// and provides methos for operating on dApps
type AppMemoryManager struct {
	*MemoryManager
}

// Apps is a MemoryManager method that provides an access point for Apps
func (tmm *MemoryManager) Apps() memory.AppMemory {
	return &AppMemoryManager{
		MemoryManager: tmm,
	}
}

// Get receives a query string (format = 'x.y.z') and iterates through the
// memory tree until it finds the dApp which name is equal to the last query element.
// The root dApp is returned if the query string is an empty string.
// If the specified dApp is found, it is returned. Otherwise, returns an error.
func (amm *AppMemoryManager) Get(query string) (*meta.App, error) {
	logger.Info("trying to get a dApp", zap.String("dApp", query))

	if query == "" {
		return amm.root, nil
	}

	reference := strings.Split(query, ".")
	err := ierrors.
		NewError().
		NotFound().
		Message("dApp not found for given query: %s", query).
		Build()

	nextApp := amm.root
	if nextApp != nil {
		for _, element := range reference {
			nextApp = nextApp.Spec.Apps[element]
			if nextApp == nil {
				logger.Error("unable to find dApp for given query",
					zap.String("query", query))
				return nil, err
			}
		}
		return nextApp, nil
	}

	logger.Error("root dApp is empty")
	return nil, err
}

// Create instantiates a new dApp in the given context.
// If the dApp's information is invalid, returns an error. The same goes for an invalid context.
// In case of context being an empty string, the dApp is created inside the root dApp.
func (amm *AppMemoryManager) Create(context string, app *meta.App) error {
	logger.Info("trying to create a dApp",
		zap.String("dApp", app.Meta.Name),
		zap.String("context", context))

	parentApp, err := amm.Get(context)
	if err != nil {
		return err
	}

	if _, ok := parentApp.Spec.Apps[app.Meta.Name]; ok {
		logger.Error("unable to create dApp for it already exists")
		return ierrors.NewError().InvalidApp().Message("this app already exists in parentApp").Build()
	}

	logger.Debug("checking dApp structure")
	appErr := amm.checkApp(app, parentApp)
	if appErr != nil {
		logger.Error("unable to create dApp - invalid structure")
		return appErr
	}

	logger.Debug("adding dApp to the memory tree")
	amm.addAppInTree(app, parentApp)

	logger.Debug("trying to resolve dApp boundaries")
	appErr = amm.recursiveBoundaryValidation(app)
	if appErr != nil {
		logger.Error("unable to resolve dApps boundaries")
		return appErr
	}

	logger.Debug("updating connected dApps and Aliases to resolved Channels")
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
	logger.Info("trying to delete a dApp", zap.String("dApp query", query))

	if query == "" {
		logger.Error("unable to delete root dApp")
		return ierrors.NewError().BadRequest().Message("can't delete root dApp").Build()
	}

	logger.Debug("getting dApp to be deleted")
	app, err := amm.Get(query)
	if err != nil {
		return err
	}
	parent, errParent := getParentApp(query)
	if errParent != nil {
		return errParent
	}

	logger.Debug("updating Channels to which the dApp was connected")

	amm.removeFromParentBoundary(app, parent)

	logger.Debug("removing dApp from its parents 'Apps' structure",
		zap.String("dApp", app.Meta.Name),
		zap.String("parent dApp", parent.Meta.Name))

	delete(parent.Spec.Apps, app.Meta.Name)

	return nil
}

// Update receives a pointer to a dApp and the path to where this dApp is inside the memory tree.
// If the current dApp is found and the new structure is valid, it's updated.
// Otherwise, returns an error.
func (amm *AppMemoryManager) Update(query string, app *meta.App) error {
	logger.Info("trying to update a dApp",
		zap.String("dApp", app.Meta.Name),
		zap.String("in context", query))

	logger.Debug("getting dApp to be updated")
	currentApp, err := amm.Get(query)
	if err != nil {
		return err
	}

	logger.Debug("validating new dApp structure")
	if currentApp.Meta.Name != app.Meta.Name {
		logger.Error("unable to change a dApps name when updating it")
		return ierrors.NewError().InvalidName().Message("dApp's name mustn't change when updating").Build()
	}
	if !nodeIsEmpty(app.Spec.Node) && !(len(app.Spec.Apps) == 0) {
		logger.Error("a Node can't contain child dApps")
		return ierrors.NewError().InvalidApp().Message("dApp mustn't have a Node and other dApps at the same time").Build()
	}

	parent, errParent := getParentApp(query)
	if errParent != nil {
		return errParent
	}

	appErr := amm.checkApp(app, parent)
	if appErr != nil {
		logger.Error("unable to update dApp - invalid structure")
		return appErr
	}

	logger.Debug("deleting old dApp")
	amm.removeFromParentBoundary(app, parent)

	delete(parent.Spec.Apps, currentApp.Meta.Name)

	logger.Debug("creating new dApp")
	amm.addAppInTree(app, parent)

	return nil
}

// AppRootGetter returns a getter that gets apps from the root structure of the app, without the current changes.
// The getter does not allow changes in the structure, just visualization.
type AppRootGetter struct {
	tree *meta.App
}

// Get receives a query string (format = 'x.y.z') and iterates through the
// memory tree until it finds the dApp which name is equal to the last query element.
// The tree root dApp is returned if the query string is an empty string.
// If the specified dApp is found, it is returned. Otherwise, returns an error.
// This method is used to get the structure as it is in the cluster, before any modifications.
func (amm *AppRootGetter) Get(query string) (*meta.App, error) {
	logger.Info("trying to get a dApp (Root Getter)", zap.String("dApp", query))

	if query == "" {
		return amm.tree, nil
	}

	reference := strings.Split(query, ".")
	err := ierrors.NewError().NotFound().Message("dApp not found for given query: " + query).Build()

	nextApp := amm.tree
	if nextApp != nil {
		for _, element := range reference {
			nextApp = nextApp.Spec.Apps[element]
			if nextApp == nil {
				logger.Error("unable to find dApp for given query (Root Getter)",
					zap.String("query", query))
				return nil, err
			}
		}
		return nextApp, nil
	}

	logger.Error("root dApp is empty (Root Getter)")
	return nil, err
}

//ResolveBoundary is the recursive method that resolves connections for dApp boundaries
func (amm *AppMemoryManager) ResolveBoundary(app *meta.App) (map[string]string, error) {
	logger.Debug("resolving dApp boundary",
		zap.String("dApp", app.Meta.Name))

	boundaries := make(map[string]string)
	unresolved := metautils.StrSet{}
	for _, bound := range app.Spec.Boundary.Input.Union(app.Spec.Boundary.Output) {
		boundaries[bound] = fmt.Sprintf("%s.%s", app.Meta.Name, bound)
		unresolved[bound] = true
	}
	parentApp, err := amm.MemoryManager.Apps().Get(app.Meta.Parent)
	if err != nil {
		return nil, err
	}

	logger.Debug("recursively resolving dApp boundaries",
		zap.String("dApp", app.Meta.Name),
		zap.Any("boundaries", boundaries))

	err = amm.recursivelyResolve(parentApp, boundaries, unresolved)
	if err != nil {
		logger.Error("couldn't resolve boundaries for given dApp",
			zap.String("dApp", app.Meta.Name),
			zap.String("parent dApp", parentApp.Meta.Name),
			zap.Any("boundaries", boundaries))
		return nil, err
	}
	return boundaries, nil
}

func (amm *AppMemoryManager) recursivelyResolve(app *meta.App, boundaries map[string]string, unresolved metautils.StrSet) error {
	merr := ierrors.MultiError{
		Errors: []error{},
	}
	if len(unresolved) == 0 {
		return nil
	}
	for key := range unresolved {
		val := boundaries[key]
		if alias, ok := app.Spec.Aliases[val]; ok { //resolve in aliases
			val = alias.Target //setup for alias resolve
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
		merr.Add(ierrors.NewError().Message("invalid boundary: %s invalid", key).Build())
		delete(unresolved, key)

	}
	if !merr.Empty() {
		// throwing erros for boundaries couldn't be resolved because of some invalid boundary
		for key := range unresolved {
			merr.Add(ierrors.NewError().Message("invalid boundary: %s unresolved", key).Build())
		}
		return &merr
	}
	parentApp, err := amm.MemoryManager.Apps().Get(app.Meta.Parent)
	if err != nil {
		return err
	}
	return amm.recursivelyResolve(parentApp, boundaries, unresolved)
}

func (amm *AppMemoryManager) removeFromParentBoundary(app, parent *meta.App) {
	logger.Debug("removing dApp from parent's Channels connected apps list",
		zap.String("dApp", app.Meta.Name),
		zap.String("parent", parent.Meta.Name))

	appBoundary := utils.StringSliceUnion(app.Spec.Boundary.Input, app.Spec.Boundary.Output)
	resolution, _ := amm.ResolveBoundary(app)
	for _, chName := range appBoundary {
		resolved := resolution[chName]
		_, chName, _ := metautils.RemoveLastPartInScope(resolved)
		if _, ok := parent.Spec.Channels[chName]; ok {
			parent.Spec.Channels[chName].ConnectedApps = utils.
				Remove(parent.Spec.Channels[chName].ConnectedApps, app.Meta.Name)
		}
	}
}
