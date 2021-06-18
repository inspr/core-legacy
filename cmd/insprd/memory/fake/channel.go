package fake

import (
	"fmt"

	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
)

// Channels - mocks the implementation of the ChannelMemory interface methods
type Channels struct {
	*MemManager
	fail     error
	channels map[string]*meta.Channel
}

// Get - simple mock
func (ch *Channels) Get(scope, name string) (*meta.Channel, error) {
	if ch.fail != nil {
		return nil, ch.fail
	}
	query := fmt.Sprintf("%s.%s", scope, name)
	ct, ok := ch.channels[query]
	if !ok {
		return nil, ierrors.
			NewError().
			NotFound().
			Message("channel %s not found", query).
			Build()
	}
	return ct, nil
}

// Create - simple mock
func (ch *Channels) Create(scope string, channel *meta.Channel) error {
	if ch.fail != nil {
		return ch.fail
	}
	query := fmt.Sprintf("%s.%s", scope, channel.Meta.Name)
	_, ok := ch.channels[query]
	if ok {
		return ierrors.
			NewError().
			AlreadyExists().
			Message("channel %s already exists", query).
			Build()
	}
	ch.channels[query] = channel
	return nil
}

// Delete - simple mock
func (ch *Channels) Delete(scope, name string) error {
	if ch.fail != nil {
		return ch.fail
	}
	query := fmt.Sprintf("%s.%s", scope, name)
	_, ok := ch.channels[query]
	if !ok {
		return ierrors.
			NewError().
			NotFound().
			Message("channel %s not found", query).
			Build()
	}

	delete(ch.channels, query)
	return nil
}

// Update - simple mock
func (ch *Channels) Update(scope string, channel *meta.Channel) error {
	if ch.fail != nil {
		return ch.fail
	}
	query := fmt.Sprintf("%s.%s", scope, channel.Meta.Name)
	_, ok := ch.channels[query]
	if !ok {
		return ierrors.
			NewError().
			NotFound().
			Message("channel %s not found", query).
			Build()
	}
	ch.channels[query] = channel
	return nil
}
