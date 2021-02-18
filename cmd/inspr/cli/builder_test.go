package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestNewCmd(t *testing.T) {
	type args struct {
		use string
	}
	tests := []struct {
		name string
		args args
		want Builder
	}{
		{
			name: "creation_test",
			args: args{use: "mock_use"},
			want: &builder{
				cmd: cobra.Command{
					Use: "mock_use",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCmd(tt.args.use); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builder_WithDescription(t *testing.T) {
	type args struct {
		description string
	}
	tests := []struct {
		name string
		b    *builder
		args args
		want Builder
	}{
		{
			name: "cobra_cmd_description",
			b:    &builder{cmd: cobra.Command{}},
			args: args{description: "mock_description"},
			want: &builder{cmd: cobra.Command{
				Short: "mock_description",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.WithDescription(tt.args.description); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builder.WithDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builder_WithLongDescription(t *testing.T) {
	type args struct {
		long string
	}
	tests := []struct {
		name string
		b    *builder
		args args
		want Builder
	}{
		{
			name: "cobra_cmd_long_description",
			b:    &builder{cmd: cobra.Command{}},
			args: args{long: "mock_long_description"},
			want: &builder{cmd: cobra.Command{
				Long: "mock_long_description",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.WithLongDescription(tt.args.long); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builder.WithLongDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builder_WithExample(t *testing.T) {
	type args struct {
		comment string
		command string
	}
	tests := []struct {
		name string
		b    *builder
		args args
		want Builder
	}{
		{
			name: "cobra_example",
			b:    &builder{cmd: cobra.Command{}},
			args: args{
				comment: "mock_comment",
				command: "mock_command",
			},
			want: &builder{cmd: cobra.Command{
				Example: fmt.Sprintf("  # %s\n inspr %s\n", "mock_comment", "mock_command"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.WithExample(tt.args.comment, tt.args.command); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builder.WithExample() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builder_WithExamples(t *testing.T) {
	t.Run("build_with_multiple_examples", func(t *testing.T) {
		cmd := NewCmd("").WithExample("comment1", "run --flag1").WithExample("comment2", "run --flag2").NoArgs(nil)
		expected := "  # comment1\n inspr run --flag1\n\n  # comment2\n inspr run --flag2\n"

		if !reflect.DeepEqual(cmd.Example, expected) {
			t.Errorf("builder.WithExample() = %v, want %v", cmd.Example, expected)
		}
	})
}

func Test_builder_NoArgs(t *testing.T) {
	t.Run("testing NoArgs", func(t *testing.T) {
		cmd := NewCmd("").NoArgs(nil)

		err := cmd.Args(cmd, []string{})
		if err != nil {
			t.Errorf("expected nil and received %v", err)
		}
		err = cmd.Args(cmd, []string{"extract arg"})
		if err == nil {
			t.Errorf("expected error and receiver nil")
		}

	})
}

func Test_builder_ExactArgs(t *testing.T) {
	cmd := NewCmd("").ExactArgs(1, nil)

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error and received nil")
	}

	err = cmd.Args(cmd, []string{"valid"})
	if err != nil {
		t.Errorf("expected nil and received %v", err)
	}

	err = cmd.Args(cmd, []string{"valid", "extra"})
	if err == nil {
		t.Error("expected error and received nil")
	}
}

func Test_builder_NewCmdError(t *testing.T) {
	cmd := NewCmd("").NoArgs(func(ctx context.Context, out io.Writer) error {
		return errors.New("expected error")
	})

	err := cmd.RunE(nil, nil)

	if err == nil {
		t.Error("expected error received nil")
	}
}

func Test_builder_CmdOutput(t *testing.T) {
	var buf bytes.Buffer
	cmd := NewCmd("").ExactArgs(1, func(ctx context.Context, out io.Writer, args []string) error {
		fmt.Fprintf(out, "test output: %v\n", args)
		return nil
	})
	cmd.SetOutput(&buf)

	err := cmd.RunE(nil, []string{"arg1"})
	if err != nil {
		t.Errorf("expected nil and received %v\n", err)
	}

	expected := "test output: [arg1]\n"
	if !reflect.DeepEqual(buf.String(), expected) {
		t.Errorf("expected %v, received %v", expected, buf.String())
	}
}

func Test_builder_NewCmdWithFlags(t *testing.T) {
	cmd := NewCmd("").WithFlagAdder(func(flagSet *pflag.FlagSet) {
		flagSet.Bool("test", false, "usage")
	}).NoArgs(nil)

	flags := listFlags(cmd.Flags())

	if len(flags) != 1 {
		t.Errorf("expected flags to be of length 1, found %v", len(flags))
	}

	if "usage" != flags["test"].Usage {
		t.Errorf("expected 'usage', got %v", flags["test"].Usage)
	}
}

func Test_builder_CmdHidden(t *testing.T) {
	cmd := NewCmd("").NoArgs(nil)
	if cmd.Hidden != false {
		t.Error("Expected cmd.Hidden to be false")
	}
	cmd = NewCmd("").Hidden().NoArgs(nil)
	if cmd.Hidden != true {
		t.Error("Expected cmd.Hidden to be true")
	}
}

// listFlags - returns a map of flags
func listFlags(set *pflag.FlagSet) map[string]*pflag.Flag {
	flagsByName := make(map[string]*pflag.Flag)

	set.VisitAll(func(f *pflag.Flag) {
		flagsByName[f.Name] = f
	})

	return flagsByName
}
