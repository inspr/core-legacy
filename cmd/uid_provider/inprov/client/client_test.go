package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/inspr/inspr/cmd/uid_provider/api/models"
	"github.com/inspr/inspr/cmd/uid_provider/client"
	"github.com/inspr/inspr/pkg/rest"
	"github.com/inspr/inspr/pkg/rest/request"
)

// Response Codes
// correct login - 200
// not found - 404
// incorrect password - 401
func TestClient_Login(t *testing.T) {

	type args struct {
		ctx context.Context
		uid string
		pwd string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		err     error
		wantErr bool
	}{
		{
			name: "correct request",
			args: args{
				ctx: context.Background(),
				uid: "this is a uid",
				pwd: "this is a password",
			},
			want: "this is a token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(rest.Handler(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("method is not POST")
				}

				if r.URL.Path != "/login" {
					t.Errorf("path is not login, is %v", r.URL.Path)
				}

				type ReceivedDataLogin struct {
					UID string
					Pwd string
				}
				receiver := ReceivedDataLogin{}
				decoder := json.NewDecoder(r.Body)
				err := decoder.Decode(&receiver)
				if err != nil {
					t.Errorf("error in decoding body: %v", err)
				}

				if receiver.Pwd != tt.args.pwd {
					t.Errorf("password does not match parameter: %v != %v", receiver.Pwd, tt.args.pwd)
				}
				if receiver.UID != tt.args.uid {
					t.Errorf("UID does not match parameter: %v != %v", receiver.UID, tt.args.uid)
				}
				encoder := json.NewEncoder(w)
				encoder.Encode(tt.want)
			}).JSON().Post())

			c := &Client{
				rc: request.NewJSONClient(server.URL),
			}
			got, err := c.Login(tt.args.ctx, tt.args.uid, tt.args.pwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestClient_CreateUser(t *testing.T) {

	type args struct {
		ctx     context.Context
		uid     string
		pwd     string
		newUser client.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "correct functionality",
			args: args{
				ctx: context.Background(),
				uid: "this is an uid",
				pwd: "this is a password",
				newUser: client.User{
					UID:  "this is a new UID",
					Role: 2,
					Scope: []string{
						"scope.1",
						"scope.2",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(rest.Handler(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("method is not POST")
				}

				if r.URL.Path != "/newusr" {
					t.Errorf("path is not newusr, is %v", r.URL.Path)
				}

				receiver := models.ReceivedDataCreate{}
				decoder := json.NewDecoder(r.Body)
				err := decoder.Decode(&receiver)
				if err != nil {
					t.Errorf("error in decoding body: %v", err)
				}

				if !reflect.DeepEqual(receiver.User, tt.args.newUser) {
					t.Errorf("body does not match request \n%v\n!=\n%v", receiver.User, tt.args.newUser)
				}

				if receiver.UID != tt.args.uid {
					t.Errorf("uid is not the same %v != %v", receiver.UID, tt.args.uid)
				}
				if receiver.Password != tt.args.pwd {
					t.Errorf("pwd is not the same %v != %v", receiver.Password, tt.args.pwd)
				}

			}).JSON().Post())

			c := &Client{
				rc: request.NewJSONClient(server.URL),
			}
			if err := c.CreateUser(tt.args.ctx, tt.args.uid, tt.args.pwd, tt.args.newUser); (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_DeleteUser(t *testing.T) {
	type args struct {
		ctx            context.Context
		uid            string
		pwd            string
		usrToBeDeleted string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "correct functionality",
			args: args{
				ctx:            context.Background(),
				uid:            "this is a uid",
				pwd:            "this is a password",
				usrToBeDeleted: "this is the user to be deleted",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(rest.Handler(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("method is not POST")
				}

				if r.URL.Path != "/deleteuser" {
					t.Errorf("path is not deleteusr, is %v", r.URL.Path)
				}

				receiver := models.ReceivedDataDelete{}
				decoder := json.NewDecoder(r.Body)
				err := decoder.Decode(&receiver)
				if err != nil {
					t.Errorf("error in decoding body: %v", err)
				}

				if receiver.UserToBeDeleted != tt.args.usrToBeDeleted {
					t.Errorf("user to be deleted does not match \n%v\n!=\n%v", receiver.UserToBeDeleted, tt.args.usrToBeDeleted)
				}

				if receiver.UID != tt.args.uid {
					t.Errorf("uid is not the same %v != %v", receiver.UID, tt.args.uid)
				}

				if receiver.Password != tt.args.pwd {
					t.Errorf("password is not the same %v != %v", receiver.Password, tt.args.pwd)
				}

			}).JSON().Post())

			c := &Client{
				rc: request.NewJSONClient(server.URL),
			}

			if err := c.DeleteUser(tt.args.ctx, tt.args.uid, tt.args.pwd, tt.args.usrToBeDeleted); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_UpdatePassword(t *testing.T) {
	type args struct {
		ctx            context.Context
		uid            string
		pwd            string
		usrToBeUpdated string
		newPwd         string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "correct functionality",
			args: args{
				ctx:            context.Background(),
				uid:            "this is a uid",
				pwd:            "this is a password",
				usrToBeUpdated: "this is the user to be updated",
				newPwd:         "this is the new password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(rest.Handler(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("method is not POST")
				}

				if r.URL.Path != "/updatepwd" {
					t.Errorf("path is not updatepwd, is %v", r.URL.Path)
				}

				receiver := models.ReceivedDataUpdate{}
				decoder := json.NewDecoder(r.Body)
				err := decoder.Decode(&receiver)
				if err != nil {
					t.Errorf("error in decoding body: %v", err)
				}

				if receiver.NewPassword != tt.args.newPwd {
					t.Errorf("new password does not match \n%v\n!=\n%v", receiver.NewPassword, tt.args.newPwd)
				}

				if receiver.UID != tt.args.uid {
					t.Errorf("uid is not the same %v != %v", receiver.UID, tt.args.uid)
				}

				if receiver.UserToBeUpdated != tt.args.usrToBeUpdated {
					t.Errorf("user to be updated does not match \n%v\n!=\n%v", receiver.UserToBeUpdated, tt.args.usrToBeUpdated)
				}

				if receiver.Password != tt.args.pwd {
					t.Errorf("password is not the same %v != %v", receiver.Password, tt.args.pwd)
				}

			}).JSON().Post())

			c := &Client{
				rc: request.NewJSONClient(server.URL),
			}
			if err := c.UpdatePassword(tt.args.ctx, tt.args.uid, tt.args.pwd, tt.args.usrToBeUpdated, tt.args.newPwd); (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
