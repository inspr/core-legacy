package models

// Message - data to be put in the message
type Message struct {
	Data interface{} `json:"data,omitempty"`
}

// BrokerData represents a an http request structure
type BrokerData struct {
	Message Message `json:"message"`
	Channel string  `json:"channel"`
}
