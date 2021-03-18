package fake

import (
	"context"
	"fmt"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

// NodeOperator mock
type NodeOperator struct {
	nodes map[string]*meta.App
	err   error
}

// NewNodeOperator returns a mocked node operator that returns err on every function if err is not nil
func NewNodeOperator(err error) operators.NodeOperatorInterface {
	return &NodeOperator{
		nodes: make(map[string]*meta.App),
		err:   err,
	}
}

// CreateNode mock
func (o *NodeOperator) CreateNode(ctx context.Context, app *meta.App) (*meta.Node, error) {

	if o.err != nil {
		return nil, o.err
	}
	if _, ok := o.nodes[app.Meta.Parent+app.Meta.Name]; ok {
		return nil, ierrors.NewError().AlreadyExists().Message("node already exists").Build()
	}
	o.nodes[app.Meta.Parent+app.Meta.Name] = app
	return &app.Spec.Node, nil
}

// GetNode mock
func (o *NodeOperator) GetNode(ctx context.Context, app *meta.App) (*meta.Node, error) {
	if o.err != nil {
		return nil, o.err
	}
	node, ok := o.nodes[app.Meta.Parent+app.Meta.Name]
	if !ok {
		return nil, ierrors.NewError().NotFound().Message("node not found").Build()
	}
	return &node.Spec.Node, nil
}

// UpdateNode mock
func (o *NodeOperator) UpdateNode(ctx context.Context, app *meta.App) (*meta.Node, error) {
	if o.err != nil {
		return nil, o.err
	}
	if _, ok := o.nodes[app.Meta.Parent+app.Meta.Name]; !ok {
		return nil, ierrors.NewError().NotFound().Message("node not found").Build()
	}
	o.nodes[app.Meta.Parent+app.Meta.Name] = app
	return &app.Spec.Node, nil
}

// DeleteNode mock
func (o *NodeOperator) DeleteNode(ctx context.Context, nodeContext string, nodeName string) error {
	if o.err != nil {
		return o.err
	}
	_, ok := o.nodes[nodeContext+nodeName]
	if !ok {
		return ierrors.NewError().NotFound().Message(fmt.Sprintf("node %s not found", nodeContext+nodeName)).Build()
	}
	delete(o.nodes, nodeContext+nodeName)
	return nil
}

// GetAllNodes mock
func (o *NodeOperator) GetAllNodes() (ret []meta.Node) {

	for _, app := range o.nodes {
		ret = append(ret, app.Spec.Node)
	}
	return
}
