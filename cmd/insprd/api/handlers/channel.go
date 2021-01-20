package handler

import (
	"net/http"

	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

type channelHandler struct{}

// ChannelHandler is an empty object to separate route functions
var ChannelHandler channelHandler = struct{}{}

// GetAllChannels returns all channels
func (ch *channelHandler) GetAllChannels() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// CreateInfo informs the data needed to create a channel
func (ch *channelHandler) CreateInfo() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// CreateChannel todo doc
func (ch *channelHandler) CreateChannel() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// GetChannelByRef todo doc
func (ch *channelHandler) GetChannelByRef() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// UpdateChannel todo doc
func (ch *channelHandler) UpdateChannel() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// GetAllChannels todo doc
func (ch *channelHandler) DeleteChannel() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}
