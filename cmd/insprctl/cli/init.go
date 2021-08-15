package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
	"inspr.dev/inspr/pkg/cmd"
	"inspr.dev/inspr/pkg/ierrors"
)

type insprConfiguration struct {
	ServerIP     string `yaml:"serverip"`
	DefaultScope string `yaml:"scope"`
	ServerHost   string `yaml:"host"`
}
type initOptionsDT struct {
	folder string
}

var initOptions initOptionsDT

var initCommand = cmd.NewCmd("init").
	WithDescription("Initialize the CLI configuration").
	WithFlags(
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
		fmt.Print("enter insprd IP or URL (localhost:8080):")
		fmt.Scanln(&config.ServerIP)
		if !strings.HasPrefix(config.ServerIP, "http") {
			config.ServerIP = fmt.Sprintf("http://%s", config.ServerIP)
		}
		if config.ServerIP == "" {
			config.ServerIP = "http://localhost:8080"
		}
		fmt.Print("Opitional config: insprd host (example.inspr.dev):")
		fmt.Scanln(&config.ServerHost)
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
					return ierrors.Wrap(
						err,
						fmt.Sprintf("error processing the %v", defaultFolder),
					)
				} else {
					return ierrors.New("default folder already defined as a file. did you name something .inspr in your home folder?")
				}
			}

			output, _ = os.OpenFile(
				filepath.Join(defaultFolder, "config"),
				os.O_TRUNC|os.O_WRONLY|os.O_CREATE,
				0644,
			)
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
