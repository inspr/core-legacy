package handler

import (
	"net/http"

	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

type channelTypeHandler struct{}

// ChannelTypeHandler is an empty object to separate route functions
var ChannelTypeHandler channelTypeHandler = struct{}{}

// GetAllChannelTypes returns all ChannelTypes that exists
func (ch *channelTypeHandler) GetAllChannelTypes() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// CreateInfo informs the data needed to create a App
func (ch *channelTypeHandler) CreateInfo() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// CreateChannelType todo doc
func (ch *channelTypeHandler) CreateChannelType() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// GetChannelTypeByRef todo doc
func (ch *channelTypeHandler) GetChannelTypeByRef() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// UpdateChannelType todo doc
func (ch *channelTypeHandler) UpdateChannelType() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// DeleteChannelType todo doc
func (ch *channelTypeHandler) DeleteChannelType() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}
