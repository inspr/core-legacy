package brokers

import (
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta/brokers"
)

// GetAll returns an array containing all currently configured brokers
func (bmm *BrokerMemoryManager) GetAll() brokers.BrokerStatusArray {
	return brokers.BrokerStatusArray(bmm.get().Available.ToArray())
}

// GetDefault returns the broker configured as default
func (bmm *BrokerMemoryManager) GetDefault() brokers.BrokerStatus {
	return brokers.BrokerStatus(bmm.get().Default)
}

func (bmm *BrokerMemoryManager) get() *brokers.Brokers {
	if bmm.broker == nil {
		panic("broker status memory unavailible")
	}
	return bmm.broker
}

// Create configures a new broker on insprd
func (bmm *BrokerMemoryManager) Create(broker brokers.BrokerStatus, config brokers.BrokerConfiguration) error {
	if ok := bmm.get().Available[string(broker)]; ok {
		return ierrors.NewError().Message("broker %s is already configured on memory", broker).Build()
	}
	//configure the sidecarFactory for the given broker
	//if succesful:
	bmm.get().Available[string(broker)] = true
	return nil
}

// SetDefault sets a previoulsy configured broker as insprd's default broker
func (bmm *BrokerMemoryManager) SetDefault(broker brokers.BrokerStatus) error {
	if ok := bmm.get().Available[string(broker)]; !ok {
		return ierrors.NewError().Message("broker %s is not configured on memory", broker).Build()
	}

	bmm.get().Default = (broker)
	return nil
}

// Factory provides the struct implementation for Sidecarfactory
func (bmm *BrokerMemoryManager) Factory() SidecarManager {
	return bmm.factory
}
