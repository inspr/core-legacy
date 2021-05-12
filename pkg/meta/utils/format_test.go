package utils

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/disiqueira/gotree"
	"github.com/inspr/inspr/pkg/meta"
	"github.com/inspr/inspr/pkg/utils"
)

var (
	wantAppTree = "mock_app" +
		"\n└── Meta" +
		"\n│   ├── Name: mock_app" +
		"\n│   ├── Parent: mock_parent" +
		"\n│   ├── Reference: mock_ref" +
		"\n│   ├── Annotations" +
		"\n│       └── mock_key: mock_val" +
		"\n└── Spec" +
		"\n    └── Boundary" +
		"\n        └── Input" +
		"\n        │   ├── input1" +
		"\n        │   ├── input2" +
		"\n        └── Output" +
		"\n            └── output1" +
		"\n            └── output2" +
		"\n\n"

	channelTree = "channel_name" +
		"\n└── Meta" +
		"\n│   ├── Name: channel_name" +
		"\n└── Spec" +
		"\n│   ├── Type: ct_meta" +
		"\n└── ConnectedApps" +
		"\n    └── a" +
		"\n    └── b" +
		"\n    └── c" +
		"\n\n"

	TypeTree = "ct_meta" +
		"\n└── Meta" +
		"\n│   ├── Name: ct_meta" +
		"\n└── Spec" +
		"\n    └── Schema: {\"type\":\"int\"}" +
		"\n\n"

	aliasTree = "alias_name" +
		"\n└── Meta" +
		"\n│   ├── Name: alias_name" +
		"\n└── Target: alias_target" +
		"\n\n"
)

func TestPrintAppTree(t *testing.T) {

	type args struct {
		app *meta.App
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			name: "basic_app_tree",
			args: args{
				app: &meta.App{
					Meta: meta.Metadata{
						Name:      "mock_app",
						Reference: "mock_ref",
						Annotations: map[string]string{
							"mock_key": "mock_val",
						},
						Parent: "mock_parent",
					},
					Spec: meta.AppSpec{
						Node:     meta.Node{},
						Apps:     map[string]*meta.App{},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  utils.StringArray{"input1", "input2"},
							Output: utils.StringArray{"output1", "output2"},
						},
					},
				},
			},
			wantOut: wantAppTree,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			PrintAppTree(tt.args.app, out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf(
					"PrintAppTree() = \n%v, want \n%v",
					gotOut,
					tt.wantOut,
				)
			}
		})
	}
}

func TestPrintChannelTree(t *testing.T) {
	type args struct {
		ch *meta.Channel
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			name: "basic_channel_tree",
			args: args{
				ch: &meta.Channel{
					Meta: meta.Metadata{
						Name: "channel_name",
					},
					Spec: meta.ChannelSpec{
						Type: "ct_meta",
					},
					ConnectedApps: []string{"a", "b", "c"},
				},
			},
			wantOut: channelTree,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			PrintChannelTree(tt.args.ch, out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf(
					"PrintChannelTree() = \n%v, want \n%v",
					gotOut,
					tt.wantOut,
				)
			}
		})
	}
}

func TestPrintTypeTree(t *testing.T) {
	type args struct {
		ct *meta.Type
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			name: "basic_Type_tree",
			args: args{
				ct: &meta.Type{
					Meta: meta.Metadata{
						Name: "ct_meta",
					},
					Schema: "{\"type\":\"int\"}",
				},
			},
			wantOut: TypeTree,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			PrintTypeTree(tt.args.ct, out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf(
					"PrintTypeTree() = \n%v, want \n%v",
					gotOut,
					tt.wantOut,
				)
			}
		})
	}
}

func Test_populateMeta(t *testing.T) {
	// vars
	var treeArg = gotree.New("mock_tree")
	var metaArg *meta.Metadata
	var wantTree = gotree.New("mock_tree")

	// preparation for test
	tName := "populate_test"
	metaArg = &meta.Metadata{
		Name:      "mock_name",
		Reference: "mock_reference",
		Annotations: map[string]string{
			"mock_key": "mock_value",
		},
		Parent: "mock_parent",
		UUID:   "mock_SHA256",
	}
	wantTree.Add("Name: " + metaArg.Name)
	wantTree.Add("Parent: " + metaArg.Parent)
	wantTree.Add("Reference: " + metaArg.Reference)
	wantTree.Add("UUID: " + metaArg.UUID)
	annotations := wantTree.Add("Annotations")
	for noteName, note := range metaArg.Annotations {
		annotations.Add(noteName + ": " + note)
	}

	// tests
	t.Run(tName, func(t *testing.T) {
		populateMeta(treeArg, metaArg)
		if !reflect.DeepEqual(treeArg, wantTree) {
			t.Errorf(
				"populateMeta() = %v, want %v",
				treeArg,
				wantTree,
			)
		}
	})
}

func TestPrintAliasTree(t *testing.T) {
	type args struct {
		al *meta.Alias
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			name: "basic_alias_tree",
			args: args{
				al: &meta.Alias{
					Meta: meta.Metadata{
						Name: "alias_name",
					},
					Target: "alias_target",
				},
			},
			wantOut: aliasTree,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			PrintAliasTree(tt.args.al, out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf(
					"PrintAliasTree() = \n%v, want \n%v",
					gotOut,
					tt.wantOut,
				)
			}
		})
	}
}
