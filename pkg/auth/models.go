package auth

// ResfreshDO is the body type expected by UID provider to refresh a payload
type ResfreshDO struct {
	RefreshToken []byte `json:"refreshtoken"`
}

// JwtDO is a data output type for master's token handler
type JwtDO struct {
	Token []byte `json:"token"`
}

//Payload is information caried by a Inspr acceess token
type Payload struct {
	UID string `json:"uid"`
	// Permissions is a map where key is the Scope and values are permissions
	Permissions map[string][]string `json:"permissions"`
	Refresh     []byte              `json:"refresh"`
	RefreshURL  string              `json:"refreshurl"`
}

// All Permissions possible values
const (
	CreateDapp    string = "create:dapp"
	CreateChannel string = "create:channel"
	CreateType    string = "create:ctype"
	CreateAlias   string = "create:alias"

	GetDapp    string = "get:dapp"
	GetChannel string = "get:channel"
	GetType    string = "get:ctype"
	GetAlias   string = "get:alias"

	UpdateDapp    string = "update:dapp"
	UpdateChannel string = "update:channel"
	UpdateType    string = "update:ctype"
	UpdateAlias   string = "update:alias"

	DeleteDapp    string = "delete:dapp"
	DeleteChannel string = "delete:channel"
	DeleteType    string = "delete:ctype"
	DeleteAlias   string = "delete:alias"

	CreateToken string = "create:token"
)
