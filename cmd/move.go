package cmd

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/pinpt/esp/internal/client"
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
		ec := client.New(client.EspClient{Backend: "ssm"})
		p := ec.Move(common.MoveCommand{
			Source:      args[0],
			Destination: args[1],
		})
		src := aurora.BrightYellow(p.Source)
		dest := aurora.BrightYellow(p.Destination)
		fmt.Printf("%s => %s\n", src, dest)
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)
}
