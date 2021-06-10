package fake

import (
	"github.com/inspr/inspr/cmd/insprd/operators"
	"github.com/inspr/inspr/pkg/meta"
	metabrokers "github.com/inspr/inspr/pkg/meta/brokers"
)

// Operator mock
type Operator struct {
	nodes    *NodeOperator
	channels *ChannelOperator
}

// Channels mock
func (f *Operator) Channels() operators.ChannelOperatorInterface {
	return f.channels
}

// Nodes mock
func (f *Operator) Nodes() operators.NodeOperatorInterface {
	return f.nodes
}

func (f *Operator) SetBrokerOperator(config metabrokers.BrokerConfiguration) error {
	return nil
}

// NewFakeOperator creates a simple operator that only acts in memory
func NewFakeOperator() operators.OperatorInterface {
	return &Operator{
		nodes: &NodeOperator{
			nodes: make(map[string]*meta.App),
		},
		channels: &ChannelOperator{
			channels: make(map[string]*meta.Channel),
		},
	}
}
