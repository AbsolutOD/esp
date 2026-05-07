package cmd

import (
	"errors"

	"github.com/AbsolutOD/esp/internal/common"
	"github.com/spf13/cobra"
)

// cpCmd represents the cp command
var copyCmd = &cobra.Command{
	Use:     "copy [OPTIONS] SRC_SSM_PATH DEST_SSM_PATH",
	Aliases: []string{"cp"},
	Short:   "Copy a SSM Param from its current path to a new SSM Path",
	Long:    "Copy SSM value from an existing path to a new path.\n",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		if args[0] == "" {
			return errors.New("source can not be empty")
		}
		if args[1] == "" {
			return errors.New("destination can not be empty")
		}
		cc := common.CopyCommand{
			Source:      args[0],
			Destination: args[1],
		}
		if _, err := c.Copy(cc); err != nil {
			return err
		}
		return nil
	},
	Example: "esp cp /ssm/path/key /ssm/new/path/key",
}

func init() {
	rootCmd.AddCommand(copyCmd)

	// cpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
