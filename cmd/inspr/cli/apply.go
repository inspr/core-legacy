// TODO -> Change to Walk, and add yaml unmarshall before doSomething

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
		WithDescription("applies changes to the connected cluster").
		WithLongDescription("apply takes a file or a directory and applies the defined components on the connected cluster").
		WithExample("Applies a structure component defined in a file", "apply -f app.yaml").
		WithExample("Applies components defined in a specific folder", "apply -k randfolder/").
		WithExample("Applies a structure component defined in a specific scope", "apply -f app.yaml --scope app1.app2").
		WithCommonFlags().
		WithFlags([]*cmd.Flag{
			{
				Name:          "file",
				Usage:         "inspr apply -f ctype.yaml",
				Shorthand:     "f",
				Value:         &cmd.InsprOptions.AppliedFileStructure,
				DefValue:      "",
				FlagAddMethod: "",
				DefinedOn:     []string{"apply"},
			},
			{
				Name:          "folder",
				Usage:         "inspr apply -k randfolder/",
				Shorthand:     "k",
				Value:         &cmd.InsprOptions.AppliedFolderStructure,
				DefValue:      "",
				FlagAddMethod: "",
				DefinedOn:     []string{"apply"},
			},
		}).
		NoArgs(doApply)

	applyCmd.MarkFlagFilename("file", "yaml", "yml")
	applyCmd.MarkFlagDirname("folder")

	return applyCmd
}

func doApply(_ context.Context, out io.Writer) error {
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
			fmt.Println(err.Error())
			return err
		}
	}

	ignoredFiles := applyValidFiles(path, files)

	if len(ignoredFiles) > 0 {
		printIgnoredFiles(ignoredFiles)
	}

	return nil
}

func isYaml(file string) bool {
	tempStr := strings.Split(file, ".")
	return tempStr[len(tempStr)-1] == "yaml" || tempStr[len(tempStr)-1] == "yml"
}

func printIgnoredFiles(ignoredFiles []string) {
	fmt.Println("The following files were ignored: ")
	for _, file := range ignoredFiles {
		fmt.Println(file)
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

func applyValidFiles(path string, files []string) []string {
	var ignoredFiles []string

	for _, file := range files {
		if isYaml(file) {
			comp := meta.Component{}
			f, err := ioutil.ReadFile(path + file)
			if err != nil {
				ignoredFiles = append(ignoredFiles, file)
				continue
			}
			err = yaml.Unmarshal(f, &comp)
			if err != nil || comp.APIVersion == "" || comp.Kind == "" {
				ignoredFiles = append(ignoredFiles, file)
				continue
			}
			funcs[comp](f)
		} else {
			ignoredFiles = append(ignoredFiles, file)
		}
	}

	return ignoredFiles
}

var funcs = map[meta.Component]func([]byte){
	{APIVersion: "v1", Kind: "channel"}: func(s []byte) {
		ch := meta.Channel{}

		yaml.Unmarshal(s, &ch)
		fmt.Println(ch)
	},
	{APIVersion: "v1", Kind: "app"}: func(s []byte) {
		ch := meta.App{}

		yaml.Unmarshal(s, &ch)
		fmt.Println(ch)
	},
}
