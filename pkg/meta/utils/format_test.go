package utils

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/disiqueira/gotree"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/utils"
)

var (
	wantAppTree = "mock_app" +
		"\n└── Meta" +
		"\n│   ├── Name: mock_app" +
		"\n│   ├── Parent: mock_parent" +
		"\n│   ├── Reference: mock_ref" +
		"\n│   ├── SHA256: " +
		"\n│   ├── Annotations" +
		"\n│       └── mock_key: mock_val" +
		"\n└── Spec" +
		"\n    └── Apps" +
		"\n    └── Channels" +
		"\n    └── ChannelTypes" +
		"\n    └── Node: " +
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
		"\n│   ├── Parent: " +
		"\n│   ├── Reference: " +
		"\n│   ├── SHA256: " +
		"\n│   ├── Annotations" +
		"\n└── Spec" +
		"\n│   ├── Type: ct_meta" +
		"\n└── ConnectedApps" +
		"\n    └── a" +
		"\n    └── b" +
		"\n    └── c" +
		"\n\n"

	channelTypeTree = "ct_meta" +
		"\n└── Meta" +
		"\n│   ├── Name: ct_meta" +
		"\n│   ├── Parent: " +
		"\n│   ├── Reference: " +
		"\n│   ├── SHA256: " +
		"\n│   ├── Annotations" +
		"\n└── Spec" +
		"\n│   ├── Schema: {\"type\":\"int\"}" +
		"\n└── ConnectedChannels" +
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
						Node:         meta.Node{},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
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

func TestPrintChannelTypeTree(t *testing.T) {
	type args struct {
		ct *meta.ChannelType
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			name: "basic_channelType_tree",
			args: args{
				ct: &meta.ChannelType{
					Meta: meta.Metadata{
						Name: "ct_meta",
					},
					Schema: "{\"type\":\"int\"}",
				},
			},
			wantOut: channelTypeTree,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			PrintChannelTypeTree(tt.args.ct, out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf(
					"PrintChannelTypeTree() = \n%v, want \n%v",
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
		SHA256: "mock_SHA256",
	}
	wantTree.Add("Name: " + metaArg.Name)
	wantTree.Add("Parent: " + metaArg.Parent)
	wantTree.Add("Reference: " + metaArg.Reference)
	wantTree.Add("SHA256: " + metaArg.SHA256)
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
