package cmd

import (
	"context"

	build "inspr.dev/inspr/pkg/cmd"
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
).WithFlags(
	&build.Flag{
		Name:     "username",
		Usage:    "username of the user to be deleted",
		Value:    &deleteUsrOptions.username,
		DefValue: "",
	},
).ExactArgs(2, deleteAction)

func deleteAction(ctx context.Context, inputArgs []string) error {

	return cl.DeleteUser(ctx, inputArgs[0], inputArgs[1], deleteUsrOptions.username)
}
