package utils

import (
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

func TestArrMakeSet(t *testing.T) {
	type args struct {
		strings []string
	}
	tests := []struct {
		name string
		args args
		want StrSet
	}{
		{
			name: "Creates a set from a string slice",
			args: args{
				strings: []string{"its", "a", "set"},
			},
			want: StrSet{"its": true, "a": true, "set": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := MakeStrSet(tt.args.strings); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrMakeSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrSet_ArrAppendSet(t *testing.T) {
	type args struct {
		strings []string
	}
	tests := []struct {
		name string
		set  *StrSet
		args args
	}{
		{
			name: "Appends string slice to a set",
			args: args{
				strings: []string{"set", "indeed"},
			},
			set: &StrSet{"its": true, "a": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			argSet, _ := MakeStrSet(tt.args.strings)
			tt.set.AppendSet(argSet)
			aux := *tt.set
			if !aux["set"] || !aux["indeed"] || !aux["its"] || !aux["a"] {
				t.Errorf("ArrAppendSet() got %v", aux)
			}
		})
	}
}

func TestArrDisjuncSet(t *testing.T) {
	type args struct {
		arr1 []string
		arr2 []string
	}
	tests := []struct {
		name string
		args args
		want StrSet
	}{
		{
			name: "Disjuction between arrays that returns a set",
			args: args{
				arr1: []string{"its"},
				arr2: []string{"set", "its"},
			},
			want: StrSet{"set": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set1, _ := MakeStrSet(tt.args.arr1)
			set2, _ := MakeStrSet(tt.args.arr2)
			if got := DisjunctSet(set1, set2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArrDisjuncSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppMakeSet(t *testing.T) {
	type args struct {
		apps MApps
	}
	tests := []struct {
		name string
		args args
		want StrSet
	}{
		{
			name: "Disjuction between arrays that returns a set",
			args: args{
				apps: MApps{"app1": &meta.App{}, "app2": &meta.App{}},
			},
			want: StrSet{"app1": true, "app2": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := MakeStrSet(tt.args.apps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppMakeSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrSet_AppAppendSet(t *testing.T) {
	type args struct {
		apps MApps
	}
	tests := []struct {
		name string
		set  *StrSet
		args args
	}{
		{
			name: "Appends string slice to a set",
			args: args{
				apps: MApps{"app1": &meta.App{}, "app2": &meta.App{}},
			},
			set: &StrSet{"app1": true, "app2": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			argSet, _ := MakeStrSet(tt.args.apps)
			tt.set.AppendSet(argSet)
			aux := *tt.set
			if !aux["app1"] || !aux["app2"] {
				t.Errorf("AppAppendSet() got %v", aux)
			}
		})
	}
}

func TestAppDisjuncSet(t *testing.T) {
	type args struct {
		apps1 MApps
		apps2 MApps
	}
	tests := []struct {
		name string
		args args
		want StrSet
	}{
		{
			name: "Disjuction between arrays that returns a set",
			args: args{
				apps1: MApps{"app1": &meta.App{}},
				apps2: MApps{"app1": &meta.App{}, "app2": &meta.App{}},
			},
			want: StrSet{"app2": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set1, _ := MakeStrSet(tt.args.apps1)
			set2, _ := MakeStrSet(tt.args.apps2)
			if got := DisjunctSet(set1, set2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppDisjuncSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppIntersecSet(t *testing.T) {
	type args struct {
		apps1 MApps
		apps2 MApps
	}
	tests := []struct {
		name string
		args args
		want StrSet
	}{
		{
			name: "Disjuction between arrays that returns a set",
			args: args{
				apps1: MApps{"app1": &meta.App{}},
				apps2: MApps{"app1": &meta.App{}, "app2": &meta.App{}},
			},
			want: StrSet{"app1": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set1, _ := MakeStrSet(tt.args.apps1)
			set2, _ := MakeStrSet(tt.args.apps2)
			if got := IntersectSet(set1, set2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppIntersecSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChsMakeSet(t *testing.T) {
	type args struct {
		chans MChannels
	}
	tests := []struct {
		name string
		args args
		want StrSet
	}{
		{
			name: "Disjuction between chan arrays that returns a set",
			args: args{
				chans: MChannels{"ch1": &meta.Channel{}, "ch2": &meta.Channel{}},
			},
			want: StrSet{"ch1": true, "ch2": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := MakeStrSet(tt.args.chans); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChsMakeSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrSet_ChsAppendSet(t *testing.T) {
	type args struct {
		apps MChannels
	}
	tests := []struct {
		name string
		set  *StrSet
		args args
	}{
		{
			name: "Appends channel slice to a set",
			args: args{
				apps: MChannels{"ch1": &meta.Channel{}, "ch2": &meta.Channel{}},
			},
			set: &StrSet{"ch1": true, "ch2": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			argSet, _ := MakeStrSet(tt.args.apps)
			tt.set.AppendSet(argSet)
			aux := *tt.set
			if !aux["ch1"] || !aux["ch2"] {
				t.Errorf("ChsAppendSet() got %v", aux)
			}
		})
	}
}

func TestChsDisjuncSet(t *testing.T) {
	type args struct {
		apps1 MChannels
		apps2 MChannels
	}
	tests := []struct {
		name string
		args args
		want StrSet
	}{
		{
			name: "Disjuction between arrays that returns a set",
			args: args{
				apps1: MChannels{"ch1": &meta.Channel{}},
				apps2: MChannels{"ch1": &meta.Channel{}, "ch2": &meta.Channel{}},
			},
			want: StrSet{"ch2": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set1, _ := MakeStrSet(tt.args.apps1)
			set2, _ := MakeStrSet(tt.args.apps2)
			if got := DisjunctSet(set1, set2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChsDisjuncSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChsIntersecSet(t *testing.T) {
	type args struct {
		apps1 MChannels
		apps2 MChannels
	}
	tests := []struct {
		name string
		args args
		want StrSet
	}{
		{
			name: "Disjuction between arrays that returns a set",
			args: args{
				apps1: MChannels{"ch1": &meta.Channel{}},
				apps2: MChannels{"ch1": &meta.Channel{}, "ch2": &meta.Channel{}},
			},
			want: StrSet{"ch1": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set1, _ := MakeStrSet(tt.args.apps1)
			set2, _ := MakeStrSet(tt.args.apps2)
			if got := IntersectSet(set1, set2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChsIntersecSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypesMakeSet(t *testing.T) {
	type args struct {
		types MTypes
	}
	tests := []struct {
		name string
		args args
		want StrSet
	}{
		{
			name: "Disjuction between ctypes arrays that returns a set",
			args: args{
				types: MTypes{"ct1": &meta.ChannelType{}, "ct2": &meta.ChannelType{}},
			},
			want: StrSet{"ct1": true, "ct2": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := MakeStrSet(tt.args.types); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TypesMakeSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrSet_TypesAppendSet(t *testing.T) {
	type args struct {
		types MTypes
	}
	tests := []struct {
		name string
		set  *StrSet
		args args
	}{
		{
			name: "Appends ctypes slice to a set",
			args: args{
				types: MTypes{"ct1": &meta.ChannelType{}, "ct2": &meta.ChannelType{}},
			},
			set: &StrSet{"ct1": true, "ct2": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			argSet, _ := MakeStrSet(tt.args.types)
			tt.set.AppendSet(argSet)
			aux := *tt.set
			if !aux["ct1"] || !aux["ct2"] {
				t.Errorf("TypesAppendSet() got %v", aux)
			}
		})
	}
}

func TestTypesDisjuncSet(t *testing.T) {
	type args struct {
		types1 MTypes
		types2 MTypes
	}
	tests := []struct {
		name string
		args args
		want StrSet
	}{
		{
			name: "Disjuction between ct arrays that returns a set",
			args: args{
				types1: MTypes{"ct1": &meta.ChannelType{}},
				types2: MTypes{"ct1": &meta.ChannelType{}, "ct2": &meta.ChannelType{}},
			},
			want: StrSet{"ct2": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set1, _ := MakeStrSet(tt.args.types1)
			set2, _ := MakeStrSet(tt.args.types2)
			if got := DisjunctSet(set1, set2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TypesDisjuncSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypesIntersecSet(t *testing.T) {
	type args struct {
		types1 MTypes
		types2 MTypes
	}
	tests := []struct {
		name string
		args args
		want StrSet
	}{
		{
			name: "Intersection between ct arrays that returns a set",
			args: args{
				types1: MTypes{"ct1": &meta.ChannelType{}},
				types2: MTypes{"ct1": &meta.ChannelType{}, "ct2": &meta.ChannelType{}},
			},
			want: StrSet{"ct1": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set1, _ := MakeStrSet(tt.args.types1)
			set2, _ := MakeStrSet(tt.args.types2)
			if got := IntersectSet(set1, set2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TypesIntersecSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrMakeSet(t *testing.T) {
	type args struct {
		strings MStr
	}
	tests := []struct {
		name string
		args args
		want StrSet
	}{
		{
			name: "Creates a set from a string slice",
			args: args{
				strings: MStr{"str1": "one", "str2": "two"},
			},
			want: StrSet{"str1": true, "str2": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := MakeStrSet(tt.args.strings); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrMakeSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrSet_StrAppendSet(t *testing.T) {
	type args struct {
		strings MStr
	}
	tests := []struct {
		name string
		set  *StrSet
		args args
	}{
		{
			name: "Creates a set from a string slice",
			args: args{
				strings: MStr{"str1": "one", "str2": "two"},
			},
			set: &StrSet{"str1": true, "str2": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			argSet, _ := MakeStrSet(tt.args.strings)
			tt.set.AppendSet(argSet)
			aux := *tt.set
			if !aux["str1"] || !aux["str2"] {
				t.Errorf("StrAppendSet() got %v", aux)
			}
		})
	}
}

func TestStrDisjuncSet(t *testing.T) {
	type args struct {
		strings1 MStr
		strings2 MStr
	}
	tests := []struct {
		name string
		args args
		want StrSet
	}{
		{
			name: "Creates a set from a string slice",
			args: args{
				strings1: MStr{"str1": "one"},
				strings2: MStr{"str1": "one", "str2": "two"},
			},
			want: StrSet{"str2": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set1, _ := MakeStrSet(tt.args.strings1)
			set2, _ := MakeStrSet(tt.args.strings2)
			if got := DisjunctSet(set1, set2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrDisjuncSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeStrSet(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    StrSet
		wantErr bool
	}{
		{
			name: "Type not supported - it should return an error",
			args: args{
				obj: 3,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeStrSet(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeStrSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeStrSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
