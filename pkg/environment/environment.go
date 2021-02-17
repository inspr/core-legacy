package environment

import (
	"os"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
)

// InsprEnvVars represents the current inspr environment
type InsprEnvVars struct {
	InputChannels    string
	OutputChannels   string
	UnixSocketAddr   string
	SidecarImage     string
	InsprAppContext  string
	InsprEnvironment string
	InsprAppID       string
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
			InsprAppID:       getEnv("INSPR_APP_ID"),
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
		SidecarImage:     getEnv("INSPR_SIDECAR_IMAGE"),
		InsprAppContext:  getEnv("INSPR_APP_CTX"),
		InsprEnvironment: getEnv("INSPR_ENV"),
		InsprAppID:       getEnv("INSPR_APP_ID"),
	}
	return env
}

// RecoverEnvironmentErrors recovers environment errors when instantiating the environment. It sends any recovered
// errors to the channel in the parameter of the function.
//
// The channel that is passed through to the parameter must have at least one buffer spot, so that the error can be easily consumed.
func RecoverEnvironmentErrors(errch chan<- error) {
	err := recover()
	if err != nil {
		errch <- ierrors.NewError().Message(err.(string)).Build()
	}
	errch <- nil
}
