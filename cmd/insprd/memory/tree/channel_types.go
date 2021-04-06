package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"go.uber.org/zap"
)

// ChannelTypeMemoryManager implements the ChannelType interface
// and provides methos for operating on ChannelTypes
type ChannelTypeMemoryManager struct {
	*MemoryManager
}

// ChannelTypes is a MemoryManager method that provides an access point for ChannelTypes
func (tmm *MemoryManager) ChannelTypes() memory.ChannelTypeMemory {
	return &ChannelTypeMemoryManager{
		MemoryManager: tmm,
	}
}

// Create creates, if it doesn't already exist, a new Channel Type for a given app.
// ct: ChannetType to be created.
// context: Path to reference app (x.y.z...)
func (ctm *ChannelTypeMemoryManager) Create(context string, ct *meta.ChannelType) error {
	logger.Info("trying to create a Channel Type",
		zap.String("channelType", ct.Meta.Name),
		zap.String("context", context))

	logger.Debug("validating Channel Type structure")
	nameErr := utils.StructureNameIsValid(ct.Meta.Name)
	if nameErr != nil {
		logger.Error("invalid Channel Type name",
			zap.String("ctype", ct.Meta.Name))
		return ierrors.NewError().InnerError(nameErr).Message(nameErr.Error()).Build()
	}

	logger.Debug("checking if Channel Type already exists",
		zap.String("channel", ct.Meta.Name),
		zap.String("context", context))

	_, err := ctm.Get(context, ct.Meta.Name)
	if err == nil {
		logger.Error("channel Type already exists")
		return ierrors.NewError().AlreadyExists().
			Message("target app already has a '" + ct.Meta.Name + "' ChannelType").Build()
	}

	logger.Debug("getting Channel Type parent dApp")
	parentApp, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		newError := ierrors.NewError().InnerError(err).InvalidChannel().
			Message("couldn't create channel type " + ct.Meta.Name + "\n" + err.Error()).
			Build()
		return newError
	}

	logger.Debug("adding Channel Type to dApp",
		zap.String("channelType", ct.Meta.Name),
		zap.String("context", parentApp.Meta.Name))
	if parentApp.Spec.ChannelTypes == nil {
		parentApp.Spec.ChannelTypes = map[string]*meta.ChannelType{}
	}
	ct.Meta = utils.InjectUUID(ct.Meta)
	parentApp.Spec.ChannelTypes[ct.Meta.Name] = ct
	return nil
}

// Get returns, if it exists, a specific Channel Type from a given app.
// ctName: Name of desired Channel Type.
// context: Path to reference app (x.y.z...)
func (ctm *ChannelTypeMemoryManager) Get(context string, ctName string) (*meta.ChannelType, error) {
	logger.Info("trying to get a Channel Type",
		zap.String("channelType", ctName),
		zap.String("context", context))

	parentApp, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		return nil, ierrors.NewError().BadRequest().InnerError(err).
			Message("target dApp doesn't exist").Build()
	}

	if parentApp.Spec.ChannelTypes != nil {
		if ct, ok := parentApp.Spec.ChannelTypes[ctName]; ok {
			return ct, nil
		}
	}

	logger.Error("unable to get Channel Type in given context",
		zap.String("ctype", ctName),
		zap.String("context", context))

	return nil, ierrors.NewError().NotFound().
		Message("channelType not found for given query").
		Build()
}

// Delete deletes, if it exists, a Channel Type from a given app.
// ctName: Name of desired Channel Type.
// context: Path to reference app (x.y.z...)
func (ctm *ChannelTypeMemoryManager) Delete(context string, ctName string) error {
	logger.Info("trying to delete a Channel Type",
		zap.String("channelType", ctName),
		zap.String("context", context))

	curCt, err := ctm.Get(context, ctName)
	if curCt == nil || err != nil {
		return ierrors.NewError().BadRequest().
			Message("target app doesn't contain a '" + context + "' ChannelType").Build()
	}

	logger.Debug("checking if Channel Type can be deleted")
	if len(curCt.ConnectedChannels) > 0 {
		logger.Error("unable to delete Channel Type for it's being used",
			zap.Any("connected channels", curCt.ConnectedChannels))

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

	logger.Debug("removing Channel Type from its parents 'ChannelTypes' structure",
		zap.String("channelType", ctName),
		zap.String("dApp", parentApp.Meta.Name))

	delete(parentApp.Spec.ChannelTypes, ctName)

	return nil
}

// Update updates, if it exists, a Channel Type of a given app.
// ct: Updated ChannetType to be updated on app
// context: Path to reference app (x.y.z...)
func (ctm *ChannelTypeMemoryManager) Update(context string, ct *meta.ChannelType) error {
	logger.Info("trying to update a Channel Type",
		zap.String("channelType", ct.Meta.Name),
		zap.String("context", context))

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

	logger.Debug("replacing old Channel Type with the new one in dApps 'ChannelTypes",
		zap.String("channel", ct.Meta.Name),
		zap.String("dApp", parentApp.Meta.Name))

	parentApp.Spec.ChannelTypes[ct.Meta.Name] = ct
	return nil
}

// ChannelTypeRootGetter returns a getter that gets channel types from the root structure of the app, without the current changes.
// The getter does not allow changes in the structure, just visualization.
type ChannelTypeRootGetter struct {
}

// Get receives a query string (format = 'x.y.z') and iterates through the
// memory tree until it finds the Channel Type which name is equal to the last query element.
// If the specified Channel Type is found, it is returned. Otherwise, returns an error.
// This method is used to get the structure as it is in the cluster, before any modifications.
func (amm *ChannelTypeRootGetter) Get(context string, ctName string) (*meta.ChannelType, error) {
	logger.Info("trying to get a Channel Type (Root Getter)",
		zap.String("channelType", ctName),
		zap.String("context", context))

	parentApp, err := GetTreeMemory().Root().Apps().Get(context)
	if err != nil {
		return nil, ierrors.NewError().BadRequest().InnerError(err).
			Message("target dApp doesn't exist (Root Getter)").Build()
	}

	if parentApp.Spec.ChannelTypes != nil {
		if ch, ok := parentApp.Spec.ChannelTypes[ctName]; ok {
			return ch, nil
		}
	}

	logger.Error("unable to get Channel Type in given context (Root Getter)",
		zap.String("ctype", ctName),
		zap.String("context", context))

	return nil, ierrors.NewError().NotFound().
		Message("channelType not found for given query (Root Getter)").
		Build()
}
