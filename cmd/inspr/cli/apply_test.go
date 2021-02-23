package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gopkg.in/yaml.v2"
)

func createYaml() (string, meta.Channel) {
	channel := meta.Channel{
		Meta: meta.Metadata{
			Name:        "mock_name",
			Reference:   "mock_reference",
			Annotations: map[string]string{},
			Parent:      "mock_parent",
			SHA256:      "mock_sha256",
		},
		Spec:          meta.ChannelSpec{Type: "mock_type"},
		ConnectedApps: []string{"a", "b", "c"},
	}
	data, _ := yaml.Marshal(&channel)
	return string(data), channel
}

// TestNewApplyCmd is mainly for improving test coverage,
// it was really tested by instantiating Inspr's CLI
func TestNewApplyCmd(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		{
			name: "Creates a new Cobra command",
			want: &cobra.Command{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewApplyCmd()
			if got == nil {
				t.Errorf("NewApplyCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isYaml(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Given file is yaml",
			args: args{
				file: "itsAYaml.yaml",
			},
			want: true,
		},
		{
			name: "Given file is yml",
			args: args{
				file: "itsAYml.yml",
			},
			want: true,
		},
		{
			name: "Given file is another extention",
			args: args{
				file: "itsNotAYaml.txt",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isYaml(tt.args.file); got != tt.want {
				t.Errorf("isYaml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_printAppliedFiles(t *testing.T) {
	type args struct {
		appliedFiles []applied
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			name: "Prints a valid file",
			args: args{
				[]applied{{
					file: "aFile.yaml",
					component: meta.Component{
						Kind:       "randKind",
						APIVersion: "v1",
					},
				}},
			},
			wantOut: "Applying: \naFile.yaml | randKind | v1\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			printAppliedFiles(tt.args.appliedFiles, out)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("printAppliedFiles() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
