package cli

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configScope    = "scope"
	configServerIP = "serverIP"
)

var defaultValues map[string]string = map[string]string{
	configScope:    "",
	configServerIP: "http://127.0.0.1:8080",
}

// initViperConfig - sets defaults values and where is the file in which new values can be read
func initViperConfig() {
	// specifies the path in which the config file present
	viper.AddConfigPath("$HOME/.inspr/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	for k, v := range defaultValues {
		viper.SetDefault(k, v)
	}
}

// createViperConfig - creates the folder and or file of the inspr's viper config
//
// if they already a file the createConfig will truncate it before writing
func createViperConfig(configPath string) error {
	// creates config file
	err := viper.WriteConfigAs(configPath)
	if err != nil {
		return err
	}
	return nil
}

// createInsprConfigFolder - creates the folder of the inspr's config, it only
// creates the folder if already doesn't exists
func createInsprConfigFolder(path string) error {
	// creates folder
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0777); err != nil { // perm 0666
			return err
		}
	}

	return nil
}

// readConfig - reads the inspr's viper config, in case it didn't
// found any, it creates one with the defaults values
func readViperConfig(baseDir string) error {
	folderDir := filepath.Join(baseDir, ".inspr")
	configDir := filepath.Join(folderDir, "config")

	if _, err := os.Stat(folderDir); os.IsNotExist(err) {
		if createErr := createInsprConfigFolder(folderDir); createErr != nil {
			return createErr
		}
	}

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if configErr := createViperConfig(configDir); configErr != nil {
			return configErr
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

// changeViperValues - changes the values of the viper configuration
// and saves it in the config file of inspr, if the file is not created
// it will return an error.
func changeViperValues(key string, value interface{}) error {
	viper.Set(key, value)
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func existingKeys() []string {
	return viper.GetViper().AllKeys()
}
