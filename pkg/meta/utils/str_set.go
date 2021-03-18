package utils

import "gitlab.inspr.dev/inspr/core/pkg/meta"

var exists = true

//MApps is a string->App map
type MApps map[string]*meta.App

// MAliases is a string->Alias map
type MAliases map[string]*meta.Alias

//MChannels is a string->channel map
type MChannels map[string]*meta.Channel

//MTypes is a string->ChannelType map
type MTypes map[string]*meta.ChannelType

//MStr is a string->string map
type MStr map[string]string

//StrSet implements a set of strings.
type StrSet map[string]bool

//ArrMakeSet creates a StrSet from a []string.
func ArrMakeSet(strings []string) StrSet {
	set := make(StrSet)
	for _, str := range strings {
		set[str] = exists
	}
	return set
}

//ArrAppendSet extends a StrSet with a []string.
func (set *StrSet) ArrAppendSet(strings []string) {
	for _, str := range strings {
		(*set)[str] = exists
	}
}

//AppendSet extends a StrSeet with a string
func (set *StrSet) AppendSet(str string) {
	(*set)[str] = exists
}

//ArrDisjuncSet returns the disjunction set between two []string.
func ArrDisjuncSet(arr1 []string, arr2 []string) StrSet {
	set := make(StrSet)
	setTemp := ArrMakeSet(arr1)
	set1 := ArrMakeSet(arr1)
	set2 := ArrMakeSet(arr2)
	setTemp.ArrAppendSet(arr2)

	for k := range setTemp {
		first := set1[k]
		second := set2[k]
		if first != second {
			set[k] = exists
		}
	}

	return set
}

//AppMakeSet creates a StrSet from a Apps map.
func AppMakeSet(apps MApps) StrSet {
	set := make(StrSet)

	for k := range apps {
		set[k] = exists
	}

	return set
}

//AppAppendSet extends a StrSet with a Apps map.
func (set *StrSet) AppAppendSet(apps MApps) {
	for k := range apps {
		(*set)[k] = exists
	}
}

//AliasesMakeSet creates a StrSet from a Aliasess map.
func AliasesMakeSet(apps MAliases) StrSet {
	set := make(StrSet)

	for k := range apps {
		set[k] = exists
	}

	return set
}

//AliasAppendSet extends a StrSet with a Apps map.
func (set *StrSet) AliasAppendSet(apps MAliases) {
	for k := range apps {
		(*set)[k] = exists
	}
}

//AppDisjuncSet returns the disjunction set between two Apps maps.
func AppDisjuncSet(apps1 MApps, apps2 MApps) StrSet {
	set := make(StrSet)

	setTemp := AppMakeSet(apps1)
	setTemp.AppAppendSet(apps2)

	for k := range setTemp {
		_, first := apps1[k]
		_, second := apps2[k]
		if first != second {
			set[k] = exists
		}
	}

	return set
}

//AliasDisjuctSet returns the disjunction set between two Apps maps.
func AliasDisjuctSet(apps1 MAliases, apps2 MAliases) StrSet {
	set := make(StrSet)

	setTemp := AliasesMakeSet(apps1)
	setTemp.AliasAppendSet(apps2)

	for k := range setTemp {
		_, first := apps1[k]
		_, second := apps2[k]
		if first != second {
			set[k] = exists
		}
	}

	return set
}

//AliasIntersecSet returns the intersection set between two Aliases maps.
func AliasIntersecSet(apps1 MAliases, apps2 MAliases) StrSet {
	set := make(StrSet)

	setTemp := AliasesMakeSet(apps1)
	setTemp.AliasAppendSet(apps2)

	for k := range setTemp {
		_, first := apps1[k]
		_, second := apps2[k]
		if first && second {
			set[k] = exists
		}
	}

	return set
}

//AppIntersecSet returns the intersection set between two Apps maps.
func AppIntersecSet(apps1 MApps, apps2 MApps) StrSet {
	set := make(StrSet)

	setTemp := AppMakeSet(apps1)
	setTemp.AppAppendSet(apps2)

	for k := range setTemp {
		_, first := apps1[k]
		_, second := apps2[k]
		if first && second {
			set[k] = exists
		}
	}

	return set
}

//ChsMakeSet creates a StrSet from a Channels map.
func ChsMakeSet(chans MChannels) StrSet {
	set := make(StrSet)

	for k := range chans {
		set[k] = exists
	}

	return set
}

//ChsAppendSet extends a StrSet with a Channels map.
func (set *StrSet) ChsAppendSet(chans MChannels) {
	for k := range chans {
		(*set)[k] = exists
	}
}

//ChsDisjuncSet returns the disjunction set between two Channels maps.
func ChsDisjuncSet(chans1 MChannels, chans2 MChannels) StrSet {
	set := make(StrSet)

	setTemp := ChsMakeSet(chans1)
	setTemp.ChsAppendSet(chans2)

	for k := range setTemp {
		_, first := chans1[k]
		_, second := chans2[k]
		if first != second {
			set[k] = exists
		}
	}

	return set
}

//ChsIntersecSet returns the intersection set between two Channels maps.
func ChsIntersecSet(chans1 MChannels, chans2 MChannels) StrSet {
	set := make(StrSet)

	setTemp := ChsMakeSet(chans1)
	setTemp.ChsAppendSet(chans2)

	for k := range setTemp {
		_, first := chans1[k]
		_, second := chans2[k]
		if first && second {
			set[k] = exists
		}
	}

	return set
}

//TypesMakeSet creates a StrSet from a Types map.
func TypesMakeSet(types MTypes) StrSet {
	set := make(StrSet)

	for k := range types {
		set[k] = exists
	}

	return set
}

//TypesAppendSet extends a StrSet with a Types map.
func (set *StrSet) TypesAppendSet(types MTypes) {
	for k := range types {
		(*set)[k] = exists
	}
}

//TypesDisjuncSet returns the disjunction set between two Types maps.
func TypesDisjuncSet(types1 MTypes, types2 MTypes) StrSet {
	set := make(StrSet)

	setTemp := TypesMakeSet(types1)
	setTemp.TypesAppendSet(types2)

	for k := range setTemp {
		_, first := types1[k]
		_, second := types2[k]
		if first != second {
			set[k] = exists
		}
	}

	return set
}

//TypesIntersecSet returns the intersection set between two Types maps.
func TypesIntersecSet(types1 MTypes, types2 MTypes) StrSet {
	set := make(StrSet)

	setTemp := TypesMakeSet(types1)
	setTemp.TypesAppendSet(types2)

	for k := range setTemp {
		_, first := types1[k]
		_, second := types2[k]
		if first && second {
			set[k] = exists
		}
	}

	return set
}

//StrMakeSet creates a StrSet from a Str map.
func StrMakeSet(strings MStr) StrSet {
	set := make(StrSet)

	for k := range strings {
		set[k] = exists
	}

	return set
}

//StrAppendSet extends a StrSet with a Str map.
func (set *StrSet) StrAppendSet(strings MStr) {
	for k := range strings {
		(*set)[k] = exists
	}
}

//StrDisjuncSet returns the disjunction set between two Str maps.
func StrDisjuncSet(strings1 MStr, strings2 MStr) StrSet {
	set := make(StrSet)

	setTemp := StrMakeSet(strings1)
	setTemp.StrAppendSet(strings2)

	for k := range setTemp {
		_, first := strings1[k]
		_, second := strings2[k]
		if first != second {
			set[k] = exists
		}
	}

	return set
}
