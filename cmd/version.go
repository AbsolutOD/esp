package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Version of esp",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runVersion()
		},
	}
}

func versionString() string {
	return fmt.Sprintf("esp %s (commit %s, built %s)", version, commit, date)
}

func runVersion() error {
	fmt.Println(versionString())
	return nil
}
