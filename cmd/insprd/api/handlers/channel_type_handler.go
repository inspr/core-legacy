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

// HandleCreate - returns the handle function that
// manages the creation of a channelType
func (cth *ChannelTypeHandler) HandleCreate() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		cth.Memory.InitTransaction()

		err = cth.Memory.ChannelTypes().Create(data.Ctx, &data.ChannelType)
		if err != nil {
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}
		diff, err := cth.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		if !data.DryRun {
			defer cth.Memory.Commit()
		} else {
			defer cth.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}

// HandleGet - return a handle function that obtains
// a channelType by the reference given
func (cth *ChannelTypeHandler) HandleGet() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		cth.Memory.InitTransaction()

		channelType, err := cth.Memory.Root().ChannelTypes().Get(data.Ctx, data.CtName)
		if err != nil {
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		defer cth.Memory.Cancel()

		rest.JSON(w, http.StatusOK, channelType)
	}
	return rest.Handler(handler)
}

// HandleUpdate - returns a handle function that
// updates the channelType with the parameters given in the request
func (cth *ChannelTypeHandler) HandleUpdate() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		cth.Memory.InitTransaction()

		err = cth.Memory.ChannelTypes().Update(data.Ctx, &data.ChannelType)
		if err != nil {
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		diff, err := cth.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		if !data.DryRun {
			err = cth.applyChangesInDiff(diff)
			if err != nil {
				rest.ERROR(w, err)
				cth.Memory.Cancel()
				return
			}
		}

		if !data.DryRun {
			defer cth.Memory.Commit()
		} else {
			defer cth.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}

// HandleDelete - returns a handle function that
// deletes the channelType of the given path
func (cth *ChannelTypeHandler) HandleDelete() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		cth.Memory.InitTransaction()

		err = cth.Memory.ChannelTypes().Delete(data.Ctx, data.CtName)
		if err != nil {
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}
		diff, err := cth.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		if !data.DryRun {
			defer cth.Memory.Commit()
		} else {
			defer cth.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}
