package kafka

import "os"

// KafkaEnvironment represents the current inspr environment
type KafkaEnvironment struct {
	KafkaBootstrapServers string
}

var env *KafkaEnvironment

// GetEnvironment returns the current inspr environment
func GetEnvironment() *KafkaEnvironment {
	if env == nil {
		env = &KafkaEnvironment{
			KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS"),
		}
	}
	return env
}

func getEnv(name string) string {
	if value, exists := os.LookupEnv(name); exists {
		return value
	}
	panic("[ENV VAR] " + name + " not found.")
}
