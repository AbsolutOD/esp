package cmd

import (
	"errors"

	"github.com/AbsolutOD/esp/internal/client"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/spf13/cobra"
)

func newCopyCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:     "copy [OPTIONS] SRC_SSM_PATH DEST_SSM_PATH",
		Aliases: []string{"cp"},
		Short:   "Copy a SSM Param from its current path to a new SSM Path",
		Long:    "Copy SSM value from an existing path to a new path.\n",
		Args:    cobra.ExactArgs(2),
		Example: "esp cp /ssm/path/key /ssm/new/path/key",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runCopy(args, a.Client)
		},
	}
}

func runCopy(args []string, c *client.EspClient) error {
	if args[0] == "" {
		return errors.New("source can not be empty")
	}
	if args[1] == "" {
		return errors.New("destination can not be empty")
	}
	_, err := c.Copy(common.CopyCommand{
		Source:      args[0],
		Destination: args[1],
	})
	return err
}
