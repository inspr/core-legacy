package kafkasc

import (
	"errors"

	"github.com/linkedin/goavro"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
	"go.uber.org/zap"
)

// creates Avro codec based on given schema
func getCodec(schema string) (*goavro.Codec, error) {
	logger.Debug("getting Avro codec given a schema",
		zap.String("schema", schema))

	codec, errCreateCodec := goavro.NewCodec(schema)
	if errCreateCodec != nil {
		logger.Error("unable to get Avro codec", zap.Any("error", errCreateCodec))
		return nil, errors.New("[KAFKA_PRODUCE_CODEC] " + errCreateCodec.Error())
	}

	return codec, nil
}

// returns the channel's channel type schema
func (ch messageChannel) getSchema() (string, error) {
	logger.Debug("getting Channel schema")

	name, _ := utils.JoinScopes(ch.appCtx, ch.channel)
	schema, err := environment.GetSchema(name)
	if err != nil {
		logger.Error("unable to get Channel schema", zap.Any("error", err))
		return "", ierrors.NewError().InnerError(err).Message(err.Error()).Build()
	}

	return schema, nil
}

func (ch messageChannel) decode(messageEncoded []byte) (interface{}, error) {
	logger.Debug("decoding received Kafka message")

	schema, errGetSchema := ch.getSchema()
	if errGetSchema != nil {
		return nil, errGetSchema
	}

	codec, errCreateCodec := getCodec(schema)
	if errCreateCodec != nil {
		return nil, errCreateCodec
	}

	message, _, errDecoding := codec.NativeFromBinary(messageEncoded)
	if errDecoding != nil {
		logger.Error("unable to decode Kafka message", zap.Any("error", errDecoding))
		return nil, errors.New("[DECODE] " + errDecoding.Error())
	}

	return message, nil
}

func (ch messageChannel) encode(message interface{}) ([]byte, error) {
	logger.Debug("encoding Kafka message")

	schema, errGetSchema := ch.getSchema()
	if errGetSchema != nil {
		return nil, errGetSchema
	}

	codec, errCreateCodec := getCodec(schema)
	if errCreateCodec != nil {
		return nil, errCreateCodec
	}

	messageEncoded, errParseAvro := codec.BinaryFromNative(nil, message)
	if errParseAvro != nil {
		logger.Error("unable to encode Kafka message", zap.Any("error", errParseAvro))
		return nil, errors.New("[ENCODE] " + errParseAvro.Error())
	}

	return messageEncoded, nil
}