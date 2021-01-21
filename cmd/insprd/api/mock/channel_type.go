package mock

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// ChannelTypes todo doc
type ChannelTypes struct {
	memory.ChannelTypeMemory
}

// GetChannelType todo doc
func (chType *ChannelTypes) GetChannelType(query string) (*meta.ChannelType, error) {
	return &meta.ChannelType{}, nil
}

// CreateChannelType todo doc
func (chType *ChannelTypes) CreateChannelType(ct *meta.ChannelType, context string) error {
	return nil
}

// DeleteChannelType todo doc
func (chType *ChannelTypes) DeleteChannelType(query string) error {
	return nil
}

// UpdateChannelType todo doc
func (chType *ChannelTypes) UpdateChannelType(ct *meta.ChannelType, query string) error {
	return nil
}
