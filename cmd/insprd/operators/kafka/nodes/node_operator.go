package nodes

import (
	"context"
	"os"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/operators"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"

	kubeApp "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

//NodeOperator defines a node operations interface.
type NodeOperator struct {
	clientSet kubernetes.Interface
	memory    memory.Manager
}

func (no *NodeOperator) retrieveKube() v1.DeploymentInterface {
	appsNamespace := getK8SVariables().AppsNamespace
	return no.clientSet.AppsV1().Deployments(appsNamespace)
}

// GetNode returns the node with the given name, if it exists.
// Otherwise, returns an error
func (no *NodeOperator) GetNode(ctx context.Context, app *meta.App) (*meta.Node, error) {
	kube := no.retrieveKube()
	deployName := toDeploymentName(app)

	logger.Info("getting Node from k8s deployment",
		zap.String("deployment name", deployName))

	dep, err := kube.Get(deployName, metav1.GetOptions{})
	if err != nil {
		logger.Error("unable to find k8s deployment")
		return &meta.Node{}, ierrors.NewError().Message(err.Error()).Build()
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
	logger.Info("deploying a Node structure in k8s",
		zap.Any("node", app))

	logger.Debug("converting dApp to the k8s deployment to be created")
	var deploy *kubeApp.Deployment
	kube := no.retrieveKube()
	deploy = no.dAppToDeployment(app)

	logger.Debug("creating the k8s deployment")
	dep, err := kube.Create(deploy)
	if err != nil {
		logger.Error("unable to create the k8s deployment")
		return &meta.Node{}, ierrors.NewError().Message(err.Error()).Build()
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
	logger.Info("updating a Node structure in k8s",
		zap.Any("node", app))

	logger.Debug("converting dApp to the k8s deployment to be updated")
	var deploy *kubeApp.Deployment
	kube := no.retrieveKube()
	deploy = no.dAppToDeployment(app)

	logger.Debug("updating the k8s deployment")
	dep, err := kube.Update(deploy)
	if err != nil {
		logger.Error("unable to update the k8s deployment")
		return &meta.Node{}, ierrors.NewError().Message(err.Error()).Build()
	}

	node, err := toNode(dep)
	if err != nil {
		return &meta.Node{}, err
	}
	return &node, nil
}

// DeleteNode deletes node with given name, if it exists. Otherwise, returns an error
func (no *NodeOperator) DeleteNode(ctx context.Context, nodeContext string, nodeName string) error {
	logger.Info("deleting a Node structure in k8s",
		zap.String("node", nodeName),
		zap.String("context", nodeContext))

	var deployName string
	kube := no.retrieveKube()

	logger.Debug("getting name of the k8s deployment to be deleted")
	scope, _ := utils.JoinScopes(nodeContext, nodeName)
	app, _ := no.memory.Root().Apps().Get(scope)
	deployName = toDeploymentName(app)

	logger.Debug("deleting the k8s deployment",
		zap.String("deployment", deployName))
	err := kube.Delete(deployName, &metav1.DeleteOptions{})

	if err != nil {
		logger.Error("unable to delete the k8s deployment")
		return ierrors.NewError().Message(err.Error()).Build()
	}
	return nil
}

// GetAllNodes returns a list of all the active nodes in the deployment, if there are any
func (no *NodeOperator) GetAllNodes() []meta.Node {
	logger.Info("getting all Nodes in k8s deployments")

	var nodes []meta.Node
	kube := no.retrieveKube()
	list, _ := kube.List(metav1.ListOptions{})
	for _, item := range list.Items {
		node, _ := toNode(&item)
		nodes = append(nodes, node)
	}
	return nodes
}

// NewOperator initializes a k8s based kafka node operator with in cluster configuration
func NewOperator(memory memory.Manager) (nop *NodeOperator, err error) {
	nop = &NodeOperator{
		memory: memory,
	}
	if _, exists := os.LookupEnv("DEBUG"); exists {
		logger.Info("initializing node operator with debug configs")
		nop.clientSet = fake.NewSimpleClientset()
	} else {
		logger.Info("initializing node operator with production configs")
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}

		nop.clientSet, err = kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}
	}
	return nop, nil
}
