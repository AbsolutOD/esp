package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of esp",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		fmt.Println("ESP version 0.2.0")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
