package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/rest"
)

// AliasHandler - contains handlers that uses the AliasMemory interface methods
type AliasHandler struct {
	*Handler
	logger *zap.Logger
}

// NewAliasHandler - returns the handle function that
// manages the creation of an Alias
func (handler *Handler) NewAliasHandler() *AliasHandler {
	return &AliasHandler{
		Handler: handler,
		logger:  logger.With(zap.String("section", "api"), zap.String("sub-section", "aliases")),
	}
}

// HandleCreate - returns the handle function that
// manages the creation of an Alias
func (ah *AliasHandler) HandleCreate() rest.Handler {
	l := ah.logger.With(zap.String("operation", "create"))
	l.Info("received alias create request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AliasDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			l.Error("unable to decode Alias create request data", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		l = l.With(
			zap.Any("alias", data.Alias),
			zap.String("targed", data.Target),
			zap.String("scope", scope),
			zap.Bool("dry-run", data.DryRun),
		)
		l.Debug("initiating Alias create transaction")
		ah.Memory.Tree().InitTransaction()

		err = ah.Memory.Tree().Alias().Create(scope, data.Target, &data.Alias)
		if err != nil {
			l.Error("unable to create Alias", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		changes, err := ah.Memory.Tree().GetTransactionChanges()
		if err != nil {
			l.Error("unable to get Alias create request changes", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		if !data.DryRun {
			l.Debug("applying Alias create changes in diff")
			err = ah.applyChangesInDiff(changes)
			if err != nil {
				l.Error("unable to apply Alias create changes in diff", zap.Error(err))
				rest.ERROR(w, err)
				ah.Memory.Tree().Cancel()
				return
			}

			l.Info("committing Alias create changes")
			defer ah.Memory.Tree().Commit()
		} else {
			l.Debug("cancelling Alias create changes")
			defer ah.Memory.Tree().Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

// HandleGet - return a handle function that obtains
// a Alias by the reference given
func (ah *AliasHandler) HandleGet() rest.Handler {
	l := ah.logger.With(zap.String("operation", "get"))
	l.Info("received alias retrieval request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AliasQueryDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			l.Error("unable to decode Alias get request data", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		l = l.With(
			zap.String("alias key", data.Key),
			zap.String("scope", scope),
			zap.Bool("dry-run", data.DryRun),
		)
		l.Debug("initiating Alias get transaction")
		ah.Memory.Tree().InitTransaction()

		app, err := ah.Memory.Tree().Perm().Alias().Get(scope, data.Key)
		if err != nil {
			l.Error("unable to get Alias", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		defer ah.Memory.Tree().Cancel()

		rest.JSON(w, http.StatusOK, app)
	}
	return rest.Handler(handler)
}

// HandleUpdate - returns a handle function that
// updates the Alias with the parameters given in the request
func (ah *AliasHandler) HandleUpdate() rest.Handler {
	l := ah.logger.With(zap.String("operation", "update"))
	l.Info("received alias update request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AliasDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			l.Error("unable to decode Alias update request data", zap.Error(err))
			rest.ERROR(w, err)
			return
		}
		l = l.With(
			zap.Any("alias", data.Alias),
			zap.String("targed", data.Target),
			zap.String("scope", scope),
			zap.Bool("dry-run", data.DryRun),
		)

		l.Debug("initiating Alias update transaction")
		ah.Memory.Tree().InitTransaction()

		err = ah.Memory.Tree().Alias().Update(scope, data.Target, &data.Alias)
		if err != nil {
			l.Error("unable to update Alias", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		changes, err := ah.Memory.Tree().GetTransactionChanges()
		if err != nil {
			l.Error("unable to get Alias update request changes", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		if !data.DryRun {
			l.Debug("applying Alias update changes in diff")
			err = ah.applyChangesInDiff(changes)
			if err != nil {
				l.Error("unable to apply Alias update changes in diff", zap.Error(err))
				rest.ERROR(w, err)
				ah.Memory.Tree().Cancel()
				return
			}

			l.Info("committing Alias update changes")
			defer ah.Memory.Tree().Commit()
		} else {
			l.Debug("cancelling Alias update changes")
			defer ah.Memory.Tree().Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

// HandleDelete - returns a handle function that
// deletes the Alias of the given path
func (ah *AliasHandler) HandleDelete() rest.Handler {
	l := ah.logger.With(zap.String("operation", "update"))
	l.Info("handling Alias delete request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AliasQueryDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			l.Error("unable to decode Alias delete request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}
		l = l.With(
			zap.String("alias key", data.Key),
			zap.String("scope", scope),
			zap.Bool("dry-run", data.DryRun),
		)
		l.Debug("initiating Alias delete transaction")
		ah.Memory.Tree().InitTransaction()

		err = ah.Memory.Tree().Alias().Delete(scope, data.Key)
		if err != nil {
			l.Error("unable to delete Alias", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		changes, err := ah.Memory.Tree().Alias().GetTransactionChanges()
		if err != nil {
			l.Error("unable to get Alias delete request changes", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		if !data.DryRun {
			l.Debug("applying Alias delete changes in diff")
			err = ah.applyChangesInDiff(changes)
			if err != nil {
				l.Error("unable to apply Alias delete changes in diff", zap.Error(err))
				rest.ERROR(w, err)
				ah.Memory.Tree().Cancel()
				return
			}

			l.Info("committing Alias delete changes")
			defer ah.Memory.Tree().Commit()
		} else {
			l.Debug("cancelling Alias delete changes")
			defer ah.Memory.Tree().Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}
