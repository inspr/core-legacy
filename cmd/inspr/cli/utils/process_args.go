package utils

import (
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
)

func ProcessArg(arg, scope string) (string, string, error) {
	path := scope
	var component string

	if err := utils.StructureNameIsValid(arg); err != nil {
		if !utils.IsValidScope(arg) {
			return "", "", ierrors.NewError().Message("invalid scope").BadRequest().Build()
		}

		newScope, lastName, err := utils.RemoveLastPartInScope(arg)
		if err != nil {
			return "", "", ierrors.NewError().Message("invalid scope").BadRequest().Build()
		}

		separator := ""
		if scope != "" {
			separator = "."
		}

		path = path + separator + newScope
		component = lastName

	} else {
		component = arg
	}
	return path, component, nil
}
