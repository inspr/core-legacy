package sidecar

// Message represents a Inspr message
type Message struct {
	Commit  bool        `json:"commit,omitempty"`
	Channel string      `json:"channel,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}
