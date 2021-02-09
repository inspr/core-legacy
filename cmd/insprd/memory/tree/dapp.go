package tree

import (
	"strings"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

// AppMemoryManager implements the App interface
// and provides methos for operating on dApps
type AppMemoryManager struct {
	root *meta.App
}

// Apps is a MemoryManager method that provides an access point for Apps
func (tmm *MemoryManager) Apps() memory.AppMemory {
	return &AppMemoryManager{
		root: tmm.root,
	}
}

// GetApp recieves a query string (format = 'x.y.z') and iterates through the
// memory tree until it finds the dApp which name is equal to the last query element.
// The root app is returned if the query string is an empty string.
// If the specified dApp is found, it is returned. Otherwise, returns an error.
func (amm *AppMemoryManager) GetApp(query string) (*meta.App, error) {
	if query == "" {
		return amm.root, nil
	}

	reference := strings.Split(query, ".")
	err := ierrors.NewError().NotFound().Message("dApp not found for given query " + query).Build()

	nxtApp := amm.root
	if nxtApp != nil {
		for _, element := range reference {
			nxtApp = nxtApp.Spec.Apps[element]
			if nxtApp == nil {
				return nil, err
			}
		}
		return nxtApp, nil
	}

	return nil, err
}

// CreateApp instantiates a new dApp in the given context.
// If the dApp's information is invalid, returns an error. The same goes for an invalid context.
// In case of context being an empty string, the dApp is created inside the root dApp.
func (amm *AppMemoryManager) CreateApp(app *meta.App, context string) error {
	if strings.Contains(app.Meta.Name, ".") {
		ierrors.NewError().InvalidName().Message("invalid character '.' in dApp's name").Build()
	}

	parentApp, err := amm.GetApp(context)
	if err != nil {
		return err
	}

	if _, ok := parentApp.Spec.Apps[app.Meta.Name]; ok {
		return ierrors.NewError().InvalidApp().Message("this app already exists in parentApp").Build()
	}

	appErr := amm.checkApp(app, parentApp)
	if appErr != nil {
		return appErr
	}
	amm.addAppInTree(app, parentApp)
	return nil
}

// DeleteApp receives a query and searches for the specified dApp through the tree.
// If the dApp is found and it doesn't have any dApps insite of it, it's deleted.
// If it has other dApps inside of itself, those dApps are deleted recursively.
// Channels and Channel Types inside the dApps to be deleted are also deleted
// dApp's reference inside of it's parent is also deleted.
// In case of dApp not found an error is returned.
func (amm *AppMemoryManager) DeleteApp(query string) error {
	if query == "" {
		return ierrors.NewError().BadRequest().Message("can't delete root dApp").Build()
	}

	app, err := amm.GetApp(query)
	if err != nil {
		return err
	}
	parent, errParent := getParentApp(query)
	if errParent != nil {
		return errParent
	}

	appBoundary := utils.StringSliceUnion(app.Spec.Boundary.Input, app.Spec.Boundary.Output)

	for _, chName := range appBoundary {
		parent.Spec.Channels[chName].ConnectedApps = utils.
			Remove(parent.Spec.Channels[chName].ConnectedApps, app.Meta.Name)
	}

	delete(parent.Spec.Apps, app.Meta.Name)

	return nil
}

// UpdateApp receives a pointer to a dApp and the path to where this dApp is inside the memory tree.
// If the current dApp is found and the new structure is valid, it's updated.
// Otherwise, returns an error.
func (amm *AppMemoryManager) UpdateApp(app *meta.App, query string) error {
	currentApp, err := amm.GetApp(query)
	if err != nil {
		return err
	}

	if currentApp.Meta.Name != app.Meta.Name {
		return ierrors.NewError().InvalidName().Message("dApp's name mustn't change when updating").Build()
	}
	if !nodeIsEmpty(app.Spec.Node) && !(len(app.Spec.Apps) == 0) {
		return ierrors.NewError().InvalidApp().Message("dApp mustn't have a Node and other dApps at the same time").Build()
	}

	parent, errParent := getParentApp(query)
	if errParent != nil {
		return errParent
	}

	appErr := amm.checkApp(app, parent)
	if appErr != nil {
		return appErr
	}

	appBoundary := utils.StringSliceUnion(currentApp.Spec.Boundary.Input, currentApp.Spec.Boundary.Output)

	for _, chName := range appBoundary {
		parent.Spec.Channels[chName].ConnectedApps = utils.
			Remove(parent.Spec.Channels[chName].ConnectedApps, currentApp.Meta.Name)
	}

	delete(parent.Spec.Apps, currentApp.Meta.Name)

	amm.addAppInTree(app, parent)

	return nil
}
