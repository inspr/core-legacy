package tree

import (
	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	metautils "inspr.dev/inspr/pkg/meta/utils"
	"inspr.dev/inspr/pkg/utils"
)

// ChannelMemoryManager implements the channel interface and
// provides methods for operating on Channels
type ChannelMemoryManager struct {
	*treeMemoryManager
}

// Channels return a pointer to ChannelMemoryManager
func (tmm *treeMemoryManager) Channels() ChannelMemory {
	return &ChannelMemoryManager{
		treeMemoryManager: tmm,
	}
}

// Get receives a scope and a channel name. The scope defines
// the path to an App. If this App has a pointer to a channel that has the
// same name as the name passed as an argument, the pointer to that channel is returned
func (chh *ChannelMemoryManager) Get(scope, name string) (*meta.Channel, error) {
	logger.Info("trying to get a Channel",
		zap.String("channel", name),
		zap.String("scope", scope))

	parentApp, err := chh.Apps().Get(scope)
	if err != nil {
		newError := ierrors.
			NewError().
			InnerError(err).
			NotFound().
			Message("channel not found, the scope '%v' is invalid", scope).
			Build()
		return nil, newError
	}

	if parentApp.Spec.Channels != nil {
		if ch, ok := parentApp.Spec.Channels[name]; ok {
			return ch, nil
		}
	}

	logger.Debug("unable to get Channel in given scope",
		zap.String("channel", name),
		zap.String("scope", scope))

	newError := ierrors.
		NewError().
		NotFound().
		Message("channel not found").
		Build()
	return nil, newError
}

// Create receives a scope that defines a path to the App
// in which to add a pointer to the channel passed as an argument
func (chh *ChannelMemoryManager) Create(scope string, ch *meta.Channel) error {
	logger.Info("trying to create a Channel",
		zap.String("channel", ch.Meta.Name),
		zap.String("scope", scope))

	nameErr := metautils.StructureNameIsValid(ch.Meta.Name)
	if nameErr != nil {
		logger.Error("invalid Channel name",
			zap.String("channel", ch.Meta.Name))
		return ierrors.NewError().InnerError(nameErr).Message(nameErr.Error()).Build()
	}

	logger.Debug("checking if Channel already exists",
		zap.String("channel", ch.Meta.Name),
		zap.String("scope", scope))

	chAlreadyExist, _ := chh.Get(scope, ch.Meta.Name)
	if chAlreadyExist != nil {
		logger.Error("channel already exists")
		return ierrors.NewError().AlreadyExists().
			Message("channel with name %v already exists in the scope %v", ch.Meta.Name, scope).
			Build()
	}

	logger.Debug("getting Channel parent dApp")
	parentApp, err := chh.Apps().Get(scope)
	if err != nil {
		newError := ierrors.NewError().InnerError(err).InvalidChannel().
			Message("couldn't create channel %v : %v", ch.Meta.Name, err.Error()).
			Build()
		return newError
	}

	logger.Debug("checking if Channel's type is valid")
	if _, ok := parentApp.Spec.Types[ch.Spec.Type]; !ok {
		logger.Error("channel's type is invalid")
		return ierrors.NewError().InvalidChannel().Message("references a Type that doesn't exist").Build()
	}

	connectedChannels := parentApp.Spec.Types[ch.Spec.Type].ConnectedChannels
	if !utils.Includes(connectedChannels, ch.Meta.Name) {
		connectedChannels = append(connectedChannels, ch.Meta.Name)
		parentApp.Spec.Types[ch.Spec.Type].ConnectedChannels = connectedChannels
	}

	logger.Debug("channel broker priority list", zap.Any("list", ch.Spec.BrokerPriorityList))

	broker, err := SelectBrokerFromPriorityList(ch.Spec.BrokerPriorityList)
	if err != nil {
		return err
	}

	logger.Debug("channel selected broker", zap.Any("broker", broker))
	ch.Spec.SelectedBroker = broker

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

// Delete receives a scope and a channel name. The scope
// defines the path to the App that will have the Delete. If the App
// has a pointer to a channel that has the same name as the name passed
// as an argument, that pointer is removed from the list of App channels
func (chh *ChannelMemoryManager) Delete(scope, name string) error {
	logger.Info("trying to delete a Channel",
		zap.String("channel", name),
		zap.String("scope", scope))

	channel, err := chh.Get(scope, name)

	if err != nil {
		newError := ierrors.
			NewError().
			InnerError(err).
			NotFound().
			Message("channel %s not found", name).
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

	parentApp, _ := chh.Apps().Get(scope)

	insprType := parentApp.Spec.Types[channel.Spec.Type]

	logger.Debug("removing Channel from Type connected channels list",
		zap.String("channel", name),
		zap.String("type", insprType.Meta.Name))

	insprType.ConnectedChannels = utils.Remove(
		insprType.ConnectedChannels,
		channel.Meta.Name,
	)

	logger.Debug("removing Channel from its parents 'Channels' structure",
		zap.String("channel", name),
		zap.String("dApp", parentApp.Meta.Name))

	delete(parentApp.Spec.Channels, name)

	return nil
}

// Update receives a scope and a channel pointer. The scope
// defines the path to the App that will have the Update. If the App has
// a channel pointer that has the same name as that passed as an argument,
// this pointer will be replaced by the new one
func (chh *ChannelMemoryManager) Update(scope string, ch *meta.Channel) error {
	logger.Info("trying to update a Channel",
		zap.String("channel", ch.Meta.Name),
		zap.String("scope", scope))

	oldCh, err := chh.Get(scope, ch.Meta.Name)
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

	parentApp, _ := chh.Apps().Get(scope)

	logger.Debug("validating new Channel structure")

	if _, ok := parentApp.Spec.Types[ch.Spec.Type]; !ok {
		logger.Error("unable to create Channel for it references an invalid Type",
			zap.String("invalid Type", ch.Spec.Type))

		return ierrors.NewError().InvalidChannel().Message("references a Type that doesn't exist").Build()
	}

	logger.Debug("replacing old Channel with the new one in dApps 'Channels'",
		zap.String("channel", ch.Meta.Name),
		zap.String("dApp", parentApp.Meta.Name))

	parentApp.Spec.Channels[ch.Meta.Name] = ch

	return nil
}

// ChannelPermTreeGetter returns a getter that gets channels from the root structure of the app, without the current changes.
// The getter does not allow changes in the structure, just visualization.
type ChannelPermTreeGetter struct {
}

// Get receives a query string (format = 'x.y.z') and iterates through the
// memory tree until it finds the Channel which name is equal to the last query element.
// If the specified Channel is found, it is returned. Otherwise, returns an error.
// This method is used to get the structure as it is in the cluster, before any modifications.
func (cmm *ChannelPermTreeGetter) Get(scope, name string) (*meta.Channel, error) {
	logger.Info("trying to get a Channel (Root Getter)",
		zap.String("channel", name),
		zap.String("scope", scope))

	parentApp, err := GetTreeMemory().Tree().Apps().Get(scope)
	if err != nil {
		newError := ierrors.
			NewError().
			InnerError(err).
			NotFound().
			Message("the '%v' scope is invalid", scope).Build()
		return nil, newError
	}

	if parentApp.Spec.Channels != nil {
		if ch, ok := parentApp.Spec.Channels[name]; ok {
			return ch, nil
		}
	}

	logger.Error("unable to get Channel in given scope (Root Getter)",
		zap.String("type", name),
		zap.String("scope", scope))

	newError := ierrors.
		NewError().
		NotFound().
		Message("channel not found (Root Getter)").
		Build()
	return nil, newError
}
