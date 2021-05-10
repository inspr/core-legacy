package handler

import (
	"encoding/json"
	"net/http"

	"github.com/inspr/inspr/pkg/api/models"
	"github.com/inspr/inspr/pkg/rest"
	"go.uber.org/zap"
)

// TypeHandler - contains handlers that uses the
// TypeMemory interface methods
type TypeHandler struct {
	*Handler
}

// NewTypeHandler - returns the handle function that
// manages the creation of a Type
func (handler *Handler) NewTypeHandler() *TypeHandler {
	return &TypeHandler{
		handler,
	}
}

// HandleCreate - returns the handle function that
// manages the creation of a Type
func (cth *TypeHandler) HandleCreate() rest.Handler {
	logger.Info("handling Type create request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.TypeDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode Type create request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating Type create transaction")
		cth.Memory.InitTransaction()

		err = cth.Memory.Types().Create(scope, &data.Type)
		if err != nil {
			logger.Error("unable to create Channel Type",
				zap.String("ctype", data.Type.Meta.Name),
				zap.String("context", scope),
				zap.Any("error", err))
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		diff, err := cth.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get Type create request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Info("committing Type create changes")
			defer cth.Memory.Commit()
		} else {
			logger.Info("canceling Type create changes")
			defer cth.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}

// HandleGet - return a handle function that obtains
// a Type by the reference given
func (cth *TypeHandler) HandleGet() rest.Handler {
	logger.Info("handling Type get request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.TypeQueryDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode Type get request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating Type get transaction")
		cth.Memory.InitTransaction()

		Type, err := cth.Memory.Root().Types().Get(scope, data.CtName)
		if err != nil {
			logger.Error("unable to get Type",
				zap.String("ctype", data.CtName),
				zap.String("context", scope),
				zap.Any("error", err))
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		defer cth.Memory.Cancel()

		rest.JSON(w, http.StatusOK, Type)
	}
	return rest.Handler(handler)
}

// HandleUpdate - returns a handle function that
// updates the Type with the parameters given in the request
func (cth *TypeHandler) HandleUpdate() rest.Handler {
	logger.Info("handling Type update request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.TypeDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode Type update request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating Type update transaction")
		cth.Memory.InitTransaction()

		err = cth.Memory.Types().Update(scope, &data.Type)
		if err != nil {
			logger.Error("unable to update Channel Type",
				zap.String("ctype", data.Type.Meta.Name),
				zap.String("context", scope),
				zap.Any("error", err))
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		diff, err := cth.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get Type update request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Debug("applying Type update changes in diff")
			err = cth.applyChangesInDiff(diff)
			if err != nil {
				logger.Error("unable to apply Type update changes in diff",
					zap.Any("error", err))
				rest.ERROR(w, err)
				cth.Memory.Cancel()
				return
			}

			logger.Info("committing Type update changes")
			defer cth.Memory.Commit()
		} else {
			logger.Info("canceling Type update changes")
			defer cth.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}

// HandleDelete - returns a handle function that
// deletes the Type of the given path
func (cth *TypeHandler) HandleDelete() rest.Handler {
	logger.Info("handling Type delete request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.TypeQueryDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode Type delete request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating Type delete transaction")
		cth.Memory.InitTransaction()

		err = cth.Memory.Types().Delete(scope, data.CtName)
		if err != nil {
			logger.Error("unable to delete Type",
				zap.String("ctype", data.CtName),
				zap.String("context", scope),
				zap.Any("error", err))
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		diff, err := cth.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get Type delete request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			cth.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Info("committing Type delete changes")
			defer cth.Memory.Commit()
		} else {
			logger.Info("canceling Type delete changes")
			defer cth.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}
