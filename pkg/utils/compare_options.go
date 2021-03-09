package utils

import (
	"fmt"
	"reflect"

	"github.com/google/go-cmp/cmp"
)

// GetMapCompareOptions - returns opts for cmd.Equal
// ignores all slices and compares maps
func GetMapCompareOptions() cmp.Options {
	evaluatesMap := cmp.Comparer(func(l, r interface{}) bool {
		mapLeft := reflect.ValueOf(l)
		mapRight := reflect.ValueOf(r)

		if mapLeft.Len() != mapRight.Len() {
			return false
		}

		for _, e := range mapLeft.MapKeys() {
			a := mapLeft.MapIndex(e)
			b := mapRight.MapIndex(e)
			fmt.Println(a, b)
			if a.String() != b.String() {
				return false
			}
		}
		return true
	})

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
		}, cmp.Comparer(func(_, _ interface{}) bool { return true })),
	}
	return opts
}
