package models

//Payload is information caried by a Inspr acceess token
type Payload struct {
	UID        string   `json:"uid"`
	Role       int      `json:"role"`
	Scope      []string `json:"scope"`
	Refresh    []byte   `json:"refresh"`
	RefreshURL string   `json:"refreshurl"`
}
