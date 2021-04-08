package tree

import (
	"strings"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"go.uber.org/zap"
)

// AliasMemoryManager implements the Alias interface
// and provides methos for operating on Aliass
type AliasMemoryManager struct {
	*MemoryManager
}

// Alias is a MemoryManager method that provides an access point for Alias
func (tmm *MemoryManager) Alias() memory.AliasMemory {
	return &AliasMemoryManager{
		MemoryManager: tmm,
	}
}

// Get receives a context and an alias key. The context defines
// the path to a dApp. If this dApp has a pointer to a alias that has the
// same key as the key passed as an argument, the pointer to that alias is returned
func (amm *AliasMemoryManager) Get(context, aliasKey string) (*meta.Alias, error) {
	logger.Info("trying to get an Alias",
		zap.String("alias", aliasKey),
		zap.String("context", context))

	app, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		logger.Error("unable to get Alias")
		return nil, err
	}

	// check if alias key exist in context
	if _, ok := app.Spec.Aliases[aliasKey]; !ok {
		return nil, ierrors.
			NewError().
			BadRequest().
			Message("alias not found for the given key %v", aliasKey).
			Build()
	}

	//return alias
	return app.Spec.Aliases[aliasKey], nil
}

// Create receives a context that defines a path to the dApp.
// The new alias will be created inside this dApp's parent.
func (amm *AliasMemoryManager) Create(context, targetBoundary string, alias *meta.Alias) error {
	logger.Info("trying to create an Alias",
		zap.Any("alias", alias),
		zap.String("context", context))
	// get app from context
	app, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		return err
	}

	logger.Debug("checking if dApp boundary is valid for given Alias")
	// check if targetBoundary exists in app
	appBound := app.Spec.Boundary
	if !appBound.Input.Contains(targetBoundary) && !appBound.Output.Contains(targetBoundary) {
		logger.Error("invalid dApp boundary for Alias",
			zap.Any("alias", alias),
			zap.String("targeted boundary", targetBoundary))
		return ierrors.NewError().BadRequest().Message("target boundary doesn't exist in %v", app.Meta.Name).Build()
	}

	// get parentApp of app
	parentApp, _ := getParentApp(context)

	targetChannel := alias.Target

	logger.Debug("checking if Alias targeted Channel exists")
	// check if targetChannel exists in channels or boundaries of parentApp
	err = validTargetChannel(parentApp, targetChannel)
	if err != nil {
		logger.Error("unable to create Alias - invalid targeted channel or boundary in parent dApp",
			zap.Any("alias", alias),
			zap.String("parent dApp", parentApp.Meta.Name),
			zap.String("targeted boundary", targetChannel))
		return err
	}

	aliasKey := app.Meta.Name + "." + targetBoundary

	logger.Debug("checking if Alias already exists")
	// check if alias is already there
	if _, ok := parentApp.Spec.Aliases[aliasKey]; ok {
		logger.Error("alias already exists")
		return ierrors.NewError().BadRequest().Message("alias already exists in parent dApp").Build()
	}

	alias.Meta = utils.InjectUUID(alias.Meta)

	logger.Debug("adding Alias to dApp",
		zap.Any("alias", alias),
		zap.String("dApp", parentApp.Meta.Name))
	// add new alias to Aliases list in parentApp
	parentApp.Spec.Aliases[aliasKey] = alias

	return nil

}

