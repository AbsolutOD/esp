package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/AbsolutOD/esp/internal/app"
	"github.com/AbsolutOD/esp/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// App holds the runtime dependencies threaded into every subcommand.
// Config is populated at construction; Client is populated by the root
// command's PersistentPreRunE after env-var validation succeeds.
type App struct {
	Config  *app.Config
	Client  *client.EspClient
	Verbose bool
}

// Execute is the entry point invoked by main. It builds the App,
// constructs the cobra tree around it, and runs.
func Execute() {
	a := &App{Config: app.New(false)}
	a.Config.Backend = "ssm"
	root := newRootCmd(a)
	root.AddCommand(
		newGetCmd(a),
		newPutCmd(a),
		newListCmd(a),
		newCopyCmd(a),
		newMoveCmd(a),
		newDeleteCmd(a),
		newInitCmd(a),
		newVersionCmd(),
	)
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

// newRootCmd builds the root cobra command bound to the given App.
func newRootCmd(a *App) *cobra.Command {
	root := &cobra.Command{
		Use:   "esp",
		Short: "A utility to browse and export SSM Parameter values into different formats.",
	}
	cobra.OnInitialize(func() { configureLogging(a.Verbose) })
	root.PersistentPreRunE = persistentPreRunE(a)
	root.PersistentFlags().StringVarP(&a.Config.Env, "env", "e", "", "Declare the env to work on.")
	root.PersistentFlags().StringVarP(&a.Config.Backend, "backend", "b", "ssm", "Set which backend to use.")
	root.PersistentFlags().BoolVar(&a.Verbose, "verbose", false, "Show more output")
	return root
}

// configureLogging sets the slog default handler. No I/O, no failure
// path. Runs via cobra.OnInitialize, which is skipped for --help.
func configureLogging(verbose bool) {
	level := slog.LevelWarn
	if verbose {
		level = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
}

// persistentPreRunE returns the closure that validates AWS env vars,
// constructs the backend client, and reads the .espFile. Cobra
// short-circuits before PersistentPreRunE for --help, so help still
// works without AWS credentials.
func persistentPreRunE(a *App) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		if _, ok := os.LookupEnv("AWS_DEFAULT_REGION"); !ok {
			return fmt.Errorf("AWS_DEFAULT_REGION environment variable is not set")
		}
		if _, ok := os.LookupEnv("AWS_PROFILE"); !ok {
			return fmt.Errorf("AWS_PROFILE environment variable is not set")
		}

		c, err := client.New(a.Config)
		if err != nil {
			return err
		}
		a.Client = c

		viper.SetConfigName(a.Config.Filename)
		viper.AddConfigPath(a.Config.Path)

		if err := viper.ReadInConfig(); err == nil {
			a.Config.IsEspProject = true
		}

		if a.Config.IsEspProject {
			if err := viper.Unmarshal(a.Config); err != nil {
				return fmt.Errorf("parsing %s: %w", a.Config.Filename, err)
			}
			if err := cmd.Root().MarkFlagRequired("env"); err != nil {
				return fmt.Errorf("marking --env required: %w", err)
			}
		}
		return nil
	}
}
