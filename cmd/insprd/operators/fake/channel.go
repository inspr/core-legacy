package fake

import (
	"context"
	"strings"

	"inspr.dev/inspr/cmd/insprd/operators"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
)

// ChannelOperator mock
type ChannelOperator struct {
	channels map[string]*meta.Channel
	err      error
}

// NewChannelOperator mock
func NewChannelOperator(err error) operators.ChannelOperatorInterface {
	return ChannelOperator{
		channels: make(map[string]*meta.Channel),
		err:      err,
	}
}

// Create mock
func (o ChannelOperator) Create(
	ctx context.Context,
	context string,
	ch *meta.Channel,
) error {
	if o.err != nil {
		return o.err
	}
	if _, ok := o.channels[context+ch.Meta.Name]; ok {
		return ierrors.New("channel already exists").AlreadyExists()
	}
	o.channels[context+ch.Meta.Name] = ch
	return nil
}

// Get mock
func (o ChannelOperator) Get(
	ctx context.Context,
	context string,
	name string,
) (*meta.Channel, error) {
	if o.err != nil {
		return nil, o.err
	}
	channelKey := context + name
	ch, ok := o.channels[channelKey]
	if !ok {
		return nil, ierrors.New("channel %s not found", channelKey).NotFound()
	}
	return ch, nil
}

// Update mock
func (o ChannelOperator) Update(
	ctx context.Context,
	context string,
	ch *meta.Channel,
) error {
	if o.err != nil {
		return o.err
	}

	channelKey := context + ch.Meta.Name
	if _, ok := o.channels[channelKey]; !ok {
		return ierrors.New("channel %s not found", channelKey).NotFound()
	}
	o.channels[channelKey] = ch
	return nil
}

// Delete mock
func (o ChannelOperator) Delete(
	ctx context.Context,
	context string,
	name string,
) error {
	if o.err != nil {
		return o.err
	}

	channelKey := context + name
	_, ok := o.channels[channelKey]
	if !ok {
		return ierrors.New("channel %s not found", channelKey).NotFound()
	}
	delete(o.channels, channelKey)
	return nil
}

// GetAll mock
func (o ChannelOperator) GetAll(
	_ context.Context,
	context string,
) (ret []*meta.Channel, err error) {

	for _, ch := range o.channels {
		if strings.HasPrefix(ch.Meta.Parent, context) {
			ret = append(ret, ch)
		}
	}
	return
}
