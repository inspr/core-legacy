package cmd

import (
	"context"
	"io"
	"os"

	"github.com/inspr/inspr/cmd/uid_provider/client"
	build "github.com/inspr/inspr/pkg/cmd"
)

var cl client.UIDClient

type loginOptionsDT struct {
	output string
	stdout bool
}

var loginOptions = loginOptionsDT{}

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
}).ExactArgs(2, login)

func login(c context.Context, s []string) error {

	var err error
	var output io.Writer

	if loginOptions.stdout {
		output = os.Stdout
	} else {
		output, err = os.OpenFile(loginOptions.output, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return doStuff(c, s[0], s[1], output)
}

func doStuff(ctx context.Context, login, password string, output io.Writer) error {
	token, err := cl.Login(ctx, login, password)
	if err != nil {
		return err
	}

	output.Write([]byte(token))
	return nil
}
