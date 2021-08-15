package cli

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
	"inspr.dev/inspr/pkg/cmd/utils"
)

// TestNewInsprCommand is mainly for improving test coverage,
// it was really tested by instantiating Inspr's CLI
func TestNewInsprCommand(t *testing.T) {
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
			out := &bytes.Buffer{}
			err := &bytes.Buffer{}
			got := NewInsprCommand(out, err, "")
			if got == nil {
				t.Errorf("NewInsprCommand() = %v", got)
			}
		})
	}
}

func prepareConfig(f string, c struct {
	ServerIP string `yaml:"serverip"`
	Scope    string `yaml:"scope"`
}) {
	os.Mkdir(filepath.Join(f, ".inspr"), 0755)
	f = filepath.Join(f, ".inspr", "config")
	arr, _ := yaml.Marshal(c)
	ioutil.WriteFile(f, arr, 0644)
}

var defaultConfig struct {
	ServerIP string `yaml:"serverip"`
	Scope    string `yaml:"scope"`
} = struct {
	ServerIP string "yaml:\"serverip\""
	Scope    string "yaml:\"scope\""
}{
	ServerIP: "http://localhost:8080",
}

func Test_mainCmdPreRun(t *testing.T) {
	prepareToken(t)
	folder := t.TempDir()
	prev := os.Getenv("HOME")
	os.Setenv("HOME", folder)
	defer os.Setenv("HOME", prev)
	utils.InitViperConfig()
	prepareConfig(folder, defaultConfig)

	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "pre run test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString("")
			if err := mainCmdPreRun(NewInsprCommand(buf, buf, ""), tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf(
					"mainCmdPreRun() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}
