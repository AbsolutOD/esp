package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pinpt/esp/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var awsProfile string
var awsRegion string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "esp",
	Short: "A utility to browse and export SSM Parameter values into different formats.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.esp.yaml)")

	// set AWS Profile
	profile := utils.GetEnv("AWS_PROFILE", "default")
	region := utils.GetEnv("AWS_REGION", "us-east-1")
	rootCmd.PersistentFlags().StringVar(&awsProfile, "profile", profile, "Set the AWS profile to use.")
	rootCmd.PersistentFlags().StringVar(&awsRegion, "region", region, "Set the AWS region to use.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".esp" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".esp")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
