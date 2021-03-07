package handler

import (
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

// ChannelHandler - contains handlers that uses the ChannelMemory interface methods
type ChannelHandler struct {
	Handler
}

// NewChannelHandler exports
func NewChannelHandler(handler Handler) *ChannelHandler {
	return &ChannelHandler{
		handler,
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
			rest.ERROR(w, err)
			return
		}
		ch.Memory.InitTransaction()
		if !data.DryRun {
			defer ch.Memory.Commit()
		} else {
			defer ch.Memory.Cancel()
		}

		err = ch.Memory.Channels().CreateChannel(data.Ctx, &data.Channel)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		changes, err := ch.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		if !data.DryRun {
			err = ch.applyChangesInDiff(changes)
			if err != nil {
				rest.ERROR(w, err)
				return
			}
		}

		rest.JSON(w, http.StatusOK, changes)
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
			rest.ERROR(w, err)
			return
		}

		ch.Memory.InitTransaction()
		defer ch.Memory.Cancel()

		channel, err := ch.Memory.Root().Channels().Get(data.Ctx, data.ChName)
		if err != nil {
			rest.ERROR(w, err)
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
			rest.ERROR(w, err)
			return
		}
		ch.Memory.InitTransaction()
		if !data.DryRun {
			defer ch.Memory.Commit()
		} else {
			defer ch.Memory.Cancel()
		}
		err = ch.Memory.Channels().UpdateChannel(data.Ctx, &data.Channel)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		changes, err := ch.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		if !data.DryRun {
			err = ch.applyChangesInDiff(changes)

			if err != nil {
				rest.ERROR(w, err)
				return
			}
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

// HandleDeleteChannel - returns a handle function that
// deletes the channel of the given path
func (ch *ChannelHandler) HandleDeleteChannel() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelQueryDI{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, err)
			return
		}

		ch.Memory.InitTransaction()
		if !data.DryRun {
			defer ch.Memory.Commit()
		} else {
			defer ch.Memory.Cancel()
		}

		err = ch.Memory.Channels().DeleteChannel(data.Ctx, data.ChName)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		changes, err := ch.Memory.Channels().GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		if !data.DryRun {
			err = ch.applyChangesInDiff(changes)

			if err != nil {
				rest.ERROR(w, err)
				return
			}
		}
		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}
