package tree

import (
	"strings"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils"
)

// AliasMemoryManager implements the Alias interface
// and provides methos for operating on Aliass
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
func (amm *AliasMemoryManager) Get(scope, aliasKey string) (*meta.Alias, error) {
	l := amm.logger.With(
		zap.String("operation", "get"),
		zap.String("alias", aliasKey),
		zap.String("scope", scope),
	)
	l.Debug("trying to get an Alias")

	app, err := amm.treeMemoryManager.Apps().Get(scope)
	if err != nil {
		l.Debug("unable to get Alias")
		return nil, err
	}

	// check if alias key exist in scope
	if _, ok := app.Spec.Aliases[aliasKey]; !ok {
		l.Debug("alias not found for the key")
		return nil, ierrors.
			New("alias not found for the given key %v", aliasKey).
			BadRequest()
	}

	//return alias
	return app.Spec.Aliases[aliasKey], nil
}

// Create receives a scope that defines a path to the dApp.
// The new alias will be created inside this dApp's parent.
func (amm *AliasMemoryManager) Create(scope, targetBoundary string, alias *meta.Alias) error {
	l := amm.logger.With(
		zap.String("operation", "create"),
		zap.Any("alias", alias),
		zap.String("scope", scope),
	)
	// get app from scope
	app, err := amm.Apps().Get(scope)
	if err != nil {
		l.Debug("parent app not found")
		return err
	}

	l.Debug("checking if dApp boundary is valid for given Alias")
	// check if targetBoundary exists in app
	appBound := app.Spec.Boundary
	if !appBound.Input.Contains(targetBoundary) && !appBound.Output.Contains(targetBoundary) {
		l.Debug("invalid dApp boundary for Alias",
			zap.String("targeted boundary", targetBoundary))
		return ierrors.New(
			"target boundary doesn't exist in %v",
			app.Meta.Name,
		).BadRequest()
	}

	// get parentApp of app
	parentApp, _ := getParentApp(scope, amm.treeMemoryManager)

	targetChannel := alias.Target

	l.Debug("checking if Alias targeted Channel exists")
	// check if targetChannel exists in channels or boundaries of parentApp
	err = validTargetChannel(parentApp, targetChannel)
	if err != nil {
		l.Debug("unable to create Alias - invalid targeted channel or boundary in parent dApp",
			zap.String("parent", parentApp.Meta.Name),
			zap.String("targeted boundary", targetChannel))
		return err
	}

	aliasKey := app.Meta.Name + "." + targetBoundary

	l.Debug("checking if Alias already exists")
	// check if alias is already there
	if _, ok := parentApp.Spec.Aliases[aliasKey]; ok {
		l.Debug("alias already exists")
		return ierrors.New("alias already exists in parent dApp").BadRequest()
	}

	alias.Meta = utils.InjectUUID(alias.Meta)

	l.Debug("adding Alias to dApp",
		zap.String("dApp", parentApp.Meta.Name))
	// add new alias to Aliases list in parentApp
	parentApp.Spec.Aliases[aliasKey] = alias

	return nil

}

// Update receives a scope a alias key and a alias. The scope
// defines the path to the dApp that contains the Alias. If the dApp has
// a alias that has the given alias key passed as an argument,
// that alias will be replaced by the new alias
func (amm *AliasMemoryManager) Update(scope, aliasKey string, alias *meta.Alias) error {
	l := amm.logger.With(
		zap.String("operation", "update"),
		zap.Any("alias", alias),
		zap.String("scope", scope),
	)
	l.Debug("trying to update an Alias")

	l.Debug("checking if Alias to be updated exists in given scope")
	// check if alias key exist in scope
	oldAlias, err := amm.Get(scope, aliasKey)
	if err != nil {
		newError := ierrors.Wrap(
			ierrors.From(err).NotFound(),
			"alias '%s' not found on scope '%s'", aliasKey, scope,
		)
		return newError
	}
	parentApp, _ := amm.treeMemoryManager.Apps().Get(scope)

	l.Debug("validating Alias")
	// valid target channel
	err = validTargetChannel(parentApp, alias.Target)
	if err != nil {
		l.Debug("unable to update Alias - invalid targeted channel or boundary in parent dApp",
			zap.String("parent", parentApp.Meta.Name),
			zap.String("targeted boundary", alias.Target))
		return err
	}

	alias.Meta.UUID = oldAlias.Meta.UUID
	l.Debug("replacing old Alias with the new one",
		zap.String("parent", parentApp.Meta.Name))
	//update alias
	parentApp.Spec.Aliases[aliasKey] = alias

	return nil
}

