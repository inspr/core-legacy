package nodes

import (
	"context"
	"fmt"
	"strings"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
	"gitlab.inspr.dev/inspr/core/pkg/enviroment"
	"gitlab.inspr.dev/inspr/core/pkg/meta"

	kubeApp "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type NodeOperator struct {
	clientSet kubernetes.Interface
}

func (nop *NodeOperator) retrieveKube() v1.DeploymentInterface {
	appsNamespace := GetK8SVariables().AppsNamespace
	return nop.clientSet.AppsV1().Deployments(appsNamespace)
}

// GetNode returns the node with the given name, if it exists.
// Otherwise, returns an error
func (nop *NodeOperator) GetNode(ctx context.Context, app *meta.App) (*meta.Node, error) {
	kube := nop.retrieveKube()
	insprEnv := enviroment.GetEnvironment().InsprEnvironment
	nodeName := parseNodeName(insprEnv, app.Meta.Parent, app.Spec.Node.Meta.Name)
	_, err := kube.Get(nodeName, metav1.GetOptions{})
	//
	return &app.Spec.Node, err
}

// Nodes is a NodeOperatorInterface that provides methods for node manipulation
func (nop *NodeOperator) Nodes() operators.NodeOperatorInterface {
	return &NodeOperator{}
}

// CreateNode deploys a new node structure, if it's information is valid.
// Otherwise, returns an error
func (nop *NodeOperator) CreateNode(ctx context.Context, app *meta.App) error {
	return nil
}

// UpdateNode updates a node that already exists, if the new structure is valid.
// Otherwise, returns an error.
func (nop *NodeOperator) UpdateNode(ctx context.Context, app *meta.App) error {
	var deploy *kubeApp.Deployment
	kube := nop.retrieveKube()
	// deploy = translatenodetodeploy
	_, err := kube.Update(deploy)

	return err
}

// DeleteNode deletes node with given name, if it exists. Otherwise, returns an error
func (nop *NodeOperator) DeleteNode(ctx context.Context, nodeContext string, nodeName string) error {
	return nil
}

// GetAllNodes returns a list of all the active nodes in the deployment, if there are any
func (nop *NodeOperator) GetAllNodes() []*meta.Node {
	kube := nop.retrieveKube()
	kube.List(metav1.ListOptions{})
	return nil

}

func parseNodeName(insprEnv string, context string, name string) string {
	s := fmt.Sprintf("%s.%s.%s", insprEnv, context, name)
	if s[0] == '.' {
		s = s[1:]
	}
	return strings.ToLower(s)
}
