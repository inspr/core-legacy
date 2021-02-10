package environment

import (
	"strings"

	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

// IsInInputChannel - checks if a channel exists in the insprEnv.InputChannels
func (insprEnv *InsprEnvironmentVariables) IsInInputChannel(channel string) bool {
	channelsList := insprEnv.GetInputChannelList()
	return utils.Includes(channelsList, channel)
}

// IsInOutputChannel - checks if a channel exists in the insprEnv.OutputChannels
func (insprEnv *InsprEnvironmentVariables) IsInOutputChannel(channel string) bool {
	channelsList := insprEnv.GetOutputChannelList()
	return utils.Includes(channelsList, channel)
}

// GetInputChannelList returns a string list with the channels in insprEnv.InputChannels
func (insprEnv *InsprEnvironmentVariables) GetInputChannelList() []string {
	channelList := strings.Split(insprEnv.InputChannels, ";")
	return channelList
}

// GetOutputChannelList returns a string list with the channels in insprEnv.OutputChannels
func (insprEnv *InsprEnvironmentVariables) GetOutputChannelList() []string {
	channelList := strings.Split(insprEnv.OutputChannels, ";")
	return channelList
}
