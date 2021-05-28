package environment

import (
	"os"
	"strings"

	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta/brokers"
	"github.com/inspr/inspr/pkg/utils"
)

// IsInChannelBoundary - checks if given channel exists in given boundary
func IsInChannelBoundary(channel string, boundary []brokers.ChannelBroker) bool {
	channelsList := GetChannelBoundaryList(boundary)
	return utils.Includes(channelsList, channel)
}

// GetChannelBoundaryList returns a string slice with the channels in given boundary
func GetChannelBoundaryList(channels []brokers.ChannelBroker) utils.StringArray {
	if channels == nil {
		return utils.StringArray{}
	}
	boundary := utils.StringArray{}
	for _, channel := range channels {
		boundary = append(boundary, channel.ChName)
	}
	return boundary
}

// GetResolvedBoundaryChannelList gets the list of resolved channels from given boundary
func GetResolvedBoundaryChannelList(channels []brokers.ChannelBroker) utils.StringArray {
	boundary := GetChannelBoundaryList(channels)
	return boundary.Map(func(s string) string {
		resolved, _ := GetResolvedChannel(s, channels, nil)
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
	return GetChannelBoundaryList(GetOutputChannels())
}

// InputChannelList returns a list of input channels
func InputChannelList() utils.StringArray {
	return GetChannelBoundaryList(GetInputChannels())
}

// OutputBrokerChannnels returns a list of input channels
func OutputBrokerChannnels(broker string) utils.StringArray {
	return GetChannelBoundaryList(GetOutputChannels()) // WRONG IMPLEMENTATION
}

// InputBrokerChannels returns a list of input channels
func InputBrokerChannels(broker string) utils.StringArray {
	return GetChannelBoundaryList(GetInputChannels()) // WRONG IMPLEMENTATION
}

// GetResolvedChannel gets a resolved channel from a channel name
func GetResolvedChannel(channel string, inputChan, outputChan []brokers.ChannelBroker) (string, error) {
	if IsInChannelBoundary(channel, inputChan) || IsInChannelBoundary(channel, outputChan) {
		return os.Getenv(channel + "_RESOLVED"), nil
	}
	return "", ierrors.NewError().
		InvalidChannel().
		Message("channel " + channel + " not listed as an input or output").
		Build()
}

func getChannelList(channelList string) []brokers.ChannelBroker {
	if channelList == "" {
		return nil
	}
	channelBrokers := utils.StringArray(strings.Split(channelList, ";"))
	channels := []brokers.ChannelBroker{}
	channelBrokers.Map(func(channel string) string {
		data := strings.Split(channel, "_")
		bla := brokers.ChannelBroker{
			ChName: data[0],
			Broker: data[1],
		}
		channels = append(channels, bla)
		return ""
	})
	return channels
}

func filterChannelsByBroker(broker string, channels []brokers.ChannelBroker) utils.StringArray {
	brokerChannels := utils.StringArray{}
	for _, channel := range channels {
		if channel.Broker == broker {
			brokerChannels = append(brokerChannels, channel.ChName)
		}
	}
	return brokerChannels
}
