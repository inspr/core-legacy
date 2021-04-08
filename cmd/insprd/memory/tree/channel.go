package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	metautils "gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
	"go.uber.org/zap"
)

// ChannelMemoryManager implements the channel interface and
// provides methods for operating on Channels
type ChannelMemoryManager struct {
	*MemoryManager
}

// Channels return a pointer to ChannelMemoryManager
func (tmm *MemoryManager) Channels() memory.ChannelMemory {
	return &ChannelMemoryManager{
		MemoryManager: tmm,
	}
}

// Get receives a context and a channel name. The context defines
// the path to an App. If this App has a pointer to a channel that has the
// same name as the name passed as an argument, the pointer to that channel is returned
func (chh *ChannelMemoryManager) Get(context string, chName string) (*meta.Channel, error) {
	logger.Info("trying to get a Channel",
		zap.String("channel", chName),
		zap.String("context", context))

	parentApp, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		newError := ierrors.
			NewError().
			InnerError(err).
			NotFound().
			Message(
				"channel not found, the context '%v' is invalid",
				context,
			).
			Build()
		return nil, newError
	}

	if parentApp.Spec.Channels != nil {
		if ch, ok := parentApp.Spec.Channels[chName]; ok {
			return ch, nil
		}
	}

	logger.Error("unable to get Channel in given context",
		zap.String("ctype", chName),
		zap.String("context", context))

	newError := ierrors.
		NewError().
		NotFound().
		Message("channel not found").
		Build()
	return nil, newError
}

// Create receives a context that defines a path to the App
// in which to add a pointer to the channel passed as an argument
func (chh *ChannelMemoryManager) Create(context string, ch *meta.Channel) error {
	logger.Info("trying to create a Channel",
		zap.String("channel", ch.Meta.Name),
		zap.String("context", context))

	nameErr := metautils.StructureNameIsValid(ch.Meta.Name)
	if nameErr != nil {
		logger.Error("invalid Channel name",
			zap.String("ctype", ch.Meta.Name))
		return ierrors.NewError().InnerError(nameErr).Message(nameErr.Error()).Build()
	}

	logger.Debug("checking if Channel already exists",
		zap.String("channel", ch.Meta.Name),
		zap.String("context", context))

	chAlreadyExist, _ := chh.Get(context, ch.Meta.Name)
	if chAlreadyExist != nil {
		logger.Error("channel already exists")
		return ierrors.NewError().AlreadyExists().Message("channel with name " + ch.Meta.Name + " already exists in the context " + context).Build()
	}

	logger.Debug("getting Channel parent dApp")
	parentApp, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		newError := ierrors.NewError().InnerError(err).InvalidChannel().
			Message("couldn't create channel " + ch.Meta.Name + "\n" + err.Error()).
			Build()
		return newError
	}

	logger.Debug("checking if Channel's type is valid")
	if _, ok := parentApp.Spec.ChannelTypes[ch.Spec.Type]; !ok {
		logger.Error("channel's type is invalid")
		return ierrors.NewError().InvalidChannel().Message("references a Channel Type that doesn't exist").Build()
	}

	connectedChannels := parentApp.Spec.ChannelTypes[ch.Spec.Type].ConnectedChannels
	if !utils.Includes(connectedChannels, ch.Meta.Name) {
		connectedChannels = append(connectedChannels, ch.Meta.Name)
		parentApp.Spec.ChannelTypes[ch.Spec.Type].ConnectedChannels = connectedChannels
	}

	logger.Debug("adding Channel to dApp",
		zap.String("channel", ch.Meta.Name),
		zap.String("dApp", parentApp.Meta.Name))
	if parentApp.Spec.Channels == nil {
		parentApp.Spec.Channels = map[string]*meta.Channel{}
	}

	ch.Meta = metautils.InjectUUID(ch.Meta)

	parentApp.Spec.Channels[ch.Meta.Name] = ch

	return nil
}