// Delete receives a scope and a alias key. The scope
// defines the path to the dApp that cointains the Alias to be deleted. If the dApp
// has an alias that has the same key as the key passed as an argument, that alias
// is removed from the dApp Aliases only if it's not being used
func (amm *AliasMemoryManager) Delete(scope, aliasKey string) error {
	l := amm.logger.With(
		zap.String("operation", "delete"),
		zap.Any("alias", aliasKey),
		zap.String("scope", scope),
	)
	l.Debug("trying to delete an Alias")
	// get app from scope
	app, err := amm.treeMemoryManager.Apps().Get(scope)
	if err != nil {
		return err
	}

	l.Debug("checking if Alias to be deleted exists in given scope")
	// check if alias key exist in scope
	if _, ok := app.Spec.Aliases[aliasKey]; !ok {
		return ierrors.
			New("alias not found for the given key %v", aliasKey).
			BadRequest()
	}

	childName := strings.Split(aliasKey, ".")[0]
	target := strings.Split(aliasKey, ".")[1]

	l.Debug("checking if Alias can be deleted")
	// check if its being used by a child app
	if childApp, ok := app.Spec.Apps[childName]; ok {
		childBound := childApp.Spec.Boundary
		if childBound.Input.Contains(target) || childBound.Output.Contains(target) {
			l.Debug("unable to delete Alias that is being used by a dApp")
			return ierrors.New(
				"can't delete the alias since it's being used by a child app",
			).BadRequest()
		}
	}

	l.Debug("removing Alias from its parents 'Aliases' structure")

	// delete alias
	delete(app.Spec.Aliases, aliasKey)

	return nil
}

// AliasPermTreeGetter returns a getter that gets alias from the root structure of the app, without the current changes.
// The getter does not allow changes in the structure, just visualization.
type AliasPermTreeGetter struct {
	*PermTreeGetter
	logs *zap.Logger
}

// Get receives a scope and a alias key. The scope defines
// the path to an App. If this App has a pointer to a alias that has the
// same key as the key passed as an argument, the pointer to that alias is returned
// This method is used to get the structure as it is in the cluster, before any modifications.
func (amm *AliasPermTreeGetter) Get(scope, aliasKey string) (*meta.Alias, error) {
	l := amm.logs.With(
		zap.String("operation", "root-get"),
		zap.String("alias", aliasKey),
		zap.String("scope", scope),
	)
	l.Debug("trying to get an Alias")
	// get app from scope
	app, err := amm.Apps().Get(scope)
	if err != nil {
		l.Debug("unable to get Alias")
		return nil, err
	}

	// check if alias key exist in scope
	if _, ok := app.Spec.Aliases[aliasKey]; !ok {
		l.Debug("alias doesn't exist")
		return nil, ierrors.New(
			"alias not found for the given key %v", aliasKey,
		).BadRequest()
	}

	//return alias
	return app.Spec.Aliases[aliasKey], nil
}

func validTargetChannel(app *meta.App, targetChannel string) error {
	logger.Debug("validating if Alias targets a valid Channel")
	parentBound := app.Spec.Boundary
	if _, ok := app.Spec.Channels[targetChannel]; !ok && !parentBound.Input.Contains(targetChannel) && !parentBound.Output.Contains(targetChannel) {
		logger.Debug("alias targets an invalid Channel")
		return ierrors.New("channel doesn't exist in app").BadRequest()
	}

	return nil
}
