package operators

import (
	"inspr.dev/inspr/cmd/insprd/memory"
	"inspr.dev/inspr/cmd/insprd/memory/brokers"
	"inspr.dev/inspr/cmd/insprd/operators/nodes"
	"inspr.dev/inspr/pkg/auth"
)

// Operator is an operator for creating channels and nodes inside kubernetes
// that communicate via Sidecars. The operators need two environment variables
type Operator struct {
	channels *GenOp
	nodes    *nodes.NodeOperator
	mem      memory.Manager
}

// Nodes returns the nodes that communicate via sidecars inside kubernetes
func (op *Operator) Nodes() NodeOperatorInterface {
	return op.nodes
}

// Channels returns the Channels Operator Interface for a given node
func (op *Operator) Channels() ChannelOperatorInterface {
	return op.channels
}

// NewOperator creates a node operator.
func NewOperator(memory memory.Manager, authenticator auth.Auth, broker brokers.Manager) (OperatorInterface, error) {
	var err error

	chOp := NewGeneralOperator(broker, memory)

	nOp, err := nodes.NewNodeOperator(memory, authenticator, broker)
	if err != nil {
		return nil, err
	}

	return &Operator{
		channels: chOp,
		nodes:    nOp,
		mem:      memory,
	}, err
}
