package cmd

import "github.com/spf13/cobra"

func newInitCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initializes the current directory to be an ESP based application.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			a.Config.InitQuestions()
			return nil
		},
	}
}
