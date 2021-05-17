package tree

import (
	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta/brokers"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
)

// BrokersMemoryManager implements the methods described by the BrokersInterface
type BrokersMemoryManager struct {
}

// Brokers is a MemoryManager method that provides an access point for Alias
func (mm *MemoryManager) Brokers() memory.BrokerInterface {
	return &BrokersMemoryManager{}
}

var broker *brokers.Brokers

// GetAll returns an array containing all currently configured brokers
func (bmm *BrokersMemoryManager) GetAll() brokers.BrokerStatusArray {
	return brokers.BrokerStatusArray(bmm.get().Available.ToArray())
}

// GetDefault returns the broker configured as default
func (bmm *BrokersMemoryManager) GetDefault() brokers.BrokerStatus {
	return brokers.BrokerStatus(bmm.get().Default)
}

func (bmm *BrokersMemoryManager) get() *brokers.Brokers {
	if broker == nil {
		broker = &brokers.Brokers{
			Available: make(metautils.StrSet),
		}
	}
	return broker
}

// Create configures a new broker on insprd
func (bmm *BrokersMemoryManager) Create(broker brokers.BrokerStatus, config brokers.BrokerConfiguration) error {
	if ok := bmm.get().Available[string(broker)]; ok {
		return ierrors.NewError().Message("broker %s is already configured on memory", broker).Build()
	}
	//configure the sidecarFactory for the given broker
	//if succesful:
	bmm.get().Available[string(broker)] = true
	return nil
}

// SetDefault sets a previoulsy configured broker as insprd's default broker
func (bmm *BrokersMemoryManager) SetDefault(broker brokers.BrokerStatus) error {
	if ok := bmm.get().Available[string(broker)]; !ok {
		return ierrors.NewError().Message("broker %s is not configured on memory", broker).Build()
	}

	bmm.get().Default = (broker)
	return nil
}
