package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
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

// HandleCreateApp - handler that generates the rest.Handle
// func to manage the http request
func (ah *AppHandler) HandleCreateApp() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, err)
			return
		}
		ah.Memory.InitTransaction()
		if !data.DryRun {
			defer ah.Memory.Commit()
		} else {
			defer ah.Memory.Cancel()
		}

		err = ah.Memory.Apps().CreateApp(data.Ctx, &data.App)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		changes, err := ah.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		if !data.DryRun {
			err = ah.applyChangesInDiff(changes)
			if err != nil {
				rest.ERROR(w, err)
			}
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

func (handler *Handler) applyChangesInDiff(changes diff.Changelog) error {
	errs := ierrors.MultiError{
		Errors: []error{},
	}
	errs.Add(changes.ForEachDiffFiltered(handler.diffReactions...))
	errs.Add(changes.ForEachFiltered(handler.changeReactions...))
	if errs.Empty() {
		return nil
	}

	return ierrors.NewError().Message(errs.Error()).Build()
}

// HandleGetAppByRef - handler that generates the rest.Handle
// func to manage the http request
func (ah *AppHandler) HandleGetAppByRef() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, err)
			return
		}

		ah.Memory.InitTransaction()
		defer ah.Memory.Cancel()

		app, err := ah.Memory.Root().Apps().Get(data.Ctx)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		rest.JSON(w, http.StatusOK, app)
	}
	return rest.Handler(handler)
}

// HandleUpdateApp - handler that generates the rest.Handle
// func to manage the http request
func (ah *AppHandler) HandleUpdateApp() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, err)
			return
		}

		ah.Memory.InitTransaction()
		if !data.DryRun {
			defer ah.Memory.Commit()
		} else {
			defer ah.Memory.Cancel()
		}
		err = ah.Memory.Apps().UpdateApp(data.Ctx, &data.App)
		if err != nil {
			rest.ERROR(w, err)
		}
		changes, err := ah.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		if !data.DryRun {
			err = ah.applyChangesInDiff(changes)

			if err != nil {
				rest.ERROR(w, err)
				return
			}
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

// HandleDeleteApp - handler that generates the rest.Handle
// func to manage the http request
func (ah *AppHandler) HandleDeleteApp() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppQueryDI{}
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, err)
			return
		}
		ah.Memory.InitTransaction()
		if !data.DryRun {
			defer ah.Memory.Commit()
		} else {
			defer ah.Memory.Cancel()
		}
		_, err = ah.Memory.Apps().Get(data.Ctx)
		if err != nil {
			rest.ERROR(w, err)
		}

		err = ah.Memory.Apps().DeleteApp(data.Ctx)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		changes, err := ah.Memory.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		if !data.DryRun {
			err = ah.applyChangesInDiff(changes)

			if err != nil {
				rest.ERROR(w, err)
				return
			}
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

func (handler *Handler) addDiffReactor(op ...diff.DifferenceReaction) {
	if handler.diffReactions == nil {
		handler.diffReactions = []diff.DifferenceReaction{}
	}
	handler.diffReactions = append(handler.diffReactions, op...)
}

func (handler *Handler) addChangeReactor(op ...diff.ChangeReaction) {
	if handler.changeReactions == nil {
		handler.changeReactions = []diff.ChangeReaction{}
	}
	handler.changeReactions = append(handler.changeReactions, op...)
}
