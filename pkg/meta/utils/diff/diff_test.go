package diff

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"text/tabwriter"

	"github.com/google/go-cmp/cmp"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/meta/utils"

	pkgUtils "gitlab.inspr.dev/inspr/core/pkg/utils"
)

func TestDiff(t *testing.T) {
	type args struct {
		appOrig *meta.App
		appCurr *meta.App
	}
	tests := []struct {
		name    string
		args    args
		want    Changelog
		wantErr bool
	}{
		{
			name: "Cascading changes test",
			args: args{
				appOrig: getMockRootApp(),
				appCurr: getMockRootApp2(),
			},
			want: Changelog{
				{
					Context:   "app2.app3",
					Kind:      NodeKind,
					Operation: Update,
					Diff: []Difference{
						{
							Field:     "Spec.Node.Spec.Image",
							From:      "imageNodeApp3",
							To:        "imageNodeApp3diff",
							Kind:      NodeKind,
							Operation: Update,
						},
					},
				},
				{
					Context:   "",
					Kind:      MetaKind | AnnotationKind | AppKind | ChannelKind | ChannelTypeKind,
					Operation: Create | Update | Delete,
					Diff: []Difference{
						{
							Field:     "Meta.SHA256",
							From:      "1",
							To:        "2",
							Kind:      MetaKind | AppKind,
							Operation: Update,
						},
						{
							Field:     "Meta.Annotations[an1]",
							From:      "<nil>",
							To:        "a",
							Kind:      AnnotationKind | AppKind | MetaKind,
							Name:      "an1",
							Operation: Create,
						},
						{
							Field:     "Meta.Annotations[an2]",
							From:      "<nil>",
							To:        "b",
							Kind:      MetaKind | AppKind | AnnotationKind,
							Operation: Create,
							Name:      "an2",
						},
						{
							Field:     "Spec.Apps[app1]",
							From:      "{...}",
							To:        "<nil>",
							Kind:      AppKind,
							Operation: Delete,
							Name:      "app1",
						},
						{
							Field:     "Spec.Channels[ch2]",
							From:      "{...}",
							To:        "<nil>",
							Kind:      ChannelKind,
							Operation: Delete,
							Name:      "ch2",
						},
						{
							Field:     "Spec.Channels[ch1].Meta.Reference",
							From:      "root.ch1",
							To:        "root.ch1diff",
							Kind:      MetaKind | ChannelKind,
							Operation: Update,
							Name:      "ch1",
						},
						{
							Field:     "Spec.ChannelTypes[ct2]",
							From:      "{...}",
							To:        "<nil>",
							Kind:      ChannelTypeKind,
							Operation: Delete,
							Name:      "ct2",
						},
						{
							Field:     "Spec.ChannelTypes[ct1].Meta.Reference",
							From:      "root.ct1",
							To:        "root.ct1diff",
							Kind:      MetaKind | ChannelTypeKind,
							Operation: Update,
							Name:      "ct1",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Diff(tt.args.appOrig, tt.args.appCurr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Diff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !equalChangelogs(got, tt.want) {
				gots, _ := json.MarshalIndent(got, "", "    ")
				wants, _ := json.MarshalIndent(tt.want, "", "    ")
				t.Errorf("Diff() = \n%s, want \n%s", gots, wants)
			}
		})
	}
}

func TestChange_diffAppSpec(t *testing.T) {
	type fields struct {
		Context string
		Diff    []Difference
	}
	type args struct {
		specOrig meta.AppSpec
		specCurr meta.AppSpec
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    Change
	}{
		{
			name:   "Uchanged Specs",
			fields: fields{},
			args: args{
				specOrig: meta.AppSpec{},
				specCurr: meta.AppSpec{},
			},
			want: Change{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			change := &Change{
				Context:   tt.fields.Context,
				Diff:      tt.fields.Diff,
				changelog: &Changelog{},
			}

			if err := change.diffAppSpec(tt.args.specOrig, tt.args.specCurr); (err != nil) != tt.wantErr {
				t.Errorf("Change.diffAppSpec() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !equalChanges(*change, tt.want) {
				t.Errorf("Changelog.diff() = %v, want %v", *change, tt.want)
			}
		})
	}
}

func TestChange_diffNodes(t *testing.T) {
	type fields struct {
		Context string
		Diff    []Difference
	}
	type args struct {
		nodeOrig meta.Node
		nodeCurr meta.Node
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    Change
	}{
		{
			name:   "Uchanged nodes",
			fields: fields{},
			args: args{
				nodeOrig: meta.Node{
					Meta: meta.Metadata{},
					Spec: meta.NodeSpec{
						Image: "image",
					},
				},
				nodeCurr: meta.Node{
					Meta: meta.Metadata{},
					Spec: meta.NodeSpec{
						Image: "image",
					},
				},
			},
			wantErr: false,
			want:    Change{},
		},
		{
			name:   "Valid change on nodes",
			fields: fields{},
			args: args{
				nodeOrig: meta.Node{
					Meta: meta.Metadata{},
					Spec: meta.NodeSpec{
						Image: "image",
					},
				},
				nodeCurr: meta.Node{
					Meta: meta.Metadata{},
					Spec: meta.NodeSpec{
						Image: "imagediff",
					},
				},
			},
			wantErr: false,
			want: Change{
				Kind:      NodeKind,
				Operation: Update,
				Diff: []Difference{
					{
						Field:     "Spec.Node.Spec.Image",
						From:      "image",
						To:        "imagediff",
						Kind:      NodeKind,
						Operation: Update,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			change := &Change{
				Context: tt.fields.Context,
				Diff:    tt.fields.Diff,
			}
			if err := change.diffNodes(tt.args.nodeOrig, tt.args.nodeCurr); (err != nil) != tt.wantErr {
				t.Errorf("Change.diffNodes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !equalChanges(*change, tt.want) {
				t.Errorf("Changelog.diff() = %v, want %v", *change, tt.want)
			}
		})
	}
}
func TestChange_diffBoudaries(t *testing.T) {
	type fields struct {
		Context string
		Diff    []Difference
	}
	type args struct {
		boundOrig meta.AppBoundary
		boundCurr meta.AppBoundary
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Change
	}{
		{
			name:   "Unchanged Boundaries",
			fields: fields{},
			args: args{
				boundOrig: meta.AppBoundary{
					Input: []string{
						"a",
						"b",
						"c",
					},
					Output: []string{
						"a",
						"b",
						"c",
					},
				},
				boundCurr: meta.AppBoundary{
					Input: []string{
						"a",
						"b",
						"c",
					},
					Output: []string{
						"a",
						"b",
						"c",
					},
				},
			},
			want: Change{},
		},
		{
			name:   "Valid changes on Boundaries",
			fields: fields{},
			args: args{
				boundOrig: meta.AppBoundary{
					Input: []string{
						"a",
						"b",
						"c",
					},
					Output: []string{
						"a",
						"b",
						"c",
					},
				},
				boundCurr: meta.AppBoundary{
					Input: []string{
						"a",
						"b",
						"d",
					},
					Output: []string{
						"a",
						"b",
						"d",
					},
				},
			},
			want: Change{
				Kind:      BoundaryKind,
				Operation: Create | Delete,
				Diff: []Difference{
					{
						Field:     "Spec.Boundary.Input",
						From:      "c",
						To:        "<nil>",
						Kind:      BoundaryKind,
						Operation: Delete,
						Name:      "c",
					},
					{
						Field:     "Spec.Boundary.Input",
						From:      "<nil>",
						To:        "d",
						Kind:      BoundaryKind,
						Operation: Create,
						Name:      "d",
					},
					{
						Field:     "Spec.Boundary.Output",
						From:      "c",
						To:        "<nil>",
						Kind:      BoundaryKind,
						Operation: Delete,
						Name:      "c",
					},
					{
						Field:     "Spec.Boundary.Output",
						From:      "<nil>",
						To:        "d",
						Kind:      BoundaryKind,
						Operation: Create,
						Name:      "d",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			change := &Change{
				Context: tt.fields.Context,
				Diff:    tt.fields.Diff,
			}
			change.diffBoudaries(tt.args.boundOrig, tt.args.boundCurr)
			if !equalChanges(*change, tt.want) {
				t.Errorf("Changelog.diff() = %v, want %v", *change, tt.want)
			}
		})
	}
}

func TestChange_diffApps(t *testing.T) {
	type fields struct {
		Context string
		Diff    []Difference
	}
	type args struct {
		appsOrig utils.MApps
		appsCurr utils.MApps
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Change
	}{
		{
			name:   "Unchanged Apps",
			fields: fields{},
			args: args{
				appsOrig: utils.MApps{
					"app1": {
						Meta: meta.Metadata{},
						Spec: meta.AppSpec{
							Node: meta.Node{
								Meta: meta.Metadata{},
								Spec: meta.NodeSpec{
									Image: "imageNodeApp1",
								},
							},
						},
					},
				},
				appsCurr: utils.MApps{
					"app1": {
						Meta: meta.Metadata{},
						Spec: meta.AppSpec{
							Node: meta.Node{
								Meta: meta.Metadata{},
								Spec: meta.NodeSpec{
									Image: "imageNodeApp1",
								},
							},
						},
					},
				},
			},
			want: Change{},
		},
		{
			name:   "Valid changes on Apps",
			fields: fields{},
			args: args{
				appsOrig: utils.MApps{
					"app1": {
						Meta: meta.Metadata{},
						Spec: meta.AppSpec{},
					},
					"app3": {
						Meta: meta.Metadata{},
						Spec: meta.AppSpec{},
					},
				},
				appsCurr: utils.MApps{
					"app1": {
						Meta: meta.Metadata{},
						Spec: meta.AppSpec{},
					},
					"app2": {
						Meta: meta.Metadata{},
						Spec: meta.AppSpec{},
					},
				},
			},
			want: Change{
				Kind:      AppKind,
				Operation: Create | Delete,
				Diff: []Difference{
					{
						Field:     "Spec.Apps[app2]",
						From:      "<nil>",
						To:        "{...}",
						Kind:      AppKind,
						Operation: Create,
						Name:      "app2",
					},
					{
						Field:     "Spec.Apps[app3]",
						From:      "{...}",
						To:        "<nil>",
						Kind:      AppKind,
						Operation: Delete,
						Name:      "app3",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			change := &Change{
				Context:   tt.fields.Context,
				Diff:      tt.fields.Diff,
				changelog: &Changelog{},
			}
			change.diffApps(tt.args.appsOrig, tt.args.appsCurr)
			if !equalChanges(*change, tt.want) {
				t.Errorf("Changelog.diff() = %v, want %v", *change, tt.want)
			}
		})
	}
}

func TestChange_diffChannels(t *testing.T) {
	type fields struct {
		Context string
		Diff    []Difference
	}
	type args struct {
		chOrig utils.MChannels
		chCurr utils.MChannels
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    Change
	}{
		{
			name:   "Unchanged Channel Types",
			fields: fields{},
			args: args{
				chOrig: utils.MChannels{
					"ch1": &meta.Channel{
						Meta: meta.Metadata{},
						Spec: meta.ChannelSpec{
							Type: "type",
						},
					},
				},
				chCurr: utils.MChannels{
					"ch1": &meta.Channel{
						Meta: meta.Metadata{},
						Spec: meta.ChannelSpec{
							Type: "type",
						},
					},
				},
			},
			wantErr: false,
			want:    Change{},
		},
		{
			name:   "Valid changes on Channel Types",
			fields: fields{},
			args: args{
				chOrig: utils.MChannels{
					"ch1": &meta.Channel{
						Meta: meta.Metadata{},
						Spec: meta.ChannelSpec{
							Type: "type",
						},
					},
				},
				chCurr: utils.MChannels{
					"ch1": &meta.Channel{
						Meta: meta.Metadata{},
						Spec: meta.ChannelSpec{
							Type: "typediff",
						},
					},
				},
			},
			wantErr: false,
			want: Change{
				Kind:      ChannelKind,
				Operation: Update,
				Diff: []Difference{
					{
						Field:     "Spec.Channels[ch1].Spec.Type",
						From:      "type",
						To:        "typediff",
						Kind:      ChannelKind,
						Operation: Update,
						Name:      "ch1",
					},
				},
			},
		},
		{
			name:   "Channel deleted",
			fields: fields{},
			args: args{
				chOrig: utils.MChannels{
					"ch1": &meta.Channel{
						Meta: meta.Metadata{},
						Spec: meta.ChannelSpec{
							Type: "type",
						},
					},
				},
				chCurr: utils.MChannels{},
			},
			wantErr: false,
			want: Change{
				Kind:      ChannelKind,
				Operation: Delete,
				Diff: []Difference{
					{
						Field:     "Spec.Channels[ch1]",
						From:      "{...}",
						To:        "<nil>",
						Kind:      ChannelKind,
						Operation: Delete,
						Name:      "ch1",
					},
				},
			},
		},
		{
			name:   "Channel created",
			fields: fields{},
			args: args{
				chOrig: utils.MChannels{},
				chCurr: utils.MChannels{
					"ch1": &meta.Channel{
						Meta: meta.Metadata{},
						Spec: meta.ChannelSpec{
							Type: "type",
						},
					},
				},
			},
			wantErr: false,
			want: Change{
				Kind:      ChannelKind,
				Operation: Create,
				Diff: []Difference{
					{
						Field:     "Spec.Channels[ch1]",
						From:      "<nil>",
						To:        "{...}",
						Kind:      ChannelKind,
						Operation: Create,
						Name:      "ch1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			change := &Change{
				Context: tt.fields.Context,
				Diff:    tt.fields.Diff,
			}
			if err := change.diffChannels(tt.args.chOrig, tt.args.chCurr); (err != nil) != tt.wantErr {
				t.Errorf("Change.diffChannels() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !equalChanges(*change, tt.want) {
				t.Errorf("Changelog.diff() = %v, want %v", *change, tt.want)
			}
		})
	}
}

func TestChange_diffChannelTypes(t *testing.T) {
	type fields struct {
		Context string
		Diff    []Difference
	}
	type args struct {
		chtOrig utils.MTypes
		chtCurr utils.MTypes
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    Change
	}{
		{
			name:   "Unchanged Channel Types",
			fields: fields{},
			args: args{
				chtOrig: utils.MTypes{
					"ct1": &meta.ChannelType{
						Meta:   meta.Metadata{},
						Schema: "",
					},
				},
				chtCurr: utils.MTypes{
					"ct1": &meta.ChannelType{
						Meta:   meta.Metadata{},
						Schema: "",
					},
				},
			},
			wantErr: false,
			want:    Change{},
		},
		{
			name:   "Valid changes on Channel Types",
			fields: fields{},
			args: args{
				chtOrig: utils.MTypes{
					"ct1": &meta.ChannelType{
						Meta: meta.Metadata{
							Name:        "ct1",
							Reference:   "root.ct1",
							Annotations: map[string]string{},
							Parent:      "root",
							SHA256:      "",
						},
						Schema: string([]byte{0, 1, 0, 1, 0, 0, 1, 1, 1, 0}),
					},
				},
				chtCurr: utils.MTypes{
					"ct1": &meta.ChannelType{
						Meta: meta.Metadata{
							Name:        "ct1",
							Reference:   "root.ct1",
							Annotations: map[string]string{},
							Parent:      "root",
							SHA256:      "",
						},
						Schema: string([]byte{0, 1, 0, 1, 0, 1, 1, 1, 1, 1}),
					},
				},
			},
			wantErr: false,
			want: Change{
				Kind:      ChannelTypeKind,
				Operation: Update,
				Diff: []Difference{
					{
						Field:     "Spec.ChannelTypes[ct1].Spec.Schema",
						From:      string([]byte{0, 1, 0, 1, 0, 0, 1, 1, 1, 0}),
						To:        string([]byte{0, 1, 0, 1, 0, 1, 1, 1, 1, 1}),
						Kind:      ChannelTypeKind,
						Operation: Update,
						Name:      "ct1",
					},
				},
			},
		},
		{
			name:   "create channel type",
			fields: fields{},
			args: args{
				chtOrig: utils.MTypes{},
				chtCurr: utils.MTypes{
					"ct1": &meta.ChannelType{
						Meta: meta.Metadata{
							Name:        "ct1",
							Reference:   "root.ct1",
							Annotations: map[string]string{},
							Parent:      "root",
							SHA256:      "",
						},
						Schema: string([]byte{0, 1, 0, 1, 0, 1, 1, 1, 1, 1}),
					},
				},
			},
			wantErr: false,
			want: Change{
				Kind:      ChannelTypeKind,
				Operation: Create,
				Diff: []Difference{
					{
						Field:     "Spec.ChannelTypes[ct1]",
						From:      "<nil>",
						To:        "{...}",
						Kind:      ChannelTypeKind,
						Operation: Create,
						Name:      "ct1",
					},
				},
			},
		},
		{
			name:   "delete channel type",
			fields: fields{},
			args: args{
				chtOrig: utils.MTypes{
					"ct1": &meta.ChannelType{
						Meta: meta.Metadata{
							Name:        "ct1",
							Reference:   "root.ct1",
							Annotations: map[string]string{},
							Parent:      "root",
							SHA256:      "",
						},
						Schema: string([]byte{0, 1, 0, 1, 0, 1, 1, 1, 1, 1}),
					},
				},
				chtCurr: utils.MTypes{},
			},
			wantErr: false,
			want: Change{
				Kind:      ChannelTypeKind,
				Operation: Delete,
				Diff: []Difference{
					{
						Field:     "Spec.ChannelTypes[ct1]",
						From:      "{...}",
						To:        "<nil>",
						Kind:      ChannelTypeKind,
						Operation: Delete,
						Name:      "ct1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			change := &Change{
				Context: tt.fields.Context,
				Diff:    tt.fields.Diff,
			}
			if err := change.diffChannelTypes(tt.args.chtOrig, tt.args.chtCurr); (err != nil) != tt.wantErr {
				t.Errorf("Change.diffChannelTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !equalChanges(*change, tt.want) {
				t.Errorf("Changelog.diff() = %v, want %v", *change, tt.want)
			}
		})
	}
}

func TestChange_diffMetadata(t *testing.T) {
	type fields struct {
		Context string
		Diff    []Difference
	}
	type args struct {
		parentKind    Kind
		metaOrig      meta.Metadata
		metaCurr      meta.Metadata
		ctx           string
		parentElement string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    Change
	}{
		{
			name: "Ucnhanged Metadata",
			fields: fields{
				Context: "",
			},
			args: args{
				metaOrig: meta.Metadata{
					Name:      "name",
					Reference: "reference",
					Parent:    "",
					SHA256:    "1234567890",
				},
				metaCurr: meta.Metadata{
					Name:      "name",
					Reference: "reference",
					Parent:    "",
					SHA256:    "1234567890",
				},
				ctx:        "",
				parentKind: AppKind,
			},
			wantErr: false,
			want:    Change{},
		},
		{
			name: "Valid changed Metadata",
			fields: fields{
				Context: "",
			},
			args: args{
				metaOrig: meta.Metadata{
					Name:      "name",
					Reference: "reference",
					Parent:    "",
					SHA256:    "1234567890",
				},
				metaCurr: meta.Metadata{
					Name:      "name",
					Reference: "referenceDiff",
					Parent:    "",
					SHA256:    "1234567890Diff",
					Annotations: map[string]string{
						"1": "1",
					},
				},
				ctx:        "",
				parentKind: NodeKind,
			},
			wantErr: false,
			want: Change{
				Kind:      MetaKind | NodeKind | AnnotationKind,
				Operation: Create | Update,
				Diff: []Difference{
					{
						Field:     "Meta.Reference",
						From:      "reference",
						To:        "referenceDiff",
						Kind:      MetaKind | NodeKind,
						Operation: Update,
					},
					{
						Field:     "Meta.SHA256",
						From:      "1234567890",
						To:        "1234567890Diff",
						Kind:      MetaKind | NodeKind,
						Operation: Update,
					},
					{
						Field:     "Meta.Annotations[1]",
						From:      "<nil>",
						To:        "1",
						Kind:      MetaKind | NodeKind | AnnotationKind,
						Operation: Create,
						Name:      "1",
					},
				},
			},
		},
		{
			name: "Change parent Metadata error",
			fields: fields{
				Context: "",
			},
			args: args{
				metaOrig: meta.Metadata{
					Name:      "name",
					Reference: "reference",
					Parent:    "",
					SHA256:    "1234567890",
				},
				metaCurr: meta.Metadata{
					Name:      "name",
					Reference: "reference",
					Parent:    "Err",
					SHA256:    "1234567890",
				},
				ctx:        "",
				parentKind: AppKind,
			},
			wantErr: false,
			want: Change{
				Kind:      MetaKind | AppKind,
				Operation: Update,
				Diff: []Difference{
					{
						Field:     "Meta.Parent",
						From:      "",
						To:        "Err",
						Kind:      MetaKind | AppKind,
						Operation: Update,
					},
				},
			},
		},
		{
			name: "Change name Metadata error",
			fields: fields{
				Context: "",
			},
			args: args{
				metaOrig: meta.Metadata{
					Name:      "name",
					Reference: "reference",
					Parent:    "",
					SHA256:    "1234567890",
				},
				metaCurr: meta.Metadata{
					Name:      "Err",
					Reference: "reference",
					Parent:    "",
					SHA256:    "1234567890",
				},
				ctx:        "",
				parentKind: ChannelKind,
			},
			wantErr: false,
			want: Change{
				Kind:      MetaKind | ChannelKind,
				Operation: Update,
				Diff: []Difference{
					{
						Field:     "Meta.Name",
						From:      "name",
						To:        "Err",
						Kind:      MetaKind | ChannelKind,
						Operation: Update,
					},
				},
			},
		},
		{
			name: "delete annotation update",
			fields: fields{
				Context: "",
			},
			args: args{
				metaOrig: meta.Metadata{
					Name:      "name",
					Reference: "reference",
					Parent:    "",
					SHA256:    "1234567890",
					Annotations: map[string]string{
						"1": "1",
						"2": "1",
					},
				},
				metaCurr: meta.Metadata{
					Name:      "name",
					Reference: "reference",
					Parent:    "",
					SHA256:    "1234567890",
					Annotations: map[string]string{
						"1": "1",
					},
				},
				ctx:        "",
				parentKind: NodeKind,
			},
			wantErr: false,
			want: Change{
				Operation: Delete,
				Kind:      MetaKind | NodeKind | AnnotationKind,
				Diff: []Difference{
					{
						Field:     "Meta.Annotations[2]",
						From:      "1",
						To:        "<nil>",
						Kind:      MetaKind | NodeKind | AnnotationKind,
						Operation: Delete,
						Name:      "2",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			change := &Change{
				Context:   tt.fields.Context,
				Diff:      tt.fields.Diff,
				changelog: &Changelog{},
			}
			if err := change.diffMetadata(tt.args.parentElement, tt.args.parentKind, tt.args.metaOrig, tt.args.metaCurr, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Change.diffMetadata() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !equalChanges(*change, tt.want) {
				t.Errorf("Changelog.diff() = %v, want %v", *change, tt.want)
			}
		})
	}
}

func TestChangelog_Print(t *testing.T) {
	expectedOut := bytes.NewBufferString("On: abc\n")
	w := tabwriter.NewWriter(expectedOut, 12, 0, 3, ' ', tabwriter.Debug)

	// generating the expected output
	fmt.Fprintln(w, "Field\t From\t To")
	fmt.Fprintf(
		w,
		"%s\t %s\t %s\n",
		"mock_field",
		"unmocked",
		"mocked",
	)
	w.Flush()

	tests := []struct {
		name    string
		cl      Changelog
		wantOut string
	}{
		{
			name: "basic_changelog_print",
			cl: Changelog{
				Change{
					Context: "abc",
					Diff: []Difference{
						{
							Field: "mock_field",
							From:  "unmocked",
							To:    "mocked",
						},
					},
				},
			},
			wantOut: expectedOut.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			tt.cl.Print(out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("Changelog.Print() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestChangelog_diff(t *testing.T) {
	type args struct {
		appOrig *meta.App
		appCurr *meta.App
		ctx     string
	}
	tests := []struct {
		name    string
		cl      Changelog
		args    args
		want    Changelog
		wantErr bool
	}{
		{
			name: "diff_output_test",
			cl:   Changelog{},
			args: args{
				appOrig: getMockRootApp(),
				appCurr: getMockRootApp2(),
				ctx:     "",
			},
			want: Changelog{
				{
					Context: "*",
					Diff: []Difference{
						{
							Field: "Meta.SHA256",
							From:  "1",
							To:    "2",
						},
						{
							Field: "Meta.Annotations[an1]",
							From:  "<nil>",
							To:    "a",
						},
						{
							Field: "Meta.Annotations[an2]",
							From:  "<nil>",
							To:    "b",
						},
						{
							Field: "Spec.Apps[app1]",
							From:  "{...}",
							To:    "<nil>",
						},
						{
							Field: "Spec.Channels[ch2]",
							From:  "{...}",
							To:    "<nil>",
						},
						{
							Field: "Spec.Channels[ch1].Meta.Reference",
							From:  "root.ch1",
							To:    "root.ch1diff",
						},
						{
							Field: "Spec.ChannelTypes[ct2]",
							From:  "{...}",
							To:    "<nil>",
						},
						{
							Field: "Spec.ChannelTypes[ct1].Meta.Reference",
							From:  "root.ct1",
							To:    "root.ct1diff",
						},
					},
				},
				{
					Context: "*.Spec.Apps.app2.Spec.Apps.app3",
					Diff: []Difference{
						{
							Field: "Spec.Node.Spec.Image",
							From:  "imageNodeApp3",
							To:    "imageNodeApp3diff",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "diff_output_test_want_error",
			cl:   Changelog{},
			args: args{
				appOrig: &meta.App{Meta: meta.Metadata{Name: "abc"}},
				appCurr: &meta.App{Meta: meta.Metadata{Name: "dbc"}},
				ctx:     "",
			},
			want:    Changelog{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cl.diff(
				tt.args.appOrig,
				tt.args.appCurr,
				tt.args.ctx,
			)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Changelog.diff() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}

			if !cmp.Equal(got, tt.want, pkgUtils.GeneralCompareOptions()) {
				t.Errorf(
					"Changelog.diff() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}

func getMockRootApp() *meta.App {
	root := meta.App{
		Meta: meta.Metadata{
			Name:        "",
			Reference:   "",
			Annotations: map[string]string{},
			Parent:      "",
			SHA256:      "1",
		},
		Spec: meta.AppSpec{
			Node: meta.Node{},
			Apps: map[string]*meta.App{
				"app1": {
					Meta: meta.Metadata{
						Name:        "app1",
						Reference:   "app1",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{
							Meta: meta.Metadata{
								Name:        "nodeApp1",
								Reference:   "app1.nodeApp1",
								Annotations: map[string]string{},
								Parent:      "app1",
								SHA256:      "",
							},
							Spec: meta.NodeSpec{
								Image: "imageNodeApp1",
							},
						},
						Apps:         map[string]*meta.App{},
						Channels:     map[string]*meta.Channel{},
						ChannelTypes: map[string]*meta.ChannelType{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1"},
							Output: []string{"ch2"},
						},
					},
				},
				"app2": {
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "app2",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{},
						Apps: map[string]*meta.App{
							"app3": {
								Meta: meta.Metadata{
									Name:        "app3",
									Reference:   "app2.app3",
									Annotations: map[string]string{},
									Parent:      "app2",
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "nodeApp3",
											Reference:   "app3.nodeApp2",
											Annotations: map[string]string{},
											Parent:      "app3",
											SHA256:      "",
										},
										Spec: meta.NodeSpec{
											Image: "imageNodeApp3",
										},
									},
									Apps:         map[string]*meta.App{},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
									Boundary: meta.AppBoundary{
										Input:  []string{"ch1"},
										Output: []string{"ch2"},
									},
								},
							},
						},
						Channels: map[string]*meta.Channel{
							"ch1app2": {
								Meta: meta.Metadata{
									Name:   "ch1app2",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
							"ch2app2": {
								Meta: meta.Metadata{
									Name:   "ch2app2",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						ChannelTypes: map[string]*meta.ChannelType{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1"},
							Output: []string{"ch2"},
						},
					},
				},
			},
			Channels: map[string]*meta.Channel{
				"ch1": {
					Meta: meta.Metadata{
						Name:        "ch1",
						Reference:   "root.ch1",
						Annotations: map[string]string{},
						Parent:      "root",
						SHA256:      "",
					},
				},
				"ch2": {
					Meta: meta.Metadata{
						Name:        "ch2",
						Reference:   "root.ch1diff",
						Annotations: map[string]string{},
						Parent:      "root",
						SHA256:      "",
					},
				},
			},
			ChannelTypes: map[string]*meta.ChannelType{
				"ct1": {
					Meta: meta.Metadata{
						Name:        "ct1",
						Reference:   "root.ct1",
						Annotations: map[string]string{},
						Parent:      "root",
						SHA256:      "",
					},
					Schema: "",
				},
				"ct2": {
					Meta: meta.Metadata{
						Name:        "ct2",
						Reference:   "root.ct2",
						Annotations: map[string]string{},
						Parent:      "root",
						SHA256:      "",
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

func getMockRootApp2() *meta.App {
	root := meta.App{
		Meta: meta.Metadata{
			Name:      "",
			Reference: "",
			Annotations: map[string]string{
				"an1": "a",
				"an2": "b",
			},
			Parent: "",
			SHA256: "2",
		},
		Spec: meta.AppSpec{
			Node: meta.Node{},
			Apps: map[string]*meta.App{
				"app2": {
					Meta: meta.Metadata{
						Name:        "app2",
						Reference:   "app2",
						Annotations: map[string]string{},
						Parent:      "",
						SHA256:      "",
					},
					Spec: meta.AppSpec{
						Node: meta.Node{},
						Apps: map[string]*meta.App{
							"app3": {
								Meta: meta.Metadata{
									Name:        "app3",
									Reference:   "app2.app3",
									Annotations: map[string]string{},
									Parent:      "app2",
									SHA256:      "",
								},
								Spec: meta.AppSpec{
									Node: meta.Node{
										Meta: meta.Metadata{
											Name:        "nodeApp3",
											Reference:   "app3.nodeApp2",
											Annotations: map[string]string{},
											Parent:      "app3",
											SHA256:      "",
										},
										Spec: meta.NodeSpec{
											Image: "imageNodeApp3diff",
										},
									},
									Apps:         map[string]*meta.App{},
									Channels:     map[string]*meta.Channel{},
									ChannelTypes: map[string]*meta.ChannelType{},
									Boundary: meta.AppBoundary{
										Input:  []string{"ch1"},
										Output: []string{"ch2"},
									},
								},
							},
						},
						Channels: map[string]*meta.Channel{
							"ch1app2": {
								Meta: meta.Metadata{
									Name:   "ch1app2",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
							"ch2app2": {
								Meta: meta.Metadata{
									Name:   "ch2app2",
									Parent: "",
								},
								Spec: meta.ChannelSpec{},
							},
						},
						ChannelTypes: map[string]*meta.ChannelType{},
						Boundary: meta.AppBoundary{
							Input:  []string{"ch1"},
							Output: []string{"ch2"},
						},
					},
				},
			},
			Channels: map[string]*meta.Channel{
				"ch1": {
					Meta: meta.Metadata{
						Name:        "ch1",
						Reference:   "root.ch1diff",
						Annotations: map[string]string{},
						Parent:      "root",
						SHA256:      "",
					},
				},
			},
			ChannelTypes: map[string]*meta.ChannelType{
				"ct1": {
					Meta: meta.Metadata{
						Name:        "ct1",
						Reference:   "root.ct1diff",
						Annotations: map[string]string{},
						Parent:      "root",
						SHA256:      "",
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

func equalChangelogs(changelog1 Changelog, changelog2 Changelog) bool {
	if len(changelog1) != len(changelog2) {
		return false
	}
	map1 := make(map[string]map[Difference]int)
	map2 := make(map[string]map[Difference]int)
	for _, change := range changelog1 {
		map1[change.Context] = make(map[Difference]int)
		for _, diff := range change.Diff {
			map1[change.Context][diff]++
		}
	}

	for _, change := range changelog2 {
		map2[change.Context] = make(map[Difference]int)
		for _, diff := range change.Diff {
			map2[change.Context][diff]++
		}
	}

	for key, val := range map1 {
		if len(map1[key]) != len(map2[key]) {
			return false
		}

		for key2, val2 := range val {
			if map2[key][key2] != val2 {
				return false
			}
		}
	}
	return true
}

func equalChanges(change1 Change, change2 Change) bool {
	if change1.Context != change2.Context {
		return false
	}
	if change1.Kind != change2.Kind {
		return false
	}
	if change1.Operation != change2.Operation {
		return false
	}
	if len(change1.Diff) != len(change2.Diff) {
		return false
	}
	map1 := make(map[Difference]int)
	map2 := make(map[Difference]int)
	for _, diff := range change1.Diff {
		map1[diff]++
	}

	for _, diff := range change2.Diff {
		map2[diff]++
	}

	for key, val := range map1 {
		if map2[key] != val {
			return false
		}
	}
	return true
}
