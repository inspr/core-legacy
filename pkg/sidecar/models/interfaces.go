package models

// BROKER's interfaces

// Reader reads from a message broker
type Reader interface {
	ReadMessage(channel string) BrokerResponse
	CommitMessage(channel string) error
}

// Writer writes messages in a message broker
type Writer interface {
	WriteMessage(channel string, msg interface{}) error
}
