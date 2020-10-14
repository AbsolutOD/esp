package cmd

import (
	"fmt"
	"github.com/pinpt/esp/internal/app"
	"github.com/pinpt/esp/internal/client"
	jww "github.com/spf13/jwalterweatherman"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var verbose bool
var esp *app.Config
var c *client.EspClient

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
	esp = app.New(false)
	// setting this to SSM just to make the interface nicer, since we only have the SSM backend
	esp.Backend = "ssm"
	c = client.New(esp)
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
	rootCmd.PersistentFlags().StringVarP(&esp.Env, "env", "e", "", "Declare the env to work on.")
	rootCmd.PersistentFlags().StringVarP(&esp.Backend, "backend", "b", "ssm", "Set which backend to use.")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Show more output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Just setting for debugging

	if verbose {
		jww.SetStdoutThreshold(jww.LevelInfo)
	}
	viper.SetConfigName(esp.Filename)
	viper.AddConfigPath(esp.Path)

	// If a config file is found, read it in and mark that this is an ESP project.
	if err := viper.ReadInConfig(); err == nil {
		esp.IsEspProject = true
	}

	if esp.IsEspProject {
		if err := viper.Unmarshal(&esp); err != nil {
			fmt.Printf("Error parsing the %s\n", esp.Filename)
		}
		// not going to force this at the moment.  I will add this when I have a second backend
		/*if err := rootCmd.MarkFlagRequired("backend"); err != nil {
			//fmt.Printf("There is an %s.yaml defined, so you need to set --env arg.\n", esp.Filename)
			os.Exit(3)
		}*/
		if err := rootCmd.MarkFlagRequired("env"); err != nil {
			//fmt.Printf("There is an %s.yaml defined, so you need to set --env arg.\n", esp.Filename)
			os.Exit(3)
		}
	}
}
