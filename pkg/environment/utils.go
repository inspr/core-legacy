package environment

import (
	"os"
	"strings"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

// IsInChannelBoundary - checks if a channel exists in the insprEnv.OutputChannels
func IsInChannelBoundary(channel, outputChan string) bool {
	channelsList := GetChannelBoundaryList(outputChan)
	return utils.Includes(channelsList, channel)
}

// GetChannelBoundaryList returns a string list with the channels in insprEnv.OutputChannels
func GetChannelBoundaryList(channels string) []string {
	if channels == "" {
		return []string{}

// GetResolvedBoundaryChannelList gets the list of resolved channels from the input boundary
func GetResolvedBoundaryChannelList(channels string) utils.StringArray {
	arr := utils.StringArray(GetChannelBoundaryList(channels))
	return arr.Map(func(s string) string {
		resolved, _ := GetResolvedChannel(s, channels, "")
		return resolved
	})
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
