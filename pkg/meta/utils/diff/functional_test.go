package diff

import (
	"reflect"
	"strings"
	"testing"
)

func TestChangelog_ForEach(t *testing.T) {
	var contexts []string
	type args struct {
		f func(c Change) error
	}
	tests := []struct {
		name     string
		c        Changelog
		args     args
		wantFunc func(t *testing.T)
	}{
		{
			"get context from changes",
			Changelog{
				Change{
					Context: "context1",
				},
				Change{
					Context: "context2",
				},
				Change{
					Context: "context3",
				},
				Change{
					Context: "context4",
				},
			},
			args{
				f: func(c Change) error {
					contexts = append(contexts, c.Context)
					return nil
				},
			},
			func(t *testing.T) {
				want := []string{"context1", "context2", "context3", "context4"}
				if !reflect.DeepEqual(want, contexts) {
					t.Errorf("changelog for each %v != %v", want, contexts)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.ForEach(tt.args.f)
			tt.wantFunc(t)
		})
	}
}

func TestChangelog_FilterDiffs(t *testing.T) {
	type args struct {
		comp func(c string, d Difference) bool
	}
	tests := []struct {
		name string
		c    Changelog
		args args
		want Changelog
	}{
		{
			name: "filter by operation",
			c: Changelog{
				Change{
					Context: "context",
					Diff: []Difference{
						{
							Name:      "diff1",
							Operation: Create | Delete,
						},
						{
							Name:      "diff2",
							Operation: Update,
						},
						{
							Name:      "diff3",
							Operation: Update | Create,
						},
					},
				},
				Change{
					Context: "context2",
					Diff: []Difference{
						{
							Name:      "diff4",
							Operation: Create,
						},
						{
							Name:      "diff5",
							Operation: Delete,
						},
						{
							Name:      "diff6",
							Operation: Update,
						},
					},
				},
			},
			args: args{
				comp: func(c string, d Difference) bool {
					return d.Operation&Update > 0
				},
			},
			want: Changelog{
				Change{
					Context: "context",
					Diff: []Difference{
						{
							Name:      "diff2",
							Operation: Update,
						},
						{
							Name:      "diff3",
							Operation: Update | Create,
						},
					},
				},
				Change{
					Context: "context2",
					Diff: []Difference{
						{
							Name:      "diff6",
							Operation: Update,
						},
					},
				},
			},
		},
		{
			name: "filter by name",
			c: Changelog{
				Change{
					Context: "context",
					Diff: []Difference{
						{
							Name:      "next_diff1",
							Operation: Create | Delete,
						},
						{
							Name:      "diff2",
							Operation: Update,
						},
						{
							Name:      "diff3",
							Operation: Update | Create,
						},
					},
				},
				Change{
					Context: "context2",
					Diff: []Difference{
						{
							Name:      "diff4",
							Operation: Create,
						},
						{
							Name:      "next_diff5",
							Operation: Delete,
						},
						{
							Name:      "diff6",
							Operation: Update,
						},
					},
				},
			},
			args: args{
				comp: func(c string, d Difference) bool {
					return strings.HasPrefix(d.Name, "next")
				},
			},
			want: Changelog{
				Change{
					Context: "context",
					Diff: []Difference{
						{
							Name:      "next_diff1",
							Operation: Create | Delete,
						},
					},
				},
				Change{
					Context: "context2",
					Diff: []Difference{
						{
							Name:      "next_diff5",
							Operation: Delete,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.FilterDiffs(tt.args.comp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Changelog.FilterDiffs() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestChangelog_FilterDiffsByKind(t *testing.T) {
	type args struct {
		kind Kind
	}
	tests := []struct {
		name string
		c    Changelog
		args args
		want Changelog
	}{
		{
			name: "single kind",
			c: Changelog{
				Change{
					Context: "context1",
					Diff: []Difference{
						{
							Name: "diff1",
							Kind: ChannelKind | MetaKind,
						},
						{
							Name: "diff2",
							Kind: NodeKind,
						},
						{
							Name: "diff3",
							Kind: AppKind | MetaKind,
						},
						{
							Name: "diff4",
							Kind: ChannelTypeKind | MetaKind,
						},
					},
				},
				Change{
					Context: "context2",
					Diff: []Difference{
						{
							Name: "diff5",
							Kind: ChannelKind,
						},
						{
							Name: "diff6",
							Kind: NodeKind,
						},
						{
							Name: "diff7",
							Kind: AnnotationKind,
						},
						{
							Name: "diff8",
							Kind: ChannelTypeKind | MetaKind,
						},
					},
				},
			},
			args: args{
				kind: MetaKind,
			},
			want: Changelog{
				Change{
					Context: "context1",
					Diff: []Difference{
						{
							Name: "diff1",
							Kind: ChannelKind | MetaKind,
						},
						{
							Name: "diff3",
							Kind: AppKind | MetaKind,
						},
						{
							Name: "diff4",
							Kind: ChannelTypeKind | MetaKind,
						},
					},
				},
				Change{
					Context: "context2",
					Diff: []Difference{
						{
							Name: "diff8",
							Kind: ChannelTypeKind | MetaKind,
						},
					},
				},
			},
		},
		{
			name: "multiple kinds",
			c: Changelog{
				Change{
					Context: "context1",
					Diff: []Difference{
						{
							Name: "diff1",
							Kind: ChannelKind | MetaKind,
						},
						{
							Name: "diff2",
							Kind: NodeKind,
						},
						{
							Name: "diff3",
							Kind: AppKind | MetaKind,
						},
						{
							Name: "diff4",
							Kind: ChannelTypeKind | MetaKind,
						},
					},
				},
				Change{
					Context: "context2",
					Diff: []Difference{
						{
							Name: "diff5",
							Kind: ChannelKind,
						},
						{
							Name: "diff6",
							Kind: NodeKind,
						},
						{
							Name: "diff7",
							Kind: AnnotationKind,
						},
						{
							Name: "diff8",
							Kind: ChannelTypeKind | MetaKind,
						},
					},
				},
			},
			args: args{
				kind: ChannelKind | AnnotationKind,
			},
			want: Changelog{
				Change{
					Context: "context1",
					Diff: []Difference{
						{
							Name: "diff1",
							Kind: ChannelKind | MetaKind,
						},
					},
				},
				Change{
					Context: "context2",
					Diff: []Difference{
						{
							Name: "diff5",
							Kind: ChannelKind,
						},
						{
							Name: "diff7",
							Kind: AnnotationKind,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.FilterDiffsByKind(tt.args.kind); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Changelog.FilterDiffsByKind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDifferenceKindOperation(t *testing.T) {
	var didDo int
	type args struct {
		kind  Kind
		apply func(scope string, d Difference) error
	}
	tests := []struct {
		name     string
		args     args
		wantFunc func(t *testing.T, got DifferenceReaction)
	}{
		{
			name: "one kind",
			args: args{
				kind: AnnotationKind,
				apply: func(scope string, d Difference) error {
					didDo = 1
					return nil
				},
			},
			wantFunc: func(t *testing.T, got DifferenceReaction) {
				if !got.filter("", Difference{Kind: AnnotationKind}) {
					t.Error("did not apply filter correctly")
				}
				if got.filter("", Difference{Kind: ^AnnotationKind}) {
					t.Error("did not apply filter correctly")
				}
				got.operation("", Difference{})
				if didDo != 1 {
					t.Error("did not apply correctly")
				}
			},
		},

		{
			name: "multiple kinds",
			args: args{
				kind: AnnotationKind | AppKind,
				apply: func(scope string, d Difference) error {
					didDo = 2
					return nil
				},
			},
			wantFunc: func(t *testing.T, got DifferenceReaction) {
				if !got.filter("", Difference{Kind: AnnotationKind}) {
					t.Error("did not apply filter correctly")
				}
				if !got.filter("", Difference{Kind: AppKind}) {
					t.Error("did not apply filter correctly")
				}
				if got.filter("", Difference{Kind: ^(AnnotationKind | AppKind)}) {
					t.Error("did not apply filter correctly")
				}
				got.operation("", Difference{})
				if didDo != 2 {
					t.Error("did not apply correctly")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDifferenceKindReaction(tt.args.kind, tt.args.apply)
			tt.wantFunc(t, got)
		})
	}
}
