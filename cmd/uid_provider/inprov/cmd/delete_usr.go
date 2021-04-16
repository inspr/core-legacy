package cmd

import (
	"context"

	build "github.com/inspr/inspr/pkg/cmd"
)

type deleteUsrOptionsDT struct {
	username string
}

var deleteUsrOptions = deleteUsrOptionsDT{}

var deleteUserCmd = build.NewCmd("delete").WithDescription(
	"Delete a user from the Inspr UID provider",
).WithExample(
	"delete a user given credentials",
	"inprov delete --username userToBeDeleted username password",
).WithFlags([]*build.Flag{
	{
		Name:     "username",
		Usage:    "username of the user to be deleted",
		Value:    &deleteUsrOptions.username,
		DefValue: "",
	},
}).ExactArgs(2, deleteAction)

func deleteAction(c context.Context, s []string) error {

	return cl.DeleteUser(c, s[0], s[1], deleteUsrOptions.username)
}
