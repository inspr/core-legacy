package utils

import (
	"reflect"

	"github.com/google/go-cmp/cmp"
)

var (
	evaluatesMap = cmp.Comparer(func(l, r interface{}) bool {
		mapLeft := reflect.ValueOf(l)
		mapRight := reflect.ValueOf(r)

		if mapLeft.Len() != mapRight.Len() {
			return false
		}

		for _, e := range mapLeft.MapKeys() {
			a := mapLeft.MapIndex(e)
			b := mapRight.MapIndex(e)
			if a.String() != b.String() {
				return false
			}
		}
		return true
	})

	alwaysTrue = cmp.Comparer(func(_, _ interface{}) bool { return true })
)

// GetMapCompareOptions - returns opts for cmd.Equal
// ignores all slices and compares maps
func GetMapCompareOptions() cmp.Options {

	opts := cmp.Options{
		cmp.FilterValues(func(x, y interface{}) bool {
			vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
			flag := (vx.IsValid() && vy.IsValid() && vx.Type() == vy.Type()) && (vx.Kind() == reflect.Map)
			return flag
		}, evaluatesMap),

		cmp.FilterValues(func(x, y interface{}) bool {
			vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
			flag := (vx.IsValid() && vy.IsValid() && vx.Type() == vy.Type()) && (vx.Kind() == reflect.Slice)
			return flag
		}, alwaysTrue),
	}
	return opts
}

// GeneralCompareOptions - returns opts for cmd.Equal
// ignores all slices and maps inside a struct.
//
// While not being ideal it removes the maps and slices from the reflect.Equal.
func GeneralCompareOptions() cmp.Options {
	opts := cmp.Options{
		cmp.FilterValues(func(x, y interface{}) bool {
			vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
			return (vx.IsValid() && vy.IsValid() && vx.Type() == vy.Type()) &&
				(vx.Kind() == reflect.Slice || vx.Kind() == reflect.Map)
		}, alwaysTrue),
	}
	return opts
}
