package environment

import (
	"fmt"
	"os"
	"strings"

	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta/brokers"
	"github.com/inspr/inspr/pkg/utils"
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
			InputChannels:    getRawInputChannels(),
			OutputChannels:   getRawInputChannels(),
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
		InputChannels:    getRawInputChannels(),
		OutputChannels:   getRawOutputChannels(),
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

// GetInputChannels returns the input channels
func GetInputChannels() []brokers.ChannelBroker {
	return getChannelList(getRawInputChannels())
}

// GetOutputChannels returns the output channels
func GetOutputChannels() []brokers.ChannelBroker {
	return getChannelList(getRawOutputChannels())
}

// GetInputBrokerChannels returns environment variable which contains the input channels
func GetInputBrokerChannels(broker string) utils.StringArray {
	channels := getChannelList(getRawInputChannels())
	return filterChannelsByBroker(broker, channels)
}

// GetOutputBrokerChannels returns environment variable which contains the output channels
func GetOutputBrokerChannels(broker string) utils.StringArray {
	channels := getChannelList(getRawOutputChannels())
	return filterChannelsByBroker(broker, channels)
}

// getRawInputChannels returns environment variable which contains the input channels
func getRawInputChannels() string {
	return getEnv("INSPR_INPUT_CHANNELS")
}

// getRawOutputChannels returns environment variable which contains the output channels
func getRawOutputChannels() string {
	return getEnv("INSPR_OUTPUT_CHANNELS")
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

// GetBrokerWritePort returns environment variable that contains given broker's write port
func GetBrokerWritePort(broker string) string {
	return getEnv(fmt.Sprintf("INSPR_SIDECAR_%s_WRITE_PORT", strings.ToUpper(broker)))
}

// GetBrokerReadPort returns environment variable that contains given broker's read port
func GetBrokerReadPort(broker string) string {
	return getEnv(fmt.Sprintf("INSPR_SIDECAR_%s_READ_PORT", strings.ToUpper(broker)))
}

// GetBrokerSpecificSidecarAddr returns environment variable that contains given broker's
// sidecar address
func GetBrokerSpecificSidecarAddr(broker string) string {
	return getEnv(fmt.Sprintf("INSPR_SIDECAR_%s_ADDR", strings.ToUpper(broker)))
}
