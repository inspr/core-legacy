package cmd

import (
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Flag defines a INSPR CLI flag which contains a list of
// subcommands the flag belongs to in `DefinedOn` field.
type Flag struct {
	Name               string
	Shorthand          string
	Usage              string
	Value              interface{}
	DefValue           interface{}
	DefValuePerCommand map[string]interface{}
	FlagAddMethod      string
	DefinedOn          []string
	Hidden             bool

	pflag *pflag.Flag
}

// flagRegistry is a list of all Inspr CLI flags.
// When adding a new flag to the registry, please specify the
// command/commands to which the flag belongs in `DefinedOn` field.
// If the flag is a global flag, or belongs to all the subcommands,
/// specify "all"
// FlagAddMethod is method which defines a flag value with specified
// name, default value, and usage string. e.g. `StringVar`, `BoolVar`
var flagRegistry = []Flag{
	{
		Name:          "scope",
		Shorthand:     "s",
		Usage:         "inspr [command] --scope app1.app2",
		Value:         &InsprOptions.Scope,
		DefValue:      "",
		FlagAddMethod: "",
		DefinedOn:     []string{"all"},
	},
	{
		Name:          "dry-run",
		Shorthand:     "d",
		Usage:         "inspr [command] --dry-run",
		Value:         &InsprOptions.DryRun,
		DefValue:      false,
		FlagAddMethod: "BoolVar",
		DefinedOn:     []string{"apply"},
	},
	{
		Name:          "specificTag",
		Usage:         "tag used only by hidden subcommand",
		Value:         &InsprOptions.SampleFlagValue,
		DefValue:      "",
		FlagAddMethod: "",
		DefinedOn:     []string{"hidden"},
	},
}

func methodNameByType(v reflect.Value) string {
	t := v.Type().Kind()
	switch t {
	case reflect.Bool:
		return "BoolVar"
	case reflect.String:
		return "StringVar"
	case reflect.Slice:
		return "StringSliceVar"
	case reflect.Struct:
		return "Var"
	case reflect.Ptr:
		return methodNameByType(reflect.Indirect(v))
	}
	return ""
}

// Flag return a pflag.Flag from the insprCMD-flag
func (fl *Flag) Flag() *pflag.Flag {
	if fl.pflag != nil {
		return fl.pflag
	}

	methodName := fl.FlagAddMethod
	if methodName == "" {
		methodName = methodNameByType(reflect.ValueOf(fl.Value))
	}
	inputs := []interface{}{fl.Value, fl.Name}
	if methodName != "Var" {
		inputs = append(inputs, fl.DefValue)
	}
	inputs = append(inputs, fl.Usage)

	fs := pflag.NewFlagSet(fl.Name, pflag.ContinueOnError)

	reflect.ValueOf(fs).MethodByName(methodName).Call(reflectValueOf(inputs))
	f := fs.Lookup(fl.Name)
	f.Shorthand = fl.Shorthand
	f.Hidden = fl.Hidden

	fl.pflag = f
	return f
}

func reflectValueOf(values []interface{}) []reflect.Value {
	var results []reflect.Value
	for _, v := range values {
		results = append(results, reflect.ValueOf(v))
	}
	return results
}

// ParseFlags - adds flags to the cmd given
func ParseFlags(cmd *cobra.Command, flags []*Flag) {
	// Update default values.
	for _, fl := range flags {
		flag := cmd.Flag(fl.Name)
		if fl.DefValuePerCommand != nil {
			if defValue, present := fl.DefValuePerCommand[cmd.Use]; present {
				if !flag.Changed {
					flag.Value.Set(fmt.Sprintf("%v", defValue))
				}
			}
		}
	}
}

// AddFlags adds to the command the common flags that are annotated with the command name.
func AddFlags(cmd *cobra.Command) {
	var flagsForCommand []*Flag

	for i := range flagRegistry {
		fl := &flagRegistry[i]
		if !hasCmdAnnotation(cmd.Use, fl.DefinedOn) {
			continue
		}

		cmd.Flags().AddFlag(fl.Flag())

		flagsForCommand = append(flagsForCommand, fl)
	}

	// Apply command-specific default values to flags.
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		ParseFlags(cmd, flagsForCommand)
		// Since PersistentPreRunE replaces the parent's PersistentPreRunE,
		// make sure we call it, if it is set.
		if parent := cmd.Parent(); parent != nil {
			if preRun := parent.PersistentPreRunE; preRun != nil {
				if err := preRun(cmd, args); err != nil {
					return err
				}
			} else if preRun := parent.PersistentPreRun; preRun != nil {
				preRun(cmd, args)
			}
		}

		return nil
	}
}

func hasCmdAnnotation(cmdName string, annotations []string) bool {
	for _, a := range annotations {
		if cmdName == a || a == "all" {
			return true
		}
	}
	return false
}
