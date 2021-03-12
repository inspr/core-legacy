package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
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

// CreateAlias DOC TODO
func (amm *AliasMemoryManager) CreateAlias(query string, targetBoundary string, alias *meta.Alias) error {
	// get app from query
	app, err := GetTreeMemory().Apps().Get(query)
	if err != nil {
		return err
	}

	// check if targetBoundary exists in app
	appBound := app.Spec.Boundary
	if !appBound.Input.Contains(targetBoundary) && !appBound.Output.Contains(targetBoundary) {
		return ierrors.NewError().BadRequest().Message("target boundary doesn't exist in " + app.Meta.Name).Build()
	}

	// get parentApp of app
	parentApp, err := getParentApp(query)
	if err != nil {
		return err
	}

	targetChannel := alias.Target

	// check if targetChannel exists in channels or boundaries of parentApp
	err = validTargetChannel(parentApp, targetChannel)
	if err != nil {
		return err
	}

	aliasKey := app.Meta.Name + "." + targetBoundary

	// check if alias is already there
	if _, ok := parentApp.Spec.Aliases[aliasKey]; ok {
		return ierrors.NewError().BadRequest().Message("alias already exists in parent app").Build()
	}

	// add new alias to Aliases list in parentApp
	parentApp.Spec.Aliases[aliasKey] = alias

	return nil

}

// UpdateAlias DOC TODO
func (amm *AliasMemoryManager) UpdateAlias(context string, aliasKey string, alias *meta.Alias) error {
	// get app from context
	app, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		return err
	}

	// check if alias key exist in context
	if _, ok := app.Spec.Aliases[aliasKey]; !ok {
		return ierrors.NewError().BadRequest().Message("alias not found for the given key " + aliasKey).Build()
	}

	// valid target channel
	err = validTargetChannel(app, alias.Target)
	if err != nil {
		return err
	}

	//update alias
	app.Spec.Aliases[aliasKey] = alias

	return nil
}

func validTargetChannel(parentApp *meta.App, targetChannel string) error {
	parentBound := parentApp.Spec.Boundary
	if _, ok := parentApp.Spec.Channels[targetChannel]; !ok && !parentBound.Input.Contains(targetChannel) && !parentBound.Output.Contains(targetChannel) {
		return ierrors.NewError().BadRequest().Message("channel doesn't exist in app").Build()
	}

	return nil
}
