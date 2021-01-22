package tree

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// ChannelMemoryManager DOC TODO
type ChannelMemoryManager struct {
	root *meta.App
}

// Channels Doc TODO
func (tmm *TreeMemoryManager) Channels() memory.ChannelMemory {
	return &ChannelMemoryManager{
		root: tmm.root,
	}
}

// GetChannel DOC TODO
func (chh *ChannelMemoryManager) GetChannel(context string, chName string) (*meta.Channel, error) {
	parentApp, err := GetTreeMemory().Apps().GetApp(context)
	if err != nil {
		newError := ierrors.NewError().InnerError(err).NotFound().Message("Channel not found").Build()
		return nil, newError
	}

	ch := parentApp.Spec.Channels[chName]
	if ch != nil {
		return ch, nil
	}

	newError := ierrors.NewError().NotFound().Message("Channel not found").Build()
	return nil, newError
}

// CreateChannel DOC TODO
func (chh *ChannelMemoryManager) CreateChannel(ch *meta.Channel, context string) error {

	// Check if channel already exists
	chAlreadyExist, _ := chh.GetChannel(context, ch.Meta.Name)
	if chAlreadyExist != nil {
		return ierrors.NewError().AlreadyExists().Message("Channel with name " + ch.Meta.Name + " already exists in the context " + context).Build()
	}

	// Get context app to add the channel to it
	parentApp, err := GetTreeMemory().Apps().GetApp(context)
	if err != nil {
		newError := ierrors.NewError().InnerError(err).InvalidChannel().Message("App for channel creation not found").Build()
		return newError
	}

	// Validate Channel Structure
	if !validateChannel(ch) {
		return ierrors.NewError().InvalidChannel().Message("Invalid Channel Structure").Build()
	}

	// Add pointer to channel in the app
	parentApp.Spec.Channels[ch.Meta.Name] = ch

	return nil
}

// DeleteChannel DOC TODO
func (chh *ChannelMemoryManager) DeleteChannel(context string, chName string) error {

	// Get channel
	_, err := chh.GetChannel(context, chName)
	if err != nil {
		return err
	}

	// Get context app to delete the channel from
	parentApp, _ := GetTreeMemory().Apps().GetApp(context)

	delete(parentApp.Spec.Channels, chName)

	return nil
}

// UpdateChannel DOC TODO
func (chh *ChannelMemoryManager) UpdateChannel(ch *meta.Channel, context string) error {

	// Check if channel exists
	_, err := chh.GetChannel(context, ch.Meta.Name)
	if err != nil {
		return err
	}

	// Get context app to update the channel from
	parentApp, _ := GetTreeMemory().Apps().GetApp(context)

	// Update channel
	parentApp.Spec.Channels[ch.Meta.Name] = ch

	return nil
}

func validateChannel(ch *meta.Channel) bool {
	return true
}
