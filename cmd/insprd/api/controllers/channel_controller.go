package controller

import (
	ierrors "gitlab.inspr.dev/inspr/core/pkg/error"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// GetChannel returns the channels found in the path './channel/path'
func (s *Server) GetChannel(query string) (*meta.Channel, error) {
	channel, err := s.MemoryManager.Channels().GetChannel(query)
	if err != nil {
		serverErr := ierrors.NewError().InnerError(err).Message(
			"Couldn't get a channel in this path",
		).InternalServer().Build()
		return nil, serverErr
	}
	return channel, nil
}

// CreateChannel creates a new channel in the stored structure of the cluster
func (s *Server) CreateChannel(ch *meta.Channel, ctx string) error {
	err := s.MemoryManager.Channels().CreateChannel(ch, ctx)
	if err != nil {
		serverErr := ierrors.NewError().InnerError(err).Message(
			"Couldn't create a channel in this context and with these values",
		).InternalServer().Build()
		return serverErr
	}
	return nil
}

// DeleteChannel todo doc
func (s *Server) DeleteChannel(query string) error {
	err := s.MemoryManager.Channels().DeleteChannel(query)
	if err != nil {
		serverErr := ierrors.NewError().InnerError(err).Message(
			"Couldn't delete the channel in this path",
		).InternalServer().Build()
		return serverErr
	}
	return nil
}

// UpdateChannel todo doc
func (s *Server) UpdateChannel(ch *meta.Channel, ctx string) error {
	err := s.MemoryManager.Channels().CreateChannel(ch, ctx)
	if err != nil {
		serverErr := ierrors.NewError().InnerError(err).Message(
			"Couldn't modify the channel in this context",
		).InternalServer().Build()
		return serverErr
	}
	return nil
}
