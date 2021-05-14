package tree

import (
	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta/brokers"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
	"github.com/inspr/inspr/pkg/utils"
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
func (bmm *BrokersMemoryManager) GetAll() utils.StringArray {
	return bmm.get().Availible.ToArray()
}

// GetDefault returns the broker configured as default
func (bmm *BrokersMemoryManager) GetDefault() string {
	return bmm.get().Default
}

func (bmm *BrokersMemoryManager) get() *brokers.Brokers {
	if broker == nil {
		broker = &brokers.Brokers{
			Availible: make(metautils.StrSet),
		}
	}
	return broker
}

// Create configures a new broker on insprd
func (bmm *BrokersMemoryManager) Create(broker string, config interface{}) error {
	if ok := bmm.get().Availible[broker]; ok {
		return ierrors.NewError().Message("broker %s is already configured on memory", broker).Build()
	}
	//configure the sidecarFactory for the given broker
	//if succesful:
	bmm.get().Availible[broker] = true
	return nil
}

// SetDefault sets a previoulsy configured broker as insprd's default broker
func (bmm *BrokersMemoryManager) SetDefault(broker string) error {
	if ok := bmm.get().Availible[broker]; !ok {
		return ierrors.NewError().Message("broker %s is not configured on memory", broker).Build()
	}

	bmm.get().Default = broker
	return nil
}
