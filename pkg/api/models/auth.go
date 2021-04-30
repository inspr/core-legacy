package models

// AuthDI - Data Input format for authorization requests
type AuthDI struct {
	Token []byte `json:"token"`
}
