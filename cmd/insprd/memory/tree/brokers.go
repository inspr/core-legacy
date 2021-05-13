package tree

import (
	"github.com/inspr/inspr/cmd/insprd/memory"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta/brokers"
	metautils "github.com/inspr/inspr/pkg/meta/utils"
	"github.com/inspr/inspr/pkg/utils"
)

type BrokersMemoryManager struct {
}

// Brokers is a MemoryManager method that provides an access point for Alias
func (tmm *MemoryManager) Brokers() memory.BrokerInterface {
	return &BrokersMemoryManager{}
}

var bro *brokers.Brokers

func (bmm *BrokersMemoryManager) GetAll() utils.StringArray {
	return bmm.get().Availible.ToArray()
}

func (bmm *BrokersMemoryManager) GetDefault() string {
	return bmm.get().Default
}

func (bmm *BrokersMemoryManager) get() *brokers.Brokers {
	if bro == nil {
		bro = &brokers.Brokers{
			Availible: make(metautils.StrSet),
		}
	}
	return bro
}

func (bmm *BrokersMemoryManager) Create(broker string, config interface{}) error {
	if ok := bmm.get().Availible[broker]; ok {
		return ierrors.NewError().Message("error: %s is already configured on memory", broker).Build()
	}
	//configure the sidecarFactory for the given broker
	//if succesful:
	bmm.get().Availible[broker] = true
	return nil
}

func (bmm *BrokersMemoryManager) SetDefault(broker string) error {
	if ok := bmm.get().Availible[broker]; !ok {
		return ierrors.NewError().Message("error: %s is not configured on memory", broker).Build()
	}

	bmm.get().Default = broker
	return nil
}
