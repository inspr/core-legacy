package operators

import "gitlab.inspr.dev/inspr/core/pkg/meta"

// ChannelOperatorInterface is responsible for handling the following methods
//
// 	- `Get`: returns a channel from the DApp of the given context
//	- `GetAll`: return all channels from the DApp of the given context
// 	- `Create`: creates a channel in the DApp of the given context
// 	- `Update`: updates a channel in the DApp of the given context
// 	- `Delete`: deletes a channel of the specified name in DApp of the given context
type ChannelOperatorInterface interface {
	Get(ctx string, name string) (*meta.Channel, error)
	GetAll(ctx string) ([]*meta.Channel, error)
	Create(ctx string, channel *meta.Channel) error
	Update(ctx string, channel *meta.Channel) error
	Delete(ctx string, name string) error
}

// ChannelOperator is the struct that has the elements
// of ChannelOperator interface implemented
type ChannelOperator struct{}

// NewChannelOperator returns a new ChannelOperator, struct
// that implements the ChannelOperator interface methods
func NewChannelOperator() *ChannelOperator {
	return &ChannelOperator{}
}

// Get - returns a channel given the search context and it's name
func (co *ChannelOperator) Get(ctx string, name string) (*meta.Channel, error) {
	return &meta.Channel{}, nil
}

// GetAll - returns all channels in the given context
func (co *ChannelOperator) GetAll(ctx string) ([]*meta.Channel, error) {
	return []*meta.Channel{{}, {}}, nil
}

// Create - adds the channel passed in the parameter into the context passed
// before doing so some tests are made to be sure no conflicts happen
func (co *ChannelOperator) Create(ctx string, channel *meta.Channel) error {
	return nil
}

// Update - updates the channel in the context given, the way it identifies
// which channel to change is through channel.Meta.Name
func (co *ChannelOperator) Update(ctx string, channel *meta.Channel) error {
	return nil
}

// Delete - deletes the channel in the context given and with the name specified
func (co *ChannelOperator) Delete(ctx string, name string) error {
	return nil
}
