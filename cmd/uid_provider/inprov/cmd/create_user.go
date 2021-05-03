package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/inspr/inspr/cmd/uid_provider/client"
	build "github.com/inspr/inspr/pkg/cmd"
	"gopkg.in/yaml.v2"
)

type createUserOptionsDT struct {
	username    string
	password    string
	scopes      []string
	yaml        string
	json        string
	permissions []string
}

var createUsrOptions = createUserOptionsDT{}

var createUserCmd = build.NewCmd(
	"create { -yaml || -json || -u USR | -p PWD | -s SCOPES } <username> <password>",
).WithDescription(
	"Creates a new user on the Insprd UID provider.",
).WithExample(
	"create a new user directly from the cli",
	"inprov create --username newUsername --password newPwd -s \"\" -s app1.app2 username password",
).WithExample(
	"create a new user directly from a YAML file",
	"inprov create --yaml user.yaml username password",
).WithExample(
	"create a new user directly from a json file",
	"inprov create --json user.json username password",
).WithFlags([]*build.Flag{
	{
		Name:      "username",
		Shorthand: "u",
		Usage:     "set the username of the user that will be created",
		Value:     &createUsrOptions.username,
		DefValue:  "",
	},
	{
		Name:      "password",
		Shorthand: "p",
		Usage:     "set the password of the user that will be created",
		Value:     &createUsrOptions.password,
		DefValue:  "",
	},
	{
		Name:      "scopes",
		Shorthand: "s",
		Usage:     "add a scope to the user permissions",
		Value:     &createUsrOptions.scopes,
		DefValue:  []string{},
	},
	{
		Name:     "yaml",
		Usage:    "read the user definition from a YAML file",
		Value:    &createUsrOptions.yaml,
		DefValue: "",
	},

	{
		Name:     "json",
		Usage:    "read the user definition from a JSON file",
		Value:    &createUsrOptions.json,
		DefValue: "",
	},
}).ExactArgs(2, createUser)

func createUser(ctx context.Context, inputArgs []string) error {
	var err error
	var usr client.User
	if createUsrOptions.yaml != "" {
		file, err := os.Open(createUsrOptions.yaml)
		if err != nil {
			return err
		}
		dec := yaml.NewDecoder(file)
		err = dec.Decode(&usr)
		if err != nil {
			return err
		}
	} else if createUsrOptions.json != "" {
		file, err := os.Open(createUsrOptions.json)
		if err != nil {
			return err
		}
		dec := json.NewDecoder(file)
		err = dec.Decode(&usr)
		if err != nil {
			return err
		}
	} else {
		usr.UID = createUsrOptions.username

		usr.Password = createUsrOptions.password

		usr.Permissions = createUsrOptions.permissions
		usr.Scope = createUsrOptions.scopes
	}

	if usr.UID == "" {
		return errors.New("username not informed")
	}
	if usr.Password == "" {
		return errors.New("password not informed")
	}

	err = cl.CreateUser(ctx, inputArgs[0], inputArgs[1], usr)
	return err
}
