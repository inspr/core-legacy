package environment

import "os"

// InsprEnvironment represents the current inspr environment
type InsprEnvironment struct {
	NodeID                string
	Namespace             string
	RegistryURL           string
	InputChannels         string
	OutputChannels        string
	LogChannel            string
	Environment           string
	Domain                string
	AppsSubdomain         string
	AppsNamespace         string
	Subdomain             string
	AppsTLS               bool
	KafkaBootstrapServers string
	ChannelSidecarImage   string
	KafkaSidecarImage     string
}

var env *InsprEnvironment

// GetEnvironment returns the current inspr environment
func GetEnvironment() *InsprEnvironment {
	if env == nil {
		env = &InsprEnvironment{
			NodeID:                getEnv("INSPR_NODE_ID"),
			Namespace:             getEnv("INSPR_NAMESPACE"),
			LogChannel:            getEnv("INSPR_LOG_CHANNEL"),
			RegistryURL:           getEnv("INSPR_REGISTRY_URL"),
			InputChannels:         getEnv("INSPR_INPUT_CHANNELS"),
			OutputChannels:        getEnv("INSPR_OUTPUT_CHANNELS"),
			Environment:           getEnv("INSPR_ENVIRONMENT"),
			Domain:                getEnv("INSPR_DOMAIN"),
			AppsSubdomain:         getEnv("INSPR_APPS_SUBDOMAIN"),
			AppsNamespace:         getEnv("INSPR_APPS_NAMESPACE"),
			KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS"),
			Subdomain:             getEnv("INSPR_SUBDOMAIN"),
			ChannelSidecarImage:   getEnv("INSPR_CHANNEL_SIDECAR"),
			KafkaSidecarImage:     getEnv("INSPR_KAFKA_SIDECAR_IMAGE"),
			AppsTLS: func() bool {
				ret := getEnv("INSPR_APPS_TLS")
				if ret == "true" {
					return true
				}
				return false
			}(),
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
