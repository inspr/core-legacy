package cli

import (
	"github.com/spf13/viper"
)

func initViper() {
	viper.SetConfigName("inspr_config")
	// viper.SetConfigFile("yaml")

	viper.SetDefault("port", "8080")

	// searches for the config in these folders
	viper.AddConfigPath("$HOME")
	// viper.AddConfigPath(".")

	viper.AutomaticEnv()
}
