package tree

import (
	"go.uber.org/zap"
	apimodels "inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	metautils "inspr.dev/inspr/pkg/meta/utils"
	"inspr.dev/inspr/pkg/utils"
)

// ChannelMemoryManager implements the channel interface and
// provides methods for operating on Channels
type ChannelMemoryManager struct {
	*treeMemoryManager
	logger *zap.Logger
}

// Channels return a pointer to ChannelMemoryManager
func (tmm *treeMemoryManager) Channels() ChannelMemory {
	return &ChannelMemoryManager{
		treeMemoryManager: tmm,
		logger:            logger.With(zap.String("subSection", "channels")),
	}
}

// Get receives a scope and a channel name. The scope defines
// the path to an App. If this App has a pointer to a channel that has the
// same name as the name passed as an argument, the pointer to that channel is returned
func (chh *ChannelMemoryManager) Get(scope, name string) (*meta.Channel, error) {
	l := chh.logger.With(
		zap.String("operation", "get"),
		zap.String("channel", name),
		zap.String("scope", scope),
	)
	l.Debug("recevid channel retrieval request")

	parentApp, err := chh.Apps().Get(scope)
	if err != nil {
		logger.Debug("unable to find channel")
		newError := ierrors.Wrap(
			err,
			"channel not found, the scope '%v' is invalid", scope,
		)
		return nil, newError
	}

	if parentApp.Spec.Channels != nil {
		if ch, ok := parentApp.Spec.Channels[name]; ok {
			l.Debug("channel found")
			return ch, nil
		}
	}

	l.Debug("unable to get Channel in given scope")

	newError := ierrors.New("channel not found").NotFound()
	return nil, newError
}

// Create receives a scope that defines a path to the App
// in which to add a pointer to the channel passed as an argument
func (chh *ChannelMemoryManager) Create(scope string, ch *meta.Channel, brokers *apimodels.BrokersDI) error {
	l := chh.logger.With(
		zap.String("operation", "create"),
		zap.String("channel", ch.Meta.Name),
		zap.String("scope", scope),
	)
	l.Debug("trying to create a Channel")

	nameErr := metautils.StructureNameIsValid(ch.Meta.Name)
	if nameErr != nil {
		l.Debug("invalid Channel name",
			zap.String("channel", ch.Meta.Name))
		return ierrors.Wrap(nameErr, "failed to create Channel")
	}

	l.Debug("checking if Channel already exists")

	chAlreadyExist, _ := chh.Get(scope, ch.Meta.Name)
	if chAlreadyExist != nil {
		l.Debug("channel already exists")
		return ierrors.New(
			"channel with name %v already exists in the scope %v",
			ch.Meta.Name, scope,
		).AlreadyExists()
	}

	l.Debug("getting Channel parent dApp")
	parentApp, err := chh.Apps().Get(scope)
	if err != nil {
		newError := ierrors.Wrap(
			err,
			"couldn't create channel %v", ch.Meta.Name,
		)
		return newError
	}

	l.Debug("checking if Channel's type is valid")
	if _, ok := parentApp.Spec.Types[ch.Spec.Type]; !ok {
		l.Debug("channel's type is invalid")
		return ierrors.New(
			"references a Type that doesn't exist",
		).InvalidChannel()
	}

	connectedChannels := parentApp.Spec.Types[ch.Spec.Type].ConnectedChannels
	if !utils.Includes(connectedChannels, ch.Meta.Name) {
		connectedChannels = append(connectedChannels, ch.Meta.Name)
		parentApp.Spec.Types[ch.Spec.Type].ConnectedChannels = connectedChannels
	}

	l.Debug("channel broker priority list", zap.Any("list", ch.Spec.BrokerPriorityList))

	broker, err := SelectBrokerFromPriorityList(ch.Spec.BrokerPriorityList, brokers)
	if err != nil {
		return err
	}

	l.Debug("channel selected broker", zap.Any("broker", broker))
	ch.Spec.SelectedBroker = broker

	l.Debug("adding Channel to dApp")
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
	l := chh.logger.With(
		zap.String("operation", "delete"),
		zap.String("channel", name),
		zap.String("scope", scope),
	)
	l.Debug("trying to delete a Channel")

	channel, err := chh.Get(scope, name)

	if err != nil {
		newError := ierrors.Wrap(
			err,
			"channel %s not found", name,
		)
		return newError
	}

	l.Debug("checking if Channel can be deleted")
	if len(channel.ConnectedApps) > 0 || len(channel.ConnectedAliases) > 0 {
		l.Debug("unable to delete Channel for it's being used",
			zap.Any("connected-dapps", channel.ConnectedApps),
			zap.Any("connected-aliases", channel.ConnectedAliases))

		return ierrors.New(
			"channel cannot be deleted as it is being used by other apps",
		).BadRequest()
	}

	parentApp, _ := chh.Apps().Get(scope)

	insprType := parentApp.Spec.Types[channel.Spec.Type]

	l.Debug("removing Channel from Type connected channels list",
		zap.String("type", insprType.Meta.Name))

	insprType.ConnectedChannels = utils.Remove(
		insprType.ConnectedChannels,
		channel.Meta.Name,
	)

	l.Debug("removing Channel from its parents 'Channels' structure")

	delete(parentApp.Spec.Channels, name)

	return nil
}

