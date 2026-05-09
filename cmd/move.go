package cmd

import (
	"fmt"

	"github.com/AbsolutOD/esp/internal/client"
	"github.com/AbsolutOD/esp/internal/common"
	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
)

func newMoveCmd(a *App) *cobra.Command {
	return &cobra.Command{
		Use:     "move [path]",
		Aliases: []string{"mv"},
		Short:   "move a parameter by path in SSM",
		Long:    `Allows you to move a specific ssm parameter with an exact path.`,
		Args:    cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			return runMove(args, a.Client)
		},
	}
}

func runMove(args []string, c *client.EspClient) error {
	p, err := c.Move(common.MoveCommand{Source: args[0], Destination: args[1]})
	if err != nil {
		return err
	}
	fmt.Printf("%s => %s\n", aurora.BrightYellow(p.Source), aurora.BrightYellow(p.Destination))
	return nil
}
