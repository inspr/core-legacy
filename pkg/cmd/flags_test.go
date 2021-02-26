package cmd

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestHasCmdAnnotation(t *testing.T) {
	tests := []struct {
		description string
		cmd         string
		definedOn   []string
		expected    bool
	}{
		{
			description: "flag has command annotations",
			cmd:         "build",
			definedOn:   []string{"build", "events"},
			expected:    true,
		},
		{
			description: "flag does not have command annotations",
			cmd:         "build",
			definedOn:   []string{"some"},
		},
		{
			description: "flag has all annotations",
			cmd:         "build",
			definedOn:   []string{"all"},
			expected:    true,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			hasAnnotation := hasCmdAnnotation(test.cmd, test.definedOn)

			if !reflect.DeepEqual(test.expected, hasAnnotation) {
				t.Errorf("got %v, expected %v", test.expected, hasAnnotation)
			}
		})
	}
}

func Test_methodNameByType(t *testing.T) {
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := methodNameByType(tt.args.v); got != tt.want {
				t.Errorf("methodNameByType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddFlagsSmoke(t *testing.T) {
	// Collect all commands that have common flags.
	commands := map[string]bool{}
	for _, fr := range flagRegistry {
		for _, command := range fr.DefinedOn {
			commands[command] = true
		}
	}

	// Make sure AddFlags() works for every command.
	for command := range commands {
		AddFlags(&cobra.Command{
			Use:   command,
			Short: "Test command for smoke testing",
		})
	}
}