// Update receives a scope and a channel pointer. The scope
// defines the path to the App that will have the Update. If the App has
// a channel pointer that has the same name as that passed as an argument,
// this pointer will be replaced by the new one
func (chh *ChannelMemoryManager) Update(scope string, ch *meta.Channel) error {
	l := logger.With(
		zap.String("operation", "update"),
		zap.String("channel", ch.Meta.Name),
		zap.String("scope", scope),
	)
	l.Debug("trying to update a Channel")

	oldCh, err := chh.Get(scope, ch.Meta.Name)
	if err != nil {
		newError := ierrors.Wrap(
			err,
			"channel %s not found", ch.Meta.Name,
		)
		return newError
	}

	ch.ConnectedApps = oldCh.ConnectedApps
	ch.ConnectedAliases = oldCh.ConnectedAliases
	ch.Meta.UUID = oldCh.Meta.UUID

	ch.Spec.SelectedBroker = oldCh.Spec.SelectedBroker

	parentApp, _ := chh.Apps().Get(scope)

	l.Debug("validating new Channel structure")

	if _, ok := parentApp.Spec.Types[ch.Spec.Type]; !ok {
		l.Debug("unable to create Channel for it references an invalid Type",
			zap.String("type", ch.Spec.Type))

		return ierrors.New(
			"references a Type that doesn't exist",
		).InvalidChannel()
	}

	l.Debug("replacing old Channel with the new one in dApps 'Channels'")

	parentApp.Spec.Channels[ch.Meta.Name] = ch

	return nil
}

// ChannelPermTreeGetter returns a getter that gets channels from the root structure of the app, without the current changes.
// The getter does not allow changes in the structure, just visualization.
type ChannelPermTreeGetter struct {
	*PermTreeGetter
	logs *zap.Logger
}

// Get receives a query string (format = 'x.y.z') and iterates through the
// memory tree until it finds the Channel which name is equal to the last query element.
// If the specified Channel is found, it is returned. Otherwise, returns an error.
// This method is used to get the structure as it is in the cluster, before any modifications.
func (cmm *ChannelPermTreeGetter) Get(scope, name string) (*meta.Channel, error) {
	l := cmm.logs.With(
		zap.String("operation", "get-root"),
		zap.String("channel", name),
		zap.String("scope", scope),
	)
	l.Debug("trying to get a Channel ")

	parentApp, err := cmm.Apps().Get(scope)
	if err != nil {
		newError := ierrors.Wrap(
			err,
			"the '%v' scope is invalid", scope,
		)
		return nil, newError
	}

	if parentApp.Spec.Channels != nil {
		if ch, ok := parentApp.Spec.Channels[name]; ok {
			return ch, nil
		}
	}

	l.Debug("unable to get Channel in given scope ")

	newError := ierrors.New("channel not found ").NotFound()
	return nil, newError
}
