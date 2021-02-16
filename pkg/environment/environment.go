package environment

import "os"

// InsprEnvVars represents the current inspr environment
type InsprEnvVars struct {
	InputChannels    string
	OutputChannels   string
	UnixSocketAddr   string
	SidecarImage     string
	InsprAppContext  string
	InsprEnvironment string
}

var env *InsprEnvVars

// GetEnvironment returns the current inspr environment
func GetEnvironment() *InsprEnvVars {
	if env == nil {
		env = &InsprEnvVars{
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

// RefreshEnviromentVariables "refreshes" the value of inspr environment variables.
// This was develop for testing and probably sholdn't be used in other cases.
func RefreshEnviromentVariables() *InsprEnvVars {
	env = &InsprEnvVars{
		InputChannels:    getEnv("INSPR_INPUT_CHANNELS"),
		OutputChannels:   getEnv("INSPR_OUTPUT_CHANNELS"),
		UnixSocketAddr:   getEnv("INSPR_UNIX_SOCKET"),
		InsprAppContext:  getEnv("INSPR_APP_CTX"),
		InsprEnvironment: getEnv("INSPR_ENV"),
	}
	return env
}
