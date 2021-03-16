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
	holder := 1
	holderPointer := &holder
	type args struct {
		v reflect.Value
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "type_int",
			args: args{v: reflect.ValueOf(1)},
			want: "IntVar",
		},
		{
			name: "type_bool",
			args: args{v: reflect.ValueOf(true)},
			want: "BoolVar",
		},
		{
			name: "type_string",
			args: args{v: reflect.ValueOf("mock")},
			want: "StringVar",
		},
		{
			name: "type_string_slice",
			args: args{v: reflect.ValueOf([]string{"a", "b"})},
			want: "StringSliceVar",
		},
		{
			name: "type_struct",
			args: args{
				v: reflect.ValueOf(struct{ x int }{x: 10}),
			},
			want: "Var",
		},
		{
			name: "type_pointer",
			args: args{
				v: reflect.ValueOf(holderPointer),
			},
			want: "IntVar",
		},
		{
			name: "type_double",
			args: args{
				v: reflect.ValueOf(2.0),
			},
			want: "",
		},
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
