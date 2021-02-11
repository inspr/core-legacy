package kafka

import "os"

// Environment represents the current inspr environment
type Environment struct {
	KafkaBootstrapServers string
	KafkaAutoOffsetReset  string
}

var env *Environment

// GetEnvironment returns the current inspr environment
func GetEnvironment() *Environment {
	if env == nil {
		env = &Environment{
			KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS"),
			KafkaAutoOffsetReset:  getEnv("KAFKA_AUTO_OFFSET_RESET"),
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
		KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS"),
		KafkaAutoOffsetReset:  getEnv("KAFKA_AUTO_OFFSET_RESET"),
	}
	return env
}
