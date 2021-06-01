package brokers

import "github.com/inspr/inspr/pkg/sidecars/models"

// GetSidecarConnectionVars returns port enviroment variable
// names for each possible broker
func GetSidecarConnectionVars(broker string) *models.ConnectionVariables {
	switch broker {
	case Kafka:
		return &models.ConnectionVariables{
			ReadEnvVar:  "INSPR_LBSIDECAR_READ_PORT",
			WriteEnvVar: "INSPR_SIDECAR_KAFKA_WRITE_PORT",
		}
	default:
		return nil
	}
}
