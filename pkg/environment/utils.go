package environment

import (
	"os"
	"strings"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

// IsInBoundaryChannel - checks if a channel exists in the insprEnv.OutputChannels
func IsInChannelBoundary(channel, outputChan string) bool {
	channelsList := GetChannelBoundaryList(outputChan)
	return utils.Includes(channelsList, channel)
}

// GetChannelList returns a string list with the channels in insprEnv.OutputChannels
func GetChannelBoundaryList(channels string) []string {
	if channels == "" {
		return []string{}
	}
	arr := strings.Split(channels, ";")
	return arr
}

// GetSchema returns a channel's schema, if the channel exists
func GetSchema(channel, inputChan, outputChan string) (string, error) {
	if IsInChannelBoundary(channel, inputChan) || IsInChannelBoundary(channel, outputChan) {
		return os.Getenv(channel + "_SCHEMA"), nil
	}
	return "", ierrors.NewError().
		InvalidChannel().
		Message("channel " + channel + " not listed as an input or output").
		Build()
}
