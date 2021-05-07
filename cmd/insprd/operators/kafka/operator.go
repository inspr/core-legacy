package kafka

import (
	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/cmd/insprd/operators"
	"github.com/inspr/inspr/cmd/insprd/operators/kafka/channels"
	"github.com/inspr/inspr/cmd/insprd/operators/kafka/nodes"
	"github.com/inspr/inspr/pkg/auth"
)

// Operator is an operator for creating channels and nodes inside kubernetes
// that communicate via Kafka. The operators need two environment variables
//
// KAFKA_BOOTSTRAP_SERVERS - tells the operators and sidecars where to connect to the kafka broker
//
// KAFKA_OFFSET_RESET - tells the operators and sidecars what that configuration should be when creating
// readers and writers
type Operator struct {
	channels operators.ChannelOperatorInterface
	nodes    *nodes.NodeOperator
}

// Nodes creates nodes that communicate via kafka inside kubernetes
func (op *Operator) Nodes() operators.NodeOperatorInterface {
	return op.nodes
}

// Channels creates channels inside kafka
func (op *Operator) Channels() operators.ChannelOperatorInterface {
	return op.channels
}

// NewKafkaOperator creates a kafka operator.
//
// View Operator
func NewKafkaOperator(memory memory.Manager, authenticator auth.Auth) (operators.OperatorInterface, error) {
	var err error
	var chOp operators.ChannelOperatorInterface
	chOp, err = channels.NewOperator(memory)
	if err != nil {
		return nil, err
	}
	nOp, err := nodes.NewOperator(memory, authenticator)
	if err != nil {
		return nil, err
	}

	return &Operator{
		channels: chOp,
		nodes:    nOp,
	}, err
}
