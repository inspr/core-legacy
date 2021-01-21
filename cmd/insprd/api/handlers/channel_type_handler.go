package handler

import (
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

// ChannelTypeHandler todo doc
type ChannelTypeHandler struct {
	memory.ChannelTypeMemory
}

// NewChannelTypeHandler todo fix
func NewChannelTypeHandler(memManager memory.Manager) *ChannelTypeHandler {
	return &ChannelTypeHandler{
		ChannelTypeMemory: memManager.ChannelTypes(),
	}
}

// HandleGetAllChannelTypes returns all ChannelTypes that exists
func (ch *ChannelTypeHandler) HandleGetAllChannelTypes() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// HandleCreateInfo informs the data needed to create a App
func (ch *ChannelTypeHandler) HandleCreateInfo() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// HandleCreateChannelType todo doc
func (ch *ChannelTypeHandler) HandleCreateChannelType() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// HandleGetChannelTypeByRef todo doc
func (ch *ChannelTypeHandler) HandleGetChannelTypeByRef() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// UpdateChannelType todo doc
func (ch *ChannelTypeHandler) UpdateChannelType() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// HandleDeleteChannelType todo doc
func (ch *ChannelTypeHandler) HandleDeleteChannelType() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}
