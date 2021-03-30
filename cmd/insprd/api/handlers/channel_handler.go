package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

// ChannelHandler - contains handlers that uses the ChannelMemory interface methods
type ChannelHandler struct {
	*Handler
}

// NewChannelHandler exports
func (handler *Handler) NewChannelHandler() *ChannelHandler {
	return &ChannelHandler{
		handler,
	}
}

// HandleCreate - returns the handle function that
// manages the creation of a channel
func (ch *ChannelHandler) HandleCreate() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		ch.Memory.InitTransaction()

		err = ch.Memory.Channels().Create(data.Ctx, &data.Channel)
		if err != nil {
			rest.ERROR(w, err)
			ch.Memory.Cancel()
			return
		}

		changes, err := ch.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			ch.Memory.Cancel()
			return
		}

		if !data.DryRun {
			err = ch.applyChangesInDiff(changes)
			if err != nil {
				rest.ERROR(w, err)
				ch.Memory.Cancel()
				return
			}
		}

		if !data.DryRun {
			defer ch.Memory.Commit()
		} else {
			defer ch.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

// HandleGet - return a handle function that obtains
// a channel by the reference given
func (ch *ChannelHandler) HandleGet() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		ch.Memory.InitTransaction()

		channel, err := ch.Memory.Root().Channels().Get(data.Ctx, data.ChName)
		if err != nil {
			rest.ERROR(w, err)
			ch.Memory.Cancel()
			return
		}

		defer ch.Memory.Cancel()

		rest.JSON(w, http.StatusOK, channel)
	}
	return rest.Handler(handler)
}

// HandleUpdate - returns a handle function that
// updates the channel with the parameters given in the request
func (ch *ChannelHandler) HandleUpdate() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		ch.Memory.InitTransaction()

		err = ch.Memory.Channels().Update(data.Ctx, &data.Channel)
		if err != nil {
			rest.ERROR(w, err)
			ch.Memory.Cancel()
			return
		}
		changes, err := ch.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			ch.Memory.Cancel()
			return
		}

		if !data.DryRun {
			err = ch.applyChangesInDiff(changes)

			if err != nil {
				rest.ERROR(w, err)
				ch.Memory.Cancel()
				return
			}
		}

		if !data.DryRun {
			defer ch.Memory.Commit()
		} else {
			defer ch.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

// HandleDelete - returns a handle function that
// deletes the channel of the given path
func (ch *ChannelHandler) HandleDelete() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelQueryDI{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		ch.Memory.InitTransaction()

		err = ch.Memory.Channels().Delete(data.Ctx, data.ChName)
		if err != nil {
			rest.ERROR(w, err)
			ch.Memory.Cancel()
			return
		}

		changes, err := ch.Memory.Channels().GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			ch.Memory.Cancel()
			return
		}
		if !data.DryRun {
			err = ch.applyChangesInDiff(changes)

			if err != nil {
				rest.ERROR(w, err)
				ch.Memory.Cancel()
				return
			}
		}

		if !data.DryRun {
			defer ch.Memory.Commit()
		} else {
			defer ch.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}
