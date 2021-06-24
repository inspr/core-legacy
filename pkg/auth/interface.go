package auth

//Auth is the inteface for interacting with the Authentication service
type Auth interface {
	Validate(token []byte) (*Payload, []byte, error)
	Tokenize(load Payload) ([]byte, error)
	Init(key string, load Payload) ([]byte, error)
	Refresh(token []byte) ([]byte, error)
}
