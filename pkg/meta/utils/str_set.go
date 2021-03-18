package utils

import (
	"reflect"

	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

const exists = true

//MApps is a string->App map
type MApps map[string]*meta.App

//MChannels is a string->channel map
type MChannels map[string]*meta.Channel

//MTypes is a string->ChannelType map
type MTypes map[string]*meta.ChannelType

//MStr is a string->string map
type MStr map[string]string

//StrSet implements a set of strings.
type StrSet map[string]bool

//MakeStrSet creates a StrSet from a supported type.
func MakeStrSet(obj interface{}) (StrSet, error) {
	set := make(StrSet)

	switch obj.(type) {
	case MApps:
		value := obj.(MApps)
		for k := range value {
			set[k] = exists
		}
		return set, nil

	case MChannels:
		value := obj.(MChannels)
		for k := range value {
			set[k] = exists
		}
		return set, nil

	case MTypes:
		value := obj.(MTypes)
		for k := range value {
			set[k] = exists
		}
		return set, nil

	case MStr:
		value := obj.(MStr)
		for k := range value {
			set[k] = exists
		}
		return set, nil

	case []string:
		value := obj.([]string)
		for _, str := range value {
			set[str] = exists
		}
		return set, nil

	case StrSet:
		value := obj.(StrSet)
		for k := range value {
			set[k] = exists
		}
		return set, nil

	case utils.StringArray:
		value := obj.(utils.StringArray)
		for _, str := range value {
			set[str] = exists
		}
		return set, nil

	case map[string]string:
		value := obj.(map[string]string)
		for k := range value {
			set[k] = exists
		}
		return set, nil

	default:
		objType := reflect.TypeOf(obj)
		return nil, ierrors.NewError().
			InternalServer().
			Message("error while making set: '" + objType.Name() + "' type not supported").
			Build()
	}
}

//AppendSet extends a StrSet with another StrSet.
func (set *StrSet) AppendSet(target StrSet) {
	for k := range target {
		(*set)[k] = exists
	}
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
