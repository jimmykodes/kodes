package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "kodes",
	Short: "A collection of CLI tools",
}

func init() {
	cobra.OnInitialize(initConfig)
}

func Execute() error {
	return rootCmd.Execute()
}

func initConfig() {
	configFile := os.Getenv("KODES_CONFIG_FILE")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		homeDir, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(homeDir)
		viper.SetConfigName(".kodes")
		viper.SetConfigType("yaml")
	}
	err := viper.ReadInConfig()
	if err != nil && configFile != "" {
		// only printing an error if the user supplied a non-standard config file.
		// Otherwise, just assume one doesn't exist and likely is not needed.
		fmt.Fprintf(os.Stderr, "Could not find config file: %s", configFile)
	}
}
