package brokers

import (
	"github.com/inspr/inspr/cmd/sidecars"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta/brokers"
	"github.com/inspr/inspr/pkg/sidecars/models"
)

// GetAll returns an array containing all currently configured brokers
func (bmm *BrokerMemoryManager) GetAll() (brokers.BrokerStatusArray, error) {
	mem, err := bmm.get()
	if err != nil {
		return nil, err
	}
	return brokers.BrokerStatusArray(mem.Available.ToArray()), nil
}

// GetDefault returns the broker configured as default
func (bmm *BrokerMemoryManager) GetDefault() (*brokers.BrokerStatus, error) {
	mem, err := bmm.get()
	if err != nil {
		return nil, err
	}
	var status brokers.BrokerStatus = mem.Default
	return &status, nil
}

func (bmm *BrokerMemoryManager) get() (*brokers.Brokers, error) {
	if bmm.broker == nil {
		return nil, ierrors.NewError().Message("broker status memory is empty").Build()
	}
	return bmm.broker, nil
}

// Create configures a new broker on insprd
func (bmm *BrokerMemoryManager) Create(broker brokers.BrokerStatus, config brokers.BrokerConfiguration) error {
	mem, err := bmm.get()
	if err != nil {
		return err
	}

	if ok := mem.Available[string(broker)]; ok {
		return ierrors.NewError().Message("broker %s is already configured on memory", broker).Build()
	}

	var factory models.SidecarFactory
	switch string(broker) {
	case brokers.Kafka:
		factory = sidecars.KafkaToDeployment(config.(sidecars.KafkaConfig))
	default:
		return ierrors.NewError().Message("broker %s is not valid", broker).Build()
	}

	err = bmm.Factory().Subscribe(string(broker), factory)

	if err != nil {
		return err
	}

	mem.Available[string(broker)] = true
	return nil
}

// SetDefault sets a previously configured broker as insprd's default broker
func (bmm *BrokerMemoryManager) SetDefault(broker brokers.BrokerStatus) error {
	mem, err := bmm.get()
	if err != nil {
		return err
	}

	if ok := mem.Available[string(broker)]; !ok {
		return ierrors.NewError().Message("broker %s is not configured on memory", broker).Build()
	}

	mem.Default = (broker)
	return nil
}

// Factory provides the struct implementation for Sidecarfactory
func (bmm *BrokerMemoryManager) Factory() SidecarManager {
	return bmm.factory
}
