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
// manages the creation of a channelType
func NewChannelTypeHandler(memManager memory.Manager) *ChannelTypeHandler {
	return &ChannelTypeHandler{
		ChannelTypeMemory: memManager.ChannelTypes(),
	}
}

// HandleCreateChannelType - returns the handle function that
// manages the creation of a channelType
func (cth *ChannelTypeHandler) HandleCreateChannelType() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		cth.InitTransaction()
		if !data.DryRun {
			defer cth.Commit()
		} else {
			defer cth.Cancel()
		}
		err = cth.CreateChannelType(&data.ChannelType, data.Ctx)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		diff, err := cth.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}

// HandleGetChannelTypeByRef - return a handle function that obtains
// a channelType by the reference given
func (cth *ChannelTypeHandler) HandleGetChannelTypeByRef() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
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

// HandleUpdateChannelType - returns a handle function that
// updates the channelType with the parameters given in the request
func (cth *ChannelTypeHandler) HandleUpdateChannelType() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		cth.InitTransaction()
		if !data.DryRun {
			defer cth.Commit()
		} else {
			defer cth.Cancel()
		}
		err = cth.UpdateChannelType(&data.ChannelType, data.Ctx)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		diff, err := cth.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}

// HandleDeleteChannelType - returns a handle function that
// deletes the channelType of the given path
func (cth *ChannelTypeHandler) HandleDeleteChannelType() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, http.StatusBadRequest, err)
			return
		}

		cth.InitTransaction()
		if !data.DryRun {
			defer cth.Commit()
		} else {
			defer cth.Cancel()
		}
		err = cth.DeleteChannelType(data.Ctx, data.CtName)
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		diff, err := cth.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}
