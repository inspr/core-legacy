package cli

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
	cliutils "inspr.dev/inspr/pkg/cmd/utils"
	"inspr.dev/inspr/pkg/ierrors"
	"inspr.dev/inspr/pkg/meta"
)

const (
	filePath = "filetest.yaml"
)

func createDAppYaml() string {
	comp := meta.Component{
		Kind:       "dapp",
		APIVersion: "v1",
	}
	data, _ := yaml.Marshal(&comp)
	return string(data)
}

func createChannelYaml() string {
	comp := meta.Component{
		Kind:       "channel",
		APIVersion: "v1",
	}
	data, _ := yaml.Marshal(&comp)
	return string(data)
}

func createTypeYaml() string {
	comp := meta.Component{
		Kind:       "type",
		APIVersion: "v1",
	}
	data, _ := yaml.Marshal(&comp)
	return string(data)
}

func createAliasYaml() string {
	comp := meta.Component{
		Kind:       "alias",
		APIVersion: "v1",
	}
	data, _ := yaml.Marshal(&comp)
	return string(data)
}

func createInvalidYaml() string {
	comp := meta.Component{
		Kind:       "none",
		APIVersion: "",
	}
	data, _ := yaml.Marshal(&comp)
	return string(data)
}

func getCurrentFilesInFolder() []string {
	var files []string
	folder, _ := ioutil.ReadDir(".")

	for _, file := range folder {
		files = append(files, file.Name())
	}
	return files
}

// TestNewApplyCmd is mainly for improving test coverage,
// it was really tested by instantiating Inspr's CLI
func TestNewApplyCmd(t *testing.T) {
	prepareToken(t)
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
	prepareToken(t)
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
			name: "Given file is another extension",
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
	prepareToken(t)
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
					fileName: "aFile.yaml",
					component: meta.Component{
						Kind:       "randKind",
						APIVersion: "v1",
					},
				}},
			},
			wantOut: "\nApplied:\naFile.yaml | randKind | v1\n",
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
	prepareToken(t)
	defer os.Remove(filePath)
	yamlString := createDAppYaml()

	// creates a file with the expected syntax
	ioutil.WriteFile(
		filePath,
		[]byte(yamlString),
		os.ModePerm,
	)

	tests := []struct {
		name           string
		flagsAndArgs   []string
		expectedOutput string
	}{
		{
			name:           "Should apply the file",
			flagsAndArgs:   []string{"-f", "filetest.yaml"},
			expectedOutput: "\nApplied:\nfiletest.yaml | dapp | v1\n",
		},
		{
			name:           "Too many flags, should raise an error",
			flagsAndArgs:   []string{"-f", "example", "-k", "example"},
			expectedOutput: "Invalid command call\nFor help, type 'insprctl apply --help'\n",
		},
		{
			name:           "No files applied",
			flagsAndArgs:   []string{"-f", "example.yaml"},
			expectedOutput: "No files were applied\nFiles to be applied must be .yaml or .yml\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetFactory().Subscribe(meta.Component{
				APIVersion: "v1",
				Kind:       "dapp",
			},
				func(b []byte, out io.Writer) error {
					return nil
				})

			cmd := NewApplyCmd()
			buf := bytes.NewBufferString("")

			cliutils.SetOutput(buf)
			cmd.SetArgs(tt.flagsAndArgs)

			cmd.Execute()
			got := buf.String()

			if !reflect.DeepEqual(got, tt.expectedOutput) {
				t.Errorf(
					"doApply() = %v, want\n%v",
					got,
					tt.expectedOutput,
				)
			}
		})
	}
}

