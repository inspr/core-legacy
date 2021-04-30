package nodes

import (
	"context"
	"math"
	"os"
	"strconv"

	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/meta/utils"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"

	kubeApp "k8s.io/api/apps/v1"
	kubeCore "k8s.io/api/core/v1"
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
}

func (no *NodeOperator) services() cv1.ServiceInterface {
	appsNamespace := getK8SVariables().AppsNamespace
	return no.clientSet.CoreV1().Services(appsNamespace)
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

	logger.Debug("converting dApp to the k8s deployment to be created")
	var deploy *kubeApp.Deployment
	var svc *kubeCore.Service

	kube := no.retrieveKube()
	services := no.services()

	deploy = no.dAppToDeployment(app)
	svc = dappToService(app)
	logger.Debug("creating the k8s deployment")
	dep, err := kube.Create(deploy)
	if err != nil {
		logger.Error("unable to create the k8s deployment")
		return nil, ierrors.NewError().Message(err.Error()).Build()
	}

	_, err = services.Create(svc)
	if err != nil {
		logger.Error("unable to create the k8s service", zap.Any("error", err))
		return nil, ierrors.NewError().InnerError(err).Message("unable to create kubernetes service").Build()
	}
	node, err := toNode(dep)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// UpdateNode updates a node that already exists, if the new structure is valid.
// Otherwise, returns an error.
func (no *NodeOperator) UpdateNode(ctx context.Context, app *meta.App) (*meta.Node, error) {
	logger.Info("updating a Node structure in k8s",
		zap.Any("node", app))

	logger.Debug("converting dApp to the k8s deployment to be updated")
	var deploy *kubeApp.Deployment
	var svc *kubeCore.Service

	kube := no.retrieveKube()
	services := no.services()

	deploy = no.dAppToDeployment(app)
	svc = dappToService(app)

	logger.Debug("updating the k8s deployment")
	dep, err := kube.Update(deploy)
	if err != nil {
		logger.Error("unable to update the k8s deployment")
		return nil, ierrors.NewError().Message(err.Error()).Build()
	}
	_, err = services.Update(svc)
	if err != nil {
		logger.Error("unable to update the k8s service")
		return nil, ierrors.NewError().InnerError(err).Message("unable to update kubernetes service").Build()
	}

	node, err := toNode(dep)
	if err != nil {
		return nil, err
	}
	return node, nil
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
	svcs := no.services()
	err = svcs.Delete(deployName, &metav1.DeleteOptions{})

	if err != nil {
		logger.Error("unable to delete the k8s service")
		return ierrors.NewError().Message(err.Error()).Build()
	}
	return nil
}

// NewOperator initializes a k8s based kafka node operator with in cluster configuration
func NewOperator(memory memory.Manager) (nop *NodeOperator, err error) {
	sidp, err := strconv.Atoi(os.Getenv("INSPR_SIDECAR_PORT"))
	if err != nil {
		if err != nil {
			panic(err)
		}

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
