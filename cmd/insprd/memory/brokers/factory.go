package brokers

import (
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/operator/k8s"
	"github.com/inspr/inspr/pkg/sidecar_old/models"
)

type brokerFactory map[string]models.SidecarFactory

var factories brokerFactory

// AbstractBrokerFactory singleton and abstract factory
// implementation for SidecarInterface
type AbstractBrokerFactory struct {
}

// Subscribe includes a broker specific factory on the Abstract broker factory
func (abf *AbstractBrokerFactory) Subscribe(broker string, factory models.SidecarFactory) error {
	if factories == nil {
		factories = make(brokerFactory)
	}
	if _, ok := factories[broker]; !ok {
		factories[broker] = factory
		return nil
	}
	return ierrors.NewError().Message("%s broker already subscribed", broker).Build()
}

// Get returns a factory for the specifyed broker
func (abf *AbstractBrokerFactory) Get(broker string) (models.SidecarFactory, error) {
	if factories == nil {
		return nil, ierrors.NewError().Message("no brokers are allowed").Build()
	}
	if factory, ok := factories[broker]; ok {
		return factory, nil
	}
	return nil, ierrors.NewError().Message("%s broker not allowed", broker).Build()
}

func getAvailiblePorts() *models.SidecarConnections {
	return nil
}

func getAllSidecarNames(app *meta.App) []string {
	return nil
}

func getAllSidecarsDeployments(app *meta.App) []k8s.DeploymentOption {
	var ret []k8s.DeploymentOption
	for _, broker := range getAllSidecarNames(app) {
		if option, ok := factories[broker]; ok {
			ret = append(ret, option(app, getAvailiblePorts()))
		} else {
			panic("broker not allowed")
		}
	}
	return ret
}
