package cmd

import (
	"fmt"

	"github.com/logrusorgru/aurora/v4"
	"github.com/pinpt/esp/internal/common"
	"github.com/spf13/cobra"
)

// deleteCmd gets the parameter from the backend store
var deleteCmd = &cobra.Command{
	Use:     "delete [path]",
	Aliases: []string{"rm"},
	Short:   "Delete a parameter by path in SSM",
	Long:    `Allows you to delete a specific ssm parameter with an exact path.`,
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		param, err := c.Delete(common.DeleteInput{
			Name: args[0],
		})
		if err != nil {
			return err
		}
		name := aurora.BrightYellow(param)
		fmt.Printf("Deleted: %s\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
