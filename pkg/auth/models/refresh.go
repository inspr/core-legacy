package models

// ResfreshDO is the body type expected by UID provider to refresh a payload
type ResfreshDO struct {
	RefreshToken []byte `json:"refreshtoken"`
}

// ResfreshDI is a data input type, expected format for Refresh's enpoint body
type ResfreshDI struct {
	RefreshToken []byte `json:"refreshtoken"`
	RefreshURL   string `json:"refreshurl"`
}
