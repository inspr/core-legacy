package brokers

import (
	"github.com/inspr/inspr/pkg/meta/brokers"
	"github.com/inspr/inspr/pkg/sidecar_old/models"
)

// Manager is the interface that allows the management
// of the system's message brokers.

// Manager is the interface that allows for interaction
// with the systems multiple brokers and its configurations.
type Manager interface {
	GetAll() brokers.BrokerStatusArray
	GetDefault() brokers.BrokerStatus
	Create(broker brokers.BrokerStatus, config brokers.BrokerConfiguration) error
	SetDefault(broker brokers.BrokerStatus) error
	Factory() SidecarManager
}

// SidecarManager is the interface that allows the build and deployment of
// available brokers
type SidecarManager interface {
	Get(broker string) (models.SidecarFactory, error)
	Subscribe(broker string, factory models.SidecarFactory) error
}
