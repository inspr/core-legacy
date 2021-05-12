package kafkasc

import (
	"errors"

	"github.com/inspr/inspr/pkg/environment"
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/sidecar/models"
	"github.com/linkedin/goavro"
	"go.uber.org/zap"
)

type kafkaTopic string

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

// returns the channel's Type schema
func (ch kafkaTopic) getSchema() (string, error) {
	logger.Debug("getting Channel schema")
	schema, err := environment.GetSchema(string(ch))
	if err != nil {
		logger.Error("unable to get Channel schema", zap.Any("error", err))

		return "", ierrors.NewError().InnerError(err).Message(err.Error()).Build()
	}

	return schema, nil
}

func (ch kafkaTopic) decode(messageEncoded []byte) (interface{}, error) {
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

func (ch kafkaTopic) encode(message interface{}) ([]byte, error) {
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

func (ch kafkaTopic) readMessage(value []byte) (models.BrokerData, error) {

	// Decoding Message
	message, errDecode := ch.decode(value)
	if errDecode != nil {
		return models.BrokerData{}, errDecode
	}

	channelName := ch

	return models.BrokerData{Message: message, Channel: string(channelName)}, nil
}
