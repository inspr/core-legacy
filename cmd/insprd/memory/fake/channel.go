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
func (ch *Channels) Get(context string, name string) (*meta.Channel, error) {
	if ch.fail != nil {
		return nil, ch.fail
	}
	query := fmt.Sprintf("%s.%s", context, name)
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
func (ch *Channels) Create(context string, ct *meta.Channel) error {
	if ch.fail != nil {
		return ch.fail
	}
	query := fmt.Sprintf("%s.%s", context, ct.Meta.Name)
	_, ok := ch.channels[query]
	if ok {
		return ierrors.
			NewError().
			AlreadyExists().
			Message("channel %s already exists", query).
			Build()
	}
	ch.channels[query] = ct
	return nil
}

// Delete - simple mock
func (ch *Channels) Delete(context string, name string) error {
	if ch.fail != nil {
		return ch.fail
	}
	query := fmt.Sprintf("%s.%s", context, name)
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
func (ch *Channels) Update(context string, ct *meta.Channel) error {
	if ch.fail != nil {
		return ch.fail
	}
	query := fmt.Sprintf("%s.%s", context, ct.Meta.Name)
	_, ok := ch.channels[query]
	if !ok {
		return ierrors.
			NewError().
			NotFound().
			Message("channel %s not found", query).
			Build()
	}
	ch.channels[query] = ct
	return nil
}
