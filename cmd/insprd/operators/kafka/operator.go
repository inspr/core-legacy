package kafka

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators/kafka/channels"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators/kafka/nodes"
)

type KafkaOperator struct {
	channels *channels.ChannelOperator
	nodes    *nodes.NodeOperator
}

func (op *KafkaOperator) Nodes() operators.NodeOperatorInterface {
	return op.nodes
}

func (op *KafkaOperator) Channels() operators.ChannelOperatorInterface {
	return op.channels
}

func NewKafkaOperator() (operators.OperatorInterface, error) {
	chOp, err := channels.NewOperator()
	if err != nil {
		return nil, err
	}

	nOp, err := nodes.NewOperator()
	if err != nil {
		return nil, err
	}

	return &KafkaOperator{
		channels: chOp,
		nodes:    nOp,
	}, err
}
