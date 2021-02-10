package environment

import "os"

// InsprEnvironment represents the current inspr environment
type InsprEnvironment struct {
	InputChannels    string
	OutputChannels   string
	UnixSocketAddr   string
	SidecarImage     string
	InsprAppContext  string
	InsprEnvironment string
}

var env *InsprEnvironment

// GetEnvironment returns the current inspr environment
func GetEnvironment() *InsprEnvironment {
	if env == nil {
		env = &InsprEnvironment{
			InputChannels:    getEnv("INSPR_INPUT_CHANNELS"),
			OutputChannels:   getEnv("INSPR_OUTPUT_CHANNELS"),
			UnixSocketAddr:   getEnv("INSPR_UNIX_SOCKET"),
			SidecarImage:     getEnv("INSPR_SIDECAR_IMAGE"),
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
