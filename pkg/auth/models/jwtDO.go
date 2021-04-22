package models

// JwtDO is a data output type for master's token handler
type JwtDO struct {
	Token []byte `json:"token"`
}
