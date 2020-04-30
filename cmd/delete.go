package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/pinpt/esp/internal/client"
	"github.com/pinpt/esp/internal/common"
	"github.com/spf13/cobra"
)

// deleteCmd gets the parameter from the backend store
var deleteCmd = &cobra.Command{
	Use:   "delete [path]",
	Aliases: []string{"rm"},
	Short: "Delete a parameter by path in SSM",
	Long:  `Allows you to delete a specific ssm parameter with an exact path.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ec := client.New(client.EspClient{ Backend: "ssm" })
		param := ec.Delete(common.DeleteInput{
			Name:    args[0],
		})
		name := aurora.BrightYellow(param)
		fmt.Printf("Deleted: %s\n", name)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
