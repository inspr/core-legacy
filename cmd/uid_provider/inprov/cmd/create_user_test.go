package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/inspr/inspr/cmd/uid_provider/client"
	"gopkg.in/yaml.v2"
)

type mockCl struct {
	createUser     func(context.Context, string, string, client.User) error
	deleteUser     func(context.Context, string, string, string) error
	login          func(context.Context, string, string) (string, error)
	updatePassword func(context.Context, string, string, string, string) error
}

func (cl mockCl) CreateUser(ctx context.Context, uid, pwd string, newUsr client.User) error {
	return cl.createUser(ctx, uid, pwd, newUsr)
}
func (cl mockCl) DeleteUser(ctx context.Context, uid, pwd string, tbd string) error {
	return cl.deleteUser(ctx, uid, pwd, tbd)
}
func (cl mockCl) Login(ctx context.Context, uid string, password string) (string, error) {
	return cl.login(ctx, uid, password)
}
func (cl mockCl) UpdatePassword(c context.Context, uid, pwd string, tobeUpdated string, newPwd string) error {
	return cl.updatePassword(c, uid, pwd, tobeUpdated, newPwd)
}

func Test_createUser(t *testing.T) {
	type args struct {
		c context.Context
		s []string
	}
	tests := []struct {
		name    string
		args    args
		usr     client.User
		wantErr bool
		before  func()
		after   func()
		options createUserOptionsDT
	}{
		{
			usr: client.User{
				UID:      "this is a uid",
				Role:     12,
				Scope:    []string{"scope1", "scope2"},
				Password: "password",
			},
			name: "yaml user creation",
			args: args{
				c: context.Background(),
				s: []string{"user", "password"},
			},
			options: createUserOptionsDT{yaml: "user.yaml"},
			before: func() {
				usr := client.User{
					UID:      "this is a uid",
					Role:     12,
					Scope:    []string{"scope1", "scope2"},
					Password: "password",
				}
				buf, _ := yaml.Marshal(usr)
				ioutil.WriteFile("user.yaml", buf, os.ModePerm)
			},
			after: func() {
				os.Remove("user.yaml")
			},
			wantErr: false,
		},
		{
			usr: client.User{
				UID:      "this is a uid",
				Role:     12,
				Scope:    []string{"scope1", "scope2"},
				Password: "password",
			},
			name: "yaml user creation -- invalid file name",
			args: args{
				c: context.Background(),
				s: []string{"user", "password"},
			},
			options: createUserOptionsDT{yaml: "inexistant-file.yaml"},
			before: func() {
				usr := client.User{
					UID:      "this is a uid",
					Role:     12,
					Scope:    []string{"scope1", "scope2"},
					Password: "password",
				}
				buf, _ := yaml.Marshal(usr)
				ioutil.WriteFile("user.yaml", buf, os.ModePerm)
			},
			after: func() {
				os.Remove("user.yaml")
			},
			wantErr: true,
		},
		{
			usr: client.User{
				UID:      "this is a uid",
				Role:     12,
				Scope:    []string{"scope1", "scope2"},
				Password: "password",
			},
			name: "yaml user creation -- invalid file format",
			args: args{
				c: context.Background(),
				s: []string{"user", "password"},
			},
			options: createUserOptionsDT{yaml: "user.yaml"},
			before: func() {
				usr := client.User{
					UID:      "this is a uid",
					Role:     12,
					Scope:    []string{"scope1", "scope2"},
					Password: "password",
				}
				buf, _ := json.Marshal(usr)
				ioutil.WriteFile("user.yaml", buf, os.ModePerm)
			},
			after: func() {
				os.Remove("user.yaml")
			},
			wantErr: true,
		},
		{
			usr: client.User{
				UID:      "this is a uid",
				Role:     12,
				Scope:    []string{"scope1", "scope2"},
				Password: "password",
			},
			name: "json user creation",
			args: args{
				c: context.Background(),
				s: []string{"user", "password"},
			},
			options: createUserOptionsDT{json: "user.json"},
			before: func() {
				usr := client.User{
					UID:      "this is a uid",
					Role:     12,
					Scope:    []string{"scope1", "scope2"},
					Password: "password",
				}
				buf, _ := json.Marshal(usr)
				ioutil.WriteFile("user.json", buf, os.ModePerm)
			},
			after: func() {
				os.Remove("user.json")
			},
			wantErr: false,
		},
		{
			usr: client.User{
				UID:      "this is a uid",
				Role:     12,
				Scope:    []string{"scope1", "scope2"},
				Password: "password",
			},
			name: "json user creation -- invalid file format",
			args: args{
				c: context.Background(),
				s: []string{"user", "password"},
			},
			options: createUserOptionsDT{json: "user.json"},
			before: func() {
				usr := client.User{
					UID:      "this is a uid",
					Role:     12,
					Scope:    []string{"scope1", "scope2"},
					Password: "password",
				}
				buf, _ := yaml.Marshal(usr)
				ioutil.WriteFile("user.json", buf, os.ModePerm)
			},
			after: func() {
				os.Remove("user.json")
			},
			wantErr: true,
		},

		{
			usr: client.User{
				UID:      "this is a uid",
				Role:     12,
				Scope:    []string{"scope1", "scope2"},
				Password: "password",
			},
			name: "flag user creation",
			args: args{
				c: context.Background(),
				s: []string{"user", "password"},
			},
			options: createUserOptionsDT{
				username: "this is a uid",
				password: "password",
				scopes:   []string{"scope1", "scope2"},
				role:     12,
			},
			wantErr: false,
		},

		{
			usr: client.User{
				UID:      "this is a uid",
				Role:     12,
				Scope:    []string{"scope1", "scope2"},
				Password: "password",
			},
			name: "flag user creation -- bad uid",
			args: args{
				c: context.Background(),
				s: []string{"user", "password"},
			},
			options: createUserOptionsDT{
				username: "",
				password: "password",
				scopes:   []string{"scope1", "scope2"},
				role:     12,
			},
			wantErr: true,
		},

		{
			usr: client.User{
				UID:      "this is a uid",
				Role:     12,
				Scope:    []string{"scope1", "scope2"},
				Password: "password",
			},
			name: "flag user creation -- bad password",
			args: args{
				c: context.Background(),
				s: []string{"user", "password"},
			},
			options: createUserOptionsDT{
				username: "this is a uid",
				password: "",
				scopes:   []string{"scope1", "scope2"},
				role:     12,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before()
			}
			if tt.after != nil {
				defer tt.after()
			}
			createUsrOptions = tt.options
			cl = mockCl{
				createUser: func(c context.Context, uid, pwd string, u client.User) error {
					if !reflect.DeepEqual(tt.usr, u) {
						t.Errorf("user is different than intended. %v\n!=\n%v", tt.usr, u)
					}
					if uid != tt.args.s[0] {
						t.Errorf("uid informed is different than intended %v\n!=\n%v", uid, tt.args.s[0])
					}

					if pwd != tt.args.s[1] {
						t.Errorf("password informed is different than intended %v\n!=\n%v", pwd, tt.args.s[1])
					}
					return nil
				},
			}
			if err := createUser(tt.args.c, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("createUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
