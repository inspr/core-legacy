package kafka

import "os"

// Environment represents the current inspr environment
type Environment struct {
	KafkaBootstrapServers string
}

var env *Environment

// GetEnvironment returns the current inspr environment
func GetEnvironment() *Environment {
	if env == nil {
		env = &Environment{
			KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS"),
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
