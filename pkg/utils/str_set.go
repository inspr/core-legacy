package utils

import "gitlab.inspr.dev/inspr/core/pkg/meta"

var exists = true

//Apps is a string->App map
type Apps map[string]*meta.App

//Channels is a string->channel map
type Channels map[string]*meta.Channel

//Types is a string->ChannelType map
type Types map[string]*meta.ChannelType

//Str is a string->string map
type Str map[string]string

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
func AppMakeSet(apps Apps) StrSet {
	set := make(StrSet)

	for k := range apps {
		set[k] = exists
	}

	return set
}

//AppAppendSet extends a StrSet with a Apps map.
func (set *StrSet) AppAppendSet(apps Apps) {
	for k := range apps {
		(*set)[k] = exists
	}
}

//AppDisjuncSet returns the disjunction set between two Apps maps.
func AppDisjuncSet(apps1 Apps, apps2 Apps) StrSet {
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

//AppIntersecSet returns the intersection set between two Apps maps.
func AppIntersecSet(apps1 Apps, apps2 Apps) StrSet {
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
func ChsMakeSet(apps Channels) StrSet {
	set := make(StrSet)

	for k := range apps {
		set[k] = exists
	}

	return set
}

//ChsAppendSet extends a StrSet with a Channels map.
func (set *StrSet) ChsAppendSet(apps Channels) {
	for k := range apps {
		(*set)[k] = exists
	}
}

//ChsDisjuncSet returns the disjunction set between two Channels maps.
func ChsDisjuncSet(apps1 Channels, apps2 Channels) StrSet {
	set := make(StrSet)

	setTemp := ChsMakeSet(apps1)
	setTemp.ChsAppendSet(apps2)

	for k := range setTemp {
		_, first := apps1[k]
		_, second := apps2[k]
		if first != second {
			set[k] = exists
		}
	}

	return set
}

//ChsIntersecSet returns the intersection set between two Channels maps.
func ChsIntersecSet(apps1 Channels, apps2 Channels) StrSet {
	set := make(StrSet)

	setTemp := ChsMakeSet(apps1)
	setTemp.ChsAppendSet(apps2)

	for k := range setTemp {
		_, first := apps1[k]
		_, second := apps2[k]
		if first && second {
			set[k] = exists
		}
	}

	return set
}

//TypesMakeSet creates a StrSet from a Types map.
func TypesMakeSet(types Types) StrSet {
	set := make(StrSet)

	for k := range types {
		set[k] = exists
	}

	return set
}

//TypesAppendSet extends a StrSet with a Types map.
func (set *StrSet) TypesAppendSet(types Types) {
	for k := range types {
		(*set)[k] = exists
	}
}

//TypesDisjuncSet returns the disjunction set between two Types maps.
func TypesDisjuncSet(types1 Types, types2 Types) StrSet {
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
func TypesIntersecSet(types1 Types, types2 Types) StrSet {
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
func StrMakeSet(strings Str) StrSet {
	set := make(StrSet)

	for k := range strings {
		set[k] = exists
	}

	return set
}

//StrAppendSet extends a StrSet with a Str map.
func (set *StrSet) StrAppendSet(strings Str) {
	for k := range strings {
		(*set)[k] = exists
	}
}

//StrDisjuncSet returns the disjunction set between two Str maps.
func StrDisjuncSet(strings1 Str, strigns2 Str) StrSet {
	set := make(StrSet)

	setTemp := StrMakeSet(strings1)
	setTemp.StrAppendSet(strigns2)

	for k := range setTemp {
		_, first := strings1[k]
		_, second := strigns2[k]
		if first != second {
			set[k] = exists
		}
	}

	return set
}
