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
	logger *zap.Logger
}

// NewTypeHandler - returns the handle function that
// manages the creation of a Type
func (handler *Handler) NewTypeHandler() *TypeHandler {
	return &TypeHandler{
		Handler: handler,
		logger:  logger.With(zap.String("sub-section", "type")),
	}
}

// HandleCreate - returns the handle function that
// manages the creation of a Type
func (th *TypeHandler) HandleCreate() rest.Handler {
	l := th.logger.With(zap.String("operation", "create"))
	l.Info("received type get request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.TypeDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			l.Error("unable to decode Type create request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}
		l = l.With(
			zap.String("type", data.Type.Meta.Name),
			zap.String("scope", scope),
			zap.Bool("dry-run", data.DryRun),
		)

		l.Debug("initiating Type create transaction")
		th.Memory.Tree().InitTransaction()

		err = th.Memory.Tree().Types().Create(scope, &data.Type)
		if err != nil {
			l.Error("unable to create Type", zap.Error(err))
			rest.ERROR(w, err)
			th.Memory.Tree().Cancel()
			return
		}

		diff, err := th.Memory.Tree().GetTransactionChanges()
		if err != nil {
			l.Error("unable to get Type create request changes", zap.Error(err))
			rest.ERROR(w, err)
			th.Memory.Tree().Cancel()
			return
		}

		if !data.DryRun {
			l.Info("committing Type create changes")
			defer th.Memory.Tree().Commit()
		} else {
			l.Debug("canceling Type create changes")
			defer th.Memory.Tree().Cancel()
		}

		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}

// HandleGet - return a handle function that obtains
// a Type by the reference given
func (th *TypeHandler) HandleGet() rest.Handler {
	l := th.logger.With(zap.String("operation", "get"))
	l.Info("handling Type get request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.TypeQueryDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			l.Error("unable to decode Type get request data", zap.Error(err))
			rest.ERROR(w, err)
			return
		}
		l = l.With(
			zap.String("type-name", data.TypeName),
			zap.String("scope", scope),
			zap.Bool("dry-run", data.DryRun),
		)

		l.Debug("initiating Type get transaction")
		th.Memory.Tree().InitTransaction()

		insprType, err := th.Memory.Tree().Perm().Types().Get(scope, data.TypeName)
		if err != nil {
			l.Error("unable to get Type", zap.Error(err))
			rest.ERROR(w, err)
			th.Memory.Tree().Cancel()
			return
		}

		defer th.Memory.Tree().Cancel()

		rest.JSON(w, http.StatusOK, insprType)
	}
	return rest.Handler(handler)
}

// HandleUpdate - returns a handle function that
// updates the Type with the parameters given in the request
func (th *TypeHandler) HandleUpdate() rest.Handler {
	l := th.logger.With(zap.String("operation", "update"))
	l.Info("handling Type update request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.TypeDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			l.Error("unable to decode Type update request data", zap.Error(err))
			rest.ERROR(w, err)
			return
		}
		l = l.With(
			zap.String("type", data.Type.Meta.Name),
			zap.String("scope", scope),
			zap.Bool("dry-run", data.DryRun),
		)

		l.Debug("initiating Type update transaction")
		th.Memory.Tree().InitTransaction()

		err = th.Memory.Tree().Types().Update(scope, &data.Type)
		if err != nil {
			l.Error("unable to update Type", zap.Error(err))
			rest.ERROR(w, err)
			th.Memory.Tree().Cancel()
			return
		}

		diff, err := th.Memory.Tree().GetTransactionChanges()
		if err != nil {
			l.Error("unable to get Type update request changes", zap.Error(err))
			rest.ERROR(w, err)
			th.Memory.Tree().Cancel()
			return
		}

		if !data.DryRun {
			l.Debug("applying Type update changes in diff")
			err = th.applyChangesInDiff(diff)
			if err != nil {
				l.Error("unable to apply Type update changes in diff", zap.Error(err))
				rest.ERROR(w, err)
				th.Memory.Tree().Cancel()
				return
			}

			l.Info("committing Type update changes")
			defer th.Memory.Tree().Commit()
		} else {
			l.Debug("canceling Type update changes")
			defer th.Memory.Tree().Cancel()
		}

		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}

// HandleDelete - returns a handle function that
// deletes the Type of the given path
func (th *TypeHandler) HandleDelete() rest.Handler {
	l := th.logger.With(zap.String("operation", "delete"))
	l.Info("handling Type delete request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.TypeQueryDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			l.Error("unable to decode Type delete request data", zap.Error(err))
			rest.ERROR(w, err)
			return
		}
		l = l.With(
			zap.String("type", data.TypeName),
			zap.String("scope", scope),
			zap.Bool("dry-run", data.DryRun),
		)

		l.Debug("initiating Type delete transaction")
		th.Memory.Tree().InitTransaction()

		err = th.Memory.Tree().Types().Delete(scope, data.TypeName)
		if err != nil {
			l.Error("unable to delete Type", zap.Error(err))
			rest.ERROR(w, err)
			th.Memory.Tree().Cancel()
			return
		}

		diff, err := th.Memory.Tree().GetTransactionChanges()
		if err != nil {
			l.Error("unable to get Type delete request changes", zap.Error(err))
			rest.ERROR(w, err)
			th.Memory.Tree().Cancel()
			return
		}

		if !data.DryRun {
			l.Info("committing Type delete changes")
			defer th.Memory.Tree().Commit()
		} else {
			l.Debug("canceling Type delete changes")
			defer th.Memory.Tree().Cancel()
		}

		rest.JSON(w, http.StatusOK, diff)
	}
	return rest.Handler(handler)
}
