package kafkasc

import (
	"errors"

	"github.com/linkedin/goavro"
	"gitlab.inspr.dev/inspr/core/pkg/environment"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
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
func (ch messageChannel) getSchema() (string, error) {
	name, _ := utils.JoinScopes(ch.appCtx, ch.channel)
	schema, err := environment.GetSchema(name)
	if err != nil {
		return "", ierrors.NewError().InnerError(err).Message(err.Error()).Build()
	}

	return schema, nil
}

func (ch messageChannel) decode(messageEncoded []byte) (interface{}, error) {

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
		return nil, errors.New("[DECODE] " + errDecoding.Error())
	}

	return message, nil
}

func (ch messageChannel) encode(message interface{}) ([]byte, error) {
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
		return nil, errors.New("[ENCODE] " + errParseAvro.Error())
	}

	return messageEncoded, nil
}
