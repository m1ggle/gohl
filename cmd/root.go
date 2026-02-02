package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "gohl",
	Short: "gohl — a small MySQL CLI",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Define persistent flags here so they are available to both login and query commands
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gohl.json)")
	rootCmd.PersistentFlags().StringP("host", "H", "127.0.0.1", "MySQL host")
	rootCmd.PersistentFlags().StringP("port", "P", "3306", "MySQL port")
	rootCmd.PersistentFlags().StringP("user", "u", "root", "MySQL user")
	rootCmd.PersistentFlags().StringP("pass", "p", "", "MySQL password (if empty will prompt)")
	rootCmd.PersistentFlags().StringP("db", "d", "", "Database name")

	// Bind flags to viper
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("pass", rootCmd.PersistentFlags().Lookup("pass"))
	viper.BindPFlag("db", rootCmd.PersistentFlags().Lookup("db"))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory and current directory with name "gohl" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.AddConfigPath("../conf")
		viper.SetConfigType("yaml")
		viper.SetConfigName("ghl")

		// Also support .gohl.json in home
		viper.SetConfigName(".gohl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
