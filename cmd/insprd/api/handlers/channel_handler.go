package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
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
		// not implemented yet
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
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}
		data := struct {
			Channel meta.Channel `json:"channel"`
			Ctx     string       `json:"ctx"`
		}{}
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

// HandleGetChannelByRef todo doc
func (ch *ChannelHandler) HandleGetChannelByRef() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}
		data := struct {
			Query string `json:"query"`
		}{}
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

// HandleUpdateChannel todo doc
func (ch *ChannelHandler) HandleUpdateChannel() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}
		data := struct {
			Channel meta.Channel `json:"channel"`
			Ctx     string       `json:"ctx"`
		}{}
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

// HandleDeleteChannel todo doc
func (ch *ChannelHandler) HandleDeleteChannel() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}
		data := struct {
			Query string `json:"query"`
		}{}
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
