package cmd

import (
	"fmt"

	"github.com/AbsolutOD/esp/internal/client"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
)

func newDeleteCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:     "delete [path]",
		Aliases: []string{"rm"},
		Short:   "Delete a parameter by path in SSM",
		Long:    `Allows you to delete a specific ssm parameter with an exact path.`,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runDelete(args, a.Client)
		},
	}
}

func runDelete(args []string, c *client.EspClient) error {
	name, err := c.Delete(common.DeleteInput{Name: args[0]})
	if err != nil {
		return err
	}
	fmt.Printf("Deleted: %s\n", aurora.BrightYellow(name))
	return nil
}
