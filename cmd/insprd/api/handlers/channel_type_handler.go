package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

// ChannelTypeHandler - contains handlers that uses the
// ChannelTypeMemory interface methods
type ChannelTypeHandler struct {
	memory.ChannelTypeMemory
}

// NewChannelTypeHandler - returns the handle function that
// manages the creation of a channel
func NewChannelTypeHandler(memManager memory.Manager) *ChannelTypeHandler {
	return &ChannelTypeHandler{
		ChannelTypeMemory: memManager.ChannelTypes(),
	}
}

// HandleCreateChannelType todo doc
func (cth *ChannelTypeHandler) HandleCreateChannelType() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeDI{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)
		if err != nil || !data.Setup {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		err = cth.CreateChannelType(&data.ChannelType, data.Ctx)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}

// HandleGetChannelTypeByRef todo doc
func (cth *ChannelTypeHandler) HandleGetChannelTypeByRef() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeQueryDI{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)
		if err != nil || !data.Setup {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		channelType, err := cth.GetChannelType(data.Ctx, data.CtName)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		rest.JSON(w, http.StatusOK, channelType)
	}
	return rest.Handler(handler)
}

// HandleUpdateChannelType todo doc
func (cth *ChannelTypeHandler) HandleUpdateChannelType() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeDI{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)
		if err != nil || !data.Setup {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		err = cth.UpdateChannelType(&data.ChannelType, data.Ctx)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}

// HandleDeleteChannelType todo doc
func (cth *ChannelTypeHandler) HandleDeleteChannelType() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeQueryDI{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)
		if err != nil || !data.Setup {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		err = cth.DeleteChannelType(data.Ctx, data.CtName)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return rest.Handler(handler)
}
