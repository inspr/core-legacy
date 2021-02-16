package kafkasc

import (
	"fmt"

	"gitlab.inspr.dev/inspr/core/pkg/environment"
)

type messageChannel struct {
	channel string
	appCtx  string
	prefix  string
}

// returns specified topic's channel
func fromTopic(topic string) messageChannel {
	msgChan := messageChannel{
		prefix: environment.GetEnvironment().InsprEnvironment,
		appCtx: environment.GetEnvironment().InsprAppContext,
	}

	if msgChan.prefix == "" {
		msgChan.channel = topic[len("inspr-"+msgChan.appCtx+"-"):]
	} else {
		msgChan.channel = topic[len("inspr-"+msgChan.prefix+"-"+msgChan.appCtx+"-"):]
	}
	return msgChan
}

// returns a topic name based on given channel
func toTopic(channel string) string {
	var topic string

	if environment.GetEnvironment().InsprEnvironment == "" {
		topic = fmt.Sprintf("inspr-%s-%s", environment.GetEnvironment().InsprAppContext, channel)
	} else {
		topic = fmt.Sprintf(
			"inspr-%s-%s-%s",
			environment.GetEnvironment().InsprEnvironment,
			environment.GetEnvironment().InsprAppContext,
			channel,
		)
	}

	return topic
}
