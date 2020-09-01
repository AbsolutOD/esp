package cmd

import (
	"errors"
	"github.com/pinpt/esp/internal/app"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
func initCmd() *cobra.Command {
	a := app.New(IsEspProject)
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initializes the current directory to be an ESP based application.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if IsEspProject {
				return errors.New(".espFile already exists\n")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			a.InitQuestions()
		},
	}

	return  initCmd
}

func init() {
	initCmd := initCmd()
	rootCmd.AddCommand(initCmd)
}
