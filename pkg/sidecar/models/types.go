package models

// Message represents a Inspr message
type Message struct {
	Commit  bool        `json:"commit,omitempty"`
	Channel string      `json:"channel,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

// RequestBody represents a an http request structure
type RequestBody struct {
	Message Message `json:"message"`
	Channel string  `json:"channel"`
}
