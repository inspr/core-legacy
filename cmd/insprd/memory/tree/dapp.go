package tree

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	apimodels "inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	metautils "inspr.dev/inspr/pkg/meta/utils"
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

	l.Debug("trying to resolve dApp boundaries and updating connected dApps to resolved Channels and Routes")
	appErr = amm.recursiveBoundaryValidation(app)
	if appErr != nil {
		l.Debug("unable to resolve dApps boundaries - refusing request")
		return appErr
	}

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

	l.Debug("checking if some of the dapp resource is used by the parent")
	if amm.isAppUsed(app, parent) {
		l.Debug("unable to delete dapp for it's being used")
		return ierrors.New(
			"dapp cannot be deleted as it is being used by the parent apps",
		).BadRequest()
	}

	l.Debug("removing dApp from its parent")
	delete(parent.Spec.Apps, app.Meta.Name)

	return nil
}

// isAppUsed checks if the current dapp has some alias, channel or route that is being used by
// the parent dapp
func (amm *AppMemoryManager) isAppUsed(app, parent *meta.App) bool {
	for _, alias := range parent.Spec.Aliases {
		if alias.Source == app.Meta.Name {
			for chName := range app.Spec.Channels {
				if alias.Resource == chName {
					return true
				}
			}

			for routeName := range app.Spec.Routes {
				if alias.Resource == routeName {
					return true
				}
			}

			for aliasName := range app.Spec.Aliases {
				if alias.Resource == aliasName {
					return true
				}
			}
		}
	}

	return false
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

func (amm *AppMemoryManager) ResolveBoundary(app *meta.App, usePermTree bool) (map[string]string, map[string]string, error) {
	parent, _ := amm.GetParent(app, usePermTree)

	merr := ierrors.MultiError{
		Errors: []error{},
	}

	resolvedRoutes := make(map[string]string)
	boundaries := app.Spec.Boundary.Routes
	for _, bound := range boundaries {
		resolvedBound, err := amm.recursivelyResolveUp(parent, app.Meta.Name, bound, usePermTree, "route")
		if err != nil {
			merr.Add(ierrors.Wrap(err, fmt.Sprintf("invalid route boundary: %s invalid", bound)))
		}
		resolvedRoutes[bound] = resolvedBound
	}

	resolvedChannels := make(map[string]string)
	boundaries = app.Spec.Boundary.Channels.Input.Union(app.Spec.Boundary.Channels.Output)
	for _, bound := range boundaries {
		resolvedBound, err := amm.recursivelyResolveUp(parent, app.Meta.Name, bound, usePermTree, "channel")
		if err != nil {
			merr.Add(ierrors.Wrap(err, fmt.Sprintf("invalid channel boundary: %s invalid", bound)))
		}
		resolvedChannels[bound] = resolvedBound
	}

	if !merr.Empty() {
		return nil, nil, &merr
	}

	return resolvedRoutes, resolvedChannels, nil

}

func (amm *AppMemoryManager) recursivelyResolveUp(app *meta.App, requester string, resource string, usePermTree bool, resourceType string) (string, error) {
	scope, ok := amm.checkForResource(app, resource, resourceType)
	if ok {
		return scope, nil
	}

	alias, ok := app.Spec.Aliases[resource]

	if ok && alias.Destination == requester {
		if alias.Source == "" {
			parent, err := amm.GetParent(app, usePermTree)
			if err != nil {
				return "", err
			}
			return amm.recursivelyResolveUp(parent, app.Meta.Name, alias.Resource, usePermTree, resourceType)

		} else {
			// THIS IS THE MEDIUM POINT -> IT ONLY PASSES ONE TIME HERE
			if child, ok := app.Spec.Apps[alias.Source]; ok {
				return amm.recursivelyResolveDown(child, app.Meta.Name, alias.Resource, resourceType)
			}
		}

	}

	return "", ierrors.New("cannot find resource %v", resource)
}

func (amm *AppMemoryManager) recursivelyResolveDown(app *meta.App, requester string, resource string, resourceType string) (string, error) {
	scope, ok := amm.checkForResource(app, resource, resourceType)
	if ok {
		return scope, nil
	}

	alias, ok := app.Spec.Aliases[resource]

	if ok && alias.Destination != "" {
		return "", ierrors.New("cannot find resource %v", resource)
	}

	if ok && alias.Source != "" {
		if child, ok := app.Spec.Apps[alias.Source]; ok {
			return amm.recursivelyResolveDown(child, app.Meta.Name, alias.Resource, resourceType)
		}
	}

	return "", ierrors.New("cannot find resource %v", resource)

}

func (amm *AppMemoryManager) checkForResource(app *meta.App, resource, resourceType string) (string, bool) {
	if resourceType == "route" {
		if route, ok := app.Spec.Routes[resource]; ok {
			scope, _ := metautils.JoinScopes(app.Meta.Parent, app.Meta.Name)
			scope, _ = metautils.JoinScopes(scope, route.Meta.Name)
			return scope, true
		}
	}

	if resourceType == "channel" {
		if channel, ok := app.Spec.Channels[resource]; ok {
			scope, _ := metautils.JoinScopes(app.Meta.Parent, app.Meta.Name)
			scope, _ = metautils.JoinScopes(scope, channel.Meta.Name)
			return scope, true
		}
	}

	return "", false
}

func (amm *AppMemoryManager) GetParent(app *meta.App, usePermTree bool) (*meta.App, error) {
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

	return parentApp, nil
}
