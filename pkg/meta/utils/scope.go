package utils

import (
	"strings"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
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

/*
RemoveLastPartInScope removes the last name defined in the scope
and returns the new scope and the element that was removed
*/
func RemoveLastPartInScope(scope string) (string, string, error) {
	if !IsValidScope(scope) {
		return "", "", ierrors.NewError().Message("invalid scope: " + scope).InvalidName().Build()
	}

	names := strings.Split(scope, ".")
	lastName := names[len(names)-1]
	names = names[:len(names)-1]

	newScope := strings.Join(names, ".")

	return newScope, lastName, nil

}

/*
JoinScopes join two scopes and return the new scope
*/
func JoinScopes(s1, s2 string) (string, error) {
	if !IsValidScope(s1) || !IsValidScope(s2) {
		return "", ierrors.NewError().Message("invalid scope in args").InvalidName().Build()
	}

	separator := ""
	if s1 != "" {
		separator = "."
	}
	newScope := s1 + separator + s2

	if !IsValidScope(newScope) {
		return "", ierrors.NewError().Message("invalid result scope. Make sure the second scope is not the root (empty)").InvalidName().Build()
	}

	return newScope, nil
}
