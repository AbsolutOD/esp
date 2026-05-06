package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/pinpt/esp/internal/app"
	"github.com/pinpt/esp/internal/client"
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
		os.Exit(1)
	}
}

func init() {
	esp = app.New(false)
	// setting this to SSM just to make the interface nicer, since we only have the SSM backend
	esp.Backend = "ssm"
	cobra.OnInitialize(configureLogging)
	rootCmd.PersistentPreRunE = persistentPreRun

	// CLI args
	rootCmd.PersistentFlags().StringVarP(&esp.Env, "env", "e", "", "Declare the env to work on.")
	rootCmd.PersistentFlags().StringVarP(&esp.Backend, "backend", "b", "ssm", "Set which backend to use.")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Show more output")
}

// configureLogging sets the slog default handler. No I/O, no failure
// path. Runs via cobra.OnInitialize, which is skipped for --help.
func configureLogging() {
	level := slog.LevelWarn
	if verbose {
		level = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
}

// persistentPreRun validates AWS env vars, constructs the backend
// client, reads the .espFile, and (when one is present) marks --env
// required. Cobra short-circuits before PersistentPreRunE for --help,
// so help output still works without AWS credentials.
func persistentPreRun(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	if _, ok := os.LookupEnv("AWS_DEFAULT_REGION"); !ok {
		return fmt.Errorf("AWS_DEFAULT_REGION environment variable is not set")
	}
	if _, ok := os.LookupEnv("AWS_PROFILE"); !ok {
		return fmt.Errorf("AWS_PROFILE environment variable is not set")
	}

	var err error
	c, err = client.New(esp)
	if err != nil {
		return err
	}

	viper.SetConfigName(esp.Filename)
	viper.AddConfigPath(esp.Path)

	// If a config file is found, read it in and mark that this is an ESP project.
	if err := viper.ReadInConfig(); err == nil {
		esp.IsEspProject = true
	}

	if esp.IsEspProject {
		if err := viper.Unmarshal(&esp); err != nil {
			return fmt.Errorf("parsing %s: %w", esp.Filename, err)
		}
		if err := rootCmd.MarkFlagRequired("env"); err != nil {
			return fmt.Errorf("marking --env required: %w", err)
		}
	}
	return nil
}
