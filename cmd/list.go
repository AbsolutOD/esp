package cmd

import (
	"fmt"
	"github.com/pinpt/esp/internal/common"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func displayParams(p []common.EspParam) {
	for _, param := range p {
		name := aurora.BrightYellow(param.Name)
		fmt.Printf("%s: %s\n", name, param.Value)
	}
}

func getPath(a []string) string {
	if len(a) == 0 {
		return esp.GetAppPath()
	}
	return a[0]
}

// listCmd represents the list command
func listCmd() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:     "list [path]",
		Aliases: []string{"ls"},
		Short:   "Recursively list a SSM path if given.",
		Long:    `The list command gives you an easy way to recursively get all SSM parameters with a base path.
If you have a .espFile.yaml in the current directory this command will list all params under the project path.`,
		//Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			decrypt, _ := cmd.Flags().GetBool("decrypt")
			params := c.ListParams(common.ListParamInput{
				Path:      getPath(args),
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
