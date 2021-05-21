package models

// BrokerMessage is the struct that represents the client's request format
type BrokerMessage struct {
	Message interface{} `json:"message"`
}

// ConnectionVariables
type ConnectionVariables struct {
	ReadVar  string
	WriteVar string
}
