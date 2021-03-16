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
func GetInputChannelList(inputChan string) utils.StringArray {
	if inputChan == "" {
		return utils.StringArray{}
	}
	arr := strings.Split(inputChan, ";")
	return arr
}

// GetResolvedInputChannelList gets the list of resolved channels from the input boundary
func GetResolvedInputChannelList(inputChan string) utils.StringArray {
	arr := utils.StringArray(GetInputChannelList(inputChan))
	return arr.Map(func(s string) string {
		resolved, _ := GetResolvedChannel(s, inputChan, "")
		return resolved
	})
}

// GetResolvedOutputChannelList gets the list of resolved channels from the output boundary
func GetResolvedOutputChannelList(outputChan string) utils.StringArray {
	arr := utils.StringArray(GetOutputChannelList(outputChan))
	return arr.Map(func(s string) string {
		resolved, _ := GetResolvedChannel(s, outputChan, "")
		return resolved
	})
}

// GetOutputChannelList returns a string list with the channels in insprEnv.OutputChannels
func GetOutputChannelList(outputChan string) utils.StringArray {
	if outputChan == "" {
		return utils.StringArray{}
	}
	arr := strings.Split(outputChan, ";")
	return arr
}

// GetSchema returns a channel's schema, if the channel exists
func GetSchema(channel string) (string, error) {
	schema, ok := os.LookupEnv(channel + "_SCHEMA")
	if !ok {
		return "", ierrors.NewError().NotFound().Message("schema for channel %s not found", channel).Build()
	}
	return schema, nil
}

// GetResolvedChannel gets a resolved channel from a channel name
func GetResolvedChannel(channel, inputChan, outputChan string) (string, error) {
	if IsInInputChannel(channel, inputChan) || IsInOutputChannel(channel, outputChan) {
		return os.Getenv(channel + "_RESOLVED"), nil
	}
	return "", ierrors.NewError().
		InvalidChannel().
		Message("channel " + channel + " not listed as an input or output").
		Build()
}