// Update receives a context a alias key and a alias. The context
// defines the path to the dApp that contains the Alias. If the dApp has
// a alias that has the given alias key passed as an argument,
// that alias will be replaced by the new alias
func (amm *AliasMemoryManager) Update(context, aliasKey string, alias *meta.Alias) error {
	logger.Info("trying to update an Alias",
		zap.Any("alias", alias),
		zap.String("context", context))

	logger.Debug("checking if Alias to be updated exists in given context")
	// check if alias key exist in context
	oldAlias, err := amm.Get(context, aliasKey)
	if err != nil {
		newError := ierrors.NewError().
			InnerError(err).
			NotFound().
			Message("alias '%s' not found on context '%s'", aliasKey, context).
			Build()
		return newError
	}
	parentApp, _ := GetTreeMemory().Apps().Get(context)

	logger.Debug("validating Alias")
	// valid target channel
	err = validTargetChannel(parentApp, alias.Target)
	if err != nil {
		logger.Error("unable to update Alias - invalid targeted channel or boundary in parent dApp",
			zap.Any("alias", alias),
			zap.String("parent dApp", parentApp.Meta.Name),
			zap.String("targeted boundary", alias.Target))
		return err
	}

	alias.Meta.UUID = oldAlias.Meta.UUID
	logger.Debug("replacing old Alias with the new one in dApps 'Aliases'",
		zap.Any("alias", alias),
		zap.String("dApp", parentApp.Meta.Name))
	//update alias
	parentApp.Spec.Aliases[aliasKey] = alias

	return nil
}

// Delete receives a context and a alias key. The context
// defines the path to the dApp that cointains the Alias to be deleted. If the dApp
// has an alias that has the same key as the key passed as an argument, that alias
// is removed from the dApp Aliases only if it's not being used
func (amm *AliasMemoryManager) Delete(context, aliasKey string) error {
	logger.Info("trying to delete an Alias",
		zap.Any("alias", aliasKey),
		zap.String("context", context))
	// get app from context
	app, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		return err
	}

	logger.Debug("checking if Alias to be deleted exists in given context")
	// check if alias key exist in context
	if _, ok := app.Spec.Aliases[aliasKey]; !ok {
		return ierrors.
			NewError().
			BadRequest().
			Message("alias not found for the given key %v", aliasKey).
			Build()
	}

	childName := strings.Split(aliasKey, ".")[0]
	target := strings.Split(aliasKey, ".")[1]

	logger.Debug("checking if Alias can be deleted")
	// check if its being used by a child app
	if childApp, ok := app.Spec.Apps[childName]; ok {
		childBound := childApp.Spec.Boundary
		if childBound.Input.Contains(target) || childBound.Output.Contains(target) {
			logger.Error("unable to delete Alias that is being used by a dApp")
			return ierrors.NewError().BadRequest().Message("can't delete the alias since it's being used by a child app").Build()
		}
	}

	logger.Debug("removing Alias from its parents 'Aliases' structure",
		zap.String("alias", aliasKey),
		zap.String("dApp", app.Meta.Name))

	// delete alias
	delete(app.Spec.Aliases, aliasKey)

	return nil
}

// AliasRootGetter returns a getter that gets alias from the root structure of the app, without the current changes.
// The getter does not allow changes in the structure, just visualization.
type AliasRootGetter struct {
}

// Get receives a context and a alias key. The context defines
// the path to an App. If this App has a pointer to a alias that has the
// same key as the key passed as an argument, the pointer to that alias is returned
// This method is used to get the structure as it is in the cluster, before any modifications.
func (amm *AliasRootGetter) Get(context, aliasKey string) (*meta.Alias, error) {
	logger.Info("trying to get an Alias (Root Getter)",
		zap.String("alias", aliasKey),
		zap.String("context", context))
	// get app from context
	app, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		logger.Error("unable to get Alias (Root Getter)")
		return nil, err
	}

	// check if alias key exist in context
	if _, ok := app.Spec.Aliases[aliasKey]; !ok {
		logger.Error("alias doesn't exists (Root Getter)")
		return nil, ierrors.NewError().BadRequest().Message("alias not found for the given key %v", aliasKey).Build()
	}

	//return alias
	return app.Spec.Aliases[aliasKey], nil
}

func validTargetChannel(app *meta.App, targetChannel string) error {
	logger.Debug("validating if Alias targets a valid Channel")
	parentBound := app.Spec.Boundary
	if _, ok := app.Spec.Channels[targetChannel]; !ok && !parentBound.Input.Contains(targetChannel) && !parentBound.Output.Contains(targetChannel) {
		logger.Error("alias targets an invalid Channel")
		return ierrors.NewError().BadRequest().Message("channel doesn't exist in app").Build()
	}

	return nil
}