// Delete receives a context and a channel name. The context
// defines the path to the App that will have the Delete. If the App
// has a pointer to a channel that has the same name as the name passed
// as an argument, that pointer is removed from the list of App channels
func (chh *ChannelMemoryManager) Delete(context string, chName string) error {
	logger.Info("trying to delete a Channel",
		zap.String("channel", chName),
		zap.String("context", context))

	channel, err := chh.Get(context, chName)

	if err != nil {
		newError := ierrors.
			NewError().
			InnerError(err).
			NotFound().
			Message("channel %s not found", chName).
			Build()
		return newError
	}

	logger.Debug("checking if Channel can be deleted")
	if len(channel.ConnectedApps) > 0 || len(channel.ConnectedAliases) > 0 {
		logger.Error("unable to delete Channel for it's being used",
			zap.Any("connected dApps", channel.ConnectedApps),
			zap.Any("connected Aliases", channel.ConnectedAliases))

		return ierrors.NewError().
			BadRequest().
			Message("channel cannot be deleted as it is being used by other apps").
			Build()
	}

	parentApp, _ := GetTreeMemory().Apps().Get(context)

	channelType := parentApp.Spec.ChannelTypes[channel.Spec.Type]

	logger.Debug("removing Channel from ChannelType connected channels list",
		zap.String("channel", chName),
		zap.String("channelType", channelType.Meta.Name))

	channelType.ConnectedChannels = utils.Remove(channelType.ConnectedChannels, channel.Meta.Name)

	logger.Debug("removing Channel from its parents 'Channels' structure",
		zap.String("channel", chName),
		zap.String("dApp", parentApp.Meta.Name))

	delete(parentApp.Spec.Channels, chName)

	return nil
}

// Update receives a context and a channel pointer. The context
// defines the path to the App that will have the Update. If the App has
// a channel pointer that has the same name as that passed as an argument,
// this pointer will be replaced by the new one
func (chh *ChannelMemoryManager) Update(context string, ch *meta.Channel) error {
	logger.Info("trying to update a Channel",
		zap.String("channel", ch.Meta.Name),
		zap.String("context", context))

	oldCh, err := chh.Get(context, ch.Meta.Name)
	if err != nil {
		newError := ierrors.
			NewError().
			InnerError(err).
			NotFound().
			Message("channel %s not found", ch.Meta.Name).
			Build()
		return newError
	}

	ch.ConnectedApps = oldCh.ConnectedApps
	ch.ConnectedAliases = oldCh.ConnectedAliases
	ch.Meta.UUID = oldCh.Meta.UUID

	parentApp, _ := GetTreeMemory().Apps().Get(context)

	logger.Debug("validating new Channel structure")

	if _, ok := parentApp.Spec.ChannelTypes[ch.Spec.Type]; !ok {
		logger.Error("unable to create Channel for it references an invalid Channel Type",
			zap.String("invalid Channel Type", ch.Spec.Type))

		return ierrors.NewError().InvalidChannel().Message("references a Channel Type that doesn't exist").Build()
	}

	logger.Debug("replacing old Channel with the new one in dApps 'Channels'",
		zap.String("channel", ch.Meta.Name),
		zap.String("dApp", parentApp.Meta.Name))

	parentApp.Spec.Channels[ch.Meta.Name] = ch

	return nil
}

// ChannelRootGetter returns a getter that gets channels from the root structure of the app, without the current changes.
// The getter does not allow changes in the structure, just visualization.
type ChannelRootGetter struct {
}

// Get receives a query string (format = 'x.y.z') and iterates through the
// memory tree until it finds the Channel which name is equal to the last query element.
// If the specified Channel is found, it is returned. Otherwise, returns an error.
// This method is used to get the structure as it is in the cluster, before any modifications.
func (amm *ChannelRootGetter) Get(context string, chName string) (*meta.Channel, error) {
	logger.Info("trying to get a Channel (Root Getter)",
		zap.String("channel", chName),
		zap.String("context", context))

	parentApp, err := GetTreeMemory().Root().Apps().Get(context)
	if err != nil {
		newError := ierrors.
			NewError().
			InnerError(err).
			NotFound().
			Message("the '%v' context is invalid", context).Build()
		return nil, newError
	}

	if parentApp.Spec.Channels != nil {
		if ch, ok := parentApp.Spec.Channels[chName]; ok {
			return ch, nil
		}
	}

	logger.Error("unable to get Channel in given context (Root Getter)",
		zap.String("ctype", chName),
		zap.String("context", context))

	newError := ierrors.
		NewError().
		NotFound().
		Message("channel not found (Root Getter)").
		Build()
	return nil, newError
}
