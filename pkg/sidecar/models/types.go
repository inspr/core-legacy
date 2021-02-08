package models

// RECEIVED BY BROKER's INTERFACE METHODS

// BrokerResponse - is the struct that represents
// the return of the interface methods
type BrokerResponse struct {
	Data interface{} `json:"data,omitempty"`
}

// RECEIVED BY DAPP CLIENT

// BodyMessage - data to be put in the message
type BodyMessage struct {
	Data interface{} `json:"data,omitempty"`
}

// RequestBody represents a an http request structure
type RequestBody struct {
	Message BodyMessage `json:"message"`
	Channel string      `json:"channel"`
}
