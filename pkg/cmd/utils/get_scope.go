package utils

import (
	"strings"

	"inspr.dev/inspr/pkg/cmd"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta/utils"
)

// GetScope retreives the path to be used as base scope for an Insprd request.
// Takes into consideration viper config and scope flag.
func GetScope() (string, error) {
	defaultScope := GetConfiguredScope()
	scope := defaultScope

	if cmd.InsprOptions.Scope != "" {
		cmd.InsprOptions.Scope = strings.TrimSuffix(cmd.InsprOptions.Scope, ".")
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
