package sidecar

// Reader reads from a message broker
type Reader interface {
	ReadMessage() (*string, interface{}, error)
	Commit() error
}

// Writer writes messages in a message broker
type Writer interface {
	WriteMessage(message interface{}, channel string) error
}
