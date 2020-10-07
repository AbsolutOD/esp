package cmd

import (
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
func initCmd() *cobra.Command {
	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initializes the current directory to be an ESP based application.",
		Run: func(cmd *cobra.Command, args []string) {
			App.InitQuestions()
		},
	}

	return  initCmd
}

func init() {
	initCmd := initCmd()
	rootCmd.AddCommand(initCmd)
}
