package brokers

import (
	"github.com/inspr/inspr/cmd/sidecars"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta/brokers"
	"github.com/inspr/inspr/pkg/sidecars/models"
	"github.com/inspr/inspr/pkg/utils"
)

// GetAll returns an array containing all currently configured brokers
func (bmm *BrokerMemoryManager) GetAll() (utils.StringArray, error) {
	mem, err := bmm.get()
	if err != nil {
		return nil, err
	}
	return mem.Available.Brokers(), nil
}

// GetDefault returns the broker configured as default
func (bmm *BrokerMemoryManager) GetDefault() (string, error) {
	mem, err := bmm.get()
	if err != nil {
		return "", err
	}
	return mem.Default, nil
}

func (bmm *BrokerMemoryManager) get() (*brokers.Brokers, error) {
	if bmm.broker == nil {
		return nil, ierrors.NewError().Message("broker status memory is empty").Build()
	}
	return bmm.broker, nil
}

// Create configures a new broker on insprd
func (bmm *BrokerMemoryManager) Create(broker string, config brokers.BrokerConfiguration) error {
	mem, err := bmm.get()
	if err != nil {
		return err
	}

	if _, ok := mem.Available[string(broker)]; ok {
		return ierrors.NewError().Message("broker %s is already configured on memory", broker).Build()
	}

	var factory models.SidecarFactory
	switch string(broker) {
	case brokers.Kafka:
		// bmm.
		factory = sidecars.KafkaToDeployment(config.(sidecars.KafkaConfig))
	default:
		return ierrors.NewError().Message("broker %s is not supported", broker).Build()
	}

	err = bmm.Factory().Subscribe(string(broker), factory)

	if err != nil {
		return err
	}

	mem.Available[string(broker)] = config
	return nil
}

// SetDefault sets a previously configured broker as insprd's default broker
func (bmm *BrokerMemoryManager) SetDefault(broker string) error {
	mem, err := bmm.get()
	if err != nil {
		return err
	}

	if _, ok := mem.Available[string(broker)]; !ok {
		return ierrors.NewError().Message("broker %s is not configured on memory", broker).Build()
	}

	mem.Default = (broker)
	return nil
}

// Factory provides the struct implementation for Sidecarfactory
func (bmm *BrokerMemoryManager) Factory() SidecarManager {
	return bmm.factory
}
