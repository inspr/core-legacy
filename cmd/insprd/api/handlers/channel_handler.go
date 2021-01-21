package handler

import (
	"encoding/json"
	"fmt"
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
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Couldn't process the request body")
			return
		}
		data := struct {
			Channel meta.Channel `json:"channel"`
			Ctx     string       `json:"ctx"`
		}{}
		json.Unmarshal(body, &data)

		err = ch.CreateChannel(&data.Channel, data.Ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
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
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Couldn't process the request body")
			return
		}
		data := struct {
			Query string `json:"query"`
		}{}
		json.Unmarshal(body, &data)
		app, err := ch.GetChannel(data.Query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		}

		// respond with json
		fmt.Fprintf(w, "%v", app)
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}

// HandleUpdateChannel todo doc
func (ch *ChannelHandler) HandleUpdateChannel() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Couldn't process the request body")
			return
		}
		data := struct {
			Channel meta.Channel `json:"channel"`
			Ctx     string       `json:"ctx"`
		}{}
		json.Unmarshal(body, &data)

		err = ch.UpdateChannel(&data.Channel, data.Ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
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
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Couldn't process the request body")
			return
		}
		data := struct {
			Query string `json:"query"`
		}{}
		json.Unmarshal(body, &data)
		err = ch.DeleteChannel(data.Query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}
