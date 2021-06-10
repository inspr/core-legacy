package operators

import (
	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/cmd/insprd/memory/brokers"
	"github.com/inspr/inspr/cmd/insprd/operators/nodes"
	"github.com/inspr/inspr/pkg/auth"
	metabrokers "github.com/inspr/inspr/pkg/meta/brokers"
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

func (op *Operator) SetBrokerOperator(config metabrokers.BrokerConfiguration) error {
	return op.channels.SetOperator(config, op.mem)
}

// NewOperator creates a node operator.
func NewOperator(memory memory.Manager, authenticator auth.Auth, broker brokers.Manager) (OperatorInterface, error) {
	var err error

	chOp := NewGeneralOperator(memory)

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
