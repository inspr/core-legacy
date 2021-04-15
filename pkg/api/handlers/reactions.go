package handler

import (
	"context"

	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta/utils"
	"github.com/inspr/inspr/pkg/meta/utils/diff"
)

var createdNodes func(handler *Handler) diff.ChangeReaction = func(handler *Handler) diff.ChangeReaction {
	return diff.NewChangeReaction(
		func(c diff.Change) bool {
			_, errFrom := handler.Memory.Root().Apps().Get(c.Context)
			to, errTo := handler.Memory.Apps().Get(c.Context)
			return (errFrom != nil && errTo == nil && to.Spec.Node.Spec.Image != "")
		},
		func(c diff.Change) error {
			to, _ := handler.Memory.Apps().Get(c.Context)
			_, err := handler.Operator.Nodes().CreateNode(context.Background(), to)
			return err
		},
	)
}

// apply this on deleted channels
var deletedChannels func(handler *Handler) diff.DifferenceReaction = func(handler *Handler) diff.DifferenceReaction {
	return diff.NewDifferenceReaction(
		func(_ string, d diff.Difference) bool {
			// if the diff is the diff of a channel and the channel has been deleted
			return d.Kind&diff.ChannelKind > 0 && d.Operation&diff.Delete > 0
		},
		func(scope string, d diff.Difference) error {
			return handler.Operator.Channels().Delete(context.Background(), scope, d.Name) // delete the channel from the cluster
		},
	)
}

// apply this on created channels
var createdChannels func(handler *Handler) diff.DifferenceReaction = func(handler *Handler) diff.DifferenceReaction {
	return diff.NewDifferenceReaction(
		func(scope string, d diff.Difference) bool {
			// if the diff is the diff of a channel and the channel has been created
			return d.Kind&diff.ChannelKind > 0 && d.Operation&diff.Create > 0
		},
		func(scope string, d diff.Difference) error {
			ch, _ := handler.Memory.Channels().Get(scope, d.Name)                      // get the actual channel definition from memory
			return handler.Operator.Channels().Create(context.Background(), scope, ch) // apply to the cluster

		},
	)
}

// apply this on deleted apps
var deletedApps func(handler *Handler) diff.DifferenceReaction = func(handler *Handler) diff.DifferenceReaction {
	return diff.NewDifferenceReaction(
		func(scope string, d diff.Difference) bool {
			// if the diff is the diff of an app and the app has been created
			return d.Kind&diff.AppKind > 0 && d.Operation&diff.Delete > 0
		},
		func(scope string, d diff.Difference) error {
			scope, _ = utils.JoinScopes(scope, d.Name)
			app, err := handler.Memory.Root().Apps().Get(scope) // get the app definition from the cluster
			if err != nil {
				return err
			}
			return handler.deleteApp(app) // delete app recursively (all nodes and channels defined) from the cluster
		},
	)
}

// apply this on updated channel types
var updatedChannelTypes func(handler *Handler) diff.DifferenceReaction = func(handler *Handler) diff.DifferenceReaction {
	return diff.NewDifferenceReaction(
		func(scope string, d diff.Difference) bool {
			// if the diff is for a channel type and the channel type has been updated
			return d.Kind&diff.ChannelTypeKind > 0 && d.Operation&diff.Update > 0
		},
		func(scope string, d diff.Difference) error {
			errors := ierrors.MultiError{
				Errors: []error{},
			}
			ct, _ := handler.Memory.ChannelTypes().Get(scope, d.Name)

			for _, channelName := range ct.ConnectedChannels { // for each channel connected to the channel type
				channel, _ := handler.Memory.Channels().Get(scope, channelName)

				for _, appName := range channel.ConnectedApps { // for each app connected to each channel
					newScope, _ := utils.JoinScopes(scope, appName)
					app, _ := handler.Memory.Apps().Get(newScope) // get the app definition from memory

					if app.Spec.Node.Spec.Image != "" { // if the app is a node, update it
						_, err := handler.Operator.Nodes().UpdateNode(context.Background(), app)
						if err != nil {
							errors.Add(err)
						}
					}
				}
			}
			if len(errors.Errors) > 0 {
				return &errors
			}
			return nil
		},
	)
}

// apply this on updated channels
var updatedChannels func(handler *Handler) diff.DifferenceReaction = func(handler *Handler) diff.DifferenceReaction {
	return diff.NewDifferenceReaction(
		func(scope string, d diff.Difference) bool {
			// if the diff is for a channel and the channel has been updated
			return d.Kind&diff.ChannelKind > 0 && d.Operation&diff.Update > 0
		},
		func(scope string, d diff.Difference) error {
			errs := ierrors.MultiError{
				Errors: []error{},
			}
			channel, _ := handler.Memory.Channels().Get(scope, d.Name)
			err := handler.Operator.Channels().Update(context.Background(), scope, channel)
			if err != nil {
				return err
			}
			// this updates the connected nodes, so that the environment variables are consistent with
			// the channel definition
			for _, appName := range channel.ConnectedApps { // for each app connected to each channel
				app, _ := handler.Memory.Apps().Get(scope + "." + appName)

				if app.Spec.Node.Spec.Image != "" { // if the app is a node, update it
					_, err := handler.Operator.Nodes().UpdateNode(context.Background(), app)
					if err != nil {
						errs.Add(err)
					}
				}
			}
			if !errs.Empty() {
				return &errs
			}
			return nil
		},
	)
}

// apply this to updated nodes
var updatedNodes func(handler *Handler) diff.ChangeReaction = func(handler *Handler) diff.ChangeReaction {
	return diff.NewChangeReaction(
		func(c diff.Change) bool {
			from, _ := handler.Memory.Root().Apps().Get(c.Context)
			// if there is a change in a given context and that context is a node
			return from != nil && from.Spec.Node.Spec.Image != ""
		},
		func(c diff.Change) error {
			errs := ierrors.MultiError{
				Errors: []error{},
			}
			to, _ := handler.Memory.Apps().Get(c.Context)
			if to == nil || to.Spec.Node.Spec.Image == "" {
				return nil
			}
			_, err := handler.Operator.Nodes().UpdateNode(context.Background(), to) // update it in the cluster
			if err != nil {
				errs.Add(err)
			}
			if !errs.Empty() {
				return &errs
			}
			return nil
		},
	)
}

// apply this to updated aliases
var updatedAliases func(handler *Handler) diff.DifferenceReaction = func(handler *Handler) diff.DifferenceReaction {
	return diff.NewDifferenceReaction(
		func(scope string, d diff.Difference) bool {
			return (d.Kind|diff.AliasKind > 0) && (d.Operation&diff.Update > 0)
		},
		func(scope string, d diff.Difference) error {
			appName, boundaryName, _ := utils.RemoveLastPartInScope(d.Name)
			newScope, _ := utils.JoinScopes(scope, appName)
			app, err := handler.Memory.Apps().Get(newScope)
			if err == nil && app.Spec.Boundary.Input.Union(app.Spec.Boundary.Output).Contains(boundaryName) {
				_, err := handler.Operator.Nodes().UpdateNode(context.Background(), app)
				if err != nil {
					return err
				}
			}
			return nil
		},
	)
}

func (h *Handler) initReactions() {
	h.addChangeReactor(
		updatedNodes(h),
		createdNodes(h),
	)

	h.addDiffReactor(
		createdChannels(h),
		deletedChannels(h),
		deletedApps(h),
		updatedChannels(h),
		updatedChannelTypes(h),
		updatedAliases(h),
	)
}
