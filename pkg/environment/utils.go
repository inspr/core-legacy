package environment

import (
	"os"
	"strings"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

// IsInInputChannel - checks if a channel exists in the insprEnv.InputChannels
func (insprEnv *InsprEnvVars) IsInInputChannel(channel string) bool {
	channelsList := insprEnv.GetInputChannelList()
	return utils.Includes(channelsList, channel)
}

// IsInOutputChannel - checks if a channel exists in the insprEnv.OutputChannels
func (insprEnv *InsprEnvVars) IsInOutputChannel(channel string) bool {
	channelsList := insprEnv.GetOutputChannelList()
	return utils.Includes(channelsList, channel)
}

// GetInputChannelList returns a string list with the channels in insprEnv.InputChannels
func (insprEnv *InsprEnvVars) GetInputChannelList() []string {
	if insprEnv.InputChannels == "" {
		return []string{}
	}
	arr := strings.Split(insprEnv.InputChannels, ";")
	return arr[:len(arr)-1]
}

// GetOutputChannelList returns a string list with the channels in insprEnv.OutputChannels
func (insprEnv *InsprEnvVars) GetOutputChannelList() []string {
	if insprEnv.OutputChannels == "" {
		return []string{}
	}
	arr := strings.Split(insprEnv.OutputChannels, ";")
	return arr[:len(arr)-1]
}

// GetSchema returns a channel's schema, if the channel exists
func (insprEnv *InsprEnvVars) GetSchema(channel string) (string, error) {
	if insprEnv.IsInInputChannel(channel) || insprEnv.IsInOutputChannel(channel) {
		return os.Getenv(channel + "_SCHEMA"), nil
	}
	return "", ierrors.NewError().
		InvalidChannel().
		Message("channel " + channel + " not listed as an input or output").
		Build()
}
