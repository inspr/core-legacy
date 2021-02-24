package utils

import (
	"strings"

	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

/*
IsValidScope checks if the given scope is of the type
'name1.name2.name3'
*/
func IsValidScope(scope string) bool {
	if len(scope) == 0 {
		return true
	}

	names := strings.Split(scope, ".")
	if utils.Includes(names, "") {
		return false
	}

	for _, name := range names {
		if err := StructureNameIsValid(name); err != nil {
			return false
		}
	}

	return true
}
