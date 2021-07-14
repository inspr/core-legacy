package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/term"
	"inspr.dev/inspr/cmd/uid_provider/client"
	"inspr.dev/inspr/pkg/cmd"
)

var cl client.UIDClient

type loginOptionsDT struct {
	output   string
	stdout   bool
	user     string
	password string
}

var loginOptions = loginOptionsDT{}

var loginCmd = cmd.NewCmd("login").WithDescription(
	"Log in to the Inspr UID provider and get a token.",
).WithLongDescription(`
login is the command responsible for associating the insprctl operations
with an account on the UID Provider on the cluster.
`).WithExample(
	"log in with your user and password",
	"inprov login",
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
	&cmd.Flag{
		Name:      "user",
		Shorthand: "u",
		Usage:     "set the user for the login",
		Value:     &loginOptions.user,
		DefValue:  "",
	},
	&cmd.Flag{
		Name:      "password",
		Shorthand: "p",
		Usage:     "set the user's password for the login",
		Value:     &loginOptions.password,
		DefValue:  "",
	},
).NoArgs(loginAction)

func loginAction(c context.Context) error {

	if loginOptions.user == "" {
		if loginOptions.password != "" {
			return fmt.Errorf("invalid user")
		}
		fmt.Print("Username: ")
		fmt.Scanln(&loginOptions.user)
	}
	if loginOptions.password == "" {
		fmt.Print("Password: ")
		password, err := term.ReadPassword(0)
		fmt.Println()
		if err != nil {
			return err
		}
		loginOptions.password = string(password)
	}

	return login(c, loginOptions.user, loginOptions.password)
}

func login(ctx context.Context, login, password string) error {
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
	token, err := cl.Login(ctx, login, password)
	if err != nil {
		return err
	}

	fmt.Println("Successfully logged in!")

	output.Write([]byte(token))
	return nil
}
