package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/inspr/inspr/pkg/cmd"
	"github.com/inspr/inspr/pkg/ierrors"
	"gopkg.in/yaml.v2"
)

type insprConfiguration struct {
	ServerIP     string `yaml:"serverip"`
	DefaultScope string `yaml:"scope"`
}
type initOptionsDT struct {
	folder string
}

var initOptions initOptionsDT

var initCommand = cmd.NewCmd("init").WithFlags(
	&cmd.Flag{
		Name:      "file",
		Shorthand: "f",
		DefValue:  "",
		Usage:     "set the value for storing the configuration file",
		Value:     &initOptions.folder,
	},
).NoArgs(
	func(c context.Context) error {
		config := insprConfiguration{}
		fmt.Print("enter insprd host (http://localhost:8080):")
		fmt.Scanln(&config.ServerIP)
		if config.ServerIP == "" {
			config.ServerIP = "http://localhost:8080"
		}
		fmt.Print("enter default scope (\"\"):")
		fmt.Scanln(&config.DefaultScope)
		var output *os.File
		file := initOptions.folder

		if file == "" {

			defaultFolder, _ := os.UserHomeDir()
			defaultFolder = filepath.Join(defaultFolder, ".inspr")
			if f, err := os.Stat(defaultFolder); err != nil || !f.IsDir() {
				if os.IsNotExist(err) {
					os.Mkdir(defaultFolder, os.ModePerm)
				} else if err != nil {
					return ierrors.NewError().Message(err.Error()).Build()
				} else {
					return errors.New("default folder already defined as a file. did you name something .inspr in your home folder?")
				}
			}

			output, _ = os.OpenFile(filepath.Join(defaultFolder, "config"), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
		} else {
			var err error
			output, err = os.OpenFile(file, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return err
			}
		}

		encoder := yaml.NewEncoder(output)
		return encoder.Encode(config)
	},
)
