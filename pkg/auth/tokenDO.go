package auth

// TokenDO is a data output type for master's token handler
type TokenDO struct {
	Token []byte `json:"token"`
}
