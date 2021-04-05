package auth

//Auth is the inteface for interacting with the Authentication service
type Auth interface {
	Validade(token string) (bool, error)
	Tokenize(load payload) (string, error)
}
