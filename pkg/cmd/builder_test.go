package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
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
		cmd := NewCmd("").
			WithExample("comment1", "run --flag1").
			WithExample("comment2", "run --flag2").
			NoArgs(nil)
		expected := "  # comment1\n inspr run --flag1\n\n  # comment2\n inspr run --flag2\n"

		if !reflect.DeepEqual(cmd.Example, expected) {
			t.Errorf(
				"builder.WithExample() = %v, want %v",
				cmd.Example,
				expected,
			)
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
			t.Errorf("expected error and received nil")
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

func Test_builder_MinimumArgs(t *testing.T) {
	cmd := NewCmd("").MinimumArgs(1, nil)

	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("expected error and received nil")
	}

	err = cmd.Args(cmd, []string{"valid"})
	if err != nil {
		t.Errorf("expected nil and received %v", err)
	}

	err = cmd.Args(cmd, []string{"valid", "extra"})
	if err != nil {
		t.Errorf("expected nil and received %v", err)
	}
}

func Test_builder_NewCmdError(t *testing.T) {
	cmd := NewCmd("").NoArgs(func(ctx context.Context) error {
		return errors.New("expected error")
	})

	err := cmd.RunE(nil, nil)

	if err == nil {
		t.Error("expected error received nil")
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

	if flags["test"].Usage != "usage" {
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

func Test_builder_WithCommonFlags(t *testing.T) {
	type fields struct {
		cmd cobra.Command
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "common_flags_test",
			fields: fields{
				cmd: cobra.Command{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &builder{
				cmd: tt.fields.cmd,
			}

			// commonFlags create a flagSet, previously was null
			got := b.WithCommonFlags().NoArgs(nil)
			if got.Flags() == nil {
				t.Errorf(
					"builder.WithCommonFlags() = %v, want not nil",
					got,
				)
			}
		})
	}
}

func Test_builder_WithFlagAdder(t *testing.T) {
	bufResp := bytes.NewBufferString("")
	// gets all flags names attributed to all
	expectedString := ""
	for _, f := range flagRegistry {
		for _, location := range f.DefinedOn {
			if location == "all" {
				expectedString += f.Name + "\n"
				break
			}
		}
	}
	type fields struct {
		cmd cobra.Command
	}
	type args struct {
		adder func(*pflag.FlagSet)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "flag_adder_test",
			fields: fields{
				cmd: *NewCmd("mock").WithCommonFlags().NoArgs(nil),
			},
			args: args{
				adder: func(p *pflag.FlagSet) {
					p.SortFlags = true
					p.VisitAll(func(flag *pflag.Flag) {
						fmt.Fprintln(bufResp, flag.Name)
					})
				},
			},
			want: expectedString,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bufResp = bytes.NewBufferString("")
			b := &builder{
				cmd: tt.fields.cmd,
			}
			_ = b.WithFlagAdder(tt.args.adder)

			if !reflect.DeepEqual(len(bufResp.String()), len(tt.want)) {
				t.Errorf(
					"builder.WithFlagAdder() = %v, want %v",
					bufResp,
					tt.want)
			}
		})
	}
}

func Test_builder_WithFlags(t *testing.T) {
	mockVar := 1
	mockString := "teste"
	type fields struct {
		cmd cobra.Command
	}
	type args struct {
		flags []*Flag
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "withFlags_test",
			fields: fields{
				cmd: cobra.Command{},
			},
			args: args{
				flags: []*Flag{
					{
						Name:          "mock",
						Shorthand:     "m",
						Usage:         "no usage",
						Value:         &mockVar,
						DefValue:      1,
						FlagAddMethod: "",
						DefinedOn:     []string{"all"},
					},
					{
						Name:          "mockString",
						Usage:         "no usage",
						Value:         &mockString,
						DefValue:      "teste",
						FlagAddMethod: "",
						DefinedOn:     []string{"all"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &builder{
				cmd: tt.fields.cmd,
			}

			got := b.NoArgs(nil)
			if got.Flags() == nil {
				t.Errorf(
					"builder.WithFlags(), flags shouldn't exist before calling the WithFlags()",
				)
			}

			got = b.WithFlags(tt.args.flags...).NoArgs(nil)
			if got.Flags() == nil {
				t.Errorf(
					"builder.WithFlags() = %v, want not nil",
					got.Flags(),
				)
			}
		})
	}
}

func Test_builder_Hidden(t *testing.T) {
	type fields struct {
		cmd cobra.Command
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "test_hidden_func",
			fields: fields{
				cmd: cobra.Command{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &builder{
				cmd: tt.fields.cmd,
			}
			got := b.Hidden().NoArgs(nil)
			if !reflect.DeepEqual(got.Hidden, tt.want) {
				t.Errorf(
					"builder.Hidden() = %v, want %v",
					got.Hidden,
					tt.want,
				)
			}
		})
	}
}

func Test_builder_AddSubCommand(t *testing.T) {
	type fields struct {
		cmd cobra.Command
	}
	type args struct {
		cmds []*cobra.Command
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "test_addSubCmd",
			fields: fields{cmd: cobra.Command{}},
			args: args{
				cmds: []*cobra.Command{
					{
						Use: "cmd_1",
					},
					{
						Use: "cmd_2",
					},
					{
						Use: "cmd_3",
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &builder{
				cmd: tt.fields.cmd,
			}
			got := b.AddSubCommand(tt.args.cmds...).NoArgs(nil)

			if !reflect.DeepEqual(got.HasSubCommands(), tt.want) {
				t.Errorf(
					"builder.AddSubCommand() = %v, want %v",
					got.HasSubCommands(),
					tt.want,
				)
			}
		})
	}
}

func Test_builder_Super(t *testing.T) {
	type fields struct {
		cmd cobra.Command
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test_super_func",
			fields: fields{
				cmd: cobra.Command{
					Use: "mock_usage",
				},
			},
			want: "Usage:\n  mock_usage\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bufResp := bytes.NewBufferString("")
			tt.fields.cmd.SetOutput(bufResp)

			b := &builder{
				cmd: tt.fields.cmd,
			}

			got := b.Super()
			_ = got.RunE(got, []string{})
			if !reflect.DeepEqual(bufResp.String(), tt.want) {
				t.Errorf(
					"builder.Super() = %v, want %v",
					bufResp,
					tt.want,
				)
			}
		})
	}
}

func Test_handleWellKnownErrors(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "handling_error",
			args: args{
				err: errors.New("new error"),
			},
			wantErr: true,
		},
		{
			name: "handling_no_error",
			args: args{
				err: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := handleWellKnownErrors(tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("handleWellKnownErrors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_builder_WithAliases(t *testing.T) {
	type fields struct {
		cmd cobra.Command
	}
	type args struct {
		alias []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "testing_alias_creation",
			fields: fields{
				cmd: cobra.Command{},
			},
			args: args{alias: []string{"alias_1", "alias_2"}},
		},
		{
			name: "testing_no_alias_given",
			fields: fields{
				cmd: cobra.Command{},
			},
			args: args{alias: []string{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &builder{
				cmd: tt.fields.cmd,
			}
			got := b.WithAliases(tt.args.alias).NoArgs(nil)

			if !reflect.DeepEqual(got.Aliases, tt.args.alias) {
				t.Errorf(
					"builder.WithAliases() = %v, want %v",
					got.Aliases,
					tt.args.alias,
				)
			}
		})
	}
}
