package utils

import (
	"reflect"
	"strings"
	"testing"

	test "gitlab.inspr.dev/inspr/core/pkg/testutils"
	"gotest.tools/assert/cmp"
	kubeCore "k8s.io/api/core/v1"
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
		name          string
		args          args
		want          []string
		checkFunction func(t *testing.T, got []string)
	}{
		{
			name: "It should return the union of two slices (without repeated elements)",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"a", "b", "d"},
			},
			want: []string{"a", "b", "c", "d"},
			checkFunction: func(t *testing.T, got []string) {
				if !(Includes(got, "a") && Includes(got, "b") && Includes(got, "c") && Includes(got, "d") && len(got) == 4) {
					t.Errorf("StringSliceUnion() = %v, want %v", got, "[a b c d]")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.checkFunction(t, StringSliceUnion(tt.args.a, tt.args.b))
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

func TestStringArray_Map(t *testing.T) {
	type args struct {
		f func(string) string
	}
	tests := []struct {
		name string
		c    StringArray
		args args
		want StringArray
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Map(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringArray.Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringArray_Union(t *testing.T) {
	type args struct {
		other StringArray
	}
	tests := []struct {
		name string
		c    StringArray
		args args
		want StringArray
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Union(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringArray.Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringArray_Contains(t *testing.T) {
	type args struct {
		item string
	}
	tests := []struct {
		name string
		c    StringArray
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Contains(tt.args.item); got != tt.want {
				t.Errorf("StringArray.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringArray_Join(t *testing.T) {
	type args struct {
		sep string
	}
	tests := []struct {
		name string
		c    StringArray
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Join(tt.args.sep); got != tt.want {
				t.Errorf("StringArray.Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseToK8sArrEnv(t *testing.T) {
	type args struct {
		arrappEnv EnvironmentMap
	}
	tests := []struct {
		name string
		args args
		want []kubeCore.EnvVar
	}{
		{
			name: "successful_test",
			args: args{
				arrappEnv: map[string]string{
					"key_1": "value_1",
					"key_2": "value_2",
					"key_3": "value_3",
				},
			},
			want: []kubeCore.EnvVar{
				{
					Name:  "key_1",
					Value: "value_1",
				},
				{
					Name:  "key_2",
					Value: "value_2",
				},
				{
					Name:  "key_3",
					Value: "value_3",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.arrappEnv.ParseToK8sArrEnv()

			var comp cmp.Comparison = cmp.DeepEqual(got, tt.want, test.GetMapCompareOptions())
			if !comp().Success() {
				t.Errorf("parseToK8sArrEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
