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
func (cth *ChannelTypeHandler) HandleGetAllChannelTypes() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// not implemented yet
	}
	return rest.Handler(handler)
}

// HandleCreateInfo informs the data needed to create a App
func (cth *ChannelTypeHandler) HandleCreateInfo() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		// not implemented yet
	}
	return rest.Handler(handler)
}

// HandleCreateChannelType todo doc
func (cth *ChannelTypeHandler) HandleCreateChannelType() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Couldn't process the request body")
			return
		}
		data := struct {
			ChannelType meta.ChannelType `json:"channel-type"`
			Ctx         string           `json:"ctx"`
		}{}
		json.Unmarshal(body, &data)

		err = cth.CreateChannelType(&data.ChannelType, data.Ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}

// HandleGetChannelTypeByRef todo doc
func (cth *ChannelTypeHandler) HandleGetChannelTypeByRef() rest.Handler {
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
		app, err := cth.GetChannelType(data.Query)
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

// HandleUpdateChannelType todo doc
func (cth *ChannelTypeHandler) HandleUpdateChannelType() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Couldn't process the request body")
			return
		}
		data := struct {
			ChannelType meta.ChannelType `json:"channel-type"`
			Ctx         string           `json:"ctx"`
		}{}
		json.Unmarshal(body, &data)

		err = cth.UpdateChannelType(&data.ChannelType, data.Ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}

// HandleDeleteChannelType todo doc
func (cth *ChannelTypeHandler) HandleDeleteChannelType() rest.Handler {
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
		err = cth.DeleteChannelType(data.Query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}
