package handler

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
	"go.uber.org/zap"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = zap.NewDevelopment(zap.Fields(zap.String("section", "insprd-api-handlers")))
}

// Handler is a general handler for inspr routes. It contains the necessary components
// for managing components on each route.
type Handler struct {
	Memory          memory.Manager
	Operator        operators.OperatorInterface
	diffReactions   []diff.DifferenceReaction
	changeReactions []diff.ChangeReaction
}

// NewHandler creates a handler from a memory manager and an operator. It also initializes the reactors for
// changes on the cluster.
func NewHandler(memory memory.Manager, operator operators.OperatorInterface) *Handler {
	logger.Info("creating new Insprd API handler")
	h := Handler{
		Memory:          memory,
		Operator:        operator,
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
