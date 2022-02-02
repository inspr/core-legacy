package tree

import (
	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils"
)

// AliasMemoryManager implements the Alias interface
// and provides methods for operating on Aliases
type AliasMemoryManager struct {
	*treeMemoryManager
	logger *zap.Logger
}

// Alias is a MemoryManager method that provides an access point for Alias
func (tmm *treeMemoryManager) Alias() AliasMemory {
	return &AliasMemoryManager{
		treeMemoryManager: tmm,
		logger:            logger.With(zap.String("subSection", "alias")),
	}
}

// Get receives a scope and an alias key. The scope defines
// the path to a dApp. If this dApp has a pointer to a alias that has the
// same key as the key passed as an argument, the pointer to that alias is returned
func (amm *AliasMemoryManager) Get(scope, name string) (*meta.Alias, error) {
	l := amm.logger.With(
		zap.String("operation", "get"),
		zap.String("alias", name),
		zap.String("scope", scope),
	)
	l.Debug("trying to get an Alias")

	app, err := amm.treeMemoryManager.Apps().Get(scope)
	if err != nil {
		l.Debug("unable to get Alias because the app was not found")
		return nil, err
	}

	if _, ok := app.Spec.Aliases[name]; !ok {
		l.Debug("alias not found with the given name")
		return nil, ierrors.New(
			"alias not found with the given name %v", name,
		).NotFound()
	}

	return app.Spec.Aliases[name], nil
}

// Create receives a scope that defines a path to the dapp and an Alias
// to be added in this dapp
func (amm *AliasMemoryManager) Create(scope string, alias *meta.Alias) error {
	l := amm.logger.With(
		zap.String("operation", "create"),
		zap.Any("alias", alias),
		zap.String("scope", scope),
	)
	l.Debug("trying to create an Alias")

	app, err := amm.Apps().Get(scope)
	if err != nil {
		l.Debug("unable to create Alias because the app was not found")
		return err
	}

	if _, ok := app.Spec.Aliases[alias.Meta.Name]; ok {
		l.Debug("alias already exists")
		return ierrors.New("alias already exists in dApp").AlreadyExists()
	}

	err = amm.CheckSource(scope, app, alias)
	if err != nil {
		return err
	}

	err = amm.CheckDestination(app, alias)
	if err != nil {
		return err
	}

	alias.Meta = utils.InjectUUID(alias.Meta)
	l.Debug("adding Alias to dApp",
		zap.String("dApp", app.Meta.Name))
	app.Spec.Aliases[alias.Meta.Name] = alias

	return nil
}

// Update receives a scope and a alias object. The scope
// defines the path to the dApp that contains the Alias. If the dApp has
// an alias that has the same name as the one passed as an argument,
// that alias will be replaced by the new alias
func (amm *AliasMemoryManager) Update(scope string, alias *meta.Alias) error {
	l := amm.logger.With(
		zap.String("operation", "update"),
		zap.Any("alias", alias),
		zap.String("scope", scope),
	)
	l.Debug("trying to update an Alias")

	app, err := amm.treeMemoryManager.Apps().Get(scope)
	if err != nil {
		l.Debug("unable to update Alias because the app was not found")
		return err
	}

	selectedAlias, ok := app.Spec.Aliases[alias.Meta.Name]
	if !ok {
		l.Debug("alias was not found")
		return ierrors.New("alias was not found in dApp").NotFound()
	}

	err = amm.CheckSource(scope, app, alias)
	if err != nil {
		return err
	}

	err = amm.CheckDestination(app, alias)
	if err != nil {
		return err
	}

	alias.Meta.UUID = selectedAlias.Meta.UUID
	l.Debug("replacing old Alias with the new one",
		zap.String("dapp", app.Meta.Name))
	app.Spec.Aliases[alias.Meta.Name] = alias

	return nil

}

