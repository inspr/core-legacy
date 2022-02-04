package models

// BrokerMessage is the struct that represents the client's request format
type BrokerMessage struct {
	Data interface{} `json:"data"`
}

// ConnectionVariables is the structure resposible for storing
// enviroment variable names regarding connection ports for sidecars
type ConnectionVariables struct {
	ReadEnvVar  string
	WriteEnvVar string
}

// BrokerHandler is a structure dedicated to agregating the interfaces of a broker.
// This object implements the BrokerInterface.
type BrokerHandler struct {
	Broker string
	reader Reader
	writer Writer
}

func NewBrokerHandler(broker string, reader Reader, writer Writer) *BrokerHandler {
	return &BrokerHandler{
		Broker: broker,
		reader: reader,
		writer: writer,
	}
}

func (bh *BrokerHandler) Writer() Writer {
	return bh.writer
}

func (bh *BrokerHandler) Reader() Reader {
	return bh.reader
}
