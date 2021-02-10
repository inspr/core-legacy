package utils

import (
	"reflect"
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	type args struct {
		vs []string
		t  string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "It should return the right position of 'a'",
			args: args{
				vs: []string{"a", "b", "c"},
				t:  "a",
			},
			want: 0,
		},
		{
			name: "It should return the right position of 'b'",
			args: args{
				vs: []string{"a", "b", "c"},
				t:  "b",
			},
			want: 1,
		},
		{
			name: "It should return the right position of 'c'",
			args: args{
				vs: []string{"a", "b", "c"},
				t:  "c",
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Index(tt.args.vs, tt.args.t); got != tt.want {
				t.Errorf("Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIncludes(t *testing.T) {
	type args struct {
		vs []string
		t  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "It should return true, because t exists in vs",
			args: args{
				vs: []string{"a", "b", "c"},
				t:  "c",
			},
			want: true,
		},
		{
			name: "It should return false, because t doens't exist in vs",
			args: args{
				vs: []string{"a", "b", "c"},
				t:  "d",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Includes(tt.args.vs, tt.args.t); got != tt.want {
				t.Errorf("Includes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	type args struct {
		vs []string
		t  string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "It should return vs but without t",
			args: args{
				vs: []string{"a", "b", "c"},
				t:  "c",
			},
			want: []string{"a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Remove(tt.args.vs, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSliceUnion(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "It should return the union of two slices (without repeated elements)",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"a", "b", "d"},
			},
			want: []string{"a", "b", "c", "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceUnion(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSliceUnion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	type args struct {
		vs []string
		f  func(string) string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "It should apply a function to each element of the slice",
			args: args{
				vs: []string{"a", "b", "c"},
				f:  strings.ToUpper,
			},
			want: []string{"A", "B", "C"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.vs, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
