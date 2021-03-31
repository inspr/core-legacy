package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
)

/*
ChannelTypeMemoryManager implements the ChannelType interface
and provides methos for operating on ChannelTypes
*/
type ChannelTypeMemoryManager struct {
	*MemoryManager
}

/*
ChannelTypes is a MemoryManager method that provides an access point for ChannelTypes
*/
func (tmm *MemoryManager) ChannelTypes() memory.ChannelTypeMemory {
	return &ChannelTypeMemoryManager{
		MemoryManager: tmm,
	}
}

/*
Create creates, if it doesn't already exist, a new ChannellType for a given app.
ct: ChannetType to be created.
context: Path to reference app (x.y.z...)
*/
func (ctm *ChannelTypeMemoryManager) Create(context string, ct *meta.ChannelType) error {
	nameErr := utils.StructureNameIsValid(ct.Meta.Name)
	if nameErr != nil {
		return ierrors.NewError().InnerError(nameErr).Message(nameErr.Error()).Build()
	}

	_, err := ctm.Get(context, ct.Meta.Name)
	if err == nil {
		return ierrors.NewError().AlreadyExists().
			Message("target app already has a '" + ct.Meta.Name + "' ChannelType").Build()
	}

	parentApp, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		newError := ierrors.NewError().InnerError(err).InvalidChannel().
			Message("couldn't create channel type " + ct.Meta.Name + "\n" + err.Error()).
			Build()
		return newError
	}

	if parentApp.Spec.ChannelTypes == nil {
		parentApp.Spec.ChannelTypes = map[string]*meta.ChannelType{}
	}

	parentApp.Spec.ChannelTypes[ct.Meta.Name] = ct
	return nil
}

/*
Get returns, if it exists, a specific ChannellType from a given app.
ctName: Name of desired Channel Type.
context: Path to reference app (x.y.z...)
*/
func (ctm *ChannelTypeMemoryManager) Get(context string, ctName string) (*meta.ChannelType, error) {
	parentApp, err := GetTreeMemory().Apps().Get(context)
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
Delete deletes, if it exists, a ChannellType from a given app.
ctName: Name of desired Channel Type.
context: Path to reference app (x.y.z...)
*/
func (ctm *ChannelTypeMemoryManager) Delete(context string, ctName string) error {
	curCt, err := ctm.Get(context, ctName)
	if curCt == nil || err != nil {
		return ierrors.NewError().BadRequest().
			Message("target app doesn't contain a '" + context + "' ChannelType").Build()
	}

	if len(curCt.ConnectedChannels) > 0 {
		return ierrors.NewError().
			BadRequest().
			Message("channelType cannot be deleted as it is being used by other channels").
			Build()
	}

	parentApp, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		return ierrors.NewError().InternalServer().InnerError(err).
			Message("target app doesn't exist").Build()
	}

	delete(parentApp.Spec.ChannelTypes, ctName)

	_, err = ctm.Get(context, ctName)
	if err == nil {
		return ierrors.NewError().InternalServer().
			Message("couldn't delete '" + context + "' ChannelType from target app").Build()
	}
	return nil
}

/*
Update updates, if it exists, a ChannellType of a given app.
ct: Updated ChannetType to be updated on app
context: Path to reference app (x.y.z...)
*/
func (ctm *ChannelTypeMemoryManager) Update(context string, ct *meta.ChannelType) error {

	oldChType, err := ctm.Get(context, ct.Meta.Name)
	if err != nil {
		return ierrors.NewError().BadRequest().
			Message("target app doesn't contain a '" + context + "' ChannelType").Build()
	}

	ct.ConnectedChannels = oldChType.ConnectedChannels

	parentApp, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		return ierrors.NewError().InternalServer().InnerError(err).
			Message("target app doesn't exist").Build()
	}

	parentApp.Spec.ChannelTypes[ct.Meta.Name] = ct
	return nil
}

// ChannelTypeRootGetter returns a getter that gets channel types from the root structure of the app, without the current changes.
// The getter does not allow changes in the structure, just visualization.
type ChannelTypeRootGetter struct {
}

// Get receives a query string (format = 'x.y.z') and iterates through the
// memory tree until it finds the dChannelType which name is equal to the last query element.
// The tree ChannelType is returned if the query string is an empty string.
// If the specified dChannelType is found, it is returned. Otherwise, returns an error.
func (amm *ChannelTypeRootGetter) Get(context string, name string) (*meta.ChannelType, error) {
	parentApp, err := GetTreeMemory().Root().Apps().Get(context)
	if err != nil {
		newError := ierrors.
			NewError().
			InnerError(err).
			NotFound().
			Message(
				"ChannelType not found, the context '%v' is invalid",
				context,
			).
			Build()
		return nil, newError
	}

	if parentApp.Spec.ChannelTypes != nil {
		if ch, ok := parentApp.Spec.ChannelTypes[name]; ok {
			return ch, nil
		}
	}

	newError := ierrors.
		NewError().
		NotFound().
		Message("channelType '%v' not found", name).
		Build()
	return nil, newError
}
