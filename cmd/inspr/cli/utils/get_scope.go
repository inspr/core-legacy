package utils

import (
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"
)

func GetScope() (string, error) {
	defaultScope := GetConfiguredScope()
	scope := defaultScope

	if cmd.InsprOptions.Scope != "" {
		if utils.IsValidScope(cmd.InsprOptions.Scope) {
			scope = cmd.InsprOptions.Scope
		} else {
			return "", ierrors.NewError().BadRequest().Message("invalid scope").Build()
		}
	}

	return scope, nil
}
