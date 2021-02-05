package tree

import (
	"reflect"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

/*
ChannelMemoryManager implements the channel interface and
provides methods for operating on Channels
*/
type ChannelMemoryManager struct {
	root *meta.App
}

// Channels return a pointer to ChannelMemoryManager
func (tmm *MemoryManager) Channels() memory.ChannelMemory {
	return &ChannelMemoryManager{
		root: tmm.root,
	}
}

/*
GetChannel receives a context and a channel name. The context defines
the path to an App. If this App has a pointer to a channel that has the
same name as the name passed as an argument, the pointer to that channel is returned
*/
func (chh *ChannelMemoryManager) GetChannel(context string, chName string) (*meta.Channel, error) {
	parentApp, err := GetTreeMemory().Apps().GetApp(context)
	if err != nil {
		newError := ierrors.NewError().InnerError(err).NotFound().Message("channel was not found because the app context has an error").Build()
		return nil, newError
	}

	if parentApp.Spec.Channels != nil {
		if ch, ok := parentApp.Spec.Channels[chName]; ok {
			return ch, nil
		}
	}

	newError := ierrors.NewError().NotFound().Message("channel not found").Build()
	return nil, newError
}

/*
CreateChannel receives a context that defines a path to the App
in which to add a pointer to the channel passed as an argument
*/
func (chh *ChannelMemoryManager) CreateChannel(context string, ch *meta.Channel) error {
	chAlreadyExist, _ := chh.GetChannel(context, ch.Meta.Name)
	if chAlreadyExist != nil {
		return ierrors.NewError().AlreadyExists().Message("channel with name " + ch.Meta.Name + " already exists in the context " + context).Build()
	}

	parentApp, err := GetTreeMemory().Apps().GetApp(context)
	if err != nil {
		newError := ierrors.NewError().InnerError(err).InvalidChannel().Message("app for channel creation not found").Build()
		return newError
	}

	if _, ok := parentApp.Spec.ChannelTypes[ch.Spec.Type]; !ok {
		return ierrors.NewError().InvalidChannel().Message("references a Channel Type that doesn't exist").Build()
	}

	connectedChannels := parentApp.Spec.ChannelTypes[ch.Spec.Type].ConnectedChannels
	if !utils.Include(connectedChannels, ch.Meta.Name) {
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
DeleteChannel receives a context and a channel name. The context
defines the path to the App that will have the Delete. If the App
has a pointer to a channel that has the same name as the name passed
as an argument, that pointer is removed from the list of App channels
*/
func (chh *ChannelMemoryManager) DeleteChannel(context string, chName string) error {
	channel, err := chh.GetChannel(context, chName)
	if err != nil {
		newError := ierrors.NewError().InnerError(err).NotFound().Message("channel not found").Build()
		return newError
	}

	if len(channel.ConnectedApps) > 0 {
		return ierrors.NewError().
			BadRequest().
			Message("channel cannot be deleted as it is being used by other apps").
			Build()
	}

	parentApp, _ := GetTreeMemory().Apps().GetApp(context)

	delete(parentApp.Spec.Channels, chName)

	return nil
}

/*
UpdateChannel receives a context and a channel pointer. The context
defines the path to the App that will have the Update. If the App has
a channel pointer that has the same name as that passed as an argument,
this pointer will be replaced by the new one
*/
func (chh *ChannelMemoryManager) UpdateChannel(context string, ch *meta.Channel) error {
	oldCh, err := chh.GetChannel(context, ch.Meta.Name)
	if err != nil {
		newError := ierrors.NewError().InnerError(err).NotFound().Message("channel not found").Build()
		return newError
	}

	if ok := reflect.DeepEqual(oldCh.ConnectedApps, ch.ConnectedApps); !ok {
		return ierrors.NewError().
			InvalidChannel().
			Message("new channel must have the same connectedApps as old channel").
			Build()
	}

	parentApp, _ := GetTreeMemory().Apps().GetApp(context)

	if _, ok := parentApp.Spec.ChannelTypes[ch.Spec.Type]; !ok {
		return ierrors.NewError().InvalidChannel().Message("references a Channel Type that doesn't exist").Build()
	}

	parentApp.Spec.Channels[ch.Meta.Name] = ch

	return nil
}
