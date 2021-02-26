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

// ChannelTypeHandler - contains handlers that uses the
// ChannelTypeMemory interface methods
type ChannelTypeHandler struct {
	mem memory.Manager
	op  operators.OperatorInterface
}

// NewChannelTypeHandler - returns the handle function that
// manages the creation of a channelType
func NewChannelTypeHandler(memManager memory.Manager, op operators.OperatorInterface) *ChannelTypeHandler {
	return &ChannelTypeHandler{
		mem: memManager,
		op:  op,
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

		cth.mem.InitTransaction()
		if !data.DryRun {
			defer cth.mem.Commit()
		} else {
			defer cth.mem.Cancel()
		}
		err = cth.mem.ChannelTypes().CreateChannelType(data.Ctx, &data.ChannelType)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		diff, err := cth.mem.GetTransactionChanges()
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

		cth.mem.InitTransaction()
		defer cth.mem.Cancel()

		channelType, err := cth.mem.ChannelTypes().Get(data.Ctx, data.CtName)
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

		cth.mem.InitTransaction()
		if !data.DryRun {
			defer cth.mem.Commit()
		} else {
			defer cth.mem.Cancel()
		}

		err = cth.mem.ChannelTypes().UpdateChannelType(data.Ctx, &data.ChannelType)
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		diff, err := cth.mem.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		if !data.DryRun {
			var errs string
			ct, _ := cth.mem.ChannelTypes().Get(data.Ctx, data.ChannelType.Meta.Name)
			for _, chName := range ct.ConnectedChannels {
				ch, _ := cth.mem.Channels().Get(data.Ctx, chName)
				err = cth.op.Channels().Update(context.Background(), data.Ctx, ch)
				if err != nil {
					errs += err.Error() + "\n"
					continue
				}

				for _, appName := range ch.ConnectedApps {
					app, _ := cth.mem.Apps().Get(ch.Meta.Parent + "." + appName)
					_, err = cth.op.Nodes().UpdateNode(context.Background(), app)
					if err != nil {
						errs += err.Error() + "\n"
					}
				}
			}

			if errs != "" {
				rest.ERROR(w, ierrors.NewError().Message(errs).InternalServer().Build())
				return
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

		cth.mem.InitTransaction()
		if !data.DryRun {
			defer cth.mem.Commit()
		} else {
			defer cth.mem.Cancel()
		}
		err = cth.mem.ChannelTypes().DeleteChannelType(data.Ctx, data.CtName)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		diff, err := cth.mem.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}
