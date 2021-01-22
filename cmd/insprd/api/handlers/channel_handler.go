package handler

import (
	"encoding/json"
	"io/ioutil"
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

// HandleCreateInfo informs the data needed to create a channel
func (ch *ChannelHandler) HandleCreateInfo() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
	}
	return rest.Handler(handler)
}

// HandleCreateChannel - returns the handle function that
// manages the creation of a channel
func (ch *ChannelHandler) HandleCreateChannel() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}
		data := models.ChannelDI{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		err = ch.CreateChannel(&data.Channel, data.Ctx)
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
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}
		data := models.ChannelQueryDI{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}
		channel, err := ch.GetChannel(data.Query)
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
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}
		data := models.ChannelDI{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		err = ch.UpdateChannel(&data.Channel, data.Ctx)
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
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}
		data := models.ChannelQueryDI{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}
		err = ch.DeleteChannel(data.Query)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}
