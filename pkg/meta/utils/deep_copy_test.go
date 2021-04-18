package utils

import (
	"reflect"
	"testing"

	"inspr.dev/inspr/pkg/meta"
)

func getMockedTree() *meta.App {
	root := meta.App{
		Meta: meta.Metadata{
			Name:        "",
			Reference:   "",
			Annotations: map[string]string{},
			Parent:      "",
			UUID:        "",
		},
		Spec: meta.AppSpec{
			Node:     meta.Node{},
			Apps:     map[string]*meta.App{},
			Channels: map[string]*meta.Channel{},
			ChannelTypes: map[string]*meta.ChannelType{
				"ct1": {
					Meta: meta.Metadata{
						Name:        "ct1",
						Reference:   "ct1",
						Annotations: map[string]string{},
						Parent:      "",
						UUID:        "",
					},
					Schema: "",
				},
			},
			Boundary: meta.AppBoundary{
				Input:  []string{},
				Output: []string{},
			},
		},
	}
	return &root
}

func TestDeepCopy(t *testing.T) {

	var metaApp *meta.App
	var stringArr []string
	var integer int

	type args struct {
		orig interface{}
		dest interface{}
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		checkFunction func(t *testing.T)
	}{
		{
			name: "Deep Copy: meta.App",
			args: args{
				orig: getMockedTree(),
				dest: &metaApp,
			},
			wantErr: false,
			checkFunction: func(t *testing.T) {
				if !reflect.DeepEqual(getMockedTree(), metaApp) {
					t.Errorf("DeepCopy() got = %v, wantErr %v", metaApp, getMockedTree())
				}
			},
		},
		{
			name: "Deep Copy: string array",
			args: args{
				orig: []string{"A", "B", "C"},
				dest: &stringArr,
			},
			wantErr: false,
			checkFunction: func(t *testing.T) {
				if !reflect.DeepEqual([]string{"A", "B", "C"}, stringArr) {
					t.Errorf("DeepCopy() got = %v, wantErr %v", stringArr, []string{"A", "B", "C"})
				}
			},
		},
		{
			name: "Deep Copy: integer",
			args: args{
				orig: 11,
				dest: &integer,
			},
			wantErr: false,
			checkFunction: func(t *testing.T) {
				if !reflect.DeepEqual(11, integer) {
					t.Errorf("DeepCopy() got = %v, wantErr %v", integer, 11)
				}
			},
		},
		{
			name: "Types dont match - it should return an error",
			args: args{
				orig: 11,
				dest: &metaApp,
			},
			wantErr: true,
		},
		{
			name: "dest its not a pointer - it should return an error",
			args: args{
				orig: 11,
				dest: integer,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeepCopy(tt.args.orig, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("DeepCopy() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.checkFunction != nil {
				tt.checkFunction(t)
			}
		})
	}
}
