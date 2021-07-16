package kafkasc

import (
	"os"

	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/logs"
)

// Environment represents the current inspr environment
type Environment struct {
	KafkaBootstrapServers string
	KafkaAutoOffsetReset  string
}

var env *Environment
var logger *zap.Logger

// init is called after all the variable declarations in the package have evaluated
// their initializers, and those are evaluated only after all the imported packages
// have been initialized
func init() {
	logger, _ = logs.Logger(zap.Fields(zap.String("section", "kafka-sidecar")))
}

// GetKafkaEnvironment returns the current inspr environment
func GetKafkaEnvironment() *Environment {
	if env == nil {
		env = &Environment{
			KafkaBootstrapServers: getEnv("INSPR_SIDECAR_KAFKA_BOOTSTRAP_SERVERS"),
			KafkaAutoOffsetReset:  getEnv("INSPR_SIDECAR_KAFKA_AUTO_OFFSET_RESET"),
		}
	}
	return env
}

func getEnv(name string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	panic("[ENV VAR] " + name + " not found")
}

// RefreshEnviromentVariables "refreshes" the value of kafka environment variables.
// This was develop for testing and probably sholdn't be used in other cases.
func RefreshEnviromentVariables() *Environment {
	env = &Environment{
		KafkaBootstrapServers: getEnv("INSPR_SIDECAR_KAFKA_BOOTSTRAP_SERVERS"),
		KafkaAutoOffsetReset:  getEnv("INSPR_SIDECAR_KAFKA_AUTO_OFFSET_RESET"),
	}
	return env
}
