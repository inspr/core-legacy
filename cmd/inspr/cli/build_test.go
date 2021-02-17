package cli

import (
	"context"
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
				Example: "\n" + fmt.Sprintf("  # %s\n  inspr %s\n", "mock_comment", "mock_command"),
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

// TODO remove?
func Test_builder_WithCommonFlags(t *testing.T) {
	tests := []struct {
		name string
		b    *builder
		want Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.WithCommonFlags(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builder.WithCommonFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builder_WithFlagAdder(t *testing.T) {
	type args struct {
		adder func(*pflag.FlagSet)
	}
	tests := []struct {
		name string
		b    *builder
		args args
		want Builder
	}{
		{
			name: "cobra_flagAdder",
			b:    &builder{cmd: cobra.Command{}},
			args: args{
				adder: func(fl *pflag.FlagSet) {

				},
			},
			want: &builder{cmd: cobra.Command{
				Long: "mock_long_description",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.WithFlagAdder(tt.args.adder); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builder.WithFlagAdder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builder_WithFlags(t *testing.T) {
	type args struct {
		flags []*Flag
	}
	tests := []struct {
		name string
		b    *builder
		args args
		want Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.WithFlags(tt.args.flags); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builder.WithFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builder_Hidden(t *testing.T) {
	tests := []struct {
		name string
		b    *builder
		want Builder
	}{
		{
			name: "cobra_hidden",
			b:    &builder{cmd: cobra.Command{}},
			want: &builder{cmd: cobra.Command{
				Hidden: true,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Hidden(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builder.Hidden() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builder_ExactArgs(t *testing.T) {
	type args struct {
		argCount int
		action   func(context.Context, io.Writer, []string) error
	}
	tests := []struct {
		name string
		b    *builder
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.ExactArgs(tt.args.argCount, tt.args.action); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builder.ExactArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_builder_NoArgs(t *testing.T) {
	type args struct {
		action func(context.Context, io.Writer) error
	}
	tests := []struct {
		name string
		b    *builder
		args args
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.NoArgs(tt.args.action); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("builder.NoArgs() = %v, want %v", got, tt.want)
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := handleWellKnownErrors(tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("handleWellKnownErrors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
