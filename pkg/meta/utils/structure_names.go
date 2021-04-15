package utils

import (
	"regexp"

	"github.com/inspr/inspr/pkg/ierrors"
)

// StructureNameIsValid checks if the given name is valid for naming Channels, CTypes and dApps
func StructureNameIsValid(name string) error {
	if len(name) == 0 || len(name) >= 64 {
		return ierrors.NewError().BadRequest().Message("invalid name length, must be (0 < length < 64)").Build()
	}
	qnameCharFmt := "[A-Za-z0-9]"
	qnameExtCharFmt := "[-A-Za-z0-9_]"
	qualifiedNameFmt := "(" + qnameCharFmt + qnameExtCharFmt + "*)?" + qnameCharFmt
	var r = regexp.MustCompile("^" + qualifiedNameFmt + "$")

	if r.MatchString(name) {
		return nil
	}
	return ierrors.NewError().BadRequest().Message("invalid character in structure's name").Build()
}
