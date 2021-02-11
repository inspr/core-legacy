package fake

import (
	"fmt"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// Channels - mocks the implementation of the ChannelMemory interface methods
type Channels struct {
	fail     error
	channels map[string]*meta.Channel
}

// GetChannel - simple mock
func (ch *Channels) GetChannel(context string, name string) (*meta.Channel, error) {
	if ch.fail != nil {
		return nil, ch.fail
	}
	query := fmt.Sprintf("%s.%s", context, name)
	ct, ok := ch.channels[query]
	if !ok {
		return nil, ierrors.NewError().NotFound().Message(fmt.Sprintf("channel type %s not found", query)).Build()
	}
	return ct, nil
}

// CreateChannel - simple mock
func (ch *Channels) CreateChannel(context string, ct *meta.Channel) error {
	if ch.fail != nil {
		return ch.fail
	}
	query := fmt.Sprintf("%s.%s", context, ct.Meta.Name)
	_, ok := ch.channels[query]
	if ok {
		return ierrors.NewError().AlreadyExists().Message(fmt.Sprintf("channel type %s already exists", query)).Build()
	}
	ch.channels[query] = ct
	return nil
}

// DeleteChannel - simple mock
func (ch *Channels) DeleteChannel(context string, name string) error {
	if ch.fail != nil {
		return ch.fail
	}
	query := fmt.Sprintf("%s.%s", context, name)
	_, ok := ch.channels[query]
	if !ok {
		return ierrors.NewError().NotFound().Message(fmt.Sprintf("channel type %s not found", query)).Build()
	}

	delete(ch.channels, query)
	return nil
}

// UpdateChannel - simple mock
func (ch *Channels) UpdateChannel(context string, ct *meta.Channel) error {
	if ch.fail != nil {
		return ch.fail
	}
	query := fmt.Sprintf("%s.%s", context, ct.Meta.Name)
	_, ok := ch.channels[query]
	if !ok {
		return ierrors.NewError().NotFound().Message(fmt.Sprintf("channel type %s not found", query)).Build()
	}
	ch.channels[query] = ct
	return nil
}
