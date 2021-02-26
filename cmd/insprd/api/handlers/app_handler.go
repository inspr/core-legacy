package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/api/models"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
	"gitlab.inspr.dev/inspr/core/pkg/rest"
)

// AppHandler - contains handlers that uses the AppMemory interface methods
type AppHandler struct {
	memory.AppMemory
	mem memory.Manager
	op  operators.OperatorInterface
}

// NewAppHandler - generates a new AppHandler through the memoryManager interface
func NewAppHandler(memManager memory.Manager, op operators.OperatorInterface) *AppHandler {
	return &AppHandler{
		AppMemory: memManager.Apps(),
		op:        op,
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
		ah.InitTransaction()
		if !data.DryRun {
			defer ah.Commit()
		} else {
			defer ah.Cancel()
		}
		err = ah.CreateApp(data.Ctx, &data.App)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		changes, err := ah.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		errs := ""
		changes.ForEach(func(change diff.Change) {
			app, _ := ah.Get(change.Context)

			if change.Kind|diff.ChannelKind > 0 { // if there is any difference in channels
				change.FilterKind(diff.ChannelKind).ForEach(func(d diff.Difference) {
					err := ah.op.Channels().Create(context.Background(), change.Context, app.Spec.Channels[d.Name])
					if err != nil {
						errs += err.Error() + "\n"
					}
				})
			} else if app.Spec.Node.Spec.Image != "" {

				_, err := ah.op.Nodes().CreateNode(context.Background(), app)
				if err != nil {
					errs += err.Error() + "\n"
				}
			}
		})
		if errs != "" {
			rest.ERROR(w, ierrors.NewError().InternalServer().Message(errs).Build())
			return
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

// HandleGetAppByRef - handler that generates the rest.Handle
// func to manage the http request
func (ah *AppHandler) HandleGetAppByRef() rest.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		data := models.AppQueryDI{}
		context.Background()
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&data)
		if err != nil || !data.Valid {
			rest.ERROR(w, err)
			return
		}

		ah.InitTransaction()
		defer ah.Cancel()

		app, err := ah.Get(data.Ctx)
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

		ah.InitTransaction()
		if !data.DryRun {
			defer ah.Commit()
		} else {
			defer ah.Cancel()
		}
		err = ah.UpdateApp(data.Ctx, &data.App)
		if err != nil {
			rest.ERROR(w, err)
		}
		changes, err := ah.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			return
		}

		var errs string // all errors from all operations

		// apply this on deleted channels
		deletedChannels := diff.NewDifferenceOperation(
			func(_ string, d diff.Difference) bool {
				// if the diff is the diff of a channel and the channel has been deleted
				return d.Kind&diff.ChannelKind > 0 && d.Operation&diff.Delete > 0
			},
			func(scope string, d diff.Difference) {
				err := ah.op.Channels().Delete(context.Background(), scope, d.Name) // delete the channel from the cluster
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
				ch, _ := ah.mem.Channels().Get(scope, d.Name)                   // get the actual channel definition from memory
				err := ah.op.Channels().Create(context.Background(), scope, ch) // apply to the cluster
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
				app, err := ah.Get(scope + "." + d.Name) // get the app definition from the cluster
				if err != nil {
					errs += err.Error() + "\n"
					return
				}
				errs += ah.deleteApp(app) // delete app recursively (all nodes and channels defined) from the cluster
			},
		)

		// apply this on updated channel types
		updatedChannelTypes := diff.NewDifferenceOperation(
			func(scope string, d diff.Difference) bool {
				// if the diff is for a channel type and the channel type has been updated
				return d.Kind&diff.ChannelTypeKind > 0 && d.Operation&diff.Update > 0
			},
			func(scope string, d diff.Difference) {
				ct, _ := ah.mem.ChannelTypes().Get(scope, d.Name)

				for _, channelName := range ct.ConnectedChannels { // for each channel connected to the channel type
					channel, _ := ah.mem.Channels().Get(scope, channelName)

					for _, appName := range channel.ConnectedApps { // for each app connected to each channel
						app, _ := ah.mem.Apps().Get(scope + "." + appName) // get the app definition from memory

						if app.Spec.Node.Spec.Image != "" { // if the app is a node, update it
							_, err := ah.op.Nodes().UpdateNode(context.Background(), app)
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
				channel, _ := ah.mem.Channels().Get(scope, d.Name)
				err := ah.op.Channels().Update(context.Background(), scope, channel)
				if err != nil {
					errs += err.Error() + "\n"
					return
				}
				// this updates the connected nodes, so that the environment variables are consistent with
				// the channel definition
				for _, appName := range channel.ConnectedApps { // for each app connected to each channel
					app, _ := ah.mem.Apps().Get(scope + "." + appName)

					if app.Spec.Node.Spec.Image != "" { // if the app is a node, update it
						_, err := ah.op.Nodes().UpdateNode(context.Background(), app)
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
				app, _ := ah.mem.Apps().Get(c.Context)
				// if there is a change in a given context and that context is a node
				return app.Spec.Node.Spec.Image != ""
			},
			func(c diff.Change) {

				app, _ := ah.mem.Apps().Get(c.Context)                        // get the app from memory
				_, err := ah.op.Nodes().UpdateNode(context.Background(), app) // update it in the cluster
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
			rest.ERROR(w, ierrors.NewError().InternalServer().Message(errs).Build())
			return
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
		ah.InitTransaction()
		if !data.DryRun {
			defer ah.Commit()
		} else {
			defer ah.Cancel()
		}
		app, err := ah.Get(data.Ctx)
		if err != nil {
			rest.ERROR(w, err)
		}
		err = ah.DeleteApp(data.Ctx)
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		changes, err := ah.GetTransactionChanges()
		if err != nil {
			rest.ERROR(w, err)
			return
		}
		errs := ah.deleteApp(app)

		if errs != "" {
			rest.ERROR(w, ierrors.NewError().InternalServer().Message(errs).Build())
			return
		}

		rest.JSON(w, http.StatusOK, changes)
	}
	return rest.Handler(handler)
}

func (ah *AppHandler) deleteApp(app *meta.App) (errs string) {
	for _, subApp := range app.Spec.Apps {
		if app.Spec.Node.Meta.Name != "" {
			err := ah.op.Nodes().DeleteNode(context.Background(), app.Meta.Parent, app.Spec.Node.Meta.Name)
			if err != nil {
				errs += err.Error()
			}

		} else {
			errs += ah.deleteApp(subApp)
		}
	}
	for c := range app.Spec.Channels {
		err := ah.op.Channels().Delete(context.Background(), app.Meta.Parent+"."+app.Meta.Name, c)
		if err != nil {
			errs += err.Error()
		}
	}
	return
}
