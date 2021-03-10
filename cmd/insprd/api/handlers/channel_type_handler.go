package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

// ChannelTypeHandler - contains handlers that uses the
// ChannelTypeMemory interface methods
type ChannelTypeHandler struct {
	*Handler
}

// NewChannelTypeHandler - returns the handle function that
// manages the creation of a channelType
func (handler *Handler) NewChannelTypeHandler() *ChannelTypeHandler {
	return &ChannelTypeHandler{
		handler,
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
			rest.ERROR(w, err)
			return
		}

		cth.Memory.InitTransaction()
		if !data.DryRun {
			defer cth.Memory.Commit()
		} else {
			defer cth.Memory.Cancel()
		}
		err = cth.Memory.ChannelTypes().CreateChannelType(data.Ctx, &data.ChannelType)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		diff, err := cth.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
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
			rest.ERROR(w, err)
			return
		}

		cth.Memory.InitTransaction()
		defer cth.Memory.Cancel()

		channelType, err := cth.Memory.Root().ChannelTypes().Get(data.Ctx, data.CtName)
		if err != nil {
			rest.ERROR(w, err)
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
			rest.ERROR(w, err)
			return
		}

		cth.Memory.InitTransaction()
		if !data.DryRun {
			defer cth.Memory.Commit()
		} else {
			defer cth.Memory.Cancel()
		}

		err = cth.Memory.ChannelTypes().UpdateChannelType(data.Ctx, &data.ChannelType)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		diff, err := cth.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		if !data.DryRun {
			err = cth.applyChangesInDiff(diff)
			if err != nil {
				rest.ERROR(w, err)
			}
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
			rest.ERROR(w, err)
			return
		}

		cth.Memory.InitTransaction()
		if !data.DryRun {
			defer cth.Memory.Commit()
		} else {
			defer cth.Memory.Cancel()
		}
		err = cth.Memory.ChannelTypes().DeleteChannelType(data.Ctx, data.CtName)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		diff, err := cth.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}
