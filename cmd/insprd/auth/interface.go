package auth

//Auth is the inteface for interacting with the Authentication service
type Auth interface {
	Validade(token string) (string, error)
	Tokenize(load payload) (string, error)
}
