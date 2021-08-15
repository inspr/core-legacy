package utils

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetMapCompareOptions(t *testing.T) {
	tests := []struct {
		name string
		want cmp.Options
	}{
		{
			name: "testing if returns the right amount of filters",
			want: cmp.Options{
				cmp.FilterValues(func(x, y interface{}) bool {
					vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
					flag := (vx.IsValid() && vy.IsValid() && vx.Type() == vy.Type()) &&
						(vx.Kind() == reflect.Map)
					return flag
				}, cmp.Comparer(func(l, r interface{}) bool { return true })),

				cmp.FilterValues(func(x, y interface{}) bool {
					vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
					flag := (vx.IsValid() && vy.IsValid() && vx.Type() == vy.Type()) &&
						(vx.Kind() == reflect.Slice)
					return flag
				}, cmp.Comparer(func(l, r interface{}) bool { return true })),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := GetMapCompareOptions()
			if len(got) != len(tt.want) {
				t.Errorf("GetMapCompareOptions() - length = %v, want %v",
					len(got),
					len(tt.want),
				)
			}

			// maps comparison
			mapA := map[string]string{"a": "something"}
			mapB := map[string]string{"b": "something"}
			if cmp.Equal(mapA, mapB, got) != false {
				t.Errorf("comparing MAPS, GetMapcompare => %v, expect %v",
					cmp.Equal(mapA, mapB, got),
					false,
				)
			}

			mapA = map[string]string{"a": "something"}
			mapB = map[string]string{"a": "something"}
			if cmp.Equal(mapA, mapB, got) != true {
				t.Errorf("comparing MAPS, GetMapcompare => %v, expect %v",
					cmp.Equal(mapA, mapB, got),
					true,
				)
			}

			// slices should always come as true
			sliceA := []int{1, 2, 3}
			sliceB := []int{3, 2, 6}
			if cmp.Equal(
				sliceA,
				sliceB,
				got,
			) != cmp.Equal(
				sliceA,
				sliceB,
				tt.want,
			) {
				t.Errorf("comparing SLICES, GetMapcompare => %v, expect %v",
					cmp.Equal(mapA, mapB, got),
					cmp.Equal(mapA, mapB, tt.want),
				)
			}

			// should be the same result as reflect.DeepEqual()
			if (cmp.Equal(1, 2, got) != cmp.Equal(1, 2, tt.want)) ||
				(cmp.Equal(1, 2, got) != reflect.DeepEqual(1, 2)) {

				t.Errorf("comparing NOT maps, GetMapcompare => %v, expect %v",
					cmp.Equal(mapA, mapB, got),
					cmp.Equal(mapA, mapB, tt.want),
				)
			}

		})
	}
}

func TestGeneralCompareOptions(t *testing.T) {
	tests := []struct {
		name  string
		left  interface{}
		right interface{}
		want  bool
	}{
		{
			name:  "basic_values_cmp_equal",
			left:  1,
			right: 1,
			want:  true,
		},
		{
			name:  "basic_values_cmp_notEqual",
			left:  "not_something",
			right: "something",
			want:  false,
		},
		{
			name:  "slice_cmp_always_equal",
			left:  []int{1, 2, 3},
			right: []int{1, 2},
			want:  true,
		},
		{
			name:  "map_cmp_always_equal",
			left:  map[int]int{1: 2, 2: 3},
			right: map[int]int{5: 6, 7: 8},
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cmpComparator := GeneralCompareOptions()
			got := cmp.Equal(tt.left, tt.right, cmpComparator)

			if got != tt.want {
				t.Errorf("GeneralCompareOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
