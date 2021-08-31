package models

import "inspr.dev/inspr/cmd/uidp/client"

// ReceivedDataCreate defines the expected data structure to be
// received as a response when calling the method Create(...)
// of RedisManager interface
type ReceivedDataCreate struct {
	UID      string      `json:"uid"`
	Password string      `json:"password"`
	User     client.User `json:"user"`
}

// ReceivedDataDelete defines the expected data structure to be
// received as a response when calling the method Delete(...)
// of RedisManager interface
type ReceivedDataDelete struct {
	UID             string `json:"uid"`
	Password        string `json:"password"`
	UserToBeDeleted string `json:"username"`
}

// ReceivedDataUpdate defines the expected data structure to be
// received as a response when calling the method Update(...)
// of RedisManager interface
type ReceivedDataUpdate struct {
	UID             string `json:"uid"`
	Password        string `json:"password"`
	UserToBeUpdated string `json:"username"`
	NewPassword     string `json:"userpassword"`
}

// ReceivedDataLogin defines the expected data structure to be
// received as a response when calling the method Login(...)
// of RedisManager interface
type ReceivedDataLogin struct {
	UID      string `json:"uid"`
	Password string `json:"password"`
}

// ReceivedDataRefresh defines the expected data structure to be
// received as a response when calling the method Refresh(...)
// of RedisManager interface
type ReceivedDataRefresh struct {
	RefreshToken []byte `json:"refreshtoken"`
}
