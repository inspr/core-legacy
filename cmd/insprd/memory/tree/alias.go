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
func (amm *AliasMemoryManager) CreateAlias(context string, targetBoundary string, targetChannel string) error {
	// get app from context
	app, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		return err
	}

	// check if targetBoundary exists in app
	appBound := app.Spec.Boundary
	if !appBound.Input.Contains(targetBoundary) && !appBound.Output.Contains(targetBoundary) {
		return ierrors.NewError().BadRequest().Message("target boundary doesn't exist in " + app.Meta.Name).Build()
	}

	// get parentApp of app
	parentApp, err := getParentApp(context)
	if err != nil {
		return err
	}

	// check if targetChannel exists in channels or boundaries of parentApp
	parentBound := parentApp.Spec.Boundary
	if _, ok := parentApp.Spec.Channels[targetChannel]; !ok && !parentBound.Input.Contains(targetChannel) && !parentBound.Output.Contains(targetChannel) {
		return ierrors.NewError().BadRequest().Message("channel doesn't exist in parent app").Build()
	}

	aliasKey := app.Meta.Name + "." + targetBoundary

	// check if alias is already there
	if _, ok := parentApp.Spec.Aliases[aliasKey]; ok {
		return ierrors.NewError().BadRequest().Message("alias already exists in parent app").Build()
	}

	// create alias structure
	newAlias := &meta.Alias{
		Target: targetChannel,
	}

	// add new alias to Aliases list in parentApp
	parentApp.Spec.Aliases[aliasKey] = newAlias

	return nil

}
