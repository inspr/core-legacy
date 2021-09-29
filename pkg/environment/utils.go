package environment

import (
	"os"
	"strings"

	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/meta/brokers"
	"inspr.dev/inspr/pkg/utils"
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
		return "", ierrors.New("schema for channel %s not found", channel).NotFound()
	}
	return schema, nil
}

// OutputChannelList returns a list of input channels
func OutputChannelList() utils.StringArray {
	return GetChannelBoundaryList(GetOutputChannelsData())
}

// InputChannelList returns a list of input channels
func InputChannelList() utils.StringArray {
	return GetChannelBoundaryList(GetInputChannelsData())
}

// OutputBrokerChannnels returns a list of input channels
func OutputBrokerChannnels(broker string) utils.StringArray {
	return GetChannelBoundaryList(GetOutputChannelsData()) // WRONG IMPLEMENTATION
}

// InputBrokerChannels returns a list of input channels
func InputBrokerChannels(broker string) utils.StringArray {
	return GetChannelBoundaryList(GetInputChannelsData()) // WRONG IMPLEMENTATION
}

// GetResolvedChannel gets a resolved channel from a channel name
func GetResolvedChannel(channel string, inputChan, outputChan []brokers.ChannelBroker) (string, error) {
	if IsInChannelBoundary(channel, inputChan) || IsInChannelBoundary(channel, outputChan) {
		return os.Getenv(channel + "_RESOLVED"), nil
	}
	return "", ierrors.New(
		"channel " + channel + " not listed as an input or output",
	).InvalidChannel()
}

// GetRouteData returns the data of a resolved route
func GetRouteData(route string) (*meta.RouteConnection, error) {
	value, exists := os.LookupEnv(route + "_ROUTE")
	if !exists {
		return nil, ierrors.New("invalid route: %s", route).BadRequest()
	}
	data := strings.Split(value, ";")
	return &meta.RouteConnection{
		Address:   data[0],
		Endpoints: data[1:],
	}, nil
}

func getChannelData(channelList string) []brokers.ChannelBroker {
	if channelList == "" {
		return nil
	}
	channelBrokers := utils.StringArray(strings.Split(channelList, ";"))
	channels := []brokers.ChannelBroker{}
	channelBrokers.Map(func(channel string) string {
		data := strings.Split(channel, "@")
		channelData := brokers.ChannelBroker{
			ChName: data[0],
			Broker: data[1],
		}
		channels = append(channels, channelData)
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
