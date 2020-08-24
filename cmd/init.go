package cmd

import (
	"errors"
	"fmt"
	"github.com/pinpt/esp/internal/app"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
func initCmd() *cobra.Command {
	a := app.New()
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initializes the current directory to be an ESP based application.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if app.CheckForConfigFile() {
				fmt.Print(".espFile already exists")
				return errors.New(".espFile already exists")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			a.Init()
		},
	}

	return  initCmd
}

func init() {
	initCmd := initCmd()
	rootCmd.AddCommand(initCmd)
}
