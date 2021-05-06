package nodes

import (
	"context"
	"os"

	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/pkg/auth"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/utils"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	cv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

//NodeOperator defines a node operations interface.
type NodeOperator struct {
	clientSet kubernetes.Interface
	memory    memory.Manager
	auth      auth.Auth
}

// Secrets returns the secret interface of the node operator
func (no *NodeOperator) Secrets() cv1.SecretInterface {
	appsNamespace := getK8SVariables().AppsNamespace
	return no.clientSet.CoreV1().Secrets(appsNamespace)
}

// Services returns the service interface for the node operator
func (no *NodeOperator) Services() cv1.ServiceInterface {
	appsNamespace := getK8SVariables().AppsNamespace
	return no.clientSet.CoreV1().Services(appsNamespace)
}

// Deployments returns the deployment interface for the k8s operator
func (no *NodeOperator) Deployments() v1.DeploymentInterface {
	appsNamespace := getK8SVariables().AppsNamespace
	return no.clientSet.AppsV1().Deployments(appsNamespace)
}

// CreateNode deploys a new node structure, if it's information is valid.
// Otherwise, returns an error
func (no *NodeOperator) CreateNode(ctx context.Context, app *meta.App) (*meta.Node, error) {
	logger.Info("deploying a Node structure in k8s",
		zap.Any("node", app))

	for _, applicable := range no.dappApplications(app) {
		err := applicable.create(no)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// UpdateNode updates a node that already exists, if the new structure is valid.
// Otherwise, returns an error.
func (no *NodeOperator) UpdateNode(ctx context.Context, app *meta.App) (*meta.Node, error) {
	logger.Info("deploying a Node structure in k8s",
		zap.Any("node", app))

	for _, applicable := range no.dappApplications(app) {
		err := applicable.update(no)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// DeleteNode deletes node with given name, if it exists. Otherwise, returns an error
func (no *NodeOperator) DeleteNode(ctx context.Context, nodeContext string, nodeName string) error {
	logger.Info("deleting a Node structure in k8s",
		zap.String("node", nodeName),
		zap.String("context", nodeContext))

	logger.Debug("getting name of the k8s deployment to be deleted")
	scope, _ := utils.JoinScopes(nodeContext, nodeName)
	app, _ := no.memory.Root().Apps().Get(scope)

	logger.Info("deploying a Node structure in k8s",
		zap.Any("node", app))

	for _, applicable := range no.dappApplications(app) {
		err := applicable.del(no)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewOperator initializes a k8s based kafka node operator with in cluster configuration
func NewOperator(memory memory.Manager, a auth.Auth) (nop *NodeOperator, err error) {
	nop = &NodeOperator{
		memory: memory,
		auth:   a,
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
