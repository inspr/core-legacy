package utils

import (
	"regexp"
	"strings"

	"inspr.dev/inspr/pkg/ierrors"
)

// StructureNameIsValid checks if the given name is valid for naming Channels, types and dApps
func StructureNameIsValid(name string) error {
	if len(name) == 0 || len(name) >= 64 {
		return ierrors.New(
			"invalid name length, must be (0 < length < 64)",
		).BadRequest()
	}
	qnameCharFmt := "[A-Za-z0-9]"
	qnameExtCharFmt := "[-A-Za-z0-9_]"
	qualifiedNameFmt := "(" + qnameCharFmt + qnameExtCharFmt + "*)?" + qnameCharFmt
	var r = regexp.MustCompile("^" + qualifiedNameFmt + "$")

	if r.MatchString(name) {
		return nil
	}
	return ierrors.New(
		"invalid character in structure's name",
	).BadRequest()
}

// AliasNameIsValid checks if the given name is valid for naming aliasses
func AliasNameIsValid(name string) error {
	names := strings.Split(name, ".")
	if len(names) != 2 || names[len(names)-1] == "" {
		return ierrors.New("invalid alias name structure").BadRequest()
	}
	err := StructureNameIsValid(names[0])
	if err != nil {
		return err
	}

	return StructureNameIsValid(names[1])
}
