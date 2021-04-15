package models

import "github.com/inspr/inspr/cmd/uid_provider/client"

type ReceivedDataCreate struct {
	UID      string      `json:"uid"`
	Password string      `json:"password"`
	User     client.User `json:"user"`
}

type ReceivedDataDelete struct {
	UID             string `json:"uid"`
	Password        string `json:"password"`
	UserToBeDeleted string `json:"username"`
}

type ReceivedDataUpdate struct {
	UID             string `json:"uid"`
	Password        string `json:"password"`
	UserToBeUpdated string `json:"username"`
	NewPassword     string `json:"userpassword"`
}

type ReceivedDataLogin struct {
	UID      string `json:"uid"`
	Password string `json:"password"`
}

type ReceivedDataRefresh struct {
	RefreshToken []byte `json:"refreshtoken"`
}
