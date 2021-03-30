package auth

type Auth interface {
	Validade(token string) (bool, error)
	Login(usr, pwd string) (User, error)
	Register(usr User)
}
