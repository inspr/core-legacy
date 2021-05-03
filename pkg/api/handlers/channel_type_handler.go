package handler

import (
	"encoding/json"
	"net/http"

	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/rest"
	"go.uber.org/zap"
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
	logger.Info("handling Channel Type create request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode Channel Type create request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating Channel Type create transaction")
		cth.Memory.InitTransaction()

		err = cth.Memory.ChannelTypes().Create(scope, &data.ChannelType)
		if err != nil {
			logger.Error("unable to create Channel Type",
				zap.String("ctype", data.ChannelType.Meta.Name),
				zap.String("context", scope),
				zap.Any("error", err))
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		diff, err := cth.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get Channel Type create request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Info("committing Channel Type create changes")
			defer cth.Memory.Commit()
		} else {
			logger.Info("canceling Channel Type create changes")
			defer cth.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}

// HandleGet - return a handle function that obtains
// a channelType by the reference given
func (cth *ChannelTypeHandler) HandleGet() rest.Handler {
	logger.Info("handling Channel Type get request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeQueryDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode Channel Type get request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating Channel Type get transaction")
		cth.Memory.InitTransaction()

		channelType, err := cth.Memory.Root().ChannelTypes().Get(scope, data.CtName)
		if err != nil {
			logger.Error("unable to get Channel Type",
				zap.String("ctype", data.CtName),
				zap.String("context", scope),
				zap.Any("error", err))
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
	logger.Info("handling Channel Type update request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode Channel Type update request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating Channel Type update transaction")
		cth.Memory.InitTransaction()

		err = cth.Memory.ChannelTypes().Update(scope, &data.ChannelType)
		if err != nil {
			logger.Error("unable to update Channel Type",
				zap.String("ctype", data.ChannelType.Meta.Name),
				zap.String("context", scope),
				zap.Any("error", err))
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		diff, err := cth.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get Channel Type update request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Debug("applying Channel Type update changes in diff")
			err = cth.applyChangesInDiff(diff)
			if err != nil {
				logger.Error("unable to apply Channel Type update changes in diff",
					zap.Any("error", err))
				rest.ERROR(w, err)
				cth.Memory.Cancel()
				return
			}

			logger.Info("committing Channel Type update changes")
			defer cth.Memory.Commit()
		} else {
			logger.Info("canceling Channel Type update changes")
			defer cth.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}

// HandleDelete - returns a handle function that
// deletes the channelType of the given path
func (cth *ChannelTypeHandler) HandleDelete() rest.Handler {
	logger.Info("handling Channel Type delete request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelTypeQueryDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode Channel Type delete request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating Channel Type delete transaction")
		cth.Memory.InitTransaction()

		err = cth.Memory.ChannelTypes().Delete(scope, data.CtName)
		if err != nil {
			logger.Error("unable to delete Channel Type",
				zap.String("ctype", data.CtName),
				zap.String("context", scope),
				zap.Any("error", err))
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		diff, err := cth.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get Channel Type delete request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Info("committing Channel Type delete changes")
			defer cth.Memory.Commit()
		} else {
			logger.Info("canceling Channel Type delete changes")
			defer cth.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}
