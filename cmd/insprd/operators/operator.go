package operators

import (
	"go.uber.org/zap"
	"inspr.dev/inspr/cmd/insprd/memory"
	"inspr.dev/inspr/cmd/insprd/memory/tree"
	"inspr.dev/inspr/cmd/insprd/operators/nodes"
	"inspr.dev/inspr/pkg/auth"
	"inspr.dev/inspr/pkg/logs"
)

var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = logs.Logger(zap.Fields(zap.String("section", "memory-tree")))
	// logger, _ = zap.NewDevelopment(zap.Fields(zap.String("section", "operators")))
	// logger = zap.NewNop()
}

// Operator is an operator for creating channels and nodes inside kubernetes
// that communicate via Sidecars. The operators need two environment variables
type Operator struct {
	channels *GenOp
	nodes    *nodes.NodeOperator
	mem      tree.Manager
}

// Nodes returns the nodes that communicate via sidecars inside kubernetes
func (op *Operator) Nodes() NodeOperatorInterface {
	logger.Info("summoning Nodes Operator")
	return op.nodes
}

// Channels returns the Channels Operator Interface for a given node
func (op *Operator) Channels() ChannelOperatorInterface {
	logger.Info("summoning Channels Operator")
	return op.channels
}

// NewOperator creates a node operator.
func NewOperator(memory memory.Manager, authenticator auth.Auth) (OperatorInterface, error) {
	var err error

	chOp := NewGeneralOperator(memory.Brokers(), memory.Tree())

	nOp, err := nodes.NewNodeOperator(memory.Tree(), authenticator, memory.Brokers())
	if err != nil {
		return nil, err
	}

	return &Operator{
		channels: chOp,
		nodes:    nOp,
		mem:      memory.Tree(),
	}, err
}
