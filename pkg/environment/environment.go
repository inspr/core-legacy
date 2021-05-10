package environment

import (
	"os"

	"github.com/inspr/inspr/pkg/ierrors"
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
			InputChannels:    GetInputChannels(),
			OutputChannels:   GetOutputChannels(),
			SidecarImage:     GetSidecarImage(),
			InsprAppContext:  GetInsprAppContext(),
			InsprEnvironment: GetInsprEnvironment(),
			InsprAppID:       GetInsprAppID(),
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
		InputChannels:    GetInputChannels(),
		OutputChannels:   GetOutputChannels(),
		SidecarImage:     GetSidecarImage(),
		InsprAppContext:  GetInsprAppContext(),
		InsprEnvironment: GetInsprEnvironment(),
		InsprAppID:       GetInsprAppID(),
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

// GetInputChannels returns environment variable which contains the input channels
func GetInputChannels() string {
	return getEnv("INSPR_INPUT_CHANNELS")
}

// GetOutputChannels returns environment variable which contains the output channels
func GetOutputChannels() string {
	return getEnv("INSPR_OUTPUT_CHANNELS")
}

// GetUnixSocketAddress returns environment variable which contains the unix socket address
func GetUnixSocketAddress() string {
	return getEnv("INSPR_UNIX_SOCKET")
}

// GetSidecarImage returns environment variable which contains the sidecar image reference
func GetSidecarImage() string {
	return getEnv("INSPR_SIDECAR_IMAGE")
}

// GetInsprAppContext returns environment variable which contains the current dApp context
func GetInsprAppContext() string {
	return getEnv("INSPR_APP_CTX")
}

// GetInsprEnvironment returns Inspr's current environment (test, production, qa, etc.)
func GetInsprEnvironment() string {
	return getEnv("INSPR_ENV")
}

// GetInsprAppID returns environment variable which contains the current dApp's ID
func GetInsprAppID() string {
	return getEnv("INSPR_APP_ID")
}
