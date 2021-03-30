package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

// AppHandler - contains handlers that uses the AppMemory interface methods
type AppHandler struct {
	*Handler
}

// NewAppHandler - generates a new AppHandler through the memoryManager interface
func (handler *Handler) NewAppHandler() *AppHandler {
	return &AppHandler{
		handler,
	}
}

// HandleCreate - handler that generates the rest.Handle
// func to manage the http request
func (ah *AppHandler) HandleCreate() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}
		ah.Memory.InitTransaction()

		err = ah.Memory.Apps().Create(data.Ctx, &data.App)
		if err != nil {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}
		changes, err := ah.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		if !data.DryRun {
			err = ah.applyChangesInDiff(changes)
			if err != nil {
				rest.ERROR(w, err)
				ah.Memory.Cancel()
				return
			}
		}

		if !data.DryRun {
			defer ah.Memory.Commit()
		} else {
			defer ah.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}

	return rest.Handler(handler)
}

// HandleGet - handler that generates the rest.Handle
// func to manage the http request
func (ah *AppHandler) HandleGet() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, err)
			return
		}

		ah.Memory.InitTransaction()

		app, err := ah.Memory.Root().Apps().Get(data.Ctx)
		if err != nil {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		defer ah.Memory.Cancel()

		rest.JSON(w, http.StatusOK, app)
	}
	return rest.Handler(handler)
}

// HandleUpdate - handler that generates the rest.Handle
// func to manage the http request
func (ah *AppHandler) HandleUpdate() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, err)
			return
		}

		ah.Memory.InitTransaction()

		err = ah.Memory.Apps().Update(data.Ctx, &data.App)
		if err != nil {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}
		changes, err := ah.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		if !data.DryRun {
			err = ah.applyChangesInDiff(changes)

			if err != nil {
				rest.ERROR(w, err)
				ah.Memory.Cancel()
				return
			}
		}

		if !data.DryRun {
			defer ah.Memory.Commit()
		} else {
			defer ah.Memory.Cancel()
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

// HandleDelete - handler that generates the rest.Handle
// func to manage the http request
func (ah *AppHandler) HandleDelete() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, err)
			return
		}
		ah.Memory.InitTransaction()

		_, err = ah.Memory.Apps().Get(data.Ctx)
		if err != nil {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}

		err = ah.Memory.Apps().Delete(data.Ctx)
		if err != nil {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}
		changes, err := ah.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			ah.Memory.Cancel()
			return
		}
		if !data.DryRun {
			err = ah.applyChangesInDiff(changes)

			if err != nil {
				rest.ERROR(w, err)
				ah.Memory.Cancel()
				return
			}
		}

		if !data.DryRun {
			defer ah.Memory.Commit()
		} else {
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
