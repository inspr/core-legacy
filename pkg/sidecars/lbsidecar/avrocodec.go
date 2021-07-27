package lbsidecar

import (
	"github.com/linkedin/goavro"
	"go.uber.org/zap"
	"inspr.dev/inspr/pkg/environment"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/sidecars/models"
)

// readMessage receives an Avro-encoded message and returns a models.BrokerMessage
// structure that contains the decoded Avro message
func readMessage(ch string, value []byte) (models.BrokerMessage, error) {
	logger.Info("reading message from channel",
		zap.String("channel", ch))

	// Decoding Message
	data, err := decode(ch, value)
	if err != nil {
		return models.BrokerMessage{}, err
	}

	return models.BrokerMessage{Data: data}, nil
}

func encode(ch string, message interface{}) ([]byte, error) {
	logger.Info("encoding message to Avro")

	schema, err := getSchema(ch)
	if err != nil {
		return nil, err
	}

	codec, err := getCodec(schema)
	if err != nil {
		return nil, err
	}

	messageEncoded, err := codec.BinaryFromNative(nil, message)
	if err != nil {
		logger.Error("unable to encode message", zap.Any("error", err))

		return nil, ierrors.New("[ENCODE] %v", err.Error())
	}

	return messageEncoded, nil
}

func decode(ch string, messageEncoded []byte) (interface{}, error) {
	logger.Debug("decoding received message")

	schema, err := getSchema(ch)
	if err != nil {
		return nil, err
	}

	codec, err := getCodec(schema)
	if err != nil {
		return nil, err
	}

	message, _, err := codec.NativeFromBinary(messageEncoded)
	if err != nil {
		logger.Error("unable to decode message", zap.Any("error", err))

		return nil, ierrors.New("[DECODE] %v", err.Error())
	}

	return message, nil
}

// returns the channel type's schema
func getSchema(ch string) (string, error) {
	logger.Debug("getting channel schema",
		zap.String("channel", ch))

	schema, err := environment.GetSchema(ch)
	if err != nil {
		logger.Error("unable to get channel schema", zap.Any("error", err))

		return "", ierrors.From(err)
	}

	return schema, nil
}

// creates Avro codec based on given schema
func getCodec(schema string) (*goavro.Codec, error) {
	logger.Debug("getting Avro codec given a schema",
		zap.String("schema", schema))

	codec, err := goavro.NewCodec(schema)
	if err != nil {
		logger.Error("unable to get Avro codec", zap.Any("error", err))

		return nil, ierrors.From(err)
	}

	return codec, nil
}
