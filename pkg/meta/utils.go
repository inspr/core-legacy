package meta

import (
	"regexp"
	"strings"

	kubeCore "k8s.io/api/core/v1"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
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

// StringArray is an array of strings with functional and set-like helper methods
type StringArray []string

// Map maps a given function into another string array
func (c StringArray) Map(f func(string) string) StringArray {
	return utils.Map(c, f)
}

// Union returns the union of a string array with another
func (c StringArray) Union(other StringArray) StringArray {
	return utils.StringSliceUnion(c, other)
}

// Contains returns whether or not an array contains an item
func (c StringArray) Contains(item string) bool {
	return utils.Includes(c, item)
}

// Join joins a string array with a given separator, returning the string generated
func (c StringArray) Join(sep string) string {
	return strings.Join(c, sep)
}

// EnvironmentMap is a type for environment variables represented as a map
type EnvironmentMap map[string]string

// ParseToK8sArrEnv parses the map into an array of kubernetes' environment variables
func (m EnvironmentMap) ParseToK8sArrEnv() []kubeCore.EnvVar {
	var arrEnv []kubeCore.EnvVar
	for key, val := range m {
		arrEnv = append(arrEnv, kubeCore.EnvVar{
			Name:  key,
			Value: val,
		})
	}
	return arrEnv
}
