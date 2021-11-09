package utils

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/disiqueira/gotree"
	"inspr.dev/inspr/pkg/meta"
	"inspr.dev/inspr/pkg/utils"
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
		"\n    │   ├── Input" +
		"\n    │   │   ├── input1" +
		"\n    │   │   ├── input2" +
		"\n    │   ├── Output" +
		"\n    │       └── output1" +
		"\n    │       └── output2" +
		"\n    └── Auth" +
		"\n        └── Scope: " +
		"\n\n"

	channelTree = "channel_name" +
		"\n└── Meta" +
		"\n│   ├── Name: channel_name" +
		"\n└── Spec" +
		"\n│   ├── Type: ct_meta" +
		"\n│   ├── SelectedBroker: " +
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
						Routes:   map[string]*meta.RouteConnection{},
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

func Test_addAppsTree(t *testing.T) {
	type args struct {
		spec gotree.Tree
		app  *meta.App
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Add app test",
			args: args{
				spec: gotree.New("myapp"),
				app: &meta.App{
					Meta: meta.Metadata{
						Name: "myapp",
					},
					Spec: meta.AppSpec{
						Apps: map[string]*meta.App{
							"app1": {},
						},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
						Aliases: map[string]*meta.Alias{},
						Routes:  map[string]*meta.RouteConnection{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.args.spec.Items())
			addAppsTree(tt.args.spec, tt.args.app)

			apps, err := findByName(tt.args.spec.Items(), "Apps")
			if err != nil {
				t.Errorf("AppTree.findByName() error = %v", err)
			}
			app1, err := findByName(apps.Items(), "app1")
			if err != nil {
				t.Errorf("AppTree.findByName() error = %v", err)
			}
			fmt.Println(app1.Text())
		})
	}
}

func Test_addChannelsTree(t *testing.T) {
	type args struct {
		spec gotree.Tree
		app  *meta.App
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Add channel test",
			args: args{
				spec: gotree.New("myapp"),
				app: &meta.App{
					Meta: meta.Metadata{
						Name: "myapp",
					},
					Spec: meta.AppSpec{
						Apps: map[string]*meta.App{},
						Channels: map[string]*meta.Channel{
							"channel1": {},
						},
						Types: map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
						Aliases: map[string]*meta.Alias{},
						Routes:  map[string]*meta.RouteConnection{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.args.spec.Items())
			addChannelsTree(tt.args.spec, tt.args.app)

			Channels, err := findByName(tt.args.spec.Items(), "Channels")
			if err != nil {
				t.Errorf("AppTree.findByName() error = %v", err)
			}
			channel1, err := findByName(Channels.Items(), "channel1")
			if err != nil {
				t.Errorf("AppTree.findByName() error = %v", err)
			}
			fmt.Println(channel1.Text())
		})
	}
}

func Test_addTypesTree(t *testing.T) {
	type args struct {
		spec gotree.Tree
		app  *meta.App
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Add type test",
			args: args{
				spec: gotree.New("myapp"),
				app: &meta.App{
					Meta: meta.Metadata{
						Name: "myapp",
					},
					Spec: meta.AppSpec{
						Apps:     map[string]*meta.App{},
						Channels: map[string]*meta.Channel{},
						Types: map[string]*meta.Type{
							"type1": {},
						},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
						Aliases: map[string]*meta.Alias{},
						Routes:  map[string]*meta.RouteConnection{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.args.spec.Items())
			addTypesTree(tt.args.spec, tt.args.app)

			Types, err := findByName(tt.args.spec.Items(), "Types")
			if err != nil {
				t.Errorf("AppTree.findByName() error = %v", err)
			}
			type1, err := findByName(Types.Items(), "type1")
			if err != nil {
				t.Errorf("AppTree.findByName() error = %v", err)
			}
			fmt.Println(type1.Text())
		})
	}
}

func Test_addAliasesTree(t *testing.T) {
	type args struct {
		spec gotree.Tree
		app  *meta.App
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Add aliases test",
			args: args{
				spec: gotree.New("myapp"),
				app: &meta.App{
					Meta: meta.Metadata{
						Name: "myapp",
					},
					Spec: meta.AppSpec{
						Apps:     map[string]*meta.App{},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
						Aliases: map[string]*meta.Alias{
							"myalias": {
								Target: "myawesometarget",
							},
						},
						Routes: map[string]*meta.RouteConnection{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.args.spec.Items())
			addAliasesTree(tt.args.spec, tt.args.app)

			aliases, err := findByName(tt.args.spec.Items(), "Aliases")
			if err != nil {
				t.Errorf("AliasesTree.findByName() error = %v", err)
			}

			myalias, err := findByName(aliases.Items(), "myalias")
			if err != nil {
				t.Errorf("AliasesTree.findByName() error = %v", err)
			}

			mytarget, err := findByName(myalias.Items(), "Target: myawesometarget")
			if err != nil {
				t.Errorf("AliasesTree.findByName() error = %v", err)
			}

			fmt.Println(mytarget.Text())

		})
	}
}

func Test_addRoutesTree(t *testing.T) {
	type args struct {
		spec gotree.Tree
		app  *meta.App
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Add route test",
			args: args{
				spec: gotree.New("myapp"),
				app: &meta.App{
					Meta: meta.Metadata{
						Name: "myapp",
					},
					Spec: meta.AppSpec{
						Apps:     map[string]*meta.App{},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
						Aliases: map[string]*meta.Alias{},
						Routes: map[string]*meta.RouteConnection{
							"myroute": {
								Endpoints: []string{"endpoint1"},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.args.spec.Items())
			addRoutesTree(tt.args.spec, tt.args.app)

			Routes, err := findByName(tt.args.spec.Items(), "Routes")
			if err != nil {
				t.Errorf("AliasesTree.findByName() error = %v", err)
			}

			myroute, err := findByName(Routes.Items(), "myroute")
			if err != nil {
				t.Errorf("AliasesTree.findByName() error = %v", err)
			}

			endpoint, err := findByName(myroute.Items(), "endpoint1")
			if err != nil {
				t.Errorf("AliasesTree.findByName() error = %v", err)
			}

			fmt.Println(endpoint.Text())

		})
	}
}

func Test_addPermissionsTree(t *testing.T) {
	type args struct {
		auth gotree.Tree
		app  *meta.App
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Add permition test",
			args: args{
				auth: gotree.New("myapp"),
				app: &meta.App{
					Meta: meta.Metadata{
						Name: "myapp",
					},
					Spec: meta.AppSpec{
						Apps:     map[string]*meta.App{},
						Channels: map[string]*meta.Channel{},
						Types:    map[string]*meta.Type{},
						Boundary: meta.AppBoundary{
							Input:  []string{},
							Output: []string{},
						},
						Aliases: map[string]*meta.Alias{},
						Routes:  map[string]*meta.RouteConnection{},
						Auth: meta.AppAuth{
							Permissions: []string{"permition1"},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.args.auth.Items())
			addPermissionsTree(tt.args.auth, tt.args.app)

			Permissions, err := findByName(tt.args.auth.Items(), "Permissions")
			if err != nil {
				t.Errorf("AppTree.findByName() error = %v", err)
			}
			permition1, err := findByName(Permissions.Items(), "permition1")
			if err != nil {
				t.Errorf("AppTree.findByName() error = %v", err)
			}
			fmt.Println(permition1.Text())
		})
	}
}

func findByName(treeArr []gotree.Tree, name string) (gotree.Tree, error) {
	for _, item := range treeArr {
		if item.Text() == name {
			return item, nil
		}
	}
	return nil, errors.New("cannot find " + name)
}
