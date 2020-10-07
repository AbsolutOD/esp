package cmd

import (
	"fmt"
	"github.com/pinpt/esp/internal/app"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var App = app.New(false)
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
	rootCmd.PersistentFlags().StringVarP(&App.Env, "env", "e", "", "Declare the env to work on.")
	if App.IsEspProject {
		if err := rootCmd.MarkFlagRequired("env"); err != nil {
			fmt.Println("There is an .espFile.yaml defined, so you need to set --env arg.")
			os.Exit(3)
		}
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Get current working directory.
	/*cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}*/

	// Search config in home directory with name ".esp" (without extension).
	//fmt.Printf("CWD: %s\n", cwd)
	viper.SetConfigName(".espFile.yaml")
	viper.AddConfigPath(".")

	// If a config file is found, read it in and mark that this is an ESP project.
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Errorf("unable to load espFile: %s", err)
	} else {
		fmt.Println("Loaded the espFie")
		App.IsEspProject = true
	}
	if err := viper.Unmarshal(&App); err != nil {
		fmt.Println("Error parsing the .espFile.yaml")
	}
	fmt.Println("End of initializing.")
}