// Delete receives a scope and a alias name. The scope
// defines the path to the dApp that cointains the Alias to be deleted. If the dApp
// has an alias that has the same key as the key passed as an argument, that alias
// is removed from the dApp Aliases only if it's not being used
func (amm *AliasMemoryManager) Delete(scope, name string) error {
	l := amm.logger.With(
		zap.String("operation", "delete"),
		zap.Any("alias", name),
		zap.String("scope", scope),
	)
	l.Debug("trying to delete an Alias")

	app, err := amm.treeMemoryManager.Apps().Get(scope)
	if err != nil {
		return ierrors.Wrap(err, "cannot delete alias because the dapp was not found")
	}

	alias, ok := app.Spec.Aliases[name]
	if !ok {
		return ierrors.New(
			"alias not found with the given name %v", name,
		).NotFound()
	}

	if amm.isAliasUsed(app, alias) {
		l.Debug("unable to delete Alias for it's being used")
		return ierrors.New(
			"alias cannot be deleted as it is being used by other apps",
		).BadRequest()
	}

	l.Debug("removing Alias")
	delete(app.Spec.Aliases, name)

	return nil
}

func (amm *AliasMemoryManager) isAliasUsed(app *meta.App, alias *meta.Alias) bool {
	if alias.Destination != "" {
		child := app.Spec.Apps[alias.Destination]
		if child.Spec.Boundary.Channels.Input.Contains(alias.Meta.Name) || child.Spec.Boundary.Channels.Output.Contains(alias.Meta.Name) {
			return true
		}

		for _, childAlias := range child.Spec.Aliases {
			if childAlias.Resource == alias.Meta.Name {
				return true
			}
		}
	} else {
		parent, _ := amm.Apps().Get(app.Meta.Parent)
		for _, parentAlias := range parent.Spec.Aliases {
			if parentAlias.Source == app.Meta.Name && parentAlias.Resource == alias.Meta.Name {
				return true
			}
		}
	}

	return false
}

// AliasPermTreeGetter returns a getter that gets alias from the root structure of the app, without the current changes.
// The getter does not allow changes in the structure, just visualization.
type AliasPermTreeGetter struct {
	*PermTreeGetter
	logs *zap.Logger
}

func (amm *AliasPermTreeGetter) Get(scope, name string) (*meta.Alias, error) {
	l := amm.logs.With(
		zap.String("operation", "root-get"),
		zap.String("alias", name),
		zap.String("scope", scope),
	)
	l.Debug("trying to get an Alias")

	app, err := amm.Apps().Get(scope)
	if err != nil {
		l.Debug("unable to get Alias because the app was not found")
		return nil, err
	}

	if _, ok := app.Spec.Aliases[name]; !ok {
		l.Debug("alias not found with the given name")
		return nil, ierrors.New(
			"alias not found with the given name %v", name,
		).NotFound()
	}

	return app.Spec.Aliases[name], nil
}

func (amm *AliasMemoryManager) CheckSource(scope string, app *meta.App, alias *meta.Alias) error {
	var source *meta.App
	if alias.Source == "" {
		parentApp, err := getParentApp(scope, amm.treeMemoryManager)
		if err != nil {
			return ierrors.Wrap(err, "cannot find parent dapp")
		}
		source = parentApp
	} else {
		childApp, ok := app.Spec.Apps[alias.Source]
		if !ok {
			return ierrors.New("cannot find source child dapp with the name '%v'", alias.Source).NotFound()
		}
		source = childApp
	}

	_, hasChannel := source.Spec.Channels[alias.Resource]
	_, hasRoute := source.Spec.Routes[alias.Resource]
	if !hasChannel && !hasRoute {
		return ierrors.New("cannot find resource with the name '%v'", alias.Resource).NotFound()
	}

	return nil
}

func (amm *AliasMemoryManager) CheckDestination(app *meta.App, alias *meta.Alias) error {
	if alias.Destination != "" {
		_, ok := app.Spec.Apps[alias.Destination]
		if !ok {
			return ierrors.New("cannot find destination child dapp with the name '%v'", alias.Destination).NotFound()
		}
	}
	return nil
}
