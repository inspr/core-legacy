package cli

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
	"gopkg.in/yaml.v2"
)

// NewApplyCmd - mock subcommand
func NewApplyCmd() *cobra.Command {
	applyCmd := cmd.NewCmd("apply").
		WithDescription("Applies changes to the connected cluster").
		WithLongDescription("Apply takes a file or a directory and applies the defined components on the connected cluster.\n"+
			"It can be called with the flag --update for updating instead of creating a new dApp.\n"+
			"It can be called with the flag --dry-run so the changes that would be made are shown, but not applied on the cluster").
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
				Usage:         "inspr apply [path] --update",
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
	fmt.Println(cmd.InsprOptions.DryRun)
	fmt.Println(cmd.InsprOptions.Update)
	var files []string
	var path string
	hasFileFlag := (cmd.InsprOptions.AppliedFileStructure != "")
	hasFolderFlag := (cmd.InsprOptions.AppliedFolderStructure != "")
	if hasFileFlag == hasFolderFlag {
		fmt.Fprint(out, "Specified file/folder path is invalid\n")
		return ierrors.NewError().Message("invalid flag arguments").Build()
	}

	if hasFileFlag {
		filePath := strings.Split(cmd.InsprOptions.AppliedFileStructure, "/")
		if len(filePath) == 1 {
			files = append(files, filePath[0])
		} else {
			path = strings.Join(filePath[:len(filePath)-1], "/") + "/"
			files = append(files, filePath[len(filePath)-1])
		}
	} else {
		path = cmd.InsprOptions.AppliedFolderStructure
		err := getFilesFromFolder(cmd.InsprOptions.AppliedFolderStructure, &files)
		if err != nil {
			fmt.Fprint(out, err.Error())
			return err
		}
	}

	appliedFiles := applyValidFiles(path, files)

	if len(appliedFiles) > 0 {
		printAppliedFiles(appliedFiles, out)
	} else {
		fmt.Fprint(out, "No files were applied\n")
	}

	return nil
}

func isYaml(file string) bool {
	tempStr := strings.Split(file, ".")
	return tempStr[len(tempStr)-1] == "yaml" || tempStr[len(tempStr)-1] == "yml"
}

func printAppliedFiles(appliedFiles []applied, out io.Writer) {
	fmt.Fprint(out, "Applying: \n")
	for _, file := range appliedFiles {
		fmt.Fprint(out, file.file+" | "+file.component.Kind+" | "+file.component.APIVersion+"\n")
	}
}

func getFilesFromFolder(path string, files *[]string) error {
	folder, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range folder {
		*files = append(*files, file.Name())
	}
	return nil
}

func applyValidFiles(path string, files []string) []applied {
	var appliedFiles []applied

	for _, file := range files {
		if isYaml(file) {
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
			err = apply(f)
			if err != nil {
				continue
			}
			appliedFiles = append(appliedFiles, applied{file: file, component: comp})

		}
	}

	return appliedFiles
}
