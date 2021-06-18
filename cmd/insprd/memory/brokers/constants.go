package brokers

import (
	"inspr.dev/inspr/pkg/meta/brokers"
	"inspr.dev/inspr/pkg/sidecars/models"
)

// GetSidecarConnectionVars returns port enviroment variable
// names for each possible broker
func GetSidecarConnectionVars(broker string) *models.ConnectionVariables {
	switch broker {
	case brokers.Kafka:
		return &models.ConnectionVariables{
			ReadEnvVar:  "INSPR_LBSIDECAR_READ_PORT",
			WriteEnvVar: "INSPR_SIDECAR_KAFKA_WRITE_PORT",
		}
	default:
		return nil
	}
}
