package models

// BrokerMessage is the struct that represents the client's request format
type BrokerMessage struct {
	Data interface{} `json:"data"`
}
