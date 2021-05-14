package brokers

import "github.com/inspr/inspr/pkg/meta/brokers"

type Manager interface {
	Brokers() BrokerInterface
}

// BrokerInterface is the interface tht allows for interaction
// with the systems multiple brokers
type BrokerInterface interface {
	GetAll() brokers.BrokerStatusArray
	GetDefault() brokers.BrokerStatus
	Create(broker brokers.BrokerStatus, config brokers.BrokerConfiguration) error
	SetDefault(broker brokers.BrokerStatus) error
}
