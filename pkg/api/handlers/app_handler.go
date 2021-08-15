package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/auth"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/utils"
	"inspr.dev/inspr/pkg/rest"
)

// AppHandler - contains handlers that uses the AppMemory interface methods
type AppHandler struct {
	*Handler
	logger *zap.Logger
}

// NewAppHandler - returns the handle function that
// manages the creation of a dApp
func (handler *Handler) NewAppHandler() *AppHandler {
	return &AppHandler{
		Handler: handler,
		logger: logger.With(
			zap.String("section", "api"),
			zap.String("subsection", "dapps"),
		),
	}
}

// HandleCreate - returns the handle function that
// manages the creation of a dApp
func (ah *AppHandler) HandleCreate() rest.Handler {
	l := ah.logger.With(zap.String("opeartion", "create"))
	l.Info("received create dapp request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			l.Error("unable to decode dApp create request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}
		l = l.With(
			zap.String("dapp", data.App.Meta.Name),
			zap.String("scope", scope),
			zap.Bool("dry-run", data.DryRun),
		)

		l.Debug("initiating dApp create transaction")
		ah.Memory.Tree().InitTransaction()

		brokers, err := ah.Memory.Brokers().Get()
		if err != nil {
			l.Error("unable to get broker data", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		err = ah.Memory.Tree().Apps().Create(scope, &data.App, brokers)
		if err != nil {
			l.Error("unable to create Channel", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		changes, err := ah.Memory.Tree().GetTransactionChanges()
		if err != nil {
			l.Error("unable to get dApp create request changes", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		if !data.DryRun {
			l.Debug("applying changes to the cluster")
			err = ah.applyChangesInDiff(changes)
			if err != nil {
				l.Error(
					"unable to apply dApp create changes in diff",
					zap.Error(err),
				)
				rest.ERROR(w, err)
				ah.Memory.Tree().Cancel()
				return
			}

			l.Info("changes completed - committing...")
			defer ah.Memory.Tree().Commit()
		} else {
			l.Debug("cancelling dApp create changes")
			defer ah.Memory.Tree().Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}

	return rest.Handler(handler)
}

// HandleGet - return a handle function that obtains
// a dApp by the reference query given
func (ah *AppHandler) HandleGet() rest.Handler {
	l := ah.logger.With(zap.String("opeartion", "get"))
	l.Info("received dapp get request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppQueryDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			l.Error("unable to decode dApp get request data", zap.Error(err))
			rest.ERROR(w, err)
			return
		}
		l = l.With(zap.String("scope", scope), zap.Bool("dry-run", data.DryRun))

		l.Debug("initiating dApp get transaction")
		ah.Memory.Tree().InitTransaction()

		app, err := ah.Memory.Tree().Apps().Get(scope)
		if err != nil {
			l.Error("unable to get dApp", zap.Error(err))
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
// updates the dApp with the parameters given in the request
func (ah *AppHandler) HandleUpdate() rest.Handler {
	l := ah.logger.With(zap.String("opeartion", "update"))
	l.Info("received dapp update request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			l.Error("unable to decode dApp update request data", zap.Error(err))
			rest.ERROR(w, err)
			return
		}

		l = l.With(
			zap.String("dapp", data.App.Meta.Name),
			zap.String("scope", scope),
			zap.Bool("dry-run", data.DryRun),
		)
		l.Debug("initiating dApp update transaction")
		ah.Memory.Tree().InitTransaction()

		brokers, err := ah.Memory.Brokers().Get()
		if err != nil {
			l.Error("unable to get broker data", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		err = ah.Memory.Tree().Apps().Update(scope, &data.App, brokers)
		if err != nil {
			l.Error("unable to update dApp", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		changes, err := ah.Memory.Tree().GetTransactionChanges()
		if err != nil {
			l.Error("unable to get dApp update request changes", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		if !data.DryRun {
			l.Debug("applying dApp update changes in diff")
			err = ah.applyChangesInDiff(changes)
			if err != nil {
				l.Error(
					"unable to apply dApp update changes in diff",
					zap.Error(err),
				)
				rest.ERROR(w, err)
				ah.Memory.Tree().Cancel()
				return
			}

			l.Info("committing dApp update changes")
			defer ah.Memory.Tree().Commit()
		} else {
			l.Debug("cancelling dApp update changes")
			defer ah.Memory.Tree().Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

// HandleDelete - returns a handle function that
// deletes the dApp of the given path
func (ah *AppHandler) HandleDelete() rest.Handler {
	l := ah.logger.With(zap.String("opeartion", "delete"))
	l.Info("received dapp delete request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppQueryDI{}
		scope := r.Header.Get(rest.HeaderScopeKey)

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			l.Error("unable to decode dApp delete request data", zap.Error(err))
			rest.ERROR(w, err)
			return
		}
		l = l.With(zap.String("scope", scope), zap.Bool("dry-run", data.DryRun))

		l.Info("initiating dApp delete transaction")
		ah.Memory.Tree().InitTransaction()

		err = ah.Memory.Tree().Apps().Delete(scope)
		if err != nil {
			l.Error("unable to delete dApp", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		changes, err := ah.Memory.Tree().GetTransactionChanges()
		if err != nil {
			l.Error("unable to get dApp delete request changes", zap.Error(err))
			rest.ERROR(w, err)
			ah.Memory.Tree().Cancel()
			return
		}

		if !data.DryRun {
			l.Debug("applying dApp delete changes in diff")
			err = ah.applyChangesInDiff(changes)
			if err != nil {
				l.Error(
					"unable to apply dApp delete changes in diff",
					zap.Error(err),
				)
				rest.ERROR(w, err)
				ah.Memory.Tree().Cancel()
				return
			}

			l.Info("committing dApp delete changes")
			defer ah.Memory.Tree().Commit()
		} else {
			l.Debug("cancelling dApp delete changes")
			defer ah.Memory.Tree().Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

func (handler *Handler) deleteApp(app *meta.App) error {
	errs := ierrors.MultiError{
		Errors: []error{},
	}
	if app.Spec.Node.Spec.Image != "" {
		err := handler.Operator.Nodes().
			DeleteNode(context.Background(), app.Meta.Parent, app.Meta.Name)
		if err != nil {
			errs.Add(err)
		}

	} else {
		for _, subApp := range app.Spec.Apps {
			errs.Add(handler.deleteApp(subApp))
		}
		for c := range app.Spec.Channels {
			scope, _ := utils.JoinScopes(app.Meta.Parent, app.Meta.Name)

			err := handler.Operator.Channels().Delete(context.Background(), scope, c)
			if err != nil {
				errs.Add(err)
			}
		}
	}
	if len(errs.Errors) > 0 {
		return &errs
	}
	return nil
}

// GetAuth returns the handler's Auth interface
func (handler *Handler) GetAuth() auth.Auth {
	return handler.Auth
}
