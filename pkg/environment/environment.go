package environment

import "os"

// InsprEnvironmentVariables represents the current inspr environment
type InsprEnvironmentVariables struct {
	InputChannels    string
	OutputChannels   string
	UnixSocketAddr   string
	InsprAppContext  string
	InsprEnvironment string
}

var env *InsprEnvironmentVariables

// GetEnvironment returns the current inspr environment
func GetEnvironment() *InsprEnvironmentVariables {
	if env == nil {
		env = &InsprEnvironmentVariables{
			InputChannels:    getEnv("INSPR_INPUT_CHANNELS"),
			OutputChannels:   getEnv("INSPR_OUTPUT_CHANNELS"),
			UnixSocketAddr:   getEnv("INSPR_UNIX_SOCKET"),
			InsprAppContext:  getEnv("INSPR_APP_CTX"),
			InsprEnvironment: getEnv("INSPR_ENV"),
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
