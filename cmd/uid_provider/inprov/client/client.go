package client

import (
	"context"

	"github.com/inspr/inspr/cmd/uid_provider/api/models"
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

	var resp string
	err := c.rc.Send(ctx, "/login", "POST", models.ReceivedDataLogin{UID: uid, Password: pwd}, &resp)
	if err != nil {
		return "", err
	}
	return resp, nil
}

// CreateUser creates a user in inspr's UID provider.
func (c *Client) CreateUser(ctx context.Context, uid, pwd string, newUser client.User) error {

	var resp interface{}
	err := c.rc.Send(ctx, "/newusr", "POST", models.ReceivedDataCreate{UID: uid, Password: pwd, User: newUser}, resp)
	return err
}

// DeleteUser deletes a user in inspr's UID provider
func (c *Client) DeleteUser(ctx context.Context, uid, pwd, usrToBeDeleted string) error {

	var resp interface{}
	err := c.rc.Send(ctx, "/deleteuser", "POST", models.ReceivedDataDelete{UID: uid, Password: pwd, UserToBeDeleted: usrToBeDeleted}, resp)
	return err
}

// UpdatePassword updates a user's password on inspr's uid provider.
func (c *Client) UpdatePassword(ctx context.Context, uid, pwd, usrToBeUpdated, newPwd string) error {

	var resp interface{}
	err := c.rc.Send(ctx, "/updatepwd", "POST", models.ReceivedDataUpdate{UID: uid, Password: pwd, UserToBeUpdated: usrToBeUpdated, NewPassword: newPwd}, resp)
	return err
}
