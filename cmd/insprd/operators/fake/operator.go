package fake

import (
	"inspr.dev/inspr/cmd/insprd/operators"
	"inspr.dev/inspr/pkg/meta"
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
