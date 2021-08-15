package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Option is an option to be applied on the construction of a command
type Option func(*cobra.Command)

// Builder is used to build cobra commands.
// it contains all the methods to manipulate a command
type Builder interface {
	WithAliases(...string) Builder
	WithDescription(description string) Builder
	WithLongDescription(long string) Builder
	WithExample(comment, command string) Builder
	WithFlagAdder(adder func(*pflag.FlagSet)) Builder
	WithFlags(...*Flag) Builder
	WithCommonFlags() Builder
	Hidden() Builder
	ExactArgs(
		argCount int,
		action func(context.Context, []string) error,
	) *cobra.Command
	MinimumArgs(
		argCount int,
		action func(context.Context, []string) error,
	) *cobra.Command
	NoArgs(action func(context.Context) error) *cobra.Command
	AddSubCommand(cmds ...*cobra.Command) Builder
	Version(version string) Builder
	WithOptions(...Option) Builder
	ValidArgsFunc(
		validation func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective),
	) Builder
	WithRequiredFlag(string) Builder
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

// WithRequiredFlag sets a flag as required when creating the command
func (b *builder) WithRequiredFlag(flag string) Builder {
	b.cmd.MarkFlagRequired(flag)
	return b
}

// ValidArgsFunc adds a validation function to the arguments of a command. This is useful for completion
func (b *builder) ValidArgsFunc(
	validation func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective),
) Builder {
	b.cmd.ValidArgsFunction = validation
	return b
}

// WithOptions adds custom options to the command. These options are functions that
// cause some change in the command
func (b *builder) WithOptions(options ...Option) Builder {
	for _, opt := range options {
		if opt != nil {
			opt(&b.cmd)
		}
	}
	return b
}

// Version adds the version to the cli
func (b *builder) Version(v string) Builder {
	b.cmd.Version = v
	return b
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
// it will show ' #comment_given insprctl  #subcommand_example '
func (b *builder) WithExample(comment, command string) Builder {
	if b.cmd.Example != "" {
		b.cmd.Example += "\n"
	}
	b.cmd.Example += fmt.Sprintf("  # %s\n insprctl %s\n", comment, command)
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
// through the function
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
func (b *builder) WithFlags(flags ...*Flag) Builder {
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

// AddSubCommand adds subcommands to the command
func (b *builder) AddSubCommand(cmds ...*cobra.Command) Builder {
	for _, cmd := range cmds {
		b.cmd.AddCommand(cmd)
	}
	return b
}

// ExactArgs - imposes the condition in which the function will only be executed
// if the exact amount of arguments are given, if didn't received the proper args
// it will show the cmd.help() content
func (b *builder) ExactArgs(
	argCount int,
	action func(context.Context, []string) error,
) *cobra.Command {
	f := b.cmd.ValidArgsFunction
	b.cmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) >= argCount {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return f(cmd, args, toComplete)
	}
	b.cmd.Args = cobra.ExactArgs(argCount)
	b.cmd.RunE = func(_ *cobra.Command, args []string) error {
		return (action(b.cmd.Context(), args))
	}
	return &b.cmd
}

// MinimumArgs - imposes the condition in which the function will only be executed
// if the minimum amount of arguments are given, if didn't received the proper args
// it will show the cmd.help() content
func (b *builder) MinimumArgs(
	minArgs int,
	action func(context.Context, []string) error,
) *cobra.Command {
	b.cmd.Args = cobra.MinimumNArgs(minArgs)
	b.cmd.RunE = func(_ *cobra.Command, args []string) error {
		return (action(b.cmd.Context(), args))
	}
	return &b.cmd
}

// NoArgs - runs the function if no args are given, in case of the user inserting
// an argument the cmd.help() content will be shown
func (b *builder) NoArgs(action func(context.Context) error) *cobra.Command {
	b.cmd.Args = cobra.NoArgs
	b.cmd.RunE = func(*cobra.Command, []string) error {
		return (action(b.cmd.Context()))
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

// WithAliases adds command aliases
func (b *builder) WithAliases(alias ...string) Builder {
	b.cmd.Aliases = alias
	return b
}
