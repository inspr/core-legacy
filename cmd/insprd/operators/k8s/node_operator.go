package operators

// NodeOperator defines a NodeOperatorInterface
type NodeOperator struct {
	node *meta.Node
}

// GetNodes returns the node with the given name, if it exists.
// Otherwise, returns an error
func (no *NodeOperator) GetNodes(ctx context.context, context string, nodeName string) (*meta.Node, error) {
	kube := pipe.retrieveKube()

	for {
		msg := <-pipe.getChan
		node := msg.Node

		dep, err := kube.Get(node.ToDeployment(), metav1.GetOptions{})
		code := treatK8sError(err, cerrors.InsprErrorCode(6))
		if code != 0 {
			insprError := generateInsprError(5, err, "error in node back end update")
			pipe.responseChan <- generateInsprResponseMessge(insprError, msg.Spec)
			continue
		}

		node.NodeStatus, err = pipe.getNodeStatus(dep)
		fmt.Println(node.NodeStatus)
		if err != nil {
			logrus.Println("unable to get pod status")
			insprError := generateInsprError(5, err, "unable to get pod status")
			pipe.responseChan <- generateInsprResponseMessge(insprError, msg.Spec)
		}
		_, err = pipe.registry.PutNode(node)
		if err != nil {
			insprError := generateInsprError(5, err, "error in node back end update")
			pipe.responseChan <- generateInsprResponseMessge(insprError, msg.Spec)
		} else {
			pipe.responseChan <- generateInsprResponseMessge(nil, msg.Spec)
		}
	}
}

// Nodes is a NodeOperatorInterface that provides methods for node manipulation
func (no *NodeOperator) Nodes() NodeOperatorInterface {
	return &NodeOperator{
		node: &meta.Node{},
	}
}

// CreateNode deploys a new node structure, if it's information is valid.
// Otherwise, returns an error
func (no *NodeOperator) CreateNode(ctx context.context, context string, node *meta.Node) error {
	return nil
}

// UpdateNode updates a node that already exists, if the new structure is valid.
// Otherwise, returns an error.
func (no *NodeOperator) UpdateNode(ctx context.context, context string, node *meta.Node) error {
	return nil
}

// DeleteNode deletes node with given name, if it exists. Otherwise, returns an error
func (no *NodeOperator) DeleteNode(ctx context.context, context string, nodeName string) error {
	return nil
}

// GetAllNodes returns a list of all the active nodes in the deployment, if there are any
func (no *NodeOperator) GetAllNodes() []*meta.Node {
	return nil
}
