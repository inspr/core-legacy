package brokers

import (
	"github.com/inspr/inspr/pkg/meta/brokers"
	"github.com/inspr/inspr/pkg/sidecar/models"
)

// Manager is the interface that allows the management
// of the system's message brokers.
type Manager interface {
	Brokers() BrokerInterface
}

// BrokerInterface is the interface that allows for interaction
// with the systems multiple brokers and its conficurations.
type BrokerInterface interface {
	GetAll() brokers.BrokerStatusArray
	GetDefault() brokers.BrokerStatus
	Create(broker brokers.BrokerStatus, config brokers.BrokerConfiguration) error
	SetDefault(broker brokers.BrokerStatus) error
	Factory() SidecarInterface
}

// SidecarInterface is the interface that allows the build and deployment of
// available brokers
type SidecarInterface interface {
	Get(broker string) (models.SidecarFactory, error)
	Subscribe(broker string, factory models.SidecarFactory) error
}
