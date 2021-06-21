package cli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	cliutils "inspr.dev/inspr/pkg/cmd/utils"
)

func TestNewConfigChangeCmd(t *testing.T) {
	prepareToken(t)
	tests := []struct {
		name          string
		checkFunction func(t *testing.T, got *cobra.Command)
	}{
		{
			name: "It should create a new config change command",
			checkFunction: func(t *testing.T, got *cobra.Command) {
				if got == nil {
					t.Errorf("NewConfigChangeCmd() not created successfully")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewConfigChangeCmd()
			if tt.checkFunction != nil {
				tt.checkFunction(t, got)
			}
		})
	}
}

func TestNewListConfig(t *testing.T) {
	prepareToken(t)
	tests := []struct {
		name          string
		checkFunction func(t *testing.T, got *cobra.Command)
	}{
		{
			name: "It should create a new config list command",
			checkFunction: func(t *testing.T, got *cobra.Command) {
				if got == nil {
					t.Errorf("NewListConfigCmd() not created successfully")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewListConfig()
			if tt.checkFunction != nil {
				tt.checkFunction(t, got)
			}
		})
	}
}

func Test_doConfigChange(t *testing.T) {
	prepareToken(t)
	folder := setupViperTest(t)
	os.Setenv("HOME", folder)
	defer os.Unsetenv("HOME")
	cliutils.ReadDefaultConfig()

	bufResp := bytes.NewBufferString("")
	fmt.Fprintf(bufResp, "Success: insprctl config [%v] changed to '%v'\n", "key_example", "new_value")
	outResp, _ := ioutil.ReadAll(bufResp)

	bufResp2 := bytes.NewBufferString("")
	fmt.Fprintf(bufResp2, "error: key inserted does not exist in the insprctl config\n")
	cliutils.SetOutput(bufResp2)
	printExistingKeys()
	outResp2, _ := ioutil.ReadAll(bufResp2)

	tests := []struct {
		name           string
		flagsAndArgs   []string
		expectedOutput []byte
	}{
		{
			name:           "Key doens't exist",
			flagsAndArgs:   []string{"invalid_key", "some_value"},
			expectedOutput: outResp2,
		},
		{
			name:           "Should make the change and print the successful message",
			flagsAndArgs:   []string{"key_example", "new_value"},
			expectedOutput: outResp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewConfigChangeCmd()
			buf := bytes.NewBufferString("")

			cliutils.SetOutput(buf)
			cmd.SetArgs(tt.flagsAndArgs)

			cmd.Execute()
			got, _ := ioutil.ReadAll(buf)

			if len(got) != len(tt.expectedOutput) {
				t.Errorf("doConfigChange() = %v, want %v", string(got), string(tt.expectedOutput))
			}
		})
	}
}
func setupViperTest(t *testing.T) string {
	folder := t.TempDir()
	config := struct {
		KeyExample string `yaml:"key_example"`
	}{}
	mars, _ := yaml.Marshal(config)
	err := os.MkdirAll(filepath.Join(folder, ".inspr"), 0755)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = os.WriteFile(filepath.Join(folder, ".inspr", "config"), mars, 0644)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	return folder
}
func Test_doListConfig(t *testing.T) {
	prepareToken(t)
	bufResp := bytes.NewBufferString("")
	cliutils.SetOutput(bufResp)
	printExistingKeys()
	outResp, _ := ioutil.ReadAll(bufResp)

	tests := []struct {
		name           string
		flagsAndArgs   []string
		expectedOutput []byte
	}{
		{
			name:           "Should list the current config",
			flagsAndArgs:   []string{},
			expectedOutput: outResp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewListConfig()
			buf := bytes.NewBufferString("")

			cliutils.SetOutput(buf)
			cmd.SetArgs(tt.flagsAndArgs)

			cmd.Execute()
			got, _ := ioutil.ReadAll(buf)

			if len(got) != len(tt.expectedOutput) {
				t.Errorf("doConfigChange() = %v, want %v", string(got), string(tt.expectedOutput))
			}

		})
	}
}
