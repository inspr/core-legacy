package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/rest"
)

// ChannelHandler - contains handlers that uses the ChannelMemory interface methods
type ChannelHandler struct {
	*Handler
}

// NewChannelHandler - returns the handle function that
// manages the creation of a Channel
func (handler *Handler) NewChannelHandler() *ChannelHandler {
	return &ChannelHandler{
		handler,
	}
}

// HandleCreate - returns the handle function that
// manages the creation of a Channel
func (ch *ChannelHandler) HandleCreate() rest.Handler {
	logger.Info("handling Channel create request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode Channel create request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating Channel create transaction")
		ch.Memory.InitTransaction()

		err = ch.Memory.Channels().Create(scope, &data.Channel)
		if err != nil {
			logger.Error("unable to create Channel",
				zap.String("channel", data.Channel.Meta.Name),
				zap.String("scope", scope),
				zap.Any("error", err))
			rest.ERROR(w, err)
			ch.Memory.Cancel()
			return
		}

		changes, err := ch.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get Channel create request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			ch.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Debug("applying Channel create changes in diff")
			err = ch.applyChangesInDiff(changes)
			if err != nil {
				logger.Error("unable to apply Channel create changes in diff",
					zap.Any("error", err))
				rest.ERROR(w, err)
				ch.Memory.Cancel()
				return
			}

			logger.Info("committing Channel create changes")
			defer ch.Memory.Commit()
		} else {
			logger.Info("cancelling Channel create changes")
			defer ch.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

// HandleGet - return a handle function that obtains
// a Channel by the reference given
func (ch *ChannelHandler) HandleGet() rest.Handler {
	logger.Info("handling Channel get request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelQueryDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode Channel get request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating Channel get transaction")
		ch.Memory.InitTransaction()

		channel, err := ch.Memory.Tree().Channels().Get(scope, data.ChName)
		if err != nil {
			logger.Error("unable to get Channel",
				zap.String("channel", data.ChName),
				zap.String("scope", scope),
				zap.Any("error", err))
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
// updates the Channel with the parameters given in the request
func (ch *ChannelHandler) HandleUpdate() rest.Handler {
	logger.Info("handling Channel update request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode Channel update request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating Channel update transaction")
		ch.Memory.InitTransaction()

		err = ch.Memory.Channels().Update(scope, &data.Channel)
		if err != nil {
			logger.Error("unable to update Channel",
				zap.String("channel", data.Channel.Meta.Name),
				zap.String("scope", scope),
				zap.Any("error", err))
			rest.ERROR(w, err)
			ch.Memory.Cancel()
			return
		}

		changes, err := ch.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get Channel update request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			ch.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Debug("applying Channel update changes in diff")
			err = ch.applyChangesInDiff(changes)
			if err != nil {
				logger.Error("unable to apply Channel update changes in diff",
					zap.Any("error", err))
				rest.ERROR(w, err)
				ch.Memory.Cancel()
				return
			}

			logger.Info("committing Channel update changes")
			defer ch.Memory.Commit()
		} else {
			logger.Info("cancelling Channel update changes")
			defer ch.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

// HandleDelete - returns a handle function that
// deletes the Channel of the given path
func (ch *ChannelHandler) HandleDelete() rest.Handler {
	logger.Info("handling Channel delete request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.ChannelQueryDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode Channel delete request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating Channel delete transaction")
		ch.Memory.InitTransaction()

		err = ch.Memory.Channels().Delete(scope, data.ChName)
		if err != nil {
			logger.Error("unable to delete Channel",
				zap.String("channel", data.ChName),
				zap.String("scope", scope),
				zap.Any("error", err))
			rest.ERROR(w, err)
			ch.Memory.Cancel()
			return
		}

		changes, err := ch.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get Channel delete request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			ch.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Debug("applying Channel delete changes in diff")
			err = ch.applyChangesInDiff(changes)
			if err != nil {
				logger.Error("unable to apply Channel delete changes in diff",
					zap.Any("error", err))
				rest.ERROR(w, err)
				ch.Memory.Cancel()
				return
			}

			logger.Info("committing Channel delete changes")
			defer ch.Memory.Commit()
		} else {
			logger.Info("cancelling Channel delete changes")
			defer ch.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}
