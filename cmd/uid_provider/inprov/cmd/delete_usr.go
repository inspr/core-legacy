package cmd

import (
	"context"

	build "github.com/inspr/inspr/pkg/cmd"
)

var deleteUsrOptions struct {
	username string
} = struct {
	username string
}{}

var deleteUserCmd = build.NewCmd("delete").WithDescription(
	"Delete a user from the Inspr UID provider",
).WithExample(
	"delete a user given credentials",
	"inprov login --username userToBeDeleted username password",
).WithFlags([]*build.Flag{
	{
		Name:     "username",
		Usage:    "username of the user to be deleted",
		Value:    &deleteUsrOptions.username,
		DefValue: "",
	},
}).ExactArgs(2, func(c context.Context, s []string) error {

	return cl.DeleteUser(c, s[0], s[1], deleteUsrOptions.username)
})
