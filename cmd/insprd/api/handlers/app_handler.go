package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/inspr/inspr/cmd/insprd/api/models"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/utils"
	"github.com/inspr/inspr/pkg/rest"
	"go.uber.org/zap"
)

// AppHandler - contains handlers that uses the AppMemory interface methods
type AppHandler struct {
	*Handler
}

// NewAppHandler - returns the handle function that
// manages the creation of a dApp
func (handler *Handler) NewAppHandler() *AppHandler {
	return &AppHandler{
		handler,
	}
}

// HandleCreate - returns the handle function that
// manages the creation of a dApp
func (ah *AppHandler) HandleCreate() rest.Handler {
	logger.Info("handling dApp create request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppDI{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode dApp create request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		logger.Debug("initiating dApp create transaction")
		ah.Memory.InitTransaction()

		err = ah.Memory.Apps().Create(data.Ctx, &data.App)
		if err != nil {
			logger.Error("unable to create Channel",
				zap.String("dApp", data.App.Meta.Name),
				zap.String("context", data.Ctx),
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		changes, err := ah.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get dApp create request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Debug("applying dApp create changes in diff")
			err = ah.applyChangesInDiff(changes)
			if err != nil {
				logger.Error("unable to apply dApp create changes in diff",
					zap.Any("error", err))
				rest.ERROR(w, err)
				ah.Memory.Cancel()
				return
			}

			logger.Info("committing dApp create changes")
			defer ah.Memory.Commit()
		} else {
			logger.Info("cancelling dApp create changes")
			defer ah.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}

	return rest.Handler(handler)
}

// HandleGet - return a handle function that obtains
// a dApp by the reference query given
func (ah *AppHandler) HandleGet() rest.Handler {
	logger.Info("handling dApp get request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppQueryDI{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode dApp get request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating dApp get transaction")
		ah.Memory.InitTransaction()

		app, err := ah.Memory.Root().Apps().Get(data.Ctx)
		if err != nil {
			logger.Error("unable to get dApp",
				zap.String("dApp query", data.Ctx),
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
// updates the dApp with the parameters given in the request
func (ah *AppHandler) HandleUpdate() rest.Handler {
	logger.Info("handling dApp update request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppDI{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode dApp update request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating dApp update transaction")
		ah.Memory.InitTransaction()

		err = ah.Memory.Apps().Update(data.Ctx, &data.App)
		if err != nil {
			logger.Error("unable to update dApp",
				zap.String("dApp", data.App.Meta.Name),
				zap.String("context", data.Ctx),
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		changes, err := ah.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get dApp update request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Debug("applying dApp update changes in diff")
			err = ah.applyChangesInDiff(changes)
			if err != nil {
				logger.Error("unable to apply dApp update changes in diff",
					zap.Any("error", err))
				rest.ERROR(w, err)
				ah.Memory.Cancel()
				return
			}

			logger.Info("committing dApp update changes")
			defer ah.Memory.Commit()
		} else {
			logger.Info("cancelling dApp update changes")
			defer ah.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

// HandleDelete - returns a handle function that
// deletes the dApp of the given path
func (ah *AppHandler) HandleDelete() rest.Handler {
	logger.Info("handling dApp delete request")
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppQueryDI{}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			logger.Error("unable to decode dApp delete request data",
				zap.Any("error", err))
			rest.ERROR(w, err)
			return
		}

		logger.Debug("initiating dApp delete transaction")
		ah.Memory.InitTransaction()

		err = ah.Memory.Apps().Delete(data.Ctx)
		if err != nil {
			logger.Error("unable to delete dApp",
				zap.String("dApp query", data.Ctx),
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		changes, err := ah.Memory.GetTransactionChanges()
		if err != nil {
			logger.Error("unable to get dApp delete request changes",
				zap.Any("error", err))
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		if !data.DryRun {
			logger.Debug("applying dApp delete changes in diff")
			err = ah.applyChangesInDiff(changes)
			if err != nil {
				logger.Error("unable to apply dApp delete changes in diff",
					zap.Any("error", err))
				rest.ERROR(w, err)
				ah.Memory.Cancel()
				return
			}

			logger.Info("committing Channel create changes")
			defer ah.Memory.Commit()
		} else {
			logger.Info("cancelling Channel create changes")
			defer ah.Memory.Cancel()
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
		err := handler.Operator.Nodes().DeleteNode(context.Background(), app.Meta.Parent, app.Meta.Name)
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
