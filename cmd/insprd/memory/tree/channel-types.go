package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

/*
ChannelTypeMemoryManager implements the ChannelType interface
and provides methos for operating on ChannelTypes
*/
type ChannelTypeMemoryManager struct {
	root *meta.App
}

/*
ChannelTypes is a TreeMemoryManager method that provides an access point for ChannelTypes
*/
func (tmm *TreeMemoryManager) ChannelTypes() memory.ChannelTypeMemory {
	return &ChannelTypeMemoryManager{
		root: tmm.root,
	}
}

/*
CreateChannelType creates, if it doesn't already exist, a new ChannellType for a given app.
ct: ChannetType to be created.
context: Path to reference app (x.y.z...)
*/
func (ctm *ChannelTypeMemoryManager) CreateChannelType(ct *meta.ChannelType, context string) error {

	_, err := ctm.GetChannelType(context, ct.Meta.Name)
	if err == nil {
		return ierrors.NewError().AlreadyExists().
			Message("Target app already has a '" + ct.Meta.Name + "' ChannelType").Build()
	}

	parentApp, err := GetTreeMemory().Apps().GetApp(context)
	if err != nil {
		return err
	}

	parentApp.Spec.ChannelTypes[ct.Meta.Name] = ct
	return nil
}

/*
GetChannelType returns, if it exists, a specific ChannellType from a given app.
ctName: Name of desired Channel Type.
context: Path to reference app (x.y.z...)
*/
func (ctm *ChannelTypeMemoryManager) GetChannelType(context string, ctName string) (*meta.ChannelType, error) {
	parentApp, err := GetTreeMemory().Apps().GetApp(context)
	if err != nil {
		return nil, ierrors.NewError().BadRequest().InnerError(err).
			Message("Target app doesn't exist").Build()

	}

	if ct, ok := parentApp.Spec.ChannelTypes[ctName]; ok {
		return ct, nil
	}

	err = ierrors.NewError().NotFound().Message("No ChannelType found for query.").Build()
	return nil, err
}

/*
DeleteChannelType deletes, if it exists, a ChannellType from a given app.
ctName: Name of desired Channel Type.
context: Path to reference app (x.y.z...)
*/
func (ctm *ChannelTypeMemoryManager) DeleteChannelType(context string, ctName string) error {
	curCt, err := ctm.GetChannelType(context, ctName)
	if curCt == nil || err != nil {
		return ierrors.NewError().BadRequest().
			Message("Target app doesn't contain a '" + context + "' ChannelType").Build()
	}

	parentApp, err := GetTreeMemory().Apps().GetApp(context)
	if err != nil {
		return ierrors.NewError().InternalServer().InnerError(err).
			Message("Target app doesn't exist").Build()
	}

	delete(parentApp.Spec.ChannelTypes, ctName)

	curCt, err = ctm.GetChannelType(context, ctName)
	if curCt != nil {
		return ierrors.NewError().InternalServer().
			Message("Couldn't delete '" + context + "' ChannelType from target app").Build()
	}
	return nil
}

/*
UpdateChannelType updates, if it exists, a ChannellType of a given app.
ct: Updated ChannetType to be updated on app
context: Path to reference app (x.y.z...)
*/
func (ctm *ChannelTypeMemoryManager) UpdateChannelType(ct *meta.ChannelType, context string) error {

	_, err := ctm.GetChannelType(context, ct.Meta.Name)
	if err != nil {
		return ierrors.NewError().BadRequest().
			Message("Target app doesn't contain a '" + context + "' ChannelType").Build()
	}

	parentApp, err := GetTreeMemory().Apps().GetApp(context)
	if err != nil {
		return ierrors.NewError().InternalServer().InnerError(err).
			Message("Target app doesn't exist").Build()
	}

	parentApp.Spec.ChannelTypes[ct.Meta.Name] = ct
	return nil
}
