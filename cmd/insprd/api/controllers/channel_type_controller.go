package controller

import (
	ierrors "gitlab.inspr.dev/inspr/core/pkg/error"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// GetChannelType todo doc
func (s *Server) GetChannelType(query string) (*meta.ChannelType, error) {
	channelType, err := s.MemoryManager.ChannelTypes().GetChannelType(query)
	if err != nil {
		serverErr := ierrors.NewError().InnerError(err).Message(
			"Couldn't get a channel-type in this path",
		).InternalServer().Build()
		return nil, serverErr
	}
	return channelType, nil
}

// CreateChannelType todo doc
func (s *Server) CreateChannelType(ct *meta.ChannelType, ctx string) error {
	err := s.MemoryManager.ChannelTypes().CreateChannelType(ct, ctx)
	if err != nil {
		serverErr := ierrors.NewError().InnerError(err).Message(
			"Couldn't create a channel-type in this context and with these values",
		).InternalServer().Build()
		return serverErr
	}
	return nil
}

// DeleteChannelType todo doc
func (s *Server) DeleteChannelType(query string) error {
	err := s.MemoryManager.ChannelTypes().DeleteChannelType(query)
	if err != nil {
		serverErr := ierrors.NewError().InnerError(err).Message(
			"Couldn't delete the channel-type in this path",
		).InternalServer().Build()
		return serverErr
	}
	return nil
}

// UpdateChannelType todo doc
func (s *Server) UpdateChannelType(ct *meta.ChannelType, ctx string) error {
	err := s.MemoryManager.ChannelTypes().UpdateChannelType(ct, ctx)
	if err != nil {
		serverErr := ierrors.NewError().InnerError(err).Message(
			"Couldn't modify the channel-type in this context",
		).InternalServer().Build()
		return serverErr
	}
	return nil
}
