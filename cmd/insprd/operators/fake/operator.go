package fake

import (
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

type FakeOperator struct {
	nodes    NodeOperator
	channels ChannelOperator
}

func (f *FakeOperator) Channels() operators.ChannelOperatorInterface {
	return f.channels
}
func (f *FakeOperator) Nodes() operators.NodeOperatorInterface {
	return f.nodes
}
func NewFakeOperator() operators.OperatorInterface {
	return &FakeOperator{
		nodes: NodeOperator{
			nodes: make(map[string]*meta.App),
		},
		channels: ChannelOperator{
			channels: make(map[string]*meta.Channel),
		},
	}
}
