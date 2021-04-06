package kafkasc

import (
	"fmt"
	"strings"

	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"go.uber.org/zap"
)

type messageChannel struct {
	channel string
	appCtx  string
	prefix  string
}

func fromResolvedChannel(channel string) (messageChannel, error) {
	logger.Debug("getting Channel as a messageChannel structure")

	ctx, name, err := utils.RemoveLastPartInScope(channel)
	if err != nil {
		logger.Error("unable to get Channel name an context",
			zap.Any("error", err))
		return messageChannel{}, err
	}
	return messageChannel{
		appCtx:  ctx,
		channel: name,
	}, nil
}

// returns specified topic's channel
func fromTopic(topic string) messageChannel {
	logger.Debug("getting Channel given a Kafka Topic")

	msgChan := messageChannel{
		prefix: environment.GetInsprEnvironment(),
		appCtx: environment.GetInsprAppContext(),
	}
	splitTopic := strings.Split(topic, "-")
	msgChan.channel = splitTopic[len(splitTopic)-1]
	msgChan.appCtx = splitTopic[len(splitTopic)-2]
	return msgChan
}

// returns a topic name based on a message channel
func (ch messageChannel) toTopic() string {
	logger.Debug("getting Kafka Topic name given a Channel name")

	var topic string
	ctx, name := ch.appCtx, ch.channel

	if environment.GetInsprEnvironment() == "" {
		topic = fmt.Sprintf("inspr-%s-%s", ctx, name)
	} else {
		topic = fmt.Sprintf(
			"inspr-%s-%s-%s",
			environment.GetInsprEnvironment(),
			ctx,
			name,
		)
	}

	return topic
}
