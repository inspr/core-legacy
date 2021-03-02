package environment

import (
	"os"
	"strings"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

// IsInInputChannel - checks if a channel exists in the insprEnv.InputChannels
func IsInInputChannel(channel, inputChan string) bool {
	channelsList := GetInputChannelList(inputChan)
	return utils.Includes(channelsList, channel)
}

// IsInOutputChannel - checks if a channel exists in the insprEnv.OutputChannels
func IsInOutputChannel(channel, outputChan string) bool {
	channelsList := GetOutputChannelList(outputChan)
	return utils.Includes(channelsList, channel)
}

// GetInputChannelList returns a string list with the channels in insprEnv.InputChannels
func GetInputChannelList(inputChan string) []string {
	if inputChan == "" {
		return []string{}
	}
	arr := strings.Split(inputChan, ";")
	return arr
}

// GetOutputChannelList returns a string list with the channels in insprEnv.OutputChannels
func GetOutputChannelList(outputChan string) []string {
	if outputChan == "" {
		return []string{}
	}
	arr := strings.Split(outputChan, ";")
	return arr
}

// GetSchema returns a channel's schema, if the channel exists
func GetSchema(channel, inputChan, outputChan string) (string, error) {
	if IsInInputChannel(channel, inputChan) || IsInOutputChannel(channel, outputChan) {
		return os.Getenv(channel + "_SCHEMA"), nil
	}
	return "", ierrors.NewError().
		InvalidChannel().
		Message("channel " + channel + " not listed as an input or output").
		Build()
}
