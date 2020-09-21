package cmd

import (
	"fmt"
	"github.com/pinpt/esp/internal/common"

	"github.com/logrusorgru/aurora"
	"github.com/pinpt/esp/internal/client"

	"github.com/spf13/cobra"
)

func displayParams(p []common.EspParam) {
	for _, param := range p {
		name := aurora.BrightYellow(param.Name)
		fmt.Printf("%s: %s\n", name, param.Value)
	}
}

// listCmd represents the list command
func listCmd() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:     "list [path]",
		Aliases: []string{"ls"},
		Short:   "Recursively list a SSM path.",
		Long:    `The list command gives you an easy way to recursively get all SSM parameters with a base path.`,
		//Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ec := client.New(client.EspClient{Backend: "ssm"})
			decrypt, _ := cmd.Flags().GetBool("decrypt")
			params := ec.ListParams(common.ListParamInput{
				Path:      args[0],
				Decrypt:   decrypt,
				Recursive: true,
			})
			displayParams(params)
		},
	}
	return listCmd
}

func init() {
	listCmd := listCmd()
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolP("decrypt", "d", false, "Decrypt SSM secure strings.")
	listCmd.Flags().BoolP("path", "p", false, "Path to list parameters.")
	//listCmd.Flags().BoolP("recursive", "r", false, "Recursively get params from sub dirs.")
}
