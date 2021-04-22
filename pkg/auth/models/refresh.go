package models

// ResfreshDO is the body type expected by UID provider to refresh a payload
type ResfreshDO struct {
	RefreshToken []byte `json:"refreshtoken"`
}
