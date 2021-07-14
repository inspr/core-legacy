package brokers

import (
	apimodels "inspr.dev/inspr/pkg/api/models"
	"inspr.dev/inspr/pkg/meta/brokers"
	"inspr.dev/inspr/pkg/sidecars/models"
)

// Manager is the interface that allows the management
// of the system's message brokers.

// Manager is the interface that allows for interaction
// with the systems multiple brokers and its configurations.
type Manager interface {
	Get() (*apimodels.BrokersDI, error)
	Create(config brokers.BrokerConfiguration) error
	SetDefault(broker string) error
	Factory() SidecarManager
	Configs(broker string) (brokers.BrokerConfiguration, error)
}

// SidecarManager is the interface that allows the build and deployment of
// available brokers
type SidecarManager interface {
	Get(broker string) (models.SidecarFactory, error)
	Subscribe(broker string, factory models.SidecarFactory) error
}
