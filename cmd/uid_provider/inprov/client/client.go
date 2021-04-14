package client

import (
	"context"

	"github.com/inspr/inspr/cmd/uid_provider/client"
	"github.com/inspr/inspr/pkg/rest/request"
	"github.com/spf13/viper"
)

// Client is the client for communicating with the in-cluster uidp
type Client struct {
	rc *request.Client
}

// NewClient creates a new client for communicating with inspr's UID provider.
func NewClient() *Client {
	return &Client{
		rc: request.NewJSONClient(viper.GetString("url")),
	}
}

// Login creates a request to log in to a uid provider. It returns a signed token for
// communicating with the insprd cluster in question.
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

// CreateUser creates a user in inspr's UID provider.
func (c *Client) CreateUser(ctx context.Context, uid string, newUser client.User) error {
	type ReceivedDataCreate struct {
		UID string
		Usr client.User
	}

	var resp interface{}
	err := c.rc.Send(ctx, "/newusr", "POST", ReceivedDataCreate{uid, newUser}, resp)
	return err
}

// DeleteUser deletes a user in inspr's UID provider
func (c *Client) DeleteUser(ctx context.Context, uid, usrToBeDeleted string) error {
	type ReceivedDataDelete struct {
		UID            string
		UsrToBeDeleted string
	}

	var resp interface{}
	err := c.rc.Send(ctx, "/deleteuser", "POST", ReceivedDataDelete{uid, usrToBeDeleted}, resp)
	return err
}

// UpdatePassword updates a user's password on inspr's uid provider.
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
