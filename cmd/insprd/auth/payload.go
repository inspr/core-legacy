package auth

type payload struct {
	UID string
	Role int
	Scope string[]
	Refresh string
}