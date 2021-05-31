package models

import "github.com/inspr/inspr/cmd/insprd/memory/brokers"

// GetSidecarConnectionVars returns port enviroment variable
// names for each possible broker
func GetSidecarConnectionVars(broker string) *ConnectionVariables {
	switch broker {
	case brokers.Kafka:
		return &ConnectionVariables{
			ReadEnvVar:  "INSPR_LBSIDECAR_READ_PORT",
			WriteEnvVar: "INSPR_SIDECAR_KAFKA_WRITE_PORT",
		}
	default:
		return nil
	}
}
