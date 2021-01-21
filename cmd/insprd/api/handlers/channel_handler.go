package handler

import (
	"fmt"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

// ChannelHandler todo doc
type ChannelHandler struct {
	memory.ChannelMemory
}

// NewChannelHandler exports
func NewChannelHandler(memManager memory.Manager) *ChannelHandler {
	return &ChannelHandler{
		ChannelMemory: memManager.Channels(),
	}
}

// HandleGetAllChannels returns all channels
func (ch *ChannelHandler) HandleGetAllChannels() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// HandleCreateInfo informs the data needed to create a channel
func (ch *ChannelHandler) HandleCreateInfo() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// HandleCreateChannel todo doc
func (ch *ChannelHandler) HandleCreateChannel() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		// couldn't parse the request body
		if err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
		}

		fmt.Fprintf(w, "Create Channel\n")
		fmt.Fprintf(w, "Info received %v\n", r.Form.Get("info"))
	}
	return rest.Handler(handler)
}

// HandleGetChannelByRef todo doc
func (ch *ChannelHandler) HandleGetChannelByRef() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		// couldn't parse the request body
		if err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
		}

		fmt.Fprintf(w, "Get Channel By Ref\n")
		fmt.Fprintf(w, "Id received %v\n", r.Form.Get("id"))
	}
	return rest.Handler(handler)
}

// HandleUpdateChannel todo doc
func (ch *ChannelHandler) HandleUpdateChannel() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// HandleDeleteChannel todo doc
func (ch *ChannelHandler) HandleDeleteChannel() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}
