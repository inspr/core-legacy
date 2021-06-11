package brokers

import (
	"github.com/inspr/inspr/pkg/meta/brokers"
	"github.com/inspr/inspr/pkg/sidecars/models"
	"github.com/inspr/inspr/pkg/utils"
)

// Manager is the interface that allows the management
// of the system's message brokers.

// Manager is the interface that allows for interaction
// with the systems multiple brokers and its configurations.
type Manager interface {
	GetAll() (utils.StringArray, error)
	GetDefault() (string, error)
	Create(broker string, config brokers.BrokerConfiguration) error
	SetDefault(broker string) error
	Factory() SidecarManager
}

// SidecarManager is the interface that allows the build and deployment of
// available brokers
type SidecarManager interface {
	Get(broker string) (models.SidecarFactory, error)
	Subscribe(broker string, factory models.SidecarFactory) error
}
