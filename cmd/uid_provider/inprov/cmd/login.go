package cmd

import (
	"context"
	"io"
	"os"

	"github.com/inspr/inspr/cmd/uid_provider/client"
	build "github.com/inspr/inspr/pkg/cmd"
)

var cl client.UIDClient

var loginOptions struct {
	output string
	stdout bool
} = struct {
	output string
	stdout bool
}{}

var loginCmd = build.NewCmd("login").WithDescription(
	"Log in to the Inspr UID provider and get a token.",
).WithExample(
	"log in with your user and password",
	"inprov login usr pwd",
).WithFlags([]*build.Flag{
	{
		Name:      "output",
		Shorthand: "o",
		Usage:     "set the output file for the returned token",
		Value:     &loginOptions.output,
		DefValue:  "",
	},
	{
		Name:     "stdout",
		Usage:    "set the output of the token to stdout",
		Value:    &loginOptions.stdout,
		DefValue: false,
	},
}).ExactArgs(2, func(c context.Context, s []string) error {

	var err error
	var output io.Writer
	if loginOptions.stdout {
		output = os.Stdout
	} else {
		output, err = os.OpenFile(loginOptions.output, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			return err
		}
	}

	token, err := cl.Login(c, s[0], s[1])
	if err != nil {
		return err
	}

	_, err = output.Write([]byte(token))
	return err
})