func Test_getFilesFromFolder(t *testing.T) {
	prepareToken(t)
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
			want:    getCurrentFilesInFolder(),
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
	prepareToken(t)
	defer os.Remove(filePath)
	tempFiles := []string{filePath}
	yamlString := createDAppYaml()

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
		name    string
		args    args
		want    []applied
		funcErr error
		errMsg  string
	}{
		{
			name: "Get file from current folder",
			args: args{
				path:  "",
				files: tempFiles,
			},
			want: []applied{{
				fileName: filePath,
				component: meta.Component{
					Kind:       "dapp",
					APIVersion: "v1",
				},
				content: []byte(yamlString),
			}},
			funcErr: nil,
			errMsg:  "",
		},
		{
			name: "Unauthorized_error",
			args: args{
				path:  "",
				files: tempFiles,
			},
			want:    nil,
			funcErr: ierrors.New("unauthorized").Unauthorized(),
			errMsg:  "failed to authenticate with the cluster. Is your token configured correctly?\n",
		},
		{
			name: "forbidden_error",
			args: args{
				path:  "",
				files: tempFiles,
			},
			want:    nil,
			funcErr: ierrors.New("forbidden").Forbidden(),
			errMsg:  "forbidden operation, please check for the scope.\n",
		},
		{
			name: "default_ierror_message",
			args: args{
				path:  "",
				files: tempFiles,
			},
			want:    nil,
			funcErr: ierrors.New("default_error").BadRequest(),
			errMsg:  "unexpected inspr error: default_error\n",
		},
		{
			name: "Unknown_error",
			args: args{
				path:  "",
				files: tempFiles,
			},
			want:    nil,
			funcErr: errors.New("unknown_Error"),
			errMsg:  "non inspr error: unknown_Error\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			applyFactory = nil

			GetFactory().Subscribe(
				meta.Component{
					Kind:       "dapp",
					APIVersion: "v1",
				},
				func(b []byte, out io.Writer) error {
					return tt.funcErr
				},
			)

			got := applyValidFiles(tt.args.path, tt.args.files, out)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"applyValidFiles() = %v, want %v",
					got,
					tt.want,
				)
			}

			if tt.funcErr != nil && tt.errMsg != out.String() {
				t.Errorf(
					"error messages incompatible\n got => %v\n want => %v\n",
					out.String(),
					tt.errMsg,
				)
			}
		})
	}
}

func Test_getOrderedFiles(t *testing.T) {
	prepareToken(t)
	defer os.Remove("app.yml")
	defer os.Remove("ch.yml")
	defer os.Remove("ct.yml")
	defer os.Remove("al.yml")
	defer os.Remove("invalid.yml")
	tempFiles := []string{"app.yml", "invalid.yml",
		"ch.yml", "ct.yml", "al.yml"}
	// creates a file with the expected syntax
	ioutil.WriteFile(
		"app.yml",
		[]byte(createDAppYaml()),
		os.ModePerm,
	)
	ioutil.WriteFile(
		"ch.yml",
		[]byte(createChannelYaml()),
		os.ModePerm,
	)
	ioutil.WriteFile(
		"ct.yml",
		[]byte(createTypeYaml()),
		os.ModePerm,
	)
	ioutil.WriteFile(
		"al.yml",
		[]byte(createAliasYaml()),
		os.ModePerm,
	)
	ioutil.WriteFile(
		"invalid.yml",
		[]byte(createInvalidYaml()),
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
			name: "Return ordered files",
			args: args{
				path:  ".",
				files: tempFiles,
			},
			want: orderedContent(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOrderedFiles(tt.args.path, tt.args.files); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOrderedFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func orderedContent() []applied {
	ordered := []applied{
		{
			fileName: "app.yml",
			component: meta.Component{
				Kind:       "dapp",
				APIVersion: "v1",
			},
			content: []byte(createDAppYaml()),
		},
		{
			fileName: "ct.yml",
			component: meta.Component{
				Kind:       "type",
				APIVersion: "v1",
			},
			content: []byte(createTypeYaml()),
		},
		{
			fileName: "ch.yml",
			component: meta.Component{
				Kind:       "channel",
				APIVersion: "v1",
			},
			content: []byte(createChannelYaml()),
		},
		{
			fileName: "al.yml",
			component: meta.Component{
				Kind:       "alias",
				APIVersion: "v1",
			},
			content: []byte(createAliasYaml()),
		},
	}

	return ordered
}
