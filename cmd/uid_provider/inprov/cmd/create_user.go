package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"gitlab.inspr.dev/inspr/core/cmd/uid_provider/client"
	build "gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gopkg.in/yaml.v2"
)

var createUsrOptions struct {
	username string
	password string
	scopes   []string
	yaml     string
	json     string
	role     int
} = struct {
	username string
	password string
	scopes   []string
	yaml     string
	json     string
	role     int
}{}

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
}).ExactArgs(2, func(c context.Context, s []string) error {

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
		file, err := os.Open(createUsrOptions.yaml)
		if err != nil {
			return err
		}
		dec := json.NewDecoder(file)
		err = dec.Decode(&usr)
		if err != nil {
			return err
		}
	} else {
		if createUsrOptions.username == "" {
			return errors.New("username not informed")
		}
		usr.UID = createUsrOptions.username

		if createUsrOptions.password == "" {
			return errors.New("password not informed")
		}
		usr.Password = createUsrOptions.password

		usr.Role = createUsrOptions.role
		usr.Scope = createUsrOptions.scopes
	}

	err = cl.CreateUser(c, s[0], usr)
	return err
})
