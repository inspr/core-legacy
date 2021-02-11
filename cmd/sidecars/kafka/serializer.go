package kafka

import (
	"errors"

	"github.com/linkedin/goavro"
	"gitlab.inspr.dev/inspr/core/cmd/insprd/memory/tree"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
)

// creates Avro codec based on given schema
func getCodec(schema string) (*goavro.Codec, error) {
	codec, errCreateCodec := goavro.NewCodec(schema)
	if errCreateCodec != nil {
		return nil, errors.New("[KAFKA_PRODUCE_CODEC] " + errCreateCodec.Error())
	}

	return codec, nil
}

// returns the channel's channel type schema
func getSchema(channel, context string) (string, error) {
	memChan, err := tree.GetTreeMemory().Channels().GetChannel(context, channel)
	if err != nil {
		return "", ierrors.NewError().InvalidChannel().Message(err.Error()).Build()
	}

	memChanType, err := tree.GetTreeMemory().ChannelTypes().GetChannelType(context, memChan.Spec.Type)
	// this error shouldn't occur given that a channel cannot be created without a valid channel type
	if err != nil {
		return "", ierrors.NewError().InvalidChannelType().Message(err.Error()).Build()
	}

	return string(memChanType.Schema), nil
}

func decode(messageEncoded []byte, channel string) (interface{}, error) {
	schema, errGetSchema := getSchema(channel, environment.GetEnvironment().InsprAppContext)
	if errGetSchema != nil {
		return nil, errGetSchema
	}

	codec, errCreateCodec := getCodec(schema)
	if errCreateCodec != nil {
		return nil, errCreateCodec
	}

	message, _, errDecoding := codec.NativeFromBinary(messageEncoded)
	if errDecoding != nil {
		return nil, errors.New("[DECODE] " + errDecoding.Error())
	}

	return message, nil
}

func encode(message interface{}, channel string) ([]byte, error) {
	schema, errGetSchema := getSchema(channel, environment.GetEnvironment().InsprAppContext)
	if errGetSchema != nil {
		return nil, errGetSchema
	}

	codec, errCreateCodec := getCodec(schema)
	if errCreateCodec != nil {
		return nil, errCreateCodec
	}

	messageEncoded, errParseAvro := codec.BinaryFromNative(nil, message)
	if errParseAvro != nil {
		return nil, errors.New("[ENCODE] " + errParseAvro.Error())
	}

	return messageEncoded, nil
}
