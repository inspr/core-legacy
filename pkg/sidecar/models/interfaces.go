package models

// Reader reads from a message broker
type Reader interface {
	ReadMessage(channel string) (Message, error)
	CommitMessage(channel string) error
}

// Writer writes messages in a message broker
type Writer interface {
	WriteMessage(channel string, msg Message) error
}
