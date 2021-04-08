package client

import (
	"context"

	"gitlab.inspr.dev/inspr/core/cmd/insprd/auth"
)

type User struct {
	UID      string
	Role     int
	Scope    []string
	Password string
}

type RedisManager interface {
	Login(ctx context.Context, uid, pwd string) (string, error)
	CreateUser(ctx context.Context, uid string, newUser User) error
	DeleteUser(ctx context.Context, uid, usrToBeDeleted string) error
	UpdatePassword(ctx context.Context, uid, usrToBeUpdated, newPwd string) error
	RefreshToken(ctx context.Context, refreshToken string) (auth.Payload, error)
}
