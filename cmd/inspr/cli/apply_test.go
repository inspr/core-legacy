package cli

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gopkg.in/yaml.v2"
)

const (
	filePath = "filetest.yaml"
)

func createYaml() string {
	comp := meta.Component{
		Kind:       "app",
		APIVersion: "v1",
	}
	data, _ := yaml.Marshal(&comp)
	return string(data)
}

// TestNewApplyCmd is mainly for improving test coverage,
// it was really tested by instantiating Inspr's CLI
func TestNewApplyCmd(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Creates a new Cobra command",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewApplyCmd()
			if got == nil {
				t.Errorf("NewApplyCmd() = %v", got)
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

func Test_doApply(t *testing.T) {
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			if err := doApply(tt.args.in0, out); (err != nil) != tt.wantErr {
				t.Errorf("doApply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("doApply() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_getFilesFromFolder(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []string
	}{
		{
			name: "Get file from current folder",
			args: args{
				path: ".",
			},
			wantErr: false,
			want: []string{"apply.go", "apply_factory.go", "apply_factory_test.go",
				"apply_test.go", "cli.go", "cli_test.go", "hidden_cmd.go", "mock_cmd.go"},
		},
		{
			name: "Invalid - path doesn't exist",
			args: args{
				path: "invalid/",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFilesFromFolder(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFilesFromFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFilesFromFolder() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func Test_applyValidFiles(t *testing.T) {
	tempFiles := []string{filePath}
	yamlString := createYaml()
	// creates a file with the expected syntax
	ioutil.WriteFile(
		filePath,
		[]byte(yamlString),
		os.ModePerm,
	)

	type args struct {
		path  string
		files []string
	}
	tests := []struct {
		name string
		args args
		want []applied
	}{
		{
			name: "Get file from current folder",
			args: args{
				path:  "",
				files: tempFiles,
			},
			want: []applied{{
				file: filePath,
				component: meta.Component{
					Kind:       "app",
					APIVersion: "v1",
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			GetFactory().Subscribe(meta.Component{
				APIVersion: "v1",
				Kind:       "app",
			},
				func(b []byte) error {
					ch := meta.Channel{}

					yaml.Unmarshal(b, &ch)
					fmt.Println(ch)

					return nil
				})
			if got := applyValidFiles(tt.args.path, tt.args.files, out); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applyValidFiles() = %v, want %v", got, tt.want)
			}
		})
	}
	os.Remove(filePath)
}
