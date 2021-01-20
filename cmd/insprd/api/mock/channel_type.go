package repos

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// ChannelTypes todo doc
type ChannelTypes struct {
	memory.ChannelTypeMemory
}

// GetChannelType todo doc
func (chType *ChannelTypes) GetChannelType(ref string) (*meta.ChannelType, error) {
	return &meta.ChannelType{}, nil
}

// CreateChannelType todo doc
func (chType *ChannelTypes) CreateChannelType(ct *meta.ChannelType) error {
	return nil
}

// DeleteChannelType todo doc
func (chType *ChannelTypes) DeleteChannelType(ref string) error {
	return nil
}

// UpdateChannelType todo doc
func (chType *ChannelTypes) UpdateChannelType(ct *meta.ChannelType, ref string) error {
	return nil
}
