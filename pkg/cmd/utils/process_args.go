package utils

import (
	"fmt"

	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta/utils"
)

// CheckEmptyArgs receives the args of a cli command and returns a error in case any
// of them are empty
func CheckEmptyArgs(args map[string]string) error {
	var err error = nil
	for k, v := range args {
		if v == "" {
			errorMessage := fmt.Sprintf("arg '%v' is empty", k)
			if err == nil {
				err = ierrors.New(errorMessage).InvalidArgs()
			} else {
				err = ierrors.Wrap(err, errorMessage)
			}
		}
	}
	return err
}

//ProcessArg is responsible for separating a path into an component name and it's parent's path.
// < path, name, error >
func ProcessArg(arg, scope string) (string, string, error) {
	path := scope
	var component string

	if err := utils.StructureNameIsValid(arg); err != nil {
		if !utils.IsValidScope(arg) {
			return "", "", ierrors.New("invalid scope").BadRequest()
		}

		newScope, lastName, _ := utils.RemoveLastPartInScope(arg)
		path, _ = utils.JoinScopes(path, newScope)

		component = lastName
	} else {
		component = arg
	}
	return path, component, nil
}

//ProcessAliasArg is responsible for separating a path into an alias name and it's parent's path.
// < path, name, error >
func ProcessAliasArg(arg, scope string) (string, string, error) {
	path := scope
	var alias string

	if err := utils.AliasNameIsValid(arg); err != nil {
		if !utils.IsValidScope(arg) {
			return "", "", ierrors.New("invalid scope").BadRequest()
		}

		newScope, lastName, _ := utils.RemoveAliasInScope(arg)
		path, _ = utils.JoinScopes(path, newScope)

		alias = lastName
	} else {
		alias = arg
	}
	return path, alias, nil
}
