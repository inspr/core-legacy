package utils

import (
	"strings"

	"github.com/inspr/inspr/pkg/ierrors"
	"github.com/inspr/inspr/pkg/utils"
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

	if s2 == "" {
		return s1, nil
	}

	separator := ""
	if s1 != "" {
		separator = "."
	}
	newScope := s1 + separator + s2

	return newScope, nil
}

/*
RemoveAliasInScope removes the two last names defined in the scope
and returns the new scope and the alias that was removed
*/
func RemoveAliasInScope(scope string) (string, string, error) {
	if !IsValidScope(scope) {
		return "", "", ierrors.NewError().Message("invalid scope: %s", scope).InvalidName().Build()
	}

	names := strings.Split(scope, ".")
	aliasNames := names[len(names)-2:]
	names = names[:len(names)-2]

	alias := strings.Join(aliasNames, ".")
	newScope := strings.Join(names, ".")

	return newScope, alias, nil

}
