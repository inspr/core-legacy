package handler

import (
	"encoding/json"
	"net/http"

	"github.com/inspr/inspr/cmd/insprd/api/models"
	"github.com/inspr/inspr/pkg/rest"
	"go.uber.org/zap"
)

// AliasHandler - contains handlers that uses the AliasMemory interface methods
type AliasHandler struct {
	*Handler
}

// NewAliasHandler - returns the handle function that
// manages the creation of an Alias
func (handler *Handler) NewAliasHandler() *AliasHandler {
	return &AliasHandler{
		handler,
	}
}

// HandleCreate - returns the handle function that
// manages the creation of an Alias
func (ah *AliasHandler) HandleCreate() rest.Handler {
	logger.Info("handling Alias create request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AliasDI{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode Alias create request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		logger.Debug("initiating Alias create transaction")
		ah.Memory.InitTransaction()

		err = ah.Memory.Alias().Create(data.Ctx, data.Target, &data.Alias)
		if err != nil {
			logger.Error("unable to create Alias",
				zap.Any("alias", data.Alias),
				zap.String("targed", data.Target),
				zap.String("context", data.Ctx),
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		changes, err := ah.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get Alias create request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Debug("applying Alias create changes in diff")
			err = ah.applyChangesInDiff(changes)
			if err != nil {
				logger.Error("unable to apply Alias create changes in diff",
					zap.Any("error", err))
				rest.ERROR(w, err)
				ah.Memory.Cancel()
				return
			}

			logger.Info("committing Alias create changes")
			defer ah.Memory.Commit()
		} else {
			logger.Info("cancelling Alias create changes")
			defer ah.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

// HandleGet - return a handle function that obtains
// a Alias by the reference given
func (ah *AliasHandler) HandleGet() rest.Handler {
	logger.Info("handling Alias get request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AliasQueryDI{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode Alias get request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		logger.Debug("initiating Alias get transaction")
		ah.Memory.InitTransaction()

		app, err := ah.Memory.Root().Alias().Get(data.Ctx, data.Key)
		if err != nil {
			logger.Error("unable to get Alias",
				zap.String("alias key", data.Key),
				zap.String("context", data.Ctx),
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		defer ah.Memory.Cancel()

		rest.JSON(w, http.StatusOK, app)
	}
	return rest.Handler(handler)
}

// HandleUpdate - returns a handle function that
// updates the Alias with the parameters given in the request
func (ah *AliasHandler) HandleUpdate() rest.Handler {
	logger.Info("handling Alias update request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AliasDI{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode Alias update request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating Alias update transaction")
		ah.Memory.InitTransaction()

		err = ah.Memory.Alias().Update(data.Ctx, data.Target, &data.Alias)
		if err != nil {
			logger.Error("unable to update Alias",
				zap.Any("alias", data.Alias),
				zap.String("targed", data.Target),
				zap.String("context", data.Ctx),
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		changes, err := ah.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get Alias update request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Debug("applying Alias update changes in diff")
			err = ah.applyChangesInDiff(changes)
			if err != nil {
				logger.Error("unable to apply Alias update changes in diff",
					zap.Any("error", err))
				rest.ERROR(w, err)
				ah.Memory.Cancel()
				return
			}

			logger.Info("committing Alias update changes")
			defer ah.Memory.Commit()
		} else {
			logger.Info("cancelling Alias update changes")
			defer ah.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

// HandleDelete - returns a handle function that
// deletes the Alias of the given path
func (ah *AliasHandler) HandleDelete() rest.Handler {
	logger.Info("handling Alias delete request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AliasQueryDI{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode Alias delete request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating Alias delete transaction")
		ah.Memory.InitTransaction()

		err = ah.Memory.Alias().Delete(data.Ctx, data.Key)
		if err != nil {
			logger.Error("unable to delete Alias",
				zap.String("alias key", data.Key),
				zap.String("context", data.Ctx),
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		changes, err := ah.Memory.Alias().GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get Alias delete request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Debug("applying Alias delete changes in diff")
			err = ah.applyChangesInDiff(changes)
			if err != nil {
				logger.Error("unable to apply Alias delete changes in diff",
					zap.Any("error", err))
				rest.ERROR(w, err)
				ah.Memory.Cancel()
				return
			}

			logger.Info("committing Alias create changes")
			defer ah.Memory.Commit()
		} else {
			logger.Info("cancelling Alias create changes")
			defer ah.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}
