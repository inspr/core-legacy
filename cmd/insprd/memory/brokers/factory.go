package brokers

import (
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/operator/k8s"
	"github.com/inspr/inspr/pkg/sidecar/models"
)

type brokerFactory map[string]models.SidecarFactory

var f brokerFactory

// AbstractBrokerFactory singleton and abstract factory
// implementation for SidecarInterface
type AbstractBrokerFactory struct {
}

// Subscribe includes a broker specific factory on the Abstract broker factory
func (abf *AbstractBrokerFactory) Subscribe(broker string, factory models.SidecarFactory) error {
	if f == nil {
		f = make(brokerFactory)
	}
	if _, ok := f[broker]; !ok {
		f[broker] = factory
		return nil
	}
	return ierrors.NewError().Message("%s broker already subscribed", broker).Build()
}

// Get returns a factory for the specifyed broker
func (abf *AbstractBrokerFactory) Get(broker string) (models.SidecarFactory, error) {
	if f == nil {
		return nil, ierrors.NewError().Message("no brokers are allowed").Build()
	}
	if factory, ok := f[broker]; ok {
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
		if option, ok := f[broker]; ok {
			ret = append(ret, option(app, getAvailiblePorts()))
		} else {
			panic("broker not allowed")
		}
	}
	return ret
}
