package models

//Payload is information caried by a Inspr acceess token
type Payload struct {
	UID        string   `json:"uid"`
	Permission []string `json:"permission"`
	Scope      []string `json:"scope"`
	Refresh    []byte   `json:"refresh"`
	RefreshURL string   `json:"refreshurl"`
}

// All Permission possible values
const (
	CreateDapp        string = "create:dapp"
	CreateChannel     string = "create:channel"
	CreateChannelType string = "create:ctype"
	CreateAlias       string = "create:alias"

	GetDapp        string = "get:dapp"
	GetChannel     string = "get:channel"
	GetChannelType string = "get:ctype"
	GetAlias       string = "get:alias"

	UpdateDapp        string = "update:dapp"
	UpdateChannel     string = "update:channel"
	UpdateChannelType string = "update:ctype"
	UpdateAlias       string = "update:alias"

	DeleteDapp        string = "delete:dapp"
	DeleteChannel     string = "delete:channel"
	DeleteChannelType string = "delete:ctype"
	DeleteAlias       string = "delete:alias"

	CreateToken string = "create:token"
)
