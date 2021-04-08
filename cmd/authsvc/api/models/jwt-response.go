package models

// JwtDO is the data output type for jwt access token
type JwtDO struct {
	Token []byte `json:"token"`
}
