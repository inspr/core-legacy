package utils

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

//GetConfiguredServerIP is responsible for returning config value for serverIp.
//Avoids having to constants public.
func GetConfiguredServerIP() string {
	return viper.GetString(configServerIP)
}

//GetConfiguredScope is responsible for returning config value for scope.
//Avoids having to constants public.
func GetConfiguredScope() string {
	return viper.GetString(configScope)
}

//InitViperConfig - sets defaults values and where is the file in which new values can be read
func InitViperConfig() {
	viper.SetConfigType("yaml")
	for k, v := range defaultValues {
		viper.SetDefault(k, v)
	}
}

// createViperConfig - creates the folder and or file of the inspr's viper config
//
// if they already a file the createConfig will truncate it before writing
func createViperConfig(path string) error {
	// creates config file
	err := viper.WriteConfigAs(path)
	if err != nil {
		return err
	}
	return nil
}

// createInsprConfigFolder - creates the folder of the inspr's config, it only
// creates the folder if already doesn't exists
func createInsprConfigFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0777); err != nil { // perm 0666
			return err
		}
	}

	return nil
}

// ConfigFile is the currently loaded config file
var ConfigFile string

// ReadDefaultConfig reads the default insprctl configuration
func ReadDefaultConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	return ReadViperConfig(home)
}

// ReadConfigFromFile reads a config from a file
func ReadConfigFromFile(file string) error {
	ConfigFile = file
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	return viper.ReadConfig(f)
}

// ReadViperConfig - reads the inspr's viper config, in case it didn't
// found any, it creates one with the defaults values
func ReadViperConfig(basePath string) error {
	folderPath := filepath.Join(basePath, ".inspr")
	filePath := filepath.Join(folderPath, "config")

	if err := ReadConfigFromFile(filePath); err != nil {
		return err
	}
	return nil
}

// ChangeViperValues - changes the values of the viper configuration
// and saves it in the config file of inspr, if the file is not created
// it will return an error.
func ChangeViperValues(key string, value interface{}) error {
	viper.Set(key, value)
	if err := viper.WriteConfigAs(ConfigFile); err != nil {
		return err
	}
	if key == configServerIP {
		setGlobalClient()
	}

	return nil
}

// ExistsKey - informs to the user if the key passed exists in the
// default keys that are saved in the insprctl config file
func ExistsKey(key string) bool {
	return viper.IsSet(key)
}

// ExistingKeys - returns to the user all available keys in viper's configs.
func ExistingKeys() []string {
	return viper.GetViper().AllKeys()
}
