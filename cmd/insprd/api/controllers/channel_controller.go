package controller

import "gitlab.inspr.dev/inspr/core/pkg/meta"

// GetAllChannels return all channels in the cluster's memory
func (s *Server) GetAllChannels() {
	// channel, err := s.MemoryManager.Channels().GetChannel(".")
	// if err != nil {
	// 	// handle error
	// }
}

// GetChannel returns the channels found in the path './channel/path'
func (s *Server) GetChannel(ref string) (*meta.Channel, error) {
	channel, err := s.MemoryManager.Channels().GetChannel(ref)
	if err != nil {
		// error package
	}
	return channel, nil
}

// CreateChannel creates a new channel in the stored structure of the cluster
func (s *Server) CreateChannel(ch *meta.Channel) error {
	err := s.MemoryManager.Channels().CreateChannel(ch)
	if err != nil {

	}
	return nil
}

// CreateChannel(ch *meta.Channel) error
// DeleteChannel(ref string) error
// UpdateChannel(ch *meta.Channel, ref string) error
