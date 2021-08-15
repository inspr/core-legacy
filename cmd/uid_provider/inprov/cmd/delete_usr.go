package cmd

import (
	"context"
	"fmt"

	build "inspr.dev/inspr/pkg/cmd"
)

type deleteUsrOptionsDT struct {
	username string
}

var deleteUsrOptions = deleteUsrOptionsDT{}

var deleteUserCmd = build.NewCmd("delete").WithDescription(
	"Delete a user from the Inspr UID provider",
).WithLongDescription(`
Delete command is responsible for executing the operation of deleting
an user in the InsprRedis instance in the cluster.

To execute this command one must specify the username of the account that
is going to be deleted and provide his credentials as well, the operation
will only be successful if the user provider has the permission to do so.
`).WithExample(
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
	err := cl.DeleteUser(
		ctx,
		inputArgs[0],
		inputArgs[1],
		deleteUsrOptions.username,
	)
	if err == nil {
		fmt.Println("Successfully deleted the user")
	}
	return err
}
