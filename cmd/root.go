package cmd

import (
	"fmt"
	"github.com/pinpt/esp/internal/app"
	jww "github.com/spf13/jwalterweatherman"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var verbose bool
var cfg *app.Config

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
	cfg = app.New(false)
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

	// CLI args
	rootCmd.PersistentFlags().StringVarP(&cfg.Env, "env", "e", "", "Declare the env to work on.")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Show more output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Just setting for debugging

	if verbose {
		jww.SetStdoutThreshold(jww.LevelInfo)
	}
	viper.SetConfigName(cfg.Filename)
	viper.AddConfigPath(cfg.Path)

	// If a config file is found, read it in and mark that this is an ESP project.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("unable to load espFile: %s\n", err)
	} else {
		cfg.IsEspProject = true
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("Error parsing the %s\n", cfg.Filename)
	}

	if cfg.IsEspProject {
		if err := rootCmd.MarkFlagRequired("env"); err != nil {
			//fmt.Printf("There is an %s.yaml defined, so you need to set --env arg.\n", cfg.Filename)
			os.Exit(3)
		}
	}
}
