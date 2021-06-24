package cmd

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"inspr.dev/inspr/cmd/uid_provider/client"
	"inspr.dev/inspr/pkg/cmd"
	build "inspr.dev/inspr/pkg/cmd"
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
).WithFlags(
	&cmd.Flag{
		Name:      "output",
		Shorthand: "o",
		Usage:     "set the output file for the returned token",
		Value:     &loginOptions.output,
		DefValue:  "",
	},
	&cmd.Flag{
		Name:     "stdout",
		Usage:    "set the output of the token to stdout",
		Value:    &loginOptions.stdout,
		DefValue: false,
	},
).ExactArgs(2, loginAction)

func loginAction(c context.Context, s []string) error {

	var err error
	var output io.Writer
	var outputPath string
	if loginOptions.output == "" {
		f, _ := os.UserHomeDir()
		outputPath = filepath.Join(f, ".inspr", "token")

	} else {
		outputPath = loginOptions.output
	}

	if loginOptions.stdout {
		output = os.Stdout
	} else {
		output, err = os.Create(outputPath)
		if err != nil {
			return err
		}
	}
	return login(c, s[0], s[1], output)
}

func login(ctx context.Context, login, password string, output io.Writer) error {
	token, err := cl.Login(ctx, login, password)
	if err != nil {
		return err
	}

	output.Write([]byte(token))
	return nil
}
