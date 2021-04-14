package client

import (
	"context"

	"github.com/inspr/inspr/cmd/uid_provider/client"
	"github.com/inspr/inspr/pkg/rest/request"
	"github.com/spf13/viper"
)

type Client struct {
	rc *request.Client
}

type UIDClient interface {
	// creates payload and sends it to insprd
	// when creating the payload, generetes the Refresh Token (cryptografado)
	Login(ctx context.Context, uid, pwd string) (string, error) // asks Insprd to generate token and saves it into file

	CreateUser(ctx context.Context, uid string, newUser client.User) error
	DeleteUser(ctx context.Context, uid, usrToBeDeleted string) error
	UpdatePassword(ctx context.Context, uid, usrToBeUpdated, newPwd string) error
}

func NewClient() *Client {
	return &Client{
		rc: request.NewJSONClient(viper.GetString("url")),
	}
}

func (c *Client) Login(ctx context.Context, uid, pwd string) (string, error) {
	type ReceivedDataLogin struct {
		UID string
		Pwd string
	}
	var resp string
	err := c.rc.Send(ctx, "/login", "POST", ReceivedDataLogin{uid, pwd}, &resp)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (c *Client) CreateUser(ctx context.Context, uid string, newUser client.User) error {
	type ReceivedDataCreate struct {
		UID string
		Usr client.User
	}

	var resp interface{}
	err := c.rc.Send(ctx, "/newusr", "POST", ReceivedDataCreate{uid, newUser}, resp)
	return err
}

func (c *Client) DeleteUser(ctx context.Context, uid, usrToBeDeleted string) error {
	type ReceivedDataDelete struct {
		UID            string
		UsrToBeDeleted string
	}

	var resp interface{}
	err := c.rc.Send(ctx, "/deleteuser", "POST", ReceivedDataDelete{uid, usrToBeDeleted}, resp)
	return err
}

func (c *Client) UpdatePassword(ctx context.Context, uid, usrToBeUpdated, newPwd string) error {
	type ReceivedDataUpdate struct {
		UID            string
		UsrToBeUpdated string
		NewPwd         string
	}

	var resp interface{}
	err := c.rc.Send(ctx, "/updatepwd", "POST", ReceivedDataUpdate{uid, usrToBeUpdated, newPwd}, resp)
	return err
}
