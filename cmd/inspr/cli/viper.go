package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	configCurrentScope = "currentScope"
	configServerIP     = "serverIP"
)

var defaultValues map[string]string = map[string]string{
	configCurrentScope: "./app1/app2",
	configServerIP:     "127.0.0.1",
}

// initConfig - sets defaults values and where is the file in which new values can be read
func initConfig() {
	// specifies the path in which the config file present
	viper.AddConfigPath("$HOME/.inspr/")
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")

	for k, v := range defaultValues {
		viper.SetDefault(k, v)
	}
}

// createConfig - creates the folder and or file of the inspr's viper config
// if they already a file the createConfig will truncate it before writing
func createConfig() error {
	homeDir := os.Getenv("HOME")
	insprDir := homeDir + "/" + ".inspr"

	if _, err := os.Stat(insprDir); os.IsNotExist(err) {
		os.Mkdir(insprDir, os.ModePerm)
	}

	bytes, err := yaml.Marshal(defaultValues)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(insprDir+"/env.yaml", bytes, 0777)
	if err != nil {
		return err
	}

	return nil
}

// readConfig - reads the inspr's viper config, in case it didn't
// found any, it creates one with the defaults values
func readConfig(out io.Writer) error {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if createErr := createConfig(); createErr != nil {
				err = createErr
				fmt.Fprintln(out, err.Error())
			} else {
				err = nil
			}
		} else {
			fmt.Fprintln(out, err.Error())
		}
		return err
	}
	return nil
}
