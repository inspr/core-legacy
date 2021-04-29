package environment

import (
	"os"
	"strings"

	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/utils"
)

// IsInChannelBoundary - checks if given channel exists in given boundary
func IsInChannelBoundary(channel, boundary string) bool {
	channelsList := GetChannelBoundaryList(boundary)
	return utils.Includes(channelsList, channel)
}

// GetChannelBoundaryList returns a string slice with the channels in given boundary
func GetChannelBoundaryList(boundary string) utils.StringArray {
	if boundary == "" {
		return utils.StringArray{}
	}
	return strings.Split(boundary, ";")
}

// GetResolvedBoundaryChannelList gets the list of resolved channels from given boundary
func GetResolvedBoundaryChannelList(boundary string) utils.StringArray {
	channels := utils.StringArray(GetChannelBoundaryList(boundary))
	return channels.Map(func(s string) string {
		resolved, _ := GetResolvedChannel(s, boundary, "")
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

// OutputChannnelList returns a list of input channels
func OutputChannnelList() utils.StringArray {
	return GetChannelBoundaryList(GetInputChannels())
}

// InputChannelList returns a list of input channels
func InputChannelList() utils.StringArray {
	return GetChannelBoundaryList(GetInputChannels())
}

// GetResolvedChannel gets a resolved channel from a channel name
func GetResolvedChannel(channel, inputChan, outputChan string) (string, error) {
	if IsInChannelBoundary(channel, inputChan) || IsInChannelBoundary(channel, outputChan) {
		return os.Getenv(channel + "_RESOLVED"), nil
	}
	return "", ierrors.NewError().
		InvalidChannel().
		Message("channel " + channel + " not listed as an input or output").
		Build()
}
