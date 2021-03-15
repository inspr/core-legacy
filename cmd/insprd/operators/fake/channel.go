package fake

import (
	"context"
	"fmt"
	"strings"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
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
func (o ChannelOperator) Create(ctx context.Context, context string, ch *meta.Channel) error {
	fmt.Println(ch.Meta.Name)
	if o.err != nil {
		return o.err
	}
	if _, ok := o.channels[context+ch.Meta.Name]; ok {
		return ierrors.NewError().AlreadyExists().Message("channel already exists").Build()
	}
	o.channels[context+ch.Meta.Name] = ch
	return nil
}

// Get mock
func (o ChannelOperator) Get(ctx context.Context, context string, name string) (*meta.Channel, error) {
	if o.err != nil {
		return nil, o.err
	}
	ch, ok := o.channels[context+name]
	if !ok {
		return nil, ierrors.NewError().NotFound().Message("channel not found").Build()
	}
	return ch, nil
}

// Update mock
func (o ChannelOperator) Update(ctx context.Context, context string, ch *meta.Channel) error {
	if o.err != nil {
		return o.err
	}
	if _, ok := o.channels[context+ch.Meta.Name]; !ok {
		return ierrors.NewError().NotFound().Message("channel not found").Build()
	}
	o.channels[context+ch.Meta.Name] = ch
	return nil
}

// Delete mock
func (o ChannelOperator) Delete(ctx context.Context, context string, name string) error {
	if o.err != nil {
		return o.err
	}
	_, ok := o.channels[context+name]
	if !ok {
		return ierrors.NewError().NotFound().Message("channel not found").Build()
	}
	delete(o.channels, context+name)
	return nil
}

// GetAll mock
func (o ChannelOperator) GetAll(_ context.Context, context string) (ret []*meta.Channel, err error) {

	for _, ch := range o.channels {
		if strings.HasPrefix(ch.Meta.Parent, context) {
			ret = append(ret, ch)
		}
	}
	return
}