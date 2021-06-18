package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/rest"
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
func (th *TypeHandler) HandleCreate() rest.Handler {
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
		th.Memory.InitTransaction()

		err = th.Memory.Types().Create(scope, &data.Type)
		if err != nil {
			logger.Error("unable to create Type",
				zap.String("type", data.Type.Meta.Name),
				zap.String("scope", scope),
				zap.Any("error", err))
			rest.ERROR(w, err)
			th.Memory.Cancel()
			return
		}

		diff, err := th.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get Type create request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			th.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Info("committing Type create changes")
			defer th.Memory.Commit()
		} else {
			logger.Info("canceling Type create changes")
			defer th.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}

// HandleGet - return a handle function that obtains
// a Type by the reference given
func (th *TypeHandler) HandleGet() rest.Handler {
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
		th.Memory.InitTransaction()

		insprType, err := th.Memory.Root().Types().Get(scope, data.TypeName)
		if err != nil {
			logger.Error("unable to get Type",
				zap.String("type-name", data.TypeName),
				zap.String("scope", scope),
				zap.Any("error", err))
			rest.ERROR(w, err)
			th.Memory.Cancel()
			return
		}

		defer th.Memory.Cancel()

		rest.JSON(w, http.StatusOK, insprType)
	}
	return rest.Handler(handler)
}

// HandleUpdate - returns a handle function that
// updates the Type with the parameters given in the request
func (th *TypeHandler) HandleUpdate() rest.Handler {
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
		th.Memory.InitTransaction()

		err = th.Memory.Types().Update(scope, &data.Type)
		if err != nil {
			logger.Error("unable to update Type",
				zap.String("type", data.Type.Meta.Name),
				zap.String("scope", scope),
				zap.Any("error", err))
			rest.ERROR(w, err)
			th.Memory.Cancel()
			return
		}

		diff, err := th.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get Type update request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			th.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Debug("applying Type update changes in diff")
			err = th.applyChangesInDiff(diff)
			if err != nil {
				logger.Error("unable to apply Type update changes in diff",
					zap.Any("error", err))
				rest.ERROR(w, err)
				th.Memory.Cancel()
				return
			}

			logger.Info("committing Type update changes")
			defer th.Memory.Commit()
		} else {
			logger.Info("canceling Type update changes")
			defer th.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}

// HandleDelete - returns a handle function that
// deletes the Type of the given path
func (th *TypeHandler) HandleDelete() rest.Handler {
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
		th.Memory.InitTransaction()

		err = th.Memory.Types().Delete(scope, data.TypeName)
		if err != nil {
			logger.Error("unable to delete Type",
				zap.String("type", data.TypeName),
				zap.String("scope", scope),
				zap.Any("error", err))
			rest.ERROR(w, err)
			th.Memory.Cancel()
			return
		}

		diff, err := th.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get Type delete request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			th.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Info("committing Type delete changes")
			defer th.Memory.Commit()
		} else {
			logger.Info("canceling Type delete changes")
			defer th.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}
