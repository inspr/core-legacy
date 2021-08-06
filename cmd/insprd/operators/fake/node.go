package fake

import (
	"context"

	"inspr.dev/inspr/cmd/insprd/operators"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
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
		return nil, ierrors.New("node already exists").AlreadyExists()
	}
	o.nodes[app.Meta.Parent+app.Meta.Name] = app
	return &app.Spec.Node, nil
}

// GetNode mock
func (o *NodeOperator) GetNode(ctx context.Context, app *meta.App) (*meta.Node, error) {
	if o.err != nil {
		return nil, o.err
	}

	nodeKey := app.Meta.Parent + app.Meta.Name
	node, ok := o.nodes[nodeKey]
	if !ok {
		return nil, ierrors.New(
			"node not found, searched for: %s", nodeKey,
		).NotFound()
	}
	return &node.Spec.Node, nil
}

// UpdateNode mock
func (o *NodeOperator) UpdateNode(ctx context.Context, app *meta.App) (*meta.Node, error) {
	if o.err != nil {
		return nil, o.err
	}

	nodeKey := app.Meta.Parent + app.Meta.Name
	if _, ok := o.nodes[nodeKey]; !ok {
		return nil, ierrors.New(
			"node not found, searched for: %s", nodeKey,
		).NotFound()
	}
	o.nodes[app.Meta.Parent+app.Meta.Name] = app
	return &app.Spec.Node, nil
}

// DeleteNode mock
func (o *NodeOperator) DeleteNode(ctx context.Context, nodeContext string, nodeName string) error {
	if o.err != nil {
		return o.err
	}

	nodeKey := nodeContext + nodeName
	_, ok := o.nodes[nodeKey]
	if !ok {
		return ierrors.New(
			"node not found, searched for: %s", nodeKey,
		).NotFound()
	}

	delete(o.nodes, nodeKey)
	return nil
}

// GetAllNodes mock
func (o *NodeOperator) GetAllNodes() (ret []meta.Node) {

	for _, app := range o.nodes {
		ret = append(ret, app.Spec.Node)
	}
	return
}
