package cli

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gopkg.in/yaml.v2"
)

var applyCmd *cobra.Command

// NewApplyCmd - mock subcommand
func NewApplyCmd() *cobra.Command {
	applyCmd = cmd.NewCmd("apply").
		WithDescription("Applies changes to the connected cluster").
		WithLongDescription(`
Apply takes a file or a directory and applies the defined components on the connected cluster.

It can be called with the flag --update for updating instead of creating a new dApp.

It can be called with the flag --dry-run so the changes that would be made are shown, but not applied on the cluster
		`).
		WithExample("Applies a structure component defined in a file", "apply -f app.yaml").
		WithExample("Applies components defined in a specific folder", "apply -k randfolder/").
		WithExample("Applies a structure component defined in a specific scope", "apply -f app.yaml --scope app1.app2").
		WithCommonFlags().
		WithFlags([]*cmd.Flag{
			{
				Name:          "file",
				Shorthand:     "f",
				Usage:         "inspr apply -f ctype.yaml",
				Value:         &cmd.InsprOptions.AppliedFileStructure,
				DefValue:      "",
				FlagAddMethod: "",
				DefinedOn:     []string{"apply"},
			},
			{
				Name:          "folder",
				Shorthand:     "k",
				Usage:         "inspr apply -k randfolder/",
				Value:         &cmd.InsprOptions.AppliedFolderStructure,
				DefValue:      "",
				FlagAddMethod: "",
				DefinedOn:     []string{"apply"},
			},
			{
				Name:          "update",
				Shorthand:     "u",
				Usage:         "inspr apply (-f FILENAME | -k DIRECTORY) --update",
				Value:         &cmd.InsprOptions.Update,
				DefValue:      false,
				FlagAddMethod: "BoolVar",
				DefinedOn:     []string{"apply"},
			},
		}).
		NoArgs(doApply)

	applyCmd.MarkFlagFilename("file", "yaml", "yml")
	applyCmd.MarkFlagDirname("folder")

	return applyCmd
}

type applied struct {
	file      string
	component meta.Component
}

func doApply(_ context.Context, out io.Writer) error {
	var files []string
	var path string
	var err error
	hasFileFlag := (cmd.InsprOptions.AppliedFileStructure != "")
	hasFolderFlag := (cmd.InsprOptions.AppliedFolderStructure != "")
	if hasFileFlag == hasFolderFlag {
		fmt.Fprintln(out, "\nGiven flags are invalid")
		fmt.Fprintln(out, "\nCommand Help")

		applyCmd.Help()
		return ierrors.NewError().Message("invalid flag arguments").Build()
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
		fmt.Fprint(out, "No files were applied\n")
	}

	return nil
}

func isYaml(file string) bool {
	tempStr := filepath.Ext(file)
	return tempStr == ".yaml" || tempStr == ".yml"
}

func printAppliedFiles(appliedFiles []applied, out io.Writer) {
	fmt.Fprint(out, "Applying: \n")
	for _, file := range appliedFiles {
		fmt.Fprint(out, file.file+" | "+file.component.Kind+" | "+file.component.APIVersion+"\n")
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

	for _, file := range files {
		if isYaml(file) {
			fmt.Println(file)
			comp := meta.Component{}
			f, err := ioutil.ReadFile(path + file)
			if err != nil {
				continue
			}
			err = yaml.Unmarshal(f, &comp)
			if err != nil || comp.APIVersion == "" || comp.Kind == "" {
				continue
			}

			apply, err := GetFactory().GetRunMethod(comp)
			if err != nil {
				continue
			}
			err = apply(f, out)
			if err != nil {
				fmt.Fprintf(out, "error while applying file '%v' :\n %v\n", file, err.Error())
				continue
			}
			appliedFiles = append(appliedFiles, applied{file: file, component: comp})

		}
	}

	return appliedFiles
}
