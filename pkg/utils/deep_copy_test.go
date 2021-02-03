package utils

import (
	"encoding/json"
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

func TestDCopy(t *testing.T) {
	type args struct {
		root *meta.App
	}
	tests := []struct {
		name    string
		args    args
		want    *meta.App
		wantErr bool
	}{
		{
			name: "Deep copy test",
			args: args{
				root: getMockedTree(),
			},
			want:    getMockedTree(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DCopy(tt.args.root)
			jsgot, _ := json.MarshalIndent(*got, "", "	")
			jswant, _ := json.MarshalIndent(*tt.want, "", "	")
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DCopy() = %v, want %v", string(jsgot), string(jswant))
			}
		})
	}
}

func getMockedTree() *meta.App {
	root := meta.App{
		Meta: meta.Metadata{
			Name:        "",
			Reference:   "",
			Annotations: map[string]string{},
			Parent:      "",
			SHA256:      "",
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
						SHA256:      "",
					},
					Schema: []byte{},
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
