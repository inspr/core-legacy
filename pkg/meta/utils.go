package meta

import (
	"regexp"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
)

// StructureNameIsValid checks if the given name is valid for naming Channels, CTypes and dApps
func StructureNameIsValid(name string) (bool, error) {
	if len(name) == 0 || len(name) >= 64 {
		return false, ierrors.NewError().BadRequest().Message("invalid name length, must be (0 < length < 64)").Build()
	}
	qnameCharFmt := "[A-Za-z0-9]"
	qnameExtCharFmt := "[-A-Za-z0-9_]"
	qualifiedNameFmt := "(" + qnameCharFmt + qnameExtCharFmt + "*)?" + qnameCharFmt
	var r = regexp.MustCompile("^" + qualifiedNameFmt + "$")

	return r.MatchString(name), nil
}
