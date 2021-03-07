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
	Handler
}

// NewAppHandler - generates a new AppHandler through the memoryManager interface
func NewAppHandler(handler Handler) *AppHandler {
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
	var errs string // all errors from all operations

	// apply this on deleted channels
	deletedChannels := diff.NewDifferenceOperation(
		func(_ string, d diff.Difference) bool {
			// if the diff is the diff of a channel and the channel has been deleted
			return d.Kind&diff.ChannelKind > 0 && d.Operation&diff.Delete > 0
		},
		func(scope string, d diff.Difference) {
			err := handler.Operator.Channels().Delete(context.Background(), scope, d.Name) // delete the channel from the cluster
			if err != nil {
				errs = err.Error() + "\n"
			}
		},
	)

	// apply this on created channels
	createdChannels := diff.NewDifferenceOperation(
		func(scope string, d diff.Difference) bool {
			// if the diff is the diff of a channel and the channel has been created
			return d.Kind&diff.ChannelKind > 0 && d.Operation&diff.Create > 0
		},
		func(scope string, d diff.Difference) {
			ch, _ := handler.Memory.Channels().Get(scope, d.Name)                      // get the actual channel definition from memory
			err := handler.Operator.Channels().Create(context.Background(), scope, ch) // apply to the cluster
			if err != nil {
				errs += err.Error() + "\n"
			}
		},
	)

	// apply this on deleted apps
	deletedApps := diff.NewDifferenceOperation(
		func(scope string, d diff.Difference) bool {
			// if the diff is the diff of an app and the app has been created
			return d.Kind&diff.AppKind > 0 && d.Operation&diff.Delete > 0
		},
		func(scope string, d diff.Difference) {
			app, err := handler.Memory.Root().Apps().Get(scope + "." + d.Name) // get the app definition from the cluster
			if err != nil {
				errs += err.Error() + "\n"
				return
			}
			errs += handler.deleteApp(app) // delete app recursively (all nodes and channels defined) from the cluster
		},
	)

	// apply this on updated channel types
	updatedChannelTypes := diff.NewDifferenceOperation(
		func(scope string, d diff.Difference) bool {
			// if the diff is for a channel type and the channel type has been updated
			return d.Kind&diff.ChannelTypeKind > 0 && d.Operation&diff.Update > 0
		},
		func(scope string, d diff.Difference) {
			ct, _ := handler.Memory.ChannelTypes().Get(scope, d.Name)

			for _, channelName := range ct.ConnectedChannels { // for each channel connected to the channel type
				channel, _ := handler.Memory.Channels().Get(scope, channelName)

				for _, appName := range channel.ConnectedApps { // for each app connected to each channel
					app, _ := handler.Memory.Apps().Get(scope + "." + appName) // get the app definition from memory

					if app.Spec.Node.Spec.Image != "" { // if the app is a node, update it
						_, err := handler.Operator.Nodes().UpdateNode(context.Background(), app)
						if err != nil {
							errs += err.Error() + "\n"
						}
					}
				}
			}
		},
	)

	// apply this on updated channels
	updatedChannels := diff.NewDifferenceOperation(
		func(scope string, d diff.Difference) bool {
			// if the diff is for a channel and the channel has been updated
			return d.Kind&diff.ChannelKind > 0 && d.Operation&diff.Update > 0
		},
		func(scope string, d diff.Difference) {
			channel, _ := handler.Memory.Channels().Get(scope, d.Name)
			err := handler.Operator.Channels().Update(context.Background(), scope, channel)
			if err != nil {
				errs += err.Error() + "\n"
				return
			}
			// this updates the connected nodes, so that the environment variables are consistent with
			// the channel definition
			for _, appName := range channel.ConnectedApps { // for each app connected to each channel
				app, _ := handler.Memory.Apps().Get(scope + "." + appName)

				if app.Spec.Node.Spec.Image != "" { // if the app is a node, update it
					_, err := handler.Operator.Nodes().UpdateNode(context.Background(), app)
					if err != nil {
						errs += err.Error()
					}
				}
			}

		},
	)

	// aply this to updated nodes
	updatedNodes := diff.NewChangeOperation(
		func(c diff.Change) bool {
			from, _ := handler.Memory.Root().Apps().Get(c.Context)
			// if there is a change in a given context and that context is a node
			return from.Spec.Node.Spec.Image != ""
		},
		func(c diff.Change) {

			to, _ := handler.Memory.Apps().Get(c.Context)
			if to == nil || to.Spec.Node.Spec.Image == "" {

				scope, name, err := utils.RemoveLastPartInScope(c.Context)
				if err != nil {
					errs += err.Error() + "\n"
					return
				}
				err = handler.Operator.Nodes().DeleteNode(context.Background(), scope, name)
				if err != nil {
					errs += err.Error() + "\n"
				}
				return
			}
			_, err := handler.Operator.Nodes().UpdateNode(context.Background(), to) // update it in the cluster
			if err != nil {
				errs += err.Error() + "\n"
			}
		},
	)

	// apply each defined operation with its filter on the created diff
	changes.ForEachFiltered(updatedNodes)
	changes.ForEachDiffFiltered(
		deletedApps,
		deletedChannels,
		updatedChannels,
		createdChannels,
		updatedChannelTypes,
	)
	if errs != "" {
		return ierrors.NewError().InternalServer().Message(errs).Build()
	}
	return nil
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

func (handler *Handler) deleteApp(app *meta.App) (errs string) {
	if app.Spec.Node.Spec.Image != "" {
		err := handler.Operator.Nodes().DeleteNode(context.Background(), app.Meta.Parent, app.Meta.Name)
		if err != nil {
			errs += err.Error()
		}

	} else {
		for _, subApp := range app.Spec.Apps {
			errs += handler.deleteApp(subApp)
		}
		for c := range app.Spec.Channels {
			scope, _ := utils.JoinScopes(app.Meta.Parent, app.Meta.Name)

			err := handler.Operator.Channels().Delete(context.Background(), scope, c)
			if err != nil {
				errs += err.Error()
			}
		}
	}

	return
}
