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
	return viper.GetString(configServerIP)
}

//InitViperConfig - sets defaults values and where is the file in which new values can be read
func InitViperConfig() {
	// specifies the path in which the config file present
	viper.AddConfigPath("$HOME/.inspr/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	for k, v := range defaultValues {
		viper.SetDefault(k, v)
	}
}

// createConfig - creates the folder and or file of the inspr's viper config
// if they already a file the createConfig will truncate it before writing
func createViperConfig() error {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// folder path
	insprFolderDir := filepath.Join(homeDir, ".inspr")

	// creates folder
	if _, err := os.Stat(insprFolderDir); os.IsNotExist(err) {
		if err := os.Mkdir(insprFolderDir, 0777); err != nil { // perm 0666
			return err
		}
	}

	// file path
	fileDir := filepath.Join(insprFolderDir, "config")

	// creates config file
	err = viper.WriteConfigAs(fileDir)
	if err != nil {
		return err
	}
	return nil
}

// ReadViperConfig - reads the inspr's viper config, in case it didn't
// found any, it creates one with the defaults values
func ReadViperConfig() error {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".inspr", "config")

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if createErr := createViperConfig(); createErr != nil {
			return createErr
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

// ChangeViperValues - changes the values of the viper configuration
// and saves it in the config file of inspr, if the file is not created
// it will return an error.
func ChangeViperValues(key string, value interface{}) error {
	viper.Set(key, value)
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

// ExistsKey - informs to the user if the key passed exists in the
// default keys that are saved in the inspr config file
func ExistsKey(key string) bool {
	for k := range defaultValues {
		if k == key {
			return true
		}
	}
	return false
}

// ExistingKeys - returns to the user all availible keys in viper's configs.
func ExistingKeys() []string {
	arr := make([]string, 0)
	for k := range defaultValues {
		arr = append(arr, k)
	}
	return arr
}
