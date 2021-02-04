package environment

import "os"

// InsprEnvironment represents the current inspr environment
type InsprEnvironment struct {
	InputChannels  string
	OutputChannels string
	UnixSocketAddr string
}

var env *InsprEnvironment

// GetEnvironment returns the current inspr environment
func GetEnvironment() *InsprEnvironment {
	if env == nil {
		env = &InsprEnvironment{
			InputChannels:  getEnv("INSPR_INPUT_CHANNELS"),
			OutputChannels: getEnv("INSPR_OUTPUT_CHANNELS"),
			UnixSocketAddr: getEnv("INSPR_UNIX_SOCKET"),
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
