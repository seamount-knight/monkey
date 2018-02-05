package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"fmt"
	"os"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use: "monkey",
	Short: "Welcome to monkey",
	Long: "Welcome to monkey",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $monkey/config.yaml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("$monkey/")
	viper.AddConfigPath(".")

	viper.SetDefault("COMPONENT", "monkey")

	// database related config variables
	viper.SetDefault("DB_ENGINE", "postgres")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_NAME", "monkeydb")
	viper.SetDefault("DB_USER", "monkey")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_CONN_TIMEOUT", 0)
	viper.SetDefault("DB_MAX_CONN", 10)
	viper.SetDefault("DB_MAX_IDLE_CONN", 1)


	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("PORT", "80")

	viper.SetDefault("GIN_MODE", "debug")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
