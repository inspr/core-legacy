package utils

import (
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
)

//GetScope retreives a path for use as base on insprd request.
//Takes into consideration viper config and scope flag.
func GetScope() (string, error) {
	defaultScope := GetConfiguredScope()
	scope := defaultScope

	if cmd.InsprOptions.Scope != "" {
		if utils.IsValidScope(cmd.InsprOptions.Scope) {
			scope = cmd.InsprOptions.Scope
		} else {
			return "", ierrors.
				NewError().
				BadRequest().
				Message("'%v' is an invalid scope", cmd.InsprOptions.Scope).
				Build()
		}
	}

	return scope, nil
}
