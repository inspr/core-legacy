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
ChannelTypes is a MemoryManager method that provides an access point for ChannelTypes
*/
func (tmm *MemoryManager) ChannelTypes() memory.ChannelTypeMemory {
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
	nameErr := meta.StructureNameIsValid(ct.Meta.Name)
	if nameErr != nil {
		return ierrors.NewError().InnerError(nameErr).Message(nameErr.Error()).Build()
	}

	_, err := ctm.GetChannelType(context, ct.Meta.Name)
	if err == nil {
		return ierrors.NewError().AlreadyExists().
			Message("target app already has a '" + ct.Meta.Name + "' ChannelType").Build()
	}

	parentApp, err := GetTreeMemory().Apps().GetApp(context)
	if err != nil {
		return err
	}

	if parentApp.Spec.ChannelTypes == nil {
		parentApp.Spec.ChannelTypes = map[string]*meta.ChannelType{}
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
			Message("target app doesn't exist").Build()

	}

	if parentApp.Spec.ChannelTypes != nil {
		if ct, ok := parentApp.Spec.ChannelTypes[ctName]; ok {
			return ct, nil
		}
	}

	err = ierrors.NewError().NotFound().Message("no ChannelType found for query.").Build()
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
			Message("target app doesn't contain a '" + context + "' ChannelType").Build()
	}

	parentApp, err := GetTreeMemory().Apps().GetApp(context)
	if err != nil {
		return ierrors.NewError().InternalServer().InnerError(err).
			Message("target app doesn't exist").Build()
	}

	delete(parentApp.Spec.ChannelTypes, ctName)

	curCt, err = ctm.GetChannelType(context, ctName)
	if curCt != nil {
		return ierrors.NewError().InternalServer().
			Message("couldn't delete '" + context + "' ChannelType from target app").Build()
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
			Message("target app doesn't contain a '" + context + "' ChannelType").Build()
	}

	parentApp, err := GetTreeMemory().Apps().GetApp(context)
	if err != nil {
		return ierrors.NewError().InternalServer().InnerError(err).
			Message("target app doesn't exist").Build()
	}

	parentApp.Spec.ChannelTypes[ct.Meta.Name] = ct
	return nil
}
