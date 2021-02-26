package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

// ChannelHandler - contains handlers that uses the ChannelMemory interface methods
type ChannelHandler struct {
	mem memory.Manager
	op  operators.OperatorInterface
}

// NewChannelHandler exports
func NewChannelHandler(memManager memory.Manager, op operators.OperatorInterface) *ChannelHandler {
	return &ChannelHandler{
		mem: memManager,
		op:  op,
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
		ch.mem.InitTransaction()
		if !data.DryRun {
			defer ch.mem.Commit()
		} else {
			defer ch.mem.Cancel()
		}

		err = ch.mem.Channels().CreateChannel(data.Ctx, &data.Channel)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		changes, err := ch.mem.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		channel, _ := ch.mem.Channels().Get(data.Ctx, data.Channel.Meta.Name)
		err = ch.op.Channels().Create(context.Background(), data.Ctx, channel)

		if err != nil {
			rest.ERROR(w, ierrors.NewError().InternalServer().InnerError(err).Message("unable to create channel in cluster").Build())
			return
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

		ch.mem.InitTransaction()
		defer ch.mem.Cancel()

		channel, err := ch.mem.Channels().Get(data.Ctx, data.ChName)
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
		ch.mem.InitTransaction()
		if !data.DryRun {
			defer ch.mem.Commit()
		} else {
			defer ch.mem.Cancel()
		}
		err = ch.mem.Channels().UpdateChannel(data.Ctx, &data.Channel)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		changes, err := ch.mem.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		channel, _ := ch.mem.Channels().Get(data.Ctx, data.Channel.Meta.Name)
		err = ch.op.Channels().Update(context.Background(), data.Ctx, channel)

		if err != nil {
			rest.ERROR(w, ierrors.NewError().InternalServer().InnerError(err).Message("unable to update channel in cluster").Build())
			return
		}

		var errs string
		for _, appName := range channel.ConnectedApps {
			app, _ := ch.mem.Apps().Get(data.Ctx + "." + appName)
			_, err = ch.op.Nodes().UpdateNode(context.Background(), app)
			if err != nil {
				errs += err.Error() + "\n"
			}
		}
		if errs != "" {
			rest.ERROR(w, ierrors.NewError().Message(errs).InternalServer().Build())
			return
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

		ch.mem.InitTransaction()
		if !data.DryRun {
			defer ch.mem.Commit()
		} else {
			defer ch.mem.Cancel()
		}

		err = ch.mem.Channels().DeleteChannel(data.Ctx, data.ChName)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		changes, err := ch.mem.Channels().GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		err = ch.op.Channels().Delete(context.Background(), data.Ctx, data.ChName)

		if err != nil {
			rest.ERROR(w, ierrors.NewError().InternalServer().InnerError(err).Message("unable to delete channel from cluster").Build())
			return
		}
		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}
