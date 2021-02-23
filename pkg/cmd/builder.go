package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
)

// Builder is used to build cobra commands.
// it contains all the methods to manipulate a command
type Builder interface {
	WithAliases([]string) Builder
	WithDescription(description string) Builder
	WithLongDescription(long string) Builder
	WithExample(comment, command string) Builder
	WithFlagAdder(adder func(*pflag.FlagSet)) Builder
	WithFlags([]*Flag) Builder
	WithCommonFlags() Builder
	Hidden() Builder
	ExactArgs(argCount int, action func(context.Context, io.Writer, []string) error) *cobra.Command
	NoArgs(action func(context.Context, io.Writer) error) *cobra.Command
	AddSubCommand(cmds ...*cobra.Command) Builder
	Super() *cobra.Command
}

// internal builder of the package, implements all the Builder interface methods
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

// WithDescription - adds a short description to the command
func (b *builder) WithDescription(description string) Builder {
	b.cmd.Short = description
	return b
}

// WithLongDescription - adds a longer description to the command
func (b *builder) WithLongDescription(long string) Builder {
	b.cmd.Long = long
	return b
}

// WithExample - adds and example of how to use the command
//
// it will show ' #comment_given inspr  #subcommand_example '
func (b *builder) WithExample(comment, command string) Builder {
	if b.cmd.Example != "" {
		b.cmd.Example += "\n"
	}
	b.cmd.Example += fmt.Sprintf("  # %s\n inspr %s\n", comment, command)
	return b
}

// WithCommonFlags - adds to the command all the appropriate flags declared
// in the flags file of the cmd pkg
//
// In general these are flags that are commonly used by more than one cmd
func (b *builder) WithCommonFlags() Builder {
	AddFlags(&b.cmd)
	return b
}

// WithFlagAdder - allows the person to set completely different flags,
// throught the function
//
// WithFlagAdder(func(f *pflag.FlagSet) {
//
// config.AddCommonFlags(f)
//
// config.AddSetUnsetFlags(f)
//
// }).
func (b *builder) WithFlagAdder(adder func(*pflag.FlagSet)) Builder {
	adder(b.cmd.Flags())
	return b
}

// WithFlags - receives a slice of flags (defined in the cmd pkg)
// adds each of the flags in the command
func (b *builder) WithFlags(flags []*Flag) Builder {
	for _, f := range flags {
		fl := f.Flag()
		b.cmd.Flags().AddFlag(fl)
	}
	b.cmd.PreRun = func(cmd *cobra.Command, args []string) {
		ParseFlags(cmd, flags)
	}
	return b
}

// Hidden - Sets the command to be hidden, to not show in the --help description
func (b *builder) Hidden() Builder {
	b.cmd.Hidden = true
	return b
}

func (b *builder) AddSubCommand(cmds ...*cobra.Command) Builder {
	for _, cmd := range cmds {
		b.cmd.AddCommand(cmd)
	}
	return b
}

// ExactArgs - imposes the condition in which the function will only be executed
// if the exact amount of arguments are given, if didn't received the proper args
// it will show the cmd.help() content
func (b *builder) ExactArgs(argCount int, action func(context.Context, io.Writer, []string) error) *cobra.Command {
	b.cmd.Args = cobra.ExactArgs(argCount)
	b.cmd.RunE = func(_ *cobra.Command, args []string) error {
		err := handleWellKnownErrors(action(b.cmd.Context(), b.cmd.OutOrStdout(), args))
		return err
	}
	return &b.cmd
}

// NoArgs - runs the function if no args are given, in case of the user inserting
// an argument the cmd.help() content will be shown
func (b *builder) NoArgs(action func(context.Context, io.Writer) error) *cobra.Command {
	b.cmd.Args = cobra.NoArgs
	b.cmd.RunE = func(*cobra.Command, []string) error {
		err := handleWellKnownErrors(action(b.cmd.Context(), b.cmd.OutOrStdout()))
		return err
	}
	return &b.cmd
}

// Super - guarantees that a command can't be executed on it's own,
// requires the use of its subcommands.
func (b *builder) Super() *cobra.Command {
	b.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	}
	return &b.cmd
}

// handleWellKnownErrors - responsible for handling the cli common errors and
// returning them in a proper manner to the user
func handleWellKnownErrors(err error) error {
	if err == nil {
		return err
	}

	return ierrors.NewError().
		Message("error with the cli").
		Build()
}

// WithAliases adds command aliases
func (b *builder) WithAliases(alias []string) Builder {
	b.cmd.Aliases = alias
	return b
}
