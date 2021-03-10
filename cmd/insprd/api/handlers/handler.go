package handler

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils/diff"
)

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
	h := Handler{
		Memory:          memory,
		Operator:        operator,
		diffReactions:   []diff.DifferenceReaction{},
		changeReactions: []diff.ChangeReaction{},
	}
	h.initReactions()
	return &h
}
