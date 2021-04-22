package cli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	cliutils "github.com/inspr/inspr/pkg/cmd/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestNewConfigChangeCmd(t *testing.T) {
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
	defer deleteMockViper()
	mockViper()

	bufResp := bytes.NewBufferString("")
	fmt.Fprintf(bufResp, "Success: inspr config [%v] changed to '%v'\n", "key_example", "new_value")
	outResp, _ := ioutil.ReadAll(bufResp)

	bufResp2 := bytes.NewBufferString("")
	fmt.Fprintf(bufResp2, "error: key inserted does not exist in the inspr config\n")
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

func Test_doListConfig(t *testing.T) {
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

func mockViper() {
	viper.SetConfigFile("test_config_cmd")
	viper.SetConfigType("yaml")
	viper.SetDefault("key_example", "value_example")
	viper.WriteConfig()
}

func deleteMockViper() {
	os.Remove("test_config_cmd")
}
