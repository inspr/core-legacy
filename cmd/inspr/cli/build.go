package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	iErrors "gitlab.inspr.dev/inspr/core/pkg/ierrors"
)

// Builder is used to build cobra commands.
type Builder interface {
	WithDescription(description string) Builder
	WithLongDescription(long string) Builder
	WithExample(comment, command string) Builder
	WithFlagAdder(adder func(*pflag.FlagSet)) Builder
	WithFlags([]*Flag) Builder
	WithCommonFlags() Builder
	Hidden() Builder
	ExactArgs(argCount int, action func(context.Context, io.Writer, []string) error) *cobra.Command
	NoArgs(action func(context.Context, io.Writer) error) *cobra.Command
}

type builder struct {
	cmd cobra.Command
}

// NewCmd creates a new command builder.
func NewCmd(use string) Builder {
	return &builder{
		cmd: cobra.Command{
			Use: use,
		},
	}
}

func (b *builder) WithDescription(description string) Builder {
	b.cmd.Short = description
	return b
}

func (b *builder) WithLongDescription(long string) Builder {
	b.cmd.Long = long
	return b
}

func (b *builder) WithExample(comment, command string) Builder {
	if b.cmd.Example != "" {
		b.cmd.Example += "\n"
	}
	b.cmd.Example += fmt.Sprintf("  # %s\n  inspr %s\n", comment, command)
	return b
}

func (b *builder) WithCommonFlags() Builder {
	AddFlags(&b.cmd)
	return b
}

func (b *builder) WithFlagAdder(adder func(*pflag.FlagSet)) Builder {
	adder(b.cmd.Flags())
	return b
}

func (b *builder) WithFlags(flags []*Flag) Builder {
	for _, f := range flags {
		fl := f.flag()
		b.cmd.Flags().AddFlag(fl)
	}
	b.cmd.PreRun = func(cmd *cobra.Command, args []string) {
		ParseFlags(cmd, flags)
	}

	return b
}

func (b *builder) Hidden() Builder {
	b.cmd.Hidden = true
	return b
}
func (b *builder) ExactArgs(argCount int, action func(context.Context, io.Writer, []string) error) *cobra.Command {
	b.cmd.Args = cobra.ExactArgs(argCount)
	b.cmd.RunE = func(_ *cobra.Command, args []string) error {
		err := handleWellKnownErrors(action(b.cmd.Context(), b.cmd.OutOrStdout(), args))
		return err
	}
	return &b.cmd
}

func (b *builder) NoArgs(action func(context.Context, io.Writer) error) *cobra.Command {
	b.cmd.Args = cobra.NoArgs
	b.cmd.RunE = func(*cobra.Command, []string) error {
		err := handleWellKnownErrors(action(b.cmd.Context(), b.cmd.OutOrStdout()))
		return err
	}
	return &b.cmd
}

func handleWellKnownErrors(err error) error {
	if err == nil {
		return err
	}
	// TODO error handler of the cli in the ierrors pkg
	return iErrors.NewError().Message("error with the cli").Build()
}
