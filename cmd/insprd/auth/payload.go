package auth

//Payload is information caried by a Inspr acceess token
type Payload struct {
	UID        string
	Role       int
	Scope      []string
	Refresh    string
	RefreshURL string
}