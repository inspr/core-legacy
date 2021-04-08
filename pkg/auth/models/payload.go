package models

//Payload is information caried by a Inspr acceess token
type Payload struct {
	UID        string   `json:"uid"`
	Role       int      `json:"role"`
	Scope      []string `json:"scope"`
	Refresh    string   `json:"refresh"`
	RefreshURL string   `json:"refreshurl"`
}
