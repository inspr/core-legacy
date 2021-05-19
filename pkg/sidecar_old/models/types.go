package models

// BrokerData represents a an http request structure
type BrokerData struct {
	Message interface{} `json:"message"`
	Channel string      `json:"channel"`
}
