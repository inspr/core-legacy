package utils

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"inspr.dev/inspr/pkg/cmd"
	"inspr.dev/inspr/pkg/meta/utils"
)

// ignores unused code for this file in the staticcheck
//lint:file-ignore U1000 Ignore all unused code

const (
	configScope    = "scope"
	configServerIP = "serverip"
	configHost     = "host"
)

var defaultValues map[string]string = map[string]string{
	configScope:    "",
	configServerIP: "http://<cluster_ip>",
	configHost:     "",
}

var flagCompletionRegistry = map[string]func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective){
	"scope": func(cm *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		toComplete = strings.TrimSuffix(toComplete, ".")
		client := GetCliClient()
		scope, err := GetScope()

		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		newScope, err := utils.JoinScopes(scope, toComplete)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		if _, err := client.Apps().Get(context.Background(), newScope); err != nil {
			newScope, _, _ = utils.RemoveLastPartInScope(newScope)
		}

		app, err := client.Apps().Get(context.Background(), newScope)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		scopes := []string{}
		for name := range app.Spec.Apps {
			newScope, _ := utils.JoinScopes(newScope, name)
			if strings.HasPrefix(newScope, toComplete) {
				scopes = append(scopes, newScope+".")
			}
		}
		return scopes, cobra.ShellCompDirectiveNoSpace
	},
}

// AddDefaultFlagCompletion adds the default completion for most used flags
func AddDefaultFlagCompletion() cmd.Option {
	return func(c *cobra.Command) {
		for name, f := range flagCompletionRegistry {
			c.RegisterFlagCompletionFunc(name, f)
		}
		c.MarkFlagFilename("token")
		c.MarkFlagFilename("config")
		c.MarkFlagFilename("file")
		c.MarkFlagDirname("folder")
	}
}

//ServetIpKey returns the key value of ConfigServerIP
//Avoids having to constants public
func ServerIpKey() string {
	return configServerIP
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

//GetConfiguredHost is responsible for returning config value for host.
//Avoids having to constants public.
func GetConfiguredHost() string {
	return viper.GetString(configHost)
}

//InitViperConfig - sets defaults values and where is the file in which new values can be read
func InitViperConfig() {
	viper.SetConfigType("yaml")
	for k, v := range defaultValues {
		viper.SetDefault(k, v)
	}
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
		SetGlobalClient()
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
