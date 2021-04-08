package client

import (
	"context"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/auth"
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
	Login(ctx context.Context, uid, pwd string) (string, error)
	CreateUser(ctx context.Context, uid string, newUser User) error
	DeleteUser(ctx context.Context, uid, usrToBeDeleted string) error
	UpdatePassword(ctx context.Context, uid, usrToBeUpdated, newPwd string) error
	RefreshToken(ctx context.Context, refreshToken string) (auth.Payload, error)
}
