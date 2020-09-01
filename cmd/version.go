package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of esp",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ESP version 0.1.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
