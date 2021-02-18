package fake

import (
	"fmt"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// ChannelTypes - mocks the implementation of the ChannelTypeMemory interface methods
type ChannelTypes struct {
	*MemManager
	fail         error
	channelTypes map[string]*meta.ChannelType
}

// GetChannelType - simple mock
func (chType *ChannelTypes) GetChannelType(context string, ctName string) (*meta.ChannelType, error) {
	if chType.fail != nil {
		return nil, chType.fail
	}
	query := fmt.Sprintf("%s.%s", context, ctName)
	ct, ok := chType.channelTypes[query]
	if !ok {
		return nil, ierrors.NewError().NotFound().Message(fmt.Sprintf("channel type %s not found", query)).Build()
	}
	return ct, nil
}

// CreateChannelType - simple mock
func (chType *ChannelTypes) CreateChannelType(context string, ct *meta.ChannelType) error {
	if chType.fail != nil {
		return chType.fail
	}
	query := fmt.Sprintf("%s.%s", context, ct.Meta.Name)
	_, ok := chType.channelTypes[query]
	if ok {
		return ierrors.NewError().AlreadyExists().Message(fmt.Sprintf("channel type %s already exists", query)).Build()
	}
	chType.channelTypes[query] = ct
	return nil
}

// DeleteChannelType - simple mock
func (chType *ChannelTypes) DeleteChannelType(context string, ctName string) error {
	if chType.fail != nil {
		return chType.fail
	}
	query := fmt.Sprintf("%s.%s", context, ctName)
	_, ok := chType.channelTypes[query]
	if !ok {
		return ierrors.NewError().NotFound().Message(fmt.Sprintf("channel type %s not found", query)).Build()
	}

	delete(chType.channelTypes, query)
	return nil
}

// UpdateChannelType - simple mock
func (chType *ChannelTypes) UpdateChannelType(context string, ct *meta.ChannelType) error {
	if chType.fail != nil {
		return chType.fail
	}
	query := fmt.Sprintf("%s.%s", context, ct.Meta.Name)
	_, ok := chType.channelTypes[query]
	if !ok {
		return ierrors.NewError().NotFound().Message(fmt.Sprintf("channel type %s not found", query)).Build()
	}
	chType.channelTypes[query] = ct
	return nil
}