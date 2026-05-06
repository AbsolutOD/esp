package cmd

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora/v4"
	"github.com/pinpt/esp/internal/common"
	"github.com/spf13/cobra"
)

// moveCmd gets the parameter from the backend store
var moveCmd = &cobra.Command{
	Use:     "move [path]",
	Aliases: []string{"mv"},
	Short:   "move a parameter by path in SSM",
	Long:    `Allows you to move a specific ssm parameter with an exact path.`,
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		p, err := c.Move(common.MoveCommand{
			Source:      args[0],
			Destination: args[1],
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		src := aurora.BrightYellow(p.Source)
		dest := aurora.BrightYellow(p.Destination)
		fmt.Printf("%s => %s\n", src, dest)
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)
}
