package brokers

import (
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/sidecars/models"
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
	return ierrors.New("%s broker already subscribed", broker)
}

// Get returns a factory for the specifyed broker
func (abf *AbstractBrokerFactory) Get(broker string) (models.SidecarFactory, error) {
	if factories == nil {
		return nil, ierrors.New("no brokers are allowed")
	}
	if factory, ok := factories[broker]; ok {
		return factory, nil
	}
	return nil, ierrors.New("%s broker not allowed", broker)
}
