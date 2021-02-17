package nodes

import (
	"context"
	"fmt"
	"strings"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"

	kubeApp "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

//NodeOperator defines a node operations interface.
type NodeOperator struct {
	clientSet kubernetes.Interface
}

func (no *NodeOperator) retrieveKube() v1.DeploymentInterface {
	appsNamespace := getK8SVariables().AppsNamespace
	return no.clientSet.AppsV1().Deployments(appsNamespace)
}

// GetNode returns the node with the given name, if it exists.
// Otherwise, returns an error
func (no *NodeOperator) GetNode(ctx context.Context, app *meta.App) (*meta.Node, error) {
	kube := no.retrieveKube()
	insprEnv := environment.GetEnvironment().InsprEnvironment
	nodeName := parseNodeName(insprEnv, app.Meta.Parent, app.Spec.Node.Meta.Name)
	dep, err := kube.Get(nodeName, metav1.GetOptions{})
	if err != nil {
		return &meta.Node{}, ierrors.NewError().Message("could't get deployment from kubernetes").InnerError(err).Build()
	}
	node, err := toNode(dep)
	if err != nil {
		return &meta.Node{}, err
	}
	return &node, nil
}

// Nodes is a NodeOperatorInterface that provides methods for node manipulation
func (no *NodeOperator) Nodes() operators.NodeOperatorInterface {
	return &NodeOperator{}
}

// CreateNode deploys a new node structure, if it's information is valid.
// Otherwise, returns an error
func (no *NodeOperator) CreateNode(ctx context.Context, app *meta.App) (*meta.Node, error) {
	var deploy *kubeApp.Deployment
	kube := no.retrieveKube()
	deploy = dAppToDeployment(app)
	dep, err := kube.Create(deploy)
	if err != nil {
		return &meta.Node{}, ierrors.NewError().Message("could't create deployment from kubernetes").InnerError(err).Build()
	}
	node, err := toNode(dep)
	if err != nil {
		return &meta.Node{}, err
	}
	return &node, nil
}

// UpdateNode updates a node that already exists, if the new structure is valid.
// Otherwise, returns an error.
func (no *NodeOperator) UpdateNode(ctx context.Context, app *meta.App) (*meta.Node, error) {
	var deploy *kubeApp.Deployment
	kube := no.retrieveKube()
	deploy = dAppToDeployment(app)
	dep, err := kube.Update(deploy)
	if err != nil {
		return &meta.Node{}, ierrors.NewError().Message("could't update deployment from kubernetes").InnerError(err).Build()
	}
	node, err := toNode(dep)
	if err != nil {
		return &meta.Node{}, err
	}
	return &node, nil
}

// DeleteNode deletes node with given name, if it exists. Otherwise, returns an error
func (no *NodeOperator) DeleteNode(ctx context.Context, nodeContext string, nodeName string) error {
	var deploy string
	kube := no.retrieveKube()
	deploy = parseNodeName(getK8SVariables().AppsNamespace, nodeContext, nodeName)
	err := kube.Delete(deploy, &metav1.DeleteOptions{})

	if err != nil {
		return ierrors.NewError().Message("could't delete deployment from kubernetes").InnerError(err).Build()
	}
	return nil
}

// GetAllNodes returns a list of all the active nodes in the deployment, if there are any
func (no *NodeOperator) GetAllNodes() []meta.Node {
	var nodes []meta.Node
	kube := no.retrieveKube()
	list, _ := kube.List(metav1.ListOptions{})
	for _, item := range list.Items {
		node, _ := toNode(&item)
		nodes = append(nodes, node)
	}
	return nodes
}

func parseNodeName(insprEnv string, context string, name string) string {
	s := fmt.Sprintf("%s.%s.%s", insprEnv, context, name)
	if s[0] == '.' {
		s = s[1:]
	}
	return strings.ToLower(s)
}
