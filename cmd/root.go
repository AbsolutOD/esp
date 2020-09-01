package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var IsEspProject = false
var Env string

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

	// check AWS Region
	if _, ok := os.LookupEnv("AWS_DEFAULT_REGION"); !ok {
		fmt.Println("Please set the AWS_DEFAULT_REGION environment variable.")
		os.Exit(1)
	}
	if _, ok := os.LookupEnv("AWS_PROFILE"); !ok {
		fmt.Println("Please set the AWS_PROFILE environment variable.")
		os.Exit(2)
	}
	rootCmd.PersistentFlags().StringVarP(&Env, "env", "e", "", "Declare the env to work on.")
	if IsEspProject {
		if err := rootCmd.MarkFlagRequired("env"); err != nil {
			fmt.Println("There is an .espFile.yaml defined, so you need to set --env arg.")
			os.Exit(3)
		}
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Get current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search config in home directory with name ".esp" (without extension).
	viper.AddConfigPath(cwd)
	viper.SetConfigName(".espFile.yaml")
	//viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in and mark that this is an ESP project.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Error reading in config file: ", viper.ConfigFileUsed())
		IsEspProject = true
	}
}
