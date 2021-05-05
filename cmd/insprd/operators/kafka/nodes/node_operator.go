package nodes

import (
	"context"
	"math"
	"os"
	"strconv"

	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/pkg/auth"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/utils"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (no *NodeOperator) secrets() cv1.SecretInterface {
	appsNamespace := getK8SVariables().AppsNamespace
	return no.clientSet.CoreV1().Secrets(appsNamespace)
}
func (no *NodeOperator) services() cv1.ServiceInterface {
	appsNamespace := getK8SVariables().AppsNamespace
	return no.clientSet.CoreV1().Services(appsNamespace)
}

func (no *NodeOperator) deployments() v1.DeploymentInterface {
	appsNamespace := getK8SVariables().AppsNamespace
	return no.clientSet.AppsV1().Deployments(appsNamespace)
}

// GetNode returns the node with the given name, if it exists.
// Otherwise, returns an error
func (no *NodeOperator) GetNode(ctx context.Context, app *meta.App) (*meta.Node, error) {
	kube := no.deployments()
	deployName := toDeploymentName(app)

	logger.Info("getting Node from k8s deployment",
		zap.String("deployment name", deployName))

	dep, err := kube.Get(deployName, metav1.GetOptions{})
	if err != nil {
		logger.Error("unable to find k8s deployment")
		return nil, ierrors.NewError().Message(err.Error()).Build()
	}

	node, err := toNode(dep)
	if err != nil {
		return nil, err
	}
	return node, nil
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
func NewOperator(memory memory.Manager) (nop *NodeOperator, err error) {
	sidp, err := strconv.Atoi(os.Getenv("INSPR_SIDECAR_PORT"))
	if err != nil {
		panic(err)
	}
	sidecarPort = int32(math.Min(float64(sidp), math.MaxInt32))
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
