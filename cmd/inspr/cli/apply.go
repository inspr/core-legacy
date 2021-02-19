package cli

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.inspr.dev/inspr/core/pkg/cmd"
	"gitlab.inspr.dev/inspr/core/pkg/ierrors"
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
	var path string
	hasFileFlag := (cmd.InsprOptions.AppliedFileStructure != "")
	hasFolderFlag := (cmd.InsprOptions.AppliedFolderStructure != "")
	if hasFileFlag == hasFolderFlag {
		fmt.Fprint(out, "Specified file/folder path is invalid\n")
		return ierrors.NewError().Message("invalid flag arguments").Build()
	}

	specificFile := ""
	if hasFileFlag {
		filePath := strings.Split(cmd.InsprOptions.AppliedFileStructure, "/")
		if len(filePath) == 1 {
			path = ""
			specificFile = filePath[0]
		} else {
			specificFile = filePath[len(filePath)-1]
			filePath = filePath[1 : len(filePath)-1]
			path = strings.Join(filePath, "/")
		}
	} else {
		filePath := strings.Split(cmd.InsprOptions.AppliedFolderStructure, "/")
		filePath = filePath[1:]
		path = strings.Join(filePath, "/")
	}

	folder, err := getFilesFromFolder("/" + path)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	validFiles, ignoredFiles := getValidFiles(folder, path, specificFile)

	doSomething(validFiles)

	if len(ignoredFiles) > 0 {
		printIgnoredFiles(ignoredFiles)
	}

	return nil
}

func getFile(name string) (io.Reader, error) {
	if name == "-" {
		return os.Stdin, nil
	}
	return os.Open(name)
}

func getValidFiles(folder []os.FileInfo, path, specificFile string) ([]io.Reader, []string) {
	var validFiles []io.Reader
	var ignoredFiles []string

	if specificFile == "" {
		for _, file := range folder {
			if isYaml(file, "") {
				validFile, err := getFile(path + "/" + file.Name())
				if err != nil {
					fmt.Println("Sabia que ia dar erro1")
					fmt.Println(err)
					ignoredFiles = append(ignoredFiles, file.Name())
				} else {
					validFiles = append(validFiles, validFile)
				}
			} else {
				ignoredFiles = append(ignoredFiles, file.Name())
			}
		}
	} else {
		if isYaml(nil, specificFile) {
			validFile, err := getFile(path + "/" + specificFile)
			if err != nil {
				fmt.Println(err)
				return nil, nil
			}
			return []io.Reader{validFile}, nil
		}

		fmt.Println("given file should be .yaml or .yml")
		return nil, nil
	}
	return validFiles, ignoredFiles
}

func isYaml(file os.FileInfo, fileStr string) bool {
	if fileStr != "" {
		tempStr := strings.Split(fileStr, ".")
		return tempStr[len(tempStr)-1] == "yaml" || tempStr[len(tempStr)-1] == "yml"
	}
	return filepath.Ext(file.Name()) == ".yaml" || filepath.Ext(file.Name()) == ".yml"
}

func printIgnoredFiles(ignoredFiles []string) {
	fmt.Println("The following files were ignored: ")
	for _, file := range ignoredFiles {
		fmt.Println(file)
	}
}

func getFilesFromFolder(path string) ([]os.FileInfo, error) {
	if path == "/" {
		path = ""
	}
	dir, err := os.Getwd()
	if err != nil {
		return nil, ierrors.NewError().InnerError(err).Message("couldn't find current directory: " + dir).Build()
	}

	folder, err := ioutil.ReadDir(dir + path)
	if err != nil {
		return nil, ierrors.NewError().InnerError(err).Message("couldn't open folder: " + dir + path).Build()
	}

	return folder, nil
}

func doSomething(files []io.Reader) error {
	if files == nil {
		fmt.Println("No changes were applied")
		return nil
	}
	for _, file := range files {
		fmt.Println(file)
	}
	return nil
}
