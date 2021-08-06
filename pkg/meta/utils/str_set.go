package utils

import (
	"reflect"

	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/utils"
)

const exists = true

//MApps is a string->App map
type MApps map[string]*meta.App

// MAliases is a string->Alias map
type MAliases map[string]*meta.Alias

//MChannels is a string->channel map
type MChannels map[string]*meta.Channel

//MTypes is a string->Type map
type MTypes map[string]*meta.Type

//MStr is a string->string map
type MStr map[string]string

//StrSet implements a set of strings.
type StrSet map[string]bool

//MakeStrSet creates a StrSet from a supported type.
func MakeStrSet(obj interface{}) (StrSet, error) {
	set := make(StrSet)

	switch value := obj.(type) {
	case MApps:
		for k := range value {
			set[k] = exists
		}
		return set, nil

	case MChannels:
		for k := range value {
			set[k] = exists
		}
		return set, nil

	case MTypes:
		for k := range value {
			set[k] = exists
		}
		return set, nil

	case MStr:
		for k := range value {
			set[k] = exists
		}
		return set, nil

	case []string:
		for _, str := range value {
			set[str] = exists
		}
		return set, nil

	case StrSet:
		for k := range value {
			set[k] = exists
		}
		return set, nil

	case utils.StringArray:
		for _, str := range value {
			set[str] = exists
		}
		return set, nil

	case map[string]string:
		for k := range value {
			set[k] = exists
		}
		return set, nil
	case MAliases:
		for k := range value {
			set[k] = exists
		}
		return set, nil

	default:
		objType := reflect.TypeOf(obj)
		return nil, ierrors.New(
			"error while making set: '" + objType.Name() + "' type not supported",
		).InternalServer()
	}
}

//AppendSet extends a StrSet with another StrSet.
func (set *StrSet) AppendSet(target StrSet) {
	for k := range target {
		(*set)[k] = exists
	}
}

//ToArray returns an array of all items in the set.
func (set *StrSet) ToArray() utils.StringArray {
	arr := utils.StringArray{}
	for k := range *set {
		arr = append(arr, k)
	}
	return arr
}

//DisjunctSet returns the disjunction set between two StrSet.
func DisjunctSet(set1, set2 StrSet) StrSet {
	setFinal := make(StrSet)

	setTemp, _ := MakeStrSet(set1)
	setTemp.AppendSet(set2)

	for k := range setTemp {
		_, first := set1[k]
		_, second := set2[k]
		if first != second {
			setFinal[k] = exists
		}
	}

	return setFinal
}

//IntersectSet returns the intersection set between two StrSet.
func IntersectSet(set1, set2 StrSet) StrSet {
	setFinal := make(StrSet)

	setTemp, _ := MakeStrSet(set1)
	setTemp.AppendSet(set2)

	for k := range setTemp {
		_, first := set1[k]
		_, second := set2[k]
		if first && second {
			setFinal[k] = exists
		}
	}

	return setFinal
}
