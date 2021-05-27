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
