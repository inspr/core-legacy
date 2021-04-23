package utils

import (
	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/meta/utils"
)

//ProcessArg is responsible for separating a path into an component name and it's parent's path.
// < path, name, error >
func ProcessArg(arg, scope string) (string, string, error) {
	path := scope
	var component string

	if err := utils.StructureNameIsValid(arg); err != nil {
		if !utils.IsValidScope(arg) {
			return "", "", ierrors.NewError().Message("invalid scope").BadRequest().Build()
		}

		newScope, lastName, _ := utils.RemoveLastPartInScope(arg)
		path, _ = utils.JoinScopes(path, newScope)

		component = lastName
	} else {
		component = arg
	}
	return path, component, nil
}