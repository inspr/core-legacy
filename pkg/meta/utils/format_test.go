package utils

import (
	"bytes"
	"testing"

	"github.com/disiqueira/gotree"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			PrintAppTree(tt.args.app, out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("PrintAppTree() = %v, want %v", gotOut, tt.wantOut)
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			PrintChannelTree(tt.args.ch, out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("PrintChannelTree() = %v, want %v", gotOut, tt.wantOut)
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			PrintChannelTypeTree(tt.args.ct, out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("PrintChannelTypeTree() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_populateMeta(t *testing.T) {
	type args struct {
		metaTree gotree.Tree
		meta     *meta.Metadata
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			populateMeta(tt.args.metaTree, tt.args.meta)
		})
	}
}
