package handler

import (
	"go.uber.org/zap"
	"inspr.dev/inspr/cmd/insprd/memory"
	"inspr.dev/inspr/cmd/insprd/operators"
	"inspr.dev/inspr/pkg/auth"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/logs"
	"inspr.dev/inspr/pkg/meta/utils/diff"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = logs.Logger(zap.Fields(zap.String("section", "api")))
}

// Handler is a general handler for inspr routes. It contains the necessary components
// for managing components on each route.
type Handler struct {
	Memory          memory.Manager
	Operator        operators.OperatorInterface
	Auth            auth.Auth
	diffReactions   []diff.DifferenceReaction
	changeReactions []diff.ChangeReaction
}

// NewHandler creates a handler from a memory manager and an operator. It also initializes the reactors for
// changes on the cluster.
func NewHandler(memory memory.Manager, operator operators.OperatorInterface, auth auth.Auth) *Handler {
	logger.Info("creating new Insprd API handler")
	h := Handler{
		Memory:          memory,
		Operator:        operator,
		Auth:            auth,
		diffReactions:   []diff.DifferenceReaction{},
		changeReactions: []diff.ChangeReaction{},
	}
	h.initReactions()
	return &h
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

func (handler *Handler) applyChangesInDiff(changes diff.Changelog) error {
	logger.Debug("trying to apply changes in diff",
		zap.Any("changes", changes))
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

// GetCancel returns the transaction cancelation function for the operations
func (handler *Handler) GetCancel() func() {
	return handler.Memory.Tree().Cancel
}
