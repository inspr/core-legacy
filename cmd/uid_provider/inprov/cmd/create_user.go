package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"gopkg.in/yaml.v2"
	"inspr.dev/inspr/cmd/uid_provider/client"
	"inspr.dev/inspr/pkg/cmd"
)

type createUserOptionsDT struct {
	username    string
	password    string
	yaml        string
	json        string
	permissions map[string][]string
}

var createUsrOptions = createUserOptionsDT{}

var createUserCmd = cmd.NewCmd(
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
).WithFlags(
	&cmd.Flag{
		Name:      "username",
		Shorthand: "u",
		Usage:     "set the username of the user that will be created",
		Value:     &createUsrOptions.username,
		DefValue:  "",
	},
	&cmd.Flag{
		Name:      "password",
		Shorthand: "p",
		Usage:     "set the password of the user that will be created",
		Value:     &createUsrOptions.password,
		DefValue:  "",
	},
	&cmd.Flag{
		Name:     "yaml",
		Usage:    "read the user definition from a YAML file",
		Value:    &createUsrOptions.yaml,
		DefValue: "",
	},
	&cmd.Flag{
		Name:     "json",
		Usage:    "read the user definition from a JSON file",
		Value:    &createUsrOptions.json,
		DefValue: "",
	},
).ExactArgs(2, createUser)

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
