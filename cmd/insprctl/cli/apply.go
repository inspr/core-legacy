package cli

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"

	"inspr.dev/inspr/pkg/cmd"
	cliutils "inspr.dev/inspr/pkg/cmd/utils"
	"inspr.dev/inspr/pkg/ierrors"

	"gopkg.in/yaml.v2"
	"inspr.dev/inspr/pkg/meta"
)

// NewApplyCmd - mock subcommand
func NewApplyCmd() *cobra.Command {
	applyCmd := cmd.NewCmd("apply").
		WithDescription("Applies changes to the connected cluster").
		WithLongDescription(`
Apply takes a file or a directory and applies the defined components on the connected cluster.

It can be called with the flag --update for updating instead of creating a new dApp.

It can be called with the flag --dry-run so the changes that would be made are shown, but not applied on the cluster
		`).
		WithExample("Applies a structure component defined in a file", "apply -f app.yaml").
		WithExample("Applies components defined in a specific folder", "apply -k randfolder/").
		WithExample("Applies a structure component defined in a specific scope", "apply -f app.yaml --scope app1.app2").
		WithFlags([]*cmd.Flag{
			{
				Name:          "file",
				Shorthand:     "f",
				Usage:         "insprctl apply -f type.yaml",
				Value:         &cmd.InsprOptions.AppliedFileStructure,
				DefValue:      "",
				FlagAddMethod: "",
				DefinedOn:     []string{"apply"},
			},
			{
				Name:          "folder",
				Shorthand:     "k",
				Usage:         "insprctl apply -k randfolder/",
				Value:         &cmd.InsprOptions.AppliedFolderStructure,
				DefValue:      "",
				FlagAddMethod: "",
				DefinedOn:     []string{"apply"},
			},
			{
				Name:          "update",
				Shorthand:     "u",
				Usage:         "insprctl apply (-f FILENAME | -k DIRECTORY) --update",
				Value:         &cmd.InsprOptions.Update,
				DefValue:      false,
				FlagAddMethod: "BoolVar",
				DefinedOn:     []string{"apply"},
			},
		}...).
		WithCommonFlags().
		WithOptions(cliutils.AddDefaultFlagCompletion()).
		NoArgs(doApply)

	applyCmd.MarkFlagFilename("file", "yaml", "yml")
	applyCmd.MarkFlagDirname("folder")

	return applyCmd
}

type applied struct {
	fileName  string
	component meta.Component
	content   []byte
}

func doApply(_ context.Context) error {
	var files []string
	var path string
	var err error
	out := cliutils.GetCliOutput()
	hasFileFlag := (cmd.InsprOptions.AppliedFileStructure != "")
	hasFolderFlag := (cmd.InsprOptions.AppliedFolderStructure != "")
	if hasFileFlag == hasFolderFlag {
		fmt.Fprintln(
			out,
			"Invalid command call\nFor help, type 'insprctl apply --help'",
		)
		return ierrors.New("invalid flag arguments")
	}

	if hasFileFlag {
		files = append(files, cmd.InsprOptions.AppliedFileStructure)
	} else {
		path = cmd.InsprOptions.AppliedFolderStructure
		files, err = getFilesFromFolder(cmd.InsprOptions.AppliedFolderStructure)
		if err != nil {
			fmt.Fprint(out, err.Error())
			return err
		}
	}

	appliedFiles := applyValidFiles(path, files, out)

	if len(appliedFiles) > 0 {
		printAppliedFiles(appliedFiles, out)
	} else {
		fmt.Fprint(out, "No files were applied\nFiles to be applied must be .yaml or .yml\n")
	}

	return nil
}

func isYaml(file string) bool {
	tempStr := filepath.Ext(file)
	return tempStr == ".yaml" || tempStr == ".yml"
}

func printAppliedFiles(appliedFiles []applied, out io.Writer) {
	fmt.Fprint(out, "\nApplied:\n")
	for _, file := range appliedFiles {
		fmt.Fprint(out, file.fileName+" | "+file.component.Kind+" | "+file.component.APIVersion+"\n")
	}
}

func getFilesFromFolder(path string) ([]string, error) {
	var files []string
	folder, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range folder {
		files = append(files, file.Name())
	}

	return files, nil
}

func applyValidFiles(path string, files []string, out io.Writer) []applied {
	var appliedFiles []applied

	filesToApply := getOrderedFiles(path, files)

	for _, file := range filesToApply {

		apply, err := GetFactory().GetRunMethod(file.component)
		if err != nil {
			continue
		}

		err = apply(file.content, out)
		if err != nil {
			ierrors.Wrap(err, file.fileName)
			cliutils.RequestErrorMessage(err, out)
			continue
		}

		appliedFiles = append(appliedFiles, file)

	}
	return appliedFiles
}

func getOrderedFiles(path string, files []string) []applied {
	var apps, channels, types, aliases []applied
	for _, file := range files {
		if isYaml(file) {
			comp := meta.Component{}

			f, err := ioutil.ReadFile(filepath.Join(path, file))
			if err != nil {
				continue
			}

			err = yaml.Unmarshal(f, &comp)
			if err != nil || comp.APIVersion == "" || comp.Kind == "" {
				continue
			}

			if comp.Kind == "dapp" {
				apps = append(apps, applied{component: comp, fileName: file, content: f})
			} else if comp.Kind == "channel" {
				channels = append(channels, applied{component: comp, fileName: file, content: f})
			} else if comp.Kind == "type" {
				types = append(types, applied{component: comp, fileName: file, content: f})
			} else if comp.Kind == "alias" {
				aliases = append(aliases, applied{component: comp, fileName: file, content: f})
			}
		}
	}
	ordered := append(apps, types...)
	ordered = append(ordered, channels...)
	ordered = append(ordered, aliases...)
	return ordered
}
