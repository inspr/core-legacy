package client

import (
	"context"

	"github.com/inspr/inspr/cmd/insprd/auth"
)

// User defines the information a user contains
type User struct {
	UID      string
	Role     int
	Scope    []string
	Password string
}

// RedisManager defines methods to manage Redis in the cluster
type RedisManager interface {
	UIDClient
	RefreshToken(ctx context.Context, refreshToken string) (auth.Payload, error)
}

// UIDClient is a client for communicating with inspr's uid provider
type UIDClient interface {
	// creates payload and sends it to insprd
	// when creating the payload, generetes the Refresh Token (cryptografado)
	Login(ctx context.Context, uid, pwd string) (string, error) // asks Insprd to generate token and saves it into file

	CreateUser(ctx context.Context, uid string, newUser User) error
	DeleteUser(ctx context.Context, uid, usrToBeDeleted string) error
	UpdatePassword(ctx context.Context, uid, usrToBeUpdated, newPwd string) error
}
