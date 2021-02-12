package models

// RECEIVED BY BROKER's INTERFACE METHODS

// BrokerResponse - is the struct that represents
// the return of the interface methods
type BrokerResponse struct {
	Data interface{} `json:"data,omitempty"`
}

// RECEIVED BY DAPP CLIENT

// Message - data to be put in the message
type Message struct {
	Data interface{} `json:"data,omitempty"`
}

// BrokerData represents a an http request structure
type BrokerData struct {
	Message Message `json:"message"`
	Channel string  `json:"channel"`
}
