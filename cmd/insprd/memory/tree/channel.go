package tree

import (
	"fmt"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	metautils "gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

/*
ChannelMemoryManager implements the channel interface and
provides methods for operating on Channels
*/
type ChannelMemoryManager struct {
	*MemoryManager
}

// Channels return a pointer to ChannelMemoryManager
func (tmm *MemoryManager) Channels() memory.ChannelMemory {
	return &ChannelMemoryManager{
		MemoryManager: tmm,
	}
}

/*
Get receives a context and a channel name. The context defines
the path to an App. If this App has a pointer to a channel that has the
same name as the name passed as an argument, the pointer to that channel is returned
*/
func (chh *ChannelMemoryManager) Get(context string, chName string) (*meta.Channel, error) {
	parentApp, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		newError := ierrors.
			NewError().
			InnerError(err).
			NotFound().
			Message("channel was not found because the app context has an error").
			Build()
		return nil, newError
	}

	if parentApp.Spec.Channels != nil {
		if ch, ok := parentApp.Spec.Channels[chName]; ok {
			return ch, nil
		}
	}

	newError := ierrors.
		NewError().
		NotFound().
		Message(fmt.Sprintf("channel %s not found", chName)).
		Build()
	return nil, newError
}

/*
Create receives a context that defines a path to the App
in which to add a pointer to the channel passed as an argument
*/
func (chh *ChannelMemoryManager) Create(context string, ch *meta.Channel) error {
	nameErr := metautils.StructureNameIsValid(ch.Meta.Name)
	if nameErr != nil {
		return ierrors.NewError().InnerError(nameErr).Message(nameErr.Error()).Build()
	}

	chAlreadyExist, _ := chh.Get(context, ch.Meta.Name)
	if chAlreadyExist != nil {
		return ierrors.NewError().AlreadyExists().Message("channel with name " + ch.Meta.Name + " already exists in the context " + context).Build()
	}

	parentApp, err := GetTreeMemory().Apps().Get(context)
	if err != nil {
		newError := ierrors.NewError().InnerError(err).InvalidChannel().
			Message("couldn't create channel " + ch.Meta.Name + "\n" + err.Error()).
			Build()
		return newError
	}

	if _, ok := parentApp.Spec.ChannelTypes[ch.Spec.Type]; !ok {
		return ierrors.NewError().InvalidChannel().Message("references a Channel Type that doesn't exist").Build()
	}

	connectedChannels := parentApp.Spec.ChannelTypes[ch.Spec.Type].ConnectedChannels
	if !utils.Includes(connectedChannels, ch.Meta.Name) {
		connectedChannels = append(connectedChannels, ch.Meta.Name)
		parentApp.Spec.ChannelTypes[ch.Spec.Type].ConnectedChannels = connectedChannels
	}

	if parentApp.Spec.Channels == nil {
		parentApp.Spec.Channels = map[string]*meta.Channel{}
	}
	parentApp.Spec.Channels[ch.Meta.Name] = ch

	return nil
}

/*
Delete receives a context and a channel name. The context
defines the path to the App that will have the Delete. If the App
has a pointer to a channel that has the same name as the name passed
as an argument, that pointer is removed from the list of App channels
*/
func (chh *ChannelMemoryManager) Delete(context string, chName string) error {
	channel, err := chh.Get(context, chName)
	if err != nil {
		newError := ierrors.
			NewError().
			InnerError(err).
			NotFound().
			Message(fmt.Sprintf("channel %s not found", chName)).
			Build()
		return newError
	}

	if len(channel.ConnectedApps) > 0 {
		return ierrors.NewError().
			BadRequest().
			Message("channel cannot be deleted as it is being used by other apps").
			Build()
	}

	parentApp, _ := GetTreeMemory().Apps().Get(context)

	channelType := parentApp.Spec.ChannelTypes[channel.Spec.Type]
	channelType.ConnectedChannels = utils.Remove(channelType.ConnectedChannels, channel.Meta.Name)

	delete(parentApp.Spec.Channels, chName)

	return nil
}

/*
Update receives a context and a channel pointer. The context
defines the path to the App that will have the Update. If the App has
a channel pointer that has the same name as that passed as an argument,
this pointer will be replaced by the new one
*/
func (chh *ChannelMemoryManager) Update(context string, ch *meta.Channel) error {
	oldCh, err := chh.Get(context, ch.Meta.Name)
	if err != nil {
		newError := ierrors.
			NewError().
			InnerError(err).
			NotFound().
			Message(fmt.Sprintf("channel %s not found", ch.Meta.Name)).
			Build()
		return newError
	}

	ch.ConnectedApps = oldCh.ConnectedApps

	parentApp, _ := GetTreeMemory().Apps().Get(context)

	if _, ok := parentApp.Spec.ChannelTypes[ch.Spec.Type]; !ok {
		return ierrors.NewError().InvalidChannel().Message("references a Channel Type that doesn't exist").Build()
	}

	parentApp.Spec.Channels[ch.Meta.Name] = ch

	return nil
}

// ChannelRootGetter returns a getter that gets channels from the root structure of the app, without the current changes.
// The getter does not allow changes in the structure, just visualization.
type ChannelRootGetter struct {
}

// Get receives a query string (format = 'x.y.z') and iterates through the
// memory tree until it finds the dChannel which name is equal to the last query element.
// The tree Channel is returned if the query string is an empty string.
// If the specified dChannel is found, it is returned. Otherwise, returns an error.
func (amm *ChannelRootGetter) Get(context string, chName string) (*meta.Channel, error) {
	parentApp, err := GetTreeMemory().Root().Apps().Get(context)
	if err != nil {
		newError := ierrors.NewError().InnerError(err).NotFound().Message("channel was not found because the app context has an error").Build()
		return nil, newError
	}

	if parentApp.Spec.Channels != nil {
		if ch, ok := parentApp.Spec.Channels[chName]; ok {
			return ch, nil
		}
	}

	newError := ierrors.
		NewError().
		NotFound().
		Message(fmt.Sprintf("channel %s not found", chName)).
		Build()
	return nil, newError
}
