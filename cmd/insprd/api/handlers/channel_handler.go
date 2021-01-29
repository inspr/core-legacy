package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

// ChannelHandler - contains handlers that uses the ChannelMemory interface methods
type ChannelHandler struct {
	memory.ChannelMemory
}

// NewChannelHandler exports
func NewChannelHandler(memManager memory.Manager) *ChannelHandler {
	return &ChannelHandler{
		ChannelMemory: memManager.Channels(),
	}
}

// HandleCreateChannel - returns the handle function that
// manages the creation of a channel
func (ch *ChannelHandler) HandleCreateChannel() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		err = ch.CreateChannel(data.Ctx, &data.Channel)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}

// HandleGetChannelByRef - return a handle function that obtains
// a channel by the reference given
func (ch *ChannelHandler) HandleGetChannelByRef() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		channel, err := ch.GetChannel(data.Ctx, data.ChName)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		rest.JSON(w, http.StatusOK, channel)
	}
	return rest.Handler(handler)
}

// HandleUpdateChannel - returns a handle function that
// updates the channel with the parameters given in the request
func (ch *ChannelHandler) HandleUpdateChannel() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		err = ch.UpdateChannel(data.Ctx, &data.Channel)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}

// HandleDeleteChannel - returns a handle function that
// deletes the channel of the given path
func (ch *ChannelHandler) HandleDeleteChannel() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		err = ch.DeleteChannel(data.Ctx, data.ChName)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}
